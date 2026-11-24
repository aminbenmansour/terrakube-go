package secretstore

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/db"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/migrations"
)

func TestBackfillEncryptsPlaintextRows(t *testing.T) {
	dsn := os.Getenv("TERRAKUBE_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TERRAKUBE_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	defer pool.Close()
	if err := db.Migrate(ctx, pool, migrations.Files, "."); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	// Backfill touches every row in the shared test database. Run it in a
	// transaction and roll back so rows other packages created are untouched.
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)
	q := store.New(tx)

	c, err := NewCipher("k1:"+testKey(1), "")
	if err != nil {
		t.Fatalf("NewCipher() = %v", err)
	}

	orgID := uuid.NewString()
	if _, err := q.CreateOrganization(ctx, store.CreateOrganizationParams{
		ID: orgID, Name: "secretstore-" + orgID[:8], ExecutionMode: "remote",
	}); err != nil {
		t.Fatalf("create organization: %v", err)
	}

	wsID := uuid.NewString()
	if _, err := q.CreateWorkspace(ctx, store.CreateWorkspaceParams{
		ID: wsID, Name: "secretstore", IacType: "terraform", OrganizationID: orgID,
	}); err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	varID := uuid.NewString()
	if _, err := q.CreateVariable(ctx, store.CreateVariableParams{
		ID: varID, VariableKey: "PLAIN", VariableValue: "plain",
		VariableCategory: "TERRAFORM", WorkspaceID: wsID,
	}); err != nil {
		t.Fatalf("create variable: %v", err)
	}

	vcsID := uuid.NewString()
	if _, err := q.CreateVCS(ctx, store.CreateVCSParams{
		ID: vcsID, Name: "vcs", Description: "test", VcsType: "GITHUB",
		ClientSecret: ptr("client"), AccessToken: ptr("token"), OrganizationID: orgID,
	}); err != nil {
		t.Fatalf("create vcs: %v", err)
	}

	changed, err := Backfill(ctx, q, c)
	if err != nil {
		t.Fatalf("Backfill() = %v", err)
	}
	if changed < 2 {
		t.Fatalf("Backfill() changed %d rows, want at least the variable and the vcs", changed)
	}

	rawVariable, err := q.GetVariable(ctx, varID)
	if err != nil {
		t.Fatalf("GetVariable() = %v", err)
	}
	assertCiphertext(t, "variable value", rawVariable.VariableValue)
	if plain, err := c.Decrypt(columnVariableValue, rawVariable.VariableValue); err != nil || plain != "plain" {
		t.Fatalf("decrypted variable = %q, %v; want plain", plain, err)
	}

	rawVCS, err := q.GetVCS(ctx, vcsID)
	if err != nil {
		t.Fatalf("GetVCS() = %v", err)
	}
	assertCiphertext(t, "client secret", *rawVCS.ClientSecret)
	assertCiphertext(t, "access token", *rawVCS.AccessToken)
	if rawVCS.PrivateKey != nil {
		t.Fatalf("private key = %q, want nil to stay nil", *rawVCS.PrivateKey)
	}

	// Other packages write plaintext rows to this shared database, so a global
	// "0 rows changed" assertion would be flaky; assert our rows are untouched.
	if _, err := Backfill(ctx, q, c); err != nil {
		t.Fatalf("second Backfill() = %v", err)
	}
	againVariable, err := q.GetVariable(ctx, varID)
	if err != nil {
		t.Fatalf("GetVariable() after second backfill = %v", err)
	}
	if againVariable.VariableValue != rawVariable.VariableValue {
		t.Fatal("second backfill re-encrypted the variable")
	}
	againVCS, err := q.GetVCS(ctx, vcsID)
	if err != nil {
		t.Fatalf("GetVCS() after second backfill = %v", err)
	}
	if *againVCS.ClientSecret != *rawVCS.ClientSecret || *againVCS.AccessToken != *rawVCS.AccessToken {
		t.Fatal("second backfill re-encrypted the vcs secrets")
	}
}

func TestBackfillWithoutCipherIsANoop(t *testing.T) {
	if changed, err := Backfill(context.Background(), nil, nil); changed != 0 || err != nil {
		t.Fatalf("Backfill(nil) = %d, %v; want 0, nil", changed, err)
	}
}
