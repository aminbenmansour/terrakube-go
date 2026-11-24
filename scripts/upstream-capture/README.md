# Capturing the upstream contract

The `/api/v1` contract is whatever `azbuilder/api-server` actually returns, not
what the OpenAPI file says. This directory runs a disposable copy of the
released upstream image, records responses verbatim, and compares them with this
API.

The recordings live in [`docs/upstream/rest/`](../../docs/upstream/rest/) and are
read by `internal/api/upstream_contract_test.go`, so a shape change fails
`make check` instead of being discovered in the UI.

```bash
./start.sh                    # postgres 5433, redis 6380, dex 5556, api 18080
./pat.py > token.txt          # a PAT the upstream API accepts
./capture.sh                  # writes docs/upstream/rest/{*.json,index.tsv}
./compare.sh                  # our API (8080) versus upstream (18080)
```

`start.sh` binds on the host network on non-default ports so it cannot collide
with `make services-up`. It seeds upstream's `demo` profile, so ids and names in
the recordings come from there; they are stable for a given image tag.

## Why a hand-minted PAT

Upstream accepts a Dex token or a PAT. Dex static passwords carry no `groups`
claim and every organization read requires the instance-owner group, so the PAT
path is the only one that works without an LDAP server. `pat.py` inserts a
`pat` row and signs an HS256 token with the `PatSecret` from `start.sh`; the
signing key is the base64url-decoded secret, not the raw string.
`PAT_GROUPS` mints a token for other group names (`PAT_GROUPS=capture-team
./pat.py`), which is how the team-derived permission captures are taken.

`token.txt` and `.secrets.env` are git-ignored. `start.sh` generates the shared
PAT and internal secrets on first run and gives them to the API; `pat.py` reads
them back from `.secrets.env` (or `PAT_SECRET`). Nothing that signs a token is
committed.

## Adding a case

1. Add a `request` line to `capture.sh` (name, path, optional curl arguments).
2. Re-run the capture against the stack.
3. Add the fixture to the table in `internal/api/upstream_contract_test.go`:
   either replay it, or list it under `deferredCaptures` with the reason. The
   test fails on a fixture that is in neither, so new captures cannot be
   silently ignored.

Statuses are compared exactly except where `docs/upstream/rest/README.md` records
a deliberate divergence.
