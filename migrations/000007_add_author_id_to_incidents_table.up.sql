ALTER TABLE incidents
ADD COLUMN author_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_incidents_author_id ON incidents(author_id);
