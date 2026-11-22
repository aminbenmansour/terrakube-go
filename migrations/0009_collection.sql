-- Variable collections: named groups of variables attached to workspaces
-- through a reference. Column names mirror upstream's 2.24.0 changelog. The
-- item table is called `item` upstream, so it is here too.
CREATE TABLE collection (
    id              varchar(36)  PRIMARY KEY,
    name            varchar(128) NOT NULL,
    description     varchar(512),
    priority        integer      NOT NULL,
    created_date    timestamptz,
    updated_date    timestamptz,
    created_by      varchar(128),
    updated_by      varchar(128),
    organization_id varchar(36)  NOT NULL REFERENCES organization (id) ON DELETE CASCADE
);

CREATE INDEX collection_organization_idx ON collection (organization_id);

CREATE TABLE item (
    id               varchar(36)   PRIMARY KEY,
    item_description varchar(512),
    item_key         varchar(3072) NOT NULL,
    item_value       varchar(3072) NOT NULL,
    sensitive        boolean       NOT NULL DEFAULT false,
    hcl              boolean       NOT NULL DEFAULT false,
    item_category    varchar(32)   NOT NULL,
    created_date     timestamptz,
    updated_date     timestamptz,
    created_by       varchar(128),
    updated_by       varchar(128),
    collection_id    varchar(36)   NOT NULL REFERENCES collection (id) ON DELETE CASCADE
);

CREATE INDEX item_collection_idx ON item (collection_id);

CREATE TABLE reference (
    id            varchar(36) PRIMARY KEY,
    created_date  timestamptz,
    updated_date  timestamptz,
    created_by    varchar(128),
    updated_by    varchar(128),
    description   varchar(512),
    workspace_id  varchar(36) NOT NULL REFERENCES workspace (id) ON DELETE CASCADE,
    collection_id varchar(36) NOT NULL REFERENCES collection (id) ON DELETE CASCADE
);

CREATE INDEX reference_collection_idx ON reference (collection_id);
CREATE INDEX reference_workspace_idx ON reference (workspace_id);
