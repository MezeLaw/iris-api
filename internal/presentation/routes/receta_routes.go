package routes

import (
	"github.com/gin-gonic/gin"

	"iris-api/internal/presentation/handlers"
)

func SetupRecetaRoutes(router *gin.Engine, recetaHandler *handlers.RecetaHandler) {
	api := router.Group("/api/v1")
	{
		recetas := api.Group("/recetas")
		{
			// CRUD básico
			recetas.POST("/", recetaHandler.CreateReceta)
			recetas.GET("/", recetaHandler.GetRecetas)
			recetas.GET("/:id", recetaHandler.GetRecetaByID)
			recetas.PUT("/:id", recetaHandler.UpdateReceta)
			recetas.DELETE("/:id", recetaHandler.DeleteReceta)
		}

		// Endpoints específicos por paciente
		pacientes := api.Group("/pacientes")
		{
			pacientes.GET("/:paciente_id/recetas", recetaHandler.GetRecetasByPacienteID)
			pacientes.GET("/:paciente_id/recetas/historial", recetaHandler.GetHistorialByPacienteID)
			pacientes.GET("/:paciente_id/recetas/alertas", recetaHandler.CheckDioptriasChange)
		}
	}
}