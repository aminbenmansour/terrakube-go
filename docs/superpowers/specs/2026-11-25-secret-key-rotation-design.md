# Secret key rotation — design

**Issues:** #159 (rotate keys, retire old ones), #54 (encryption and key
management), #158 (landed in #176).

## Intent

The keyring from #158 already names the key in every envelope and accepts
several keys with one active write key. Rotation is therefore a re-encrypt
job: read every secret value, seal anything not under the active key with the
active key, and prove that nothing remains under an old key. Operators can run
it while the API serves traffic, stop it at any point, and rerun it safely.

## Goals

- Re-encrypt every secret column under `TERRAKUBE_ENCRYPTION_ACTIVE_KEY`.
- Encrypt legacy plaintext too, so a skipped backfill does not block rotation.
- Run online: no maintenance window, no schema change, no envelope change.
- Stop and rerun safely; a rerun re-scans and skips values already current.
- Report counts per column and key, including what remains after a run.
- Fail closed on an unknown key or a corrupt value, naming only the column and
  row id, never the value.
- Dry run that reads and validates every value without writing.

## Non-goals

- KMS/HSM, key derivation, or a new envelope version.
- A rotation checkpoint table, progress UI, or resumable cursor. The job is
  idempotent, so a rescan is the cursor.
- Rotating Terraform state, plan files, logs or `history` blobs.
- Automatic rotation at startup. Rotation is an operator action.
- Changing API responses or upstream table/column names.

## Approaches considered

1. **Stateless rescan (chosen).** Page each secret table by primary key, seal
   anything not under the active key, then rescan to verify. No migration, no
   bookkeeping to corrupt, and a rerun always converges.
2. **Durable checkpoint table.** Records the last id per table so a rerun
   resumes exactly. Auditable, but adds a migration and state for a job whose
   writes are already idempotent.
3. **Stateless now, checkpoint later.** Same as 1; revisit only if real table
   sizes make a full rescan expensive.

## Decisions

### Command

New `cmd/rotatekeys` binary (`terrakube-rotatekeys`), built like the other
commands. It parses the normal configuration, builds the keyring, connects to
Postgres and calls `secretstore.Rotate`. Flags:

- `-dry-run` — read, validate and report; write nothing.

`make rotate-keys` runs it from source; `DRY_RUN=1 make rotate-keys` previews.
It is not mounted in the API and does not run during `db-migrate`.

### Scan

Five new sqlc queries page the same `(id, secret columns)` projections the
backfill uses:

```sql
SELECT id, variable_value FROM variable WHERE id > $1 ORDER BY id LIMIT $2;
```

(and the equivalents for `globalvar`, `ssh`, `vcs` and
`notification_configuration`). Ids are `varchar(36)`, so keyset paging by id is
stable; a batch never depends on `OFFSET`. Batch size defaults to 500 and is an
option so tests can force several pages.

### Per-value action

| Stored value | Action |
| --- | --- |
| empty | skip |
| plaintext | encrypt under the active key |
| `enc:v1:<active>:…` | decrypt to validate, leave as is |
| `enc:v1:<other>:…` | decrypt, re-encrypt under the active key |
| malformed / missing key / bad tag | stop with an error |

Every non-empty value is decrypted once, including ones already under the
active key, so the command validates the whole table rather than only the rows
it changes. Plaintext values encrypt without a decrypt step.

Writes reuse the existing `Set*Secret` queries: compare-and-set on the values
that were read. A concurrent API write under the active key makes the update
match nothing, which is harmless; a rerun or the verification pass covers it.
Rows with several secret columns (`vcs`) update all four in one statement,
keeping unchanged columns at their current bytes.

### Report and verification

`Rotate` returns counts per `table.column`:

- `empty`, `current` (active key), `plaintext`, and `old_keys` (source key id
  to count).

In a dry run the counts describe what would change. In a live run they describe
what was seen and rewritten. After a live run the command performs one more
dry-run scan; if any plaintext or old-key value remains, it logs the remaining
counts and exits non-zero. Decrypt failures stop either pass with the column and
row id, never the value or key material.

### Operator sequence

1. Add the new key to `TERRAKUBE_ENCRYPTION_KEYS` everywhere and point
   `TERRAKUBE_ENCRYPTION_ACTIVE_KEY` at it; restart the API and executor.
2. `DRY_RUN=1 make rotate-keys` and review the per-key counts.
3. `make rotate-keys`; it exits non-zero until nothing remains under an old key.
4. Remove the old key from the ring and restart the services.

Removing a key before step 3 makes its values unreadable; the command and
README state this. Concurrent rotations are harmless (idempotent writes with
compare-and-set) but pointless, so the docs say to run one at a time.

## Tests

- DB-backed, inside a transaction like the backfill test:
  - a value sealed under `k1` is re-encrypted under `k2`, the raw column names
    `k2`, and both keys still decrypt existing rows while the ring holds them;
  - a second run leaves the row untouched;
  - dry run reports the old key but changes no bytes;
  - a legacy plaintext row is encrypted under the active key;
  - batch size 1 with two rows exercises paging;
  - a value under a key missing from the ring returns an error and writes
    nothing;
  - verification counts reach zero after rotation.
- Cipher tests already cover multi-key rings and decrypt failures.

## Compatibility and docs

- No API, column or table change; `docs/compat.md` records the operational
  command and the rotation sequence.
- README gains the sequence and the warning about dropping keys early.
- `AGENTS.md` drops the "rotation pending" note once merged.
