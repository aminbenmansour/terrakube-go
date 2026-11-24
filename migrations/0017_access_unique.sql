-- Upstream's changelog-2.25.0-workspace-access.xml names this constraint
-- const_workspace_access, deferrable and initially deferred. 0005 created the
-- table without it, so a workspace could grant the same group twice and
-- deleting one row would silently leave the other grant in place.
--
-- A database that already holds duplicates will fail here rather than lose a
-- grant silently; remove the duplicate rows first.
ALTER TABLE access
    ADD CONSTRAINT const_workspace_access UNIQUE (workspace_id, name)
    DEFERRABLE INITIALLY DEFERRED;
