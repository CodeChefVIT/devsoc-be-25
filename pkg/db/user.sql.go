// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const banUser = `-- name: BanUser :exec
UPDATE users
SET is_banned = TRUE
WHERE email = $1
`

func (q *Queries) BanUser(ctx context.Context, email string) error {
	_, err := q.db.Exec(ctx, banUser, email)
	return err
}

const completeProfile = `-- name: CompleteProfile :exec
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
WHERE email = $1
`

type CompleteProfileParams struct {
	Email         string
	FirstName     string
	LastName      string
	PhoneNo       pgtype.Text
	Gender        string
	RegNo         *string
	GithubProfile *string
	HostelBlock   *string
	RoomNo        *string
}

func (q *Queries) CompleteProfile(ctx context.Context, arg CompleteProfileParams) error {
	_, err := q.db.Exec(ctx, completeProfile,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.PhoneNo,
		arg.Gender,
		arg.RegNo,
		arg.GithubProfile,
		arg.HostelBlock,
		arg.RoomNo,
	)
	return err
}

const createUser = `-- name: CreateUser :exec
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
)
`

type CreateUserParams struct {
	ID                uuid.UUID
	TeamID            uuid.NullUUID
	FirstName         string
	LastName          string
	Email             string
	PhoneNo           pgtype.Text
	Gender            string
	RegNo             *string
	Password          string
	Role              string
	IsLeader          bool
	IsVerified        bool
	IsBanned          bool
	IsProfileComplete bool
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.ID,
		arg.TeamID,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.PhoneNo,
		arg.Gender,
		arg.RegNo,
		arg.Password,
		arg.Role,
		arg.IsLeader,
		arg.IsVerified,
		arg.IsBanned,
		arg.IsProfileComplete,
	)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT u.id, u.team_id, u.first_name, u.last_name, u.email, u.phone_no, u.gender, u.reg_no, u.github_profile, u.password, u.role, u.is_leader, u.is_verified, u.is_banned, u.is_profile_complete, u.is_starred, u.room_no, u.hostel_block, t.round_qualified
FROM users u
JOIN teams t ON t.id = u.team_id
WHERE (u.first_name LIKE '%' || $1 || '%'
       OR u.reg_no LIKE '%' || $1 || '%'
       OR u.email LIKE '%' || $1 || '%')
  AND u.id > $2
  AND ($4 = '' OR u.gender = $4)
ORDER BY u.id
LIMIT $3
`

type GetAllUsersParams struct {
	Column1 *string
	ID      uuid.UUID
	Limit   int32
	Column4 interface{}
}

type GetAllUsersRow struct {
	ID                uuid.UUID
	TeamID            uuid.NullUUID
	FirstName         string
	LastName          string
	Email             string
	PhoneNo           pgtype.Text
	Gender            string
	RegNo             *string
	GithubProfile     *string
	Password          string
	Role              string
	IsLeader          bool
	IsVerified        bool
	IsBanned          bool
	IsProfileComplete bool
	IsStarred         bool
	RoomNo            *string
	HostelBlock       *string
	RoundQualified    pgtype.Int4
}

func (q *Queries) GetAllUsers(ctx context.Context, arg GetAllUsersParams) ([]GetAllUsersRow, error) {
	rows, err := q.db.Query(ctx, getAllUsers,
		arg.Column1,
		arg.ID,
		arg.Limit,
		arg.Column4,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUsersRow
	for rows.Next() {
		var i GetAllUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.PhoneNo,
			&i.Gender,
			&i.RegNo,
			&i.GithubProfile,
			&i.Password,
			&i.Role,
			&i.IsLeader,
			&i.IsVerified,
			&i.IsBanned,
			&i.IsProfileComplete,
			&i.IsStarred,
			&i.RoomNo,
			&i.HostelBlock,
			&i.RoundQualified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllVitians = `-- name: GetAllVitians :many
SELECT id, team_id, first_name, last_name, email, phone_no, gender, reg_no, github_profile, password, role, is_leader, is_verified, is_banned, is_profile_complete, is_starred, room_no, hostel_block FROM users WHERE is_vitian = TRUE
`

func (q *Queries) GetAllVitians(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllVitians)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.PhoneNo,
			&i.Gender,
			&i.RegNo,
			&i.GithubProfile,
			&i.Password,
			&i.Role,
			&i.IsLeader,
			&i.IsVerified,
			&i.IsBanned,
			&i.IsProfileComplete,
			&i.IsStarred,
			&i.RoomNo,
			&i.HostelBlock,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTeamLeader = `-- name: GetTeamLeader :one
SELECT id, team_id, first_name, last_name, email, phone_no, gender, reg_no, github_profile, password, role, is_leader, is_verified, is_banned, is_profile_complete, is_starred, room_no, hostel_block FROM users WHERE team_id = $1 AND is_leader = TRUE
`

func (q *Queries) GetTeamLeader(ctx context.Context, teamID uuid.NullUUID) (User, error) {
	row := q.db.QueryRow(ctx, getTeamLeader, teamID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PhoneNo,
		&i.Gender,
		&i.RegNo,
		&i.GithubProfile,
		&i.Password,
		&i.Role,
		&i.IsLeader,
		&i.IsVerified,
		&i.IsBanned,
		&i.IsProfileComplete,
		&i.IsStarred,
		&i.RoomNo,
		&i.HostelBlock,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, team_id, first_name, last_name, email, phone_no, gender, reg_no, github_profile, password, role, is_leader, is_verified, is_banned, is_profile_complete, is_starred, room_no, hostel_block FROM users WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PhoneNo,
		&i.Gender,
		&i.RegNo,
		&i.GithubProfile,
		&i.Password,
		&i.Role,
		&i.IsLeader,
		&i.IsVerified,
		&i.IsBanned,
		&i.IsProfileComplete,
		&i.IsStarred,
		&i.RoomNo,
		&i.HostelBlock,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, team_id, first_name, last_name, email, phone_no, gender, reg_no, github_profile, password, role, is_leader, is_verified, is_banned, is_profile_complete, is_starred, room_no, hostel_block FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PhoneNo,
		&i.Gender,
		&i.RegNo,
		&i.GithubProfile,
		&i.Password,
		&i.Role,
		&i.IsLeader,
		&i.IsVerified,
		&i.IsBanned,
		&i.IsProfileComplete,
		&i.IsStarred,
		&i.RoomNo,
		&i.HostelBlock,
	)
	return i, err
}

const getUserByPhoneNo = `-- name: GetUserByPhoneNo :one
SELECT id, team_id, first_name, last_name, email, phone_no, gender, reg_no, github_profile, password, role, is_leader, is_verified, is_banned, is_profile_complete, is_starred, room_no, hostel_block FROM users WHERE phone_no = $1
`

func (q *Queries) GetUserByPhoneNo(ctx context.Context, phoneNo pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUserByPhoneNo, phoneNo)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PhoneNo,
		&i.Gender,
		&i.RegNo,
		&i.GithubProfile,
		&i.Password,
		&i.Role,
		&i.IsLeader,
		&i.IsVerified,
		&i.IsBanned,
		&i.IsProfileComplete,
		&i.IsStarred,
		&i.RoomNo,
		&i.HostelBlock,
	)
	return i, err
}

const getUserByRegNo = `-- name: GetUserByRegNo :one
SELECT id, team_id, first_name, last_name, email, phone_no, gender, reg_no, github_profile, password, role, is_leader, is_verified, is_banned, is_profile_complete, is_starred, room_no, hostel_block FROM users WHERE reg_no = $1
`

func (q *Queries) GetUserByRegNo(ctx context.Context, regNo *string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByRegNo, regNo)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PhoneNo,
		&i.Gender,
		&i.RegNo,
		&i.GithubProfile,
		&i.Password,
		&i.Role,
		&i.IsLeader,
		&i.IsVerified,
		&i.IsBanned,
		&i.IsProfileComplete,
		&i.IsStarred,
		&i.RoomNo,
		&i.HostelBlock,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, team_id, first_name, last_name, email, phone_no, gender, reg_no, github_profile, password, role, is_leader, is_verified, is_banned, is_profile_complete, is_starred, room_no, hostel_block FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.PhoneNo,
			&i.Gender,
			&i.RegNo,
			&i.GithubProfile,
			&i.Password,
			&i.Role,
			&i.IsLeader,
			&i.IsVerified,
			&i.IsBanned,
			&i.IsProfileComplete,
			&i.IsStarred,
			&i.RoomNo,
			&i.HostelBlock,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersByGender = `-- name: GetUsersByGender :many
SELECT id, team_id, first_name, last_name, email, phone_no, gender, reg_no, github_profile, password, role, is_leader, is_verified, is_banned, is_profile_complete, is_starred, room_no, hostel_block FROM users WHERE gender = $1
`

func (q *Queries) GetUsersByGender(ctx context.Context, gender string) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsersByGender, gender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.PhoneNo,
			&i.Gender,
			&i.RegNo,
			&i.GithubProfile,
			&i.Password,
			&i.Role,
			&i.IsLeader,
			&i.IsVerified,
			&i.IsBanned,
			&i.IsProfileComplete,
			&i.IsStarred,
			&i.RoomNo,
			&i.HostelBlock,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersByTeamId = `-- name: GetUsersByTeamId :many
SELECT first_name, last_name, email, reg_no, phone_no FROM users WHERE team_id = $1
`

type GetUsersByTeamIdRow struct {
	FirstName string
	LastName  string
	Email     string
	RegNo     *string
	PhoneNo   pgtype.Text
}

func (q *Queries) GetUsersByTeamId(ctx context.Context, teamID uuid.NullUUID) ([]GetUsersByTeamIdRow, error) {
	rows, err := q.db.Query(ctx, getUsersByTeamId, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersByTeamIdRow
	for rows.Next() {
		var i GetUsersByTeamIdRow
		if err := rows.Scan(
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.RegNo,
			&i.PhoneNo,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unbanUser = `-- name: UnbanUser :exec
UPDATE users
SET is_banned = FALSE
WHERE email = $1
`

func (q *Queries) UnbanUser(ctx context.Context, email string) error {
	_, err := q.db.Exec(ctx, unbanUser, email)
	return err
}

const updateGitHub = `-- name: UpdateGitHub :exec
UPDATE users
SET
    github_profile = $1
WHERE email = $2
`

type UpdateGitHubParams struct {
	GithubProfile *string
	Email         string
}

func (q *Queries) UpdateGitHub(ctx context.Context, arg UpdateGitHubParams) error {
	_, err := q.db.Exec(ctx, updateGitHub, arg.GithubProfile, arg.Email)
	return err
}

const updatePassword = `-- name: UpdatePassword :exec
UPDATE users
SET password = $2
WHERE email = $1
`

type UpdatePasswordParams struct {
	Email    string
	Password string
}

func (q *Queries) UpdatePassword(ctx context.Context, arg UpdatePasswordParams) error {
	_, err := q.db.Exec(ctx, updatePassword, arg.Email, arg.Password)
	return err
}

const updateStarred = `-- name: UpdateStarred :exec
UPDATE users
SET
    is_starred = $1
WHERE email = $2
`

type UpdateStarredParams struct {
	IsStarred bool
	Email     string
}

func (q *Queries) UpdateStarred(ctx context.Context, arg UpdateStarredParams) error {
	_, err := q.db.Exec(ctx, updateStarred, arg.IsStarred, arg.Email)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET first_name = $2,
    last_name = $3,
    phone_no = $4,
    gender = $5,
    reg_no = $6,
    github_profile = $7,
    hostel_block = $8,
    room_no = $9
WHERE id = $1
`

type UpdateUserParams struct {
	ID            uuid.UUID
	FirstName     string
	LastName      string
	PhoneNo       pgtype.Text
	Gender        string
	RegNo         *string
	GithubProfile *string
	HostelBlock   *string
	RoomNo        *string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.PhoneNo,
		arg.Gender,
		arg.RegNo,
		arg.GithubProfile,
		arg.HostelBlock,
		arg.RoomNo,
	)
	return err
}

const verifyUser = `-- name: VerifyUser :exec
UPDATE users
SET is_verified = TRUE
WHERE email = $1
`

func (q *Queries) VerifyUser(ctx context.Context, email string) error {
	_, err := q.db.Exec(ctx, verifyUser, email)
	return err
}
