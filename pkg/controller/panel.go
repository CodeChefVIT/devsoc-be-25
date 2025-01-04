package controller

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UpdateScore(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	points := new(models.UpdateScore)
	if err := c.Bind(points); err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := utils.Validate.Struct(points); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "400",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	scoreid, err := uuid.Parse(id)
	if err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	teamid, err := uuid.Parse(points.TeamID)
	if err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	score := db.UpdateScoreParams{
		ID:             scoreid,
		TeamID:         teamid,
		Design:         int32(points.Design),
		Implementation: int32(points.Implementation),
		Presentation:   int32(points.Presentation),
		Round:          int32(points.Round),
	}

	err = utils.Queries.UpdateScore(ctx, score)
	if err != nil {
		logger.Errorf(logger.DatabaseError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	return utils.WriteJSON(c, http.StatusOK, "Score Updated")
}

func DeleteScore(c echo.Context) error {
	ctx := c.Request().Context()
	scoreId := c.Param("id")

	scoreUuid, err := uuid.Parse(scoreId)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = utils.Queries.DeleteScore(ctx, scoreUuid)
	if err != nil {
		logger.Errorf(logger.DatabaseError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	return utils.WriteJSON(c, http.StatusOK, "Score Deleted")
}

func CreateScore(c echo.Context) error {
	ctx := c.Request().Context()

	points := new(models.CreateScore)
	if err := c.Bind(points); err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := utils.Validate.Struct(points); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "400",
			Data:   utils.FormatValidationErrors(err),
		})
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

	return utils.WriteJSON(c, http.StatusOK, "Score created")
}
