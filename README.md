# Terrakube Go

A from-scratch Go implementation of [Terrakube](https://github.com/terrakube-io/terrakube),
the open-source Terraform/OpenTofu automation and collaboration platform.

The goal is **drop-in compatibility**: the existing Terrakube React UI, the
Terraform/OpenTofu CLI (`cloud` block and remote backend), the Terrakube
Terraform provider and the API clients should keep working while the backend is
replaced.

> Status: M2 state backend released as `v0.2.0`. M3 is in progress: discovery,
> workspaces, tags, state versions, configuration versions, runs and locking are
> implemented; the executor and run log/event streaming follow in M4.

## Requirements

- Go 1.27+
- Podman (or Docker) for Postgres, Redis and MinIO
- `terraform` and/or `tofu` for the end-to-end tests (M3+)

## Quick start

```bash
export PATH=$PATH:/usr/local/go/bin

make services-up     # Postgres, Redis, MinIO via plain podman
make db-migrate      # apply SQL migrations
make run-api         # http://localhost:8080/healthz
```

In another shell:

```bash
curl -s localhost:8080/healthz
curl -s localhost:8080/readyz
curl -s localhost:8080/version

TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/smoke.sh
TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/tfe-smoke.sh
```

The management API lives at `/api/v1` and is exercised by `scripts/smoke.sh`
(organizations, projects, workspaces and variables). The API accepts a static
bearer token from `TERRAKUBE_AUTH_TOKEN`; leave it empty to disable the check
during development. PATs and Dex OIDC replace this in M3/M5.

Build everything with `make build`; run the checks with `make check`.

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

## Roadmap

| Milestone | Scope | Status |
| --- | --- | --- |
| M0 | Repo, config, migrations, storage, queue, health, CI | done |
| M1 | JSON:API management API: org, project, workspace, variable, tag | done |
| M2 | State backend `/tfstate/v1` (executor HTTP backend) | done — [v0.2.0](https://github.com/aminbenmansour/terrakube-go/releases/tag/v0.2.0) |
| M3 | TFE v2 remote CLI workflow (`terraform init/plan/apply`) | in progress — TFE core, state/config versions and runs landed; executor-driven plan/apply and log streaming remain |
| M4 | Executor: clone, run, stream logs, store state | planned |
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
