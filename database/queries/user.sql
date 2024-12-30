-- name: GetAllUsers :many
SELECT * FROM user;

-- name: GetAllVitians :many
SELECT * FROM user WHERE is_vitian = TRUE;

-- name: GetUserByEmail :one
SELECT * FROM user WHERE Email = $1;

-- name: BanUser :exec
UPDATE "user"
SET is_banned = TRUE
WHERE email = $1;

-- name: UnbanUser :exec
UPDATE "user"
SET is_banned = FALSE
WHERE email = $1;
