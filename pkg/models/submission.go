package models

type CreateSubmissionRequest struct {
	GithubLink string `json:"github_link" validate:"required,url"`
	FigmaLink  string `json:"figma_link" validate:"required,url"`
	PptLink    string `json:"ppt_link" validate:"required,url"`
	OtherLink  string `json:"other_link" validate:"omitempty,url"`
}

type UpdateSubmissionRequest struct {
	GithubLink string `json:"github_link" validate:"required,url"`
	FigmaLink  string `json:"figma_link" validate:"required,url"`
	PptLink    string `json:"ppt_link" validate:"omitempty,url"`
	OtherLink  string `json:"other_link" validate:"omitempty,url"`
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
