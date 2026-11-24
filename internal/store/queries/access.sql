-- Workspace and project access grants. A row's name is the Dex group it
-- grants, so the by-name lookups are the RBAC checks.

-- name: ListAccessByWorkspace :many
SELECT * FROM access WHERE workspace_id = $1 ORDER BY name, id;

-- name: GetAccess :one
SELECT * FROM access WHERE id = $1;

-- name: CreateAccess :one
INSERT INTO access (
    id, name, manage_workspace, manage_state, manage_job, plan_job,
    approve_job, destroy_job, role, workspace_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateAccess :one
UPDATE access SET
    name = $2,
    manage_workspace = $3,
    manage_state = $4,
    manage_job = $5,
    plan_job = $6,
    approve_job = $7,
    destroy_job = $8,
    role = $9
WHERE id = $1
RETURNING *;

-- name: DeleteAccess :exec
DELETE FROM access WHERE id = $1;

-- name: ListProjectAccessByProject :many
SELECT * FROM project_access WHERE project_id = $1 ORDER BY name, id;

-- name: GetProjectAccess :one
SELECT * FROM project_access WHERE id = $1;

-- name: CreateProjectAccess :one
INSERT INTO project_access (
    id, name, manage_workspace, manage_state, manage_job, plan_job,
    approve_job, destroy_job, role, project_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateProjectAccess :one
UPDATE project_access SET
    name = $2,
    manage_workspace = $3,
    manage_state = $4,
    manage_job = $5,
    plan_job = $6,
    approve_job = $7,
    destroy_job = $8,
    role = $9
WHERE id = $1
RETURNING *;

-- name: DeleteProjectAccess :exec
DELETE FROM project_access WHERE id = $1;

-- name: ListAccessByWorkspaceAndNames :many
SELECT * FROM access
WHERE workspace_id = $1 AND name = ANY(sqlc.arg(names)::text[])
ORDER BY name, id;

-- name: ListProjectAccessByProjectAndNames :many
SELECT * FROM project_access
WHERE project_id = $1 AND name = ANY(sqlc.arg(names)::text[])
ORDER BY name, id;
