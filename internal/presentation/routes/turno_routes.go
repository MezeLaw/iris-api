package routes

import (
	"github.com/gin-gonic/gin"

	"iris-api/internal/presentation/handlers"
	"iris-api/internal/presentation/middleware"
)

func SetupTurnoRoutes(router *gin.Engine, turnoHandler *handlers.TurnoHandler, authMiddleware *middleware.AuthMiddleware) {
	api := router.Group("/api/v1")
	{
		// All turno routes require authentication
		turnos := api.Group("/turnos")
		turnos.Use(authMiddleware.RequireAuth())
		{
			// CRUD básico
			turnos.POST("/", turnoHandler.CreateTurno)
			turnos.GET("/", turnoHandler.GetTurnos)
			turnos.GET("/:id", turnoHandler.GetTurnoByID)
			turnos.PUT("/:id", turnoHandler.UpdateTurno)
			turnos.DELETE("/:id", turnoHandler.DeleteTurno)

			// Cambiar estado
			turnos.PATCH("/:id/estado", turnoHandler.CambiarEstado)

			// Disponibilidad
			turnos.POST("/disponibilidad", turnoHandler.CheckDisponibilidad)

			// Vistas específicas
			turnos.GET("/por-dia/:fecha", turnoHandler.GetTurnosByDia)
			turnos.GET("/por-semana/:fecha_inicio", turnoHandler.GetTurnosBySemana)
			turnos.GET("/por-profesional/:user_id", turnoHandler.GetTurnosByProfesional)
		}
	}
}
