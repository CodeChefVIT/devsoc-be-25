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
