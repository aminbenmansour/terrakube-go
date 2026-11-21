-- Organization-level variables, inherited by every workspace in the
-- organization unless the workspace defines the same key. Column names and the
-- organization+key uniqueness mirror upstream's globalvar table; upstream keeps
-- no timestamps on this one, so neither do we.
CREATE TABLE globalvar (
    id                   varchar(36) PRIMARY KEY,
    variable_key         varchar(1024) NOT NULL,
    variable_value       text         NOT NULL,
    variable_description varchar(128) NOT NULL,
    sensitive            boolean      NOT NULL DEFAULT false,
    hcl                  boolean      NOT NULL DEFAULT false,
    variable_category    varchar(32)  NOT NULL,
    organization_id      varchar(36)  NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
    CONSTRAINT const_global_var UNIQUE (organization_id, variable_key)
);

CREATE INDEX globalvar_organization_idx ON globalvar (organization_id);
