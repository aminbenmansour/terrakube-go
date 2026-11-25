-- name: ListVariablesByWorkspace :many
SELECT * FROM variable WHERE workspace_id = $1 ORDER BY variable_key;

-- name: GetVariable :one
SELECT * FROM variable WHERE id = $1;

-- name: CreateVariable :one
INSERT INTO variable (
    id, variable_key, variable_value, variable_description, variable_category,
    sensitive, hcl, incomplete, workspace_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateVariable :one
UPDATE variable
SET variable_key = $2,
    variable_value = $3,
    variable_description = $4,
    variable_category = $5,
    sensitive = $6,
    hcl = $7,
    incomplete = $8
WHERE id = $1
RETURNING *;

-- name: DeleteVariable :exec
DELETE FROM variable WHERE id = $1;

-- name: ListVariableSecrets :many
SELECT id, variable_value FROM variable;

-- name: ListVariableSecretsAfter :many
SELECT id, variable_value FROM variable WHERE id > $1 ORDER BY id LIMIT $2;

-- name: SetVariableSecret :exec
UPDATE variable SET variable_value = $2 WHERE id = $1 AND variable_value = $3;
