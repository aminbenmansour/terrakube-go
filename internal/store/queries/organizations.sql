-- name: GetOrganization :one
SELECT * FROM organization WHERE id = $1;

-- name: GetOrganizationByName :one
SELECT * FROM organization WHERE name = $1;

-- name: CreateOrganization :one
INSERT INTO organization (id, name, description, disabled, execution_mode, icon)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateOrganization :one
UPDATE organization
SET name = $2,
    description = $3,
    disabled = $4,
    execution_mode = $5,
    icon = $6
WHERE id = $1
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organization WHERE id = $1;

-- name: ListOrganizationsPage :many
SELECT * FROM organization
WHERE (sqlc.narg('allowed_ids')::text[] IS NULL OR id = ANY(sqlc.narg('allowed_ids')::text[]))
  AND (sqlc.narg('name')::text IS NULL OR lower(name) = sqlc.narg('name'))
  AND (sqlc.narg('name_prefix')::text IS NULL OR lower(name) LIKE sqlc.narg('name_prefix') ESCAPE '\')
ORDER BY lower(name), id
LIMIT sqlc.narg('limit')::bigint OFFSET sqlc.arg('offset')::bigint;

-- name: ListOrganizationsPageDesc :many
SELECT * FROM organization
WHERE (sqlc.narg('allowed_ids')::text[] IS NULL OR id = ANY(sqlc.narg('allowed_ids')::text[]))
  AND (sqlc.narg('name')::text IS NULL OR lower(name) = sqlc.narg('name'))
  AND (sqlc.narg('name_prefix')::text IS NULL OR lower(name) LIKE sqlc.narg('name_prefix') ESCAPE '\')
ORDER BY lower(name) DESC, id DESC
LIMIT sqlc.narg('limit')::bigint OFFSET sqlc.arg('offset')::bigint;

-- name: CountOrganizationsByName :one
SELECT count(*) FROM organization
WHERE (sqlc.narg('allowed_ids')::text[] IS NULL OR id = ANY(sqlc.narg('allowed_ids')::text[]))
  AND (sqlc.narg('name')::text IS NULL OR lower(name) = sqlc.narg('name'))
  AND (sqlc.narg('name_prefix')::text IS NULL OR lower(name) LIKE sqlc.narg('name_prefix') ESCAPE '\');
