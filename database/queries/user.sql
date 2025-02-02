-- name: GetAllUsers :many
SELECT u.*,t.round_qualified
FROM users u
JOIN teams t ON t.id = u.team_id
WHERE (u.first_name LIKE '%' || $1 || '%'
       OR u.reg_no LIKE '%' || $1 || '%'
       OR u.email LIKE '%' || $1 || '%')
  AND u.id > $2
ORDER BY u.id
LIMIT $3;

-- name: GetUsersByGender :many
SELECT * FROM users WHERE gender = $1;

-- name: GetUsersByTeamId :many
SELECT first_name, last_name, email, reg_no, phone_no FROM users WHERE team_id = $1;

-- name: GetUsers :many
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
    password,
    role,
    is_leader,
    is_verified,
    is_banned,
    is_profile_complete
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
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
    phone_no = $4,
    gender = $5,
    reg_no = $6,
    github_profile = $7,
    hostel_block = $8,
    room_no = $9
WHERE id = $1;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByPhoneNo :one
SELECT * FROM users WHERE phone_no = $1;

-- name: CompleteProfile :exec
UPDATE users
SET
    first_name = $2,
    last_name = $3,
    phone_no = $4,
    gender = $5,
    reg_no = $6,
    github_profile = $7,
    hostel_block = $8,
    room_no = $9,
    is_profile_complete = TRUE
WHERE email = $1;

-- name: UpdateStarred :exec
UPDATE users
SET
    is_starred = $1
WHERE email = $2;

-- name: UpdateGitHub :exec
UPDATE users
SET
    github_profile = $1
WHERE email = $2;
