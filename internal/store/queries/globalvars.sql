-- name: ListGlobalVariablesByOrganization :many
SELECT * FROM globalvar WHERE organization_id = $1 ORDER BY variable_key;

-- name: GetGlobalVariable :one
SELECT * FROM globalvar WHERE id = $1;

-- name: CreateGlobalVariable :one
INSERT INTO globalvar (
    id, variable_key, variable_value, variable_description, variable_category,
    sensitive, hcl, organization_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateGlobalVariable :one
UPDATE globalvar
SET variable_key = $2,
    variable_value = $3,
    variable_description = $4,
    variable_category = $5,
    sensitive = $6,
    hcl = $7
WHERE id = $1
RETURNING *;

-- name: DeleteGlobalVariable :exec
DELETE FROM globalvar WHERE id = $1;

-- name: ListGlobalVariableSecrets :many
SELECT id, variable_value FROM globalvar;

-- name: ListGlobalVariableSecretsAfter :many
SELECT id, variable_value FROM globalvar WHERE id > $1 ORDER BY id LIMIT $2;

-- name: SetGlobalVariableSecret :exec
UPDATE globalvar SET variable_value = $2 WHERE id = $1 AND variable_value = $3;
