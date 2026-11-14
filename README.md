# Terrakube Go

A from-scratch Go implementation of [Terrakube](https://github.com/terrakube-io/terrakube),
the open-source Terraform/OpenTofu automation and collaboration platform.

The goal is **drop-in compatibility**: the existing Terrakube React UI, the
Terraform/OpenTofu CLI (`cloud` block and remote backend), the Terrakube
Terraform provider and the API clients should keep working while the backend is
replaced.

> Status: foundation only. The service binaries start and health endpoints
> respond. The database, management API and Terraform workflow land in later
> pull requests.

## Requirements

- Go 1.27+
- Podman (or Docker) for Postgres, Redis and MinIO

## Quick start

```bash
export PATH=$PATH:/usr/local/go/bin

make services-up     # Postgres, Redis, MinIO via plain podman
make run-api         # http://localhost:8080/healthz
```

In another shell:

```bash
curl -s localhost:8080/healthz
curl -s localhost:8080/readyz
curl -s localhost:8080/version
```

Build everything with `make build`; run the checks with `make check`.

If you prefer compose, `deploy/compose/compose.yaml` mirrors the podman script.
The `app` profile starts the API too:

```bash
docker compose -f deploy/compose/compose.yaml --profile app up --build
```

## Repository layout

```
cmd/api         API server
cmd/executor    Runs Terraform/OpenTofu jobs
cmd/registry    Private module/provider registry

internal/config     Environment configuration
internal/httpapi    Shared middleware and health endpoints
internal/logging    Structured logging
internal/storage    Object-store contract
internal/version    Build metadata

deploy/             Compose stack and container image
scripts/            Local development helpers
```

## Configuration

All settings are environment variables; see `.env.example`. Where upstream
Terrakube uses an environment variable, the Go services accept it as an alias
(for example `DATABASE_URL` and `REDIS_URL`).

## Roadmap

| Milestone | Scope | Status |
| --- | --- | --- |
| M0 | Repo, config, storage, queue, health, CI | in progress |
| M1 | JSON:API management API | planned |
| M2 | State backend `/tfstate/v1` with locking | planned |
| M3 | TFE v2 remote CLI workflow (`terraform init/plan/apply`) | planned |
| M4 | Executor: clone, run, stream logs, store state | planned |
| M5 | Upstream UI compatibility, GraphQL + subscriptions | planned |
| M6 | RBAC, templates, VCS, policies, registry, parity | planned |

## Compatibility notes

- API paths, payload shapes and table names are taken from the upstream
  implementation. See `docs/compat.md` for the verified inventory.
- Upstream is Apache-2.0. Code here is written from scratch, but the API and
  schema are derived from upstream; keep attribution when distributing.
