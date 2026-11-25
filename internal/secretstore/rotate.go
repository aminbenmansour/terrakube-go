package secretstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const defaultRotationBatch = 500

// RotateOptions controls a rotation pass. DryRun reads and validates every
// value but writes nothing. BatchSize defaults to 500.
type RotateOptions struct {
	DryRun    bool
	BatchSize int
}

// RotationCounts summarizes the values seen for one column during a pass. In a
// dry run Plaintext and OldKeys are what would change; in a live run they are
// what was re-encrypted.
type RotationCounts struct {
	Empty     int
	Current   int
	Plaintext int
	OldKeys   map[string]int
}

// RotationReport is keyed by table.column.
type RotationReport map[string]*RotationCounts

// NeedsRotation reports whether a pass saw a value that is not under the
// active key.
func (r RotationReport) NeedsRotation() bool {
	for _, counts := range r {
		if counts.Plaintext > 0 {
			return true
		}
		for _, n := range counts.OldKeys {
			if n > 0 {
				return true
			}
		}
	}
	return false
}

func (r RotationReport) counts(column string) *RotationCounts {
	counts := r[column]
	if counts == nil {
		counts = &RotationCounts{OldKeys: map[string]int{}}
		r[column] = counts
	}
	return counts
}

// Rotate re-encrypts every secret value that is not under the active key. It
// pages each table by id and is idempotent: values already current are only
// decrypted to validate them, and updates compare-and-set the values that were
// read so a concurrent write is never clobbered.
func Rotate(ctx context.Context, q *store.Queries, c *Cipher, opts RotateOptions) (RotationReport, error) {
	report := RotationReport{}
	if c == nil {
		return report, fmt.Errorf("secretstore: rotation requires an encryption keyring")
	}
	if opts.BatchSize <= 0 {
		opts.BatchSize = defaultRotationBatch
	}
	for _, rotate := range []func(context.Context, *store.Queries, *Cipher, RotateOptions, RotationReport) error{
		rotateVariables,
		rotateGlobalVariables,
		rotateSSHKeys,
		rotateVCS,
		rotateNotifications,
	} {
		if err := rotate(ctx, q, c, opts, report); err != nil {
			return report, err
		}
	}
	return report, nil
}

// rotateValue seals a value under the active key when needed and records what
// it was. Every envelope is decrypted, so a corrupt or unknown-key value stops
// the pass instead of being skipped.
func rotateValue(c *Cipher, column, value string, opts RotateOptions, counts *RotationCounts) (sealed string, changed bool, err error) {
	if value == "" {
		counts.Empty++
		return value, false, nil
	}
	if !strings.HasPrefix(value, envelopePrefix) {
		counts.Plaintext++
		if opts.DryRun {
			return value, false, nil
		}
		sealed, err := c.Encrypt(column, value)
		if err != nil {
			return "", false, err
		}
		return sealed, true, nil
	}
	plain, err := c.Decrypt(column, value)
	if err != nil {
		return "", false, err
	}
	key := envelopeKey(value)
	if key == c.ActiveKeyID() {
		counts.Current++
		return value, false, nil
	}
	counts.OldKeys[key]++
	if opts.DryRun {
		return value, false, nil
	}
	sealed, err = c.Encrypt(column, plain)
	if err != nil {
		return "", false, err
	}
	return sealed, true, nil
}

func envelopeKey(value string) string {
	id, _, _ := strings.Cut(strings.TrimPrefix(value, envelopePrefix), ":")
	return id
}

func rotatePtr(c *Cipher, column string, value **string, opts RotateOptions, counts *RotationCounts) (bool, error) {
	if *value == nil {
		return false, nil
	}
	sealed, changed, err := rotateValue(c, column, **value, opts, counts)
	if err != nil {
		return false, err
	}
	if changed {
		*value = &sealed
	}
	return changed, nil
}

func rotateVariables(ctx context.Context, q *store.Queries, c *Cipher, opts RotateOptions, report RotationReport) error {
	counts := report.counts(columnVariableValue)
	last := ""
	for {
		rows, err := q.ListVariableSecretsAfter(ctx, store.ListVariableSecretsAfterParams{
			ID: last, Limit: int32(opts.BatchSize),
		})
		if err != nil {
			return err
		}
		for _, row := range rows {
			sealed, changed, err := rotateValue(c, columnVariableValue, row.VariableValue, opts, counts)
			if err != nil {
				return fmt.Errorf("rotate %s id %s: %w", columnVariableValue, row.ID, err)
			}
			if !changed {
				continue
			}
			if err := q.SetVariableSecret(ctx, store.SetVariableSecretParams{
				ID: row.ID, VariableValue: sealed, VariableValue_2: row.VariableValue,
			}); err != nil {
				return fmt.Errorf("rotate %s id %s: %w", columnVariableValue, row.ID, err)
			}
		}
		if len(rows) < opts.BatchSize {
			return nil
		}
		last = rows[len(rows)-1].ID
	}
}

func rotateGlobalVariables(ctx context.Context, q *store.Queries, c *Cipher, opts RotateOptions, report RotationReport) error {
	counts := report.counts(columnGlobalVariableValue)
	last := ""
	for {
		rows, err := q.ListGlobalVariableSecretsAfter(ctx, store.ListGlobalVariableSecretsAfterParams{
			ID: last, Limit: int32(opts.BatchSize),
		})
		if err != nil {
			return err
		}
		for _, row := range rows {
			sealed, changed, err := rotateValue(c, columnGlobalVariableValue, row.VariableValue, opts, counts)
			if err != nil {
				return fmt.Errorf("rotate %s id %s: %w", columnGlobalVariableValue, row.ID, err)
			}
			if !changed {
				continue
			}
			if err := q.SetGlobalVariableSecret(ctx, store.SetGlobalVariableSecretParams{
				ID: row.ID, VariableValue: sealed, VariableValue_2: row.VariableValue,
			}); err != nil {
				return fmt.Errorf("rotate %s id %s: %w", columnGlobalVariableValue, row.ID, err)
			}
		}
		if len(rows) < opts.BatchSize {
			return nil
		}
		last = rows[len(rows)-1].ID
	}
}

func rotateSSHKeys(ctx context.Context, q *store.Queries, c *Cipher, opts RotateOptions, report RotationReport) error {
	counts := report.counts(columnSSHPrivateKey)
	last := ""
	for {
		rows, err := q.ListSSHKeySecretsAfter(ctx, store.ListSSHKeySecretsAfterParams{
			ID: last, Limit: int32(opts.BatchSize),
		})
		if err != nil {
			return err
		}
		for _, row := range rows {
			if row.PrivateKey == nil {
				continue
			}
			sealed, changed, err := rotateValue(c, columnSSHPrivateKey, *row.PrivateKey, opts, counts)
			if err != nil {
				return fmt.Errorf("rotate %s id %s: %w", columnSSHPrivateKey, row.ID, err)
			}
			if !changed {
				continue
			}
			if err := q.SetSSHKeySecret(ctx, store.SetSSHKeySecretParams{
				ID: row.ID, PrivateKey: &sealed, PrivateKey_2: row.PrivateKey,
			}); err != nil {
				return fmt.Errorf("rotate %s id %s: %w", columnSSHPrivateKey, row.ID, err)
			}
		}
		if len(rows) < opts.BatchSize {
			return nil
		}
		last = rows[len(rows)-1].ID
	}
}

func rotateVCS(ctx context.Context, q *store.Queries, c *Cipher, opts RotateOptions, report RotationReport) error {
	last := ""
	for {
		rows, err := q.ListVCSSecretsAfter(ctx, store.ListVCSSecretsAfterParams{
			ID: last, Limit: int32(opts.BatchSize),
		})
		if err != nil {
			return err
		}
		for _, row := range rows {
			params := store.SetVCSSecretsParams{
				ID:             row.ID,
				ClientSecret:   row.ClientSecret,
				AccessToken:    row.AccessToken,
				RefreshToken:   row.RefreshToken,
				PrivateKey:     row.PrivateKey,
				ClientSecret_2: row.ClientSecret,
				AccessToken_2:  row.AccessToken,
				RefreshToken_2: row.RefreshToken,
				PrivateKey_2:   row.PrivateKey,
			}
			dirty := false
			for _, column := range []struct {
				name  string
				value **string
			}{
				{columnVCSClientSecret, &params.ClientSecret},
				{columnVCSAccessToken, &params.AccessToken},
				{columnVCSRefreshToken, &params.RefreshToken},
				{columnVCSPrivateKey, &params.PrivateKey},
			} {
				changed, err := rotatePtr(c, column.name, column.value, opts, report.counts(column.name))
				if err != nil {
					return fmt.Errorf("rotate %s id %s: %w", column.name, row.ID, err)
				}
				dirty = dirty || changed
			}
			if !dirty {
				continue
			}
			if err := q.SetVCSSecrets(ctx, params); err != nil {
				return fmt.Errorf("rotate vcs id %s: %w", row.ID, err)
			}
		}
		if len(rows) < opts.BatchSize {
			return nil
		}
		last = rows[len(rows)-1].ID
	}
}

func rotateNotifications(ctx context.Context, q *store.Queries, c *Cipher, opts RotateOptions, report RotationReport) error {
	counts := report.counts(columnNotificationSecret)
	last := ""
	for {
		rows, err := q.ListNotificationSecretsAfter(ctx, store.ListNotificationSecretsAfterParams{
			ID: last, Limit: int32(opts.BatchSize),
		})
		if err != nil {
			return err
		}
		for _, row := range rows {
			if row.SigningSecret == nil {
				continue
			}
			sealed, changed, err := rotateValue(c, columnNotificationSecret, *row.SigningSecret, opts, counts)
			if err != nil {
				return fmt.Errorf("rotate %s id %s: %w", columnNotificationSecret, row.ID, err)
			}
			if !changed {
				continue
			}
			if err := q.SetNotificationSecret(ctx, store.SetNotificationSecretParams{
				ID: row.ID, SigningSecret: &sealed, SigningSecret_2: row.SigningSecret,
			}); err != nil {
				return fmt.Errorf("rotate %s id %s: %w", columnNotificationSecret, row.ID, err)
			}
		}
		if len(rows) < opts.BatchSize {
			return nil
		}
		last = rows[len(rows)-1].ID
	}
}
