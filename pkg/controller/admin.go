package controller

import (
	"context"
	"database/sql"
	"errors"
	"log"
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

	var nextCursor uuid.NullUUID

	for _, user := range users {
		nextCursor = uuid.NullUUID{UUID: user.ID}
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Users fetched successfully",
		Data: map[string]interface{}{
			"users":       users,
			"next_cursor": nextCursor.UUID.String(),
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

	var nextCursor uuid.NullUUID

	for _, team := range teams {
		nextCursor = uuid.NullUUID{UUID: team.ID}
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Teams fetched successfully",
		Data: map[string]interface{}{
			"teams":       teams,
			"next_cursor": nextCursor.UUID.String(),
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

func GetLeaderBoard(c echo.Context) error {
	ctx := c.Request().Context()

	limitParam := c.QueryParam("limit")
	cursorParam := c.QueryParam("cursor")
	nameParam := c.QueryParam("name")

	limit := 10
	var cursor uuid.NullUUID
	var err error

	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if cursorParam != "" {
		parsedCursor, err := uuid.Parse(cursorParam)
		if err == nil {
			cursor = uuid.NullUUID{UUID: parsedCursor}
		}
	}

	rows, err := utils.Queries.GetLeaderboardWithPagination(ctx, db.GetLeaderboardWithPaginationParams{
		Column1: cursor.UUID,
		Limit:   int32(limit),
		Column2: nameParam,
	})
	if err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, &models.Response{
			Status:  "fail",
			Message: "Some error occurred",
			Data: map[string]string{
				"error": err.Error(),
			},
		})
	}

	leaderboardMap := make(map[uuid.UUID]*models.TeamLeaderboard)
	var nextCursor uuid.NullUUID

	for _, row := range rows {
		if _, exists := leaderboardMap[row.TeamID]; !exists {
			leaderboardMap[row.TeamID] = &models.TeamLeaderboard{
				TeamID:       row.TeamID,
				TeamName:     row.Name,
				Rounds:       []models.Round{},
				OverallTotal: int(row.OverallTotal),
			}
		}

		leaderboardMap[row.TeamID].Rounds = append(leaderboardMap[row.TeamID].Rounds, models.Round{
			Round:          int(row.Round),
			Design:         int(row.Design),
			Implementation: int(row.Implementation),
			Presentation:   int(row.Presentation),
			Innovation:     int(row.Innovation),
			Teamwork:       int(row.Teamwork),
			RoundTotal:     int(row.RoundTotal),
		})
		nextCursor = uuid.NullUUID{UUID: row.TeamID}
	}

	leaderBoard := make([]models.TeamLeaderboard, 0, len(leaderboardMap))
	for _, team := range leaderboardMap {
		leaderBoard = append(leaderBoard, *team)
	}

	response := map[string]interface{}{
		"leaderboard": leaderBoard,
		"next_cursor": nextCursor.UUID.String(),
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "Leaderboard fetched successfully",
		Data:    response,
	})

}

func GetAllIdeas(c echo.Context) error {
	ctx := c.Request().Context()
	limitParam := c.QueryParam("limit")
	cursor := c.QueryParam("cursor")

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

	ideas, err := utils.Queries.GetAllIdeas(ctx, db.GetAllIdeasParams{
		Limit: int32(limit),
		ID:    cursorUUID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	var nextCursor uuid.NullUUID

	for _, idea := range ideas {
		nextCursor = uuid.NullUUID{UUID: idea.TeamID, Valid: true}
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "ideas fetched successfully",
		Data: map[string]interface{}{
			"ideas":       ideas,
			"next_cursor": nextCursor,
		},
	})
}

func GetIdeasByTrack(c echo.Context) error {
	ctx := c.Request().Context()

	TrackParam := c.QueryParam("track")
	TitleParam := c.QueryParam("title")

	log.Println(TrackParam)

	if TrackParam == "" {
		TrackParam = "7"
	}

	log.Println(TrackParam)

	TrackParamInt, err := strconv.Atoi(TrackParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
	}

	if TrackParamInt < 1 || TrackParamInt > 6 {
		TrackParamInt = 7
	}

	payload := struct {
		Track int    `json:"track"`
		Title string `json:"title"`
	}{
		Track: TrackParamInt,
		Title: TitleParam,
	}

	limitParam := c.QueryParam("limit")
	cursor := c.QueryParam("cursor")

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "success",
			Message: err.Error(),
		})
	}

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

	var idea []db.Idea

	tracks := map[int]string{
		1: "AI&ML",
		2: "Finance and Fintech",
		3: "Healthcare and Education",
		4: "Digital Security",
		5: "Environment and Sustainability",
		6: "Open Innovation",
		7: "General",
	}
	log.Println("--------------")
	log.Println(payload.Title)
	log.Println("--------------")

	switch payload.Track {
	case 1:
		trackname := tracks[1]
		idea, err = utils.Queries.GetIdeasByTrack(ctx, db.GetIdeasByTrackParams{
			Column1: &payload.Title,
			Column2: &trackname,
			ID:      cursorUUID,
			Limit:   int32(limit),
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	case 2:
		trackname := tracks[2]
		idea, err = utils.Queries.GetIdeasByTrack(ctx, db.GetIdeasByTrackParams{
			Column1: &payload.Title,
			Column2: &trackname,
			ID:      cursorUUID,
			Limit:   int32(limit),
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	case 3:
		trackname := tracks[3]
		idea, err = utils.Queries.GetIdeasByTrack(ctx, db.GetIdeasByTrackParams{
			Column1: &payload.Title,
			Column2: &trackname,
			ID:      cursorUUID,
			Limit:   int32(limit),
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	case 4:
		trackname := tracks[4]
		idea, err = utils.Queries.GetIdeasByTrack(ctx, db.GetIdeasByTrackParams{
			Column1: &payload.Title,
			Column2: &trackname,
			ID:      cursorUUID,
			Limit:   int32(limit),
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	case 5:
		trackname := tracks[5]
		idea, err = utils.Queries.GetIdeasByTrack(ctx, db.GetIdeasByTrackParams{
			Column1: &payload.Title,
			Column2: &trackname,
			ID:      cursorUUID,
			Limit:   int32(limit),
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	case 6:
		trackname := tracks[6]
		idea, err = utils.Queries.GetIdeasByTrack(ctx, db.GetIdeasByTrackParams{
			Column1: &payload.Title,
			Column2: &trackname,
			ID:      cursorUUID,
			Limit:   int32(limit),
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	case 7:
		trackname := tracks[7]
		idea, err = utils.Queries.GetIdeasByTrack(ctx, db.GetIdeasByTrackParams{
			Column1: &payload.Title,
			Column2: &trackname,
			ID:      cursorUUID,
			Limit:   int32(limit),
		})
		log.Println("-----------")
		log.Println(cursorUUID)
		log.Println(payload.Title)
		log.Println("-----------")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Status:  "fail",
				Message: err.Error(),
			})
		}
	default:
		return c.JSON(http.StatusBadRequest, &models.Response{
			Status:  "fail",
			Message: "",
		})
	}

	var nextCursor uuid.NullUUID

	for _, ide := range idea {
		nextCursor = uuid.NullUUID{UUID: ide.TeamID}
	}

	return c.JSON(http.StatusOK, &models.Response{
		Status:  "success",
		Message: "ideas fetched successfully",
		Data: map[string]interface{}{
			"ideas":       idea,
			"next_cursor": nextCursor,
		},
	})

}
