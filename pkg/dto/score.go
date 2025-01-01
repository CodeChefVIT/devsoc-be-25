package dto

type ScoreDto struct {
	Design         int    `json:"design"`
	Implementation int    `json:"implementation"`
	Presentation   int    `json:"presentation"`
	Round          int    `json:"round"`
	TeamID         string `json:"team_id"`
}
