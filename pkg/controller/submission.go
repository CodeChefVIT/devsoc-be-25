package controller

import (
	"net/http"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/dto"
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetSubmission(c echo.Context) error {
	ctx := c.Request().Context()
	teamId := c.Param("teamid")

	teamUuid, err := uuid.Parse(teamId)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	submission, err := utils.Queries.GetSubmissionByTeamID(ctx, teamUuid)
	if err != nil {
		logger.Errorf(logger.InternalError, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	team_sub := dto.Submission{
		GithubLink: submission.GithubLink,
		FigmaLink:  submission.FigmaLink,
		OtherLink:  submission.OtherLink,
		TeamID:     submission.TeamID.String(),
	}

	return c.JSON(http.StatusOK, team_sub)
}
