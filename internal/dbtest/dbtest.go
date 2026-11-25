// Package dbtest gives database-backed tests a private, migrated database.
package dbtest

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aminbenmansour/terrakube-go/internal/db"
	"github.com/aminbenmansour/terrakube-go/migrations"
)

// Pool returns a pool on a scratch database whose migrations have run. The
// database is dropped when the test ends, which keeps tests that scan whole
// tables independent of a developer's local rows. The test is skipped when
// TERRAKUBE_TEST_DATABASE_URL is not set.
func Pool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TERRAKUBE_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TERRAKUBE_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	t.Cleanup(cancel)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("parse test database URL: %v", err)
	}
	name := "terrakube_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	quoted := pgx.Identifier{name}.Sanitize()
	adminCfg := poolCfg.ConnConfig.Copy()
	adminCfg.Database = "postgres"
	admin, err := pgx.ConnectConfig(ctx, adminCfg)
	if err != nil {
		t.Fatalf("connect admin database: %v", err)
	}
	if _, err := admin.Exec(ctx, "CREATE DATABASE "+quoted); err != nil {
		_ = admin.Close(context.Background())
		t.Fatalf("create test database: %v", err)
	}
	// Registered before the pool's cleanup so it runs after it: cleanups are LIFO.
	t.Cleanup(func() {
		_, _ = admin.Exec(context.Background(), "DROP DATABASE IF EXISTS "+quoted)
		_ = admin.Close(context.Background())
	})

	poolCfg.ConnConfig.Database = name
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		t.Fatalf("connect test database: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := db.Migrate(ctx, pool, migrations.Files, "."); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	return pool
}
