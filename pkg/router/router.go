package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"

	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	admin := e.Group("/users")
	admin.Use(echojwt.JWT(utils.Config.JwtSecret))
	admin.GET("/ping", controller.Ping)

	e.POST("/star", controller.CheckStarred)
	e.GET("/ping", controller.Ping)
	e.GET("/docs", controller.Docs)
}
