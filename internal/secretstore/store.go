package secretstore

import (
	"context"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// Column names double as the AEAD additional data, so ciphertext moved to
// another column fails to open.
const (
	columnVariableValue       = "variable.variable_value"
	columnGlobalVariableValue = "globalvar.variable_value"
	columnSSHPrivateKey       = "ssh.private_key"
	columnVCSClientSecret     = "vcs.client_secret"
	columnVCSAccessToken      = "vcs.access_token"
	columnVCSRefreshToken     = "vcs.refresh_token"
	columnVCSPrivateKey       = "vcs.private_key"
	columnNotificationSecret  = "notification_configuration.signing_secret"
)

type encryptedStore struct {
	store.Querier
	cipher *Cipher
}

// Wrap returns a querier that encrypts and decrypts the secret columns. It
// returns q unchanged when no cipher is configured.
func Wrap(q store.Querier, c *Cipher) store.Querier {
	if c == nil {
		return q
	}
	return &encryptedStore{Querier: q, cipher: c}
}

func (s *encryptedStore) encryptPtr(column string, v *string) (*string, error) {
	if v == nil {
		return nil, nil
	}
	sealed, err := s.cipher.Encrypt(column, *v)
	if err != nil {
		return nil, err
	}
	return &sealed, nil
}

func (s *encryptedStore) decryptPtr(column string, v *string) (*string, error) {
	if v == nil {
		return nil, nil
	}
	plain, err := s.cipher.Decrypt(column, *v)
	if err != nil {
		return nil, err
	}
	return &plain, nil
}

func (s *encryptedStore) CreateVariable(ctx context.Context, arg store.CreateVariableParams) (store.Variable, error) {
	sealed, err := s.cipher.Encrypt(columnVariableValue, arg.VariableValue)
	if err != nil {
		return store.Variable{}, err
	}
	arg.VariableValue = sealed
	row, err := s.Querier.CreateVariable(ctx, arg)
	if err != nil {
		return store.Variable{}, err
	}
	row.VariableValue, err = s.cipher.Decrypt(columnVariableValue, row.VariableValue)
	if err != nil {
		return store.Variable{}, err
	}
	return row, nil
}

func (s *encryptedStore) UpdateVariable(ctx context.Context, arg store.UpdateVariableParams) (store.Variable, error) {
	sealed, err := s.cipher.Encrypt(columnVariableValue, arg.VariableValue)
	if err != nil {
		return store.Variable{}, err
	}
	arg.VariableValue = sealed
	row, err := s.Querier.UpdateVariable(ctx, arg)
	if err != nil {
		return store.Variable{}, err
	}
	row.VariableValue, err = s.cipher.Decrypt(columnVariableValue, row.VariableValue)
	if err != nil {
		return store.Variable{}, err
	}
	return row, nil
}

func (s *encryptedStore) GetVariable(ctx context.Context, id string) (store.Variable, error) {
	row, err := s.Querier.GetVariable(ctx, id)
	if err != nil {
		return store.Variable{}, err
	}
	row.VariableValue, err = s.cipher.Decrypt(columnVariableValue, row.VariableValue)
	if err != nil {
		return store.Variable{}, err
	}
	return row, nil
}

func (s *encryptedStore) ListVariablesByWorkspace(ctx context.Context, workspaceID string) ([]store.Variable, error) {
	rows, err := s.Querier.ListVariablesByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].VariableValue, err = s.cipher.Decrypt(columnVariableValue, rows[i].VariableValue)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func (s *encryptedStore) CreateGlobalVariable(ctx context.Context, arg store.CreateGlobalVariableParams) (store.Globalvar, error) {
	sealed, err := s.cipher.Encrypt(columnGlobalVariableValue, arg.VariableValue)
	if err != nil {
		return store.Globalvar{}, err
	}
	arg.VariableValue = sealed
	row, err := s.Querier.CreateGlobalVariable(ctx, arg)
	if err != nil {
		return store.Globalvar{}, err
	}
	row.VariableValue, err = s.cipher.Decrypt(columnGlobalVariableValue, row.VariableValue)
	if err != nil {
		return store.Globalvar{}, err
	}
	return row, nil
}

func (s *encryptedStore) UpdateGlobalVariable(ctx context.Context, arg store.UpdateGlobalVariableParams) (store.Globalvar, error) {
	sealed, err := s.cipher.Encrypt(columnGlobalVariableValue, arg.VariableValue)
	if err != nil {
		return store.Globalvar{}, err
	}
	arg.VariableValue = sealed
	row, err := s.Querier.UpdateGlobalVariable(ctx, arg)
	if err != nil {
		return store.Globalvar{}, err
	}
	row.VariableValue, err = s.cipher.Decrypt(columnGlobalVariableValue, row.VariableValue)
	if err != nil {
		return store.Globalvar{}, err
	}
	return row, nil
}

func (s *encryptedStore) GetGlobalVariable(ctx context.Context, id string) (store.Globalvar, error) {
	row, err := s.Querier.GetGlobalVariable(ctx, id)
	if err != nil {
		return store.Globalvar{}, err
	}
	row.VariableValue, err = s.cipher.Decrypt(columnGlobalVariableValue, row.VariableValue)
	if err != nil {
		return store.Globalvar{}, err
	}
	return row, nil
}

func (s *encryptedStore) ListGlobalVariablesByOrganization(ctx context.Context, organizationID string) ([]store.Globalvar, error) {
	rows, err := s.Querier.ListGlobalVariablesByOrganization(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].VariableValue, err = s.cipher.Decrypt(columnGlobalVariableValue, rows[i].VariableValue)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func (s *encryptedStore) CreateSSHKey(ctx context.Context, arg store.CreateSSHKeyParams) (store.Ssh, error) {
	var err error
	arg.PrivateKey, err = s.encryptPtr(columnSSHPrivateKey, arg.PrivateKey)
	if err != nil {
		return store.Ssh{}, err
	}
	row, err := s.Querier.CreateSSHKey(ctx, arg)
	if err != nil {
		return store.Ssh{}, err
	}
	row.PrivateKey, err = s.decryptPtr(columnSSHPrivateKey, row.PrivateKey)
	if err != nil {
		return store.Ssh{}, err
	}
	return row, nil
}

func (s *encryptedStore) UpdateSSHKey(ctx context.Context, arg store.UpdateSSHKeyParams) (store.Ssh, error) {
	var err error
	arg.PrivateKey, err = s.encryptPtr(columnSSHPrivateKey, arg.PrivateKey)
	if err != nil {
		return store.Ssh{}, err
	}
	row, err := s.Querier.UpdateSSHKey(ctx, arg)
	if err != nil {
		return store.Ssh{}, err
	}
	row.PrivateKey, err = s.decryptPtr(columnSSHPrivateKey, row.PrivateKey)
	if err != nil {
		return store.Ssh{}, err
	}
	return row, nil
}

func (s *encryptedStore) GetSSHKey(ctx context.Context, id string) (store.Ssh, error) {
	row, err := s.Querier.GetSSHKey(ctx, id)
	if err != nil {
		return store.Ssh{}, err
	}
	row.PrivateKey, err = s.decryptPtr(columnSSHPrivateKey, row.PrivateKey)
	if err != nil {
		return store.Ssh{}, err
	}
	return row, nil
}

func (s *encryptedStore) ListSSHKeysByOrganization(ctx context.Context, organizationID string) ([]store.Ssh, error) {
	rows, err := s.Querier.ListSSHKeysByOrganization(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].PrivateKey, err = s.decryptPtr(columnSSHPrivateKey, rows[i].PrivateKey)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func (s *encryptedStore) CreateVCS(ctx context.Context, arg store.CreateVCSParams) (store.Vc, error) {
	var err error
	if arg.ClientSecret, err = s.encryptPtr(columnVCSClientSecret, arg.ClientSecret); err != nil {
		return store.Vc{}, err
	}
	if arg.AccessToken, err = s.encryptPtr(columnVCSAccessToken, arg.AccessToken); err != nil {
		return store.Vc{}, err
	}
	if arg.PrivateKey, err = s.encryptPtr(columnVCSPrivateKey, arg.PrivateKey); err != nil {
		return store.Vc{}, err
	}
	row, err := s.Querier.CreateVCS(ctx, arg)
	if err != nil {
		return store.Vc{}, err
	}
	return s.decryptVCS(row)
}

func (s *encryptedStore) UpdateVCS(ctx context.Context, arg store.UpdateVCSParams) (store.Vc, error) {
	var err error
	if arg.ClientSecret, err = s.encryptPtr(columnVCSClientSecret, arg.ClientSecret); err != nil {
		return store.Vc{}, err
	}
	if arg.AccessToken, err = s.encryptPtr(columnVCSAccessToken, arg.AccessToken); err != nil {
		return store.Vc{}, err
	}
	if arg.PrivateKey, err = s.encryptPtr(columnVCSPrivateKey, arg.PrivateKey); err != nil {
		return store.Vc{}, err
	}
	row, err := s.Querier.UpdateVCS(ctx, arg)
	if err != nil {
		return store.Vc{}, err
	}
	return s.decryptVCS(row)
}

func (s *encryptedStore) GetVCS(ctx context.Context, id string) (store.Vc, error) {
	row, err := s.Querier.GetVCS(ctx, id)
	if err != nil {
		return store.Vc{}, err
	}
	return s.decryptVCS(row)
}

func (s *encryptedStore) ListVCSByOrganization(ctx context.Context, organizationID string) ([]store.Vc, error) {
	rows, err := s.Querier.ListVCSByOrganization(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		if rows[i], err = s.decryptVCS(rows[i]); err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func (s *encryptedStore) decryptVCS(row store.Vc) (store.Vc, error) {
	var err error
	if row.ClientSecret, err = s.decryptPtr(columnVCSClientSecret, row.ClientSecret); err != nil {
		return store.Vc{}, err
	}
	if row.AccessToken, err = s.decryptPtr(columnVCSAccessToken, row.AccessToken); err != nil {
		return store.Vc{}, err
	}
	if row.RefreshToken, err = s.decryptPtr(columnVCSRefreshToken, row.RefreshToken); err != nil {
		return store.Vc{}, err
	}
	if row.PrivateKey, err = s.decryptPtr(columnVCSPrivateKey, row.PrivateKey); err != nil {
		return store.Vc{}, err
	}
	return row, nil
}

func (s *encryptedStore) CreateNotificationConfiguration(ctx context.Context, arg store.CreateNotificationConfigurationParams) (store.NotificationConfiguration, error) {
	var err error
	arg.SigningSecret, err = s.encryptPtr(columnNotificationSecret, arg.SigningSecret)
	if err != nil {
		return store.NotificationConfiguration{}, err
	}
	row, err := s.Querier.CreateNotificationConfiguration(ctx, arg)
	if err != nil {
		return store.NotificationConfiguration{}, err
	}
	row.SigningSecret, err = s.decryptPtr(columnNotificationSecret, row.SigningSecret)
	if err != nil {
		return store.NotificationConfiguration{}, err
	}
	return row, nil
}

func (s *encryptedStore) UpdateNotificationConfiguration(ctx context.Context, arg store.UpdateNotificationConfigurationParams) (store.NotificationConfiguration, error) {
	var err error
	arg.SigningSecret, err = s.encryptPtr(columnNotificationSecret, arg.SigningSecret)
	if err != nil {
		return store.NotificationConfiguration{}, err
	}
	row, err := s.Querier.UpdateNotificationConfiguration(ctx, arg)
	if err != nil {
		return store.NotificationConfiguration{}, err
	}
	row.SigningSecret, err = s.decryptPtr(columnNotificationSecret, row.SigningSecret)
	if err != nil {
		return store.NotificationConfiguration{}, err
	}
	return row, nil
}

func (s *encryptedStore) GetNotificationConfiguration(ctx context.Context, id string) (store.NotificationConfiguration, error) {
	row, err := s.Querier.GetNotificationConfiguration(ctx, id)
	if err != nil {
		return store.NotificationConfiguration{}, err
	}
	row.SigningSecret, err = s.decryptPtr(columnNotificationSecret, row.SigningSecret)
	if err != nil {
		return store.NotificationConfiguration{}, err
	}
	return row, nil
}

func (s *encryptedStore) ListNotificationConfigurationsByOrganization(ctx context.Context, organizationID string) ([]store.NotificationConfiguration, error) {
	rows, err := s.Querier.ListNotificationConfigurationsByOrganization(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].SigningSecret, err = s.decryptPtr(columnNotificationSecret, rows[i].SigningSecret)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func (s *encryptedStore) ListNotificationConfigurationsByWorkspace(ctx context.Context, workspaceID *string) ([]store.NotificationConfiguration, error) {
	rows, err := s.Querier.ListNotificationConfigurationsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].SigningSecret, err = s.decryptPtr(columnNotificationSecret, rows[i].SigningSecret)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}
