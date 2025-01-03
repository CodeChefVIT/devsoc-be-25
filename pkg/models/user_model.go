package models

import(

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	TeamID     *uuid.UUID `json:"team_id" db:"team_id"`
	Email      string    `json:"email" db:"email"`
	IsVitian   bool      `json:"is_vitian" db:"is_vitian"`
	RegNo      string    `json:"reg_no" db:"reg_no"`
	Password   string    `json:"password" db:"password"`
	PhoneNo    string    `json:"phone_no" db:"phone_no"`
	IsLeader   bool      `json:"is_leader" db:"is_leader"`
	College    string    `json:"college" db:"college"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
}

