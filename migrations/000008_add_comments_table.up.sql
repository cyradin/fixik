CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,

    author_id BIGINT NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    incident_id BIGINT NOT NULL
        REFERENCES incidents(id) ON DELETE CASCADE,

    text TEXT NOT NULL
        CHECK (length(text) > 0 AND length(text) <= 5000),

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_comments_author_id
    ON comments(author_id);

CREATE INDEX idx_comments_incident_id
    ON comments(incident_id);

CREATE INDEX idx_comments_incident_created
    ON comments(incident_id, created_at);