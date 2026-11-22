-- name: ListTemplatesByOrganization :many
SELECT * FROM template WHERE organization_id = $1 ORDER BY name, id;

-- name: GetTemplate :one
SELECT * FROM template WHERE id = $1;

-- name: CreateTemplate :one
INSERT INTO template (id, name, description, version, tcl, created_date,
                      updated_date, organization_id)
VALUES ($1, $2, $3, $4, $5, now(), now(), $6)
RETURNING *;

-- name: UpdateTemplate :one
UPDATE template
SET name = $2,
    description = $3,
    version = $4,
    tcl = $5,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTemplate :exec
DELETE FROM template WHERE id = $1;
