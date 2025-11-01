package routes

import (
	"github.com/gin-gonic/gin"

	"iris-api/internal/domain/entities"
	"iris-api/internal/presentation/handlers"
	"iris-api/internal/presentation/middleware"
)

func SetupReporteriaRoutes(router *gin.Engine, reporteriaHandler *handlers.ReporteriaHandler, authMiddleware *middleware.AuthMiddleware) {
	api := router.Group("/api/v1")
	{
		reportes := api.Group("/reportes")
		// Require authentication for all reports
		reportes.Use(authMiddleware.RequireAuth())
		// Only admin and optometrista can access reports
		reportes.Use(authMiddleware.RequireRole(entities.RoleAdmin, entities.RoleOptometrista))
		{
			reportes.GET("/pacientes-activos", reporteriaHandler.GetPacientesActivos)
			reportes.GET("/pacientes-inactivos", reporteriaHandler.GetPacientesInactivos)
		}
	}
}
