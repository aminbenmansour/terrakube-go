package secretstore

import (
	"context"
	"strings"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type fakeQuerier struct {
	store.Querier
	createdVariable store.CreateVariableParams
	updatedVariable store.UpdateVariableParams
	variable        store.Variable
	variables       []store.Variable
	createdSSH      store.CreateSSHKeyParams
	updatedSSH      store.UpdateSSHKeyParams
	ssh             store.Ssh
	createdVCS      store.CreateVCSParams
	updatedVCS      store.UpdateVCSParams
	vcs             store.Vc
	createdNotif    store.CreateNotificationConfigurationParams
	updatedNotif    store.UpdateNotificationConfigurationParams
	notif           store.NotificationConfiguration
}

func (f *fakeQuerier) CreateVariable(_ context.Context, arg store.CreateVariableParams) (store.Variable, error) {
	f.createdVariable = arg
	return store.Variable{ID: arg.ID, VariableValue: arg.VariableValue, WorkspaceID: arg.WorkspaceID}, nil
}

func (f *fakeQuerier) UpdateVariable(_ context.Context, arg store.UpdateVariableParams) (store.Variable, error) {
	f.updatedVariable = arg
	return store.Variable{ID: arg.ID, VariableValue: arg.VariableValue}, nil
}

func (f *fakeQuerier) GetVariable(context.Context, string) (store.Variable, error) {
	return f.variable, nil
}

func (f *fakeQuerier) ListVariablesByWorkspace(context.Context, string) ([]store.Variable, error) {
	return f.variables, nil
}

func (f *fakeQuerier) CreateSSHKey(_ context.Context, arg store.CreateSSHKeyParams) (store.Ssh, error) {
	f.createdSSH = arg
	return store.Ssh{ID: arg.ID, PrivateKey: arg.PrivateKey}, nil
}

func (f *fakeQuerier) UpdateSSHKey(_ context.Context, arg store.UpdateSSHKeyParams) (store.Ssh, error) {
	f.updatedSSH = arg
	return store.Ssh{ID: arg.ID, PrivateKey: arg.PrivateKey}, nil
}

func (f *fakeQuerier) GetSSHKey(context.Context, string) (store.Ssh, error) {
	return f.ssh, nil
}

func (f *fakeQuerier) CreateVCS(_ context.Context, arg store.CreateVCSParams) (store.Vc, error) {
	f.createdVCS = arg
	return store.Vc{
		ID:           arg.ID,
		ClientSecret: arg.ClientSecret,
		AccessToken:  arg.AccessToken,
		PrivateKey:   arg.PrivateKey,
	}, nil
}

func (f *fakeQuerier) UpdateVCS(_ context.Context, arg store.UpdateVCSParams) (store.Vc, error) {
	f.updatedVCS = arg
	return store.Vc{
		ID:           arg.ID,
		ClientSecret: arg.ClientSecret,
		AccessToken:  arg.AccessToken,
		PrivateKey:   arg.PrivateKey,
	}, nil
}

func (f *fakeQuerier) GetVCS(context.Context, string) (store.Vc, error) {
	return f.vcs, nil
}

func (f *fakeQuerier) CreateNotificationConfiguration(_ context.Context, arg store.CreateNotificationConfigurationParams) (store.NotificationConfiguration, error) {
	f.createdNotif = arg
	return store.NotificationConfiguration{ID: arg.ID, SigningSecret: arg.SigningSecret}, nil
}

func (f *fakeQuerier) UpdateNotificationConfiguration(_ context.Context, arg store.UpdateNotificationConfigurationParams) (store.NotificationConfiguration, error) {
	f.updatedNotif = arg
	return store.NotificationConfiguration{ID: arg.ID, SigningSecret: arg.SigningSecret}, nil
}

func (f *fakeQuerier) GetNotificationConfiguration(context.Context, string) (store.NotificationConfiguration, error) {
	return f.notif, nil
}

func assertCiphertext(t *testing.T, name, got string) {
	t.Helper()
	if !strings.HasPrefix(got, "enc:v1:") {
		t.Fatalf("%s = %q, want ciphertext", name, got)
	}
}

func ptr(s string) *string { return &s }

func TestWrapCreateVariableEncryptsAndDecrypts(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	q := &fakeQuerier{}
	row, err := Wrap(q, c).CreateVariable(context.Background(), store.CreateVariableParams{
		ID: "v1", VariableValue: "plain",
	})
	if err != nil {
		t.Fatalf("CreateVariable() = %v", err)
	}
	assertCiphertext(t, "stored value", q.createdVariable.VariableValue)
	if row.VariableValue != "plain" {
		t.Fatalf("returned value = %q, want plain", row.VariableValue)
	}
}

func TestWrapReadsDecryptVariables(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	sealed, _ := c.Encrypt(columnVariableValue, "secret")
	q := &fakeQuerier{
		variable:  store.Variable{ID: "v1", VariableValue: sealed},
		variables: []store.Variable{{ID: "v2", VariableValue: sealed}},
	}
	wrapped := Wrap(q, c)

	got, err := wrapped.GetVariable(context.Background(), "v1")
	if err != nil || got.VariableValue != "secret" {
		t.Fatalf("GetVariable() = %q, %v", got.VariableValue, err)
	}
	list, err := wrapped.ListVariablesByWorkspace(context.Background(), "ws")
	if err != nil || list[0].VariableValue != "secret" {
		t.Fatalf("ListVariablesByWorkspace() = %q, %v", list[0].VariableValue, err)
	}
}

func TestWrapUpdateVariableEncrypts(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	q := &fakeQuerier{}
	row, err := Wrap(q, c).UpdateVariable(context.Background(), store.UpdateVariableParams{
		ID: "v1", VariableValue: "changed",
	})
	if err != nil {
		t.Fatalf("UpdateVariable() = %v", err)
	}
	assertCiphertext(t, "updated value", q.updatedVariable.VariableValue)
	if row.VariableValue != "changed" {
		t.Fatalf("returned value = %q, want changed", row.VariableValue)
	}
}

func TestWrapSSHKeys(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	q := &fakeQuerier{}
	wrapped := Wrap(q, c)

	created, err := wrapped.CreateSSHKey(context.Background(), store.CreateSSHKeyParams{ID: "s1", PrivateKey: ptr("key")})
	if err != nil {
		t.Fatalf("CreateSSHKey() = %v", err)
	}
	assertCiphertext(t, "stored private key", *q.createdSSH.PrivateKey)
	if *created.PrivateKey != "key" {
		t.Fatalf("returned private key = %q, want key", *created.PrivateKey)
	}

	updated, err := wrapped.UpdateSSHKey(context.Background(), store.UpdateSSHKeyParams{ID: "s1", PrivateKey: ptr("new")})
	if err != nil {
		t.Fatalf("UpdateSSHKey() = %v", err)
	}
	assertCiphertext(t, "updated private key", *q.updatedSSH.PrivateKey)
	if *updated.PrivateKey != "new" {
		t.Fatalf("returned private key = %q, want new", *updated.PrivateKey)
	}

	sealed, _ := c.Encrypt(columnSSHPrivateKey, "stored")
	q.ssh = store.Ssh{ID: "s1", PrivateKey: &sealed}
	got, err := wrapped.GetSSHKey(context.Background(), "s1")
	if err != nil || *got.PrivateKey != "stored" {
		t.Fatalf("GetSSHKey() = %v, %v", got.PrivateKey, err)
	}
}

func TestWrapVCS(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	q := &fakeQuerier{}
	wrapped := Wrap(q, c)

	created, err := wrapped.CreateVCS(context.Background(), store.CreateVCSParams{
		ID: "c1", ClientSecret: ptr("client"), AccessToken: ptr("token"), PrivateKey: ptr("key"),
	})
	if err != nil {
		t.Fatalf("CreateVCS() = %v", err)
	}
	assertCiphertext(t, "stored client secret", *q.createdVCS.ClientSecret)
	assertCiphertext(t, "stored access token", *q.createdVCS.AccessToken)
	assertCiphertext(t, "stored private key", *q.createdVCS.PrivateKey)
	if *created.ClientSecret != "client" || *created.AccessToken != "token" || *created.PrivateKey != "key" {
		t.Fatalf("returned VCS secrets = %q/%q/%q, want plaintext",
			*created.ClientSecret, *created.AccessToken, *created.PrivateKey)
	}

	updated, err := wrapped.UpdateVCS(context.Background(), store.UpdateVCSParams{
		ID: "c1", ClientSecret: ptr("client2"), AccessToken: ptr("token2"), PrivateKey: ptr("key2"),
	})
	if err != nil {
		t.Fatalf("UpdateVCS() = %v", err)
	}
	assertCiphertext(t, "updated client secret", *q.updatedVCS.ClientSecret)
	if *updated.ClientSecret != "client2" {
		t.Fatalf("returned client secret = %q, want client2", *updated.ClientSecret)
	}

	sealedSecret, _ := c.Encrypt(columnVCSClientSecret, "client")
	sealedToken, _ := c.Encrypt(columnVCSAccessToken, "token")
	sealedRefresh, _ := c.Encrypt(columnVCSRefreshToken, "refresh")
	sealedKey, _ := c.Encrypt(columnVCSPrivateKey, "key")
	q.vcs = store.Vc{
		ID: "c1", ClientSecret: &sealedSecret, AccessToken: &sealedToken,
		RefreshToken: &sealedRefresh, PrivateKey: &sealedKey,
	}
	got, err := wrapped.GetVCS(context.Background(), "c1")
	if err != nil {
		t.Fatalf("GetVCS() = %v", err)
	}
	if *got.ClientSecret != "client" || *got.AccessToken != "token" ||
		*got.RefreshToken != "refresh" || *got.PrivateKey != "key" {
		t.Fatalf("GetVCS() secrets = %q/%q/%q/%q, want plaintext",
			*got.ClientSecret, *got.AccessToken, *got.RefreshToken, *got.PrivateKey)
	}
}

func TestWrapNotificationSigningSecret(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	q := &fakeQuerier{}
	wrapped := Wrap(q, c)

	created, err := wrapped.CreateNotificationConfiguration(context.Background(),
		store.CreateNotificationConfigurationParams{ID: "n1", SigningSecret: ptr("sign")})
	if err != nil {
		t.Fatalf("CreateNotificationConfiguration() = %v", err)
	}
	assertCiphertext(t, "stored signing secret", *q.createdNotif.SigningSecret)
	if *created.SigningSecret != "sign" {
		t.Fatalf("returned signing secret = %q, want sign", *created.SigningSecret)
	}

	sealed, _ := c.Encrypt(columnNotificationSecret, "sign")
	q.notif = store.NotificationConfiguration{ID: "n1", SigningSecret: &sealed}
	got, err := wrapped.GetNotificationConfiguration(context.Background(), "n1")
	if err != nil || *got.SigningSecret != "sign" {
		t.Fatalf("GetNotificationConfiguration() = %v, %v", got.SigningSecret, err)
	}
}

func TestWrapKeepsNilSecretsNil(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	q := &fakeQuerier{}
	wrapped := Wrap(q, c)

	if _, err := wrapped.CreateSSHKey(context.Background(), store.CreateSSHKeyParams{ID: "s1"}); err != nil {
		t.Fatalf("CreateSSHKey() = %v", err)
	}
	if q.createdSSH.PrivateKey != nil {
		t.Fatalf("stored private key = %v, want nil", *q.createdSSH.PrivateKey)
	}
	if _, err := wrapped.CreateNotificationConfiguration(context.Background(),
		store.CreateNotificationConfigurationParams{ID: "n1"}); err != nil {
		t.Fatalf("CreateNotificationConfiguration() = %v", err)
	}
	if q.createdNotif.SigningSecret != nil {
		t.Fatalf("stored signing secret = %v, want nil", *q.createdNotif.SigningSecret)
	}
}

func TestWrapWithoutCipherPassesThrough(t *testing.T) {
	q := &fakeQuerier{}
	wrapped := Wrap(q, nil)
	if wrapped != q {
		t.Fatal("Wrap(q, nil) did not return the same querier")
	}
	if _, err := wrapped.CreateVariable(context.Background(), store.CreateVariableParams{
		ID: "v1", VariableValue: "plain",
	}); err != nil {
		t.Fatalf("CreateVariable() = %v", err)
	}
	if q.createdVariable.VariableValue != "plain" {
		t.Fatalf("stored value = %q, want plaintext", q.createdVariable.VariableValue)
	}
}
