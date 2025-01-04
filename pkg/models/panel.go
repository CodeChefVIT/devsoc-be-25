package models

type CreatePanel struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	RegNo    string `json:"reg_no" validate:"required"`
	PhoneNo  string `json:"phone_no" validate:"required"`
	Password string `json:"password" validate:"required"`
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
