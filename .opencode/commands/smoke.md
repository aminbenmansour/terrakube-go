---
description: Start dependencies and run the API smoke test
---

Run the end-to-end smoke test without deleting data:

1. `make services-up`
2. `make db-migrate`
3. Start the API on a spare port with `TERRAKUBE_AUTH_TOKEN=dev-token` and
   `TERRAKUBE_HTTP_ADDR=:18080` (use `bin/terrakube-api` after `make build`).
4. Run
   `TERRAKUBE_API_URL=http://localhost:18080 TERRAKUBE_AUTH_TOKEN=dev-token ./scripts/smoke.sh`.
5. Stop the API process.
6. Report the result. Do not run `make services-wipe`.
