package dto

type Submission struct {
	GithubLink string `json:"github_link"`
	FigmaLink  string `json:"figma_link"`
	PptLink    string `json:"ppt_link"`
	OtherLink  string `json:"misc_link"`
	TeamID     string `json:"team_id"`
}
