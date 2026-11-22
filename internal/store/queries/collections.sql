-- name: ListCollectionsByOrganization :many
SELECT * FROM collection WHERE organization_id = $1 ORDER BY priority DESC, name;

-- name: GetCollection :one
SELECT * FROM collection WHERE id = $1;

-- name: CreateCollection :one
INSERT INTO collection (
    id, name, description, priority, created_date, updated_date, organization_id
)
VALUES ($1, $2, $3, $4, now(), now(), $5)
RETURNING *;

-- name: UpdateCollection :one
UPDATE collection
SET name = $2,
    description = $3,
    priority = $4,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCollection :exec
DELETE FROM collection WHERE id = $1;

-- name: ListItemsByCollection :many
SELECT * FROM item WHERE collection_id = $1 ORDER BY item_key;

-- name: GetItem :one
SELECT * FROM item WHERE id = $1;

-- name: CreateItem :one
INSERT INTO item (
    id, item_description, item_key, item_value, sensitive, hcl, item_category,
    created_date, updated_date, collection_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now(), $8)
RETURNING *;

-- name: UpdateItem :one
UPDATE item
SET item_description = $2,
    item_key = $3,
    item_value = $4,
    sensitive = $5,
    hcl = $6,
    item_category = $7,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM item WHERE id = $1;

-- name: ListReferencesByCollection :many
SELECT * FROM reference WHERE collection_id = $1 ORDER BY id;

-- name: ListReferencesByWorkspace :many
SELECT * FROM reference WHERE workspace_id = $1 ORDER BY id;

-- name: GetReference :one
SELECT * FROM reference WHERE id = $1;

-- name: CreateReference :one
INSERT INTO reference (
    id, description, created_date, updated_date, workspace_id, collection_id
)
VALUES ($1, $2, now(), now(), $3, $4)
RETURNING *;

-- name: DeleteReference :exec
DELETE FROM reference WHERE id = $1;
