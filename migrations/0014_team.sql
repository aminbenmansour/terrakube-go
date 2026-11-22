-- Teams as upstream's team tab models them: a name plus the booleans the RBAC
-- checks read and the role the UI presets. Nullable booleans upstream maps to
-- false are NOT NULL here, so a row always has a usable value.
CREATE TABLE team (
    id                varchar(36) PRIMARY KEY,
    name              varchar(64) NOT NULL,
    manage_workspace  boolean NOT NULL DEFAULT false,
    manage_module     boolean NOT NULL DEFAULT false,
    manage_provider   boolean NOT NULL DEFAULT false,
    manage_vcs        boolean NOT NULL DEFAULT false,
    manage_template   boolean NOT NULL DEFAULT false,
    manage_state      boolean NOT NULL DEFAULT false,
    manage_collection boolean NOT NULL DEFAULT false,
    manage_job        boolean NOT NULL DEFAULT false,
    manage_policies   boolean NOT NULL DEFAULT false,
    plan_job          boolean NOT NULL DEFAULT false,
    approve_job       boolean NOT NULL DEFAULT false,
    destroy_job       boolean NOT NULL DEFAULT true,
    role              varchar(32),
    organization_id   varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE
);

CREATE INDEX team_organization_idx ON team (organization_id);
