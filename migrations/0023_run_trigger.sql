-- Workspace run triggers: a successful apply on the source workspace starts a
-- run on the destination. Columns follow
-- api/src/main/resources/db/changelog/local/changelog-2.34.0-workspace-run-trigger.xml.
-- Both workspace foreign keys are checked: a deleted workspace is soft-deleted,
-- so the row stays and an edge must be removed explicitly.

CREATE TABLE workspace_run_trigger (
    id                     varchar(36) PRIMARY KEY,
    source_workspace_id    varchar(36) NOT NULL REFERENCES workspace (id) ON DELETE CASCADE,
    destination_workspace_id varchar(36) NOT NULL REFERENCES workspace (id),
    template_id            varchar(36) REFERENCES template (id),
    organization_id        varchar(36) NOT NULL REFERENCES organization (id),
    enabled                boolean     NOT NULL DEFAULT true,
    created_date           timestamptz,
    created_by             varchar(255),
    updated_date           timestamptz,
    updated_by             varchar(255),
    CONSTRAINT uk_run_trigger_source_dest UNIQUE (source_workspace_id, destination_workspace_id),
    CONSTRAINT chk_run_trigger_no_self CHECK (source_workspace_id <> destination_workspace_id)
);

-- Dispatch looks up enabled edges by source after every apply.
CREATE INDEX idx_run_trigger_source_enabled ON workspace_run_trigger (source_workspace_id, enabled);

-- The UI's inbound list, and the destination workspace lookups.
CREATE INDEX idx_run_trigger_destination ON workspace_run_trigger (destination_workspace_id);

-- Cycle validation loads one organization's whole graph.
CREATE INDEX idx_run_trigger_organization ON workspace_run_trigger (organization_id);
