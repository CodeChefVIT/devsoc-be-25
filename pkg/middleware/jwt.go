package middleware

import (
	"net/http"
	"os"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Protected() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	})
}

func JWTMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		SigningKey:  []byte(os.Getenv("JWT_SECRET")),
		TokenLookup: "cookie:jwt",
		SuccessHandler: func(c echo.Context) {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(utils.JWTClaims)
			c.Set("user", &db.User{
				ID:         claims.UserID,
				Name:       claims.UserName,
				TeamID:     claims.TeamID,
				Email:      claims.Email,
				IsVitian:   claims.IsVitian,
				RegNo:      claims.RegNo,
				PhoneNo:    claims.PhoneNo,
				Role:       claims.Role,
				IsLeader:   claims.IsLeader,
				College:    claims.College,
				IsVerified: claims.IsVerified,
				IsBanned:   claims.IsBanned,
			})
		},
		ErrorHandler: func(c echo.Context, err error) error {
			if err == echojwt.ErrJWTMissing {
				return c.JSON(http.StatusUnauthorized, models.Response{
					Status: "fail",
					Data: map[string]string{
						"error": "Missing or malformed JWT",
					},
				})
			}

			return c.JSON(http.StatusUnauthorized, models.Response{
				Status: "fail",
				Data: map[string]string{
					"error": "Invalid or expired token",
				},
			})
		},
	}

	return echojwt.WithConfig(config)
}
