package controller

import (
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
		return c.String(http.StatusBadRequest, err.Error())
	}

	teamScore, err := utils.Queries.GetTeamScores(ctx, teamUuid)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
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

	return c.JSON(http.StatusOK, scores)
}
