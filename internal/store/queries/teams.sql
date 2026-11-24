-- name: ListTeamsByOrganization :many
SELECT * FROM team WHERE organization_id = $1 ORDER BY name, id;

-- name: GetTeam :one
SELECT * FROM team WHERE id = $1;

-- name: CreateTeam :one
INSERT INTO team (id, name, manage_workspace, manage_module, manage_provider,
                  manage_vcs, manage_template, manage_state, manage_collection,
                  manage_job, manage_policies, plan_job, approve_job, destroy_job,
                  role, organization_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING *;

-- name: UpdateTeam :one
UPDATE team
SET name = $2,
    manage_workspace = $3,
    manage_module = $4,
    manage_provider = $5,
    manage_vcs = $6,
    manage_template = $7,
    manage_state = $8,
    manage_collection = $9,
    manage_job = $10,
    manage_policies = $11,
    plan_job = $12,
    approve_job = $13,
    destroy_job = $14,
    role = $15
WHERE id = $1
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM team WHERE id = $1;

-- name: ListTeamsByOrganizationAndNames :many
SELECT * FROM team
WHERE organization_id = $1 AND name = ANY(sqlc.arg(names)::text[])
ORDER BY name, id;
