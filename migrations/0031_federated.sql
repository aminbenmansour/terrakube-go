-- Federated (OIDC) credentials, mirroring upstream's 2.31.0 tables. A
-- credential names an issuer and audience and owns the claim conditions a token
-- must match; upstream stores no organization because the credential's name is
-- the team the token acts as.
--
-- Claims do not cascade: upstream refuses to delete a credential that still has
-- claims (409), so the UI deletes them first. The foreign key stays plain.

CREATE TABLE federated (
    id           varchar(36) PRIMARY KEY,
    name         varchar(256) NOT NULL,
    issuer_url   varchar(256),
    audience     varchar(256),
    created_date timestamptz,
    updated_date timestamptz,
    created_by   varchar(128),
    updated_by   varchar(128)
);

CREATE TABLE federated_claim (
    id           varchar(36) PRIMARY KEY,
    claim_key    varchar(256) NOT NULL,
    claim_value  varchar(1024) NOT NULL,
    federated_id varchar(36) NOT NULL REFERENCES federated (id)
);

CREATE INDEX idx_federated_claim_federated ON federated_claim (federated_id);
