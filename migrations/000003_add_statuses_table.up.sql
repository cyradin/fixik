CREATE TABLE statuses (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    sort INT NOT NULL DEFAULT 0,
    is_final BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

ALTER TABLE incidents
    ADD COLUMN status_id BIGINT REFERENCES statuses(id) NOT NULL,
    DROP COLUMN status;

CREATE INDEX idx_incidents_status_id
ON incidents(status_id);