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

func GetUserSubmission(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid user"},
		})
	}

	teamUuid := user.TeamID.UUID
	if !user.TeamID.Valid {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid team ID format"},
		})
	}

	submission, err := utils.Queries.GetSubmissionByTeamID(ctx, teamUuid)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: dto.Submission{
			GithubLink: submission.GithubLink,
			FigmaLink:  submission.FigmaLink,
			PptLink:    submission.PptLink,
			OtherLink:  submission.OtherLink,
			TeamID:     submission.TeamID.String(),
		},
	})
}

func CreateSubmission(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid user"},
		})
	}

	if !ok || !user.TeamID.Valid || !user.IsLeader {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User does not belong to any team or is not team leader",
			},
		})
	}

	var req models.CreateSubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid request body"},
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	teamUuid := user.TeamID.UUID

	submission_id, _ := uuid.NewV7()
	submission, err := utils.Queries.CreateSubmission(ctx, db.CreateSubmissionParams{
		ID:         submission_id,
		TeamID:     teamUuid,
		GithubLink: req.GithubLink,
		FigmaLink:  req.FigmaLink,
		PptLink:    req.PptLink,
		OtherLink:  req.OtherLink,
	})

	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, &models.Response{
		Status: "success",
		Data: dto.Submission{
			TeamID:     submission.TeamID.String(),
			GithubLink: submission.GithubLink,
			FigmaLink:  submission.FigmaLink,
			PptLink:    submission.PptLink,
			OtherLink:  submission.OtherLink,
		},
	})
}

func UpdateSubmission(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid user"},
		})
	}

	if !user.TeamID.Valid || !user.IsLeader {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User does not belong to any team or is not team leader",
			},
		})
	}

	var req models.UpdateSubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid request body"},
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	teamUuid := user.TeamID.UUID
	if !user.TeamID.Valid {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid team ID format"},
		})
	}

	submission, err := utils.Queries.UpdateSubmission(ctx, db.UpdateSubmissionParams{
		TeamID:     teamUuid,
		GithubLink: req.GithubLink,
		FigmaLink:  req.FigmaLink,
		PptLink:    req.PptLink,
		OtherLink:  req.OtherLink,
	})

	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: dto.Submission{
			TeamID:     submission.TeamID.String(),
			GithubLink: submission.GithubLink,
			FigmaLink:  submission.FigmaLink,
			PptLink:    submission.PptLink,
			OtherLink:  submission.OtherLink,
		},
	})
}

func DeleteSubmission(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid user"},
		})
	}

	if !user.IsLeader {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User is not team leader",
			},
		})
	}

	teamUuid := user.TeamID.UUID
	if !user.TeamID.Valid {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid team ID format"},
		})
	}

	err := utils.Queries.DeleteSubmission(ctx, teamUuid)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data:   map[string]string{"message": "Submission deleted successfully"},
	})
}
