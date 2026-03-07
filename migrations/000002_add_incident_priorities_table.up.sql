CREATE TABLE incident_priorities (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE
);

ALTER TABLE incidents
    ADD COLUMN priority_id BIGINT REFERENCES incident_priorities(id) NOT NULL,
    DROP COLUMN priority;

CREATE INDEX idx_incidents_priority_id
ON incidents(priority_id);