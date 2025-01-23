package main

import (
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/router"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	logger.InitLogger()
	utils.LoadConfig()
	utils.InitCache()
	utils.InitDB()
	utils.InitValidator()
	utils.InitMailer()

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogStatus:     true,
		LogError:      true,
		LogValuesFunc: logger.RouteLogger,
	}))
	router.TeamRoutes(e)
	router.RegisterRoutes(e)
	router.IdeaRoutes(e)
	router.AdminRoutes(e)
	router.SubmissionRoutes(e)
	router.AuthRoutes(e)
	router.PanelRoutes(e)
	router.InfoRoutes(e)
	utils.Cron()

	e.Start(":" + utils.Config.Port)
}
