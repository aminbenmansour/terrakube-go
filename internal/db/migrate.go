package db

import (
	"context"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// migrationLockKey is an arbitrary application-wide advisory lock that stops
// two API replicas from running migrations at the same time.
const migrationLockKey int64 = 0x74657272616b7562 // "terrakub"

// Migrate applies every *.sql file in dir that has not been applied yet.
// Files are applied in lexical order inside individual transactions and
// recorded in schema_migrations.
func Migrate(ctx context.Context, pool *pgxpool.Pool, fsys fs.FS, dir string) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("db: acquire migration connection: %w", err)
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, `SELECT pg_advisory_lock($1)`, migrationLockKey); err != nil {
		return fmt.Errorf("db: acquire migration lock: %w", err)
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = conn.Exec(unlockCtx, `SELECT pg_advisory_unlock($1)`, migrationLockKey)
	}()

	if _, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    text PRIMARY KEY,
			applied_at timestamptz NOT NULL DEFAULT now()
		)`); err != nil {
		return fmt.Errorf("db: create schema_migrations: %w", err)
	}

	applied, err := appliedVersions(ctx, conn)
	if err != nil {
		return err
	}

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return fmt.Errorf("db: read migrations: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)

	for _, name := range files {
		if applied[name] {
			continue
		}
		if err := applyMigration(ctx, conn, fsys, dir, name); err != nil {
			return err
		}
	}

	return nil
}

func appliedVersions(ctx context.Context, conn *pgxpool.Conn) (map[string]bool, error) {
	rows, err := conn.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("db: list applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("db: scan applied migration: %w", err)
		}
		applied[version] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: iterate applied migrations: %w", err)
	}
	return applied, nil
}

func applyMigration(ctx context.Context, conn *pgxpool.Conn, fsys fs.FS, dir, name string) error {
	body, err := fs.ReadFile(fsys, path(dir, name))
	if err != nil {
		return fmt.Errorf("db: read migration %s: %w", name, err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db: begin migration %s: %w", name, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, string(body)); err != nil {
		return fmt.Errorf("db: apply migration %s: %w", name, err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, name); err != nil {
		return fmt.Errorf("db: record migration %s: %w", name, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("db: commit migration %s: %w", name, err)
	}
	return nil
}

func path(dir, name string) string {
	if dir == "" || dir == "." {
		return name
	}
	return dir + "/" + name
}
