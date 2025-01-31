-- name: GetTeamIDByCode :one
SELECT id FROM teams WHERE code = $1;

-- name: GetTeams :many
SELECT *
FROM teams
WHERE name ILIKE '%' || $1 || '%'
  AND id > $2
ORDER BY id
LIMIT $3;


-- name: GetTeamById :one
SELECT teams.id, teams.name, teams.round_qualified, teams.code,teams.is_banned,
       score.design, score.implementation, score.presentation, score.round,
       submission.title, submission.description, submission.track, submission.github_link, submission.figma_link, submission.other_link,
       ideas.title, ideas.description, ideas.track, ideas.is_selected
FROM teams
INNER JOIN score ON score.team_id = teams.id
LEFT JOIN submission ON submission.team_id = teams.id
LEFT JOIN ideas ON ideas.team_id = teams.id
WHERE teams.id = $1;

-- name: GetTeamByTeamId :one
SELECT * FROM teams WHERE id = $1;

-- name: FindTeam :one
SELECT id,name,code,round_qualified FROM teams
WHERE code = $1
LIMIT 1;

-- name: KickMemeber :exec
UPDATE users
SET team_id = NULL
WHERE id = $1;

-- name: LeaveTeam :exec
UPDATE users
SET team_id = NULL
WHERE id = $1;

-- name: CreateTeam :one
INSERT INTO teams (
    id, name, number_of_people, round_qualified, code, is_banned
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1;

-- name: CountTeamMembers :one
SELECT COUNT(*) FROM users
WHERE team_id = $1;

-- name: AddUserToTeam :exec
UPDATE users
SET team_id = $1
WHERE id = $2;

-- name: RemoveUserFromTeam :exec
UPDATE users
SET team_id = NULL
WHERE team_id = $1 AND id = $2;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserTeam :exec
UPDATE users
SET team_id = $1, is_leader = $2
WHERE id = $3;

-- name: IncreaseCountTeam :exec
UPDATE teams
SET number_of_people = number_of_people + 1
WHERE id = $1;

-- name: DecreaseUserCountTeam :exec
UPDATE teams
SET number_of_people = number_of_people - 1
WHERE id = $1;

-- name: RemoveTeamIDFromUsers :exec
UPDATE users
SET team_id = NULL
WHERE team_id = $1;

-- name: UpdateLeader :exec
UPDATE users
SET is_leader = $1
WHERE id = $2;

-- name: UpdateTeamName :exec
UPDATE teams
SET name = $1
WHERE id = $2;

-- name: GetTeamMembers :many
SELECT first_name , last_name, is_leader, github_profile, reg_no, phone_no
FROM users
Where team_id = $1;

-- name: GetTeamUsers :many
SELECT first_name, last_name
From users
Where team_id = $1;

-- name: GetTeamUsersEmails :many
SELECT email
FROM users
WHERE team_id = $1;

-- name: BanTeam :exec
UPDATE users
SET is_banned = TRUE
WHERE id = $1;

-- name: UnBanTeam :exec
UPDATE teams
SET is_banned = FALSE
WHERE id = $1;

-- name: UpdateTeamRound :exec
UPDATE teams
SET round_qualified = $1
WHERE id = $2;

-- name: InfoQuery :many
SELECT * FROM teams INNER JOIN users ON users.team_id = teams.id WHERE teams.id = $1;
