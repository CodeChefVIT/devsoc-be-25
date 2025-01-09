package models

type SignupRequest struct {
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	PhoneNo       string `json:"phone_no" validate:"required,len=10"`
	Gender        string `json:"gender" validate:"required,len=1"`
	RegNo         string `json:"reg_no" validate:"required"`
	VitEmail      string `json:"vit_email" validate:"required,email,endswith=@vitstudent.ac.in"`
	HostelBlock   string `json:"hostel_block" validate:"required"`
	RoomNumber    int    `json:"room_no" validate:"required"`
	GithubProfile string `json:"github_profile" validate:"required,url"`
	Password      string `json:"password" validate:"required"`
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
