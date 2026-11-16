-- name: GetJob :one
SELECT * FROM job WHERE id = $1;

-- name: GetStep :one
SELECT * FROM step WHERE id = $1;
