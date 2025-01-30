-- name: GetTeamScores :many
SELECT id, team_id, design, implementation, presentation, round, points1, points2, points3, comment
FROM score
WHERE team_id = $1;

-- name: CreateScore :exec
INSERT INTO score (id, team_id, design, implementation, presentation, round, points1, points2, points3, comment)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: UpdateScore :exec
UPDATE score
SET team_id = $1, design = $2, implementation = $3, presentation = $4, round = $5, points1 = $6, points2 = $7, points3 = $8, comment = $9
WHERE id = $10;

-- name: DeleteScore :exec
DELETE FROM score
WHERE id = $1;