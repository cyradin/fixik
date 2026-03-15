ALTER TABLE incidents
    ADD COLUMN status TEXT NOT NULL,
    DROP COLUMN status_id;

DROP INDEX IF EXISTS idx_incidents_status_id;
DROP TABLE IF EXISTS statuses;