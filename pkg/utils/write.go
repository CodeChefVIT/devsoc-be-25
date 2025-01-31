package utils

import (
	//"context"
	"context"
	"fmt"
	"math/rand"
	"time"

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

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

func SendTeamEmail(ctx context.Context, emails []string) error {
	for i := 0; i < len(emails); i++ {
		err := SendEmail(emails[i], "Team Deleted", fmt.Sprint(Config.TeamDeleteTemplate))
		if err != nil {
			return nil
		}
	}
	return nil
}
