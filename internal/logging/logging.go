// Package logging builds the structured logger used by every binary.
package logging

import (
	"log/slog"
	"os"
)

// New returns a slog logger. Production environments get JSON so the output can
// be shipped to a log aggregator; development gets human-readable text.
func New(level slog.Level, environment string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if environment == "prod" || environment == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
