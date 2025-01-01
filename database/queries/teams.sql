-- name: GetTeamIDByCode :one
SELECT id FROM teams WHERE code = $1;

-- name: GetTeams :many
SELECT * FROM teams;

-- name: GetTeamById :one
SELECT * FROM teams WHERE id = $1;