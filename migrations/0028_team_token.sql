-- Team-scoped access tokens, mirroring upstream's team_token table
-- (changelog 2.16 through 2.19). A row is the metadata for one signed JWT; the
-- token itself is never stored. last_used_at is this implementation's addition
-- so operators can see stale tokens; upstream keeps no use history.

CREATE TABLE team_token (
    id            varchar(36) PRIMARY KEY,
    days          integer,
    hours         integer,
    minutes       integer,
    description   varchar(64),
    group_name    varchar(64),
    deleted       boolean NOT NULL DEFAULT false,
    last_used_at  timestamp with time zone,
    created_date  timestamp with time zone,
    updated_date  timestamp with time zone,
    created_by    varchar(128),
    updated_by    varchar(128)
);

CREATE INDEX team_token_group_name_idx ON team_token (group_name);
