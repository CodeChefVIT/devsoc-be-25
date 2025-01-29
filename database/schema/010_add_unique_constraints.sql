-- +goose Up
ALTER TABLE users ADD UNIQUE (github_profile);

ALTER TABLE teams ADD UNIQUE (name);

-- +goose Down
