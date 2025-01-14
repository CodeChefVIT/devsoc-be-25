package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func AdminRoutes(incomingRoutes *echo.Echo) {
	admin := incomingRoutes.Group("/admin")
	admin.Use(middleware.JWTMiddleware())
	admin.Use(middleware.CheckAdmin)

	admin.GET("/users", controller.GetAllUsers)
	admin.GET("/vitians", controller.GetAllVitians)
	admin.GET("/user/:email", controller.GetUsersByEmail)
	admin.POST("/ban", controller.BanUser)
	admin.POST("/unban", controller.UnbanUser)

	admin.GET("/teams", controller.GetTeams)
	admin.GET("/teams/:id", controller.GetTeamById)
	admin.GET("/team/leader/:id", controller.GetTeamLeader)
	admin.POST("/createpanel", controller.CreatePanel)

	admin.GET("/members/:id", controller.GetAllTeamMembers)

	admin.GET("/usercsv", controller.ExportUsers)
	admin.GET("/teamcsv", controller.ExportTeams)
}
