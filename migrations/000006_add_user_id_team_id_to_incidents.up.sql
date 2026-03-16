ALTER TABLE incidents
ADD COLUMN team_id BIGINT NULL REFERENCES teams(id) ON DELETE SET NULL,
ADD COLUMN user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_incidents_team_id ON incidents(team_id);
CREATE INDEX idx_incidents_user_id ON incidents(user_id);