-- name: GetTeamIDByCode :one
SELECT id FROM teams WHERE code = $1;