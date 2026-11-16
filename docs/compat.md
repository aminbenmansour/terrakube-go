# Compatibility contract

Terrakube Go is a rewrite, but not a redesign. Every surface below was verified
against the upstream Java implementation (`terrakube-io/terrakube`, `main` on
2026-11-14). When in doubt, upstream behaviour wins.

## Discovery

`GET /.well-known/terraform.json` — implemented by
`api/.../plugin/state/WellKnownWebServiceImpl.java`.

```json
{
  "login.v1": { "client": "...", "grant_types": ["authz_code"], "authz": "...", "token": "...", "ports": [10000, 10001] },
  "state.v2": "/remote/state/v2/",
  "tfe.v2": "/remote/tfe/v2/",
  "tfe.v2.1": "/remote/tfe/v2/"
}
```

Notes:

- `login.v1` comes from Dex unless the optional OAuth broker is enabled.
- Upstream advertises `state.v2` at `/remote/state/v2/` but the state handler is
  mounted at `/tfstate/v1` (`TerraformStateController`). The executor generates
  a `backend "http"` override pointing at
  `/tfstate/v1/organization/{orgId}/workspace/{wsId}/backend/terraform.tfstate`;
  the discovery `state.v2` entry is not the path the CLI uses. We mirror the
  discovery document byte-for-byte and serve state under `/tfstate/v1`.

## TFE API — `/remote/tfe/v2/`

Source: `api/.../plugin/state/RemoteTfeController.java`. All responses use
`application/vnd.api+json`.

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/ping` | Health probe used by the CLI |
| GET | `/organizations/{org}` | Organization lookup by name |
| GET | `/organizations/{org}/entitlement-set` | Feature entitlement stub |
| GET | `/organizations/{org}/capacity` | Organization capacity |
| GET | `/organizations/{org}/workspaces` | List workspaces (filters, tags) |
| POST | `/organizations/{org}/workspaces` | Create a workspace (implicit creation) |
| GET | `/organizations/{org}/workspaces/{ws}` | Workspace by name |
| PATCH | `/workspaces/{id}` | Update workspace |
| POST | `/organizations/{org}/projects` | Search projects by name (needed on create) |
| GET | `/workspaces/{id}/tag-bindings` | List tag bindings |
| PATCH | `/workspaces/{id}/tag-bindings` | Set tag bindings |
| POST | `/workspaces/{id}/relationships/tags` | Replace workspace tags |
| POST | `/workspaces/{id}/actions/lock` | Lock workspace |
| POST | `/workspaces/{id}/actions/unlock` | Unlock workspace |
| POST | `/workspaces/{id}/state-versions` | Create state version |
| GET | `/state-versions/{id}` | Read state version |
| GET | `/workspaces/{id}/current-state-version` | Current state version |
| GET | `/workspaces/{id}/current-state-version-outputs` | Outputs |
| POST | `/workspaces/{id}/configuration-versions` | Create configuration upload |
| GET | `/configuration-versions/{id}` | Configuration status |
| PUT | `/configuration-versions/{id}` | Upload configuration tarball |
| GET | `/configuration-versions/{id}/terraformContent.tar.gz` | Download configuration |
| POST | `/runs` | Create run |
| GET | `/workspaces/{id}/runs` | List workspace runs |
| GET | `/organizations/{org}/runs/queue` | Organization run queue |
| GET | `/runs/{id}` | Read run |
| GET | `/runs/{id}/run-events` | Run event stream |
| POST | `/runs/{id}/actions/apply` | Approve/apply |
| POST | `/runs/{id}/actions/discard` | Discard |
| GET | `/runs/{id}/policy-checks` | Policy checks |
| GET | `/policy-checks/{id}` | Policy check |
| GET | `/policy-checks/{id}/output` | Policy output |
| POST | `/policy-checks/{id}/actions/override` | Override policy |
| GET | `/plans/{id}` | Plan status/log reference |
| GET | `/applies/{id}` | Apply status/log reference |
| GET | `/plans/logs/{a}/{b}` | Plan log |
| GET | `/applies/logs/{a}/{b}` | Apply log |

Implemented in the M3 core slice: `ping`, organization lookup, entitlement-set,
capacity, workspace list/get/create/update, workspace lock/unlock, project
list/create, and workspace tag bindings. Configuration versions, runs, state
versions, policy checks and logs are the remaining M3 work.

Tag bindings use `GET`/`PATCH /workspaces/{id}/tag-bindings` for key/value tags
and `POST /workspaces/{id}/relationships/tags` for legacy key-only tags. Keys
are at most 128 characters, values at most 256, and a workspace has at most 10
tags. `PATCH` keeps bindings that are absent from the request, matching upstream.
Workspace creation accepts both `relationships.tags` and
`relationships.tag-bindings`.

State versions are served at `POST /workspaces/{id}/state-versions`,
`GET /state-versions/{id}` and `GET /workspaces/{id}/current-state-version`.
A supplied `state`/`json-state` is base64 decoded and stored under the M2
`state/...` layout, with a dummy completed job and step matching upstream.
Omitting `state` returns a `pending` version with hosted upload URLs backed by
`temp_archive`, as upstream does. `current-state-version` reads the current
backend state and returns its serial, Terraform version and `finalized` status.

Authentication for this surface is a bearer token. During local development the
`TERRAKUBE_AUTH_TOKEN` value is accepted; internal executor JWTs are also
accepted. `GET /ping` is intentionally unauthenticated and returns `404` with
`TFP-API-Version: 2.5` and `TFP-AppName: Terrakube`, matching upstream.

`/.well-known/terraform.json` is served with `login.v1`, `state.v2`,
`tfe.v2` and `tfe.v2.1`. The `login.v1` block points at Dex by default
(`TERRAKUBE_DEX_CLIENT_ID`, `TERRAKUBE_DEX_ISSUER_URL`); setting
`TERRAKUBE_TERRAFORM_LOGIN_ENABLED=true` switches it to the optional OAuth
broker URLs under `TERRAKUBE_PUBLIC_URL`.

## State backend — `/tfstate/v1`

Source: `api/.../plugin/storage/controller/TerraformStateController.java`,
`api/.../plugin/security/state/StateService.java` and
`executor/.../plugin/tfstate/http/HttpTerraformStateImpl.java`.

| Method | Path | Notes |
| --- | --- | --- |
| GET | `/organization/{orgId}/workspace/{wsId}/backend/terraform.tfstate` | HTTP backend read; `204` when absent, otherwise `200` `application/json` |
| POST/PUT | same | Store current state; returns `200`; validates `Content-MD5`; `400` for empty body or bad MD5; `413` above the size limit |
| DELETE | same | Delete current state; returns `200` |
| GET | `/organization/{orgId}/workspace/{wsId}/state/terraform.tfstate` | Current state |
| GET | `/organization/{orgId}/workspace/{wsId}/state/{name}.json` | Archived JSON state |
| GET | `/organization/{orgId}/workspace/{wsId}/jobId/{jobId}/step/{stepId}/terraform.tfstate` | Plan artifact |
| PUT | `/archive/{archiveId}/terraform.tfstate` | Hosted state upload; `201` on success, `403` for an unknown archive |
| PUT | `/archive/{archiveId}/terraform.json.tfstate` | JSON hosted state upload |
| PUT | `/organization/{orgId}/workspace/{wsId}/rollback/{name}.json` | Rollback; `201` on success |

Notes:

- There is **no `LOCK`/`UNLOCK` handler** upstream. The executor writes an
  HTTP backend with only `address`, so Terraform does not perform state
  locking. Workspace locking is the TFE action
  `POST /workspaces/{id}/actions/lock`.
- `POST`/`PUT` on the backend path only update the current state object. The
  `history` rows are created by the TFE state-version and executor apply flows.
- `GET` on the `/state/...` and plan-artifact paths returns an empty `200`
  response when the object is missing, because the controller returns the
  storage byte slice directly.

Auth for executor calls is HTTP Basic with username `terraform` and an
internal JWT as the password (bearer tokens are also accepted). The JWT is
issued with `iss=TerrakubeInternal`, `sub=TerrakubeInternal (EXECUTOR)`,
`aud=TerrakubeInternal`, and `organizationId`, `workspaceId`, `jobId` and
`stepId` claims. The signing key is the base64url-decoded
`io.terrakube.token.internal` secret; HMAC algorithm selection follows key
length (HS256, HS384, HS512). Missing or invalid credentials return `401`;
wrong organization/workspace scope or a non-executor token returns `403`.
Writes additionally require an active step.

Storage keys mirror upstream per backend:

- Local: `backend/{orgId}/{wsId}/terraform.tfstate`,
  `state/{orgId}/{wsId}/state/{historyId}.raw.json`,
  `state/{orgId}/{wsId}/state/{historyId}.json` and
  `state/{orgId}/{wsId}/{jobId}/{stepId}/terraformLibrary.tfPlan`.
- S3: `tfstate/{orgId}/{wsId}/terraform.tfstate`,
  `tfstate/{orgId}/{wsId}/state/{historyId}.raw.json`,
  `tfstate/{orgId}/{wsId}/state/{historyId}.json` and
  `tfstate/{orgId}/{wsId}/{jobId}/{stepId}/terraformLibrary.tfPlan`.

Archive uploads are unauthenticated upstream; the unguessable archive UUID is
the capability. A successful upload deletes the `temp_archive` row. The raw
upload writes the current state plus `{historyId}.raw.json`; the JSON upload
writes `{historyId}.json`. Rollback reads the archived JSON and raw state,
creates a new `history` row with `serial=1`, `md5=0`, `lineage=0`, writes the
new current/raw/JSON objects, and stores the new state URL in `history.output`.

## Management API — `/api/v1`

Elide-generated JSON:API. Reference specification (upstream, v2.27.0) is at
`openapi-spec/v2_27_0.yml`; resources include organization, project, workspace,
variable, tag, workspacetag, job, step, template, history, team, token, vcs,
module, provider, schedule, agent, webhook and policy objects.

Query conventions:

- `filter[type]=...` uses RSQL: `name==prod*`, `status=in=(failed,cancelled)`,
  `a==1;b==2` (AND), comma between expressions (OR).
- `sort=name`, `sort=-name`.
- `page[size]`, `page[number]`, `page[totals]` (adds `meta.totals`).
- `include=step`, `fields[workspace]=name`.
- Content type `application/vnd.api+json`; errors use a JSON:API `errors`
  array. `DELETE` must not send a body or content type.

Our initial implementation targets the resources needed by the CLI workflow and
the UI's core screens. Anything missing is filled in against the upstream spec.

### Implemented in M1

| Resource | Routes |
| --- | --- |
| organization | `GET/POST /organization`, `GET/PATCH/DELETE /organization/{id}` |
| project | `GET/POST /organization/{id}/project`, `GET/PATCH/DELETE /project/{id}`, `GET /project/{id}/workspace` |
| workspace | `GET/POST /organization/{id}/workspace`, `GET/PATCH/DELETE /workspace/{id}` |
| variable | `GET/POST /workspace/{id}/variable`, `GET/PATCH/DELETE /variable/{id}` |

Current M1 limitations (tracked for later milestones):

- Filters, sorting and paging run in memory over the full collection. This is
  fine for smoke tests but must become SQL before production use.
- Authentication is a single static bearer token (`TERRAKUBE_AUTH_TOKEN`).
  PATs, team tokens, executor tokens and Dex JWTs are not implemented yet.
- `include=`, sparse fieldsets, relationship endpoints and RSQL beyond
  `attr==value` / `attr==prefix*` are not implemented.
- Variable values are stored in plaintext. Upstream encrypts them with
  `io.terrakube.encryption.*`; this must land before anything is exposed.
- `tag` and `workspacetag` tables exist but have no handlers yet.

## Executor contract

- API enqueues a run; executor claims it.
- Executor fetches the configuration tarball, injects variables, runs
  `terraform`/`tofu`, and reports step/job status.
- Live output goes to Redis; final output and state go to object storage.
- Storage keys and the state archive layout must remain compatible with
  upstream history rows (`history.md5`, `history.serial`, `history.lineage`).

## Database

Migrations intentionally use upstream table and column names so that existing
data can be imported and the UI's assumptions hold. Core tables implemented in
`migrations/0001_core.sql`:

`organization`, `project`, `workspace`, `variable`, `tag`, `workspacetag`,
`template`, `job`, `step`, `history`.

`migrations/0002_temp_archive.sql` adds upstream's `temp_archive` table
(`id`, `type`, `history_id`) used by the hosted-state and rollback flows.

Differences from upstream at this stage:

- UUID primary keys stay `varchar(36)` as upstream.
- Audit columns use upstream names (`created_date`, `created_by`,
  `updated_date`, `updated_by`).
- Additional upstream tables (VCS, teams, tokens, registry, policies) arrive in
  later migrations rather than all at once.

## UI contract

The React UI talks REST/JSON:API plus:

- **GraphQL over HTTP** at `/graphql/api/v1`. This is not optional: the UI uses
  it for workspace and project lists, organization summaries, tags and more
  (`ui/src/modules/workspaces/workspaceService.ts`,
  `ui/src/modules/projects/projectService.ts`,
  `ui/src/modules/organizations/organizationService.ts`). Only the organization
  list currently uses REST.
- **GraphQL subscriptions** at `/subscriptions` via **graphql-ws** for live job
  updates (`ui/src/modules/api/subscriptionClient.ts`).
- **SSE** event streams for run logs (`eventStream.ts`, `Last-Event-ID`
  support).

Consequence: M5 must implement a GraphQL query/subscription surface, not just a
WebSocket. Budget for it explicitly, or plan a UI rewrite instead.

## Open questions

- GraphQL schema and subscriptions used by the UI (M5) — the biggest hidden
  cost of reusing the upstream UI.
- Exact JSON:API error objects and RSQL edge cases (differential tests).
- Terraform login/OAuth broker behaviour (M3/M5).
- Kubernetes ephemeral executor contract (M4/M6).
- Variable encryption and key management (before variable values are exposed).
