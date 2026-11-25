-- The management API lists projects and tags ordered by lower(name) with id as
-- the tiebreaker, so paging is stable. Tags may be nameless; their index and
-- order coalesce a NULL name to the empty string, which is where the in-memory
-- sort put it.
CREATE INDEX project_organization_name_lower_idx ON project (organization_id, lower(name), id);
CREATE INDEX tag_organization_name_lower_idx ON tag (organization_id, coalesce(lower(name), ''), id);
