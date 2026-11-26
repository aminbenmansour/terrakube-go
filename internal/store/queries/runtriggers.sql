-- name: GetRunTrigger :one
SELECT * FROM workspace_run_trigger WHERE id = $1;

-- name: ListRunTriggersByOrganization :many
SELECT * FROM workspace_run_trigger WHERE organization_id = $1 ORDER BY created_date, id;

-- name: CountEnabledRunTriggersBySource :one
SELECT count(*) FROM workspace_run_trigger
WHERE source_workspace_id = $1 AND enabled = true AND id <> $2;

-- name: CreateRunTrigger :one
INSERT INTO workspace_run_trigger (id, source_workspace_id, destination_workspace_id,
                                   template_id, organization_id, enabled,
                                   created_date, updated_date)
VALUES ($1, $2, $3, $4, $5, $6, now(), now())
RETURNING *;

-- name: UpdateRunTrigger :one
UPDATE workspace_run_trigger
SET source_workspace_id = $2,
    destination_workspace_id = $3,
    template_id = $4,
    organization_id = $5,
    enabled = $6,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteRunTrigger :exec
DELETE FROM workspace_run_trigger WHERE id = $1;
