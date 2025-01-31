package models

import "github.com/google/uuid"

type Team struct {
	ID             uuid.UUID `json:"team_id" db:"id"`
	Name           string    `json:"team_name" db:"name"`
	NumberOfPeople int64     `json:"number_of_people" db:"number_of_people"`
	Submission     uuid.UUID `json:"submission" db:"submission"`
	RoundQualified int       `json:"round_qualified" db:"round_qualified" default:"0"`
	Code           string    `json:"code" db:"code"`
	IsBanned       bool      `json:"is_banned" db:"is_banned" default:"false"`
}

type CreateTeam struct {
	Name string `json:"name" validate:"required"`
}

type JoinTeam struct {
	Code string `json:"code" validate:"required" `
}

type KickMember struct {
	UserID uuid.UUID `json:"id" validate:"required"`
}

type GetTeams struct {
	ID             uuid.UUID `json:"team_id" db:"id"`
	Name           string    `json:"team_name" db:"name"`
	NumberOfPeople int64     `json:"number_of_people" db:"number_of_people"`
	Submission     uuid.UUID `json:"submission" db:"submission"`
	RoundQualified int       `json:"round_qualified" db:"round_qualified" default:"0"`
	Code           string    `json:"code" db:"code"`
	IsBanned       bool      `json:"is_banned" db:"is_banned" default:"false"`
}

type LeaveTeam struct {
	Email string `json:"email" validate:"required"`
}

type DeleteTeam struct {
	UserID uuid.UUID `json:"id" validate:"required"`
}

type UpdateTeamName struct {
	Name string `json:"name" validate:"required,alphanum"`
}

type GetTeamMembers struct {
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	GithubProfile string `json:"github_profile" validate:"required,url"`
	VitEmail      string `json:"vit_email" validate:"required,email,endswith=@vitstudent.ac.in"`
	RegNo         string `json:"reg_no" validate:"required"`
	PhoneNo       string `json:"phone_no" validate:"required,len=10"`
}

type GetTeamUsers struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type BanTeam struct {
	TeamId uuid.UUID `json:"id" validate:"required"`
}

type UnBanTeam struct {
	TeamId uuid.UUID `json:"id" validate:"required"`
}

type TeamRoundQualified struct {
	TeamId         uuid.UUID `json:"id" validate:"required"`
	RoundQualified int       `json:"round_qualified" validate:"required"`
}
