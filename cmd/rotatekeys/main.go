// Command rotatekeys re-encrypts secret columns under the active keyring key.
// It is safe to stop and rerun, and -dry-run reports what would change.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/aminbenmansour/terrakube-go/internal/config"
	"github.com/aminbenmansour/terrakube-go/internal/db"
	"github.com/aminbenmansour/terrakube-go/internal/logging"
	"github.com/aminbenmansour/terrakube-go/internal/secretstore"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "report what would change without writing")
	flag.Parse()
	if err := run(*dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "terrakube-rotatekeys: %v\n", err)
		os.Exit(1)
	}
}

func run(dryRun bool) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	logger := logging.New(cfg.LogLevel, cfg.Env)

	cipher, err := secretstore.NewCipher(cfg.EncryptionKeys, cfg.EncryptionActiveKey)
	if err != nil {
		return err
	}
	if cipher == nil {
		return errors.New("TERRAKUBE_ENCRYPTION_KEYS is not set")
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	logger.Info("rotating secrets", "active_key", cipher.ActiveKeyID(), "dry_run", dryRun)
	report, err := secretstore.Rotate(ctx, store.New(pool), cipher, secretstore.RotateOptions{DryRun: dryRun})
	logReport(logger, "secret rotation", report)
	if err != nil {
		return err
	}
	if dryRun {
		return nil
	}

	verify, err := secretstore.Rotate(ctx, store.New(pool), cipher, secretstore.RotateOptions{DryRun: true})
	logReport(logger, "secret rotation verification", verify)
	if err != nil {
		return err
	}
	if verify.NeedsRotation() {
		return errors.New("values remain under an old key or in plaintext")
	}
	logger.Info("secret rotation complete", "active_key", cipher.ActiveKeyID())
	return nil
}

func logReport(logger *slog.Logger, msg string, report secretstore.RotationReport) {
	columns := make([]string, 0, len(report))
	for column := range report {
		columns = append(columns, column)
	}
	sort.Strings(columns)
	for _, column := range columns {
		counts := report[column]
		logger.Info(msg,
			"column", column,
			"empty", counts.Empty,
			"current", counts.Current,
			"plaintext", counts.Plaintext,
			"old_keys", formatOldKeys(counts.OldKeys),
		)
	}
}

func formatOldKeys(keys map[string]int) string {
	ids := make([]string, 0, len(keys))
	for id, n := range keys {
		if n > 0 {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("%s=%d", id, keys[id]))
	}
	return strings.Join(parts, " ")
}
