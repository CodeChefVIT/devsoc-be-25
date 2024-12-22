-- +goose Up
CREATE TABLE IF NOT EXISTS submission (
    id UUID NOT NULL UNIQUE,
    github_link TEXT NOT NULL DEFAULT '',
    figma_link TEXT NOT NULL DEFAULT '',
    ppt_link TEXT NOT NULL DEFAULT '',
    other_link TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS submission;
