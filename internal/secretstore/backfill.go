package secretstore

import (
	"context"
	"strings"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// needsEncryption reports whether a stored value is legacy plaintext.
func needsEncryption(c *Cipher, value string) bool {
	return c != nil && value != "" && !strings.HasPrefix(value, envelopePrefix)
}

// Backfill encrypts legacy plaintext secret columns. It is idempotent: already
// enveloped values are skipped, and each update is a compare-and-set against
// the value that was read so a concurrent write is never clobbered. It returns
// the number of rows changed.
func Backfill(ctx context.Context, q *store.Queries, c *Cipher) (int, error) {
	if c == nil {
		return 0, nil
	}
	changed := 0
	for _, backfill := range []func(context.Context, *store.Queries, *Cipher) (int, error){
		backfillVariables,
		backfillGlobalVariables,
		backfillSSHKeys,
		backfillVCS,
		backfillNotifications,
	} {
		n, err := backfill(ctx, q, c)
		changed += n
		if err != nil {
			return changed, err
		}
	}
	return changed, nil
}

func backfillVariables(ctx context.Context, q *store.Queries, c *Cipher) (int, error) {
	rows, err := q.ListVariableSecrets(ctx)
	if err != nil {
		return 0, err
	}
	changed := 0
	for _, row := range rows {
		if !needsEncryption(c, row.VariableValue) {
			continue
		}
		sealed, err := c.Encrypt(columnVariableValue, row.VariableValue)
		if err != nil {
			return changed, err
		}
		if err := q.SetVariableSecret(ctx, store.SetVariableSecretParams{
			ID:              row.ID,
			VariableValue:   sealed,
			VariableValue_2: row.VariableValue,
		}); err != nil {
			return changed, err
		}
		changed++
	}
	return changed, nil
}

func backfillGlobalVariables(ctx context.Context, q *store.Queries, c *Cipher) (int, error) {
	rows, err := q.ListGlobalVariableSecrets(ctx)
	if err != nil {
		return 0, err
	}
	changed := 0
	for _, row := range rows {
		if !needsEncryption(c, row.VariableValue) {
			continue
		}
		sealed, err := c.Encrypt(columnGlobalVariableValue, row.VariableValue)
		if err != nil {
			return changed, err
		}
		if err := q.SetGlobalVariableSecret(ctx, store.SetGlobalVariableSecretParams{
			ID:              row.ID,
			VariableValue:   sealed,
			VariableValue_2: row.VariableValue,
		}); err != nil {
			return changed, err
		}
		changed++
	}
	return changed, nil
}

func backfillSSHKeys(ctx context.Context, q *store.Queries, c *Cipher) (int, error) {
	rows, err := q.ListSSHKeySecrets(ctx)
	if err != nil {
		return 0, err
	}
	changed := 0
	for _, row := range rows {
		if row.PrivateKey == nil || !needsEncryption(c, *row.PrivateKey) {
			continue
		}
		sealed, err := c.Encrypt(columnSSHPrivateKey, *row.PrivateKey)
		if err != nil {
			return changed, err
		}
		if err := q.SetSSHKeySecret(ctx, store.SetSSHKeySecretParams{
			ID:           row.ID,
			PrivateKey:   &sealed,
			PrivateKey_2: row.PrivateKey,
		}); err != nil {
			return changed, err
		}
		changed++
	}
	return changed, nil
}

func backfillVCS(ctx context.Context, q *store.Queries, c *Cipher) (int, error) {
	rows, err := q.ListVCSSecrets(ctx)
	if err != nil {
		return 0, err
	}
	changed := 0
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
		var err error
		if params.ClientSecret, err = encryptIfNeeded(c, columnVCSClientSecret, params.ClientSecret, &dirty); err != nil {
			return changed, err
		}
		if params.AccessToken, err = encryptIfNeeded(c, columnVCSAccessToken, params.AccessToken, &dirty); err != nil {
			return changed, err
		}
		if params.RefreshToken, err = encryptIfNeeded(c, columnVCSRefreshToken, params.RefreshToken, &dirty); err != nil {
			return changed, err
		}
		if params.PrivateKey, err = encryptIfNeeded(c, columnVCSPrivateKey, params.PrivateKey, &dirty); err != nil {
			return changed, err
		}
		if !dirty {
			continue
		}
		if err := q.SetVCSSecrets(ctx, params); err != nil {
			return changed, err
		}
		changed++
	}
	return changed, nil
}

func backfillNotifications(ctx context.Context, q *store.Queries, c *Cipher) (int, error) {
	rows, err := q.ListNotificationSecrets(ctx)
	if err != nil {
		return 0, err
	}
	changed := 0
	for _, row := range rows {
		if row.SigningSecret == nil || !needsEncryption(c, *row.SigningSecret) {
			continue
		}
		sealed, err := c.Encrypt(columnNotificationSecret, *row.SigningSecret)
		if err != nil {
			return changed, err
		}
		if err := q.SetNotificationSecret(ctx, store.SetNotificationSecretParams{
			ID:              row.ID,
			SigningSecret:   &sealed,
			SigningSecret_2: row.SigningSecret,
		}); err != nil {
			return changed, err
		}
		changed++
	}
	return changed, nil
}

func encryptIfNeeded(c *Cipher, column string, value *string, dirty *bool) (*string, error) {
	if value == nil || !needsEncryption(c, *value) {
		return value, nil
	}
	sealed, err := c.Encrypt(column, *value)
	if err != nil {
		return nil, err
	}
	*dirty = true
	return &sealed, nil
}
