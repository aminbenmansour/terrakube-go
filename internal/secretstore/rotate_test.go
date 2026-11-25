package secretstore

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

func openRotationTest(t *testing.T) (context.Context, *store.Queries) {
	t.Helper()
	// Rotation scans whole tables and decrypts every envelope, so it needs a
	// database of its own rather than a developer's local rows.
	return t.Context(), store.New(dbtest.Pool(t))
}

func rotationWorkspace(t *testing.T, ctx context.Context, q *store.Queries) (orgID, wsID string) {
	t.Helper()
	orgID = uuid.NewString()
	if _, err := q.CreateOrganization(ctx, store.CreateOrganizationParams{
		ID: orgID, Name: "rotation-" + orgID[:8], ExecutionMode: "remote",
	}); err != nil {
		t.Fatalf("create organization: %v", err)
	}
	wsID = uuid.NewString()
	if _, err := q.CreateWorkspace(ctx, store.CreateWorkspaceParams{
		ID: wsID, Name: "rotation", IacType: "terraform", OrganizationID: orgID,
	}); err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	return orgID, wsID
}

func createRawVariable(t *testing.T, ctx context.Context, q *store.Queries, wsID, key, value string) string {
	t.Helper()
	return createRawVariableWithID(t, ctx, q, uuid.NewString(), wsID, key, value)
}

func createRawVariableWithID(t *testing.T, ctx context.Context, q *store.Queries, id, wsID, key, value string) string {
	t.Helper()
	if _, err := q.CreateVariable(ctx, store.CreateVariableParams{
		ID: id, VariableKey: key, VariableValue: value,
		VariableCategory: "TERRAFORM", WorkspaceID: wsID,
	}); err != nil {
		t.Fatalf("create variable: %v", err)
	}
	return id
}

func rawVariable(t *testing.T, ctx context.Context, q *store.Queries, id string) string {
	t.Helper()
	row, err := q.GetVariable(ctx, id)
	if err != nil {
		t.Fatalf("raw GetVariable() = %v", err)
	}
	return row.VariableValue
}

func assertSealedUnder(t *testing.T, name, value, key string) {
	t.Helper()
	if !strings.HasPrefix(value, envelopePrefix+key+":") {
		t.Fatalf("%s is not sealed under %s", name, key)
	}
}

func TestRotateReencryptsOldKeyAndPlaintext(t *testing.T) {
	ctx, q := openRotationTest(t)
	old, err := NewCipher("k1:"+testKey(1), "")
	if err != nil {
		t.Fatalf("NewCipher(old) = %v", err)
	}
	next, err := NewCipher("k1:"+testKey(1)+",k2:"+testKey(2), "k2")
	if err != nil {
		t.Fatalf("NewCipher(next) = %v", err)
	}
	seal := func(column, value string) string {
		t.Helper()
		sealed, err := old.Encrypt(column, value)
		if err != nil {
			t.Fatalf("old Encrypt() = %v", err)
		}
		return sealed
	}

	orgID, wsID := rotationWorkspace(t, ctx, q)
	alpha := createRawVariable(t, ctx, q, wsID, "OLD_A", seal(columnVariableValue, "alpha"))
	beta := createRawVariable(t, ctx, q, wsID, "OLD_B", seal(columnVariableValue, "beta"))
	gamma := createRawVariable(t, ctx, q, wsID, "PLAIN", "gamma")
	empty := createRawVariable(t, ctx, q, wsID, "EMPTY", "")

	vcsID := uuid.NewString()
	if _, err := q.CreateVCS(ctx, store.CreateVCSParams{
		ID: vcsID, Name: "rotation", Description: "test", VcsType: "GITHUB",
		ClientSecret:   ptr(seal(columnVCSClientSecret, "client")),
		AccessToken:    ptr(seal(columnVCSAccessToken, "token")),
		OrganizationID: orgID,
	}); err != nil {
		t.Fatalf("create vcs: %v", err)
	}

	report, err := Rotate(ctx, q, next, RotateOptions{BatchSize: 1})
	if err != nil {
		t.Fatalf("Rotate() = %v", err)
	}
	counts := report[columnVariableValue]
	if counts == nil {
		t.Fatalf("Rotate() report has no %s entry", columnVariableValue)
	}
	if counts.OldKeys["k1"] < 2 {
		t.Fatalf("Rotate() old key count = %d, want at least 2", counts.OldKeys["k1"])
	}
	if counts.Plaintext < 1 {
		t.Fatalf("Rotate() plaintext count = %d, want at least 1", counts.Plaintext)
	}
	if counts.Empty < 1 {
		t.Fatalf("Rotate() empty count = %d, want at least 1", counts.Empty)
	}

	for name, id := range map[string]string{"alpha": alpha, "beta": beta} {
		value := rawVariable(t, ctx, q, id)
		assertSealedUnder(t, name, value, "k2")
		if plain, err := next.Decrypt(columnVariableValue, value); err != nil || plain != name {
			t.Fatalf("%s does not decrypt to its original value: %v", name, err)
		}
	}
	assertSealedUnder(t, "gamma", rawVariable(t, ctx, q, gamma), "k2")
	if got := rawVariable(t, ctx, q, empty); got != "" {
		t.Fatal("empty variable value was encrypted")
	}

	rawVCS, err := q.GetVCS(ctx, vcsID)
	if err != nil {
		t.Fatalf("raw GetVCS() = %v", err)
	}
	assertSealedUnder(t, "vcs client secret", *rawVCS.ClientSecret, "k2")
	assertSealedUnder(t, "vcs access token", *rawVCS.AccessToken, "k2")
	if rawVCS.RefreshToken != nil || rawVCS.PrivateKey != nil {
		t.Fatal("nil vcs secret columns did not stay nil")
	}

	verify, err := Rotate(ctx, q, next, RotateOptions{DryRun: true, BatchSize: 1})
	if err != nil {
		t.Fatalf("verification Rotate() = %v", err)
	}
	if verify.NeedsRotation() {
		t.Fatal("verification still reports plaintext or old-key values")
	}

	before := map[string]string{
		alpha: rawVariable(t, ctx, q, alpha),
		beta:  rawVariable(t, ctx, q, beta),
		gamma: rawVariable(t, ctx, q, gamma),
		empty: rawVariable(t, ctx, q, empty),
	}
	if _, err := Rotate(ctx, q, next, RotateOptions{BatchSize: 1}); err != nil {
		t.Fatalf("second Rotate() = %v", err)
	}
	for id, want := range before {
		if got := rawVariable(t, ctx, q, id); got != want {
			t.Fatalf("second Rotate() re-encrypted row %s", id)
		}
	}
}

func TestRotateDryRunWritesNothing(t *testing.T) {
	ctx, q := openRotationTest(t)
	old, err := NewCipher("k1:"+testKey(1), "")
	if err != nil {
		t.Fatalf("NewCipher(old) = %v", err)
	}
	next, err := NewCipher("k1:"+testKey(1)+",k2:"+testKey(2), "k2")
	if err != nil {
		t.Fatalf("NewCipher(next) = %v", err)
	}
	_, wsID := rotationWorkspace(t, ctx, q)
	sealed, err := old.Encrypt(columnVariableValue, "alpha")
	if err != nil {
		t.Fatalf("old Encrypt() = %v", err)
	}
	oldID := createRawVariable(t, ctx, q, wsID, "OLD", sealed)
	plainID := createRawVariable(t, ctx, q, wsID, "PLAIN", "gamma")
	beforeOld := rawVariable(t, ctx, q, oldID)
	beforePlain := rawVariable(t, ctx, q, plainID)

	report, err := Rotate(ctx, q, next, RotateOptions{DryRun: true, BatchSize: 1})
	if err != nil {
		t.Fatalf("Rotate(dry run) = %v", err)
	}
	counts := report[columnVariableValue]
	if counts == nil || counts.OldKeys["k1"] < 1 || counts.Plaintext < 1 {
		t.Fatalf("Rotate(dry run) counts = %+v, want the old-key and plaintext rows", counts)
	}
	if !report.NeedsRotation() {
		t.Fatal("Rotate(dry run) does not report pending rotation")
	}
	if got := rawVariable(t, ctx, q, oldID); got != beforeOld {
		t.Fatal("dry run changed the old-key row")
	}
	if got := rawVariable(t, ctx, q, plainID); got != beforePlain {
		t.Fatal("dry run changed the plaintext row")
	}
}

func TestRotateStopsOnUnknownKey(t *testing.T) {
	ctx, q := openRotationTest(t)
	old, err := NewCipher("k1:"+testKey(1), "")
	if err != nil {
		t.Fatalf("NewCipher(old) = %v", err)
	}
	missing, err := NewCipher("k2:"+testKey(2), "")
	if err != nil {
		t.Fatalf("NewCipher(missing) = %v", err)
	}
	_, wsID := rotationWorkspace(t, ctx, q)
	sealed, err := old.Encrypt(columnVariableValue, "top-secret")
	if err != nil {
		t.Fatalf("old Encrypt() = %v", err)
	}
	// The all-zero id sorts first, so the scan reaches this row before any other
	// keyed row the shared test database may hold.
	id := createRawVariableWithID(t, ctx, q, "00000000-0000-0000-0000-000000000001", wsID, "UNKNOWN", sealed)
	before := rawVariable(t, ctx, q, id)

	if _, err := Rotate(ctx, q, missing, RotateOptions{BatchSize: 1}); err == nil {
		t.Fatal("Rotate() with a missing key = nil error, want error")
	} else {
		if !strings.Contains(err.Error(), columnVariableValue) || !strings.Contains(err.Error(), id) {
			t.Fatalf("Rotate() error does not name the column and row: %v", err)
		}
		if strings.Contains(err.Error(), "top-secret") {
			t.Fatalf("Rotate() error leaks the value: %v", err)
		}
	}
	if got := rawVariable(t, ctx, q, id); got != before {
		t.Fatal("failed Rotate() changed the row")
	}
}
