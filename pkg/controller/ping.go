package controller

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "pong",
		},
	})
}
