-- name: GetTeamScores :many
SELECT * FROM score WHERE team_id = $1;

-- name: CreateScore :exec
INSERT INTO score (id, team_id, design, implementation, presentation, round)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateScore :exec
UPDATE score
SET team_id = $1, design = $2, implementation = $3, presentation = $4, round = $5
WHERE id = $6;

-- name: DeleteScore :exec
DELETE FROM score
WHERE id = $1;