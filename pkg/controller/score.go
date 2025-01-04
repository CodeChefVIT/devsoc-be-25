package controller

import (
	"database/sql"
	"errors"
	"net/http"

	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetScore(c echo.Context) error {
	ctx := c.Request().Context()
	teamId := c.Param("teamid")

	teamUuid, err := uuid.Parse(teamId)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Invalid team ID format",
				"error":   err.Error()},
		})
	}

	teamScore, err := utils.Queries.GetTeamScores(ctx, teamUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "No scores found for the team",
					"error":   err.Error(),
				},
			})
		}

		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to fetch team scores",
				"error":   err.Error(),
			},
		})
	}

	scores := make([]models.GetScore, len(teamScore))
	for i := 0; i < len(teamScore); i++ {
		scores[i] = models.GetScore{
			Id:             teamScore[i].ID.String(),
			TeamID:         teamId,
			Design:         int(teamScore[i].Design),
			Implementation: int(teamScore[i].Implementation),
			Presentation:   int(teamScore[i].Implementation),
			Round:          int(teamScore[i].Round),
		}
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "Scores fetched successfully",
			"scores":  scores,
		},
	})
}
