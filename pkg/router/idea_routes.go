package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func IdeaRoutes(incomingRoutes *echo.Group) {
	idea := incomingRoutes.Group("/idea")
	idea.Use(middleware.JWTMiddleware())
	idea.Use(middleware.CheckTeamBan)
	idea.Use(middleware.CheckUserVerifiation)

	idea.POST("/create", controller.CreateIdea)
	idea.PUT("/update", controller.UpdateIdea)
	idea.GET("/", controller.GetIdea)
}
