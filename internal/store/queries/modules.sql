-- name: ListModulesByOrganization :many
SELECT * FROM module WHERE organization_id = $1 ORDER BY name, provider, id;

-- name: GetModule :one
SELECT * FROM module WHERE id = $1;

-- name: CreateModule :one
INSERT INTO module (id, name, description, provider, source, organization_id,
                    vcs_id, ssh_id, download_quantity, created_date,
                    updated_date, created_by, updated_by, tag_prefix, folder,
                    latest_version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now(), $10, $11, $12, $13, $14)
RETURNING *;

-- name: UpdateModule :one
UPDATE module
SET name = $2,
    description = $3,
    provider = $4,
    source = $5,
    vcs_id = $6,
    ssh_id = $7,
    tag_prefix = $8,
    folder = $9,
    latest_version = $10,
    updated_date = now(),
    updated_by = $11
WHERE id = $1
RETURNING *;

-- name: DeleteModule :exec
DELETE FROM module WHERE id = $1;

-- name: ListModuleVersions :many
SELECT * FROM module_version WHERE module_id = $1 ORDER BY version, id;

-- name: GetModuleVersion :one
SELECT * FROM module_version WHERE id = $1;

-- name: CreateModuleVersion :one
INSERT INTO module_version (id, module_id, version, commit_info, git_tag,
                            status, deprecation_message)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateModuleVersion :one
UPDATE module_version
SET version = $2,
    commit_info = $3,
    git_tag = $4,
    status = $5,
    deprecation_message = $6
WHERE id = $1
RETURNING *;
