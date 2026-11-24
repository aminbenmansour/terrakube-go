-- Workspace and project access grants. A row's name is the Dex group it
-- grants, so the by-name lookups are the RBAC checks.

-- name: ListAccessByWorkspace :many
SELECT * FROM access WHERE workspace_id = $1 ORDER BY name, id;

-- name: ListAccessByWorkspaceAndNames :many
SELECT * FROM access
WHERE workspace_id = $1 AND name = ANY(sqlc.arg(names)::text[])
ORDER BY name, id;

-- name: ListProjectAccessByProject :many
SELECT * FROM project_access WHERE project_id = $1 ORDER BY name, id;

-- name: ListProjectAccessByProjectAndNames :many
SELECT * FROM project_access
WHERE project_id = $1 AND name = ANY(sqlc.arg(names)::text[])
ORDER BY name, id;
