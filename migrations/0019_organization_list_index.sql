-- The management API lists organizations ordered by lower(name) with id as the
-- tiebreaker, so paging is stable. This index serves both sort directions and
-- the name equality filter; prefix filters are a residual predicate over the
-- same ordered scan.
CREATE INDEX organization_name_lower_idx ON organization (lower(name), id);
