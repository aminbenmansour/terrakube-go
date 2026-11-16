// Command api serves the Terrakube JSON:API, the Terraform TFE-compatible
// endpoints and the state backend.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aminbenmansour/terrakube-go/internal/api"
	internalauth "github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/config"
	"github.com/aminbenmansour/terrakube-go/internal/db"
	"github.com/aminbenmansour/terrakube-go/internal/httpapi"
	"github.com/aminbenmansour/terrakube-go/internal/logging"
	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfe"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
	"github.com/aminbenmansour/terrakube-go/migrations"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "terrakube-api: %v\n", err)
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if cfg.AutoMigrate {
		if err := db.Migrate(ctx, pool, migrations.Files, "."); err != nil {
			return err
		}
		logger.Info("database migrations applied")
	}

	q, err := queue.New(cfg.RedisURL)
	if err != nil {
		return err
	}
	defer q.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	err = q.Ping(pingCtx)
	cancel()
	if err != nil {
		return err
	}

	storeBackend, err := storage.New(ctx, cfg.Storage)
	if err != nil {
		return err
	}
	logger.Info("storage configured", "type", cfg.Storage.Type)

	var internalAuth *internalauth.Validator
	if cfg.InternalTokenSecret == "" {
		logger.Warn("internal token authentication is disabled: set TERRAKUBE_INTERNAL_TOKEN_SECRET to enable /tfstate/v1")
	} else {
		internalAuth, err = internalauth.NewValidator(cfg.InternalTokenSecret)
		if err != nil {
			return fmt.Errorf("configure internal token authentication: %w", err)
		}
	}

	handler := httpapi.New(httpapi.Deps{
		Logger:  logger,
		Service: "api",
		Started: time.Now(),
		DB:      pool,
		Redis:   q,
		Storage: storeBackend,
	})
	handler.Mount("/api/v1", api.New(api.Deps{
		Store:     store.New(pool),
		Logger:    logger,
		PublicURL: cfg.PublicURL,
		AuthToken: cfg.AuthToken,
	}))
	tfeDeps := tfe.Deps{
		Store:          store.New(pool),
		Storage:        storeBackend,
		Layout:         tfstate.NewLayout(cfg.Storage.Type),
		Logger:         logger,
		PublicURL:      cfg.PublicURL,
		AuthToken:      cfg.AuthToken,
		InternalAuth:   internalAuth,
		DexClientID:    cfg.DexClientID,
		DexIssuerURL:   cfg.DexIssuerURL,
		TerraformLogin: cfg.TerraformLogin,
	}
	handler.Get("/.well-known/terraform.json", tfe.Discovery(tfeDeps))
	handler.Mount("/remote/tfe/v2", tfe.New(tfeDeps))
	handler.Mount("/tfstate/v1", tfstate.New(tfstate.Deps{
		Store:             store.New(pool),
		Storage:           storeBackend,
		Logger:            logger,
		Layout:            tfstate.NewLayout(cfg.Storage.Type),
		Auth:              internalAuth,
		MaxStateSizeBytes: cfg.StateMaxSizeBytes,
		PublicURL:         cfg.PublicURL,
	}))

	return httpapi.Serve(ctx, logger, cfg.HTTPAddr, handler, cfg.ShutdownTimeout)
}
