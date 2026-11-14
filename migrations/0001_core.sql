-- Core schema for Terrakube Go.
--
-- Table and column names intentionally follow the upstream Terrakube
-- (Liquibase) schema so that existing clients, the React UI and eventually a
-- data migration keep working unchanged. Columns that upstream adds later are
-- added by subsequent migrations.

CREATE TABLE organization (
    id             varchar(36) PRIMARY KEY,
    name           varchar(64)  NOT NULL,
    description    varchar(256),
    disabled       boolean      NOT NULL DEFAULT false,
    execution_mode varchar(32)  NOT NULL DEFAULT 'remote',
    icon           varchar(256)
);

CREATE UNIQUE INDEX organization_name_uq ON organization (name);

CREATE TABLE project (
    id              varchar(36) PRIMARY KEY,
    name            varchar(64) NOT NULL,
    description     varchar(256),
    organization_id varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE
);

CREATE INDEX project_organization_idx ON project (organization_id);

CREATE TABLE workspace (
    id                               varchar(36) PRIMARY KEY,
    name                             varchar(64)  NOT NULL,
    description                      varchar(256),
    source                           varchar(512),
    branch                           varchar(512),
    folder                           varchar(256),
    last_job_status                  varchar(32)  NOT NULL DEFAULT 'NeverExecuted',
    last_job_date                    timestamptz,
    locked                           boolean      NOT NULL DEFAULT false,
    deleted                          boolean      NOT NULL DEFAULT false,
    allow_remote_apply               boolean      NOT NULL DEFAULT false,
    global_remote_state              boolean      NOT NULL DEFAULT true,
    shared_ids                       text,
    default_template                 varchar(36),
    restrict_apply_to_configured_branch boolean   NOT NULL DEFAULT false,
    require_approve_to_apply         boolean      NOT NULL DEFAULT false,
    lock_description                 varchar(256),
    iac_type                         varchar(32)  NOT NULL DEFAULT 'terraform',
    module_ssh_key                   text,
    terraform_version                varchar(32),
    execution_mode                   varchar(32),
    policy_compliance_status         varchar(32)  NOT NULL DEFAULT 'UNKNOWN',
    organization_id                  varchar(36)  NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
    project_id                       varchar(36)  REFERENCES project (id) ON DELETE SET NULL,
    created_date                     timestamptz  NOT NULL DEFAULT now(),
    created_by                       varchar(128),
    updated_date                     timestamptz,
    updated_by                       varchar(128)
);

CREATE INDEX workspace_organization_idx ON workspace (organization_id);
CREATE INDEX workspace_project_idx ON workspace (project_id);
CREATE UNIQUE INDEX workspace_org_name_uq ON workspace (organization_id, name) WHERE deleted = false;

CREATE TABLE variable (
    id                   varchar(36) PRIMARY KEY,
    variable_key         varchar(128) NOT NULL,
    variable_value       text         NOT NULL,
    variable_description varchar(256),
    variable_category    varchar(32)  NOT NULL DEFAULT 'TERRAFORM',
    sensitive            boolean      NOT NULL DEFAULT false,
    hcl                  boolean      NOT NULL DEFAULT false,
    incomplete           boolean      NOT NULL DEFAULT false,
    workspace_id         varchar(36)  NOT NULL REFERENCES workspace (id) ON DELETE CASCADE
);

CREATE INDEX variable_workspace_idx ON variable (workspace_id);

CREATE TABLE tag (
    id              varchar(36) PRIMARY KEY,
    name            varchar(128) NOT NULL,
    organization_id varchar(36)  NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
    created_date    timestamptz  NOT NULL DEFAULT now(),
    created_by      varchar(128),
    updated_date    timestamptz,
    updated_by      varchar(128)
);

CREATE UNIQUE INDEX tag_org_name_uq ON tag (organization_id, name);

CREATE TABLE workspacetag (
    id              varchar(36) PRIMARY KEY,
    tag_id          varchar(36) NOT NULL REFERENCES tag (id) ON DELETE CASCADE,
    tag_value       varchar(256),
    workspace_id    varchar(36) NOT NULL REFERENCES workspace (id) ON DELETE CASCADE,
    created_date    timestamptz NOT NULL DEFAULT now(),
    created_by      varchar(128),
    updated_date    timestamptz,
    updated_by      varchar(128),
    CONSTRAINT const_workspace_tag UNIQUE (tag_id, workspace_id)
);

CREATE INDEX workspacetag_workspace_idx ON workspacetag (workspace_id);

CREATE TABLE template (
    id              varchar(36) PRIMARY KEY,
    name            varchar(64) NOT NULL,
    description     varchar(256),
    version         varchar(32),
    tcl             text,
    organization_id varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE
);

CREATE INDEX template_organization_idx ON template (organization_id);

CREATE TABLE job (
    id                   bigserial PRIMARY KEY,
    organization_id      varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
    workspace_id         varchar(36) NOT NULL REFERENCES workspace (id) ON DELETE CASCADE,
    command              varchar(32),
    status               varchar(32) NOT NULL DEFAULT 'pending',
    output               text,
    template_reference   varchar(36),
    tcl                  text,
    via                  varchar(32) NOT NULL DEFAULT 'UI',
    auto_apply           boolean     NOT NULL DEFAULT false,
    deleted              boolean     NOT NULL DEFAULT false,
    comments             text,
    commit_id            varchar(64),
    approved_by          varchar(128),
    approved_at          timestamptz,
    plan_changes         boolean     NOT NULL DEFAULT true,
    refresh              boolean     NOT NULL DEFAULT true,
    created_date         timestamptz NOT NULL DEFAULT now(),
    created_by           varchar(128),
    updated_date         timestamptz,
    updated_by           varchar(128)
);

CREATE INDEX job_organization_idx ON job (organization_id);
CREATE INDEX job_workspace_idx ON job (workspace_id);

CREATE TABLE step (
    id          varchar(36) PRIMARY KEY,
    job_id      bigint      NOT NULL REFERENCES job (id) ON DELETE CASCADE,
    step_number integer     NOT NULL,
    name        varchar(128),
    status      varchar(32),
    log_status  varchar(32),
    output      text
);

CREATE INDEX step_job_idx ON step (job_id);

CREATE TABLE history (
    id              varchar(36) PRIMARY KEY,
    workspace_id    varchar(36) NOT NULL REFERENCES workspace (id) ON DELETE CASCADE,
    job_reference   varchar(36),
    output          text,
    serial          integer     NOT NULL DEFAULT 1,
    md5             varchar(64) NOT NULL DEFAULT '0',
    lineage         varchar(128) NOT NULL DEFAULT '0',
    created_date    timestamptz NOT NULL DEFAULT now(),
    created_by      varchar(128),
    updated_date    timestamptz,
    updated_by      varchar(128)
);

CREATE INDEX history_workspace_idx ON history (workspace_id);
