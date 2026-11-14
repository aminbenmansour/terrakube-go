-- name: ListProjectsByOrganization :many
SELECT * FROM project WHERE organization_id = $1 ORDER BY name;

-- name: GetProject :one
SELECT * FROM project WHERE id = $1;

-- name: GetProjectByName :one
SELECT * FROM project WHERE organization_id = $1 AND name = $2;

-- name: CreateProject :one
INSERT INTO project (id, name, description, organization_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateProject :one
UPDATE project
SET name = $2,
    description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM project WHERE id = $1;
