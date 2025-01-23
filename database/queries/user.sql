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
    vit_email,
    hostel_block,
    room_no,
    github_profile,
    password,
    role,
    is_leader,
    is_verified,
    is_banned,
    is_profile_complete
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
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

-- name: CompleteProfile :exec
UPDATE users
SET
    first_name = $2,
    last_name = $3,
    phone_no = $4,
    gender = $5,
    reg_no = $6,
    vit_email = $7,
    hostel_block = $8,
    room_no = $9,
    github_profile = $10,
    is_profile_complete = TRUE
WHERE email = $1;

-- Goofy ahh query hai, but kaam karega if decoded
-- SELECT
--   (json_build_object(
--     'user', json_strip_nulls(json_build_object(
--       'first_name', u.first_name,
--       'last_name', u.last_name,
--       'email', u.email,
--       'phone_no', u.phone_no,
--       'gender', u.gender,
--       'reg_no', u.reg_no,
--       'vit_email', u.vit_email,
--       'hostel_block', u.hostel_block,
--       'room_no', u.room_no,
--       'github_profile', u.github_profile,
--       'role', u.role
--     )),
--     'team', json_build_object(
--       'team_name', t.name,
--       'number_of_people', t.number_of_people,
--       'round_qualified', t.round_qualified,
--       'code', t.code,
--       'members', (
--         SELECT json_agg(json_strip_nulls(json_build_object(
--           'first_name', members.first_name,
--           'last_name', members.last_name,
--           'email', members.email,
--           'phone_no', members.phone_no,
--           'github_profile', members.github_profile,
--           'role', members.role,
--           'is_leader', members.is_leader
--         )))
--         FROM users members
--         WHERE members.team_id = t.id AND members.id != u.id
--       )
--     )
--   ))::json AS result
-- FROM
--   users u
-- JOIN
--   teams t ON u.team_id = t.id
-- WHERE
--   u.id = $1;

-- name: GetUserAndTeamDetails :many
SELECT teams.name, teams.number_of_people, teams.round_qualified, teams.code,
	users.id, users.first_name, users.last_name, users.email, users.reg_no, users.phone_no, users.gender, users.vit_email, users.hostel_block, users.room_no, users.github_profile, users.is_leader
	FROM teams
	INNER JOIN users ON users.team_id = teams.id
	LEFT JOIN submission ON submission.team_id = teams.id
	LEFT JOIN ideas ON ideas.team_id = teams.id
WHERE teams.id = $1;
