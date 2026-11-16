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
