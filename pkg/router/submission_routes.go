package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/labstack/echo/v4"
)

func SubmissionRoutes(incomingRoutes *echo.Echo) {
	submission := incomingRoutes.Group("/submission")
	submission.POST("", controller.CreateSubmission)
	submission.GET("/:teamId", controller.GetSubmission)
	submission.PUT("/:teamId", controller.UpdateSubmission)
	submission.DELETE("/:teamId", controller.DeleteSubmission)
}