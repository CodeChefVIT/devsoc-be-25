package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"
	"github.com/labstack/echo/v4"
)
func UserRoutes(incomingRoutes *echo.Echo) {
    info := incomingRoutes.Group("/info")
    info.Use(middleware.Protected())

    info.GET("/me", controller.GetDetails)
    info.PATCH("/me", controller.UpdateUser)
}