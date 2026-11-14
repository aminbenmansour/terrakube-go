// Command executor runs Terraform/OpenTofu jobs. Job processing lands in M4;
// for now it exposes health endpoints.
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
	"github.com/aminbenmansour/terrakube-go/internal/httpapi"
	"github.com/aminbenmansour/terrakube-go/internal/logging"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("executor started")
	logger.Warn("job processing is not implemented yet; see milestone M4")

	handler := httpapi.New(httpapi.Deps{
		Logger:  logger,
		Service: "executor",
		Started: time.Now(),
	})

	return httpapi.Serve(ctx, logger, cfg.ExecutorHTTPAddr, handler, cfg.ShutdownTimeout)
}
