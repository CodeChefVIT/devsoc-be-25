package dto

type Submission struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Track       string `json:"track"`
	GithubLink  string `json:"github_link"`
	FigmaLink   string `json:"figma_link"`
	PptLink     string `json:"ppt_link"`
	OtherLink   string `json:"other_link"`
	TeamID      string `json:"team_id"`
}
