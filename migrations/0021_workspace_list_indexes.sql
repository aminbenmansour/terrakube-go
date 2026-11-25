-- The management API lists workspaces ordered by lower(name) with id as the
-- tiebreaker, so paging is stable. Both list routes exclude soft-deleted rows,
-- so the indexes are partial.
CREATE INDEX workspace_organization_name_lower_idx ON workspace (organization_id, lower(name), id) WHERE deleted = false;
CREATE INDEX workspace_project_name_lower_idx ON workspace (project_id, lower(name), id) WHERE deleted = false;
