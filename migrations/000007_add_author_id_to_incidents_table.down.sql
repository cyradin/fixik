DROP INDEX IF EXISTS idx_incidents_author_id;

ALTER TABLE incidents
DROP COLUMN IF EXISTS author_id;