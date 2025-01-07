package controller

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusForbidden, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User not found",
			},
		})
	}

	logger.Infof("User: %+v", user)
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "pong",
		},
	})
}
