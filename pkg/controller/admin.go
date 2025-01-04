package controller

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(c echo.Context) error {
	ctx := context.Background()
	users, err := utils.Queries.GetAllUsers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to fetch users",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Users fetched successfully",
		"users":   users,
	})
}

func GetAllVitians(c echo.Context) error {
	users, err := utils.Queries.GetAllVitians(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to fetch users",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Users fetched successfully",
		"users":   users,
	})
}

func GetUsersByEmail(c echo.Context) error {
	email := c.Param("email")
	user, err := utils.Queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "User not found",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Some error occured",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User fetched successfully",
		"user":    user,
	})
}

func BanUser(c echo.Context) error {
	var payload models.BanUserReq

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Improper request",
			"error":   err.Error(),
		})
	}

	user, err := utils.Queries.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "User does not exist",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "some error occured",
			"error":   err.Error(),
		})
	}

	if err := utils.Queries.BanUser(context.Background(), user.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Some error occured",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user banned successfully",
	})

}

func UnbanUser(c echo.Context) error {
	var payload models.BanUserReq

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Improper request",
			"error":   err.Error(),
		})
	}

	user, err := utils.Queries.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "User does not exist",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "some error occured",
			"error":   err.Error(),
		})
	}

	if err := utils.Queries.UnbanUser(context.Background(), user.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Some error occured",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user unbanned successfully",
	})
}

func GetTeams(c echo.Context) error {
	teams, err := utils.Queries.GetTeams(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to fetch teams",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Teams fetched successfully",
		"teams":   teams,
	})
}

func GetTeamById(c echo.Context) error {
	teamIdParam := c.Param("id")
	teamId, err := uuid.Parse(teamIdParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "some error occured",
			"error":   err.Error(),
		})
	}

	team, err := utils.Queries.GetTeamById(context.Background(), teamId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "some error occured",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Team fetched successfully",
		"team":    team,
	})
}

func GetTeamLeader(c echo.Context) error {
	teamIdParam := c.Param("id")
	teamId, err := uuid.Parse(teamIdParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "some error occured",
			"error":   err.Error(),
		})
	}

	nullUUID := uuid.NullUUID{
		UUID:  teamId,
		Valid: true,
	}

	user, err := utils.Queries.GetTeamLeader(context.Background(), nullUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "some error occured",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Team leader fetched successfully",
		"user":    user,
	})
}

func CreatePanel(c echo.Context) error {
	ctx := c.Request().Context()
	panel := new(models.CreatePanel)
	if err := c.Bind(panel); err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := utils.Validate.Struct(panel); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "400",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(panel.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	panelDb := db.CreateUserParams{
		Name:     panel.Name,
		Email:    panel.Email,
		RegNo:    panel.RegNo,
		Password: string(hashedPassword),
		PhoneNo:  panel.PhoneNo,
	}
	panelDb.ID, _ = uuid.NewV7()
	panelDb.Role = "panel"
	panelDb.IsVerified = true
	panelDb.IsLeader = true
	panelDb.IsVitian = true
	panelDb.IsBanned = false
	panelDb.College = "VIT"
	panelDb.TeamID = uuid.NullUUID{
		Valid: false,
	}

	err = utils.Queries.CreateUser(ctx, panelDb)
	if err != nil {
		logger.Errorf(logger.DatabaseError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	return utils.WriteJSON(c, http.StatusOK, "Panel Created")
}
