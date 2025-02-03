package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	//"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func AdminRoutes(incomingRoutes *echo.Group) {
	admin := incomingRoutes.Group("/admin")
	// admin.Use(middleware.JWTMiddleware())
	// admin.Use(middleware.CheckAdmin)

	admin.GET("/users", controller.GetAllUsers)
	admin.GET("/user/:email", controller.GetUsersByEmail)
	admin.POST("/ban", controller.BanUser)
	admin.POST("/unban", controller.UnbanUser)
	admin.POST("/star", controller.CheckStarred)
	admin.GET("/users/:gender", controller.GetUsersByGender)

	admin.GET("/leaderboard", controller.GetLeaderBoard)

	admin.GET("/teams", controller.GetTeams)
	admin.GET("/teams/:id", controller.GetTeamById)
	admin.GET("/team/leader/:id", controller.GetTeamLeader)
	admin.POST("/createpanel", controller.CreatePanel)
	admin.GET("/teams/track/:track", controller.GetTeamsByTrack)

	admin.GET("/members/:id", controller.GetAllTeamMembers)
	admin.POST("/ban/team", controller.BanTeam)
	admin.POST("/unban/team", controller.UnBanTeam)

	admin.GET("/usercsv", controller.ExportUsers)
	admin.GET("/teamcsv", controller.ExportTeams)
	admin.PUT("/team/rounds", controller.UpdateTeamRounds)

	admin.GET("/ideas", controller.GetAllIdeas)
	admin.GET("/ideas/filter", controller.GetIdeasByTrack)
}
