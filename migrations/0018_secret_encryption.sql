-- Encrypted values carry an `enc:v1:<key>:<nonce>:` envelope and base64
-- ciphertext, which can exceed the upstream varchar widths. Widen only the
-- secret-bearing columns; names and tables stay upstream-compatible.
ALTER TABLE vcs
    ALTER COLUMN client_secret TYPE text,
    ALTER COLUMN access_token TYPE text,
    ALTER COLUMN refresh_token TYPE text,
    ALTER COLUMN private_key TYPE text;

ALTER TABLE notification_configuration
    ALTER COLUMN signing_secret TYPE text;
