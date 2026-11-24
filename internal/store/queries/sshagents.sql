-- name: ListSSHKeysByOrganization :many
SELECT * FROM ssh WHERE organization_id = $1 ORDER BY name, id;

-- name: GetSSHKey :one
SELECT * FROM ssh WHERE id = $1;

-- name: CreateSSHKey :one
INSERT INTO ssh (
    id, name, description, private_key, ssh_type,
    created_date, updated_date, organization_id
)
VALUES ($1, $2, $3, $4, $5, now(), now(), $6)
RETURNING *;

-- name: UpdateSSHKey :one
UPDATE ssh
SET name = $2,
    description = $3,
    private_key = $4,
    ssh_type = $5,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSSHKey :exec
DELETE FROM ssh WHERE id = $1;

-- name: ListAgentsByOrganization :many
SELECT * FROM agent WHERE organization_id = $1 ORDER BY name, id;

-- name: GetAgent :one
SELECT * FROM agent WHERE id = $1;

-- name: CreateAgent :one
INSERT INTO agent (id, name, url, description, organization_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteAgent :exec
DELETE FROM agent WHERE id = $1;

-- name: CountWorkspacesWithSSHKey :one
SELECT count(*) FROM workspace WHERE ssh_id = $1;

-- name: CountWorkspacesWithAgent :one
SELECT count(*) FROM workspace WHERE agent_id = $1;

-- name: ListSSHKeySecrets :many
SELECT id, private_key FROM ssh;

-- name: SetSSHKeySecret :exec
UPDATE ssh SET private_key = $2 WHERE id = $1 AND private_key = $3;
