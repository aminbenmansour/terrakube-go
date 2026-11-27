-- Workspace repository webhooks, mirroring upstream's webhook and webhook_event
-- tables after 2.34 (changelog 2.18.0 through 2.34.0). A workspace has one
-- webhook row; its triggers live in webhook_event, checked by ascending
-- priority. remote_hook_id is the provider's id for the registered hook, and
-- migrated_v2 marks the shared repository webhook introduced in 2.31.

CREATE TABLE webhook (
    id            varchar(36) PRIMARY KEY,
    workspace_id  varchar(36) NOT NULL
                  CONSTRAINT fk_webhook_workspace_id REFERENCES workspace (id) ON DELETE CASCADE,
    remote_hook_id varchar(220),
    migrated_v2   boolean NOT NULL DEFAULT false,
    sync_error    varchar(1024),
    created_date  timestamp with time zone,
    updated_date  timestamp with time zone,
    created_by    varchar(128),
    updated_by    varchar(128)
);

CREATE INDEX idx_webhook_workspace_id ON webhook (workspace_id);

CREATE TABLE webhook_event (
    id                  varchar(36) PRIMARY KEY,
    webhook_id          varchar(36) NOT NULL
                        CONSTRAINT fk_webhook_event_webhook_id REFERENCES webhook (id) ON DELETE CASCADE,
    branch              varchar(512) NOT NULL,
    path                varchar(512) NOT NULL,
    path_type           varchar(32) NOT NULL,
    event               varchar(72) NOT NULL,
    template_id         varchar(36) NOT NULL
                        CONSTRAINT fk_webhook_event_template_id REFERENCES template (id),
    priority            integer NOT NULL DEFAULT 0,
    pr_workflow_enabled boolean NOT NULL DEFAULT false,
    pr_apply_enabled    boolean NOT NULL DEFAULT false,
    created_date        timestamp with time zone,
    updated_date        timestamp with time zone,
    created_by          varchar(128),
    updated_by          varchar(128)
);
