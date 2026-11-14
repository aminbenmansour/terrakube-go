---
description: Run formatting, vet, tests and build
---

Run the project verification suite in this order:

1. `make check`
2. `make build`
3. If the change touches concurrency or shared state, `make test-race`

Report the exact failing command and the relevant error. Do not run
`make services-wipe`. Do not fix unrelated failures unless the user asks.
