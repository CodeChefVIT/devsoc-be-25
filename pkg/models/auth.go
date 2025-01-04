package models

type SignupRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Email    string `json:"email" validate:"required,email,endswith=@vitstudent.ac.in"`
	RegNo    string `json:"reg_no" validate:"required"`
	Password string `json:"password" validate:"required"`
	PhoneNo  string `json:"phone_no" validate:"required"`
}

type SendOTPRequest struct {
	Email string `json:"email" validate:"required,email,endswith=@vitstudent.ac.in"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email,endswith=@vitstudent.ac.in"`
	OTP   string `json:"otp" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,endswith=@vitstudent.ac.in"`
	Password string `json:"password" validate:"required"`
}

type UpdatePasswordRequest struct {
	Email       string `json:"email" validate:"required,email,endswith=@vitstudent.ac.in"`
	NewPassword string `json:"new_password" validate:"required"`
	OTP         string `json:"otp" validate:"required"`
}
