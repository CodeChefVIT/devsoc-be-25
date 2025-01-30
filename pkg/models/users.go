package models

type BanUserReq struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdateGithub struct {
	Github string `json:"github" validate:"required"`
}
