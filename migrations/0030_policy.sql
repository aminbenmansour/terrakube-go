-- Policy sets, their attachments and parameters, and policy exemptions,
-- mirroring upstream's 2.34.0 OPA policy tables. The evaluation tables
-- (policy_evaluation, policy_override) and workspace.policy_compliance_status
-- belong to the policy engine work, not the Settings management surface.
--
-- Child rows cascade when a policy set is deleted, which is what upstream's
-- JPA orphanRemoval does; the foreign keys in the upstream schema do not.
-- Attachments and exemptions keep plain references to workspaces, projects
-- and tags, as upstream does.

CREATE TABLE policy_set (
    id                            varchar(36) PRIMARY KEY,
    name                          varchar(255) NOT NULL,
    description                   text,
    enforcement_level             varchar(50) NOT NULL DEFAULT 'HARD_MANDATORY',
    shadow_enforcement_level      varchar(50),
    override_team                 varchar(255),
    global                        boolean NOT NULL DEFAULT false,
    repository                    varchar(1024),
    branch                        varchar(255) DEFAULT 'main',
    folder                        varchar(1024) DEFAULT '/',
    vcs_id                        varchar(36) REFERENCES vcs (id),
    organization_id               varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
    notification_configuration_id varchar(36) REFERENCES notification_configuration (id),
    created_date                  timestamptz,
    updated_date                  timestamptz,
    created_by                    varchar(255),
    updated_by                    varchar(255),
    opa_version                   varchar(50)
);

CREATE INDEX idx_policy_set_org ON policy_set (organization_id);

CREATE TABLE policy_attachment (
    id            varchar(36) PRIMARY KEY,
    policy_set_id varchar(36) NOT NULL REFERENCES policy_set (id) ON DELETE CASCADE,
    workspace_id  varchar(36) REFERENCES workspace (id),
    project_id    varchar(36) REFERENCES project (id),
    tag_id        varchar(36) REFERENCES tag (id),
    created_date  timestamptz,
    updated_date  timestamptz,
    created_by    varchar(255),
    updated_by    varchar(255)
);

CREATE INDEX idx_attachment_policy_set ON policy_attachment (policy_set_id);
CREATE INDEX idx_attachment_tag_id ON policy_attachment (tag_id);

CREATE TABLE policy_set_parameter (
    id              varchar(36) PRIMARY KEY,
    policy_set_id   varchar(36) NOT NULL REFERENCES policy_set (id) ON DELETE CASCADE,
    parameter_key   varchar(255) NOT NULL,
    parameter_value text NOT NULL,
    description     text,
    created_date    timestamptz,
    updated_date    timestamptz,
    created_by      varchar(255),
    updated_by      varchar(255)
);

CREATE INDEX idx_policy_param_set ON policy_set_parameter (policy_set_id);

CREATE TABLE policy_exemption (
    id               varchar(36) PRIMARY KEY,
    policy_set_id    varchar(36) NOT NULL REFERENCES policy_set (id) ON DELETE CASCADE,
    rule_id          varchar(255) NOT NULL,
    workspace_id     varchar(36) REFERENCES workspace (id),
    project_id       varchar(36) REFERENCES project (id),
    ticket_reference varchar(255) NOT NULL,
    justification    text NOT NULL,
    expires_at       timestamptz,
    organization_id  varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
    created_date     timestamptz,
    updated_date     timestamptz,
    created_by       varchar(255),
    updated_by       varchar(255)
);

CREATE INDEX idx_exemption_org ON policy_exemption (organization_id);
CREATE INDEX idx_exemption_ws_rule ON policy_exemption (workspace_id, rule_id);
