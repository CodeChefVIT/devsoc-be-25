-- +goose Up
CREATE TABLE submission (
    id UUID NOT NULL UNIQUE,
    title TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    track TEXT NOT NULL DEFAULT '',
    github_link TEXT NOT NULL DEFAULT '',
    figma_link TEXT NOT NULL DEFAULT '',
    other_link TEXT NOT NULL DEFAULT '',
    team_id UUID NOT NULL UNIQUE,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE submission;
