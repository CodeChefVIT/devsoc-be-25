package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"
	"github.com/labstack/echo/v4"
)
func InfoRoutes(incomingRoutes *echo.Echo) {
    info := incomingRoutes.Group("/info")
    info.Use(middleware.JWTMiddleware())

    info.GET("/me", controller.GetDetails)
    info.POST("/me", controller.UpdateUser)
}