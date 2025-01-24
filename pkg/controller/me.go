package controller

import (
	"net/http"
	"strings"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
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

	if !user.TeamID.Valid {
		userData := models.ResponseData{
			User: models.UserData{
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				Email:         user.Email,
				RegNo:         getSafeString(user.RegNo),
				PhoneNo:       user.PhoneNo,
				Gender:        user.Gender,
				VitEmail:      getSafeString(user.VitEmail),
				HostelBlock:   user.HostelBlock,
				RoomNo:        int(user.RoomNo),
				GithubProfile: user.GithubProfile,
				IsLeader:      user.IsLeader,
			},
		}

		return c.JSON(http.StatusOK, &models.Response{
			Status:  "success",
			Message: "User details fetched successfully",
			Data:    userData,
		})
	}

	teamData, err := utils.Queries.GetUserAndTeamDetails(ctx, user.TeamID.UUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "Failed to fetch user details",
		})
	}

	marshallData := Marshall(teamData, user.ID)
	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "User details fetched successfully",
		Data:    marshallData,
	})
}

func getSafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Marshall(data []db.GetUserAndTeamDetailsRow, userID uuid.UUID) models.ResponseData {
	var response models.ResponseData

	if len(data) == 0 {
		return response
	}

	entry := data[0]
	response.User = models.UserData{
		FirstName:     entry.FirstName,
		LastName:      entry.LastName,
		Email:         entry.Email,
		RegNo:         getSafeString(entry.RegNo),
		PhoneNo:       entry.PhoneNo,
		Gender:        entry.Gender,
		VitEmail:      getSafeString(entry.VitEmail),
		HostelBlock:   entry.HostelBlock,
		RoomNo:        int(entry.RoomNo),
		GithubProfile: entry.GithubProfile,
		IsLeader:      entry.IsLeader,
	}

	if entry.Name != "" {
		response.Team = models.TeamData{
			Name:           entry.Name,
			NumberOfPeople: int(entry.NumberOfPeople),
			RoundQualified: int(entry.RoundQualified.Int32),
			Code:           entry.Code,
			Members:        []models.TeamMember{},
		}

		for _, member := range data {
			if member.ID != userID {
				response.Team.Members = append(response.Team.Members, models.TeamMember{
					FirstName:     member.FirstName,
					LastName:      member.LastName,
					Email:         member.Email,
					PhoneNo:       member.PhoneNo.String,
					GithubProfile: member.GithubProfile,
					IsLeader:      member.IsLeader,
				})
			}
		}
	}

	return response
}

func UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
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
		Email:     req.Email,
		PhoneNo: pgtype.Text{
			String: req.PhoneNo,
		},
		Gender:        req.Gender,
		RegNo:         &req.RegNo,
		VitEmail:      &req.VitEmail,
		HostelBlock:   req.HostelBlock,
		RoomNo:        int32(req.RoomNumber),
		GithubProfile: req.GithubProfile,
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
		Status:  "success",
		Message: "User updated successfully",
		Data:    updatedUser,
	})
}
