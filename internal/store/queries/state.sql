-- name: GetJob :one
SELECT * FROM job WHERE id = $1;

-- name: ListJobsByOrganization :many
SELECT * FROM job
WHERE organization_id = $1 AND deleted = false
ORDER BY id DESC;

-- name: ListJobsByWorkspace :many
SELECT * FROM job
WHERE workspace_id = $1 AND deleted = false
ORDER BY id DESC;

-- name: CountJobsByStatus :many
SELECT status, count(*) AS total
FROM job
WHERE deleted = false
GROUP BY status;

-- name: GetStep :one
SELECT * FROM step WHERE id = $1;

-- name: GetHistory :one
SELECT * FROM history WHERE id = $1;

-- name: GetLatestHistory :one
SELECT * FROM history
WHERE workspace_id = $1
ORDER BY created_date DESC, id DESC
LIMIT 1;

-- name: GetHistoryByJobReference :one
SELECT * FROM history
WHERE workspace_id = $1 AND job_reference = $2
ORDER BY created_date DESC, id DESC
LIMIT 1;

-- name: GetHistoryBefore :one
SELECT * FROM history
WHERE workspace_id = $1 AND created_date < $2
ORDER BY created_date DESC, id DESC
LIMIT 1;

-- name: CreateHistory :one
INSERT INTO history (id, workspace_id, job_reference, output, serial, md5, lineage, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: CreateJob :one
INSERT INTO job (organization_id, workspace_id, command, status, output, via, auto_apply, refresh, plan_changes, created_by, tcl, template_reference,
                 triggered_by_job_id, cascade_depth, commit_id, override_branch)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING *;

-- name: GetJobByWorkspaceAndTrigger :one
SELECT * FROM job
WHERE workspace_id = $1 AND triggered_by_job_id = $2
ORDER BY id DESC
LIMIT 1;

-- name: DeleteJob :exec
-- Only for a just-created dependent whose steps could not all be written; a
-- normal job is soft-deleted so its history stays visible.
DELETE FROM job WHERE id = $1;

-- name: CreateStep :one
INSERT INTO step (id, job_id, step_number, name, status, log_status, output)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateJobStatus :one
-- Completing a job and recording its dispatch work are one statement, so a
-- crash between them cannot lose either. The run trigger event is only written
-- for a completed job with an enabled edge to a live destination; the unique
-- job_id makes a repeated completion a no-op.
--
-- A notification row is written for every active configuration in scope (the
-- workspace's own, or an organization-wide one) whose trigger names the new
-- status and whose template filter matches. The payload is rendered at
-- dispatch time, so the row carries only the status it was written for. A
-- "failed" row is held for a short grace period: a lost API-to-executor
-- response flips the job to failed briefly before the executor's own callback
-- corrects it, and the dispatcher drops the row if the live status moved on.
WITH updated AS (
    UPDATE job
    SET status = $2,
        updated_date = now()
    WHERE job.id = $1
    RETURNING job.id, job.workspace_id, job.status
), event AS (
    INSERT INTO run_trigger_event (id, job_id, status)
    SELECT gen_random_uuid()::text, updated.id, 'PENDING'
    FROM updated
    WHERE updated.status = 'completed'
      AND EXISTS (
          SELECT 1
          FROM workspace_run_trigger t
          JOIN workspace w ON w.id = t.destination_workspace_id
          WHERE t.source_workspace_id = updated.workspace_id
            AND t.enabled
            AND w.deleted = false
      )
    ON CONFLICT (job_id) DO NOTHING
    RETURNING job_id
), notification AS (
    INSERT INTO notification_outbox (id, job_id, configuration_id, status, job_status, next_attempt_at)
    SELECT gen_random_uuid()::text, updated.id, c.id, 'PENDING', updated.status,
           CASE
               WHEN updated.status = 'failed' THEN now() + make_interval(secs => 90)
           END
    FROM updated
    JOIN job j ON j.id = updated.id
    JOIN workspace w ON w.id = j.workspace_id
    JOIN notification_configuration c
      ON c.active
     AND (c.workspace_id = w.id
          OR (c.workspace_id IS NULL AND c.organization_id = w.organization_id))
    JOIN notification_trigger t
      ON t.configuration_id = c.id
     AND t.job_status = updated.status
    WHERE NOT EXISTS (
              SELECT 1 FROM notification_configuration_template ct
              WHERE ct.configuration_id = c.id
          )
       OR EXISTS (
              SELECT 1 FROM notification_configuration_template ct
              WHERE ct.configuration_id = c.id
                AND ct.template_id = j.template_reference
          )
    ON CONFLICT DO NOTHING
    RETURNING id
)
SELECT * FROM job WHERE job.id = $1;

-- name: SetJobOverrideSource :one
UPDATE job
SET override_source = $2,
    override_branch = $3,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: UpdateJobOutput :exec
UPDATE job
SET output = $2,
    updated_date = now()
WHERE id = $1;

-- name: UpdateJobPlanChanges :one
UPDATE job
SET plan_changes = $2,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: ApproveJob :one
UPDATE job
SET approved_by = $2,
    approved_at = now(),
    status = 'pending',
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: SetWorkspaceLastJob :exec
UPDATE workspace
SET last_job_status = $2,
    last_job_date = now(),
    updated_date = now()
WHERE id = $1;

-- name: ListActiveJobsByWorkspace :many
SELECT * FROM job
WHERE workspace_id = $1
  AND status IN ('pending', 'running', 'queue', 'waitingApproval')
ORDER BY id DESC;

-- name: ListActiveJobsByOrganization :many
SELECT j.*
FROM job j
JOIN organization o ON o.id = j.organization_id
WHERE o.name = $1
  AND j.status IN ('pending', 'running', 'queue', 'waitingApproval')
ORDER BY j.id ASC;

-- name: ListArchivesByHistory :many
SELECT * FROM temp_archive WHERE history_id = $1 ORDER BY type;

-- name: UpdateHistoryOutput :exec
UPDATE history
SET output = $2,
    updated_date = now()
WHERE id = $1;

-- name: GetArchive :one
SELECT ta.id, ta.type, ta.history_id, h.workspace_id, w.organization_id
FROM temp_archive ta
JOIN history h ON h.id = ta.history_id
JOIN workspace w ON w.id = h.workspace_id
WHERE ta.id = $1;

-- name: CreateArchive :one
INSERT INTO temp_archive (id, type, history_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteArchive :exec
DELETE FROM temp_archive WHERE id = $1;

-- name: ListStepsByJob :many
SELECT * FROM step WHERE job_id = $1 ORDER BY step_number ASC;

-- name: GetStepByJobAndNumber :one
SELECT * FROM step
WHERE job_id = $1 AND step_number = $2
LIMIT 1;

-- name: GetPendingStepByJob :one
-- "running" is included so an executor that reclaims a crashed job resumes the
-- interrupted step instead of treating the job as finished.
SELECT * FROM step
WHERE job_id = $1 AND status IN ('pending', 'queue', 'running')
ORDER BY step_number ASC
LIMIT 1;

-- name: UpdateStepStatus :one
UPDATE step
SET status = $2
WHERE id = $1
RETURNING *;

-- name: UpdateStepOutput :exec
UPDATE step
SET output = $2
WHERE id = $1;
