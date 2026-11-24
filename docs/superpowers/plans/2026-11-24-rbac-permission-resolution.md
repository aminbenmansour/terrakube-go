# RBAC: permission resolution and `/access-token/v1` Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Give the stock UI the permissions it asks for, computed from the caller's groups with upstream's role matrix, without changing any existing endpoint.

**Architecture:** A pure `internal/rbac` package holds upstream's `RbacV2Service` matrix over a `Flags` value built from a `team`, `access` or `project_access` row. A new `/access-token/v1` router (mounted beside `/pat/v1`, outside `/api/v1`) answers `current-teams` and the two permission endpoints from the bearer identity. No schema change: `team`, `access` and `project_access` already exist.

**Tech Stack:** Go 1.27, chi, pgx/sqlc, pgx mock-free fakes in handler tests.

**Spec:** GitHub #140 and #141; upstream `TeamTokenController`, `RbacV2Service`, `DexGroupServiceImpl` in `terrakube/api-server:2.34.3`.

## Global Constraints

- Upstream is the contract. Permission endpoints answer **202** with a plain JSON `PermissionSet` (12 booleans, camelCase), never JSON:API.
- `current-teams` answers **202** with `{"groups": [...]}` — the JWT groups claim verbatim, in order.
- Groups come from the token only. No membership table, no caching in this slice.
- The instance owner group (`TERRAKUBE_ADMIN_GROUP`, default `TERRAKUBE_ADMIN`) sets `managePermission` and `managePolicies` regardless of teams.
- Role matrix (`admin`, `write`, `plan`, `read`, `custom`; blank = `custom`; case-insensitive):
  | check | admin | write | plan | read | custom |
  | --- | --- | --- | --- | --- | --- |
  | manageWorkspace | y | y | n | n | flag |
  | manageState | y | y | n | n | flag |
  | manageModule/Provider/Vcs/Template/Collection/Policies | y | n | n | n | flag |
  | planJob | y | y | y | n | flag |
  | approveJob | y | y | n | n | flag |
  | manageJob | planJob or approveJob | | | | |
  | destroyJob | y | y | n | n | flag |
- A divergence to record: upstream's workspace permission endpoint grants the workspace `access` rows to *every* caller (it never compares names). We compare names with the caller's groups.
- Authentication is the existing bearer chain (`TERRAKUBE_AUTH_TOKEN`, Dex, PATs). The development token is treated as instance owner.

## Review Focus

- A token with no `groups` claim must get all-false permissions, not an error and not owner rights.
- A group named like a team in *another* organization must not grant anything: queries filter by organization id.
- `current-teams` must not leak the whole team table; it echoes the token only.
- Blank or mixed-case roles must behave like `custom`/lower-case, never like a missing permission.
- Owner rights must survive an unknown organization id (upstream returns `managePermission: true` there).

---

### Task 1: Role matrix and permission set

**Files:**
- Create: `internal/rbac/rbac.go`
- Test: `internal/rbac/rbac_test.go`

**Interfaces:**
- Consumes: nothing (pure package).
- Produces:
  - `type Flags struct { Role string; ManageWorkspace, ManageState, ManageModule, ManageProvider, ManageVcs, ManageTemplate, ManageCollection, ManagePolicies, PlanJob, ApproveJob, DestroyJob bool }`
  - `func FlagsFromTeam(store.Team) Flags`, `func FlagsFromAccess(store.Access) Flags`, `func FlagsFromProjectAccess(store.ProjectAccess) Flags`
  - `func CanManageWorkspace(Flags) bool`, `CanManageState`, `CanManageModule`, `CanManageProvider`, `CanManageVcs`, `CanManageTemplate`, `CanManageCollection`, `CanManagePolicies`, `CanPlanJob`, `CanApproveJob`, `CanManageJob`, `CanDestroyJob`
  - `type Permissions struct { ManageState, ManageWorkspace, ManageModule, ManageProvider, ManageVcs, ManageTemplate, ManageCollection, ManageJob, PlanJob, ApproveJob, ManagePermission, ManagePolicies bool }` with camelCase json tags
  - `func (p *Permissions) ApplyTeam(Flags)`, `ApplyAccess(Flags)`, `ApplyProjectAccess(Flags)`, `GrantOwner()`

- [ ] **Step 1: Write the failing table test**

```go
func TestRoleMatrix(t *testing.T) {
	cases := []struct {
		role  string
		check func(rbac.Flags) bool
		admin, write, plan, read, custom bool
	}{ /* one row per check from the Global Constraints table */ }
}

func TestFlagsFromRowsUseEveryBooleans(t *testing.T) { /* all flags true in, all true out, for the three row types */ }

func TestPermissionsAccumulateAndOwnerAddsOnlyTwo(t *testing.T) { /* two teams, one access; owner flips managePermission and managePolicies only */ }
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/rbac -run . -count=1`
Expected: build failure, `package rbac` does not exist.

- [ ] **Step 3: Implement `internal/rbac/rbac.go`**

A `role(f Flags) string` helper lower-cases and trims, mapping `""` to `custom`; each check is a `switch`. `ApplyTeam`/`ApplyAccess`/`ApplyProjectAccess` OR each permission the row type supports (`access` and `project_access` have no module/provider/vcs/template/collection/policies).

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/rbac -count=1`
Expected: `ok`.

- [ ] **Step 5: Commit**

```bash
git add internal/rbac
git commit -m "feat(rbac): resolve upstream team permissions"
```

---

### Task 2: Queries for teams, access and project access by name

**Files:**
- Modify: `internal/store/queries/teams.sql`
- Create: `internal/store/queries/access.sql`
- Generated: `internal/store/*.sql.go`, `internal/store/models.go` (sqlc only)

**Interfaces:**
- Produces (sqlc method names): `ListTeamsByOrganizationAndNames(ctx, ListTeamsByOrganizationAndNamesParams{OrganizationID string, Names []string})`, `ListAccessByWorkspaceAndNames(ctx, ...{WorkspaceID string, Names []string})`, `ListProjectAccessByProjectAndNames(ctx, ...{ProjectID string, Names []string})`, plus `ListAccessByWorkspace(ctx, workspaceID string)` and `ListProjectAccessByProject(ctx, projectID string)` for the tests and later slices.

- [ ] **Step 1: Add the SQL**

```sql
-- name: ListTeamsByOrganizationAndNames :many
SELECT * FROM team
WHERE organization_id = $1 AND name = ANY($2::text[])
ORDER BY name, id;

-- name: ListAccessByWorkspace :many
SELECT * FROM access WHERE workspace_id = $1 ORDER BY name, id;

-- name: ListAccessByWorkspaceAndNames :many
SELECT * FROM access
WHERE workspace_id = $1 AND name = ANY($2::text[])
ORDER BY name, id;

-- name: ListProjectAccessByProject :many
SELECT * FROM project_access WHERE project_id = $1 ORDER BY name, id;

-- name: ListProjectAccessByProjectAndNames :many
SELECT * FROM project_access
WHERE project_id = $1 AND name = ANY($2::text[])
ORDER BY name, id;
```

- [ ] **Step 2: Regenerate and build**

Run: `export PATH=/usr/local/go/bin:/home/amin/go/bin:$PATH && make generate && go build ./...`
Expected: `go build` passes; new methods exist in `internal/store`.

- [ ] **Step 3: Commit**

```bash
git add internal/store/queries/teams.sql internal/store/queries/access.sql
git commit -m "chore(store): add rbac lookup queries"
git add internal/store
git commit -m "chore(store): regenerate sqlc"
```

---

### Task 3: `/access-token/v1` permission endpoints

**Files:**
- Create: `internal/api/accesstoken.go`
- Test: `internal/api/accesstoken_test.go`
- Modify: `cmd/api/main.go` (mount), `internal/api/api.go` (add `AdminGroup` to `Deps` and default it)

**Interfaces:**
- Consumes: `internal/rbac`, `internal/store.Querier`, `auth.IdentityFromContext`, `auth.Authenticate`, `Deps.AuthToken`, `Deps.Verifiers`, `Deps.AdminGroup`.
- Produces: `func NewAccessToken(deps Deps) http.Handler` mounted at `/access-token/v1`; a `Store` interface in the file listing exactly the three lookup methods so tests can fake it.

- [ ] **Step 1: Write the failing handler tests**

Cover, with a fake store and the development token or a PAT verifier:
- `GET /teams/current-teams` → 202 `{"groups":["team-a","TERRAKUBE_ADMIN"]}` in token order.
- Org permissions with a `read` team → only `planJob:false, manageJob:false` style assertions on every field.
- Org permissions with a custom team's flags.
- Org permissions with owner group and an unknown org → `managePermission`/`managePolicies` true, the rest false.
- Workspace permissions: org team + access row + project access row accumulate; an access row whose name is not in groups is ignored.
- Unauthenticated → 401.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/api -run AccessToken -count=1`
Expected: build failure, `NewAccessToken` undefined.

- [ ] **Step 3: Implement the router**

`chi` router with `a.authenticate` reused from `api.go`; handlers fetch the identity, build groups (`identity.Groups`), load teams with `ListTeamsByOrganizationAndNames` (org id parsed with `uuid.Parse`, invalid → 400 documented divergence), accumulate `rbac.Permissions`, apply owner on `slices.Contains(groups, adminGroup)` or `identity.Issuer == "development"`, then write 202 JSON. The workspace endpoint additionally loads the workspace (`GetWorkspace`) to get `project_id`, then the access rows and project access rows **filtered by name** with `= ANY(names)`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/api -run AccessToken -count=1` then `make check`
Expected: `ok` and a clean check.

- [ ] **Step 5: Commit**

```bash
git add internal/api/accesstoken.go internal/api/accesstoken_test.go internal/api/api.go cmd/api/main.go
git commit -m "feat(api): serve the team permission endpoints"
```

---

### Task 4: Record the contract

**Files:**
- Modify: `docs/compat.md` (new subsection under Management API)
- Modify: `docs/upstream/rest/index.tsv` + captured bodies from `tk-upstream-api`

- [ ] **Step 1: Capture the three responses** from the running upstream with a PAT whose `groups` includes a `role=admin` team, plus the unknown-org case.
- [ ] **Step 2: Write the compat subsection** with the 202 status, the 12 fields, the role matrix table, owner behaviour, and the two divergences (name filtering, 400 on a bad UUID).
- [ ] **Step 3: Commit**

```bash
git add docs
git commit -m "docs: record the team permission contract"
```
