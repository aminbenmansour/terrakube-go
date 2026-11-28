-- name: ListPolicySets :many
SELECT * FROM policy_set ORDER BY name, id;

-- name: ListPolicySetsByOrganization :many
SELECT * FROM policy_set WHERE organization_id = $1 ORDER BY name, id;

-- name: ListPolicySetsByOrganizations :many
SELECT * FROM policy_set
WHERE organization_id = ANY(sqlc.arg(organization_ids)::text[])
ORDER BY name, id;

-- name: GetPolicySet :one
SELECT * FROM policy_set WHERE id = $1;

-- name: CreatePolicySet :one
INSERT INTO policy_set (id, name, description, enforcement_level,
                        shadow_enforcement_level, override_team, global,
                        repository, branch, folder, vcs_id, organization_id,
                        notification_configuration_id, created_date,
                        updated_date, created_by, updated_by, opa_version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, now(),
        now(), $14, $14, $15)
RETURNING *;

-- name: UpdatePolicySet :one
UPDATE policy_set
SET name = $2,
    description = $3,
    enforcement_level = $4,
    shadow_enforcement_level = $5,
    override_team = $6,
    global = $7,
    repository = $8,
    branch = $9,
    folder = $10,
    vcs_id = $11,
    notification_configuration_id = $12,
    opa_version = $13,
    updated_date = now(),
    updated_by = $14
WHERE id = $1
RETURNING *;

-- name: DeletePolicySet :exec
DELETE FROM policy_set WHERE id = $1;

-- name: ListPolicyAttachments :many
SELECT * FROM policy_attachment WHERE policy_set_id = $1 ORDER BY id;

-- name: GetPolicyAttachment :one
SELECT * FROM policy_attachment WHERE id = $1;

-- name: CreatePolicyAttachment :one
INSERT INTO policy_attachment (id, policy_set_id, workspace_id, project_id,
                               tag_id, created_date, updated_date, created_by,
                               updated_by)
VALUES ($1, $2, $3, $4, $5, now(), now(), $6, $6)
RETURNING *;

-- name: UpdatePolicyAttachment :one
UPDATE policy_attachment
SET workspace_id = $2,
    project_id = $3,
    tag_id = $4,
    updated_date = now(),
    updated_by = $5
WHERE id = $1
RETURNING *;

-- name: DeletePolicyAttachment :exec
DELETE FROM policy_attachment WHERE id = $1;

-- name: ListPolicySetParameters :many
SELECT * FROM policy_set_parameter WHERE policy_set_id = $1 ORDER BY parameter_key, id;

-- name: GetPolicySetParameter :one
SELECT * FROM policy_set_parameter WHERE id = $1;

-- name: CreatePolicySetParameter :one
INSERT INTO policy_set_parameter (id, policy_set_id, parameter_key,
                                  parameter_value, description, created_date,
                                  updated_date, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, now(), now(), $6, $6)
RETURNING *;

-- name: UpdatePolicySetParameter :one
UPDATE policy_set_parameter
SET parameter_key = $2,
    parameter_value = $3,
    description = $4,
    updated_date = now(),
    updated_by = $5
WHERE id = $1
RETURNING *;

-- name: DeletePolicySetParameter :exec
DELETE FROM policy_set_parameter WHERE id = $1;

-- name: ListPolicyExemptionsByOrganization :many
SELECT * FROM policy_exemption WHERE organization_id = $1 ORDER BY rule_id, id;

-- name: ListPolicyExemptionsByPolicySet :many
SELECT * FROM policy_exemption WHERE policy_set_id = $1 ORDER BY rule_id, id;

-- name: GetPolicyExemption :one
SELECT * FROM policy_exemption WHERE id = $1;

-- name: CreatePolicyExemption :one
INSERT INTO policy_exemption (id, policy_set_id, rule_id, workspace_id,
                              project_id, ticket_reference, justification,
                              expires_at, organization_id, created_date,
                              updated_date, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now(), $10, $10)
RETURNING *;

-- name: UpdatePolicyExemption :one
UPDATE policy_exemption
SET rule_id = $2,
    workspace_id = $3,
    project_id = $4,
    ticket_reference = $5,
    justification = $6,
    expires_at = $7,
    updated_date = now(),
    updated_by = $8
WHERE id = $1
RETURNING *;

-- name: DeletePolicyExemption :exec
DELETE FROM policy_exemption WHERE id = $1;
