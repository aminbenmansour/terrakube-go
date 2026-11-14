---
description: Regenerate sqlc code after SQL or schema changes
---

Regenerate the sqlc store. $ARGUMENTS

1. Edit SQL in `internal/store/queries/*.sql`; add a new migration if the
   schema changes.
2. Run `make generate`.
3. Run `go build ./...`, `go vet ./...` and `go test ./...`.
4. Review the generated diff. Never edit `internal/store/*.sql.go` by hand.
5. Update `docs/compat.md` if the public contract changed.
