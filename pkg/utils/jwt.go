package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID     uuid.UUID     `json:"user_id"`
	UserName   string        `json:"user_name"`
	TeamID     uuid.NullUUID `json:"team_id"`
	Email      string        `json:"email"`
	IsVitian   bool          `json:"is_vitian"`
	RegNo      string        `json:"reg_no"`
	PhoneNo    string        `json:"phone_no"`
	Role       string        `json:"role"`
	IsLeader   bool          `json:"is_leader"`
	College    string        `json:"college"`
	IsVerified bool          `json:"is_verified"`
	IsBanned   bool          `json:"is_banned"`
	jwt.RegisteredClaims
}

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(user *db.User) (string, error) {
	claims := JWTClaims{
		user.ID,
		user.Name,
		user.TeamID,
		user.Email,
		user.IsVitian,
		user.RegNo,
		user.PhoneNo,
		user.Role,
		user.IsLeader,
		user.College,
		user.IsVerified,
		user.IsBanned,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
