package middleware

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token, ok := user.(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusUnauthorized, models.Response{
				Status: "failed",
				Data: map[string]string{
					"message": "Jwt is invalid or missing",
				},
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, models.Response{
				Status: "failed",
				Data: map[string]string{
					"message": "malformed JWT",
				},
			})
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, models.Response{
				Status: "failed",
				Data: map[string]string{
					"message": "Access denied. Not admin",
				},
			})
		}

		return next(c)
	}
}

func CheckPanel(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token, ok := user.(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusUnauthorized, models.Response{
				Status: "failed",
				Data: map[string]string{
					"message": "Jwt is invalid or missing",
				},
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, models.Response{
				Status: "failed",
				Data: map[string]string{
					"message": "malformed JWT",
				},
			})
		}

		role, ok := claims["role"].(string)
		if !ok || role != "panel" {
			return c.JSON(http.StatusForbidden, models.Response{
				Status: "failed",
				Data: map[string]string{
					"message": "Access denied. Not admin",
				},
			})
		}

		return next(c)
	}
}
