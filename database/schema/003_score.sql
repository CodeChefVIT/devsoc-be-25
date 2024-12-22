-- +goose Up
CREATE TABLE IF NOT EXISTS score (
    id UUID NOT NULL UNIQUE,
    team_id UUID NOT NULL UNIQUE,
    design INTEGER DEFAULT 0,
    implementation INTEGER DEFAULT 0,
    presentation INTEGER DEFAULT 0,
    round INTEGER DEFAULT 0,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS score;
