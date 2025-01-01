package controller

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/dto"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func CreatePanel(c echo.Context) error {
	ctx := c.Request().Context()
	user := new(dto.UserDto)
	if err := c.Bind(user); err != nil {
		logger.Errorf(logger.ParsingError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	panel := db.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		RegNo:    user.RegNo,
		Password: string(hashedPassword),
		PhoneNo:  user.PhoneNo,
	}
	panel.ID, _ = uuid.NewV7()
	panel.Role = "panel"
	panel.IsVerified = true
	panel.IsLeader = true
	panel.IsVitian = true
	panel.IsBanned = false
	panel.College = "VIT"
	panel.TeamID = uuid.NullUUID{
		Valid: false,
	}

	err = utils.Queries.CreateUser(ctx, panel)
	if err != nil {
		logger.Errorf(logger.DatabaseError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
