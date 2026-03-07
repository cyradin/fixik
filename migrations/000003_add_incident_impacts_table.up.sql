CREATE TABLE incident_impacts (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE
);

ALTER TABLE incidents
    ADD COLUMN impact_id BIGINT REFERENCES incident_impacts(id) NOT NULL,
    DROP COLUMN impact;

CREATE INDEX idx_incidents_impact_id
ON incidents(impact_id);