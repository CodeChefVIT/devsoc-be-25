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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(c echo.Context) error {
	ctx := c.Request().Context()
	limitParam := c.QueryParam("limit")
	cursor := c.QueryParam("cursor")
	name := c.QueryParam("name")
	gender := c.QueryParam("gender")

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
		Column4: gender,
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

func GetUsersByGender(c echo.Context) error {
	ctx := c.Request().Context()
	gender := c.Param("gender")

	if gender != "M" && gender != "F" && gender != "O" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "Gender should be M or F",
		})
	}

	users, err := utils.Queries.GetUsersByGender(ctx, gender)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "users fetched successfully",
		Data:    users,
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

	team, err := utils.Queries.GetTeamByTeamId(c.Request().Context(), teamId)
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

	users, err := utils.Queries.GetUsersByTeamId(c.Request().Context(), nullUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	submission, err := utils.Queries.GetSubmissionByTeamID(c.Request().Context(), teamId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			submission = db.Submission{}
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	}

	idea, err := utils.Queries.GetIdeaByTeamID(c.Request().Context(), teamId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			idea = db.GetIdeaByTeamIDRow{}
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	}

	score, err := utils.Queries.GetTeamScores(c.Request().Context(), teamId)
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
			"team":         team,
			"team_members": users,
			"submission":   submission,
			"idea":         idea,
			"score":        score,
		},
	})
}

func GetTeamsByTrack(c echo.Context) error {
	ctx := c.Request().Context()
	track := c.Param("track")

	teams, err := utils.Queries.GetTeamByTrack(ctx, track)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "teams fetched successfully",
		Data:    teams,
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

func GetAllIdeas(c echo.Context) error {
	ctx := c.Request().Context()

	ideas, err := utils.Queries.GetAllIdeas(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "ideas fetched successfully",
		Data:    ideas,
	})
}

func GetIdeasByTrack(c echo.Context) error {
	ctx := c.Request().Context()
	num := c.Param("track")
	trackNum, err := strconv.Atoi(num)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: "couldn't convert string to int",
		})
	}

	var idea []db.Idea

	if trackNum == 1 {
		idea, err = utils.Queries.GetIdeasByTrack(ctx, "Media and Entertainment")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	} else if trackNum == 2 {
		idea, err = utils.Queries.GetIdeasByTrack(ctx, "Finance and Fintech")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	} else if trackNum == 3 {
		idea, err = utils.Queries.GetIdeasByTrack(ctx, "Healthcare and Education")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	} else if trackNum == 4 {
		idea, err = utils.Queries.GetIdeasByTrack(ctx, "Digital Security")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	} else if trackNum == 5 {
		idea, err = utils.Queries.GetIdeasByTrack(ctx, "Environment and Sustainability")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	} else if trackNum == 6 {
		idea, err = utils.Queries.GetIdeasByTrack(ctx, "Open Innovation")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	} else {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "give number from 1 to 6",
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Ideas fetched successfully",
		Data:    idea,
	})

}
