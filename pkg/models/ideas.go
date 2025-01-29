package models

type UpdateIdeaRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=10,max=1000"`
	Track       string `json:"track" validate:"required"`
	IsSelected  bool   `json:"is_selected"`
}
