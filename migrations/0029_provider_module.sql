-- Provider, provider version, implementation, module and module version,
-- mirroring upstream's provider/module tables (changelog-1.0 and 1.2 with the
-- later 2.19, 2.30 and 2.33 alterations). These are the private-registry
-- management rows the Settings and Registry screens read and write; the files
-- themselves live in storage and are served by cmd/registry.
--
-- github_app_token_id is omitted: this implementation has no GitHub App
-- credentials table yet, so there is nothing for the column to reference.

CREATE TABLE provider (
    id                 varchar(36) PRIMARY KEY,
    name               varchar(32) NOT NULL,
    description        varchar(256),
    organization_id    varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
    imported           boolean DEFAULT false,
    registry_namespace varchar(64),
    CONSTRAINT const_provider UNIQUE (organization_id, name)
        DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE version (
    id                  varchar(36) PRIMARY KEY,
    version_number      varchar(64) NOT NULL,
    protocols           varchar(64) NOT NULL,
    provider_id         varchar(36) NOT NULL REFERENCES provider (id),
    status              varchar(16) NOT NULL DEFAULT 'active',
    deprecation_message varchar(1024),
    CONSTRAINT const_version UNIQUE (provider_id, version_number)
        DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE implementation (
    id                    varchar(36) PRIMARY KEY,
    os                    varchar(32) NOT NULL,
    arch                  varchar(32),
    filename              varchar(512) NOT NULL,
    download_url          varchar(1024) NOT NULL,
    shasums_url           varchar(1024) NOT NULL,
    shasums_signature_url varchar(1024) NOT NULL,
    shasum                varchar(1024) NOT NULL,
    key_id                varchar(32) NOT NULL,
    ascii_armor           text,
    trust_signature       text NOT NULL,
    source                varchar(64) NOT NULL,
    source_url            varchar(512) NOT NULL,
    version_id            varchar(36) NOT NULL REFERENCES version (id),
    CONSTRAINT const_implementation UNIQUE (version_id, os, arch)
        DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE module (
    id                varchar(36) PRIMARY KEY,
    name              varchar(256) NOT NULL,
    description       varchar(256),
    provider          varchar(64) NOT NULL,
    source            varchar(1024) NOT NULL,
    organization_id   varchar(36) NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
    vcs_id            varchar(36) REFERENCES vcs (id),
    ssh_id            varchar(36) REFERENCES ssh (id),
    download_quantity bigint DEFAULT 0,
    created_date      timestamptz NOT NULL DEFAULT now(),
    updated_date      timestamptz NOT NULL DEFAULT now(),
    created_by        varchar(128),
    updated_by        varchar(128),
    tag_prefix        varchar(256),
    folder            varchar(256),
    latest_version    varchar(64),
    CONSTRAINT const_module UNIQUE (organization_id, name, provider)
        DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE module_version (
    id                  varchar(36) PRIMARY KEY,
    module_id           varchar(36) NOT NULL REFERENCES module (id) ON DELETE CASCADE,
    version             varchar(256) NOT NULL,
    commit_info         varchar(256) NOT NULL,
    git_tag             varchar(256),
    status              varchar(16) NOT NULL DEFAULT 'active',
    deprecation_message varchar(1024),
    CONSTRAINT uk_module_version_module_id_version UNIQUE (module_id, version)
);

CREATE INDEX provider_organization_idx ON provider (organization_id);
CREATE INDEX version_provider_idx ON version (provider_id);
CREATE INDEX implementation_version_idx ON implementation (version_id);
CREATE INDEX module_organization_idx ON module (organization_id);
CREATE INDEX idx_module_version_module_id ON module_version (module_id);
