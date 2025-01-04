package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func AdminRoutes(incomingRoutes *echo.Echo) {
	admin := incomingRoutes.Group("/admin")
	admin.Use(middleware.Protected())
	admin.Use(middleware.CheckAdmin)

	admin.GET("/users", controller.GetAllUsers)
	admin.GET("/vitians", controller.GetAllVitians)
	admin.GET("/user/:email", controller.GetUsersByEmail, middleware.CheckAdmin)
	admin.POST("/ban", controller.BanUser, middleware.CheckAdmin)
	admin.POST("/unban", controller.UnbanUser, middleware.CheckAdmin)

	admin.GET("/teams", controller.GetTeams)
	admin.GET("/teams/:id", controller.GetTeamById)
	admin.GET("/team/leader/:id", controller.GetTeamLeader)
	admin.POST("/createpanel", controller.CreatePanel)
}
