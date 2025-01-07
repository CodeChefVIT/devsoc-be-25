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
	teamId := c.Param("teamId")
	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid user"},
		})
	}

	teamUuid, err := uuid.Parse(teamId)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid team ID format"},
		})
	}

	if teamUuid == uuid.Nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User is not part of any team"},
		})
	}

	if user.TeamID.UUID.String() != teamId {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User is not part of the team"},
		})
	}
	submission, err := utils.Queries.GetSubmissionByTeamID(ctx, teamUuid)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
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
	var req models.CreateSubmissionRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid request body"},
		})
	}

	teamUuid, err := uuid.Parse(req.TeamID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid team ID format"},
		})
	}

	if teamUuid == uuid.Nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User is not part of any team"},
		})
	}

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid user"},
		})
	}

	if user.TeamID.UUID.String() != req.TeamID {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User is not part of the team"},
		})
	}

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
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, models.Response{
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
	teamId := c.Param("teamId")
	var req models.UpdateSubmissionRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid request body"},
		})
	}

	teamUuid, err := uuid.Parse(teamId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid team ID format"},
		})
	}

	if teamUuid == uuid.Nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User is not part of any team"},
		})
	}

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid user"},
		})
	}

	if user.TeamID.UUID.String() != teamId {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User is not part of the team"},
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
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
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
	teamId := c.Param("teamId")

	teamUuid, err := uuid.Parse(teamId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid team ID format"},
		})
	}

	if teamUuid == uuid.Nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User is not part of any team"},
		})
	}

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Invalid user"},
		})
	}

	if user.TeamID.UUID.String() != teamId {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "User is not part of the team"},
		})
	}

	err = utils.Queries.DeleteSubmission(ctx, teamUuid)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   map[string]string{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   map[string]string{"message": "Submission deleted successfully"},
	})
}
