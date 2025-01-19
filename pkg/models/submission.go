package models

type CreateSubmissionRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Track       string `json:"track" validate:"required"`
	GithubLink  string `json:"github_link" validate:"required,url"`
	FigmaLink   string `json:"figma_link" validate:"required,url"`
	OtherLink   string `json:"other_link" validate:"omitempty,url"`
}

type UpdateSubmissionRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Track       string `json:"track" validate:"required"`
	GithubLink  string `json:"github_link" validate:"required,url"`
	FigmaLink   string `json:"figma_link" validate:"required,url"`
	OtherLink   string `json:"other_link" validate:"omitempty,url"`
}
