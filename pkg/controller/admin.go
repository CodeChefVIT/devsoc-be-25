package controller

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(c echo.Context) error {
	ctx := c.Request().Context()
	limitParam := c.QueryParam("limit")
	cursor := c.QueryParam("cursor")

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"error": "failed to convert to integer",
			},
		})
	}

	var cursorUUID uuid.UUID
	if cursor == "" {
		cursorUUID = uuid.Nil // uuid.Nil is equivalent to "00000000-0000-0000-0000-000000000000"
	} else {
		cursorUUID, err = uuid.Parse(cursor)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid UUID for cursor",
			})
		}
	}

	users, err := utils.Queries.GetAllUsers(ctx, db.GetAllUsersParams{
		Limit: int32(limit),
		ID:    cursorUUID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to fetch users",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "Users fetched successfully",
			"users":   users,
		},
	})
}

func GetAllVitians(c echo.Context) error {
	users, err := utils.Queries.GetAllVitians(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to fetch users",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "Users fetched successfully",
			"users":   users,
		},
	})
}

func GetUsersByEmail(c echo.Context) error {
	email := c.Param("email")
	user, err := utils.Queries.GetUserByEmail(c.Request().Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, &models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "User not found",
					"error":   err.Error(),
				},
			})

		}
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}
	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "User fetched successfully",
			"user":    user,
		},
	})
}

func BanUser(c echo.Context) error {
	var payload models.BanUserReq

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Improper request",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	user, err := utils.Queries.GetUserByEmail(c.Request().Context(), payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, &models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "User does not exist",
					"error":   err.Error(),
				},
			})
		}
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Queries.BanUser(context.Background(), user.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "user banned successfully",
		},
	})
}

func UnbanUser(c echo.Context) error {
	var payload models.BanUserReq

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Improper request",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	user, err := utils.Queries.GetUserByEmail(c.Request().Context(), payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, &models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "User does not exist",
					"error":   err.Error(),
				},
			})
		}
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Queries.UnbanUser(context.Background(), user.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "user unbanned successfully",
		},
	})
}

func GetTeams(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	cursor := c.QueryParam("cursor")

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"error": "failed to convert to integer",
			},
		})
	}

	var cursorUUID uuid.UUID
	if cursor == "" {
		cursorUUID = uuid.Nil // uuid.Nil is equivalent to "00000000-0000-0000-0000-000000000000"
	} else {
		cursorUUID, err = uuid.Parse(cursor)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid UUID for cursor",
			})
		}
	}
	teams, err := utils.Queries.GetTeams(c.Request().Context(), db.GetTeamsParams{
		Limit: int32(limit),
		ID:    cursorUUID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to fetch teams",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "Teams fetched successfully",
			"teams":   teams,
		},
	})
}

func GetTeamById(c echo.Context) error {
	teamIdParam := c.Param("id")
	teamId, err := uuid.Parse(teamIdParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	team, err := utils.Queries.GetTeamById(c.Request().Context(), teamId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "Team fetched successfully",
			"team":    team,
		},
	})
}

func GetTeamLeader(c echo.Context) error {
	teamIdParam := c.Param("id")
	teamId, err := uuid.Parse(teamIdParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	nullUUID := uuid.NullUUID{
		UUID:  teamId,
		Valid: true,
	}

	user, err := utils.Queries.GetTeamLeader(c.Request().Context(), nullUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "Team leader fetched successfully",
			"user":    user,
		},
	})
}

func CreatePanel(c echo.Context) error {
	ctx := c.Request().Context()
	panel := new(models.CreatePanel)
	if err := c.Bind(panel); err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Invalid request body",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Validate.Struct(panel); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(panel.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to hash password",
				"error":   err.Error(),
			},
		})
	}

	panelDb := db.CreateUserParams{
		FirstName:     panel.FirstName,
		LastName:      panel.LastName,
		Email:         panel.Email,
		VitEmail:      panel.VitEmail,
		RegNo:         panel.RegNo,
		Password:      string(hashedPassword),
		PhoneNo:       panel.PhoneNo,
		Role:          "panel",
		IsLeader:      true,
		IsVerified:    true,
		IsBanned:      false,
		Gender:        panel.Gender,
		HostelBlock:   panel.HostelBlock,
		RoomNo:        int32(panel.RoomNumber),
		GithubProfile: panel.GithubProfile,
	}
	panelDb.ID, _ = uuid.NewV7()
	panelDb.TeamID = uuid.NullUUID{
		Valid: false,
	}

	err = utils.Queries.CreateUser(ctx, panelDb)
	if err != nil {
		logger.Errorf(logger.DatabaseError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to create user",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data:   "Panel created successfully",
	})
}

func GetAllTeamMembers(c echo.Context) error {
	teamIdParam := c.Param("id")
	teamId, err := uuid.Parse(teamIdParam)
	ctx := c.Request().Context()

	nullUUID := uuid.NullUUID{
		UUID:  teamId,
		Valid: true,
	}

	team_members, err := utils.Queries.GetTeamMembers(ctx, nullUUID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   "Cannot get Members of the team",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   team_members,
	})
}
