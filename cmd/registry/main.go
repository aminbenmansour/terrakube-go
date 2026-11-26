// Command registry serves the private Terraform module and provider registry.
// The registry protocol lands after the CLI workflow milestone.
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
		fmt.Fprintf(os.Stderr, "terrakube-registry: %v\n", err)
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

	logger.Info("registry started")
	logger.Warn("registry protocol is not implemented yet; see milestone M6")

	handler := httpapi.New(httpapi.Deps{
		Logger:       logger,
		Service:      "registry",
		Started:      time.Now(),
		RateLimit:    cfg.RateLimit,
		RateBurst:    cfg.RateBurst,
		Metrics:      httpapi.NewMetrics(),
		MetricsToken: cfg.MetricsToken,
	})

	return httpapi.Serve(ctx, logger, cfg.RegistryHTTPAddr, handler, cfg.ShutdownTimeout)
}
