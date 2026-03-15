ALTER TABLE incidents
    ADD COLUMN priority TEXT NOT NULL,
    DROP COLUMN priority_id;

DROP INDEX IF EXISTS idx_incidents_priority_id;
DROP TABLE IF EXISTS priorities;