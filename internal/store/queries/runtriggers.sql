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

-- name: ListEnabledRunTriggersBySource :many
SELECT t.* FROM workspace_run_trigger t
JOIN workspace w ON w.id = t.destination_workspace_id
WHERE t.source_workspace_id = $1 AND t.enabled = true AND w.deleted = false
ORDER BY t.created_date, t.id;

-- name: ClaimRunTriggerEvent :one
UPDATE run_trigger_event
SET status = 'PROCESSING',
    attempt_count = attempt_count + 1,
    last_attempt_at = now(),
    updated_date = now()
WHERE id = (
    SELECT id FROM run_trigger_event
    WHERE status = 'PENDING'
      AND (next_attempt_at IS NULL OR next_attempt_at <= now())
    ORDER BY created_date, id
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;

-- name: RecordRunTriggerEvent :execrows
UPDATE run_trigger_event
SET status = $3,
    last_error = $4,
    next_attempt_at = $5,
    updated_date = now()
WHERE id = $1
  AND status = 'PROCESSING'
  AND last_attempt_at = $2;

-- name: ReclaimStuckRunTriggerEvents :execrows
UPDATE run_trigger_event
SET status = 'PENDING', updated_date = now()
WHERE status = 'PROCESSING' AND last_attempt_at < $1 AND attempt_count < $2;

-- name: FailStuckRunTriggerEvents :execrows
UPDATE run_trigger_event
SET status = 'FAILED', last_error = $3, updated_date = now()
WHERE status = 'PROCESSING' AND last_attempt_at < $1 AND attempt_count >= $2;
