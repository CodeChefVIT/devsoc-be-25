package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"

	"github.com/labstack/echo/v4"
)

func TeamRoutes(incomingRoutes *echo.Echo) {
	team := incomingRoutes.Group("/team")

	team.Use(middleware.JWTMiddleware())
	team.Use(middleware.CheckUserBan)

	team.POST("/join", controller.JoinTeam)
	team.POST("/create",controller.CreateTeam)
	team.POST("/leave",controller.LeaveTeam)
	team.POST("/kick",controller.KickMemeber)
	team.POST("/delete",controller.DeleteTeam)
	team.PUT("/update",controller.UpdateTeamName)
	team.GET("/users",controller.GetAllTeamUsers)
}
