-- name: GetTagByName :one
SELECT * FROM tag WHERE organization_id = $1 AND name = $2;

-- name: GetTag :one
SELECT * FROM tag WHERE id = $1;

-- name: CreateTag :one
INSERT INTO tag (id, name, created_date, updated_date, organization_id)
VALUES ($1, $2, now(), now(), $3)
RETURNING *;

-- name: UpdateTag :one
UPDATE tag SET name = $2, updated_date = now() WHERE id = $1 RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tag WHERE id = $1;

-- name: CreateWorkspaceTag :one
INSERT INTO workspacetag (id, tag_id, tag_value, workspace_id, created_date, updated_date)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: UpdateWorkspaceTagValue :one
UPDATE workspacetag
SET tag_value = $2,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: GetWorkspaceTag :one
SELECT * FROM workspacetag WHERE id = $1;

-- name: DeleteWorkspaceTag :exec
DELETE FROM workspacetag WHERE id = $1;
