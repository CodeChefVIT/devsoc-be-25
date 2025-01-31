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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(c echo.Context) error {
	ctx := c.Request().Context()
	limitParam := c.QueryParam("limit")
	cursor := c.QueryParam("cursor")
	name := c.QueryParam("name")

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	var cursorUUID uuid.UUID
	if cursor == "" {
		cursorUUID = uuid.Nil
	} else {
		cursorUUID, err = uuid.Parse(cursor)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid UUID for cursor",
			})
		}
	}

	users, err := utils.Queries.GetAllUsers(ctx, db.GetAllUsersParams{
		Limit:   int32(limit),
		ID:      cursorUUID,
		Column1: &name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Users fetched successfully",
		Data: map[string]interface{}{
			"users": users,
		},
	})
}

func GetUsersByEmail(c echo.Context) error {
	email := c.Param("email")
	user, err := utils.Queries.GetUserByEmail(c.Request().Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "User fetched successfully",
		Data: map[string]interface{}{
			"user": user,
		},
	})
}

func BanUser(c echo.Context) error {
	var payload models.BanUserReq

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	if err := utils.Queries.BanUser(c.Request().Context(), payload.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "user bannned successfully",
		Data:    map[string]string{},
	})
}

func UnbanUser(c echo.Context) error {
	var payload models.BanUserReq

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	if err := utils.Queries.UnbanUser(context.Background(), payload.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "user unbanned successfully",
		Data:    map[string]string{},
	})
}

func GetTeams(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	cursor := c.QueryParam("cursor")
	name := c.QueryParam("name")

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	var cursorUUID uuid.UUID
	if cursor == "" {
		cursorUUID = uuid.Nil
	} else {
		cursorUUID, err = uuid.Parse(cursor)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid UUID for cursor",
			})
		}
	}
	teams, err := utils.Queries.GetTeams(c.Request().Context(), db.GetTeamsParams{
		Limit:   int32(limit),
		ID:      cursorUUID,
		Column1: &name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Teams fetched successfully",
		Data: map[string]interface{}{
			"teams": teams,
		},
	})
}

func GetTeamById(c echo.Context) error {
	teamIdParam := c.Param("id")
	teamId, err := uuid.Parse(teamIdParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	team, err := utils.Queries.GetTeamById(c.Request().Context(), teamId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Team fetched successfully",
		Data: map[string]interface{}{
			"team": team,
		},
	})
}

func GetTeamLeader(c echo.Context) error {
	teamIdParam := c.Param("id")
	teamId, err := uuid.Parse(teamIdParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	nullUUID := uuid.NullUUID{
		UUID:  teamId,
		Valid: true,
	}

	user, err := utils.Queries.GetTeamLeader(c.Request().Context(), nullUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Team leader fetched successfully",
		Data: map[string]interface{}{
			"user": user,
		},
	})
}

func CreatePanel(c echo.Context) error {
	ctx := c.Request().Context()
	panel := new(models.CreatePanel)
	if err := c.Bind(panel); err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: err.Error(),
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
			Status:  "fail",
			Message: err.Error(),
		})
	}

	panelDb := db.CreateUserParams{
		FirstName: panel.FirstName,
		LastName:  panel.LastName,
		Email:     panel.Email,
		RegNo:     &panel.RegNo,
		Password:  string(hashedPassword),
		PhoneNo: pgtype.Text{
			String: panel.PhoneNo,
		},
		Role:       "panel",
		IsLeader:   true,
		IsVerified: true,
		IsBanned:   false,
		Gender:     panel.Gender,
	}
	panelDb.ID, _ = uuid.NewV7()
	panelDb.TeamID = uuid.NullUUID{
		Valid: false,
	}

	err = utils.Queries.CreateUser(ctx, panelDb)
	if err != nil {
		logger.Errorf(logger.DatabaseError, err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Panel Created Successfully",
		Data:    map[string]interface{}{},
	})
}

func GetAllTeamMembers(c echo.Context) error {
	teamIdParam := c.Param("id")
	teamId, _ := uuid.Parse(teamIdParam)
	ctx := c.Request().Context()

	nullUUID := uuid.NullUUID{
		UUID:  teamId,
		Valid: true,
	}

	team_members, err := utils.Queries.GetTeamMembers(ctx, nullUUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Team fetched successfully",
		Data: map[string]interface{}{
			"team": team_members,
		},
	})
}

// Ban Team
func BanTeam(c echo.Context) error {
	var payload models.UnBanTeam

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	ctx := c.Request().Context()

	team, err := utils.Queries.GetTeamById(ctx, payload.TeamId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, &models.Response{
				Status: "fail",
				Data:   "Team Does not exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data:   err,
		})
	}

	if err := utils.Queries.BanTeam(ctx, team.ID); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Team Banned Successfully",
	})
}

// UnBan Team
func UnBanTeam(c echo.Context) error {
	var payload models.UnBanTeam

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	ctx := c.Request().Context()

	team, err := utils.Queries.GetTeamById(ctx, payload.TeamId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	if err := utils.Queries.UnBanTeam(ctx, team.ID); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "fail",
			Message: "Failed to unban Team",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Team UnBanned Successfully",
	})
}

func UpdateTeamRounds(c echo.Context) error {
	var payload models.TeamRoundQualified

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	ctx := c.Request().Context()

	team, err := utils.Queries.GetTeamById(ctx, payload.TeamId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	if err := utils.Queries.UpdateTeamRound(ctx, db.UpdateTeamRoundParams{
		RoundQualified: (pgtype.Int4{
			Int32: int32(payload.RoundQualified),
			Valid: true,
		}),
		ID: team.ID,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Rounds qualified by team Updated",
	})
}
