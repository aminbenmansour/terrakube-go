package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TERRAKUBE_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TERRAKUBE_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := NewPool(ctx, url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func TestLockerExcludesASecondHolder(t *testing.T) {
	pool := testPool(t)
	locker := NewLocker(pool)
	ctx := context.Background()
	jobID := time.Now().UnixNano()

	release, ok, err := locker.TryLock(ctx, jobID)
	if err != nil {
		t.Fatalf("first lock: %v", err)
	}
	if !ok {
		t.Fatal("first lock was refused")
	}

	if _, ok, err := locker.TryLock(ctx, jobID); err != nil {
		t.Fatalf("second lock: %v", err)
	} else if ok {
		t.Fatal("second lock was granted while the first was held")
	}

	release()

	second, ok, err := locker.TryLock(ctx, jobID)
	if err != nil {
		t.Fatalf("reacquire: %v", err)
	}
	if !ok {
		t.Fatal("lock was not released")
	}
	second()
}

func TestLockerAllowsDifferentJobs(t *testing.T) {
	pool := testPool(t)
	locker := NewLocker(pool)
	ctx := context.Background()
	first, second := time.Now().UnixNano(), time.Now().UnixNano()+1

	releaseA, ok, err := locker.TryLock(ctx, first)
	if err != nil || !ok {
		t.Fatalf("lock first: ok=%v err=%v", ok, err)
	}
	defer releaseA()

	releaseB, ok, err := locker.TryLock(ctx, second)
	if err != nil || !ok {
		t.Fatalf("lock second: ok=%v err=%v", ok, err)
	}
	releaseB()
}
