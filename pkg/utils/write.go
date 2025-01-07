package utils

import (
	"fmt"

	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/labstack/echo/v4"
)

type StandardResponse struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func WriteJSON(c echo.Context, status int, message interface{}) error {
	response := StandardResponse{
		Status:  status,
		Message: message,
	}
	logger.Infof(fmt.Sprintf("Response: %+v", response))
	return c.JSON(status, response)
}

func WriteError(c echo.Context, status int, err error) error {

	response := StandardResponse{
		Status: status,
		Error:  err.Error(),
	}
	logger.Errorf(fmt.Sprintf("Error Response: %+v", response))
	return c.JSON(status, response)
}
