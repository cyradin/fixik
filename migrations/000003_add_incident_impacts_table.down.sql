ALTER TABLE incidents
    ADD COLUMN impact TEXT NOT NULL,
    DROP COLUMN impact_id;

DROP INDEX IF EXISTS idx_incidents_impact_id;
DROP TABLE IF EXISTS incident_impacts;