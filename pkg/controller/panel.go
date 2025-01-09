package controller

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/dto"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetSubmission(c echo.Context) error {
	ctx := c.Request().Context()
	teamId := c.Param("teamId")

	teamUuid, err := uuid.Parse(teamId)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID format"})
	}

	submission, err := utils.Queries.GetSubmissionByTeamID(ctx, teamUuid)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.Submission{
		GithubLink: submission.GithubLink,
		FigmaLink:  submission.FigmaLink,
		PptLink:    submission.PptLink,
		OtherLink:  submission.OtherLink,
		TeamID:     submission.TeamID.String(),
	},
	)
}

func UpdateScore(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	points := new(models.UpdateScore)
	if err := c.Bind(points); err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.JSON(http.StatusBadRequest,
			map[string]string{
				"message": "Invalid request body",
				"error":   err.Error(),
			},
		)
	}

	if err := utils.Validate.Struct(points); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	scoreid, err := uuid.Parse(id)
	if err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid score ID format",
			"error":   err.Error(),
		},
		)
	}
	teamid, err := uuid.Parse(points.TeamID)
	if err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid team ID format",
			"error":   err.Error(),
		},
		)
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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to update score",
			"error":   err.Error(),
		},
		)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Score updated successfully",
	},
	)
}

func DeleteScore(c echo.Context) error {
	ctx := c.Request().Context()
	scoreId := c.Param("id")

	scoreUuid, err := uuid.Parse(scoreId)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid score ID format",
			"error":   err.Error(),
		},
		)
	}

	err = utils.Queries.DeleteScore(ctx, scoreUuid)
	if err != nil {
		logger.Errorf(logger.DatabaseError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to delete score",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "Score deleted successfully",
		},
	})
}

func CreateScore(c echo.Context) error {
	ctx := c.Request().Context()

	points := new(models.CreateScore)
	if err := c.Bind(points); err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
			"error":   err.Error(),
		},
		)
	}

	if err := utils.Validate.Struct(points); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	teamid, err := uuid.Parse(points.TeamID)
	if err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid team ID format",
			"error":   err.Error(),
		},
		)
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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to create score",
			"error":   err.Error(),
		},
		)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Score created successfully"})
}
