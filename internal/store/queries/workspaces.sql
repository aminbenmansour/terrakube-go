-- name: ListWorkspacesByOrganization :many
SELECT * FROM workspace
WHERE organization_id = $1 AND deleted = false
ORDER BY name;

-- name: ListWorkspacesByProject :many
SELECT * FROM workspace
WHERE project_id = $1 AND deleted = false
ORDER BY name;

-- name: ListTagsByOrganization :many
SELECT * FROM tag WHERE organization_id = $1 ORDER BY name;

-- name: ListWorkspaceTagsByWorkspace :many
SELECT * FROM workspacetag WHERE workspace_id = $1;

-- name: GetWorkspace :one
SELECT * FROM workspace WHERE id = $1 AND deleted = false;

-- name: GetWorkspaceByName :one
SELECT * FROM workspace
WHERE organization_id = $1 AND name = $2 AND deleted = false;

-- name: CreateWorkspace :one
INSERT INTO workspace (
    id, name, description, source, branch, folder, terraform_version, iac_type,
    execution_mode, organization_id, project_id, allow_remote_apply,
    global_remote_state, default_template, restrict_apply_to_configured_branch,
    require_approve_to_apply, created_by, module_ssh_key, shared_ids,
    vcs_id, ssh_id, agent_id
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12,
    $13, $14, $15,
    $16, $17, $18, $19,
    $20, $21, $22
)
RETURNING *;

-- name: UpdateWorkspace :one
UPDATE workspace
SET name = $2,
    description = $3,
    source = $4,
    branch = $5,
    folder = $6,
    terraform_version = $7,
    iac_type = $8,
    execution_mode = $9,
    allow_remote_apply = $10,
    global_remote_state = $11,
    default_template = $12,
    restrict_apply_to_configured_branch = $13,
    require_approve_to_apply = $14,
    project_id = $15,
    module_ssh_key = $16,
    shared_ids = $17,
    vcs_id = $18,
    ssh_id = $19,
    agent_id = $20,
    updated_date = now()
WHERE id = $1 AND deleted = false
RETURNING *;

-- name: SetWorkspaceLock :one
UPDATE workspace
SET locked = $2,
    lock_description = $3,
    updated_date = now()
WHERE id = $1 AND deleted = false
RETURNING *;

-- name: SoftDeleteWorkspace :exec
UPDATE workspace SET deleted = true, updated_date = now() WHERE id = $1;
