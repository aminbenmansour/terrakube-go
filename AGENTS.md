# Terrakube Go — agent guide

Terrakube Go is a from-scratch Go rewrite of [Terrakube](https://github.com/terrakube-io/terrakube),
the open-source Terraform/OpenTofu automation platform. The target is **drop-in
compatibility** with the upstream Java implementation: the React UI, the
Terraform/OpenTofu CLI (`cloud` block and remote backend), the Terraform
provider and existing API clients must keep working.

## Non-negotiables

- **Upstream is the specification.** Before changing an endpoint, payload,
  status code, header, table or column, check `docs/compat.md` and, when
  available, the upstream source. Do not "improve" the API shape.
- **JSON:API 1.0.** Use `application/vnd.api+json`, return the `errors` array
  on failure, and send no body or content type on `DELETE`.
- **Schema names mirror upstream.** Migrations are append-only: add a new
  numbered file instead of editing an applied one.
- **Generated code is off limits.** `internal/store/*.sql.go` is produced by
  sqlc. Change `internal/store/queries/*.sql`, then run `make generate`.
- **Secrets stay out of logs and commits.** `TERRAKUBE_AUTH_TOKEN` is a
  development convenience; PATs and Dex replace it later. Variable values are
  currently plaintext — never expose them in logs or error messages.
- **Keep attribution.** Upstream is Apache-2.0; the API and schema are derived
  from it. Preserve `LICENSE` and the compatibility notes.

## Commands

If `go` is not on `PATH` on the development machine:

```bash
export PATH=$PATH:/usr/local/go/bin
```

| Task | Command |
| --- | --- |
| Start Postgres, Redis and MinIO | `make services-up` |
| Stop / status / wipe services | `make services-down` / `make services-status` / `make services-wipe` |
| Apply migrations | `make db-migrate` |
| Run the API | `make run-api` |
| Format, vet and test | `make check` |
| Race tests | `make test-race` |
| Build all binaries | `make build` |
| Regenerate sqlc | `make generate` |
| Install pinned sqlc | `make tools` |
| API smoke test | `make smoke` (API must be running; requires `jq`) |
| TFE v2 smoke test | `make tfe-smoke` (API must be running; requires `jq`) |
| State backend smoke test | `make state-smoke` (API must be running; see script header) |

`scripts/dev-services.sh` uses plain podman. There is no `podman compose`
provider on the development machine; use the script or Docker Compose. The
MinIO image is `docker.io/pgsty/minio` because the official image is not
available in this environment. `.data/` holds local state and is git-ignored.

## Layout

- `cmd/api` — API server: JSON:API, TFE-compatible endpoints, state backend
- `cmd/executor` — Terraform/OpenTofu job runner (M4)
- `cmd/registry` — private module/provider registry (M6)
- `cmd/migrate` — applies embedded migrations and exits
- `internal/api` — `/api/v1` management resources
- `internal/tfe` — `/remote/tfe/v2` Terraform/OpenTofu cloud API
- `internal/auth` — internal JWT minting and validation for executor calls
- `internal/tfstate` — `/tfstate/v1` executor state backend
- `internal/store` — sqlc-generated queries (do not edit)
- `internal/store/queries` — hand-written SQL for sqlc
- `internal/db` — pgx pool and migration runner
- `internal/httpapi` — shared middleware and health endpoints
- `internal/jsonapi` — JSON:API codec
- `internal/storage` — local and S3 object storage
- `internal/queue` — Redis Streams transport
- `internal/config`, `internal/logging`, `internal/version`
- `migrations` — embedded, append-only SQL
- `docs/compat.md` — verified upstream contract
- `scripts` — development and smoke-test helpers

## Working rules

1. Read `docs/compat.md` before touching the API or schema.
2. Prefer small, focused commits with conventional messages (`feat:`, `fix:`,
   `docs:`, `chore:`, `test:`).
3. Keep generated-code changes in their own commit when practical.
4. Add or update tests for behavior changes.
5. Run verification before committing:
   - always: `make check` (or at least `go build ./... && go vet ./... && go test ./...`)
   - API changes: start the stack and run `scripts/smoke.sh`
   - migration changes: test against a clean database
6. Update `docs/compat.md` when an observed upstream contract changes.

## Verification details

- `make check` runs `gofmt -w`, `go vet ./...` and `go test ./...`.
- CI also runs `go test -race ./...` and `go build ./...`; run
  `make test-race` for concurrency-sensitive changes.
- Smoke test: `TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/smoke.sh`.
  To target another port, set `TERRAKUBE_API_URL=http://localhost:18080`.
- Never run `make services-wipe`, delete `.data/`, or drop a database unless
  the user explicitly asks; it destroys local state.

## Compatibility gotchas

- Discovery advertises `state.v2` at `/remote/state/v2/`, but the upstream
  state handler is mounted at `/tfstate/v1`. Verify with the real CLI before
  changing either path.
- The upstream `/tfstate/v1` handler has no `LOCK`/`UNLOCK` methods. The
  executor configures a plain HTTP backend address; workspace locking is the
  TFE action `POST /workspaces/{id}/actions/lock`. Do not invent a state-lock
  endpoint.
- The upstream UI uses GraphQL at `/graphql/api/v1` (plus graphql-ws
  `/subscriptions`) for most screens; only the organization list uses REST.
  M5 needs a GraphQL implementation, not just a WebSocket.
- Variable values are stored plaintext today; encryption is required before
  exposing the API beyond local development.
- M1 filtering, sorting and paging are in memory; move them into SQL before
  production.

## Git and PRs

- Work in a feature branch; never commit directly to `main`.
- This repository uses stacked PRs: each layer builds on the previous branch.
  See `.opencode/skills/stacked-prs/SKILL.md`.
- Do not push, force-push, rewrite shared history, or open/merge PRs without
  explicit approval.
- If a request is ambiguous or the change is cross-cutting, plan first.
