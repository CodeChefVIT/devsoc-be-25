package middleware

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"

	"github.com/labstack/echo/v4"
)

func CheckUserBan(next echo.HandlerFunc) echo.HandlerFunc {
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
		if user.IsBanned {
			return c.JSON(http.StatusForbidden, &models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "banned",
				},
			})
		}
		return next(c)
	}
}

func CheckTeamBan(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(db.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, &models.Response{
				Status:  "fail",
				Message: "unauthorized",
			})
		}
		team, err := utils.Queries.GetTeamById(c.Request().Context(), user.TeamID.UUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.Response{
				Status:  "fail",
				Message: "team not found",
			})
		}
		if team.IsBanned {
			return c.JSON(http.StatusForbidden, models.Response{
				Status:  "fail",
				Message: "Team is banned",
			})
		}
		return next(c)
	}
}

