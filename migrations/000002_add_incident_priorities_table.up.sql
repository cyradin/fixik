CREATE TABLE priorities (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

ALTER TABLE incidents
    ADD COLUMN priority_id BIGINT REFERENCES priorities(id) NOT NULL,
    DROP COLUMN priority;

CREATE INDEX idx_incidents_priority_id
ON incidents(priority_id);