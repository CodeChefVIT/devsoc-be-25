package controller

import (
	"net/http"
	"strings"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/labstack/echo/v4"
)

func GetDetails(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User not found",
			},
		})
	}

	userData, err := utils.Queries.GetUser(ctx, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to fetch user details",
				"error":   err.Error(),
			},
		})
	}

	fetchedUser := map[string]interface{}{
		"first_name":     userData.FirstName,
		"last_name":      userData.LastName,
		"email":          userData.Email,
		"phone_no":       userData.PhoneNo,
		"gender":         userData.Gender,
		"reg_no":         userData.RegNo,
		"vit_email":      userData.VitEmail,
		"hostel_block":   userData.HostelBlock,
		"room_no":        int32(userData.RoomNo),
		"github_profile": userData.GithubProfile,
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "User fetched successfully",
			"user":    fetchedUser,
		},
	})
}

type UpdateUserRequest struct {
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	PhoneNo       string `json:"phone_no" validate:"required,len=10"`
	Gender        string `json:"gender" validate:"required,len=1"`
	RegNo         string `json:"reg_no" validate:"required"`
	VitEmail      string `json:"vit_email" validate:"required,email,endswith=@vitstudent.ac.in"`
	HostelBlock   string `json:"hostel_block" validate:"required"`
	RoomNumber    int    `json:"room_no" validate:"required"`
	GithubProfile string `json:"github_profile" validate:"required,url"`
	Password      string `json:"password" validate:"required"`
}

func UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Invalid request body",
				"error":   err.Error(),
			},
		})
	}

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"error": "User not found",
			},
		})
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	if req.FirstName == "" || req.LastName == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "First name and last name cannot be empty"},
		})
	}

	if req.Gender != "M" && req.Gender != "F" && req.Gender != "O" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Gender must be M, F or O"},
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
		})
	}

	err := utils.Queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:            user.ID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		PhoneNo:       req.PhoneNo,
		Gender:        req.Gender,
		RegNo:         req.RegNo,
		VitEmail:      req.VitEmail,
		HostelBlock:   req.HostelBlock,
		RoomNo:        int32(req.RoomNumber),
		GithubProfile: req.GithubProfile,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to update user",
				"error":   err.Error(),
			},
		})
	}

	updatedUser := map[string]interface{}{
		"first_name":     req.FirstName,
		"last_name":      req.LastName,
		"email":          req.Email,
		"phone_no":       req.PhoneNo,
		"gender":         req.Gender,
		"reg_no":         req.RegNo,
		"vit_email":      req.VitEmail,
		"hostel_block":   req.HostelBlock,
		"room_no":        int32(req.RoomNumber),
		"github_profile": req.GithubProfile,
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: map[string]interface{}{
			"message": "User updated successfully",
			"user":    updatedUser,
		},
	})
}
