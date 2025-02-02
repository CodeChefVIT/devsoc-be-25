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

-- name: GetIdeaByTeamID :one
SELECT id, title, description, track
FROM ideas
WHERE team_id = $1
LIMIT 1;

-- name: GetAllIdeas :many
SELECT * FROM ideas
WHERE id > $1
ORDER BY id
LIMIT $2;

-- name: GetIdeasByTrack :many
SELECT * FROM ideas
WHERE track = $1
AND id > $2
ORDER BY id
LIMIT $3;
