package models

type BanUserReq struct {
	Email string `json:"email" validate:"required,email"`
}
