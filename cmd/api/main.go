package main

import (
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/router"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	logger.InitLogger()
	utils.InitCache()
	utils.InitDB()
	utils.InitValidator()
	utils.InitMailer()
}

func main() {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogStatus:     true,
		LogError:      true,
		LogLatency:    true,
		LogMethod:     true,
		LogValuesFunc: logger.RouteLogger,
	}))

	apiGroup := e.Group("")

	router.TeamRoutes(apiGroup)
	router.RegisterRoutes(apiGroup)
	router.IdeaRoutes(apiGroup)
	router.AdminRoutes(apiGroup)
	router.SubmissionRoutes(apiGroup)
	router.AuthRoutes(apiGroup)
	router.PanelRoutes(apiGroup)
	router.InfoRoutes(apiGroup)

	e.Start(":" + utils.Config.Port)
}
