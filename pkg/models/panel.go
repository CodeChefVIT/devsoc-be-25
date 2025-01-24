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
	Design         int    `json:"design" validate:"required,min=0,max=10"`
	Implementation int    `json:"implementation" validate:"required,min=0,max=10"`
	Presentation   int    `json:"presentation" validate:"required,min=0,max=10"`
	Round          int    `json:"round" validate:"required"`
	TeamID         string `json:"team_id" validate:"required,uuid"`
}

type CreateScore struct {
	Design         int    `json:"design" validate:"required,min=0,max=10"`
	Implementation int    `json:"implementation" validate:"required,min=0,max=10"`
	Presentation   int    `json:"presentation" validate:"required,min=0,max=10"`
	Round          int    `json:"round" validate:"required"`
	TeamID         string `json:"team_id" validate:"required,uuid"`
}
