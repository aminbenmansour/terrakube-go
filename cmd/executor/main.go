// Command executor runs Terraform/OpenTofu jobs from the Redis queue and
// exposes health endpoints.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aminbenmansour/terrakube-go/internal/config"
	"github.com/aminbenmansour/terrakube-go/internal/db"
	"github.com/aminbenmansour/terrakube-go/internal/executor"
	"github.com/aminbenmansour/terrakube-go/internal/httpapi"
	"github.com/aminbenmansour/terrakube-go/internal/jobstatus"
	"github.com/aminbenmansour/terrakube-go/internal/logging"
	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/secretstore"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "terrakube-executor: %v\n", err)
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

	cipher, err := secretstore.NewCipher(cfg.EncryptionKeys, cfg.EncryptionActiveKey)
	if err != nil {
		return fmt.Errorf("configure secret encryption: %w", err)
	}
	if cipher == nil {
		logger.Warn("secret encryption is disabled: set TERRAKUBE_ENCRYPTION_KEYS")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

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

	events, err := jobstatus.NewBus(cfg.RedisURL)
	if err != nil {
		return err
	}
	defer events.Close()

	storeBackend, err := storage.New(ctx, cfg.Storage)
	if err != nil {
		return err
	}

	runner := executor.New(executor.Deps{
		Store:               jobstatus.Wrap(secretstore.Wrap(store.New(pool), cipher), events.Publish),
		Storage:             storeBackend,
		Queue:               q,
		Locker:              db.NewLocker(pool),
		Layout:              tfstate.NewLayout(cfg.Storage.Type),
		Logger:              logger,
		PublicURL:           cfg.PublicURL,
		InternalTokenSecret: cfg.InternalTokenSecret,
		Binary:              cfg.ExecutorBinary,
		WorkDir:             cfg.ExecutorWorkDir,
		JobTimeout:          cfg.ExecutorJobTimeout,
	})

	go func() {
		if err := runner.Run(ctx); err != nil && ctx.Err() == nil {
			logger.Error("executor stopped", "error", err)
			stop()
		}
	}()

	logger.Info("executor started", "storage", cfg.Storage.Type)

	handler := httpapi.New(httpapi.Deps{
		Logger:    logger,
		Service:   "executor",
		Started:   time.Now(),
		DB:        pool,
		Redis:     q,
		Storage:   storeBackend,
		RateLimit: cfg.RateLimit,
		RateBurst: cfg.RateBurst,
	})

	return httpapi.Serve(ctx, logger, cfg.ExecutorHTTPAddr, handler, cfg.ShutdownTimeout)
}
