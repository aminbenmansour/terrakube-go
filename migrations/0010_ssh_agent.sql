-- SSH keys and self-hosted agent pools. Columns mirror upstream's 2.6.0 and
-- 2.20.0 changelogs; the workspace ssh_id/agent_id links were added in 0007.
CREATE TABLE ssh (
    id              varchar(36) PRIMARY KEY,
    name            varchar(64),
    description     varchar(64),
    private_key     text,
    ssh_type        varchar(32),
    created_date    timestamptz,
    updated_date    timestamptz,
    created_by      varchar(128),
    updated_by      varchar(128),
    organization_id varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE
);

CREATE INDEX ssh_organization_idx ON ssh (organization_id);

CREATE TABLE agent (
    id              varchar(36) PRIMARY KEY,
    name            varchar(64),
    url             varchar(256),
    description     varchar(64),
    organization_id varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE
);

CREATE INDEX agent_organization_idx ON agent (organization_id);
