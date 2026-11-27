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

Implemented in M3: `ping`, organization lookup, entitlement-set, capacity,
workspace list/get/create/update, workspace lock/unlock, project list/create,
workspace tag bindings, state versions, configuration versions, the run
lifecycle (including step creation, apply/discard and plan/apply log
capabilities) and executor-driven plan/apply. Policy-check details and the
policy override flow remain M6 work; `run-events` stays an empty list as
upstream does, and live logs are served by `/tfoutput/v1`.

Tag bindings use `GET`/`PATCH /workspaces/{id}/tag-bindings` for key/value tags
and `POST /workspaces/{id}/relationships/tags` for legacy key-only tags. Keys
are at most 128 characters, values at most 256, and a workspace has at most 10
tags. `PATCH` keeps bindings that are absent from the request, matching upstream.
Workspace creation accepts both `relationships.tags` and
`relationships.tag-bindings`.

Workspace locking uses `POST /workspaces/{id}/actions/lock` and
`/actions/unlock`. Locking preserves any existing `lock_description` (used by
PR apply locks); unlocking clears it. A locked workspace still accepts run
creation, matching upstream: the lock blocks dispatch, not the API call. There
is no `LOCK`/`UNLOCK` method on `/tfstate/v1`.

State versions are served at `POST /workspaces/{id}/state-versions`,
`GET /state-versions/{id}` and `GET /workspaces/{id}/current-state-version`.
A supplied `state`/`json-state` is base64 decoded and stored under the M2
`state/...` layout, with a dummy completed job and step matching upstream.
Omitting `state` returns a `pending` version with hosted upload URLs backed by
`temp_archive`, as upstream does. Supplying `state` also overwrites the current
backend state, matching upstream `uploadState`. `current-state-version` reads the
current backend state and returns its serial, Terraform version and `finalized`
status. `GET /workspaces/{id}/current-state-version-outputs` exposes the root
module outputs as `state-version-outputs` resources.

Configuration versions are served at
`POST /workspaces/{id}/configuration-versions` and
`GET /configuration-versions/{id}`. Creation returns `pending` with an absolute
`upload-url` carrying an HMAC-signed, one-hour ticket. The tarball is uploaded
with `PUT /configuration-versions/{id}?ticket=...`, stored at
`content/{id}/terraformContent.tar.gz`, and the version moves to `uploaded`.
`GET /configuration-versions/{id}/terraformContent.tar.gz` returns the stored
tarball. The upload and download endpoints do not require a bearer token,
matching the upstream controller.

Runs are served at `POST /runs`, `GET /runs/{id}`,
`GET /workspaces/{id}/runs`, `GET /organizations/{name}/runs/queue`,
`POST /runs/{id}/actions/apply` and `POST /runs/{id}/actions/discard`. Creation
validates the workspace and uploaded configuration version, writes the
configuration tarball URL to `job.override_source` with
`job.override_branch = remote-content`, creates the CLI steps `100` (plan), `150`
(approval) and `200` (apply), publishes the job ID to the Redis Stream
`terrakube:jobs`, and returns the upstream run shape (`run-<job id>`).
`GET /runs/{id}?include=plan,apply` adds the related plan/apply resources under
`included`, and `run-events` is exposed as an empty to-many relationship.

Plan and apply status come from the matching step status, as upstream does:
`pending` → `pending`, `queue`/`running` → `running`, `completed` → `finished`,
`failed` → `errored`. While the plan step is finished but the approval step is
still open and the plan has changes, the plan stays `running` so the CLI keeps
polling until the executor reaches the approval gate. The run's own status maps
`completed` → `finished`, `pending` → `pending`, `running`/`queue`/
`waitingApproval` → `running`, `failed` → `errored`, default → `unknown`.
Workspace and organization run listings only include `pending`, `running`,
`queue` and `waitingApproval`, matching upstream.

Apply (`POST /runs/{id}/actions/apply`) is accepted only for `waitingApproval`
or `completed` jobs. It completes the approval step, records `approved_by` and
`approved_at`, sets the job back to `pending` and re-queues it so the executor
can apply the saved plan. Discard is accepted only while the job still holds the
queue (`pending`, `approved`, `waitingApproval`, `queue`, `running`); it sets
`cancelled`, which the CLI renders as `unknown`.

Plan and apply `log-read-url` values embed an HMAC capability token instead of
an `Authorization` header, mirroring upstream's encrypted plan-id token. The CLI
fetches those URLs directly; a missing or invalid token returns `403`. Logs are
read from the step output, falling back to `job.output` for older rows, and
honor `offset`/`limit` query parameters.

Authentication for the remaining surface is a bearer token. During local development the
`TERRAKUBE_AUTH_TOKEN` value is accepted; internal executor JWTs are also
accepted. A workspace-scoped JWT (one carrying a `workspaceId` claim) is
confined to its `organizationId`/`workspaceId`: a route naming another
organization or workspace, or one resolving a run, state version or
configuration version that belongs to another workspace, answers `403`.
Executor and registry service tokens are unscoped. A workspace-scoped JWT is
only accepted while its step is active, so it expires when the job ends. `GET /ping` is intentionally
unauthenticated and returns `404` with
`TFP-API-Version: 2.5` and `TFP-AppName: Terrakube`, matching upstream.

`/.well-known/terraform.json` is served with `login.v1`, `state.v2`,
`tfe.v2` and `tfe.v2.1`. The `login.v1` block points at Dex by default
(`TERRAKUBE_DEX_CLIENT_ID`, `TERRAKUBE_DEX_ISSUER_URL`); setting
`TERRAKUBE_TERRAFORM_LOGIN_ENABLED=true` switches it to the optional OAuth
broker URLs under `TERRAKUBE_PUBLIC_URL`, with client `terraform-cli` and the
loopback port range `[10000, 10010]`. See “Terraform login — `/oauth`”.

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
Writes additionally require an active step. The configured
`TERRAKUBE_AUTH_TOKEN` (the development user token) is accepted for state reads
because the CLI downloads hosted state with its own credentials; it is never
accepted for backend writes.

Storage keys mirror upstream per backend:

- Local: `backend/{orgId}/{wsId}/terraform.tfstate`,
  `state/{orgId}/{wsId}/state/{historyId}.raw.json`,
  `state/{orgId}/{wsId}/state/{historyId}.json`,
  `state/{orgId}/{wsId}/{jobId}/{stepId}/terraformLibrary.tfPlan` and
  `tfoutput/{orgId}/{jobId}/{stepId}.tfoutput`.
- S3: `tfstate/{orgId}/{wsId}/terraform.tfstate`,
  `tfstate/{orgId}/{wsId}/state/{historyId}.raw.json`,
  `tfstate/{orgId}/{wsId}/state/{historyId}.json`,
  `tfstate/{orgId}/{wsId}/{jobId}/{stepId}/terraformLibrary.tfPlan` and
  `tfoutput/{orgId}/{jobId}/{stepId}.tfoutput`.

Archive uploads are unauthenticated upstream; the unguessable archive UUID is
the capability. A successful upload deletes the `temp_archive` row. The raw
upload writes the current state plus `{historyId}.raw.json`; the JSON upload
writes `{historyId}.json`. Both uploads are bounded by the configured state
size limit (`TERRAKUBE_STATE_MAX_SIZE_BYTES`) and return `413` beyond it.
Rollback reads the archived JSON and raw state, writes the new current/raw/JSON
objects, and only then creates a new `history` row with `serial=1`, `md5=0`,
`lineage=0` and the new state URL in `history.output`.

## Executor logs — `/tfoutput/v1`

Source: `api/.../plugin/storage/controller/TerraformOutputController.java`.

| Method | Path | Notes |
| --- | --- | --- |
| GET | `/organization/{orgId}/job/{jobId}/step/{stepId}` | Step log; `application/octet-stream`; `404` when the job/step pair is unknown or has no output |
| GET | `/organization/{orgId}/job/{jobId}/step/{stepId}/stream` | Server-Sent Events live log |

The plain endpoint returns the archived step output for finished steps and the
live Redis stream for running steps. The stream resumes from `Last-Event-ID`,
emits `id:`/`data:` frames with the Redis stream ID, a `gap` event when the
requested ID was trimmed (`recovery: "step-log"`), `: heartbeat` comments after
idle reads, and an `end` event (`{"reason":"complete"}`) once the job is
terminal. Requests are authorized with the user token or an internal executor
JWT; both endpoints verify that the job belongs to the organization in the path
and that the step belongs to the job. A workspace-scoped JWT must also match the
job's organization and workspace and is refused with `403` otherwise; it is only
accepted while its step is active.

`internal/tfoutput/stream_integration_test.go` covers the live path end to
end against Redis: a record reaches an open SSE client, a record for another
step does not, and finishing the job emits `event: end` with
`{"reason":"complete"}`.

Live records use the Redis stream named after the job ID with the fields
`stepId` and `output` (one log line per record), matching the upstream
executor. Final step output is uploaded to
`tfoutput/{orgId}/{jobId}/{stepId}.tfoutput`, and the plain endpoint serves
that object for finished steps before falling back to `step.output`. This
implementation also keeps the text inline in `step.output`; upstream stores a
`/tfoutput/v1/...` URL there. Imported upstream rows resolve through the same
object key.

## Management API — `/api/v1`

Elide-generated JSON:API. This API implements 100 routes (organization, project,
workspace, variable, globalvar, job, action, collection/item/reference,
ssh, agent, vcs, template, team, notification configuration/trigger,
tag/workspacetag and the atomic extension); the 106
REST operations the UI actually sends are listed in
`docs/upstream/ui-rest-operations.md`, with what is missing. Reference
specification (upstream, v2.27.0) is at
`openapi-spec/v2_27_0.yml`; resources include organization, project, workspace,
variable, tag, workspacetag, job, step, template, history, team, token, vcs,
module, provider, schedule, agent, webhook and policy objects.

Query conventions:

- `filter[type]=...` uses RSQL: `==`, `!=`, `=in=(...)`, `=out=(...)`,
  `=isnull=true|false`, `=ini=`/`=outi=` (case-insensitive), `=lt=`, `=le=`,
  `=gt=`, `=ge=`, `=between=(a,b)`, `=notbetween=(a,b)`, `*` as the wildcard,
  `;` (AND), `,` (OR) and parentheses.
  `==` is case-sensitive; the UI's name search uses `=ini='*term*'`. A filter
  naming an attribute we do not model is ignored whole, so it widens instead of
  hiding rows. A malformed expression is a `400`. Semicolons must be URL-encoded
  (`%3B`). On the workspace list, the related attributes `workspaceTag.tagId`,
  `workspaceTag.value` and `job.status` filter through an `EXISTS` subquery and
  accept the membership operators (`==`, `!=`, `=in=`, `=out=`, `=ini=`,
  `=outi=`), plus `=hasmember=`/`=hasnomember=` and ranges; `!=`/`=out=`/
  `=hasnomember=` become `NOT EXISTS`, so an untagged workspace matches, and
  discarded jobs are ignored. The relationship itself takes `=isempty=`:
  `workspaceTag=isempty=true` lists workspaces with no tag binding.
- `sort=name`, `sort=-name`, or comma-separated keys (`sort=name,-id`); `name`
  compares case-insensitively and `id` is the tiebreaker.
- `page[size]`, `page[number]`, `page[totals]`. Any page parameter makes Elide
  emit `meta.page` (`{number, limit}`, default limit 10000); `page[totals]` adds
  `totalPages` and `totalRecords`. No page parameter means no `meta` at all.
- `include=step`, `fields[workspace]=name`.
- Content type `application/vnd.api+json`; errors use a JSON:API `errors`
  array. `DELETE` must not send a body or content type.
- JSON:API and GraphQL request bodies are limited to 8 MiB and answer `413`
  with a JSON:API error above it. Configuration tarballs and state uploads
  each allow 100 MiB.
- `TERRAKUBE_RATE_LIMIT` caps requests per second per caller and answers `429`
  with `Retry-After`; unset means no limit, as upstream.
- Responses carry `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`,
  `Referrer-Policy: strict-origin-when-cross-origin`, and
  `Strict-Transport-Security` when the request arrives over HTTPS (directly or
  via `X-Forwarded-Proto`). A content security policy is not set because the
  bundled React UI owns its asset policy.

Our initial implementation targets the resources needed by the CLI workflow and
the UI's core screens. Anything missing is filled in against the upstream spec.

### Personal access tokens — `/pat/v1`

Captured from `azbuilder/api-server:2.34.3`. This surface is plain JSON, not
JSON:API, and is served outside `/api/v1`:

- `POST /pat/v1` `{"days":1,"description":"…"}` → `201 {"token":"<jwt>"}`. The
  token is shown once and never stored.
- `GET /pat/v1` → `202` with an array of
  `{id, days, description, deleted, source, lastUsedAt, createdDate,
  updatedDate, createdBy, updatedBy}` for the caller only.
- `DELETE /pat/v1/{id}` → `202` with no body; `404` unknown, `403` another
  user's token. Revocation is a soft delete and takes effect immediately.
- The token is an HMAC JWT with `iss`/`aud` `Terrakube`, `jti` = the `pat` row
  id, `sub`/`name` `"<name> (Token)"`, `email`, `email_verified`, `groups`,
  `iat`, and `exp` only when `days > 0`. Verification checks the signature,
  issuer, expiry and a live non-deleted row; `last_used_at` is stamped at most
  hourly.
- Only an interactive user may mint one. A PAT cannot mint another PAT.
- The signing secret is `TERRAKUBE_PAT_SECRET` (base64, ≥ 32 bytes) and the
  instance owner is `TERRAKUBE_ADMIN_GROUP` (default `TERRAKUBE_ADMIN`), which
  may revoke any token.

### Team permissions — `/access-token/v1`

Captured from `azbuilder/api-server:2.34.3` (`TeamTokenController`). The UI
calls these on every screen, so they are implemented ahead of grant CRUD. Like
`/pat/v1` this is plain JSON served outside `/api/v1`; every answer is **202**,
never JSON:API, except the errors:

- `GET /access-token/v1/teams/current-teams` →
  `202 {"groups":["team-a","TERRAKUBE_ADMIN"]}` — the token's `groups` claim,
  in order, with no membership lookup. (Upstream prefers federated group names
  when that integration is configured; this API has no federated lookup.) A
  token without groups answers an empty array.
- `GET /access-token/v1/teams/permissions/organization/{org}` → `202` with the
  twelve-boolean `PermissionSet` below. `{org}` and `{ws}` must be UUIDs; a
  malformed one is **500** upstream and **400** here.
- `GET /access-token/v1/teams/permissions/organization/{org}/workspace/{ws}`
  → the organization set plus the workspace's `access` rows and its project's
  `project_access` rows. A missing workspace is not an error: it answers the
  organization set with 202.
- No token → **401 with an empty body** upstream; this API sends a JSON:API
  error instead.

The body is `{manageState, manageWorkspace, manageModule, manageProvider,
manageVcs, manageTemplate, manageCollection, manageJob, planJob, approveJob,
managePermission, managePolicies}`. Roles decide; `custom` (and a blank or
unknown role, which is upstream's normalisation) falls back to the row's
booleans:

| check | admin | write | plan | read | custom |
| --- | --- | --- | --- | --- | --- |
| `manageWorkspace`, `manageState` | yes | yes | no | no | flag |
| `manageModule`, `manageProvider`, `manageVcs`, `manageTemplate`, `manageCollection`, `managePolicies` | yes | no | no | no | flag |
| `planJob` | yes | yes | yes | no | flag |
| `approveJob`, `destroyJob` | yes | yes | no | no | flag |
| `manageJob` | yes | yes | yes | no | `planJob` or `approveJob` |

`managePermission` mirrors `manageWorkspace` for teams. A caller whose groups
contain the instance owner group gets `managePermission` and `managePolicies`
set even when no team matches, including for an unknown organization; the other
ten stay as the teams decided. `access` and `project_access` rows have no
module/provider/VCS/template/collection/policy columns, so they only add the
five workspace permissions.

One deliberate divergence: upstream applies **every** `access` row of the
workspace to every caller (`workspace.getAccess().forEach`, no name check),
which would hand out grants to unrelated users. This API only applies the rows
whose `name` is one of the caller's groups; `project_access` rows are filtered
by name upstream too. The captured responses are in
`docs/upstream/rest/access-token-*.json`.

### Workspace and project grants

Source: `api/.../rs/workspace/access/Access.java`,
`api/.../rs/project/access/ProjectAccess.java` and the Elide routes in 2.34.3.

Both are relationship resources upstream (`@Include(rootLevel = false)`), so
their CRUD lives under the parent and is implemented here:

- `GET|POST /api/v1/organization/{org}/workspace/{ws}/access`
- `GET|PATCH|DELETE /api/v1/organization/{org}/workspace/{ws}/access/{id}`
- `GET|POST /api/v1/organization/{org}/project/{project}/projectAccess`
- `GET|PATCH|DELETE /api/v1/organization/{org}/project/{project}/projectAccess/{id}`

Create answers **201** with the resource; PATCH and DELETE answer **204** with
no body, and DELETE carries no content type. Attributes are `name`, `role` and
the five permission booleans plus `destroyJob`; `destroyJob` defaults to
**true**, the rest to **false**, and an omitted `role` stays `null`. `name` is
required for `access` (upstream's data-integrity **409**) but nullable for
`project_access`. `role` is a free string upstream does not validate.

A grant id that does not exist, or that belongs to another workspace/project,
answers **404 `Unknown identifier [id] for access|projectAccess`**. A project
outside the organization answers **404 `Unknown identifier [id] for project`**;
the workspace routes share the generic `The requested resource does not exist.`
detail for a missing or foreign workspace, with the same 404 status. Lists
answer every row for the parent, ordered by name then id; `page`, `sort` and
`filter` parameters are not applied yet, as with the other in-memory lists.
The GraphQL `access` and `projectAccess` relations read the same tables.

Divergences: a duplicate `access.name` within a workspace violates upstream's
`const_workspace_access` unique constraint (added to this schema by `0017`),
which upstream lets surface as **500**; this API answers **409**. An unknown
attribute in a PATCH is ignored
here, while upstream answers **404 `Unknown attribute ...`**, as with templates.
Writes are gated on the `manage workspace access` / `manage project access`
permissions upstream annotates, which resolve to `manageWorkspace` here (see
below). Captures are in `docs/upstream/rest/workspace-access-*.json` and
`project-access-*.json`.

### Permission enforcement

Source: Elide's security annotations on the 2.34.3 resource classes, checked
against a live upstream with owner, admin, write and read tokens.

The caller is the **instance owner** when authentication is disabled, when the
token is the development token (`TERRAKUBE_AUTH_TOKEN`), or when its `groups`
contain `TERRAKUBE_ADMIN_GROUP` (default `TERRAKUBE_ADMIN`). The owner bypasses
every check.

Every other caller must match a `team.name`, an `access.name` or a
`project_access.name` in the organization. A caller that matches none gets
**403** `{"errors":[{"detail":"ReadPermission Denied"}]}` on every
organization-scoped route, reads included; `GET /organization` answers an
empty collection rather than an error. A matching caller may read every
resource in that organization. Writes additionally need the permission
upstream annotates:

| Routes | Write permission |
| --- | --- |
| `vcs`, `ssh` | `manageVcs` |
| `template` | `manageTemplate` |
| `collection`, its items and references | `manageCollection` |
| `workspace`, `project`, `tag`, `workspaceTag`, `variable`, `notificationConfiguration`, `access`, `projectAccess` | `manageWorkspace` |
| `PATCH /organization/{org}/job/{id}` | `manageJob` |
| `organization`, `team`, `agent`, `globalvar`, `action` | owner only |

A workspace `access` row and its project's `project_access` rows add the five
workspace permissions to the team set, so a grant can unlock a write without a
team. Denied writes answer **403** with `CreatePermission Denied`,
`UpdatePermission Denied` or `DeletePermission Denied` by method, matching the
detail the UI displays.

Root-level resources (`/workspace/{id}`, `/variable/{id}`,
`/notification_configuration/{id}`, `/reference/{id}/collection`) resolve their
organization from the resource, so a token cannot reach another organization's
resource by id. `POST /operations` checks the organization or workspace each
operation names. The GraphQL collections are scoped to the caller's
organizations the same way, and the job-status subscriptions refuse a workspace
or organization outside them.

Divergence: upstream applies every `access` row of a workspace to every caller
(see the permission endpoint note above); this API only applies rows naming one
of the caller's groups, and enforces with the same filtered set.

### Terraform login — `/oauth`

Source: `api/.../plugin/token/login/OAuthBrokerController.java`,
`CliLoginService.java`, `DexExchangeClient.java` and `PkceUtil.java` in 2.34.3.

Off unless `TERRAKUBE_TERRAFORM_LOGIN_ENABLED=true`; it then needs
`TERRAKUBE_DEX_ISSUER_URL`, `TERRAKUBE_DEX_CLIENT_ID`, `TERRAKUBE_PAT_SECRET`
and `TERRAKUBE_PUBLIC_URL`. The Dex client must be public and list
`<public URL>/oauth/callback`; the CLI's loopback redirects are intercepted by
the broker, so they need no registration. `terraform login` reads the flow from
`login.v1` in `/.well-known/terraform.json`.

- `GET /oauth/authorize?client_id=terraform-cli&redirect_uri=…&response_type=code&code_challenge=…&code_challenge_method=S256&state=…`
  → `302` to Dex with the broker's own PKCE challenge. `client_id` must be
  `terraform-cli`, the method S256, and `redirect_uri` must be
  `http://{localhost,127.0.0.1,[::1]}:10000-10010/login` with no query,
  fragment or userinfo; anything else is `400`.
- `GET /oauth/callback?code=…&state=…` exchanges the code at Dex and redirects
  the browser to `/oauth/consent`. The `id_token` is verified like a UI bearer
  token with the Dex client as the required audience. `error=access_denied`
  redirects straight to the CLI loopback with `error=access_denied&state=…`.
- `GET /oauth/consent` renders the consent page (signed-in email, token name,
  days). `POST /oauth/consent` is form-urlencoded with `decision`
  (`approve`/`deny`), `days` and `name`; it requires a same-origin
  `Origin`/`Referer` and the signed `tk_cli_login` cookie, then redirects to the
  loopback with `code=…&state=…` or `error=access_denied&state=…`.
- `POST /oauth/token` (form-urlencoded `grant_type=authorization_code`, `code`,
  `code_verifier`, `redirect_uri`, `client_id`) → `200
  {"access_token","token_type":"Bearer","expires_in"}`; bad input is `400
  {"error":"invalid_grant"}`. The one-time code is stored SHA-256 hashed,
  expires after 60 s and can be redeemed once.
- Sessions live in `cli_auth_session` for 10 minutes with the upstream statuses
  (`PENDING_IDP`, `PENDING_CONSENT`, `CODE_ISSUED`, `EXCHANGED`, `DENIED`,
  `FAILED`); expired rows are pruned when the next login starts.
- The issued token is an ordinary personal access token with `source`
  `CLI_LOGIN` and `createdBy` set to the authorizing user. Days default to 30
  and are capped at 365 (`TERRAKUBE_TERRAFORM_LOGIN_DEFAULT_DAYS`,
  `TERRAKUBE_TERRAFORM_LOGIN_MAX_DAYS`).
- `tk_cli_login` is an HMAC-SHA256 cookie over the session id, keyed with
  HKDF-SHA256 of the PAT secret, `Path=/oauth`, `HttpOnly`, `Secure`,
  `SameSite=Lax`.

Verified with `tofu login localhost:8443` against a local Dex v2.44: the CLI
picked port 10004, the broker redeemed the Dex code, the token was stored and
then accepted by `/api/v1`, `/graphql/api/v1` and the PAT list.

### Verified against a live upstream (2026-11-20)

Captured from `azbuilder/api-server:2.34.3`; the responses are in
`docs/upstream/rest/` and the differences are listed there. What matters:

- `GET /organization` matches: id plus `name`, `description`, `disabled`,
  `executionMode`, `icon`, which is all the UI's organization list reads
  (`listOrganizations` in `organizationService.ts`).
- Unknown paths answer `404 {"errors":[{"detail":"Unknown collection <name>"}]}`;
  ours matches. A malformed id is `400 {"errors":[{"detail":"Invalid value: ..."}]}`;
  a well-formed id that does not exist is a 404 whose detail text differs from
  upstream's `Unknown identifier <id> for organization`.
- We are a superset in two places: `organization` accepts page parameters
  (upstream answers 400 `Cannot paginate organization`) and our error objects
  carry `status` and `title` as well as `detail`.
- We are behind in one: most resources still have no `relationships` member
  because the tables behind them do not exist yet. The workspace carries
  `organization`, `project`, `vcs`, `ssh` and `agent`, and its
  `/relationships/{relation}` PATCH route exists for those four.
- `fields[type]` is honoured: only the named attributes are returned, a
  parameter for another type removes attributes entirely, and an empty
  parameter leaves no `attributes` member, as a live 2.34.3 does. A name we do
  not implement is dropped; upstream answers `400 organization does not contain
  the fields: [name]`, but it knows attribute sets we do not model yet and
  refusing a field the UI asked for would break the screen.
- Upstream answers an unauthenticated request with **401 and an empty body**. A
  client must use the status code, not the body.
- Runs: `job` is read and PATCHed with the upstream attribute set (35 keys,
  including the dispatch and PR columns we do not have yet, as null). PATCH
  answers **204 with no body** like upstream, accepts its status enum, rejects
  anything else with `400 Invalid value: ...`, and refuses `approved` unless the
  job is `waitingApproval` (409). It calls the same `UpdateJobStatus` path as the
  TFE discard action, so `jobstatus.Wrap` still publishes to subscriptions.
- Workspace PATCH answers **204 with no body**, for attributes and relationships
  alike, as upstream does. `relationships.project` is stored on create and
  update; an id that does not exist is `404 Unknown identifier <id> for project`
  and a project from another organization is `403 Project does not belong to the
  workspace organization`. `relationships.vcs`, `.ssh` and `.agent` are validated against their
  tables for both create and update: an unknown id is `404 Unknown identifier
  <id> for vcs|ssh|agent`, and one from another organization is `403`. The same four relations have a
  `PATCH /organization/{id}/workspace/{id}/relationships/{relation}` route, the
  one the UI's settings screens call; `moduleSshKey` and `sharedIds` are stored
  as attributes.
- The recordings are replayed by `internal/api/upstream_contract_test.go`, so an
  attribute rename fails `make check` for every replayed capture; a deferred one
  keeps its reason in that test. They are regenerated with
  `scripts/upstream-capture/` against a new upstream tag; a capture that is
  neither replayed nor deferred with a reason fails the build.

### Implemented in M1

| Resource | Routes |
| --- | --- |
| organization | `GET/POST /organization`, `GET/PATCH/DELETE /organization/{id}` |
| project | `GET/POST /organization/{id}/project`, `GET/PATCH/DELETE /project/{id}`, `GET /project/{id}/workspace` |
| workspace | `GET/POST /organization/{id}/workspace`, `GET/PATCH/DELETE /workspace/{id}`, and the same pair under `/organization/{id}/workspace/{id}` |
| variable | `GET/POST /workspace/{id}/variable`, `GET/PATCH/DELETE /variable/{id}`, and the same under `/organization/{id}/workspace/{id}/variable[/{id}]` |
| job | `GET /organization/{id}/job`, `GET/PATCH /organization/{id}/job/{id}` |
| globalvar | `GET/POST /organization/{id}/globalvar`, `GET/PATCH/DELETE /organization/{id}/globalvar/{id}` |
| operations | `POST /operations` — JSON:API atomic `add`/`update` for workspace |

The organization-scoped workspace and variable routes exist because the React UI
uses them; a resource reached through the wrong organization answers 404, and a
variable from a sibling workspace is not reachable through another workspace's
path. `GET /workspace/{id}?include=variable,organization,project` returns those
in `included`, and `relationships.variable` carries the identifiers. Other
`include` names are ignored rather than failing the request, because the UI sends
relationship names whose tables do not exist yet.

Organization `globalvar` rows are inherited by every executor run in the
organization; a workspace variable with the same key in the same category wins,
which is upstream's `putIfAbsent` order. A sensitive global's value is omitted
from every response, including the create response, as upstream does. Duplicate
keys are `409`; upstream's lifecycle hook throws and Spring answers `500`, a
bug recorded in the captures rather than copied.

`action` is the catalog behind the workspace resource drawer and Settings >
Actions. Ids are the editor's slugs (`terrakube.open-ai`), not UUIDs; the list
answers the UI's `filter[action]=active==true;type=in=(...)` expression, reads
are open, and writes exist because the UI's editor sends them (upstream
restricts writes to superusers, as this API does). Creating does not
set `createdBy`/`updatedBy` because a static token carries no user identity, so
they are null; timestamps are set.

`collection`, `item` and `reference` are the Settings > Collections screen and
the read-only "Collection Variables" table on the workspace page. A collection
belongs to one organization and lists its items and references; a reference
binds a collection to a workspace in the same organization, and
`GET /reference/:id/collection?include=item` is the read the workspace page
makes for each reference. Missing required attributes answer upstream's `409`
data-integrity conflict (its Elide/primitive behaviour), an invalid category is
`400`, and a collection created without `priority` gets `0` as upstream does.
Upstream gates collection writes on a team's `manageCollection` permission,
and this API enforces the same one.
`GET /workspace/:id?include=reference` now emits the reference relationship and
`included` rows, and workspace resources carry the audit attributes and
`lastJobDate` upstream serialises.

`ssh` is the Settings > SSH Keys screen and the workspace SSH picker: list, get,
create, patch and delete under `/organization/:id/ssh`. The private key is
write-only — requests accept `privateKey`, but no response echoes it, as
upstream's `read private key` permission hides it. `agent` is the self-hosted
agent pool the Agents screen and the workspace picker use: list, create and
delete. Upstream requires no attribute beyond the JSON:API `type` on either, so
a key without a name is accepted with `name: null`; agent writes are owner-only
and ssh writes need `manageVcs`, as upstream. `GET /workspace/:id?include=ssh,agent` now emits both `included` rows. Deleting
a key or pool a workspace still points at answers upstream's data-integrity
`409` rather than leaving a dangling id. The ssh private key and variable
values are encrypted at rest; see Database below.

`team` is the Settings > Teams screen and the name the access and project-access
entries carry: list, get, create, patch and delete under
`/organization/:id/team`. `name` is the only required attribute (missing or null
answers upstream's `409`); every permission is a boolean, missing ones default to
false except `destroyJob`, which defaults to true as upstream does. `role` is a
free string upstream does not validate, and null is accepted. Upstream gates
team writes on superuser, as this API does.

`notificationConfiguration` is the Settings > Notifications screen. A
configuration is organization-wide or workspace-scoped (created under
`/organization/:id/workspace/:workspaceId/notificationConfiguration`, which
fills both ids); `name`, `channelType` and `destinationUrl` are required and a
missing or null one answers upstream's `409`. `channelType` is
`SLACK|TEAMS|WEBHOOK`; `messageStyle` is `DETAILED|SIMPLE` and defaults to
`DETAILED`; `active` defaults to true. `signingSecret` is write-only: it is
never returned, and PATCHing null clears it. Triggers live at
`/notification_configuration/:id/triggers` (create) and
`/notification_configuration/:id/triggers/:triggerId` (delete); `jobStatus` is
required (`409`) and is upstream's `JobStatus` enum (`400 Invalid value: x`).
`PATCH /notification_configuration/:id/relationships/templates` replaces the
"applies to these templates" set wholesale; an empty list means every template,
and an unknown template is `404`, as upstream. Deleting a configuration cascades
its triggers and template links. Two divergences and one gap: an unknown
attribute in a PATCH is ignored here (upstream `404 Unknown attribute`), the
workspace relationship is emitted as JSON null (upstream emits `data:null` plus
empty `triggers`/`templates`/`policySet` arrays), and the
`/notification/v1/.../test` endpoints and the outbox that drives real delivery
are not implemented yet (#52). The payload builders they will use are
implemented (#154): Slack gets Block Kit cards with a coloured attachment for
`DETAILED` and a one-line top-level `text` for `SIMPLE`; Teams gets adaptive
cards (1.4) with a fact set and an optional View Run action; a generic webhook
gets a flat JSON body (`organization`, `workspace`, `jobId`, `status`, `runUrl`,
`workspaceUrl`, `commitId`, `failureReason`, `configurationName`, absent values
null), signed with
`X-Terrakube-Signature: <lowercase hex HMAC-SHA256 of the raw body>` when
`signingSecret` is set. Writes need `manageWorkspace`, as upstream.


`template` is the Settings > Templates screen and the source of the workspace's
`defaultTemplate`: list, get, create, patch and delete under
`/organization/:id/template`. Upstream requires no attribute beyond the
JSON:API `type`, so an empty `attributes` object answers `201` with nulls; an
unknown id is `404` and a malformed one `400 Invalid value: x`. Deleting a
template leaves `workspace.default_template` alone, as upstream does (no foreign
key). Upstream gates writes on a team's `manageTemplate` permission, and this
API enforces the same one. One divergence: an unknown attribute in a PATCH is ignored
here, while upstream answers `404 Unknown attribute`.

`schedule` is the workspace Schedules tab: list and create under
`/organization/:id/workspace/:id/schedule`, get, patch and delete under
`/organization/:id/workspace/:id/schedule/:scheduleID`. Attributes are `cron`,
`description`, `templateReference`, `enabled` (default `true`) and the audit
fields; the workspace is the only relationship. Writes need `manageWorkspace`,
as upstream's schedule checks delegate to its workspace check. A
`templateReference` must be a template in the workspace's organization: an empty
value is allowed, anything else answers `403` as upstream's security check does.
An unknown schedule is `404` and one reached through another workspace is also
`404`. Cron expressions are stored verbatim — upstream's create hook logs a
Quartz parse failure and still keeps the row — and an expression this API cannot
parse is logged once and skipped by the dispatcher.

The API dispatches schedules itself (upstream uses Quartz): every 20 seconds it
scans the elapsed seconds for matches, so a non-zero seconds field fires
regardless of tick phase, and creates an ordinary plan/approval/apply run with
`via=SCHEDULE`, `created_by=serviceAccount`, the schedule's `tcl` or its
template's body and the template reference. It skips a locked workspace and a
template that has left the workspace's organization, as upstream's trigger does.
Expressions are Quartz's six or seven fields (seconds first; one of
day-of-month/day-of-week must be `?`; weekdays are 1 = Sunday through 7 =
Saturday); the calendar extensions (`L`, `W`, `#`) and month/weekday names are
not supported and are skipped rather than guessed. Only one API replica
dispatches: it holds a Postgres advisory lock and the others retry every 30
seconds.

`runTrigger` is the workspace Run Triggers tab. The UI calls the collection at
the root: `GET/POST /runTrigger` and `GET/PATCH/DELETE /runTrigger/:id`.
Attributes are `enabled` (default `true`) and the audit fields; the
relationships are `sourceWorkspace`, `destinationWorkspace` and an optional
`template`. Both ends must be live workspaces of one organization, must differ,
and a template override must belong to that organization; a malformed shape is
`403`, as upstream's structural check answers. The `organization` is derived
from the destination and is not part of the API surface. Reads are open to every
member of that organization, and a caller outside it sees neither the edge nor
its place in the collection. Creating or changing an edge needs `manageWorkspace`
on the destination and membership on the source; only the instance owner may
repoint an existing edge, so a manager's `PATCH` of `sourceWorkspace` or
`destinationWorkspace` is `403 UpdatePermission Denied` while `enabled` and
`template` are accepted and answer `204`. A duplicate pair is `409`. The
organization's graph must stay acyclic: an edge that would close a loop is
`400`, as is a source workspace past 20 enabled outbound edges. `GET` accepts
`include=sourceWorkspace,destinationWorkspace,template`; the page's
`filter[runTrigger]=sourceWorkspace.id==X,destinationWorkspace.id==X` returns
both directions in one request. `DELETE` with a declared content type is `400`,
matching upstream Elide.

Dispatch is implemented by the API rather than Quartz. The same statement that
marks a job `completed` writes one `run_trigger_event` row when the workspace
has an enabled edge to a live destination, so a crash between the two cannot
lose the dispatch. A poller on the leader replica claims due rows and, for a run
that changed state, creates one ordinary plan/approval/apply job per edge with
`via=RUN_TRIGGER`, `created_by=serviceAccount`, the edge's template or the
destination's `defaultTemplate`, and `triggered_by_job_id` / `cascade_depth`
lineage. A plan-only run never fires, and an apply whose state serial and md5
match the previous history row does not either. A chain stops at cascade depth
10, the runtime net for a cycle that slipped past the write-time check. A
destination without a template, or with a template that has left the
organization, is skipped with a warning. Enqueue failures retry exponentially
(5s to 15min, 10 attempts), and a retry re-queues a still-`pending` dependent
instead of creating a second one. Divergences: the lookback is the DB, not
Quartz's job store, and there is no kill switch or processed-event retention
sweep yet (#200).

`POST /operations` implements the JSON:API atomic operations extension the UI
uses to create and update workspaces. It supports `op: add` (with `data.lid`,
attributes and the same relationships as `POST /workspace`) and `op: update`
(addressed by `href` to the workspace); both require the atomic media type on
`Content-Type` **and** `Accept`, and answer `415` otherwise. The success shape
is upstream's `{"atomic:results":[{"data":…}]}`. Two divergences: an unsupported
`op` is `400` where upstream silently no-ops it, and a failure answers one
JSON:API error object where upstream answers an array of per-operation results.

Current M1 limitations (tracked for later milestones):

- Filters, sorts and paging run in SQL for the organization, project, tag and
  workspace lists. Sortable attributes are `name` and `id` everywhere, plus
  `lastJobDate` on workspaces; modeled filter attributes are `name`,
  `description`, `id`, and for workspaces also `lastJobStatus`,
  `policyComplianceStatus`, `project.id`, `lastJobDate`, `workspaceTag.tagId`,
  `workspaceTag.value` and `job.status`. Other relationship filters are ignored,
  as is the `action` catalog, which keeps its own in-memory two-term parser. The organization list also
  applies the caller's membership restriction in the query, so `page[totals]`
  counts only visible rows.
- Authentication is a static bearer token (`TERRAKUBE_AUTH_TOKEN`), a Dex JWT
  with `TERRAKUBE_DEX_ENABLED=true`, or a personal access token issued at
  `/pat/v1` (see above). Team and executor tokens are not implemented yet.
- `include=` is implemented only for the workspace's `variable`,
  `job`, `schedule`, `organization`, `project`, `reference`, `vcs`, `ssh` and `agent`; relationship
  writes exist only for the workspace's `project`, `vcs`, `ssh` and `agent`.
- Secret columns are encrypted at rest (see Database). `make rotate-keys`
  re-encrypts them under the active key and verifies that nothing remains
  under an old key; it is an operator command, not an API endpoint, and
  changes no schema.
- `tag` and `workspacetag` are served both through the TFE tag-binding
  endpoints and through the Elide-style management routes.

## Compatibility tests

- `internal/tfe/contract_test.go` pins the JSON:API media type, success and
  error envelopes, discovery, ping headers and the main TFE status codes.
- `internal/tfe/*_test.go` cover handler behaviour and edge cases, including
  step-driven plan/apply status and log capability tokens.
- `internal/executor/*_test.go` cover tarball extraction and traversal safety,
  variable/source parsing, log line splitting and state history archiving.
- `internal/tfoutput/handler_test.go` covers archived logs, scoping, auth,
  resume IDs and SSE framing.
- `scripts/tfe-smoke.sh` exercises the real API for discovery, workspaces, tags,
  state, configuration versions, runs, locking and projects.
- `scripts/executor-e2e.sh` (also `make executor-e2e`) drives a real OpenTofu
  plan, apply and destroy through the executor and verifies logs, state outputs
  and current state.
- `scripts/smoke.sh` and `scripts/state-smoke.sh` cover the management API and
  state backend.
- New TFE endpoints must add contract coverage in the same change.

## Executor contract

The M3 executor (`cmd/executor`, `internal/executor`) consumes the Redis Stream
`terrakube:jobs` with the consumer group `terrakube-executor`. It takes a
Postgres advisory lock on the job for the length of the run, so a duplicate
message is acknowledged and skipped instead of running the job twice; the lock
is released with the connection when an executor dies. Messages idle for longer
than the job timeout plus 15 minutes are reclaimed, so a surviving executor
picks up work left by a dead consumer without stealing a job that is still
running.

- **Configuration**: a job with a VCS source clones that repository, otherwise
  `job.override_source` holds either an object-storage key
  (`content/{configurationId}/terraformContent.tar.gz`) or the TFE download URL
  written by run creation. The executor prefers object storage, falls back to
  HTTP, and extracts the gzip tarball with path traversal checks. When the
  workspace has a `folder`, the first entry is used as the working directory.
- **Git sources**: `job.override_source`/`workspace.source` supply the
  repository and `job.commit_id`, `job.override_branch` or `workspace.branch`
  the revision. The synthetic `remote-content` branch is ignored. Clones are
  shallow (`--depth 1`), use `--no-checkout` plus a shallow fetch of the
  revision, and disable interactive credential prompts. Embedded credentials
  are stripped from error messages.
- **Archive safety**: symlinks and hardlinks in an uploaded tarball are dropped.
  Extraction is capped at 128 MiB per file, 1 GiB per archive and 100 000
  entries so a malformed upload cannot exhaust the executor host.
- **Backend**: the executor writes `terrakube_override.tf` containing only the
  HTTP backend `address` at
  `{TERRAKUBE_PUBLIC_URL}/tfstate/v1/organization/{orgId}/workspace/{wsId}/backend/terraform.tfstate`.
  It does not configure locking, matching upstream.
- **Variables**: `ENV` variables become plain environment entries; `TERRAFORM`
  variables become `TF_VAR_*`. Incomplete variables are skipped. Values are
  encrypted at rest (see Database).
- **Secrets**: values of variables marked `sensitive` and the minted internal
  JWT are replaced with `[REDACTED]` in live logs and in the stored job/step
  output. Terraform's own `(sensitive value)` masking is preserved.
- **Isolation**: each job runs in a private `0700` directory under
  `TERRAKUBE_EXECUTOR_WORK_DIR`, removed when the step finishes. Startup sweeps
  `terrakube-job-*` directories older than the job timeout plus one hour, which
  covers an executor that was killed. Terraform commands inherit only an
  allowlist of host variables (`PATH`, `HOME`, `TMPDIR`, locale, proxy and TLS
  certificate settings) plus the workspace environment; the executor's own
  database, Redis, storage and signing secrets are not passed through. The
  executor warns when it runs as root; production deployments should run it as
  an unprivileged user. Host network access and credential scope are the
  deployment's responsibility (container network policy, cloud identity), not
  something the executor sandboxes.
- **Timeouts**: `TERRAKUBE_EXECUTOR_JOB_TIMEOUT` (default `2h`, `0` disables)
  bounds configuration acquisition and every Terraform command. Commands run in
  their own process group, so a timeout or a cancelled job kills the provider
  plugins too and marks the step `failed` or `cancelled` respectively.
- **Auth**: when `TERRAKUBE_INTERNAL_TOKEN_SECRET` is set, the executor mints a
  workspace-scoped internal JWT (two-hour expiry) and passes it as
  `TF_HTTP_PASSWORD` with `TF_HTTP_USERNAME=terraform`. The API confines that
  token to the job's organization and workspace and stops accepting it when the
  step finishes.
- **Steps**: plan is step `100`, approval `150`, apply `200`. Plan runs
  `init` and `plan -detailed-exitcode -out=tfplan`; exit code `2` records
  `plan_changes=true`. The saved plan is stored under the upstream plan key and
  applied by step `200`. `is-destroy` runs `plan -destroy`. `auto-apply` jobs
  approve themselves and continue; otherwise the job moves to
  `waitingApproval`.
- **Logs**: each command line is published to the Redis stream named after the
  job with `stepId`/`output` fields. The stream is capped at 200,000 records and
  deleted when the job finishes, so Redis does not accumulate one stream per
  run; finished steps are served from their archived output. Ordering and step
  boundaries come from the stream order and `stepId`, and the
  stream ID carries the timestamp. When a step finishes, the final output is
  uploaded to `tfoutput/{orgId}/{jobId}/{stepId}.tfoutput` (the upstream S3
  key). `/tfoutput/v1` serves that object for finished steps and falls back to
  `step.output`, which still holds the same text inline (truncated at 10 MiB)
  for the TFE log endpoints; `job.output` keeps the latest step's text for the
  TFE run-log fallback. A log upload failure marks the step `failed` and
  fails the run, unless the step was already cancelled or failed.
- **History**: after a successful apply the executor reads the current state,
  uploads the raw and JSON state objects under the upstream `state/...`
  (local) or `tfstate/...` (S3) keys,
  then creates a `history` row with the state `serial`, `lineage` and MD5 and
  sets `history.output` to
  `{publicUrl}/tfstate/v1/organization/{orgId}/workspace/{wsId}/state/{historyId}.json`.
  The row is created last, so a failed upload leaves no state version pointing
  at a missing artifact. Archiving is idempotent per job reference, so a
  retried job cannot create a second row. A history upload failure marks the
  step `failed` and fails the run.
- **Reliability**: transient failures (database, storage, network) are retried
  up to three attempts by republishing the job with an `attempt` field. A
  message is acknowledged only after its work is finished or safely requeued,
  and acknowledged entries are deleted from `terrakube:jobs` so the stream does
  not grow by one entry per run;
  during shutdown it is left pending so another executor reclaims it after the
  idle window. Permanent Terraform failures are not retried. Log and
  history upload failures are treated as permanent: when a step was completing,
  they mark it `failed` and fail the run, rather than republishing a job that
  would re-run the step.
- **Cancellation**: while a command runs the executor polls the job status. A
  discard or cancel kills the command's process group, marks the current step
  `cancelled` and leaves the API's `cancelled` job status intact.
- **Resume**: a reclaimed job treats a `running` step as the next step, so a
  crashed executor resumes the interrupted step instead of finishing early.
- **Binary**: `TERRAKUBE_EXECUTOR_BINARY` overrides the CLI. Otherwise
  `iac_type` containing `tofu` selects `tofu`; other workspaces use
  `terraform` when installed and fall back to `tofu`.
- **Work dir**: `TERRAKUBE_EXECUTOR_WORK_DIR` sets the temporary root; each job
  gets a unique directory that is removed afterwards.

Divergences to keep in mind:

- `step.output` stores log text inline; upstream stores a `/tfoutput/v1/...`
  URL there. The executor also archives the text under the upstream
  `tfoutput/...` object key, which `/tfoutput/v1` prefers; the column remains
  for the TFE log endpoints and rows written before archiving.
- The executor has no policy check, cost estimation or remote (agent)
  execution yet; those are M6.

## Database

Migrations intentionally use upstream table and column names so that existing
data can be imported and the UI's assumptions hold. Core tables implemented in
`migrations/0001_core.sql`:

`organization`, `project`, `workspace`, `variable`, `tag`, `workspacetag`,
`template`, `job`, `step`, `history`.

`migrations/0002_temp_archive.sql` adds upstream's `temp_archive` table
(`id`, `type`, `history_id`) used by the hosted-state and rollback flows.

`migrations/0004_remote_content.sql` adds `job.override_source` and
`job.override_branch` (upstream stores the CLI configuration tarball URL there
with `override_branch = remote-content`) plus an index on
`step(job_id, step_number)` for step lookups.

Differences from upstream at this stage:

- UUID primary keys stay `varchar(36)` as upstream.
- Audit columns use upstream names (`created_date`, `created_by`,
  `updated_date`, `updated_by`).
- Additional upstream tables (tokens, registry, policies, the notification
  outbox) arrive in later migrations rather than all at once.
- `vcs.client_secret`, `vcs.access_token`, `vcs.refresh_token`,
  `vcs.private_key` and `notification_configuration.signing_secret` are `text`
  rather than upstream's `varchar`, because the encryption envelope below can
  exceed the original widths.

Secret columns are encrypted at rest with a versioned AES-256-GCM envelope
(`enc:v1:<key>:<nonce>:<ciphertext>`) and a per-column AAD. The JSON:API and
GraphQL responses still return plaintext, so the upstream UI and clients see
no difference. The keyring is configured with `TERRAKUBE_ENCRYPTION_KEYS`;
`make db-migrate` backfills legacy rows. `notification_configuration.destination_url`
remains plaintext because upstream treats it as configuration.

## UI contract

The React UI talks REST/JSON:API plus:

- **GraphQL over HTTP** at `/graphql/api/v1`, for ten operations: workspace and
  project lists, workspace/project access lists, notification configuration
  lists and details, and run resolution by name. Everything else, including the
  organization list, uses JSON:API REST. The operations are captured verbatim in
  `docs/upstream/ui-operations.md` and `docs/upstream/ui-operations/`.
- **GraphQL subscriptions** at `/subscriptions` via **graphql-ws** for live job
  updates (`ui/src/modules/api/subscriptionClient.ts`).
- **SSE** event streams for run logs (`eventStream.ts`, `Last-Event-ID`
  support).

Consequence: M5 needs a GraphQL query surface for ten operations over ten entity
types, not the whole Elide schema. The schema only has to cover the entities and
relationship paths in `docs/upstream/ui-operations.md`.

Implemented so far (`internal/graphql`, backed by gqlgen):

- Queries at `/graphql/api/v1` for organization, workspace, project and
  notification_configuration at the root, with every relationship connection the
  UI traverses. Relay connection shape, `pageInfo.totalRecords`, offset cursors,
  `first`/`after` as `StringOrInt`, RSQL `filter` and Elide `sort`.
- Queries are scoped to the caller's organizations once authentication is
  enabled: the root collections add an organization predicate, nested
  connections inherit the parent's scope, and a token without a matching team
  or grant sees nothing. The job-status subscriptions refuse a workspace or
  organization outside that set. The owner rule is the management API's
  (development token, admin group, or authentication disabled).
- Same bearer token as `/api/v1` (`TERRAKUBE_AUTH_TOKEN`). Resolver errors are
  logged and replaced with a generic message.
- **Subscriptions** at `/subscriptions` over graphql-ws: `jobStatusChanged`
  and `organizationJobStatusChanged`, matching upstream's
  `subscriptions.graphqls`. Status updates are published to Redis from a thin
  store wrapper, so every writer (management API and executor) notifies
  subscribers without call-site changes. The UI sends its token as a
  connection parameter, so it is checked after the socket upgrade.
  `TestSubscriptionSeesAStatusWrittenThroughTheStore` covers the whole loop
  against Redis: a status written through the wrapped store reaches an open
  socket.
- Missing tables the UI reads were added in migration `0005_graphql_entities.sql`
  (`access`, `project_access`, `notification_configuration`,
  `notification_trigger`, `notification_configuration_template`).
- **Static hosting.** Set `TERRAKUBE_UI_DIR` to an upstream UI build and the API
  serves it at `/`: `index.html` for `/` and every extensionless path (so deep
  links work), `/assets/*` immutable, a missing asset a real 404, and
  `/env-config.js` generated from `TERRAKUBE_PUBLIC_URL`,
  `TERRAKUBE_DEX_CLIENT_ID`, `TERRAKUBE_DEX_ISSUER_URL` and the build version
  with `Cache-Control: no-store`. UI and API share an origin, so no CORS
  exception is required. `/api/v1`, `/graphql/api/v1` and `/subscriptions` send
  `Access-Control-Allow-Origin: *` by default, as they always have; set
  `TERRAKUBE_CORS_ORIGINS` to a comma-separated allowlist to restrict browser
  access for a split-origin deployment. Verified against a real 2.34.3 build:
  single-origin serving passes `scripts/ui-smoke.sh`, and the split-origin
  container answers GraphQL queries with the UI's origin echoed.
- **Dex authentication.** With `TERRAKUBE_DEX_ENABLED=true` the API accepts the
  Dex JWT the UI stores and sends as `Authorization: Bearer …`, mirroring
  upstream's `DexAuthenticationManagerResolver`. The token is verified against
  the issuer's OIDC discovery document and JWKS: RS256 signature, `iss` and
  `exp` (60 s leeway). `TERRAKUBE_DEX_AUDIENCE_VALIDATION=true` additionally
  requires `aud` to contain `TERRAKUBE_DEX_CLIENT_ID` or one of
  `TERRAKUBE_DEX_EXTRA_AUDIENCES`; it is off by default, as upstream's
  `DexAudienceValidation`. Discovery is cached, signing keys are cached, and an
  unknown `kid` refetches the JWKS (at most once a minute) so Dex key rotation
  works. The same credential is accepted on `/api/v1`, `/graphql/api/v1` and
  the `/subscriptions` connection parameter. Tokens are never logged. This is
  authentication only: group/role authorization is still pending.

- **Contract tests.** `internal/graphql/contract_test.go` validates every
  captured UI operation against the SDL and then executes it against Postgres,
  so a field the UI reads and the SQL behind it are both pinned. CI runs it in
  the `integration` job, since it skips without a database.

Not implemented: GraphQL mutations and the rest of the 53 upstream entities.

## Open questions

- GraphQL mutations for the operations in `docs/upstream/ui-operations.md`.
- Exact JSON:API error objects and RSQL edge cases (differential tests).
- Terraform login/OAuth broker behaviour (M3/M5).
- Kubernetes ephemeral executor contract (M6).
- Variable encryption and key management (before variable values are exposed).
