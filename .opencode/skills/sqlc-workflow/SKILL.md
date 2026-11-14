---
name: sqlc workflow
description: Change the schema, queries and generated store safely
---

# sqlc workflow

## Rules

- Migrations are append-only. Add `migrations/000N_*.sql`; never edit applied
  files.
- Hand-written SQL lives in `internal/store/queries/*.sql`.
- Generated code in `internal/store/*.sql.go` must never be edited.
- Mirror upstream table/column names; see `docs/compat.md` and
  `migrations/0001_core.sql`.

## Steps

1. Update the schema with a new migration.
2. Add or update queries.
3. Run `make generate` (install with `make tools` if needed).
4. Run `go build ./...`, `go vet ./...` and `go test ./...`.
5. Review the generated diff for unexpected renames or dropped columns.
6. Test migrations against a clean database before claiming success.

## Verification

- `make db-migrate` against a fresh database.
- `make check`.
- API changes: `scripts/smoke.sh`.
