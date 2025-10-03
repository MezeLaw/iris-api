package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"iris-api/internal/application/usecases"
	"iris-api/internal/presentation/middleware"
)

type ReporteriaHandler struct {
	reporteriaUseCase *usecases.ReporteriaUseCase
}

func NewReporteriaHandler(reporteriaUseCase *usecases.ReporteriaUseCase) *ReporteriaHandler {
	return &ReporteriaHandler{
		reporteriaUseCase: reporteriaUseCase,
	}
}

func (h *ReporteriaHandler) GetPacientesActivos(c *gin.Context) {
	// Get client_id from authenticated user context
	clientID, exists := middleware.GetClientID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	reporte, err := h.reporteriaUseCase.GetPacientesActivos(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve active patients report",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Active patients report retrieved successfully",
		"data":    reporte,
	})
}

func (h *ReporteriaHandler) GetPacientesInactivos(c *gin.Context) {
	// Get client_id from authenticated user context
	clientID, exists := middleware.GetClientID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	reporte, err := h.reporteriaUseCase.GetPacientesInactivos(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve inactive patients report",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inactive patients report retrieved successfully",
		"data":    reporte,
	})
}
