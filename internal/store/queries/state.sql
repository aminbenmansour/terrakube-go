-- name: GetJob :one
SELECT * FROM job WHERE id = $1;

-- name: GetStep :one
SELECT * FROM step WHERE id = $1;

-- name: GetHistory :one
SELECT * FROM history WHERE id = $1;

-- name: GetLatestHistory :one
SELECT * FROM history
WHERE workspace_id = $1
ORDER BY created_date DESC, id DESC
LIMIT 1;

-- name: CreateHistory :one
INSERT INTO history (id, workspace_id, job_reference, output, serial, md5, lineage, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: CreateJob :one
INSERT INTO job (organization_id, workspace_id, command, status, output, via, auto_apply, refresh, plan_changes, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: CreateStep :one
INSERT INTO step (id, job_id, step_number, name, status, log_status, output)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateJobStatus :one
UPDATE job
SET status = $2,
    updated_date = now()
WHERE id = $1
RETURNING *;

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
SELECT * FROM step
WHERE job_id = $1 AND status IN ('pending', 'queue')
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
