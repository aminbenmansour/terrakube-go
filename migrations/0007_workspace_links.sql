-- Links from a workspace to a VCS provider, an SSH key and a self-hosted agent.
-- The referenced tables land with #48 and #120; the columns exist now so the
-- UI's create form and relationship updates round-trip the ids instead of
-- silently dropping them. No foreign keys until those tables exist.
ALTER TABLE workspace
    ADD COLUMN vcs_id   varchar(36),
    ADD COLUMN ssh_id   varchar(36),
    ADD COLUMN agent_id varchar(36);
