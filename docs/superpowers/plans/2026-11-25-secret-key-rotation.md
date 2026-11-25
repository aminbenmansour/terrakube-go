# Secret Key Rotation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Re-encrypt every encrypted secret column under the active keyring key, online and idempotently, with a dry run and a verification pass.

**Architecture:** A stateless keyset scan pages the five secret tables and re-seals anything not under the active key with the existing compare-and-set queries. A small `rotatekeys` command runs the scan, logs per-column/per-key counts, then rescans to prove nothing remains.

**Tech Stack:** Go, pgx, sqlc, PostgreSQL, slog; AES-256-GCM keyring already in `internal/secretstore`.

**Spec:** `docs/superpowers/specs/2026-11-25-secret-key-rotation-design.md`

## Global Constraints

- Upstream behavior wins; no API, table or column changes.
- Generated `internal/store/*.sql.go` is never hand-edited; run `make generate`.
- Migrations are append-only; this plan adds none.
- Never log, return or commit a secret value; errors may name the column and row id only.
- One concern per commit; generated code never shares a commit with hand-written logic.
- Never squash, never commit to `main`; use merge commits only.
- `export PATH=/usr/local/go/bin:/home/amin/go/bin:$PATH`.
- DB tests need `TERRAKUBE_TEST_DATABASE_URL='postgres://terrakube@localhost:5432/terrakube?sslmode=disable'`; Redis is not needed.

## Review Focus

- **Dropping a key too early:** an envelope under a key missing from the ring must fail the run, not silently skip the row. Task 3's `TestRotateStopsOnUnknownKey`.
- **Concurrent API writes:** a write during rotation must not be clobbered. Compare-and-set updates plus the verification pass; Task 3's second-run and verify assertions.
- **Large tables:** keyset paging by id must visit every row once. Task 3's batch-size-1 test with two rows.
- **Null vs empty:** a NULL secret column stays NULL and an empty string stays empty; neither gets an envelope. Task 3's VCS-nil and empty-variable assertions.
- **Partial failure:** a failed run leaves prior writes in place but must be safe to rerun; the error names the column and id, never the value. Task 3's unknown-key test and the command's rescan.

---

### Task 1: Paged secret scan queries

**Files:**
- Modify: `internal/store/queries/variables.sql`
- Modify: `internal/store/queries/globalvars.sql`
- Modify: `internal/store/queries/sshagents.sql`
- Modify: `internal/store/queries/vcs.sql`
- Modify: `internal/store/queries/notifications.sql`

**Interfaces:**
- Produces: five sqlc queries used by Task 3:
  - `ListVariableSecretsAfter(ctx, ListVariableSecretsAfterParams{ID, Limit}) ([]ListVariableSecretsAfterRow, error)`
  - `ListGlobalVariableSecretsAfter(...)`
  - `ListSSHKeySecretsAfter(...)`
  - `ListVCSSecretsAfter(...)`
  - `ListNotificationSecretsAfter(...)`
  Each row has `ID` plus the same secret columns as its unpaged `List*Secrets` sibling.

- [ ] **Step 1: Add the five queries next to their unpaged siblings**

```sql
-- name: ListVariableSecretsAfter :many
SELECT id, variable_value FROM variable WHERE id > $1 ORDER BY id LIMIT $2;

-- name: ListGlobalVariableSecretsAfter :many
SELECT id, variable_value FROM globalvar WHERE id > $1 ORDER BY id LIMIT $2;

-- name: ListSSHKeySecretsAfter :many
SELECT id, private_key FROM ssh WHERE id > $1 ORDER BY id LIMIT $2;

-- name: ListVCSSecretsAfter :many
SELECT id, client_secret, access_token, refresh_token, private_key
FROM vcs WHERE id > $1 ORDER BY id LIMIT $2;

-- name: ListNotificationSecretsAfter :many
SELECT id, signing_secret FROM notification_configuration
WHERE id > $1 ORDER BY id LIMIT $2;
```

- [ ] **Step 2: Verify only SQL changed**

Run: `git status --short`
Expected: exactly the five query files modified, no generated files.

- [ ] **Step 3: Commit**

```bash
git add internal/store/queries
git commit -m "feat(store): add paged secret scan queries"
```

---

### Task 2: Regenerate sqlc

**Files:**
- Modify (generated): `internal/store/variables.sql.go`, `globalvars.sql.go`, `sshagents.sql.go`, `vcs.sql.go`, `notifications.sql.go`, `querier.go` (or the files sqlc names).

**Interfaces:**
- Produces: the five generated methods from Task 1. Confirm `Limit` is `int32`; if sqlc emits another type, use that type verbatim in Task 3.

- [ ] **Step 1: Generate**

Run: `make generate`

- [ ] **Step 2: Confirm the methods and parameter types exist**

Run: `grep -n "ListVariableSecretsAfter" internal/store/*.sql.go`
Expected: the generated method plus `ListVariableSecretsAfterParams{ID string; Limit int32}`.

- [ ] **Step 3: Run checks**

Run: `make check`
Expected: exit 0.

- [ ] **Step 4: Commit generated code alone**

```bash
git add internal/store
git commit -m "chore(store): regenerate sqlc"
```

---

### Task 3: Rotation scan and re-encryption

**Files:**
- Create: `internal/secretstore/rotate.go`
- Create: `internal/secretstore/rotate_test.go`

**Interfaces:**
- Consumes: the Task 1 queries, `Cipher` (`Encrypt`, `Decrypt`, `ActiveKeyID`), the envelope constants and the `Set*Secret` queries.
- Produces:
  - `type RotateOptions struct { DryRun bool; BatchSize int }`
  - `type RotationCounts struct { Empty, Current, Plaintext int; OldKeys map[string]int }`
  - `type RotationReport map[string]*RotationCounts` with `func (r RotationReport) NeedsRotation() bool`
  - `func Rotate(ctx context.Context, q *store.Queries, c *Cipher, opts RotateOptions) (RotationReport, error)`

- [ ] **Step 1: Write the failing DB-backed tests**

Create `internal/secretstore/rotate_test.go` with the same connect/migrate/rollback scaffold as `backfill_test.go`, but begin the transaction with a repeatable-read snapshot so concurrently committed rows from other packages cannot change the counts:

```go
tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
```

Add a small fixture helper that creates an organization and workspace, then these tests:

```go
func TestRotateReencryptsOldKeyAndPlaintext(t *testing.T)
```
- old ring `k1:testKey(1)` active; new ring `k1:testKey(1),k2:testKey(2)` active `k2`.
- two variables sealed with the old cipher (`"alpha"`, `"beta"`), one legacy plaintext `"gamma"`, one empty value, all in one workspace.
- one VCS row with client secret and access token sealed under `k1`, and nil refresh token/private key.
- `Rotate(ctx, q, newCipher, RotateOptions{BatchSize: 1})` succeeds.
- `counts := report[columnVariableValue]`: `counts.OldKeys["k1"] >= 2`, `counts.Plaintext >= 1`, `counts.Empty >= 1`.
- each non-empty variable raw value has prefix `enc:v1:k2:` and `newCipher.Decrypt` returns the original; the empty one is `""`.
- the VCS client secret and token have prefix `enc:v1:k2:`; refresh token and private key are still nil.
- `verify, err := Rotate(..., RotateOptions{DryRun: true, BatchSize: 1})`; `verify.NeedsRotation()` is false.
- capture the raw bytes, run `Rotate` live again, assert every row's raw bytes are unchanged.

```go
func TestRotateDryRunWritesNothing(t *testing.T)
```
- one `k1` variable and one plaintext variable.
- capture raw bytes, `Rotate` with `DryRun: true`.
- report has `OldKeys["k1"] >= 1` and `Plaintext >= 1`; `NeedsRotation()` is true.
- raw bytes are identical to before.

```go
func TestRotateStopsOnUnknownKey(t *testing.T)
```
- a variable sealed with `k1`; rotate with a ring that only contains `k2`.
- `Rotate` returns an error naming `variable.variable_value` and the row id.
- the error does not contain the plaintext; the raw value is unchanged.

- [ ] **Step 2: Run the tests to confirm they fail to compile**

Run: `go test ./internal/secretstore -run TestRotate -v`
Expected: FAIL, `undefined: Rotate`, `undefined: RotateOptions`.

- [ ] **Step 3: Implement `rotate.go`**

Shape:

- `defaultRotationBatch = 500`; `Rotate` replaces a non-positive `BatchSize` with it.
- `Rotate` iterates the five table functions; each returns on the first error, with the partially filled report.
- A shared `rotateValue(c, column, value, opts, counts) (sealed string, changed bool, err error)`:
  - empty → `counts.Empty++`, return unchanged;
  - no `envelopePrefix` → `counts.Plaintext++`, encrypt under the active key unless dry run;
  - envelope → `c.Decrypt` first (validates every value, including current-key ones), read the key id from the envelope, `counts.Current++` when it is `c.ActiveKeyID()`, otherwise `counts.OldKeys[key]++` and re-encrypt under the active key unless dry run.
- `envelopeKey(value)` cuts the id between `envelopePrefix` and the next `:`.
- Each table function pages: call the `*After` query with `ID: last, Limit: opts.BatchSize`, stop when fewer than `Limit` rows return, otherwise set `last` to the final row's `ID`.
- Secrets are re-sealed per column so the AAD stays `table.column`; when any column of a VCS row changed, call `SetVCSSecrets` once with all four values and the read values for the compare-and-set parameters. Nil pointers stay nil and are not counted.
- Wrap per-row errors as `fmt.Errorf("rotate %s id %s: %w", column, row.ID, err)`.
- `NeedsRotation` is true when any column has `Plaintext > 0` or a positive `OldKeys` count.

- [ ] **Step 4: Run the rotation tests**

Run: `go test ./internal/secretstore -run TestRotate -v`
Expected: PASS (`TestRotate*`).

- [ ] **Step 5: Run the package and race checks**

Run: `make check && go test -race ./internal/secretstore`
Expected: exit 0.

- [ ] **Step 6: Commit**

```bash
git add internal/secretstore/rotate.go internal/secretstore/rotate_test.go
git commit -m "feat(secretstore): rotate secrets under the active key"
```

---

### Task 4: `rotatekeys` command

**Files:**
- Create: `cmd/rotatekeys/main.go`
- Modify: `Makefile`

**Interfaces:**
- Consumes: `config.Load`, `logging.New`, `db.NewPool`, `secretstore.NewCipher`, `store.New`, `secretstore.Rotate`.

- [ ] **Step 1: Add the command**

`main` defines `-dry-run` and calls `run(dryRun)`. `run`:

1. load config, build the logger, build the cipher;
2. return `errors.New("TERRAKUBE_ENCRYPTION_KEYS is not set")` when `cipher == nil`;
3. connect, `defer pool.Close()`;
4. log `active_key` and `dry_run`, then `Rotate` and log the report;
5. return the `Rotate` error after logging;
6. on a live run, run `Rotate` again with `DryRun: true`, log it as the verification, and return an error when `NeedsRotation()` is true.

Report logging sorts column names; one Info line per column with `empty`, `current`, `plaintext`, and `old_keys` rendered as space-separated sorted `id=count` values. No value is logged.

`main` prints `terrakube-rotatekeys: %v` to stderr and exits 1 on error.

- [ ] **Step 2: Add Makefile targets**

- `build: api executor registry migrate rotatekeys`
- A `rotatekeys` target that builds `$(BIN_DIR)/terrakube-rotatekeys` from `./cmd/rotatekeys`, matching the other targets.
- `rotate-keys: ## Re-encrypt secrets under the active key (DRY_RUN=1 to preview)` running `$(GO) run ./cmd/rotatekeys $(if $(DRY_RUN),-dry-run,)`.

- [ ] **Step 3: Verify the command**

Run: `make check && make rotatekeys`
Expected: exit 0 and `bin/terrakube-rotatekeys` exists.

Run without a keyring:
`env -u TERRAKUBE_ENCRYPTION_KEYS ./bin/terrakube-rotatekeys -dry-run`
Expected: exit 1 and `TERRAKUBE_ENCRYPTION_KEYS is not set`; no database writes.

Run a safe dry run against the dev database using the known test key `k1:<base64 of 32 0x01 bytes>` plus a throwaway `k2`, with `TERRAKUBE_ENCRYPTION_ACTIVE_KEY=k2`:
`DRY_RUN=1 TERRAKUBE_ENCRYPTION_KEYS='k1:...,k2:...' TERRAKUBE_ENCRYPTION_ACTIVE_KEY=k2 make rotate-keys`
Expected: exit 0, per-column lines including `old_keys=k1=<n>`, and rows unchanged (dry run).

- [ ] **Step 4: Commit**

```bash
git add cmd/rotatekeys Makefile
git commit -m "feat(cmd): add the rotate keys tool"
```

---

### Task 5: Operator documentation

**Files:**
- Modify: `README.md` (keyring paragraph)
- Modify: `docs/compat.md` (secret-storage bullet around line 668)
- Modify: `AGENTS.md` (compatibility gotcha)

- [ ] **Step 1: Document the sequence**

README: keep the keyring description, then state the rotation sequence — add the new key everywhere, point `TERRAKUBE_ENCRYPTION_ACTIVE_KEY` at it and restart the API/executor, preview with `DRY_RUN=1 make rotate-keys`, run `make rotate-keys` until it exits 0, then remove the old key and restart. State that removing a key before rotation completes makes its values unreadable and that one rotation should run at a time.

`docs/compat.md`: replace the "rotation still pending" sentence with a note that `make rotate-keys` re-encrypts under the active key, performs no schema/API change, and only needs old keys in the ring until it verifies.

`AGENTS.md`: replace the "rotation still pending" gotcha with a pointer to `make rotate-keys` and the README.

- [ ] **Step 2: Check the wording matches the command**

Run: `grep -n "rotate-keys" README.md docs/compat.md AGENTS.md`
Expected: all three files mention it, and none still says rotation is pending.

- [ ] **Step 3: Commit**

```bash
git add README.md docs/compat.md AGENTS.md
git commit -m "docs: document key rotation"
```

---

### Task 6: Graph, full verification and stacked PRs

**Files:**
- Modify (generated): `graphify-out/`

- [ ] **Step 1: Update the knowledge graph**

Run: `graphify update .`
Commit any churn separately as `chore(graph): update knowledge graph`.

- [ ] **Step 2: Full verification**

```bash
export PATH=/usr/local/go/bin:/home/amin/go/bin:$PATH
export TERRAKUBE_TEST_DATABASE_URL='postgres://terrakube@localhost:5432/terrakube?sslmode=disable'
make check
make test-race
```
Expected: exit 0 for both.

- [ ] **Step 3: Merge stack layer 1**

- branch 1: `feat/01-secret-rotation-queries` from `main`, containing Tasks 1–2;
- push, open PR 1 with base `main`, run `scripts/smoke.sh` only if the API is running (no API change, so optional);
- wait for `test`, `integration`, `executor-e2e` and GitGuardian; merge with `gh pr merge --merge` (never squash).

- [ ] **Step 4: Merge stack layer 2**

- branch 2: `feat/02-secret-key-rotation` from the branch-1 tip, containing Tasks 3–5;
- push, open PR 2 with base `feat/01-secret-rotation-queries`;
- after PR 1 merges, retarget PR 2 to `main`, wait for checks, merge with `--merge`;
- sync `main` with `git pull --ff-only`, delete both remote branches.

- [ ] **Step 5: Close the loop**

- comment on #159 with the command, the operator sequence and the verification behavior;
- comment on #54 with a link to the merged PR; keep it open only for follow-ups that remain (e.g. service-account tokens from #143).
