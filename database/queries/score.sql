-- name: GetTeamScores :many
SELECT id, team_id, design, implementation, presentation, round, innovation, teamwork, comment
FROM score
WHERE team_id = $1;

-- name: CreateScore :exec
INSERT INTO score (id, team_id, design, implementation, presentation, round,  innovation, teamwork, comment)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: UpdateScore :exec
UPDATE score
SET team_id = $1, design = $2, implementation = $3, presentation = $4, round = $5, innovation = $6, teamwork = $7, comment = $8
WHERE id = $9;

-- name: DeleteScore :exec
DELETE FROM score
WHERE id = $1;