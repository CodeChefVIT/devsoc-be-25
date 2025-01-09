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
    id,
    team_id,
    first_name,
    last_name, 
    email,
    phone_no,
    gender,
    reg_no,
    vit_email,
    hostel_block,
    room_no,
    github_profile,
    password,
    role,
    is_leader,
    is_verified,
    is_banned
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
);

-- name: GetUserByRegNo :one
SELECT * FROM users WHERE reg_no = $1;

-- name: VerifyUser :exec
UPDATE users
SET is_verified = TRUE
WHERE email = $1;

-- name: UpdatePassword :exec
UPDATE users
SET password = $2
WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users
SET first_name = $2,
    last_name = $3,
    email = $4,
    phone_no = $5,
    gender = $6,
    reg_no = $7,
    vit_email = $8,
    hostel_block = $9,
    room_no = $10,
    github_profile = $11
WHERE id = $1;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByPhoneNo :one
SELECT * FROM users WHERE phone_no = $1;

-- name: GetUserByVitEmail :one
SELECT * FROM users WHERE vit_email = $1;