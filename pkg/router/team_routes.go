package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"

	"github.com/labstack/echo/v4"
)

func TeamRoutes(incomingRoutes *echo.Echo) {
	auth := incomingRoutes.Group("/team")

	auth.Use(middleware.JWTMiddleware())

	auth.POST("/join", controller.JoinTeam)
	auth.POST("/create",controller.CreateTeam)
	auth.POST("/leave",controller.LeaveTeam)
	auth.POST("/kick",controller.KickMemeber)
	auth.POST("/delete",controller.DeleteTeam)
	auth.PUT("/update",controller.UpdateTeamName)
}
