-- +goose Up
CREATE TABLE submission (
    id UUID NOT NULL UNIQUE,
    github_link TEXT NOT NULL DEFAULT '',
    figma_link TEXT NOT NULL DEFAULT '',
    ppt_link TEXT NOT NULL DEFAULT '',
    other_link TEXT NOT NULL DEFAULT '',
    team_id UUID NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE submission;
