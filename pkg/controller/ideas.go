package controller

import (
	"context"
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

func CreateIdea(c echo.Context) error {
	user, ok := c.Get("user").(db.User)
	if !ok || !user.TeamID.Valid {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status:  "fail",
			Message: "User does not belong to any team",
		})
	}

	if !user.IsLeader {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status:  "fail",
			Message: "Only team leaders can create ideas",
		})
	}

	var input db.CreateIdeaParams
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "invalid request body",
		})
	}

	input.ID = uuid.New()
	input.TeamID = user.TeamID.UUID

	_, err := utils.Queries.GetIdeaByTeamID(context.Background(), input.TeamID)
	if err == nil {
		return c.JSON(http.StatusConflict, &models.Response{
			Status:  "fail",
			Message: "An idea already exists for this team",
		})
	}

	_, err = utils.Queries.CreateIdea(context.Background(), input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, &models.Response{
				Status:  "fail",
				Message: "Idea not found",
			})
		}

		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to create idea",
		})
	}

	return c.JSON(http.StatusCreated, &models.Response{
		Status:  "success",
		Message: "Idea created successfully",
	})
}

func UpdateIdea(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Invalid user",
		})
	}

	if !user.TeamID.Valid || !user.IsLeader {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status:  "fail",
			Message: "User does not belong to any team or is not team leader",
		})
	}

	var req models.UpdateIdeaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
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
			Status:  "fail",
			Message: "Invalid team ID format",
		})
	}

	err := utils.Queries.UpdateIdea(ctx, db.UpdateIdeaParams{
		TeamID:      teamUuid,
		Title:       req.Title,
		Description: req.Description,
		Track:       req.Track,
	})
	if err != nil {
		logger.Errorf("Failed to update idea: %v", err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to update idea",
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Idea updated successfully",
		Data: dto.Idea{
			Title:       req.Title,
			Description: req.Description,
			Track:       req.Track,
		},
	})
}

func GetIdea(c echo.Context) error {
	user, ok := c.Get("user").(db.User)
	if !ok || !user.TeamID.Valid {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status:  "fail",
			Message: "Please join a team or create one",
		})
	}

	ideas, err := utils.Queries.GetIdeaByTeamID(context.Background(), user.TeamID.UUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, &models.Response{
				Status:  "success",
				Message: "Idea not found",
			})
		}

		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "some error occurred",
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "ideas fetched successfully",
		Data: dto.Idea{
			Title:       ideas.Title,
			Description: ideas.Description,
			Track:       ideas.Track,

		},
	})
}
