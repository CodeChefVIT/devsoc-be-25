package middleware

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/labstack/echo/v4"
)

func CheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(db.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, &models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "Jwt is invalid or missing",
				},
			})
		}

		if user.Role != "admin" {
			return c.JSON(http.StatusForbidden, &models.Response{
				Status: "fail",
				Data: map[string]string{
					"error": "Access denied. Not admin",
				},
			})
		}

		return next(c)
	}
}

func CheckPanel(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(db.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, &models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "Jwt is invalid or missing",
				},
			})
		}

		if user.Role != "panel" {
			return c.JSON(http.StatusForbidden, &models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "Access denied. Not panel",
				},
			})
		}

		return next(c)
	}
}
