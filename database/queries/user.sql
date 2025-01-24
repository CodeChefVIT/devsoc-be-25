-- name: GetAllUsers :many
SELECT *
FROM users
WHERE (first_name LIKE '%' || $1 || '%'
       OR reg_no LIKE '%' || $1 || '%'
       OR email LIKE '%' || $1 || '%')
  AND id > $2
ORDER BY id
LIMIT $3;

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
    github_profile,
    password,
    role,
    is_leader,
    is_verified,
    is_banned,
    is_profile_complete
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
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
    github_profile = $8
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
    is_profile_complete = TRUE
WHERE email = $1;

-- name: UpdateStarred :exec
UPDATE users
SET
    is_starred = $1
WHERE email = $2;

-- name: GetUserAndTeamDetails :many
SELECT teams.name, teams.number_of_people, teams.round_qualified, teams.code,
	users.id, users.first_name, users.last_name, users.email, users.reg_no, users.phone_no, users.gender, users.github_profile, users.is_leader
	FROM teams
	INNER JOIN users ON users.team_id = teams.id
	LEFT JOIN submission ON submission.team_id = teams.id
	LEFT JOIN ideas ON ideas.team_id = teams.id
WHERE teams.id = $1;
