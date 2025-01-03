package models

import "github.com/google/uuid"

type Team struct {
	ID               uuid.UUID `json:"team_id" db:"id"`
	Name             string    `json:"team_name" db:"name"`
	NumberOfPeople   int64     `json:"number_of_people" db:"number_of_people"`
	Submission       uuid.UUID `json:"submission" db:"submission"`
	RoundQualified   int       `json:"round_qualified" db:"round_qualified" default:"0"`
	Code             string    `json:"code" db:"code"`
}

type CreateTeam struct {
	Name string `json:"name" validate:"required"`
}

type JoinTeam struct {
	Code string `json:"code" validate:"reuqired" `
}

type KickMember struct {
	UserID uuid.UUID `json:"id" validate:"required"`
}

type GetTeams struct {
	ID uuid.UUID `json:"team_id" db:"id"`
	Name             string    `json:"team_name" db:"name"`
	NumberOfPeople   int64     `json:"number_of_people" db:"number_of_people"`
	Submission       uuid.UUID `json:"submission" db:"submission"`
	RoundQualified   int       `json:"round_qualified" db:"round_qualified" default:"0"`
	Code             string    `json:"code" db:"code"`
}

type LeaveTeam struct {
	UserID uuid.UUID `json:"id" validate:"required"`
}