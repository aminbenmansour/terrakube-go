# Secret Encryption at Rest Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Encrypt variable values, VCS/SSH credentials and notification signing secrets at rest with a rotatable AES-256-GCM keyring, transparently to the API.

**Architecture:** A `Cipher` in `internal/secretstore` owns the keyring and the `enc:v1:<key>:<nonce>:<ciphertext>` envelope. `secretstore.Wrap` decorates `store.Querier` and overrides only the secret-bearing reads/writes, so every caller sees plaintext. `secretstore.Backfill` encrypts legacy rows through dedicated sqlc queries and runs after migrations.

**Tech Stack:** Go 1.27, `crypto/aes`, `crypto/cipher`, `encoding/base64`, sqlc/pgx, existing config and logging packages.

**Spec:** `docs/superpowers/specs/2026-11-25-secret-encryption-design.md`

## Global Constraints

- AES-256 keys only: each key must decode to exactly 32 bytes.
- Envelope prefix is `enc:v1:`; AAD is `table.column` with the exact strings below.
- Empty strings and nil pointers are stored unchanged.
- Encryption disabled (empty keyring) writes plaintext; `TERRAKUBE_ENV=production` with no keyring refuses to start.
- Never log or return a secret value, key or nonce. Errors name the column only.
- Do not edit `internal/store/*.sql.go`; edit `internal/store/queries/*.sql` and run `make generate`.
- Generated code is committed separately from hand-written code.
- Conventional commits; one concern per commit; no squash merges.

## Review Focus

- A plaintext user value that starts with `enc:v1:` must fail closed, not be returned as ciphertext; malformed envelopes error.
- Optional `*string` secrets: nil stays nil, empty stays empty, no encryption of empties.
- Backfill racing a live write: compare-and-set updates must not clobber a newer value.
- Unknown or retired key ids fail; a second key in the ring still decrypts its own ciphertext.
- Encrypt/decrypt errors contain the column name only, never the value or key.

---

### Task 1: Keyring configuration

**Files:**
- Modify: `internal/config/config.go`
- Test: `internal/config/config_test.go`

**Interfaces:**
- Consumes: existing `env` helpers.
- Produces: `Config.EncryptionKeys string` (raw comma-separated spec), `Config.EncryptionActiveKey string`.

- [ ] **Step 1: Write the failing test**

Add to `internal/config/config_test.go`:

```go
func TestProductionRequiresEncryptionKeys(t *testing.T) {
	t.Setenv("TERRAKUBE_ENV", "production")
	t.Setenv("TERRAKUBE_ENCRYPTION_KEYS", "")
	if _, err := Load(); err == nil {
		t.Fatal("Load() with production and no keys = nil error, want error")
	}
}

func TestDevelopmentAllowsMissingEncryptionKeys(t *testing.T) {
	t.Setenv("TERRAKUBE_ENV", "dev")
	t.Setenv("TERRAKUBE_ENCRYPTION_KEYS", "")
	if _, err := Load(); err != nil {
		t.Fatalf("Load() = %v, want nil", err)
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./internal/config -run 'TestProductionRequiresEncryptionKeys|TestDevelopmentAllowsMissingEncryptionKeys' -v`
Expected: `TestProductionRequiresEncryptionKeys` fails (nil error); the other passes.

- [ ] **Step 3: Add the fields and guard**

In `Config`, next to `PATSecret`:

```go
	// EncryptionKeys is the comma-separated id:base64key keyring for secret
	// columns. Empty disables encryption outside production.
	EncryptionKeys string
	// EncryptionActiveKey selects the key used for new ciphertext. Empty uses
	// the first key in the ring.
	EncryptionActiveKey string
```

In `Load`, after `PATSecret`:

```go
		EncryptionKeys:            env("TERRAKUBE_ENCRYPTION_KEYS", ""),
		EncryptionActiveKey:       env("TERRAKUBE_ENCRYPTION_ACTIVE_KEY", ""),
```

Before the `return cfg, nil` at the end of `Load`:

```go
	if cfg.Env == "production" && strings.TrimSpace(cfg.EncryptionKeys) == "" {
		return Config{}, fmt.Errorf("config: TERRAKUBE_ENCRYPTION_KEYS is required when TERRAKUBE_ENV=production")
	}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./internal/config`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat(config): add the secret encryption keyring settings"
```

---

### Task 2: Cipher and envelope

**Files:**
- Create: `internal/secretstore/cipher.go`
- Test: `internal/secretstore/cipher_test.go`

**Interfaces:**
- Consumes: `Config.EncryptionKeys`, `Config.EncryptionActiveKey`.
- Produces:
  - `func NewCipher(spec, active string) (*Cipher, error)` — `nil, nil` when `spec` is empty.
  - `func (c *Cipher) Encrypt(column, plaintext string) (string, error)`
  - `func (c *Cipher) Decrypt(column, ciphertext string) (string, error)`
  - `func (c *Cipher) ActiveKeyID() string`

- [ ] **Step 1: Write the failing test**

Create `internal/secretstore/cipher_test.go`:

```go
package secretstore

import (
	"encoding/base64"
	"strings"
	"testing"
)

func testKey(seed byte) string {
	key := make([]byte, 32)
	for i := range key {
		key[i] = seed
	}
	return base64.StdEncoding.EncodeToString(key)
}

func TestCipherRoundTrip(t *testing.T) {
	c, err := NewCipher("k1:"+testKey(1), "")
	if err != nil {
		t.Fatalf("NewCipher() = %v", err)
	}
	sealed, err := c.Encrypt("variable.variable_value", "hello")
	if err != nil {
		t.Fatalf("Encrypt() = %v", err)
	}
	if !strings.HasPrefix(sealed, "enc:v1:k1:") {
		t.Fatalf("envelope = %q, want enc:v1:k1: prefix", sealed)
	}
	got, err := c.Decrypt("variable.variable_value", sealed)
	if err != nil {
		t.Fatalf("Decrypt() = %v", err)
	}
	if got != "hello" {
		t.Fatalf("Decrypt() = %q, want hello", got)
	}
}

func TestCipherRejectsWrongColumnAndKey(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	sealed, _ := c.Encrypt("variable.variable_value", "hello")
	if _, err := c.Decrypt("globalvar.variable_value", sealed); err == nil {
		t.Fatal("Decrypt() with another column = nil error, want error")
	} else if strings.Contains(err.Error(), "hello") {
		t.Fatalf("Decrypt() error leaks the value: %v", err)
	}
	other, _ := NewCipher("k2:"+testKey(2), "")
	if _, err := other.Decrypt("variable.variable_value", sealed); err == nil {
		t.Fatal("Decrypt() with another key = nil error, want error")
	}
}

func TestCipherEmptyAndPlaintext(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	if got, err := c.Encrypt("ssh.private_key", ""); err != nil || got != "" {
		t.Fatalf("Encrypt(empty) = %q, %v; want empty, nil", got, err)
	}
	if got, err := c.Decrypt("ssh.private_key", "plain legacy"); err != nil || got != "plain legacy" {
		t.Fatalf("Decrypt(plaintext) = %q, %v; want passthrough", got, err)
	}
	if _, err := c.Decrypt("ssh.private_key", "enc:v1:k1:notbase64:notbase64"); err == nil {
		t.Fatal("Decrypt(malformed) = nil error, want error")
	}
	if _, err := c.Decrypt("ssh.private_key", "enc:v1:unknown:AAAA:AAAA"); err == nil {
		t.Fatal("Decrypt(unknown key) = nil error, want error")
	}
}

func TestCipherKeyring(t *testing.T) {
	c, err := NewCipher("old:"+testKey(1)+",new:"+testKey(2), "new")
	if err != nil {
		t.Fatalf("NewCipher() = %v", err)
	}
	if c.ActiveKeyID() != "new" {
		t.Fatalf("ActiveKeyID() = %q, want new", c.ActiveKeyID())
	}
	old, _ := NewCipher("old:"+testKey(1), "")
	sealed, _ := old.Encrypt("vcs.client_secret", "s3cret")
	if got, err := c.Decrypt("vcs.client_secret", sealed); err != nil || got != "s3cret" {
		t.Fatalf("Decrypt(old key) = %q, %v", got, err)
	}
}

func TestCipherRejectsBadConfiguration(t *testing.T) {
	if c, err := NewCipher("", ""); c != nil || err != nil {
		t.Fatalf("NewCipher(empty) = %v, %v; want nil, nil", c, err)
	}
	cases := []string{
		"nokey",
		"k1:short",
		"k1:" + testKey(1) + ",k1:" + testKey(2),
		"k1:" + testKey(1) + ",bad id:" + testKey(2),
	}
	for _, spec := range cases {
		if _, err := NewCipher(spec, ""); err == nil {
			t.Errorf("NewCipher(%q) = nil error, want error", spec)
		}
	}
	if _, err := NewCipher("k1:"+testKey(1), "missing"); err == nil {
		t.Fatal("NewCipher(active missing) = nil error, want error")
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./internal/secretstore -v`
Expected: build failure, `NewCipher` undefined.

- [ ] **Step 3: Implement the cipher**

Create `internal/secretstore/cipher.go` with:
- `const envelopePrefix = "enc:v1:"`.
- `type Cipher struct { aeads map[string]cipher.AEAD; active string }`.
- `NewCipher`: empty spec → `nil, nil`; split on `,`; split each entry on the first `:`; validate id with a small loop allowing `A-Za-z0-9_-` and length 1–32; decode with `base64.StdEncoding` and require 32 bytes; `cipher.NewGCM`; reject duplicate ids; resolve `active` (default first entry); unknown active is an error.
- `Encrypt`: empty plaintext returns `""`; 12-byte `rand.Read` nonce; `Seal(nil, nonce, []byte(plaintext), []byte(column))`; return `envelopePrefix + active + ":" + RawURLEncoding(nonce) + ":" + RawURLEncoding(sealed)`.
- `Decrypt`: plaintext (no prefix) returns unchanged; `strings.SplitN(strings.TrimPrefix(value, envelopePrefix), ":", 3)`; unknown id, bad base64, nonce length ≠ 12 and `Open` failure all return `fmt.Errorf("secretstore: decrypt %s: %w", column, err)` with no value.

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./internal/secretstore -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/secretstore/cipher.go internal/secretstore/cipher_test.go
git commit -m "feat(secretstore): add the AES-256-GCM keyring"
```

---

### Task 3: Transparent store decorator

**Files:**
- Create: `internal/secretstore/store.go`
- Test: `internal/secretstore/store_test.go`

**Interfaces:**
- Consumes: `*Cipher`, `store.Querier`.
- Produces: `func Wrap(q store.Querier, c *Cipher) store.Querier` — returns `q` unchanged when `c == nil`.

Column constants (exact strings, used as AAD):

```go
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
```

- [ ] **Step 1: Write the failing test**

Create `internal/secretstore/store_test.go` with a `fakeQuerier` embedding `store.Querier` and recording params/returning one row. Cover:
- `CreateVariable` stores a ciphertext param and returns the decrypted row.
- `GetVariable` and `ListVariablesByWorkspace` decrypt.
- `UpdateVariable` re-encrypts.
- `CreateSSHKey`/`GetSSHKey`, `CreateVCS`/`GetVCS`, `CreateNotificationConfiguration`/`GetNotificationConfiguration` do the same for their columns.
- nil secret pointers stay nil.
- `Wrap(q, nil)` returns the same querier and does not encrypt.

Use assertions like:

```go
func TestWrapCreateVariableEncryptsAndDecrypts(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	q := &fakeQuerier{}
	wrapped := Wrap(q, c)
	row, err := wrapped.CreateVariable(context.Background(), store.CreateVariableParams{
		ID: "v1", VariableValue: "plain",
	})
	if err != nil {
		t.Fatalf("CreateVariable() = %v", err)
	}
	if !strings.HasPrefix(q.createdVariable.VariableValue, "enc:v1:") {
		t.Fatalf("stored value = %q, want ciphertext", q.createdVariable.VariableValue)
	}
	if row.VariableValue != "plain" {
		t.Fatalf("returned value = %q, want plain", row.VariableValue)
	}
}
```

The fake's create methods return the params as the row; get/list return a row with a ciphertext produced by `c.Encrypt`.

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./internal/secretstore -run TestWrap -v`
Expected: build failure, `Wrap` undefined.

- [ ] **Step 3: Implement the decorator**

Create `internal/secretstore/store.go`:

```go
type encryptedStore struct {
	store.Querier
	cipher *Cipher
}

func Wrap(q store.Querier, c *Cipher) store.Querier {
	if c == nil {
		return q
	}
	return &encryptedStore{Querier: q, cipher: c}
}
```

Override exactly these methods, each encrypting params before delegating and decrypting the returned row(s):

- `ListVariablesByWorkspace`, `GetVariable`, `CreateVariable`, `UpdateVariable`
- `ListGlobalVariablesByOrganization`, `GetGlobalVariable`, `CreateGlobalVariable`, `UpdateGlobalVariable`
- `ListSSHKeysByOrganization`, `GetSSHKey`, `CreateSSHKey`, `UpdateSSHKey`
- `ListVCSByOrganization`, `GetVCS`, `CreateVCS`, `UpdateVCS` (`RefreshToken` is read-only; decrypt it, never encrypt a write because the param has no field)
- `ListNotificationConfigurationsByOrganization`, `ListNotificationConfigurationsByWorkspace`, `GetNotificationConfiguration`, `CreateNotificationConfiguration`, `UpdateNotificationConfiguration`

Helpers:

```go
func (s *encryptedStore) encryptPtr(column string, v *string) (*string, error)
func (s *encryptedStore) decryptPtr(column string, v *string) (*string, error)
```

Return the wrapped store's error unchanged; do not include values in errors.

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./internal/secretstore -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/secretstore/store.go internal/secretstore/store_test.go
git commit -m "feat(secretstore): encrypt secret columns through a querier wrapper"
```

---

### Task 4: Backfill queries

**Files:**
- Modify: `internal/store/queries/variables.sql`, `globalvars.sql`, `sshagents.sql`, `vcs.sql`, `notifications.sql`
- Generated: `internal/store/*.sql.go`, `internal/store/querier.go`

**Interfaces:**
- Produces sqlc methods `ListVariableSecrets`, `SetVariableSecret`, `ListGlobalVariableSecrets`, `SetGlobalVariableSecret`, `ListSSHKeySecrets`, `SetSSHKeySecret`, `ListVCSSecrets`, `SetVCSSecrets`, `ListNotificationSecrets`, `SetNotificationSecret`.

- [ ] **Step 1: Add the queries**

Append to each file:

```sql
-- name: ListVariableSecrets :many
SELECT id, variable_value FROM variable;

-- name: SetVariableSecret :exec
UPDATE variable SET variable_value = $2 WHERE id = $1 AND variable_value = $3;
```

Use the same shape for `globalvar` (`variable_value`), `ssh` (`private_key`) and
`notification_configuration` (`signing_secret`).

For `vcs`:

```sql
-- name: ListVCSSecrets :many
SELECT id, client_secret, access_token, refresh_token, private_key FROM vcs;

-- name: SetVCSSecrets :exec
UPDATE vcs
SET client_secret = $2, access_token = $3, refresh_token = $4, private_key = $5
WHERE id = $1
  AND client_secret IS NOT DISTINCT FROM $6
  AND access_token IS NOT DISTINCT FROM $7
  AND refresh_token IS NOT DISTINCT FROM $8
  AND private_key IS NOT DISTINCT FROM $9;
```

- [ ] **Step 2: Commit the hand-written queries**

```bash
git add internal/store/queries
git commit -m "feat(store): add secret backfill queries"
```

- [ ] **Step 3: Regenerate sqlc**

Run: `make generate`
Expected: new methods in `internal/store/*.sql.go` and `internal/store/querier.go`; no other diff.

- [ ] **Step 4: Build**

Run: `go build ./...`
Expected: success. sqlc names the second VCS set of params with a `_2` suffix; use whatever it emits.

- [ ] **Step 5: Commit the generated code alone**

```bash
git add internal/store
git commit -m "chore(store): regenerate sqlc"
```

---

### Task 5: Backfill implementation

**Files:**
- Create: `internal/secretstore/backfill.go`
- Test: `internal/secretstore/backfill_test.go`

**Interfaces:**
- Consumes: `*store.Queries`, `*Cipher`.
- Produces: `func Backfill(ctx context.Context, q *store.Queries, c *Cipher) (int, error)` — returns the number of rows changed; `0, nil` when `c == nil`.

- [ ] **Step 1: Write the failing unit test**

Create `internal/secretstore/backfill_test.go` with a test for the pure decision helper (implement it as `needsEncryption(c *Cipher, value string) bool`):

```go
func TestNeedsEncryption(t *testing.T) {
	c, _ := NewCipher("k1:"+testKey(1), "")
	if needsEncryption(c, "") {
		t.Error("empty value needs encryption")
	}
	if needsEncryption(c, "plain") == false {
		t.Error("plaintext does not need encryption")
	}
	if needsEncryption(c, "enc:v1:k1:AAAA:AAAA") {
		t.Error("enveloped value needs encryption")
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./internal/secretstore -run TestNeedsEncryption -v`
Expected: build failure, `needsEncryption` undefined.

- [ ] **Step 3: Implement the backfill**

Create `internal/secretstore/backfill.go`. For each table:

1. List rows with the `List*Secrets` query.
2. For each row, encrypt every non-nil, non-empty, non-enveloped column with `c.Encrypt(column, value)`.
3. If nothing changed, skip the row.
4. Call the matching `Set*Secret` with the original values as the compare-and-set arguments.
5. Count changed rows.

`needsEncryption` returns false when `c == nil`, the value is empty, or it has `envelopePrefix`.

- [ ] **Step 4: Run the unit test**

Run: `go test ./internal/secretstore -run TestNeedsEncryption -v`
Expected: PASS.

- [ ] **Step 5: Add the DB-backed test**

Append to `backfill_test.go`, skipped unless `TERRAKUBE_TEST_DATABASE_URL` is set:

```go
func TestBackfillEncryptsPlaintextRows(t *testing.T) {
	dsn := os.Getenv("TERRAKUBE_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TERRAKUBE_TEST_DATABASE_URL is not set")
	}
	// pool + db.Migrate(ctx, pool, migrations.Files, ".")
	// insert an organization and workspace with unique UUIDs via store.New(pool)
	// create a variable with value "plain" (raw queries, no wrapper)
	// Backfill -> 1 row changed; raw variable_value has the enc: prefix
	// Backfill again -> 0 rows changed
	// delete the rows in t.Cleanup
}
```

Add a second case for a VCS row with two plaintext secret columns to prove the multi-column compare-and-set, and a case where the row is updated between list and set is not required here (the CAS SQL is covered by the no-op second run).

- [ ] **Step 6: Run the tests**

Run: `TERRAKUBE_TEST_DATABASE_URL='postgres://terrakube@localhost:5432/terrakube?sslmode=disable' go test ./internal/secretstore -v`
Expected: PASS, including the DB test.

- [ ] **Step 7: Commit**

```bash
git add internal/secretstore/backfill.go internal/secretstore/backfill_test.go
git commit -m "feat(secretstore): backfill legacy plaintext secrets"
```

---

### Task 6: Wiring and round-trip test

**Files:**
- Modify: `cmd/api/main.go`, `cmd/executor/main.go`, `cmd/migrate/main.go`
- Test: `internal/secretstore/roundtrip_test.go`

**Interfaces:**
- Consumes: `NewCipher`, `Wrap`, `Backfill`.
- Produces: the API and executor read/write through the wrapper; migrate backfills.

- [ ] **Step 1: Write the round-trip test**

Create `internal/secretstore/roundtrip_test.go`, skipped without `TERRAKUBE_TEST_DATABASE_URL`:

```go
func TestWrappedStoreRoundTrip(t *testing.T) {
	// pool, migrations, unique org/workspace
	// cipher := NewCipher("k1:"+testKey(1), "")
	// wrapped := Wrap(store.New(pool), cipher)
	// wrapped.CreateVariable(... VariableValue: "original")
	// raw := store.New(pool).GetVariable(...) -> strings.HasPrefix(raw.VariableValue, "enc:v1:")
	// wrapped.GetVariable(...) -> "original"
}
```

- [ ] **Step 2: Run the test**

Run: `TERRAKUBE_TEST_DATABASE_URL='...' go test ./internal/secretstore -run TestWrappedStoreRoundTrip -v`
Expected: PASS. The test calls `Wrap` directly, so it proves the decorator against Postgres before the binaries use it.

- [ ] **Step 3: Wire the binaries**

In `cmd/api/main.go`, after config load:

```go
	cipher, err := secretstore.NewCipher(cfg.EncryptionKeys, cfg.EncryptionActiveKey)
	if err != nil {
		return fmt.Errorf("configure secret encryption: %w", err)
	}
	if cipher == nil {
		logger.Warn("secret encryption is disabled: set TERRAKUBE_ENCRYPTION_KEYS")
	}
```

Inside the `if cfg.AutoMigrate` block, after `db.Migrate`:

```go
		if cipher != nil {
			changed, err := secretstore.Backfill(ctx, store.New(pool), cipher)
			if err != nil {
				return fmt.Errorf("backfill secret encryption: %w", err)
			}
			if changed > 0 {
				logger.Info("encrypted plaintext secrets", "rows", changed)
			}
		}
```

Replace the store construction:

```go
	querier := jobstatus.Wrap(secretstore.Wrap(store.New(pool), cipher), events.Publish)
```

In `cmd/executor/main.go`, build the cipher the same way (no backfill) and use:

```go
		Store: jobstatus.Wrap(secretstore.Wrap(store.New(pool), cipher), events.Publish),
```

In `cmd/migrate/main.go`, build the cipher, call `Backfill` after `db.Migrate`, and log the row count or the disabled warning.

- [ ] **Step 4: Run the full checks**

Run: `make check`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add cmd/api/main.go cmd/executor/main.go cmd/migrate/main.go internal/secretstore/roundtrip_test.go
git commit -m "feat: wire secret encryption into the api, executor and migrate"
```

---

### Task 7: Documentation and graph

**Files:**
- Modify: `README.md`, `docs/compat.md`
- Generated: `graphify-out/`

**Interfaces:**
- Consumes: the env names and behavior above.
- Produces: operator documentation.

- [ ] **Step 1: Document the environment**

Add to the README configuration section:

```markdown
`TERRAKUBE_ENCRYPTION_KEYS` is a comma-separated `id:base64key` keyring
(32-byte keys) used to encrypt variable values, VCS/SSH credentials and
notification signing secrets at rest. `TERRAKUBE_ENCRYPTION_ACTIVE_KEY`
selects the key for new writes (default: the first). `make db-migrate`
encrypts existing plaintext rows; running it twice is safe. Production
refuses to start without a keyring.
```

- [ ] **Step 2: Document compatibility**

In `docs/compat.md`, under the database/limitations section, add:

```markdown
Secret columns are encrypted at rest with a versioned AES-256-GCM envelope
(`enc:v1:<key>:<nonce>:<ciphertext>`) and a per-column AAD. The JSON:API and
GraphQL responses still return plaintext, so the upstream UI and clients see
no difference. The keyring is configured with `TERRAKUBE_ENCRYPTION_KEYS`;
`make db-migrate` backfills legacy rows. `notification_configuration.destination_url`
remains plaintext because upstream treats it as configuration.
```

- [ ] **Step 3: Update the graph**

Run: `graphify update .`
Expected: `graphify-out/` changes only.

- [ ] **Step 4: Commit**

```bash
git add README.md docs/compat.md
git commit -m "docs: document secret encryption"
git add graphify-out
git commit -m "chore(graph): update knowledge graph"
```

---

### Task 8: Verification

- [ ] **Step 1: Run the full suite with the DB and race detector**

Run:

```bash
export PATH=/usr/local/go/bin:$PATH
export TERRAKUBE_TEST_DATABASE_URL='postgres://terrakube@localhost:5432/terrakube?sslmode=disable'
export TERRAKUBE_TEST_REDIS_URL='redis://localhost:6379/9'
make check
make test-race
```

Expected: PASS.

- [ ] **Step 2: Prove the API round-trip manually**

Start the API with a generated keyring, create a variable, and inspect the row:

```bash
key=$(head -c 32 /dev/urandom | base64 -w0)
TERRAKUBE_ENCRYPTION_KEYS="dev:$key" TERRAKUBE_AUTH_TOKEN=dev-token make run-api
```

Then create a variable with `POST /api/v1/workspace/{id}/variable`, read it back
through the API (original value), and query Postgres for `variable_value`
(starts with `enc:v1:dev:`).

- [ ] **Step 3: Confirm no plaintext remains**

```bash
psql 'postgres://terrakube@localhost:5432/terrakube?sslmode=disable' -c \
  "SELECT count(*) FROM variable WHERE variable_value <> '' AND variable_value NOT LIKE 'enc:v1:%';"
```

Expected: `0` after the backfill for rows created under encryption. Rows created
before the keyring was configured are encrypted by the next `make db-migrate`.

## Self-Review Notes

- Spec coverage: keyring config (Task 1), envelope and fail-closed (Task 2),
  transparent reads/writes (Task 3), backfill queries and command (Tasks 4–5),
  binary wiring and round-trip (Task 6), docs (Task 7), verification (Task 8).
- Types: `NewCipher`, `Wrap`, `Backfill`, `needsEncryption` are used with the
  same signatures in every task.
- Review Focus tests: prefix collision and malformed envelope (Task 2), nil/empty
  (Task 3), CAS/no-op (Task 5), old/unknown keys (Task 2), error text without
  values (Task 2/3 assertions).
