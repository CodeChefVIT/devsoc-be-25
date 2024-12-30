package utils

import (
	"github.com/labstack/echo/v4"
)

func WriteError (r echo.Context,status int, err error) error {
	return r.JSON(status, map[string]string{"error":err.Error()})
}

func WriteJSON(r echo.Context, status int, v any) error {
	return r.JSON(status, map[string]interface{}{"message":v})
}