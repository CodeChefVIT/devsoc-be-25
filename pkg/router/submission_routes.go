package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func SubmissionRoutes(incomingRoutes *echo.Echo) {
	submission := incomingRoutes.Group("/submission")
	submission.Use(middleware.JWTMiddleware())

	submission.POST("", controller.CreateSubmission)
	submission.GET("", controller.GetUserSubmission)
	submission.POST("/:teamId", controller.UpdateSubmission)
	submission.DELETE("/:teamId", controller.DeleteSubmission)
}
