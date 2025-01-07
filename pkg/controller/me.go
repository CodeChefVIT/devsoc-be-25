package controller

import (
	"context"
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/labstack/echo/v4"
)

type UpdateUserInput struct {
	Name     string `json:"name"`
	PhoneNo  string `json:"phone_no"`
}

func GetDetails(c echo.Context) error {
	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusForbidden, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User not found",
			},
		})
	}

	var err error
	userData, err := utils.Queries.GetUser(context.Background(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to fetch user details",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"name":     userData.Name,
			"email":    userData.Email,
			"phone_no": userData.PhoneNo,
			"college":  userData.College,
			"reg_no":   userData.RegNo,
		},
	})
}

func UpdateUser(c echo.Context) error {
	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusForbidden, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User not found",
			},
		})
	}

	var input UpdateUserInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Invalid request",
				"error":   err.Error(),
			},
		})
	}

	var updateParams db.UpdateUserParams
	updateParams.Name = input.Name
	updateParams.PhoneNo = input.PhoneNo
	updateParams.ID = user.ID

	var err error
	err = utils.Queries.UpdateUser(context.Background(), updateParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to update user",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "User updated successfully",
		},
	})
}