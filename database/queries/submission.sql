-- name: GetSubmissionByTeamID :one
SELECT * FROM submission WHERE team_id = $1;