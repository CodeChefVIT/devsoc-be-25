-- +goose Up
CREATE TABLE IF NOT EXISTS teams (
    id UUID NOT NULL UNIQUE,
    name TEXT NOT NULL,
    number_of_people INTEGER,
    users UUID NOT NULL UNIQUE,
    submission UUID UNIQUE NOT NULL,
    round_qualified INTEGER DEFAULT 0,
    code TEXT UNIQUE NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS teams;
