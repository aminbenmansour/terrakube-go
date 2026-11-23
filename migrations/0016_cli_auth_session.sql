-- Terraform login broker sessions, as upstream's cli_auth_session table. A
-- session tracks one `terraform login`: the CLI's PKCE challenge and loopback
-- redirect, the Dex round trip, the consent choice, and the one-time code the
-- CLI redeems for a personal access token. Only the code hash is stored.
CREATE TABLE cli_auth_session (
    id                varchar(36) PRIMARY KEY,
    status            varchar(20) NOT NULL,
    cli_redirect_uri  varchar(255) NOT NULL,
    cli_code_challenge varchar(128) NOT NULL,
    cli_state         varchar(255) NOT NULL,
    dex_code_verifier varchar(128),
    identity_email    varchar(255),
    identity_name     varchar(255),
    identity_groups   text,
    chosen_days       integer,
    chosen_name       varchar(64),
    auth_code_hash    varchar(64),
    code_expires_at   timestamptz,
    expires_at        timestamptz NOT NULL,
    created_date      timestamptz,
    updated_date      timestamptz,
    created_by        varchar(128),
    updated_by        varchar(128)
);

CREATE INDEX cli_auth_session_auth_code_hash_idx ON cli_auth_session (auth_code_hash);
CREATE INDEX cli_auth_session_expires_at_idx ON cli_auth_session (expires_at);
