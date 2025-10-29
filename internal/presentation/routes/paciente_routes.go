package routes

import (
	"github.com/gin-gonic/gin"

	"iris-api/internal/presentation/handlers"
	"iris-api/internal/presentation/middleware"
)

func SetupPacienteRoutes(router *gin.Engine, pacienteHandler *handlers.PacienteHandler, authMiddleware *middleware.AuthMiddleware) {
	api := router.Group("/api/v1")
	{
		// Todas las rutas de pacientes requieren autenticación
		pacientes := api.Group("/pacientes")
		pacientes.Use(authMiddleware.RequireAuth())
		{
			// CRUD Pacientes
			pacientes.POST("/", pacienteHandler.CreatePaciente)
			pacientes.GET("/", pacienteHandler.ListPacientes)
			pacientes.GET("/search", pacienteHandler.SearchPacientes)
			pacientes.GET("/:id", pacienteHandler.GetPacienteByID)
			pacientes.PUT("/:id", pacienteHandler.UpdatePaciente)
			pacientes.DELETE("/:id", pacienteHandler.DeletePaciente)

			// Antecedentes Médicos
			pacientes.PUT("/:id/antecedentes-medicos", pacienteHandler.CreateOrUpdateAntecedentesMedicos)
			pacientes.GET("/:id/antecedentes-medicos", pacienteHandler.GetAntecedentesMedicos)

			// Antecedentes Visuales
			pacientes.PUT("/:id/antecedentes-visuales", pacienteHandler.CreateOrUpdateAntecedentesVisuales)
			pacientes.GET("/:id/antecedentes-visuales", pacienteHandler.GetAntecedentesVisuales)

			// Exámenes Visuales
			pacientes.GET("/:id/examenes", pacienteHandler.GetExamenesVisualesByPaciente)
		}

		// Rutas de exámenes visuales (pueden estar asociadas a cualquier paciente)
		examenes := api.Group("/examenes-visuales")
		examenes.Use(authMiddleware.RequireAuth())
		{
			examenes.POST("/", pacienteHandler.CreateExamenVisual)
			examenes.GET("/:examenId", pacienteHandler.GetExamenVisual)
			examenes.PUT("/:examenId", pacienteHandler.UpdateExamenVisual)
			examenes.DELETE("/:examenId", pacienteHandler.DeleteExamenVisual)
			examenes.GET("/comparar", pacienteHandler.CompararExamenes)
		}
	}
}
