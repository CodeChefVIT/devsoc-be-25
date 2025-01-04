package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/labstack/echo/v4"
)

func PanelRoutes(incomingRoutes *echo.Echo) {
	panel := incomingRoutes.Group("/panel")
	panel.POST("/createscore", controller.CreateScore)
	panel.DELETE("/deletescore/:id", controller.DeleteScore)
	panel.GET("/getscore/:teamid", controller.GetScore)
	panel.PUT("/updatescore/:id", controller.UpdateScore)
}
