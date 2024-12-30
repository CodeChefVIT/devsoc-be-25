-- +goose Up
ALTER TABLE users ADD CONSTRAINT fk_users_teams FOREIGN KEY(team_id) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE score ADD CONSTRAINT fk_score_teams FOREIGN KEY(team_id) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE submission ADD CONSTRAINT fk_submission_teams FOREIGN KEY(team_id) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ideas ADD CONSTRAINT fk_ideas_teams FOREIGN KEY(team_id) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT fk_users_teams;

ALTER TABLE score
DROP CONSTRAINT fk_score_teams;

ALTER TABLE submission
DROP CONSTRAINT fk_submission_teams;

ALTER TABLE ideas
DROP CONSTRAINT fk_ideas_teams;
