-- name: GetIdea :one
SELECT * FROM ideas
WHERE id = $1 LIMIT 1;

-- name: ListIdeas :many
SELECT * FROM ideas
ORDER BY created_at DESC;

-- name: CreateIdea :one
INSERT INTO ideas (
  id, title, description, track, team_id, is_selected
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateIdea :exec
UPDATE ideas
SET title = $2,
    description = $3,
    track = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE team_id = $1;


-- name: DeleteIdea :exec
DELETE FROM ideas
WHERE id = $1;

-- +goose Up
CREATE TABLE ideas (
  id UUID NOT NULL UNIQUE,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  track TEXT NOT NULL,
  team_id UUID NOT NULL,
  is_selected BOOLEAN NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(id)
);

-- name: GetIdeaByTeamID :one
SELECT id, title, description, track
FROM ideas
WHERE team_id = $1
LIMIT 1;


-- +goose Down
DROP TABLE ideas;