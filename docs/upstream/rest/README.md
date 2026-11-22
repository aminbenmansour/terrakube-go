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
| `include=` | adds an `included` member | same for the workspace's `variable`, `organization`, `project`, `reference`, `ssh` and `agent`; other names are ignored |
| Workspace PATCH | **204 with no body**, attributes or relationships | same (fixed 2026-11-21) |
| `relationships.project` on workspace | stored; unknown id `404 Unknown identifier ...`, other organization `403 Project does not belong ...` | same |
| `relationships.vcs`/`.ssh`/`.agent` on workspace | stored, validated against their tables | same |
| `PATCH /relationships/{relation}` | `204`, project validated; vcs/ssh/agent must exist | same |
| Duplicate `globalvar` key | **500** (a lifecycle hook throws) | **409** — the unique constraint, not upstream's bug |
| `action` ids | editor slugs (`terrakube.open-ai`), not UUIDs | same |
| `action` list filter | `filter[action]=active==true;type=in=('A','B')` | same |
| `action` create / patch / delete | **201** with the resource / **204** / **204**; audit fields set from the PAT's user | same statuses; `createdBy`/`updatedBy` are null (no user in the token) |
| `ssh` | organization-scoped CRUD; `privateKey` is never returned; delete in use **409** | same; no team permission check yet (#44) |
| `agent` | list, create, delete; delete in use **409**; superuser-only writes | same routes and 409; writes open (no RBAC yet) |
| `vcs` | organization-scoped CRUD; `name`/`description`/`vcsType` required; invalid enum **400**; secrets never returned; delete unlinks its workspaces | same routes and statuses; secrets plaintext until #54 |
| `tag` | organization-scoped CRUD; optional name; duplicate name **500** | same routes and optional name; duplicate **409** (unique index) |
| `workspacetag` | list/create/patch/delete under the workspace; missing/unknown `tagId` **409**; duplicate binding **500** | same routes and 409; duplicate binding **409** |
| `collection` / `item` / `reference` | organization-scoped CRUD, `GET /reference/:id/collection?include=item`; writes gated on a team permission; missing required attribute **409**; no `priority` means **0** | same routes and statuses; writes are open (no RBAC yet) |
| `POST /operations` (atomic) | `add`/`update`; `Content-Type` **and** `Accept` must be the atomic media type or **415**; an unsupported op is a no-op 200 | same shape and 415; an unsupported op is **400** and a failure is one error object, not an array of per-op results |
| Resources | `type`, `id`, `attributes`, `relationships` | `type`, `id`, `attributes`, and `relationships` for the links that exist |
| Error objects | `detail` only | `detail`, `status`, `title` |

The differences that matter are deliberate: we are a superset where upstream is
stricter (paging on `organization`, extra error members), and leaving a
relationship out is honest — the tables behind most of them (`team`,
`policySet`, `template`, `provider`, …) do not exist yet.
Workspace `vcs`, `ssh` and `agent` ids are validated against their tables now.
See `docs/upstream/ui-rest-operations.md`.

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
