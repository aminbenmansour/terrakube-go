---
name: local dev stack
description: Start, inspect and troubleshoot Postgres, Redis and MinIO with podman
---

# Local dev stack

Use plain podman through `scripts/dev-services.sh`:

```bash
make services-up      # start containers
make services-status  # show status and ports
make services-down    # stop, keep volumes
make services-wipe    # stop and delete volumes (destructive)
```

There is no `podman compose` provider on the development machine. Use the
script for dependency containers, or Docker Compose for the full stack.

## Ports and images

- Postgres: `localhost:5432`, user/password/database `terrakube`
- Redis: `localhost:6379`
- MinIO: `localhost:9000` (API), `localhost:9001` (console)
- MinIO image: `docker.io/pgsty/minio` (the official image is unavailable)

## Notes

- Local object storage lives in `.data/`; it is git-ignored.
- Migrations run with `make db-migrate` or automatically when
  `TERRAKUBE_AUTO_MIGRATE=true`.
- Do not run `make services-wipe` unless the user explicitly asks; it deletes
  the local database, Redis data and MinIO data.
