# REST contract, captured from upstream

These files are verbatim responses from **upstream Terrakube 2.34.3** (the Java
image `azbuilder/api-server:2.34.3`), captured on 2026-11-20 with its own demo
seed data. They are the reference for `/api/v1`, not a reading of the OpenAPI
spec: where the two disagree, these win.

`index.tsv` is `name<TAB>status<TAB>path`. An empty file (`organization-no-token.json`)
is an empty body, not a missing capture.

They are not just reference material: `internal/api/upstream_contract_test.go`
replays every capture against the handlers and compares the shape (status,
top-level members, resource type, attribute key set). A capture that is
neither replayed nor explicitly deferred with a reason fails the build. The
tools that produced them are in
[`scripts/upstream-capture/`](../../../scripts/upstream-capture/README.md).

## What a live upstream actually does

| Behaviour | Upstream | This API |
| --- | --- | --- |
| List response | `{"data":[...]}`, no `links`, no `meta` | same |
| Page parameters sent | `meta.page` = `{number, limit}`, default limit **10000** | same (fixed 2026-11-20) |
| `page[totals]` | adds `totalPages`, `totalRecords` to `meta.page` | same |
| `page[size]`/`page[totals]` on **organization** | **400** `Cannot paginate organization` | 200, paginated |
| Unknown path | 404 `{"errors":[{"detail":"Unknown collection nope"}]}` | same |
| Invalid id | 400 `{"errors":[{"detail":"Invalid value: not-a-uuid"}]}` | same |
| Missing id | 404 `Unknown identifier <id> for organization` | 404, different detail text |
| No token | **401 with an empty body** | 401 with a JSON:API error |
| `fields[type]` | only the requested attributes; a parameter for another type removes attributes; an empty list removes the member; an unknown name is **400** | same, except an unknown name is dropped (we do not model every upstream attribute) |
| `include=` | adds an `included` member | ignored — the UI sends none on `/api/v1` |
| Workspace PATCH | **204 with no body**, attributes or relationships | same (fixed 2026-11-21) |
| `relationships.project` on workspace | stored; unknown id `404 Unknown identifier ...`, other organization `403 Project does not belong ...` | same |
| `relationships.vcs`/`.ssh`/`.agent` on workspace | stored, validated against their tables | ids stored; the tables do not exist yet (#48, #120), so an id is not validated |
| `PATCH /relationships/{relation}` | `204`, project validated; vcs/ssh/agent must exist | workspace project/vcs/ssh/agent only; project validated, the others not |
| Duplicate `globalvar` key | **500** (a lifecycle hook throws) | **409** — the unique constraint, not upstream's bug |
| `POST /operations` (atomic) | `add`/`update`; `Content-Type` **and** `Accept` must be the atomic media type or **415**; an unsupported op is a no-op 200 | same shape and 415; an unsupported op is **400** and a failure is one error object, not an array of per-op results |
| Resources | `type`, `id`, `attributes`, `relationships` | `type`, `id`, `attributes`, and `relationships` for the links that exist |
| Error objects | `detail` only | `detail`, `status`, `title` |

The differences that matter are deliberate: we are a superset where upstream is
stricter (paging on `organization`, extra error members), and leaving a
relationship out is honest — the tables behind most of them (`team`,
`policySet`, `template`, `tag`, `collection`, `provider`, …) do not exist yet.
Workspace `vcs`, `ssh` and `agent` ids are stored without their tables, so they
round-trip but an unknown id is not refused. See
`docs/upstream/ui-rest-operations.md`.

`fields` follows the same rule: JSON:API says a server "MUST NOT include
additional fields" when `fields` is sent, and the UI uses it on every
breadcrumb. The one deliberate gap is an unknown field name: upstream knows
attributes we have not modelled and answers 400; we drop the name instead of
breaking a screen over an attribute that does not exist here.

Note the empty 401 body: a client that tries to parse it as JSON will fail. Use
the status code, as the React UI does.

## Shapes worth copying when implementing an M6 resource

`organization-get.json` and `workspace-list-totals.json` are the two useful
references:

- Attributes are serialised **alphabetically**, camelCase (`allowRemoteApply`,
  `globalRemoteState`, `lastJobStatus`), and are always present even when null.
- A workspace carries 25 attributes and 16 relationships; the GraphQL entity
  list in `docs/upstream/graphql-entities.md` describes the same fields under
  GraphQL names.
- Relationship values are `{"data":[...], "links":{...}}` or `{"data":null}` for
  to-one, never omitted.

`workspace-list-paged.json` and `workspace-list-totals.json` are the reference
for `meta.page` with and without `page[totals]`. They document a global
`/workspace` collection this API does not serve; the UI does not use it, so the
answers live in `deferredCaptures` rather than behind a route.

## How these were captured

A disposable upstream stack (Postgres, Redis, Dex and the released API image on
the host network, ports 5433/6380/5556/18080) with its `demo` profile, which
seeds six organizations. Requests were authenticated with a hand-minted PAT:
upstream accepts `iss=Terrakube`, HS256 over the base64url-decoded `PatSecret`,
with `jti` present in the `pat` table and the owner group in `groups`. It is
reproducible with `scripts/upstream-capture/start.sh`, `pat.py` and
`capture.sh`; the stack itself is not vendored.
