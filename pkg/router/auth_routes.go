package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(incomingRoutes *echo.Echo) {
	auth := incomingRoutes.Group("/auth")

	auth.POST("/signup", controller.SignUp)
	auth.POST("/send-otp", controller.SendOTP)
	auth.POST("/verify-otp", controller.VerifyOTP)
	auth.POST("/login", controller.Login)
	auth.POST("/update-password", controller.UpdatePassword)
}
