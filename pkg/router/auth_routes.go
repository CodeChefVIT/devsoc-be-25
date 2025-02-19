package router

import (
	"github.com/CodeChefVIT/devsoc-be-24/pkg/controller"
	"github.com/CodeChefVIT/devsoc-be-24/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(incomingRoutes *echo.Group) {
	auth := incomingRoutes.Group("/auth")

	auth.POST("/signup", controller.SignUp)
	auth.POST("/complete-profile", controller.CompleteProfile, middleware.JWTMiddleware())
	auth.POST("/verify-otp", controller.VerifyOTP)
	auth.POST("/login", controller.Login)
	auth.POST("/update-password", controller.UpdatePassword)
	auth.POST("/refresh", controller.RefreshToken)
	auth.GET("/star", controller.CheckStarred, middleware.JWTMiddleware())
	auth.POST("/github", controller.UpdateGithubProfile, middleware.JWTMiddleware())
	auth.POST("/logout", controller.Logout)
	auth.POST("/resend-otp", controller.ResendOTP)

	// auth.Use(middleware.JWTMiddleware())
	// auth.GET("/ping", controller.Ping)
}
