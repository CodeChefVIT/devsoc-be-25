package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/db"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/models"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
)

func GetTeamId(c echo.Context) error {
	ctx := c.Request().Context()
	teamCode := c.Param("teamcode")

	teamId, err := utils.Queries.GetTeamIDByCode(ctx, teamCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to fetch team ID",
				"error":   err.Error(),
			},
		})
	}

	response := map[string]interface{}{
		"teamId": teamId,
	}
	return c.JSON(http.StatusOK, response)
}

//JOIN TEAM

func JoinTeam(c echo.Context) error {
	var payload models.JoinTeam
	if err := c.Bind(&payload); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}

	user := c.Get("user").(db.User)
	fmt.Println(user)

	if user.TeamID.Valid {
		return utils.WriteError(c, echo.ErrExpectationFailed.Code, errors.New("user already in a team"))
	}

	ctx := context.Background()

	team, err := utils.Queries.FindTeam(ctx, payload.Code)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("team code invalid"))
		}
		return utils.WriteError(c, echo.ErrInternalServerError.Code, errors.New("failed to get team"))
	}

	nullableTeamID := uuid.NullUUID{
		UUID:  team.ID,
		Valid: true,
	}

	count, err := utils.Queries.CountTeamMembers(ctx, nullableTeamID)
	if err != nil {
		return utils.WriteError(c, echo.ErrInternalServerError.Code, errors.New("failed to get team member count"))
	}

	if count >= 4 {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Cannot join Team. Team is already full"))
	}

	if err := utils.Queries.IncreaseCountTeam(ctx, nullableTeamID.UUID); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Failed to add member to team"))
	}

	if err := utils.Queries.AddUserToTeam(ctx, db.AddUserToTeamParams{
		TeamID: nullableTeamID,
		ID:     user.ID,
	}); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Cannot join team"))
	}

	return utils.WriteJSON(c, 200, "Team joined successfully")
}

//KICK MEMBER

func KickMemeber(c echo.Context) error {
	var payload models.KickMember
	ctx := context.Background()
	if err := c.Bind(&payload); err != nil {
		utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}

	user := c.Get("user").(db.User)

	if user.IsLeader != true {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Only leader can kick User"))
	}

	member, err := utils.Queries.GetUserByID(ctx, payload.UserID)
	if err != nil {
		return utils.WriteError(c, echo.ErrNotFound.Code, errors.New("user not found"))
	}

	nullableTeamID := uuid.NullUUID{
		UUID:  user.TeamID.UUID,
		Valid: true,
	}

	nullableMemberID := uuid.NullUUID{
		UUID:  member.TeamID.UUID,
		Valid: true,
	}

	count, err := utils.Queries.CountTeamMembers(ctx, nullableTeamID)
	if count <= 0 || err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Cannot leave Team. Team is already empty"))
	}

	if user.TeamID.UUID != member.TeamID.UUID {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("user is not a member of your team"))
	}

	if err := utils.Queries.RemoveUserFromTeam(ctx, db.RemoveUserFromTeamParams{
		TeamID: nullableMemberID,
		ID:     member.ID,
	}); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}

	if err := utils.Queries.DecreaseUserCountTeam(ctx, nullableTeamID.UUID); err != nil {
		utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Failed to leave team"))
	}

	return utils.WriteJSON(c, 200, "Member kicked successfully")
}

// CREATE TEAM

func CreateTeam(c echo.Context) error {
	var payload models.CreateTeam

	ctx := context.Background()
	if err := c.Bind(&payload); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}

	payload.Name = strings.TrimSpace(payload.Name)

	userInterface := c.Get("user")
	if userInterface == nil {
		return c.JSON(400, map[string]string{"error": "user not found in context"})
	}

	user, ok := userInterface.(db.User)
	if !ok {
		return c.JSON(500, map[string]string{"error": "failed to cast user"})
	}

	if user.TeamID.Valid {
		return utils.WriteError(c, echo.ErrExpectationFailed.Code, errors.New("user already in a team"))
	}

	fmt.Println(user.TeamID.UUID.String())

	params := db.CreateTeamParams{
		ID:             uuid.New(),
		Name:           payload.Name,
		Code:           uuid.NewString(),
		NumberOfPeople: 1,
		RoundQualified: pgtype.Int4{Int32: 0, Valid: true},
	}

	team, err := utils.Queries.CreateTeam(ctx, params)
	if err != nil {
		return err
	}

	err = utils.Queries.UpdateUserTeam(ctx, db.UpdateUserTeamParams{
		TeamID: uuid.NullUUID{
			UUID:  team.ID,
			Valid: true,
		},
		IsLeader: true,
		ID:       user.ID,
	})

	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code == "23505" {
				return utils.WriteError(c, 23505, errors.New("team name already exists"))
			}
		}
		return utils.WriteError(c, echo.ErrInternalServerError.Code, err)
	}

	return utils.WriteJSON(c, 200, team)
}

//LEAVE TEAM

func LeaveTeam(c echo.Context) error {
	var payload models.LeaveTeam
	ctx := context.Background()
	if err := c.Bind(&payload); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}

	user, ok := c.Get("user").(db.User)
	if !ok {
		return utils.WriteError(c, echo.ErrUnauthorized.Code, errors.New("unauthorized"))
	}

	if user.ID != payload.UserID {
		return utils.WriteError(c, echo.ErrForbidden.Code, errors.New("you can only leave your own team"))
	}

	if !user.TeamID.Valid || user.TeamID.UUID == uuid.Nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("you are not part of a team"))
	}

	if user.IsLeader {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Leader cannot leave the team with making eader someone else"))
	}

	if err := utils.Queries.RemoveUserFromTeam(ctx, db.RemoveUserFromTeamParams{
		ID: user.TeamID.UUID,
	}); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}

	nullableTeamID := uuid.NullUUID{
		UUID:  user.TeamID.UUID,
		Valid: true,
	}

	count, err := utils.Queries.CountTeamMembers(ctx, nullableTeamID)
	if count <= 0 || err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Cannot leave Team. Team is already empty"))
	}

	if err := utils.Queries.LeaveTeam(ctx, payload.UserID); err != nil {
		return err
	}

	if err := utils.Queries.DecreaseUserCountTeam(ctx, nullableTeamID.UUID); err != nil {
		utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Failed to leave team"))
	}

	return utils.WriteJSON(c, 200, "Team Left Successfully")
}

//DELETE TEAM

func DeleteTeam(c echo.Context) error {
	ctx := context.Background()

	user, ok := c.Get("user").(db.User)
	if !ok {
		return utils.WriteError(c, echo.ErrUnauthorized.Code, errors.New("unauthorized"))
	}

	if !user.IsLeader {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Only leader can delete the team"))
	}

	nullableTeamID := uuid.NullUUID{
		UUID:  user.TeamID.UUID,
		Valid: true,
	}

	if err := utils.Queries.RemoveTeamIDFromUsers(ctx, nullableTeamID); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}

	if err := utils.Queries.DeleteTeam(ctx, user.TeamID.UUID); err != nil {
		return utils.WriteError(c, echo.ErrBadGateway.Code, err)
	}

	if err := utils.Queries.UpdateLeader(ctx, db.UpdateLeaderParams{
		IsLeader: false,
		ID:       user.ID,
	}); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}
	user.TeamID = uuid.NullUUID{}

	return utils.WriteJSON(c, 200, "Team Deleated")
}

//update team name

func UpdateTeamName(c echo.Context) error {
	var payload models.UpdateTeamName
	ctx := context.Background()
	if err := c.Bind(&payload); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}

	user, ok := c.Get("user").(db.User)
	if !ok {
		return utils.WriteError(c, echo.ErrUnauthorized.Code, errors.New("unauthorized"))
	}

	if !user.IsLeader {
		return utils.WriteError(c, echo.ErrBadRequest.Code, errors.New("Only Leader can update name"))
	}

	if err := utils.Queries.UpdateTeamName(ctx, db.UpdateTeamNameParams{
		Name: payload.Name,
		ID:   user.TeamID.UUID,
	}); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}

	return utils.WriteJSON(c, 200, map[string]interface{}{
		"message": "team name updated",
		"data":    payload.Name,
	})
}
