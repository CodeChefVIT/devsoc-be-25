-- name: ExportAllUsers :many
SELECT
    u.id,
    u.team_id,
    u.first_name,
    u.last_name,
    u.email,
    u.phone_no,
    u.gender,
    u.reg_no,
    u.hostel_block,
    u.room_no,
    u.github_profile,
    u.password,
    u.role,
    u.is_leader,
    u.is_verified,
    u.is_banned,
    u.is_profile_complete,
    t.round_qualified
FROM
    users u
    JOIN teams t ON u.team_id = t.id;

-- name: ExportAllTeams :many
SELECT
    id,
    name,
    number_of_people,
    round_qualified,
    code
FROM
    teams;
