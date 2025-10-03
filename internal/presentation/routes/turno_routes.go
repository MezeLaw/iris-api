package routes

import (
	"github.com/gin-gonic/gin"

	"iris-api/internal/presentation/handlers"
)

func SetupTurnoRoutes(router *gin.Engine, turnoHandler *handlers.TurnoHandler) {
	api := router.Group("/api/v1")
	{
		turnos := api.Group("/turnos")
		{
			// CRUD básico
			turnos.POST("/", turnoHandler.CreateTurno)
			turnos.GET("/", turnoHandler.GetTurnos)
			turnos.GET("/:id", turnoHandler.GetTurnoByID)
			turnos.PUT("/:id", turnoHandler.UpdateTurno)
			turnos.DELETE("/:id", turnoHandler.DeleteTurno)

			// Cancelación específica
			turnos.POST("/:id/cancel", turnoHandler.CancelTurno)

			// Vistas específicas
			turnos.GET("/dia", turnoHandler.GetTurnosByDia)
			turnos.GET("/semana", turnoHandler.GetTurnosBySemana)

			// Alertas
			turnos.GET("/alertas", turnoHandler.GetProximosTurnosAlert)
		}

		// Turnos por profesional
		profesionales := api.Group("/profesionales")
		{
			profesionales.GET("/:contactologo_id/turnos", turnoHandler.GetTurnosByProfesional)
		}
	}
}
