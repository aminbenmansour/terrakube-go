-- Configuration versions uploaded through the TFE API.
--
-- Mirrors the upstream `content` table from changelog-2.10.0-remote-state.xml.
-- Configuration tarballs are stored in object storage, not in Postgres.

CREATE TABLE content (
    id              varchar(36) PRIMARY KEY,
    auto_queue_runs boolean     NOT NULL DEFAULT false,
    source          varchar(32),
    speculative     boolean,
    status          varchar(32),
    workspace_id    varchar(36) REFERENCES workspace (id) ON DELETE CASCADE
);

CREATE INDEX content_workspace_idx ON content (workspace_id);
