# UI REST operations

The UI's GraphQL surface is small and is covered by
[`ui-operations.md`](ui-operations.md). Its REST surface is not: writes, run
details, variables, tags, teams, VCS connections, the registry and policy all go
through `/api/v1` (Elide) or `/registry/v1`. The verbatim list of 106 operations
the UI sends is in [`ui-rest-operations.txt`](ui-rest-operations.txt), captured
from `ui/src` and regenerated with `scripts/ui-rest-inventory.sh`.

## Where we stand

`internal/api` currently answers 100 routes. Everything marked *missing* below
returns 404, and the page that calls it shows an error or an empty state.

| Resource | Operations | Status |
| --- | --- | --- |
| `organization` | get | implemented |
| `project` | list, create, get, patch, delete | implemented |
| `workspace` | list, create, get, patch, delete | implemented |
| `workspace` relationships (`project`, `vcs`, `ssh`, `agent`) | patch | implemented |
| `collection`, `item`, `reference` | 12 (collections, items, workspace references, `/reference/:id/collection`) | implemented |
| `job` | get, patch | implemented |
| `tag`, `workspaceTag` | 10 (tag CRUD, workspace bindings list/create/patch/delete) | implemented |
| `template` | 5 (list, get, create, patch, delete) | implemented |
| `ssh`, `agent` | 5 | implemented |
| `vcs` | 5 (list, get, create, patch, delete) | implemented |
| `team` | 5 (list, get, create, patch, delete) | implemented |
| `provider`, `module` | 19 (versions, implementations, module versions) | **missing** |
| `policy_set`, `policy_exemption` | 19 (attachments, parameters, parameters CRUD) | **missing** |
| `notificationConfiguration`, `notification_trigger` | 7 (config CRUD, triggers create/delete, templates relationship) | implemented; the two `/notification/v1/.../test` sends are **missing** (delivery, #52) |
| `federated` (OIDC credentials) | 9 | **missing** |
| registry `/registry/v1/providers` | 4 | separate service (`cmd/registry`, M6) |
| `/remote/tfe/v2/runs/:id/policy-checks` | 1 | implemented |
| `/organization/:id/globalvar` | 4 | implemented |
| `POST /operations` (atomic extension) | 2 (`add`, `update` workspace) | implemented |
| `action` | list, get, create, patch, delete | implemented |
| `/organization/:id/operations`, `/doc` | 2 | **missing** — the vocabulary endpoint the filter builder reads |

Consequence: the UI can load and browse organizations, projects, workspaces and
the workspace page's GraphQL data, but the variable editor, run details, the
registry and policy screens do not work against this API yet. That is the bulk of M5 and it is REST work, not GraphQL work.

## Order that unblocks the most UI

1. `job` get/patch and the workspace relationship patches (lock, project, VCS).
2. `provider`, `module` and the registry service.
3. `policy_set`, `policy_exemption`, `federated`.
4. The legacy `/workspace/:id/variable` editor routes are served; the collection
   screens are done.

## Measured route surface

`internal/api/routes_test.go` pins the routes the management API registers, so
the list above cannot go stale in the optimistic direction: deleting or renaming
a route fails the test until the golden list is updated with it.
