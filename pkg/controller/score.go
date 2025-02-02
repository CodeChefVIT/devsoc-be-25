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
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: err.Error()},
		)
	}

	teamScore, err := utils.Queries.GetTeamScores(ctx, teamUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, &models.Response{
				Status:  "fail",
				Message: err.Error()},
			)
		}

		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	if len(teamScore) == 0 {
		return c.JSON(http.StatusNotFound, &models.Response{
			Status:  "fail",
			Message: "No scores found for the given team ID"},
		)
	}

	scores := make([]models.GetScore, len(teamScore))
	for i := 0; i < len(teamScore); i++ {
		scores[i] = models.GetScore{
			Id:             teamScore[i].ID.String(),
			TeamID:         teamId,
			Design:         int(teamScore[i].Design),
			Implementation: int(teamScore[i].Implementation),
			Presentation:   int(teamScore[i].Presentation),
			Innovation:     int(teamScore[i].Innovation),
			Teamwork:       int(teamScore[i].Teamwork),
			Comment:        getSafeString(teamScore[i].Comment),
			Round:          int(teamScore[i].Round),
		}
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Message: "Scores fetched successfully",
		Data: map[string]interface{}{
			"scores":  scores,
		},
	})
}
