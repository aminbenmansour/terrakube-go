-- name: ListProvidersByOrganization :many
SELECT * FROM provider WHERE organization_id = $1 ORDER BY name, id;

-- name: GetProvider :one
SELECT * FROM provider WHERE id = $1;

-- name: CreateProvider :one
INSERT INTO provider (id, name, description, organization_id, imported,
                      registry_namespace)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProvider :one
UPDATE provider
SET name = $2,
    description = $3,
    imported = $4,
    registry_namespace = $5
WHERE id = $1
RETURNING *;

-- name: DeleteProvider :exec
DELETE FROM provider WHERE id = $1;

-- name: ListProviderVersions :many
SELECT * FROM version WHERE provider_id = $1 ORDER BY version_number, id;

-- name: ListProviderVersionsByOrganization :many
SELECT version.* FROM version
JOIN provider ON provider.id = version.provider_id
WHERE provider.organization_id = $1
ORDER BY version.version_number, version.id;

-- name: GetProviderVersion :one
SELECT * FROM version WHERE id = $1;

-- name: CreateProviderVersion :one
INSERT INTO version (id, version_number, protocols, provider_id, status,
                     deprecation_message)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProviderVersion :one
UPDATE version
SET version_number = $2,
    protocols = $3,
    status = $4,
    deprecation_message = $5
WHERE id = $1
RETURNING *;

-- name: DeleteProviderVersion :exec
DELETE FROM version WHERE id = $1;

-- name: ListImplementations :many
SELECT * FROM implementation WHERE version_id = $1 ORDER BY os, arch, id;

-- name: ListImplementationsByOrganization :many
SELECT implementation.* FROM implementation
JOIN version ON version.id = implementation.version_id
JOIN provider ON provider.id = version.provider_id
WHERE provider.organization_id = $1
ORDER BY implementation.os, implementation.arch, implementation.id;

-- name: GetImplementation :one
SELECT * FROM implementation WHERE id = $1;

-- name: CreateImplementation :one
INSERT INTO implementation (id, os, arch, filename, download_url, shasums_url,
                            shasums_signature_url, shasum, key_id, ascii_armor,
                            trust_signature, source, source_url, version_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: UpdateImplementation :one
UPDATE implementation
SET os = $2,
    arch = $3,
    filename = $4,
    download_url = $5,
    shasums_url = $6,
    shasums_signature_url = $7,
    shasum = $8,
    key_id = $9,
    ascii_armor = $10,
    trust_signature = $11,
    source = $12,
    source_url = $13
WHERE id = $1
RETURNING *;

-- name: DeleteImplementation :exec
DELETE FROM implementation WHERE id = $1;
