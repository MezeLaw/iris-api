package routes

import (
	"github.com/gin-gonic/gin"

	"iris-api/internal/presentation/handlers"
	"iris-api/internal/presentation/middleware"
)

func SetupAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	auth := router.Group("/api/v1/auth")
	{
		// Public routes
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)

		// Protected routes
		auth.GET("/profile", authMiddleware.RequireAuth(), authHandler.GetProfile)
	}
}
