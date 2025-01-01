-- +goose Up
CREATE TABLE score (
    id UUID NOT NULL UNIQUE,
    team_id UUID NOT NULL,
    design INTEGER NOT NULL  DEFAULT 0,
    implementation INTEGER NOT NULL DEFAULT 0,
    presentation INTEGER NOT NULL DEFAULT 0,
    round INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE score;