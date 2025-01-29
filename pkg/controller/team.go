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
			Data:   err,
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
	fmt.Println(user)  

	member, err := utils.Queries.GetUser(ctx, user.ID)
	fmt.Println(member.TeamID)
	if member.TeamID.UUID != uuid.Nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "user already in a team",
			},
		})
	}

	team, err := utils.Queries.FindTeam(ctx, payload.Code)
	fmt.Println(team)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return c.JSON(http.StatusBadRequest, models.Response{
				Status: "fail",
				Data:   err,
			})
		}
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Team does not exist",
				"error":   err.Error(),
			},
		})
	}

	nullableTeamID := uuid.NullUUID{
		UUID:  team.ID,
		Valid: true,
	}

	count, err := utils.Queries.CountTeamMembers(ctx, nullableTeamID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to get team members count",
				"error":   err.Error(),
			},
		})
	}

	if count >= 5 {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Cannot join team already full",
			},
		})
	}

	if err := utils.Queries.AddUserToTeam(ctx, db.AddUserToTeamParams{
		TeamID: nullableTeamID,
		ID:     user.ID,
	}); err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Cannot join team",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Queries.IncreaseCountTeam(ctx, nullableTeamID.UUID); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to leave team",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "Team joined successfully",
		},
	})
}

//KICK MEMBER

func KickMemeber(c echo.Context) error {

	var payload models.KickMember

	ctx := c.Request().Context()

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   err,
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
			Data:   err,
		})
	}

	if leader.IsLeader != true {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Only leaders can kick members",
			},
		})
	}

	member, err := utils.Queries.GetUserByID(ctx, payload.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "User not found",
				"error":   err.Error(),
			},
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
			Status: "fail",
			Data: map[string]string{
				"message": "Cannot leave team, already empty",
				"error":   err.Error(),
			},
		})
	}

	if user.TeamID.UUID != member.TeamID.UUID {
		return c.JSON(http.StatusBadGateway, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "User not a memebr of your team",
			},
		})
	}

	if err := utils.Queries.RemoveUserFromTeam(ctx, db.RemoveUserFromTeamParams{
		TeamID: nullableMemberID,
		ID:     member.ID,
	}); err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to leave team",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Queries.DecreaseUserCountTeam(ctx, nullableTeamID.UUID); err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to leave team",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "User kicked successfully",
		},
	})
}

// CREATE TEAM

func CreateTeam(c echo.Context) error {

	var payload models.CreateTeam

	ctx := context.Background()
	if err := c.Bind(&payload); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   validationErrors,
		})
	}

	payload.Name = strings.TrimSpace(payload.Name)

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   "unauthorized",
		})
	}
	// fmt.Println(user)  -- For testing

	member, err := utils.Queries.GetUser(ctx, user.ID)
	//fmt.Println(member.TeamID) -- For testing
	if member.TeamID.UUID != uuid.Nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "User already in a team",
				"error":   err.Error(),
			},
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
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to create Team",
				"error":   err.Error(),
			},
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
					Data: map[string]string{
						"message": "Team already exists",
						"error":   err.Error(),
					},		
				})
			}
		}
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   team,
	})
}

//LEAVE TEAM

func LeaveTeam(c echo.Context) error {

	var payload models.LeaveTeam

	ctx := c.Request().Context()

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
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
			Status: "fail",
			Data: map[string]string{
				"message": "unauthorized",
			},
		})
	}

	member, err := utils.Queries.GetUserByEmail(ctx, user.Email)

	if member.TeamID.UUID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "User not in a team",
				"error":   err.Error(),
			},
		})
	}

	nullableTeamID := uuid.NullUUID{
		UUID:  user.TeamID.UUID,
		Valid: true,
	}

	count, err := utils.Queries.CountTeamMembers(ctx, nullableTeamID)
	if count <= 0 || err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Cannot leave team, already empty",
				"error":   err.Error(),
			},
		})
	}

	if user.IsLeader {
		nullableTeamID := uuid.NullUUID{
			UUID:  user.TeamID.UUID,
			Valid: true,
		}

		emails, err := utils.Queries.GetTeamUsersEmails(ctx, user.TeamID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "Failed to get Team id's",
					"error":   err.Error(),
				},	
			})
		}

		if err := utils.Queries.RemoveTeamIDFromUsers(ctx, nullableTeamID); err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "some error occured",
					"error":   err.Error(),
				},	
			})
		}

		if err := utils.Queries.DeleteTeam(ctx, user.TeamID.UUID); err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "Failed to delete team",
					"error":   err.Error(),
				},	
			})
		}

		if err := utils.Queries.UpdateLeader(ctx, db.UpdateLeaderParams{
			IsLeader: false,
			ID:       user.ID,
		}); err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "Failed to update leader",
					"error":   err.Error(),
				},	
			})
		}
		user.TeamID = uuid.NullUUID{}

		if err := utils.SendTeamEmail(ctx, emails); err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{
				Status: "fail",
				Data: map[string]string{
					"message": "Failed sending the email",
					"error":   err.Error(),
				},
			})
		}

		return c.JSON(http.StatusOK, models.Response{
			Status: "success",
			Data: map[string]string{
				"message": "Team Left Successfully",
			},

		})
	}

	if err := utils.Queries.RemoveUserFromTeam(ctx, db.RemoveUserFromTeamParams{
		ID: user.TeamID.UUID,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed ot remove user",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Queries.LeaveTeam(ctx, user.ID); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Queries.DecreaseUserCountTeam(ctx, nullableTeamID.UUID); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to leave team",
				"error":   err.Error(),
			},
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "Team left successfully",
		},
	})
}

//DELETE TEAM

func DeleteTeam(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(db.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   "unauthorized",
		})
	}

	if !user.IsLeader {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Only leaders can delete team",
			},
		})
	}

	nullableTeamID := uuid.NullUUID{
		UUID:  user.TeamID.UUID,
		Valid: true,
	}

	emails, err := utils.Queries.GetTeamUsersEmails(ctx, user.TeamID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to get email id's",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Queries.RemoveTeamIDFromUsers(ctx, nullableTeamID); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Cannot remove user form team",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Queries.DeleteTeam(ctx, user.TeamID.UUID); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to delete team",
				"error":   err.Error(),
			},
		})
	}

	if err := utils.Queries.UpdateLeader(ctx, db.UpdateLeaderParams{
		IsLeader: false,
		ID:       user.ID,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}
	user.TeamID = uuid.NullUUID{}

	if err := utils.SendTeamEmail(ctx, emails); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Failed to send emails",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "Team deleted successfully",
		},
	})
}

//update team name

func UpdateTeamName(c echo.Context) error {

	var payload models.UpdateTeamName

	ctx := c.Request().Context()

	if err := c.Bind(&payload); err != nil {
		return utils.WriteError(c, echo.ErrBadRequest.Code, err)
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
			Data: map[string]string{
				"message": "unauthorized",
			},

		})
	}

	if !user.IsLeader {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Only leaders can update team",
			},
		})
	}

	if err := utils.Queries.UpdateTeamName(ctx, db.UpdateTeamNameParams{
		Name: payload.Name,
		ID:   user.TeamID.UUID,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "some error occured",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data: map[string]string{
			"message": "Team updated successfully",
		},
	})
}

// Get All users in a team
func GetAllTeamUsers(c echo.Context) error {

	ctx := c.Request().Context()

	user, ok := c.Get("user").(db.User)

	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   "unauthorized",
		})
	}

	team_members, err := utils.Queries.GetTeamUsers(ctx, user.TeamID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data: map[string]string{
				"message": "Cannot get members of the team",
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   team_members,
	})
}
