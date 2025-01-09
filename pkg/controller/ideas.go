package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateIdea(c echo.Context) error {
	user, ok := c.Get("user").(db.User)
	if !ok || !user.TeamID.Valid {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User does not belong to any team",
			},
		})
	}

	if !user.IsLeader {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "Only team leaders can create ideas",
			},
		})
	}

	var input db.CreateIdeaParams
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Invalid request",
				"error":   err.Error(),
			},
		})
	}

	input.ID = uuid.New()
	input.TeamID = user.TeamID.UUID

	_, err := utils.Queries.GetIdeaByTeamID(context.Background(), input.TeamID)
	if err == nil {
		return c.JSON(http.StatusConflict, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "An idea already exists for this team",
				"error":   "Duplicate idea is not allowed",
			},
		})
	}

	_, err = utils.Queries.CreateIdea(context.Background(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to create idea",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusCreated, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "Idea created successfully",
			"error":   nil,
		},
	})
}

func UpdateIdea(c echo.Context) error {

	user, ok := c.Get("user").(db.User)
	if !ok || !user.TeamID.Valid {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User does not belong to any team",
			},
		})
	}

	if !user.IsLeader {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "Only team leaders can update ideas",
			},
		})
	}

	ideaID := c.Param("id")
	id, err := uuid.Parse(ideaID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Invalid ID format",
				"error":   err.Error(),
			},
		})
	}

	existingIdea, err := utils.Queries.GetIdea(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Idea not found",
				"error":   err.Error(),
			},
		})
	}
	fmt.Println(existingIdea.TeamID)
	if existingIdea.TeamID != user.TeamID.UUID {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User's team does not match the idea's team",
			},
		})
	}

	var input db.UpdateIdeaParams
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Invalid request",
				"error":   err.Error(),
			},
		})
	}

	input.ID = id
	err = utils.Queries.UpdateIdea(context.Background(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to update idea",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "Idea updated successfully",
		},
	})
}

func GetIdea(c echo.Context) error {
	user, ok := c.Get("user").(db.User)
	if !ok || !user.TeamID.Valid {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User does not belong to any team",
			},
		})
	}

	ideas, err := utils.Queries.GetIdeaByTeamID(context.Background(), user.TeamID.UUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to fetch ideas",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "Ideas fetched successfully",
			"ideas":   ideas,
		},
	})
}
