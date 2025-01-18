package controller

import (
	"encoding/csv"
	"net/http"
	"os"
	"strconv"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/labstack/echo/v4"
)

func ExportUsers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := utils.Queries.ExportAllUsers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to fetch users"},
		})
	}

	file, err := os.Create("users.csv")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to create CSV file"},
		})
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	headers := []string{"ID", "FirstName", "LastName", "Email", "PhoneNo", "Gender", "RegNo", "TeamID", "VitEmail", "Hostel", "RoomNo", "GitHub", "Role", "IsLeader", "IsVerified", "IsBanned", "IsProfComplete"}
	if err := csvWriter.Write(headers); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to write CSV headers"},
		})
	}

	for _, user := range users {
		record := []string{
			user.ID.String(),
			user.FirstName,
			user.LastName,
			user.Email,
			user.PhoneNo.String,
			user.Gender,
			*user.RegNo,
			user.TeamID.UUID.String(),
			*user.VitEmail,
			user.HostelBlock,
			strconv.Itoa(int(user.RoomNo)),
			user.GithubProfile,
			user.Role,
			strconv.FormatBool(user.IsLeader),
			strconv.FormatBool(user.IsVerified),
			strconv.FormatBool(user.IsBanned),
			strconv.FormatBool(user.IsProfileComplete),
		}

		if err := csvWriter.Write(record); err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status: "fail",
				Data:   map[string]string{"error": "Failed to write CSV record"},
			})
		}
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to flush CSV writer"},
		})
	}

	return c.Attachment("users.csv", "users.csv")
}

func ExportTeams(c echo.Context) error {
	ctx := c.Request().Context()

	teams, err := utils.Queries.ExportAllTeams(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to fetch teams"},
		})
	}

	file, err := os.Create("teams.csv")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to create CSV file"},
		})
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	headers := []string{"ID", "TeamName", "TeamCode", "NumberOfPeople", "RoundQualified",
		"IdeaId", "IdeaTitle", "IdeaDescription", "IdeaTrack",
		"SubmissionId", "SubmissionTitle", "SubmissionDescription", "SubmissionTrack", "GitHubLink", "FigmaLink", "OtherLink",
		"ScoreId", "DesignScore", "ImplementationScore", "PresentationScore", "ScoreRound"}
	if err := csvWriter.Write(headers); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to write CSV headers"},
		})
	}

	for _, team := range teams {
		submission, _ := utils.Queries.GetSubmissionByTeamID(ctx, team.ID)
		scores, _ := utils.Queries.GetTeamScores(ctx, team.ID)
		idea, _ := utils.Queries.GetIdeaByTeamID(ctx, team.ID)

		if len(scores) != 0 {
			for _, score := range scores {
				record := []string{
					team.ID.String(),
					team.Name,
					team.Code,
					strconv.Itoa(int(team.NumberOfPeople)),
					strconv.Itoa(int(team.RoundQualified.Int32)),
					idea.ID.String(),
					idea.Title,
					idea.Description,
					idea.Track,
					submission.ID.String(),
					submission.Title,
					submission.Description,
					submission.Track,
					submission.GithubLink,
					submission.FigmaLink,
					submission.OtherLink,
					score.ID.String(),
					strconv.Itoa(int(score.Design)),
					strconv.Itoa(int(score.Implementation)),
					strconv.Itoa(int(score.Presentation)),
					strconv.Itoa(int(score.Round)),
				}

				if err := csvWriter.Write(record); err != nil {
					return c.JSON(http.StatusInternalServerError, &models.Response{
						Status: "fail",
						Data:   map[string]string{"error": "Failed to write CSV record"},
					})
				}
			}
		} else {
			record := []string{
				team.ID.String(),
				team.Name,
				team.Code,
				strconv.Itoa(int(team.NumberOfPeople)),
				strconv.Itoa(int(team.RoundQualified.Int32)),
				idea.ID.String(),
				idea.Title,
				idea.Description,
				idea.Track,
				submission.ID.String(),
				submission.Title,
				submission.Description,
				submission.Track,
				submission.GithubLink,
				submission.FigmaLink,
				submission.OtherLink,
				"NA",
				"NA",
				"NA",
				"NA",
				"NA",
			}

			if err := csvWriter.Write(record); err != nil {
				return c.JSON(http.StatusInternalServerError, &models.Response{
					Status: "fail",
					Data:   map[string]string{"error": "Failed to write CSV record"},
				})
			}
		}
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data:   map[string]string{"error": "Failed to flush CSV writer"},
		})
	}

	return c.Attachment("teams.csv", "teams.csv")
}
