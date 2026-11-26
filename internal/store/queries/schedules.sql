-- name: ListSchedulesByWorkspace :many
SELECT * FROM schedule WHERE workspace_id = $1 ORDER BY created_date, id;

-- name: GetSchedule :one
SELECT * FROM schedule WHERE id = $1;

-- name: CreateSchedule :one
INSERT INTO schedule (id, cron, description, template_reference, enabled,
                      workspace_id, created_date, updated_date)
VALUES ($1, $2, $3, $4, $5, $6, now(), now())
RETURNING *;

-- name: UpdateSchedule :one
UPDATE schedule
SET cron = $2,
    description = $3,
    template_reference = $4,
    enabled = $5,
    updated_date = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSchedule :exec
DELETE FROM schedule WHERE id = $1;
