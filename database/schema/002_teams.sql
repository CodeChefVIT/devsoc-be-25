-- +goose Up
CREATE TABLE teams (
    id UUID NOT NULL UNIQUE,
    name TEXT NOT NULL,
    number_of_people INTEGER NOT NULL,
    round_qualified INTEGER DEFAULT 0,
    code TEXT UNIQUE NOT NULL,
    is_banned BOOLEAN NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE teams;
