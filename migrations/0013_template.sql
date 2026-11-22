-- Templates existed from 0001 for the notification join; the management API
-- needs upstream's nullable name and audit columns. Upstream's 1.5 table lets a
-- create with no attributes through, so name must lose its NOT NULL.
ALTER TABLE template ALTER COLUMN name DROP NOT NULL;
ALTER TABLE template ALTER COLUMN name TYPE varchar(256);
ALTER TABLE template ADD COLUMN created_date timestamptz;
ALTER TABLE template ADD COLUMN updated_date timestamptz;
ALTER TABLE template ADD COLUMN created_by varchar(128);
ALTER TABLE template ADD COLUMN updated_by varchar(128);
