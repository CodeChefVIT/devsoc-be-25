package models

type CreatePanel struct {
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	PhoneNo       string `json:"phone_no" validate:"required,len=10"`
	Gender        string `json:"gender" validate:"required,len=1"`
	RegNo         string `json:"reg_no" validate:"required"`
	GithubProfile string `json:"github_profile" validate:"required,url"`
	Password      string `json:"password" validate:"required"`
}

type UpdateScore struct {
	Design         int    `json:"design" validate:"min=0,max=10"`
	Implementation int    `json:"implementation" validate:"min=0,max=10"`
	Presentation   int    `json:"presentation" validate:"min=0,max=10"`
	Round          int    `json:"round" validate:"required"`
	Innovation     int    `json:"innovation" validate:"min=0,max=10"`
	Teamwork       int    `json:"teamwork" validate:"min=0,max=10"`
	Comment        string `json:"comment"`
	TeamID         string `json:"team_id" validate:"uuid"`
}

type CreateScore struct {
	Design         int    `json:"design" validate:"min=0,max=10"`
	Implementation int    `json:"implementation" validate:"min=0,max=10"`
	Presentation   int    `json:"presentation" validate:"min=0,max=10"`
	Round          int    `json:"round" validate:""`
	Innovation     int    `json:"innovation" validate:"min=0,max=10"`
	Teamwork       int    `json:"teamwork" validate:"min=0,max=10"`
	Comment        string `json:"comment"`
	TeamID         string `json:"team_id" validate:"required,uuid"`
}
