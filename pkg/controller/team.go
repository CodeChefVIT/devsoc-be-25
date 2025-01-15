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
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: err,
		})
	}

	ctx := c.Request().Context()

	if err := utils.Validate.Struct(payload); err != nil {
        validationErrors := utils.FormatValidationErrors(err)
        return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   validationErrors,
		})
    }

	user := c.Get("user").(db.User)
	// fmt.Println(user)  -- For testing 

	member,err := utils.Queries.GetUser(ctx, user.ID)
	fmt.Println(member.TeamID)
	if member.TeamID.UUID != uuid.Nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: "user already in a team",
		})
	}

	team, err := utils.Queries.FindTeam(ctx, payload.Code)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return c.JSON(http.StatusBadRequest, models.Response{
				Status:"fail",
				Data:err,
			})
		}
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"failed to join team",
		})
	}

	nullableTeamID := uuid.NullUUID{
		UUID:  team.ID,
		Valid: true,
	}

	count, err := utils.Queries.CountTeamMembers(ctx, nullableTeamID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"failed to get team member count",
		})
	}

	if count >= 5 {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"Cannot join Team. Team is already full",
		})
	}


	if err := utils.Queries.AddUserToTeam(ctx, db.AddUserToTeamParams{
		TeamID: nullableTeamID,
		ID:     user.ID,
	}); err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{
			Status:"fail",
			Data:"Cannot join Team",
		})
	}

	if err := utils.Queries.IncreaseCountTeam(ctx, nullableTeamID.UUID); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:"Failed to add member to team",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:"success",
		Data:"Team joined successfully",
	})
}

//KICK MEMBER

func KickMemeber(c echo.Context) error {

	var payload models.KickMember

	ctx := c.Request().Context()

	if err := c.Bind(&payload); err!=nil{
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:err,
		})
	}

	if err := utils.Validate.Struct(payload); err != nil {
        validationErrors := utils.FormatValidationErrors(err)
        return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   validationErrors,
		})
    }

	user := c.Get("user").(db.User)

	leader, err := utils.Queries.GetUser(ctx, user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: err,
		})
	}

	if leader.IsLeader != true {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: "Only leader can kick Members",
		})
	}

	member, err := utils.Queries.GetUserByID(ctx, payload.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"User not found",
		})
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
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"Cannot leave team. Team is already empty",
		})
	}

	if user.TeamID.UUID != member.TeamID.UUID {
		return c.JSON(http.StatusBadGateway, models.Response{
			Status:"fail",
			Data:"User is not a member of your team",
		})
	}

	if err := utils.Queries.RemoveUserFromTeam(ctx, db.RemoveUserFromTeamParams{
		TeamID: nullableMemberID,
		ID:     member.ID,
	}); err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{
			Status:"fail",
			Data:"Failed to kick user",
		})
	}

	if err := utils.Queries.DecreaseUserCountTeam(ctx, nullableTeamID.UUID); err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{
			Status:"fail",
			Data:"Failed to leave Team",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:"success",
		Data:"User kicked successfully",
	})
}

// CREATE TEAM

func CreateTeam(c echo.Context) error {

	var payload models.CreateTeam

	ctx := context.Background()
	if err := c.Bind(&payload);err!=nil{
		return utils.WriteError(c, echo.ErrBadRequest.Code,err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
        validationErrors := utils.FormatValidationErrors(err)
        return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   validationErrors,
		})
    }

	payload.Name = strings.TrimSpace(payload.Name)

	user,ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"unauthorized",
		})
	}
	// fmt.Println(user)  -- For testing 

	member,err := utils.Queries.GetUser(ctx, user.ID)
	//fmt.Println(member.TeamID) -- For testing
	if member.TeamID.UUID != uuid.Nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: "user already in a team",
		})
	}

	//fmt.Println(user.TeamID.UUID.String()) -- testing

	params := db.CreateTeamParams{
		ID:             uuid.New(),
		Name:           payload.Name,
		Code:           utils.GenerateRandomString(6),
		NumberOfPeople: 1,
		RoundQualified: pgtype.Int4{Int32: 0, Valid: true},
		IsBanned:	false,
	}

	team, err := utils.Queries.CreateTeam(ctx, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:err,
		})
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
				return c.JSON(http.StatusBadRequest, models.Response{
					Status: "fail",
					Data:"Team name already exists",
				})
			}
		}
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:err,
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:"success",
		Data:team,
	})
}

//LEAVE TEAM

func LeaveTeam(c echo.Context) error {

	var payload models.LeaveTeam

	ctx := c.Request().Context()

	if err := c.Bind(&payload);err != nil{
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data: err,
		})
	}

	if err := utils.Validate.Struct(payload); err != nil {
        validationErrors := utils.FormatValidationErrors(err)
        return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   validationErrors,
		})
    }

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"unauthorized",
		})
	}

	member,err := utils.Queries.GetUser(ctx, user.ID)

	if member.TeamID.UUID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: "user not in a team",
		})
	}

	nullableTeamID := uuid.NullUUID{
		UUID:  user.TeamID.UUID,
		Valid: true,
	}

	count, err := utils.Queries.CountTeamMembers(ctx, nullableTeamID)
	if count <= 0 || err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"Cannot leave Team. Team is already empty",
		})
	}

	if user.IsLeader {
		nullableTeamID := uuid.NullUUID{
			UUID:  user.TeamID.UUID,
			Valid: true,
		}

		emails, err := utils.Queries.GetTeamUsersEmails(ctx, user.TeamID)
		if err != nil{
			return c.JSON(http.StatusBadRequest, models.Response{
				Status:"fail",
				Data:"Failed to get email id's",
			})
		}
	
		if err := utils.Queries.RemoveTeamIDFromUsers(ctx, nullableTeamID); err != nil{
			return c.JSON(http.StatusBadRequest, models.Response{
				Status:"fail",
				Data:err,
			})
		}
	
		if err := utils.Queries.DeleteTeam(ctx, user.TeamID.UUID); err != nil{
			return c.JSON(http.StatusBadRequest, models.Response{
				Status:"fail",
				Data:err,
			})
		}
	
		if err := utils.Queries.UpdateLeader(ctx, db.UpdateLeaderParams{
			IsLeader: false,
			ID: user.ID,
		});err != nil{
			return c.JSON(http.StatusBadRequest, models.Response{
				Status:"fail",
				Data:err,
			})
		}
		user.TeamID = uuid.NullUUID{}

		if err := utils.SendTeamEmail(ctx, emails); err!= nil {
			return c.JSON(http.StatusBadRequest, models.Response{
				Status:"fail",
				Data:"Failed send emails",
			})
		}

		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"success",
			Data:"Team Left and Deleted successfully",
		})
	}

	if err := utils.Queries.RemoveUserFromTeam(ctx, db.RemoveUserFromTeamParams{
		ID: user.TeamID.UUID,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:err,
		})
	}

	if err := utils.Queries.LeaveTeam(ctx, payload.UserID); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:err,
		})
	}

	if err := utils.Queries.DecreaseUserCountTeam(ctx,nullableTeamID.UUID);err!=nil{
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"Failed to leave team",
		})
	}
	return c.JSON(http.StatusOK,models.Response{
		Status:"success",
		Data:"Team Left Successfully",
	})
}

//DELETE TEAM

func DeleteTeam(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"unauthorized",
		})
	}

	if !user.IsLeader{
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"Only leader can delete the team",
		})
	}

	nullableTeamID := uuid.NullUUID{
		UUID:  user.TeamID.UUID,
		Valid: true,
	}

	emails, err := utils.Queries.GetTeamUsersEmails(ctx, user.TeamID)
	if err != nil{
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
				Data:"Failed to get email id's",
		})
	}

	if err := utils.Queries.RemoveTeamIDFromUsers(ctx, nullableTeamID); err != nil{
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:err,
		})
	}

	if err := utils.Queries.DeleteTeam(ctx, user.TeamID.UUID); err != nil{
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:err,
		})
	}

	if err := utils.Queries.UpdateLeader(ctx, db.UpdateLeaderParams{
		IsLeader: false,
		ID: user.ID,
	});err != nil{
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:err,
		})
	}
	user.TeamID = uuid.NullUUID{}

	if err := utils.SendTeamEmail(ctx, emails); err!= nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"Failed send emails",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:"success",
		Data:"Team Deleted",
	})
}

//update team name

func UpdateTeamName(c echo.Context) error {

	var payload models.UpdateTeamName

	ctx := c.Request().Context()

	if err := c.Bind(&payload);err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code,err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
        validationErrors := utils.FormatValidationErrors(err)
        return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   validationErrors,
		})
    }

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: "unauthorized",
		})
	}

	if !user.IsLeader {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"Only Leader can update me",
		})
	}

	if err := utils.Queries.UpdateTeamName(ctx, db.UpdateTeamNameParams{
		Name: payload.Name,
		ID:   user.TeamID.UUID,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:err,
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:"success",
		Data:"Team name updated successfully",
	})
}

//Get All users in a team
func GetAllTeamUsers (c echo.Context) error {

	ctx := c.Request().Context()

	user,ok := c.Get("user").(db.User)

	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"unauthorized",
		})
	}

	team_members, err := utils.Queries.GetTeamUsers(ctx,user.TeamID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:"fail",
			Data:"Cannot get Members of the team",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:"success",
		Data:team_members,
	})
}