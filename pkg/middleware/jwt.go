package middleware

import (
	"net/http"
	"os"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
			claims := token.Claims.(jwt.MapClaims)

			userId, _ := uuid.Parse(claims["user_id"].(string))

			teamId := uuid.NullUUID{}
			if claims["team_id"] != nil {
				teamId.UUID, _ = uuid.Parse(claims["team_id"].(string))
				teamId.Valid = true
			}

			c.Set("user", db.User{
				ID:         userId,
				Name:       claims["user_name"].(string),
				TeamID:     teamId,
				Email:      claims["email"].(string),
				IsVitian:   claims["is_vitian"].(bool),
				RegNo:      claims["reg_no"].(string),
				PhoneNo:    claims["phone_no"].(string),
				Role:       claims["role"].(string),
				IsLeader:   claims["is_leader"].(bool),
				College:    claims["college"].(string),
				IsVerified: claims["is_verified"].(bool),
				IsBanned:   claims["is_banned"].(bool),
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
