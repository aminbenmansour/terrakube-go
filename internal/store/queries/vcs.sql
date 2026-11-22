-- name: ListVCSByOrganization :many
SELECT * FROM vcs WHERE organization_id = $1 ORDER BY name, id;

-- name: GetVCS :one
SELECT * FROM vcs WHERE id = $1;

-- name: CreateVCS :one
INSERT INTO vcs (id, name, description, vcs_type, client_id, client_secret,
                 access_token, status, created_date, updated_date, callback,
                 redirect_url, endpoint, api_url, connection_type, private_key,
                 organization_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), now(), $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: UpdateVCS :one
UPDATE vcs
SET name = $2,
    description = $3,
    vcs_type = $4,
    client_id = $5,
    client_secret = $6,
    access_token = $7,
    status = $8,
    callback = $9,
    redirect_url = $10,
    endpoint = $11,
    api_url = $12,
    connection_type = $13,
    private_key = $14,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteVCS :exec
DELETE FROM vcs WHERE id = $1;

-- Deleting a provider in use unlinks its workspaces instead of failing, which
-- is what upstream does through the mapped relationship.
-- name: ClearWorkspaceVCS :exec
UPDATE workspace SET vcs_id = NULL WHERE vcs_id = $1;

-- name: ListWorkspaceIDsByVCS :many
SELECT id FROM workspace WHERE vcs_id = $1 ORDER BY name, id;
