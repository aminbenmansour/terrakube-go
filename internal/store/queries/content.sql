-- name: CreateContent :one
INSERT INTO content (id, workspace_id, auto_queue_runs, speculative, status, source)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetContent :one
SELECT * FROM content WHERE id = $1;

-- name: UpdateContentStatus :one
UPDATE content
SET status = $2
WHERE id = $1
RETURNING *;
