# REST contract, captured from upstream

These files are verbatim responses from **upstream Terrakube 2.34.3** (the Java
image `azbuilder/api-server:2.34.3`), captured on 2026-11-20 with its own demo
seed data. They are the reference for `/api/v1`, not a reading of the OpenAPI
spec: where the two disagree, these win.

`index.tsv` is `name<TAB>status<TAB>path`. An empty file (`organization-no-token.json`)
is an empty body, not a missing capture.

## What a live upstream actually does

| Behaviour | Upstream | This API |
| --- | --- | --- |
| List response | `{"data":[...]}`, no `links`, no `meta` | same |
| Page parameters sent | `meta.page` = `{number, limit}`, default limit **10000** | same (fixed 2026-11-20) |
| `page[totals]` | adds `totalPages`, `totalRecords` to `meta.page` | same |
| `page[size]`/`page[totals]` on **organization** | **400** `Cannot paginate organization` | 200, paginated |
| Unknown path | 404 `{"errors":[{"detail":"Unknown collection nope"}]}` | same |
| Invalid id | 400 `{"errors":[{"detail":"Invalid value: not-a-uuid"}]}` | same shape |
| No token | **401 with an empty body** | 401 with a JSON:API error |
| Resources | `type`, `id`, `attributes`, `relationships` | `type`, `id`, `attributes` only |
| Error objects | `detail` only | `detail`, `status`, `title` |

The three remaining differences are additive except the last one in each row,
which is deliberate: we are a superset where upstream is stricter, and omitting
`relationships` is honest — the tables behind them (`vcs`, `team`, `policySet`,
`ssh`, `agent`, `template`, `tag`, `collection`, `globalvar`, `provider`, …) do
not exist yet. See `docs/upstream/ui-rest-operations.md`.

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

## How these were captured

A disposable upstream stack (Postgres, Redis, Dex and the released API image on
the host network, ports 5433/6380/5556/18080) with its `demo` profile, which
seeds six organizations. Requests were authenticated with a hand-minted PAT:
upstream accepts `iss=Terrakube`, HS256 over the base64url-decoded `PatSecret`,
with `jti` present in the `pat` table and the owner group in `groups`. The
capture is reproducible from that stack; the stack itself is not vendored.
