-- Upstream's tag.name is nullable and a create without a name stores NULL (the
-- UI always sends one, but Elide does not require it).
ALTER TABLE tag ALTER COLUMN name DROP NOT NULL;
