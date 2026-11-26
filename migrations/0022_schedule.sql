-- Workspace schedules: the Schedules tab creates these and upstream's scheduler
-- turns each cron tick into a run. Columns follow
-- api/src/main/resources/db/changelog/local/changelog-1.6.xml and
-- changelog-2.0.xml; `tcl` is deprecated upstream (schedules use a template)
-- but kept so a data migration lines up.

CREATE TABLE schedule (
    id                 varchar(36) PRIMARY KEY,
    cron               varchar(128),
    description        varchar(256),
    tcl                text,
    template_reference varchar(36),
    enabled            boolean     NOT NULL DEFAULT true,
    workspace_id       varchar(36) NOT NULL REFERENCES workspace (id) ON DELETE CASCADE,
    created_date       timestamptz,
    updated_date       timestamptz,
    created_by         varchar(128),
    updated_by         varchar(128)
);

CREATE INDEX schedule_workspace_idx ON schedule (workspace_id);
