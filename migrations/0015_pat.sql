-- Personal access tokens as upstream's Pat entity stores them. The token itself
-- is a signed JWT; only this metadata row is persisted. Revoking is a soft
-- delete, checked by the bearer validator on every request.
CREATE TABLE pat (
    id           varchar(36) PRIMARY KEY,
    days         integer,
    description  varchar(64),
    deleted      boolean NOT NULL DEFAULT false,
    source       varchar(20) NOT NULL DEFAULT 'API',
    last_used_at timestamptz,
    created_date timestamptz,
    updated_date timestamptz,
    created_by   varchar(128),
    updated_by   varchar(128)
);
