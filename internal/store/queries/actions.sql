-- name: ListActions :many
SELECT * FROM action;

-- name: GetAction :one
SELECT * FROM action WHERE id = $1;

-- name: CreateAction :one
INSERT INTO action (
    id, name, label, category, description, type, action, display_criteria,
    version, active, created_date, updated_date
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now(), now())
RETURNING *;

-- name: UpdateAction :one
UPDATE action
SET name = $2,
    label = $3,
    category = $4,
    description = $5,
    type = $6,
    action = $7,
    display_criteria = $8,
    version = $9,
    active = $10,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAction :exec
DELETE FROM action WHERE id = $1;
