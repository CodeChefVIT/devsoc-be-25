package controller

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/dto"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
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

	scores := make([]dto.ScoreDto, len(teamScore))
	for i := 0; i < len(teamScore); i++ {
		scores[i] = dto.ScoreDto{
			TeamID:         teamId,
			Design:         int(teamScore[i].Design),
			Implementation: int(teamScore[i].Implementation),
			Presentation:   int(teamScore[i].Implementation),
			Round:          int(teamScore[i].Round),
		}
	}

	return c.JSON(http.StatusOK, scores)
}

func CreateScore(c echo.Context) error {
	ctx := c.Request().Context()

	points := new(dto.ScoreDto)
	if err := c.Bind(points); err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	teamid, err := uuid.Parse(points.TeamID)
	if err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	score := db.CreateScoreParams{
		TeamID:         teamid,
		Round:          int32(points.Round),
		Presentation:   int32(points.Presentation),
		Implementation: int32(points.Implementation),
		Design:         int32(points.Design),
	}
	score.ID, _ = uuid.NewV7()

	err = utils.Queries.CreateScore(ctx, score)
	if err != nil {
		logger.Errorf(logger.DatabaseError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "all good")
}
