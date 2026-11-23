-- name: CreateCliAuthSession :one
INSERT INTO cli_auth_session (id, status, cli_redirect_uri, cli_code_challenge,
                              cli_state, dex_code_verifier, expires_at,
                              created_date, updated_date)
VALUES ($1, 'PENDING_IDP', $2, $3, $4, $5, $6, now(), now())
RETURNING *;

-- name: GetCliAuthSession :one
SELECT * FROM cli_auth_session WHERE id = $1;

-- name: SetCliAuthSessionIdentity :exec
UPDATE cli_auth_session
SET status = 'PENDING_CONSENT',
    identity_email = $2,
    identity_name = $3,
    identity_groups = $4,
    dex_code_verifier = NULL,
    updated_date = now()
WHERE id = $1 AND status = 'PENDING_IDP';

-- name: FailCliAuthSession :exec
UPDATE cli_auth_session
SET status = 'FAILED', dex_code_verifier = NULL, updated_date = now()
WHERE id = $1;

-- name: IssueCliAuthCode :exec
UPDATE cli_auth_session
SET status = 'CODE_ISSUED',
    auth_code_hash = $2,
    chosen_days = $3,
    chosen_name = $4,
    code_expires_at = $5,
    updated_date = now()
WHERE id = $1 AND status = 'PENDING_CONSENT';

-- name: DenyCliAuthSession :exec
UPDATE cli_auth_session
SET status = 'DENIED', updated_date = now()
WHERE id = $1 AND status = 'PENDING_CONSENT';

-- name: GetCliAuthSessionByCodeHash :one
SELECT * FROM cli_auth_session WHERE auth_code_hash = $1;

-- name: ClaimCliAuthCode :execrows
UPDATE cli_auth_session
SET status = 'EXCHANGED', updated_date = now()
WHERE id = $1 AND status = 'CODE_ISSUED' AND code_expires_at > now();

-- name: DeleteExpiredCliAuthSessions :execrows
DELETE FROM cli_auth_session WHERE expires_at < now();
