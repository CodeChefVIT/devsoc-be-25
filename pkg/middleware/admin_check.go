package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token, ok := user.(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "JWT is invalid or missing",
				"status":  "fail",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Malformed JWT",
				"status":  "fail",
			})
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "Not an admin",
				"status":  "fail",
			})
		}

		return next(c)
	}
}
