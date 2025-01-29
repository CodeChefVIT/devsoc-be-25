package models

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CompleteProfileRequest struct {
	FirstName     string `json:"first_name" validate:"required,alphanum"`
	LastName      string `json:"last_name" validate:"required,alphanum"`
	PhoneNo       string `json:"phone_no" validate:"required,len=10"`
	Gender        string `json:"gender" validate:"required,len=1"`
	RegNo         string `json:"reg_no" validate:"required"`
	GithubProfile string `json:"github_profile" validate:"required,url"`
	HostelBlock   string `json:"hostel_block" validate:"required"`
	RoomNo        string `json:"room_no" validate:"required"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdatePasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required"`
	OTP         string `json:"otp" validate:"required"`
}
