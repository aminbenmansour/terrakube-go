-- name: GetTagByName :one
SELECT * FROM tag WHERE organization_id = $1 AND name = $2;

-- name: CreateTag :one
INSERT INTO tag (id, name, organization_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateWorkspaceTag :one
INSERT INTO workspacetag (id, tag_id, tag_value, workspace_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateWorkspaceTagValue :one
UPDATE workspacetag
SET tag_value = $2,
    updated_date = now()
WHERE id = $1
RETURNING *;
