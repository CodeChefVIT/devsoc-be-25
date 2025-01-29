-- +goose Up
ALTER TABLE users ADD CONSTRAINT unique_github_profile UNIQUE (github_profile);

ALTER TABLE users ALTER COLUMN github_profile DROP NOT NULL;

ALTER TABLE teams ADD CONSTRAINT unique_team_name UNIQUE (name);

-- +goose Down
ALTER TABLE users DROP CONSTRAINT unique_github_profile;

ALTER TABLE teams DROP CONSTRAINT unique_team_name;
