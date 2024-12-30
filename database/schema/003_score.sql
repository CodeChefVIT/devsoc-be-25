-- +goose Up
CREATE TABLE score (
    id UUID NOT NULL UNIQUE,
    team_id UUID NOT NULL,
    design INTEGER DEFAULT 0,
    implementation INTEGER DEFAULT 0,
    presentation INTEGER DEFAULT 0,
    round INTEGER DEFAULT 0,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE score;
