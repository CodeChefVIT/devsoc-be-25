package controller

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/labstack/echo/v4"
)

func GetTeamId(c echo.Context) error {
	ctx := c.Request().Context()
	teamCode := c.Param("teamcode")

	teamId, err := utils.Queries.GetTeamIDByCode(ctx, teamCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to fetch team ID",
				"error":   err.Error(),
			},
		})
	}

	response := map[string]interface{}{
		"teamId": teamId,
	}
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   response,
	})
}
