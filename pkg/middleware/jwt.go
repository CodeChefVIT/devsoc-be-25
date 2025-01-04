package middleware

import (
	"net/http"
	"os"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
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
