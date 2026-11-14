// Command api serves the Terrakube API. The management and Terraform
// endpoints land in later milestones; for now it exposes health endpoints.
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

	logger.Info("api started")
	logger.Warn("management API is not implemented yet; see milestone M1")

	handler := httpapi.New(httpapi.Deps{
		Logger:  logger,
		Service: "api",
		Started: time.Now(),
	})

	return httpapi.Serve(ctx, logger, cfg.HTTPAddr, handler, cfg.ShutdownTimeout)
}
