-- name: GetTeamScores :many
SELECT * FROM score WHERE team_id = $1;

-- name: CreateScore :exec
INSERT INTO score (id, team_id, design, implementation, presentation, round)
VALUES ($1, $2, $3, $4, $5, $6);