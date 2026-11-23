-- name: CreatePat :one
INSERT INTO pat (id, days, description, source, deleted, created_date, updated_date, created_by, updated_by)
VALUES ($1, $2, $3, $4, false, now(), now(), $5, $6)
RETURNING *;

-- name: GetPat :one
SELECT * FROM pat WHERE id = $1;

-- name: ListPatsByCreator :many
SELECT * FROM pat
WHERE created_by = $1 AND deleted = false
ORDER BY created_date DESC, id;

-- name: SoftDeletePat :exec
UPDATE pat
SET deleted = true, updated_date = now(), updated_by = $2
WHERE id = $1;

-- name: StampPatLastUsed :exec
UPDATE pat
SET last_used_at = now()
WHERE id = $1
  AND (last_used_at IS NULL OR last_used_at < now() - interval '1 hour');
