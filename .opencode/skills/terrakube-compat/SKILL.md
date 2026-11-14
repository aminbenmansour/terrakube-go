---
name: Terrakube compatibility
description: Check API, schema and CLI changes against the upstream Terrakube contract
---

# Terrakube compatibility

Use this skill when a change touches an HTTP endpoint, JSON:API payload,
database table/column, auth behavior or Terraform CLI integration.

## Workflow

1. Read `docs/compat.md`; it is the verified contract.
2. If the upstream checkout is available (for example
   `/tmp/opencode/terrakube-upstream`), inspect the Java source named in the
   document.
3. Compare method, path, request/response media type, status codes, error shape
   and auth.
4. Check table/column names and migration ordering.
5. Add or update tests and `docs/compat.md`.
6. Run `make check` and, for API changes, `scripts/smoke.sh`.

## Checklist

- [ ] JSON:API content type and `errors` array
- [ ] `DELETE` has no body or content type
- [ ] Status codes match upstream (200/201/204/404/409)
- [ ] Filtering, sorting, paging and `meta.totals`
- [ ] Relationship links and `include`
- [ ] Auth: bearer token locally; executor basic auth in the state backend
- [ ] Table/column names mirror upstream
- [ ] Real CLI behavior verified when possible

## Known discrepancies

- Discovery says `/remote/state/v2/`; the handler is `/tfstate/v1`.
- The UI needs GraphQL at `/graphql/api/v1`, not only REST.
- Variable values are plaintext pending encryption.
