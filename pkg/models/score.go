package models

import "github.com/google/uuid"

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

type Round struct {
	Round          int `json:"round"`
	Design         int `json:"design"`
	Implementation int `json:"implementation"`
	Presentation   int `json:"presentation"`
	Innovation     int `json:"innovation"`
	Teamwork       int `json:"teamwork"`
	RoundTotal     int `json:"round_total"`
}

type TeamLeaderboard struct {
	TeamID       uuid.UUID `json:"team_id"`
	TeamName     string  `json:"team_name"`
	Rounds       []Round `json:"rounds"`
	OverallTotal int     `json:"overall_total"`
}