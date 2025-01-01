package dto

type UserDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	RegNo    string `json:"reg_no"`
	PhoneNo  string `json:"phone_no"`
	College  string `json:"college"`
	Password string `json:"password"`
}
