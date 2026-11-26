# Terrakube Go

A from-scratch Go implementation of [Terrakube](https://github.com/terrakube-io/terrakube),
the open-source Terraform/OpenTofu automation and collaboration platform.

The goal is **drop-in compatibility**: the existing Terrakube React UI, the
Terraform/OpenTofu CLI (`cloud` block and remote backend), the Terrakube
Terraform provider and the API clients should keep working while the backend is
replaced.

> Status: M4 executor hardening released as `v0.4.0`. The TFE v2 CLI workflow,
> the state backend and the executor (queue reliability, VCS clone, isolation,
> durable logs and artifacts) are implemented. Upstream UI compatibility (M5)
> is next.

## Requirements

- Go 1.27+
- Podman (or Docker) for Postgres, Redis and MinIO
- `terraform` and/or `tofu` for the end-to-end tests (M3+)

## Quick start

```bash
export PATH=$PATH:/usr/local/go/bin

make services-up     # Postgres, Redis, MinIO via plain podman
make db-migrate      # apply SQL migrations

# The API and executor must share this secret: the executor mints
# workspace-scoped tokens for the /tfstate/v1 backend with it.
export TERRAKUBE_INTERNAL_TOKEN_SECRET="$(openssl rand -base64 64 | tr '+/' '-_' | tr -d '=')"
export TERRAKUBE_PUBLIC_URL=http://localhost:8080

make run-api         # http://localhost:8080/healthz
make run-executor    # http://localhost:8081/healthz
```

In another shell (with the same secret exported):

```bash
curl -s localhost:8080/healthz
curl -s localhost:8080/readyz
curl -s localhost:8080/version

TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/smoke.sh
TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/tfe-smoke.sh
TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/executor-e2e.sh
```

`scripts/executor-e2e.sh` uploads a small OpenTofu configuration and waits for
the executor to plan, apply and destroy it. It needs `tofu` (or a
`TERRAKUBE_EXECUTOR_BINARY` override) on `PATH`.

The management API lives at `/api/v1` and is exercised by `scripts/smoke.sh`
(organizations, projects, workspaces and variables). The API accepts a static
bearer token from `TERRAKUBE_AUTH_TOKEN`; leave it empty to disable the check
during development. With `TERRAKUBE_PAT_SECRET` set it also issues personal
access tokens at `/pat/v1` (the UI's Settings > Tokens screen and
`terraform login`), and with `TERRAKUBE_DEX_ENABLED=true` it accepts the Dex JWT
the UI sends.

Configuration upload tickets and log-read URLs are signed with
`TERRAKUBE_UPLOAD_TICKET_SECRET`, falling back to
`TERRAKUBE_INTERNAL_TOKEN_SECRET` and then `TERRAKUBE_AUTH_TOKEN`. Set one of
them explicitly when running more than one API replica: without a shared
secret each process invents its own key and rejects URLs minted by another.

The upstream React UI is served from the API when you point it at a build:

```bash
TERRAKUBE_UI_DIR=/path/to/terrakube/ui/build make run-api   # http://localhost:8080/
```

`index.html` is returned for deep links, `/env-config.js` is generated from
`TERRAKUBE_PUBLIC_URL`, `TERRAKUBE_DEX_CLIENT_ID` and `TERRAKUBE_DEX_ISSUER_URL`.
With `TERRAKUBE_DEX_ENABLED=true` the API validates the JWT Dex issues to the UI
against the issuer's JWKS and accepts it on the REST, GraphQL and subscription
endpoints. Audience validation is opt-in via
`TERRAKUBE_DEX_AUDIENCE_VALIDATION=true`, matching upstream. With an empty
`TERRAKUBE_AUTH_TOKEN` and Dex disabled the API accepts everything, which is the
current way to browse the UI locally.

The build is upstream's, unmodified, and it comes from
[terrakube-ui](https://github.com/aminbenmansour/terrakube-ui): a pinned pipeline
that clones `terrakube-io/terrakube` at a tag, builds it with upstream's own
Node/Yarn, and publishes an nginx image plus a tarball of the static build.
Nothing is forked, and never run `yarn build` against a patched checkout — a UI
that only works with our patches would hide exactly the API gaps this project is
about.

For the split-origin deployment, run the UI as its own container instead of
setting `TERRAKUBE_UI_DIR`:

```bash
docker compose -f deploy/compose/compose.yaml --profile app --profile ui up
```

That pulls `ghcr.io/aminbenmansour/terrakube-ui`, a private package: log in to
`ghcr.io` first, or set `TERRAKUBE_UI_IMAGE` to a local build of
[terrakube-ui](https://github.com/aminbenmansour/terrakube-ui). The image takes
its configuration from `REACT_APP_*` environment variables, so the same image
serves any stack.

`TERRAKUBE_CORS_ORIGINS` (comma-separated) restricts browser access to
`/api/v1`, `/graphql/api/v1` and `/subscriptions`; leaving it empty keeps the
permissive `Access-Control-Allow-Origin: *` default, and a UI served by this API
needs no entry at all. Set it to the UI's origin in the split-origin deployment.
Check a running stack with `./scripts/ui-smoke.sh` (needs `TERRAKUBE_UI_DIR`).

`TERRAKUBE_RATE_LIMIT` (requests per second per caller) turns on rate limiting;
it is off unless set, and `TERRAKUBE_RATE_BURST` sets the bucket size (default:
the rate). Callers are keyed by bearer token, or by address without one, and an
exceeded limit answers `429` with `Retry-After`. Keep it well above what a UI
page load and a CLI run need.

Build everything with `make build`; run the checks with `make check`, and scan
dependencies with `make vuln` (`govulncheck`; CI runs it too).

`TERRAKUBE_METRICS_TOKEN` exposes `/metrics` in Prometheus format to bearers of
that token; without it the endpoint does not exist. It reports request counts
and latency by method, route pattern and status, plus the Go runtime and process
collectors.

Security: see [SECURITY.md](SECURITY.md) to report a problem, and
[docs/threat-model.md](docs/threat-model.md) for the assets, trust boundaries,
controls and accepted risks. CI publishes a CycloneDX SBOM for each build.

If you prefer compose, `deploy/compose/compose.yaml` mirrors the podman script.
The `app` profile starts the API too:

```bash
docker compose -f deploy/compose/compose.yaml --profile app up --build
```

## Repository layout

```
cmd/api         API server: JSON:API, TFE-compatible endpoints, state backend
cmd/executor    Runs Terraform/OpenTofu jobs
cmd/registry    Private module/provider registry
cmd/migrate     Applies SQL migrations and exits

internal/config     Environment configuration
internal/db         pgx pool and migration runner
internal/storage    Local and S3 object storage
internal/queue      Redis Streams queue/log transport
internal/httpapi    Shared middleware and health endpoints
internal/jsonapi    JSON:API 1.0 codec
internal/api        /api/v1 management resources
internal/tfe        /remote/tfe/v2 Terraform/OpenTofu cloud API
internal/tfstate    /tfstate/v1 executor state backend
internal/tfoutput   /tfoutput/v1 step logs and SSE streaming
internal/executor   Job consumer: configuration, plan/apply, logs, history
internal/auth       Internal executor JWTs
internal/store      sqlc-generated queries (regenerate with `make generate`)

migrations/         Embedded SQL migrations
deploy/             Compose stack and container image
docs/compat.md      Compatibility contract and endpoint inventory
test/conformance    Tests that drive the real Terraform CLI (M3+)
```

## Configuration

All settings are environment variables; see `.env.example`. Where upstream
Terrakube uses an environment variable, the Go services accept it as an alias
(for example `DATABASE_URL` and `REDIS_URL`).

Storage defaults to the local filesystem (`TERRAKUBE_STORAGE_LOCAL_PATH`).
Set `TERRAKUBE_STORAGE_TYPE=s3` for S3/MinIO.

`TERRAKUBE_ENCRYPTION_KEYS` is a comma-separated `id:base64key` keyring
(32-byte keys) used to encrypt variable values, VCS/SSH credentials and
notification signing secrets at rest. `TERRAKUBE_ENCRYPTION_ACTIVE_KEY`
selects the key for new writes (default: the first). `make db-migrate`
encrypts existing plaintext rows; running it twice is safe. Production
refuses to start without a keyring.

To rotate keys without downtime:

1. Add the new key to `TERRAKUBE_ENCRYPTION_KEYS` on every service and point
   `TERRAKUBE_ENCRYPTION_ACTIVE_KEY` at it; restart the API and executor.
2. `DRY_RUN=1 make rotate-keys` reports what is sealed under each key.
3. `make rotate-keys` re-encrypts everything under the active key and exits
   non-zero until nothing remains under an old key. Run one rotation at a time.
4. Remove the old key from the ring and restart the services.

Back up the keyring, and never remove a key before step 3 finishes: a value
sealed with a key you removed cannot be recovered.

## Roadmap

| Milestone | Scope | Status |
| --- | --- | --- |
| M0 | Repo, config, migrations, storage, queue, health, CI | done |
| M1 | JSON:API management API: org, project, workspace, variable, tag | done |
| M2 | State backend `/tfstate/v1` (executor HTTP backend) | done — [v0.2.0](https://github.com/aminbenmansour/terrakube-go/releases/tag/v0.2.0) |
| M3 | TFE v2 remote CLI workflow (`terraform init/plan/apply`) | done — [v0.3.0](https://github.com/aminbenmansour/terrakube-go/releases/tag/v0.3.0) |
| M4 | Executor hardening: queue reliability, VCS clone, isolation, durable logs and artifacts | done — [v0.4.0](https://github.com/aminbenmansour/terrakube-go/releases/tag/v0.4.0) |
| M5 | Upstream UI compatibility, GraphQL + subscriptions | planned |
| M6 | RBAC, templates, VCS, policies, registry, parity | planned |

## Developing

- `make generate` runs `sqlc generate`. Install the pinned tool once with
  `make tools`.
- `make check` formats, vets and tests.
- `scripts/smoke.sh` exercises the management API against a running server.
- `scripts/tfe-smoke.sh` (or `make tfe-smoke`) exercises discovery, TFE v2
  organizations, workspaces, tags, state versions, configuration versions,
  runs and locking.
- `scripts/executor-e2e.sh` (or `make executor-e2e`) runs a real OpenTofu
  plan/apply/destroy through the executor and checks logs, outputs and state.
- `scripts/state-smoke.sh` exercises the state backend against a dedicated
  workspace; see the script header for the required IDs and signing key.
- `scripts/tofu-state-conformance.sh` runs a real OpenTofu HTTP backend
  `init/apply/state pull/destroy` cycle and checks that no `LOCK`/`UNLOCK`
  requests are sent.
- The upstream reference checkout used to verify compatibility lives outside
  this repo; see `docs/compat.md` for the source paths.

## Compatibility notes

- API paths, payload shapes and table names are taken from the upstream
  implementation. See `docs/compat.md` for the verified inventory.
- Upstream is Apache-2.0. Code here is written from scratch, but the API and
  schema are derived from upstream; keep attribution when distributing.
