// Command migrate applies pending SQL migrations and exits. The api binary also
// migrates on startup when TERRAKUBE_AUTO_MIGRATE is enabled.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/aminbenmansour/terrakube-go/internal/config"
	"github.com/aminbenmansour/terrakube-go/internal/db"
	"github.com/aminbenmansour/terrakube-go/internal/logging"
	"github.com/aminbenmansour/terrakube-go/migrations"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "terrakube-migrate: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := logging.New(cfg.LogLevel, cfg.Env)
	slog.SetDefault(logger)

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool, migrations.Files, "."); err != nil {
		return err
	}

	logger.Info("database migrations applied")
	return nil
}
