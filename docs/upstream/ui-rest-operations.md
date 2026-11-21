# UI REST operations

The UI's GraphQL surface is small and is covered by
[`ui-operations.md`](ui-operations.md). Its REST surface is not: writes, run
details, variables, tags, teams, VCS connections, the registry and policy all go
through `/api/v1` (Elide) or `/registry/v1`. The verbatim list of 106 operations
the UI sends is in [`ui-rest-operations.txt`](ui-rest-operations.txt), captured
from `ui/src` and regenerated with `scripts/ui-rest-inventory.sh`.

## Where we stand

`internal/api` currently answers 22 routes. Everything marked *missing* below
returns 404, and the page that calls it shows an error or an empty state.

| Resource | Operations | Status |
| --- | --- | --- |
| `organization` | get | implemented |
| `project` | list, create, get, patch, delete | implemented |
| `workspace` | list, create, get, patch, delete | implemented |
| `workspace` relationships (`project`, `vcs`, `agent`) | patch | **missing** |
| variable / `collection` | 15 (collections, items, references) | **missing** — only the legacy `/workspace/:id/variable` routes exist |
| `job` | get, patch | **missing** |
| `tag`, `workspaceTag` | list, create, list workspace tags | **missing** |
| `template` | list | **missing** |
| `team`, `ssh`, `agent`, `vcs` | 16 | **missing** |
| `provider`, `module` | 19 (versions, implementations, module versions) | **missing** |
| `policy_set`, `policy_exemption` | 19 (attachments, parameters, parameters CRUD) | **missing** |
| `notificationConfiguration` | 5 (triggers, templates relationship) | **missing** — the GraphQL reads exist, the writes do not |
| `federated` (OIDC credentials) | 9 | **missing** |
| registry `/registry/v1/providers` | 4 | separate service (`cmd/registry`, M6) |
| `/remote/tfe/v2/runs/:id/policy-checks` | 1 | implemented |
| `/organization/:id/globalvar` | 4 | implemented |
| `/organization/:id/operations`, `/doc`, `/action` | 3 | **missing** — helpers around the vocabulary endpoint |

Consequence: the UI can load and browse organizations, projects, workspaces and
the workspace page's GraphQL data, but the variable editor, run details, VCS
setup, tags, teams, the registry and policy screens do not work against this
API yet. That is the bulk of M5 and it is REST work, not GraphQL work.

## Order that unblocks the most UI

1. `collection` + items (the variable editor is on every workspace page).
2. `job` get/patch and the workspace relationship patches (lock, project, VCS).
3. `vcs`, `ssh`, `agent` (creating a workspace from a repository).
4. `tag`, `workspaceTag`, `template` (the workspace list filters).
5. `provider`, `module` and the registry service.
6. `policy_set`, `policy_exemption`, `federated`, `notificationConfiguration`.

## Measured route surface

`internal/api/routes_test.go` pins the routes the management API registers, so
the list above cannot go stale in the optimistic direction: deleting or renaming
a route fails the test until the golden list is updated with it.
