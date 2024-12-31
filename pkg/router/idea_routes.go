package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/labstack/echo/v4"
)

func IdeaRoutes(incomingRoutes *echo.Echo) {
	idea := incomingRoutes.Group("/idea")
	idea.POST("/create", controller.CreateIdea)
	idea.PUT("/update/:id", controller.UpdateIdea)
	idea.GET("", controller.GetIdea)
}
