package controller

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/labstack/echo/v4"
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
		"users":   users})
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
			"message": err.Error(),
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
			"message": err.Error(),
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
			"message": err.Error(),
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
			"message": err.Error(),
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
