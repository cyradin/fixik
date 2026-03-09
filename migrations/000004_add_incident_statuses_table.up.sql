CREATE TABLE incident_statuses (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

ALTER TABLE incidents
    ADD COLUMN status_id BIGINT REFERENCES incident_statuses(id) NOT NULL,
    DROP COLUMN status;

CREATE INDEX idx_incidents_status_id
ON incidents(status_id);