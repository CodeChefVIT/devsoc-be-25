-- +goose Up
CREATE TABLE authors (
    id BIGSERIAL PRIMARY KEY,
    name text NOT NULL,
    bio text
);

-- +goose Down
DROP TABLE authors;
