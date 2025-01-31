-- +goose Up
ALTER TABLE users
ADD COLUMN is_starred BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
