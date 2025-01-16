package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func SubmissionRoutes(incomingRoutes *echo.Echo) {
	submission := incomingRoutes.Group("/submission")
	submission.Use(middleware.JWTMiddleware())
	submission.Use(middleware.CheckTeamBan)

	submission.POST("/create", controller.CreateSubmission)
	submission.GET("/get", controller.GetUserSubmission)
	submission.POST("/update", controller.UpdateSubmission)
	submission.DELETE("/delete", controller.DeleteSubmission)
}
