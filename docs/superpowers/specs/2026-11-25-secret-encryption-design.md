# Secret encryption at rest — design

**Issues:** #158 (encrypt secret columns), #54 (variable values and key
rotation), #159 (rotation, follow-up).

## Intent

Variable values, VCS credentials, SSH private keys and notification signing
secrets are plaintext in Postgres today. Encryption must be transparent to the
upstream API: reads return the original value, writes store ciphertext, and a
missing or wrong key fails closed instead of returning ciphertext or garbage.

## Goals

- AES-256-GCM with a keyring; the envelope names the key that produced it.
- Encrypt on write and decrypt on read for the five secret-bearing tables.
- Backfill existing plaintext rows as part of `make db-migrate`.
- Fail closed when an encrypted value cannot be decrypted.
- No secret in logs, errors or JSON:API envelopes beyond what upstream returns.
- Keyring already accepts multiple keys so #159 is a rotation command, not a
  format change.

## Non-goals

- KMS/HSM integration. Keys come from the environment for now.
- Encrypting Terraform state, plan files, job logs or `history` blobs.
- Changing API shapes, column types or upstream table/column names.
- Rotation progress/resume UI; that is #159.

## Approaches considered

1. **Keyring AES-256-GCM (chosen).** One env keyring of `id:base64key`, an
   active key id for writes, and `enc:v1:<id>:<nonce>:<ciphertext>` in the
   column. Rotation re-encrypts rows in place later. Small, no extra tables.
2. **Single env key.** Simplest, but every rotation needs downtime and a new
   envelope version. Rejected because #159 is already scoped.
3. **Envelope encryption with wrapped data keys.** Strongest separation but
   adds a key table and a second layer for no current KMS requirement.

## Decisions

### Keyring configuration

- `TERRAKUBE_ENCRYPTION_KEYS` — comma-separated `id:base64key` entries. Each
  key must decode to exactly 32 bytes; ids are `[A-Za-z0-9_-]{1,32}` and unique.
- `TERRAKUBE_ENCRYPTION_ACTIVE_KEY` — id used for new writes. Defaults to the
  first entry. Must name an existing key.
- Empty keyring disables encryption: values stay plaintext and startup logs a
  warning. `TERRAKUBE_ENV=production` with an empty keyring refuses to start.
- A configured keyring that fails validation refuses to start.

### Envelope

```
enc:v1:<keyid>:<base64url(nonce)>:<base64url(ciphertext||tag)>
```

- 12-byte random nonce per value.
- AAD is `table.column` (`variable.variable_value`,
  `notification_configuration.signing_secret`, …), so a ciphertext moved to
  another column fails authentication.
- Empty strings are stored empty; nothing else is special-cased.
- A value without the `enc:` prefix is legacy plaintext and is returned as-is,
  so a rolling deploy does not break before the backfill runs. The backfill
  removes those rows.
- A value with the `enc:` prefix that has an unknown key id, bad base64, a
  short nonce or a failed tag returns an error. Callers already turn store
  errors into a generic `500`; no value or key material is logged.

### Columns

| Table | Columns |
| --- | --- |
| `variable` | `variable_value` |
| `globalvar` | `variable_value` |
| `ssh` | `private_key` |
| `vcs` | `client_secret`, `access_token`, `refresh_token`, `private_key` |
| `notification_configuration` | `signing_secret` |

`pat` and `cli_auth_session` already store only hashes. `destination_url` stays
plaintext because the UI displays it and upstream treats it as configuration;
documented as a residual risk.

### Store boundary

New package `internal/secretstore`:

- `Cipher` owns the keyring, `Encrypt(column, value)` and
  `Decrypt(column, value)`.
- `Wrap(store.Querier, *Cipher) store.Querier` embeds the querier and overrides
  only the methods that read or write those columns:
  - variables and global variables: list, get, create, update;
  - SSH keys: list, get, create, update;
  - VCS: list, get, create, update;
  - notification configurations: list by organization, list by workspace, get,
    create, update.
- Create/update encrypt their params and decrypt the `RETURNING *` row before
  returning it, so callers always see plaintext.
- Deletes, counts, workspace-ID lookups and the job-status wrapper are
  untouched.
- Wiring: `secretstore.Wrap(store.New(pool), cipher)` underneath
  `jobstatus.Wrap` in `cmd/api` and `cmd/executor`. The API, TFE, GraphQL,
  tfoutput and executor therefore share one decryption path.

### Backfill

- New sqlc queries list `(id, <secret columns>)` for each table and update one
  row's secret columns by id. They bypass the wrapper.
- `secretstore.Backfill(ctx, store.New(pool), cipher)` encrypts every non-empty
  value without the `enc:` prefix and is idempotent. It runs after migrations
  in `cmd/migrate`, and after the startup migration when
  `TERRAKUBE_AUTO_MIGRATE` is on.
- With no keyring the backfill is skipped and logged.
- A decryptable value under an old key is left alone; re-encrypting under the
  active key is #159's rotation job, which reuses the same queries.

## Tests

- Cipher: round-trip, key id in the envelope, wrong/missing key, malformed
  envelope, AAD mismatch, empty value, old-key decrypt with a two-key ring.
- Decorator: create encrypts the param and decrypts the returned row; get and
  list decrypt; update re-encrypts; legacy plaintext passes through; tampered
  ciphertext returns an error.
- Backfill: plaintext rows are encrypted, already-encrypted rows are untouched,
  a second run is a no-op. DB-backed, skipped without
  `TERRAKUBE_TEST_DATABASE_URL`.
- API round-trip: create a variable, read the raw column and assert the `enc:`
  prefix, read it through `/api/v1` and assert the original value.
- Failure: wrong key returns `500` with the generic message and no value.

## Compatibility and docs

- API responses keep returning plaintext, as upstream does; `docs/compat.md`
  gets a note that storage is encrypted and that ciphertext is an
  implementation detail.
- README documents the two environment variables and the backfill.
- `docs/compat.md` "Database" keeps the upstream column names unchanged.
