// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: export.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const exportAllTeams = `-- name: ExportAllTeams :many
SELECT id, name, number_of_people, round_qualified, code
FROM teams
`

type ExportAllTeamsRow struct {
	ID             uuid.UUID
	Name           string
	NumberOfPeople int32
	RoundQualified pgtype.Int4
	Code           string
}

func (q *Queries) ExportAllTeams(ctx context.Context) ([]ExportAllTeamsRow, error) {
	rows, err := q.db.Query(ctx, exportAllTeams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExportAllTeamsRow
	for rows.Next() {
		var i ExportAllTeamsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.NumberOfPeople,
			&i.RoundQualified,
			&i.Code,
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

const exportAllUsers = `-- name: ExportAllUsers :many
SELECT id, team_id, first_name, last_name, email, phone_no, gender, reg_no, hostel_block, room_no, github_profile, password, role, is_leader, is_verified, is_banned, is_profile_complete
FROM users
`

type ExportAllUsersRow struct {
	ID                uuid.UUID
	TeamID            uuid.NullUUID
	FirstName         string
	LastName          string
	Email             string
	PhoneNo           pgtype.Text
	Gender            string
	RegNo             *string
	HostelBlock       *string
	RoomNo            *string
	GithubProfile     *string
	Password          string
	Role              string
	IsLeader          bool
	IsVerified        bool
	IsBanned          bool
	IsProfileComplete bool
}

func (q *Queries) ExportAllUsers(ctx context.Context) ([]ExportAllUsersRow, error) {
	rows, err := q.db.Query(ctx, exportAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExportAllUsersRow
	for rows.Next() {
		var i ExportAllUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.PhoneNo,
			&i.Gender,
			&i.RegNo,
			&i.HostelBlock,
			&i.RoomNo,
			&i.GithubProfile,
			&i.Password,
			&i.Role,
			&i.IsLeader,
			&i.IsVerified,
			&i.IsBanned,
			&i.IsProfileComplete,
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
