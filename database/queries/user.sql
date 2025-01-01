-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetAllVitians :many
SELECT * FROM users WHERE is_vitian = TRUE;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: BanUser :exec
UPDATE users
SET is_banned = TRUE
WHERE email = $1;

-- name: UnbanUser :exec
UPDATE users
SET is_banned = FALSE
WHERE email = $1;

-- name: GetTeamLeader :one
SELECT * FROM users WHERE team_id = $1 AND is_leader = TRUE;

-- name: CreateUser :exec
INSERT INTO users (
    id, name, team_id, email, is_vitian, reg_no, password, phone_no, role, is_leader, college, is_verified, is_banned
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
);
