package controller

import (
	"net/http"
	"strings"

	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func GetDetails(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status:  "fail",
			Message: "User not found",
		})
	}

	hostelBlock := ""
	roomNo := ""

	if user.HostelBlock != nil {
		hostelBlock = *user.HostelBlock
	}

	if user.RoomNo != nil {
		roomNo = *user.RoomNo
	}

	res := models.ResponseData{
		User: models.UserData{
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Email:         user.Email,
			RegNo:         getSafeString(user.RegNo),
			PhoneNo:       user.PhoneNo,
			Gender:        user.Gender,
			GithubProfile: *user.GithubProfile,
			IsLeader:      user.IsLeader,
			HostelBlock:   hostelBlock,
			RoomNo:        roomNo,
		},
	}

	if !user.TeamID.Valid {
		return c.JSON(http.StatusOK, &models.Response{
			Status:  "success",
			Message: "User details fetched successfully",
			Data:    res,
		})
	}

	teamData, err := utils.Queries.InfoQuery(ctx, user.TeamID.UUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to fetch details",
		})
	}

	res.Team = models.TeamData{
		Name:           teamData[0].Name,
		NumberOfPeople: len(teamData),
		RoundQualified: int(teamData[0].RoundQualified.Int32),
		Code:           teamData[0].Code,
	}
	for _, v := range teamData {
		res.Team.Members = append(res.Team.Members, models.TeamMember{
			FirstName:     v.FirstName,
			LastName:      v.LastName,
			GithubProfile: *v.GithubProfile,
			IsLeader:      v.IsLeader,
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "User details fetched successfully",
		Data:    res,
	})
}

func getSafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		logger.Warnf(err.Error())
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Invalid request body",
		})
	}

	if err := utils.Validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Validation failed",
			Data:    utils.FormatValidationErrors(err),
		})
	}

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status:  "fail",
			Message: "User not found",
		})
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	if req.FirstName == "" || req.LastName == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "First name and last name cannot be empty",
		})
	}

	if req.Gender != "M" && req.Gender != "F" && req.Gender != "O" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Gender must be M, F or O",
		})
	}

	err := utils.Queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:        user.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		PhoneNo: pgtype.Text{
			String: req.PhoneNo,
			Valid:  true,
		},
		Gender:        req.Gender,
		RegNo:         &req.RegNo,
		GithubProfile: &req.GithubProfile,
		HostelBlock:   &req.HostelBlock,
		RoomNo:        &req.RoomNumber,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to update user",
		})
	}

	updatedUser := map[string]interface{}{
		"first_name":     req.FirstName,
		"last_name":      req.LastName,
		"phone_no":       req.PhoneNo,
		"gender":         req.Gender,
		"reg_no":         req.RegNo,
		"github_profile": req.GithubProfile,
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "User updated successfully",
		Data:    updatedUser,
	})
}
