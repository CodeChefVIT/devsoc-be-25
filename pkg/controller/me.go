package controller

import (
	"net/http"
	"strings"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserData struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	RegNo         string `json:"reg_no"`
	PhoneNo       string `json:"phone_no"`
	Gender        string `json:"gender"`
	VitEmail      string `json:"vit_email"`
	HostelBlock   string `json:"hostel_block"`
	RoomNo        int    `json:"room_no"`
	GithubProfile string `json:"github_profile"`
	IsLeader      bool   `json:"is_leader"`
}

type TeamMember struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	PhoneNo       string `json:"phone_no"`
	GithubProfile string `json:"github_profile"`
	IsLeader      bool   `json:"is_leader"`
}

type TeamData struct {
	Name           string       `json:"team_name"`
	NumberOfPeople int          `json:"number_of_people"`
	RoundQualified int          `json:"round_qualified"`
	Code           string       `json:"code"`
	Members        []TeamMember `json:"members"`
}

type ResponseData struct {
	Message string   `json:"message"`
	User    UserData `json:"user"`
	Team    TeamData `json:"team"`
}

func Marshall(data []db.GetUserAndTeamDetailsRow, userID uuid.UUID) ResponseData {
	var response ResponseData

	for _, entry := range data {
		if entry.ID == userID {
			response.User = UserData{
				FirstName:     entry.FirstName,
				LastName:      entry.LastName,
				Email:         entry.Email,
				RegNo:         entry.RegNo,
				PhoneNo:       entry.PhoneNo,
				Gender:        entry.Gender,
				VitEmail:      entry.VitEmail,
				HostelBlock:   entry.HostelBlock,
				RoomNo:        int(entry.RoomNo),
				GithubProfile: entry.GithubProfile,
				IsLeader:      entry.IsLeader,
			}

			response.Team = TeamData{
				Name:           entry.Name,
				NumberOfPeople: int(entry.NumberOfPeople),
				RoundQualified: int(entry.RoundQualified.Int32),
				Code:           entry.Code,
				Members:        []TeamMember{},
			}
			break
		}
	}

	for _, entry := range data {
		if entry.ID != userID {
			member := TeamMember{
				FirstName:     entry.FirstName,
				LastName:      entry.LastName,
				Email:         entry.Email,
				PhoneNo:       entry.PhoneNo,
				GithubProfile: entry.GithubProfile,
			}
			response.Team.Members = append(response.Team.Members, member)
		}
	}

	return response
}

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

	if user.TeamID.Valid == false {
		return c.JSON(http.StatusForbidden, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Forbidden",
				"error":   "User is not part of any team",
			},
		})
	}

	userData, err := utils.Queries.GetUserAndTeamDetails(ctx, user.TeamID.UUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to fetch user details",
				"error":   err.Error(),
			},
		})
	}

	marshallData := Marshall(userData, user.ID)
	marshallData.Message = "User details fetched successfully"
	return c.JSON(http.StatusOK, &models.Response{
		Status: "success",
		Data: marshallData,
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

	if err := utils.Validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status: "fail",
			Data:   utils.FormatValidationErrors(err),
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
