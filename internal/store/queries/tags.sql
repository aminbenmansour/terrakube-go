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

-- name: ListTagsPage :many
SELECT * FROM tag
WHERE organization_id = $1
  AND (sqlc.narg('name')::text IS NULL OR coalesce(lower(name), '') = sqlc.narg('name'))
  AND (sqlc.narg('name_prefix')::text IS NULL OR coalesce(lower(name), '') LIKE sqlc.narg('name_prefix') ESCAPE '\')
ORDER BY coalesce(lower(name), ''), id
LIMIT sqlc.narg('limit')::bigint OFFSET sqlc.arg('offset')::bigint;

-- name: ListTagsPageDesc :many
SELECT * FROM tag
WHERE organization_id = $1
  AND (sqlc.narg('name')::text IS NULL OR coalesce(lower(name), '') = sqlc.narg('name'))
  AND (sqlc.narg('name_prefix')::text IS NULL OR coalesce(lower(name), '') LIKE sqlc.narg('name_prefix') ESCAPE '\')
ORDER BY coalesce(lower(name), '') DESC, id DESC
LIMIT sqlc.narg('limit')::bigint OFFSET sqlc.arg('offset')::bigint;

-- name: CountTagsByName :one
SELECT count(*) FROM tag
WHERE organization_id = $1
  AND (sqlc.narg('name')::text IS NULL OR coalesce(lower(name), '') = sqlc.narg('name'))
  AND (sqlc.narg('name_prefix')::text IS NULL OR coalesce(lower(name), '') LIKE sqlc.narg('name_prefix') ESCAPE '\');
