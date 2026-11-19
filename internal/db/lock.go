package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Locker takes a Postgres advisory lock on a job. The lock belongs to the
// connection, so a crashed executor releases it automatically and another
// executor can resume the job.
type Locker struct {
	pool *pgxpool.Pool
}

// NewLocker builds a locker over an existing pool.
func NewLocker(pool *pgxpool.Pool) *Locker {
	return &Locker{pool: pool}
}

// TryLock reports whether the caller now owns the job. When it returns ok, the
// release function must be called to free the lock and its connection.
func (l *Locker) TryLock(ctx context.Context, jobID int64) (release func(), ok bool, err error) {
	conn, err := l.pool.Acquire(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("db: acquire lock connection: %w", err)
	}

	var locked bool
	if err := conn.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, jobID).Scan(&locked); err != nil {
		conn.Release()
		return nil, false, fmt.Errorf("db: try job lock %d: %w", jobID, err)
	}
	if !locked {
		conn.Release()
		return nil, false, nil
	}

	return func() {
		// The lock is released when the connection closes anyway; unlocking
		// first keeps the connection reusable without a round trip to the pool.
		_, _ = conn.Exec(context.Background(), `SELECT pg_advisory_unlock($1)`, jobID)
		conn.Release()
	}, true, nil
}
