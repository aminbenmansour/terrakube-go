package secretstore

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/db"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/migrations"
)

func TestWrappedStoreRoundTrip(t *testing.T) {
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

	raw := store.New(pool)
	c, err := NewCipher("k1:"+testKey(1), "")
	if err != nil {
		t.Fatalf("NewCipher() = %v", err)
	}
	wrapped := Wrap(raw, c)

	orgID := uuid.NewString()
	if _, err := raw.CreateOrganization(ctx, store.CreateOrganizationParams{
		ID: orgID, Name: "secretstore-rt-" + orgID[:8], ExecutionMode: "remote",
	}); err != nil {
		t.Fatalf("create organization: %v", err)
	}
	t.Cleanup(func() { _ = raw.DeleteOrganization(context.Background(), orgID) })

	wsID := uuid.NewString()
	if _, err := raw.CreateWorkspace(ctx, store.CreateWorkspaceParams{
		ID: wsID, Name: "secretstore", IacType: "terraform", OrganizationID: orgID,
	}); err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	varID := uuid.NewString()
	if _, err := wrapped.CreateVariable(ctx, store.CreateVariableParams{
		ID: varID, VariableKey: "SECRET", VariableValue: "original",
		VariableCategory: "TERRAFORM", WorkspaceID: wsID,
	}); err != nil {
		t.Fatalf("wrapped CreateVariable() = %v", err)
	}

	stored, err := raw.GetVariable(ctx, varID)
	if err != nil {
		t.Fatalf("raw GetVariable() = %v", err)
	}
	if !strings.HasPrefix(stored.VariableValue, envelopePrefix) {
		t.Fatalf("stored value = %q, want ciphertext", stored.VariableValue)
	}

	got, err := wrapped.GetVariable(ctx, varID)
	if err != nil {
		t.Fatalf("wrapped GetVariable() = %v", err)
	}
	if got.VariableValue != "original" {
		t.Fatalf("wrapped value = %q, want original", got.VariableValue)
	}
}
