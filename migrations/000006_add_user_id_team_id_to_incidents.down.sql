DROP INDEX IF EXISTS idx_incidents_team_id;
DROP INDEX IF EXISTS idx_incidents_user_id;

ALTER TABLE incidents
DROP COLUMN IF EXISTS team_id,
DROP COLUMN IF EXISTS user_id;