-- VCS providers. Columns mirror upstream's 1.4.0 table plus the 2.6/2.8/2.14/
-- 2.18/2.22/2.26/2.33 additions; workspace vcs_id already exists from 0007.
-- The secrets are plaintext until #54, like variable values and ssh keys.
CREATE TABLE vcs (
    id              varchar(36) PRIMARY KEY,
    name            varchar(64)  NOT NULL,
    description     varchar(264) NOT NULL,
    vcs_type        varchar(32)  NOT NULL,
    client_id       varchar(64),
    client_secret   varchar(1024),
    access_token    varchar(4096),
    refresh_token   varchar(4096),
    token_expiration timestamptz,
    status          varchar(16) DEFAULT 'COMPLETED',
    created_date    timestamptz,
    updated_date    timestamptz,
    created_by      varchar(128),
    updated_by      varchar(128),
    callback        varchar(36),
    redirect_url    varchar(256),
    endpoint        varchar(128),
    api_url         varchar(128),
    connection_type varchar(12) DEFAULT 'OAUTH',
    private_key     varchar(4098),
    organization_id varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE
);

CREATE INDEX vcs_organization_idx ON vcs (organization_id);
