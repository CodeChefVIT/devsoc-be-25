package controller

import (
	"errors"
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/dto"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func GetUserSubmission(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: "Invalid user",
		})
	}

	teamUuid := user.TeamID.UUID
	if !user.TeamID.Valid {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: "Invalid team ID format",
		})
	}

	submission, err := utils.Queries.GetSubmissionByTeamID(ctx, teamUuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, &models.Response{
				Status: "fail",
				Message: "Submission not found",
			})
		}
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: dto.Submission{
			Title:       submission.Title,
			Description: submission.Description,
			Track:       submission.Track,
			GithubLink:  submission.GithubLink,
			FigmaLink:   submission.FigmaLink,
			OtherLink:   submission.OtherLink,
			TeamID:      submission.TeamID.String(),
		},
	})
}

func CreateSubmission(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: "Invalid user",
		})
	}

	if !ok || !user.TeamID.Valid || !user.IsLeader {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Message:   "User does not belong to any team or is not team leader",
		})
	}

	var req models.CreateSubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: "Invalid request body",
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
		ID:          submission_id,
		Title:       req.Title,
		Description: req.Description,
		Track:       req.Track,
		TeamID:      teamUuid,
		GithubLink:  req.GithubLink,
		FigmaLink:   req.FigmaLink,
		OtherLink:   req.OtherLink,
	})

	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: "failed to create submission",
		})
	}

	return c.JSON(http.StatusCreated, &models.Response{
		Status: "success",
		Data: dto.Submission{
			TeamID:      submission.TeamID.String(),
			Title:       submission.Title,
			Description: submission.Description,
			Track:       submission.Track,
			GithubLink:  submission.GithubLink,
			FigmaLink:   submission.FigmaLink,
			OtherLink:   submission.OtherLink,
		},
	})
}

func UpdateSubmission(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: "Invalid user",
		})
	}

	if !user.TeamID.Valid || !user.IsLeader {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Message: "User does not belong to any team or is not team leader",
		})
	}

	var req models.UpdateSubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: "Invalid request body",
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
			Message:"Invalid team ID format",
		})
	}

	submission, err := utils.Queries.UpdateSubmission(ctx, db.UpdateSubmissionParams{
		TeamID:      teamUuid,
		Title:       req.Title,
		Description: req.Description,
		Track:       req.Track,
		GithubLink:  req.GithubLink,
		FigmaLink:   req.FigmaLink,
		OtherLink:   req.OtherLink,
	})

	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: "failed to update submission",
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: dto.Submission{
			TeamID:      submission.TeamID.String(),
			Title:       submission.Title,
			Description: submission.Description,
			Track:       submission.Track,
			GithubLink:  submission.GithubLink,
			FigmaLink:   submission.FigmaLink,
			OtherLink:   submission.OtherLink,
		},
	})
}

func DeleteSubmission(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message:   "Invalid user",
		})
	}

	if !user.IsLeader {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Message: "User is not team leader",
		})
	}

	teamUuid := user.TeamID.UUID
	if !user.TeamID.Valid {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: "Invalid team ID format",
		})
	}

	err := utils.Queries.DeleteSubmission(ctx, teamUuid)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Message: "failed to delete submission",
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Message: "Submission deleted successfully",
	})
}
