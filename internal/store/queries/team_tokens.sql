-- name: CreateTeamToken :one
INSERT INTO team_token (id, days, hours, minutes, description, group_name, deleted,
                        created_date, updated_date, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, false, now(), now(), $7, $8)
RETURNING *;

-- name: GetTeamToken :one
SELECT * FROM team_token WHERE id = $1;

-- name: ListTeamTokensByGroups :many
SELECT * FROM team_token
WHERE group_name = ANY($1::text[])
ORDER BY created_date DESC, id;

-- name: SoftDeleteTeamToken :exec
UPDATE team_token
SET deleted = true,
    updated_date = now(),
    updated_by = $2
WHERE id = $1;

-- name: StampTeamTokenLastUsed :exec
UPDATE team_token
SET last_used_at = now()
WHERE id = $1
  AND (last_used_at IS NULL OR last_used_at < now() - interval '1 hour');
