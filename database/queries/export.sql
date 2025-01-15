-- name: ExportAllUsers :many
SELECT id, team_id, first_name, last_name, email, phone_no, gender, reg_no, vit_email, hostel_block, room_no, github_profile, password, role, is_leader, is_verified, is_banned, is_profile_complete
FROM users;

-- name: ExportAllTeams :many
SELECT id, name, number_of_people, round_qualified, code
FROM teams;

