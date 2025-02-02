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

-- name: GetLeaderboardWithPagination :many
WITH RoundScores AS (
    SELECT 
        s.team_id,
        t.name,
        s.round,
        s.design,
        s.implementation,
        s.presentation,
        s.innovation,
        s.teamwork,
        (s.design + s.implementation + s.presentation + s.innovation + s.teamwork) AS round_total
    FROM score s
    JOIN teams t ON s.team_id = t.id
),
TotalScores AS (
    SELECT 
        s.team_id,
        t.name,
        SUM(s.design + s.implementation + s.presentation + s.innovation + s.teamwork) AS overall_total
    FROM score s
    JOIN teams t ON s.team_id = t.id
    GROUP BY s.team_id, t.name
)
SELECT 
    rs.team_id,
    rs.name,
    rs.round,
    rs.design,
    rs.implementation,
    rs.presentation,
    rs.innovation,
    rs.teamwork,
    rs.round_total,
    ts.overall_total
FROM RoundScores rs
JOIN TotalScores ts ON rs.team_id = ts.team_id
WHERE 
    ($1::UUID IS NULL OR rs.team_id > $1)  -- Cursor-based pagination
    AND ($2::TEXT IS NULL OR rs.name ILIKE '%' || $2 || '%') -- Optional name filter
ORDER BY ts.overall_total DESC, rs.round ASC
LIMIT $3;