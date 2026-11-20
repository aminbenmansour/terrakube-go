-- Tables the upstream UI reads through GraphQL: workspace and project access
-- entries, and notification configurations with their triggers.
--
-- Column names follow the upstream Liquibase changelogs (see
-- api/src/main/resources/db/changelog/local/changelog-2.25.0-workspace-access.xml,
-- changelog-2.32.0-project-access.xml and changelog-2.33.0-notification.xml) so
-- a data migration later does not have to rename anything.

CREATE TABLE access (
    id               varchar(36)  PRIMARY KEY,
    name             varchar(256) NOT NULL,
    manage_workspace boolean      NOT NULL DEFAULT false,
    manage_state     boolean      NOT NULL DEFAULT false,
    manage_job       boolean      NOT NULL DEFAULT false,
    plan_job         boolean      NOT NULL DEFAULT false,
    approve_job      boolean      NOT NULL DEFAULT false,
    destroy_job      boolean      NOT NULL DEFAULT true,
    role             varchar(32),
    workspace_id     varchar(36)  NOT NULL REFERENCES workspace (id) ON DELETE CASCADE
);

CREATE INDEX access_workspace_idx ON access (workspace_id);

CREATE TABLE project_access (
    id               varchar(36) PRIMARY KEY,
    name             varchar(255),
    manage_workspace boolean     NOT NULL DEFAULT false,
    manage_state     boolean     NOT NULL DEFAULT false,
    manage_job       boolean     NOT NULL DEFAULT false,
    plan_job         boolean     NOT NULL DEFAULT false,
    approve_job      boolean     NOT NULL DEFAULT false,
    destroy_job      boolean     NOT NULL DEFAULT true,
    role             varchar(255),
    project_id       varchar(36) NOT NULL REFERENCES project (id) ON DELETE CASCADE
);

CREATE INDEX project_access_project_idx ON project_access (project_id);

CREATE TABLE notification_configuration (
    id              varchar(36)   PRIMARY KEY,
    name            varchar(255)  NOT NULL,
    description     text,
    organization_id varchar(36)   NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
    workspace_id    varchar(36)   REFERENCES workspace (id) ON DELETE CASCADE,
    channel_type    varchar(20)   NOT NULL,
    destination_url varchar(2048) NOT NULL,
    signing_secret  varchar(255),
    active          boolean       NOT NULL DEFAULT true,
    message_style   varchar(20)   NOT NULL DEFAULT 'DETAILED',
    created_date    timestamptz,
    created_by      varchar(128),
    updated_date    timestamptz,
    updated_by      varchar(128)
);

CREATE INDEX notification_configuration_workspace_idx ON notification_configuration (workspace_id);
CREATE INDEX notification_configuration_org_workspace_idx ON notification_configuration (organization_id, workspace_id);

CREATE TABLE notification_trigger (
    id               varchar(36) PRIMARY KEY,
    job_status       varchar(20) NOT NULL,
    configuration_id varchar(36) NOT NULL REFERENCES notification_configuration (id) ON DELETE CASCADE,
    created_date     timestamptz,
    created_by       varchar(128),
    updated_date     timestamptz,
    updated_by       varchar(128)
);

CREATE INDEX notification_trigger_configuration_idx ON notification_trigger (configuration_id);

CREATE TABLE notification_configuration_template (
    configuration_id varchar(36) NOT NULL REFERENCES notification_configuration (id) ON DELETE CASCADE,
    template_id      varchar(36) NOT NULL REFERENCES template (id) ON DELETE CASCADE,
    PRIMARY KEY (configuration_id, template_id)
);
