package routes

import (
	"github.com/gin-gonic/gin"

	"iris-api/internal/presentation/handlers"
)

func SetupUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}
}
