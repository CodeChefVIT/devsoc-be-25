-- +goose Up
ALTER TABLE score
ADD COLUMN innovation INTEGER NOT NULL DEFAULT 0,
ADD COLUMN teamwork INTEGER NOT NULL DEFAULT 0,
ADD COLUMN comment TEXT;

-- +goose Down
ALTER TABLE score
DROP COLUMN innovation,
DROP COLUMN teamwork,
DROP COLUMN comment;