package models

type GetScore struct {
	Id             string `json:"id" validate:"required"`
	Design         int    `json:"design" validate:"required,min=0,max=10"`
	Implementation int    `json:"implementation" validate:"required,min=0,max=10"`
	Presentation   int    `json:"presentation" validate:"required,min=0,max=10"`
	Round          int    `json:"round" validate:"required"`
	Innovation     int    `json:"innovation" validate:"required,min=0,max=10"`
	Teamwork       int    `json:"teamwork" validate:"required,min=0,max=10"`
	Comment        string `json:"comment"`
	TeamID         string `json:"team_id" validate:"required,uuid"`
}
