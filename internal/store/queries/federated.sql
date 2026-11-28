-- name: ListFederated :many
SELECT * FROM federated ORDER BY name, id;

-- name: GetFederated :one
SELECT * FROM federated WHERE id = $1;

-- name: CreateFederated :one
INSERT INTO federated (id, name, issuer_url, audience, created_date,
                       updated_date, created_by, updated_by)
VALUES ($1, $2, $3, $4, now(), now(), $5, $5)
RETURNING *;

-- name: UpdateFederated :one
UPDATE federated
SET name = $2,
    issuer_url = $3,
    audience = $4,
    updated_date = now(),
    updated_by = $5
WHERE id = $1
RETURNING *;

-- name: DeleteFederated :exec
DELETE FROM federated WHERE id = $1;

-- name: ListFederatedClaims :many
SELECT * FROM federated_claim WHERE federated_id = $1 ORDER BY claim_key, id;

-- name: ListAllFederatedClaims :many
SELECT * FROM federated_claim ORDER BY federated_id, claim_key, id;

-- name: GetFederatedClaim :one
SELECT * FROM federated_claim WHERE id = $1;

-- name: CreateFederatedClaim :one
INSERT INTO federated_claim (id, claim_key, claim_value, federated_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateFederatedClaim :one
UPDATE federated_claim
SET claim_key = $2,
    claim_value = $3
WHERE id = $1
RETURNING *;

-- name: DeleteFederatedClaim :exec
DELETE FROM federated_claim WHERE id = $1;
