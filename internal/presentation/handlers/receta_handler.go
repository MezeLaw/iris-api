package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"iris-api/internal/application/usecases"
	"iris-api/internal/domain/entities"
)

type RecetaHandler struct {
	recetaUseCase usecases.RecetaUseCase
}

func NewRecetaHandler(recetaUseCase usecases.RecetaUseCase) *RecetaHandler {
	return &RecetaHandler{
		recetaUseCase: recetaUseCase,
	}
}

func (h *RecetaHandler) CreateReceta(c *gin.Context) {
	var req entities.CreateRecetaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	receta, err := h.recetaUseCase.CreateReceta(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create receta",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Receta created successfully",
		"data":    receta,
	})
}

func (h *RecetaHandler) GetRecetaByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid receta ID",
			"details": "Receta ID must be a number",
		})
		return
	}

	receta, err := h.recetaUseCase.GetRecetaByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Receta not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Receta retrieved successfully",
		"data":    receta,
	})
}

func (h *RecetaHandler) GetRecetas(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	offsetParam := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	response, err := h.recetaUseCase.GetRecetas(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve recetas",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Recetas retrieved successfully",
		"data":    response,
	})
}

func (h *RecetaHandler) GetRecetasByPacienteID(c *gin.Context) {
	pacienteIDParam := c.Param("paciente_id")
	pacienteID, err := strconv.ParseInt(pacienteIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	limitParam := c.DefaultQuery("limit", "10")
	offsetParam := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	response, err := h.recetaUseCase.GetRecetasByPacienteID(c.Request.Context(), pacienteID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve recetas for paciente",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Recetas retrieved successfully",
		"data":    response,
	})
}

func (h *RecetaHandler) GetHistorialByPacienteID(c *gin.Context) {
	pacienteIDParam := c.Param("paciente_id")
	pacienteID, err := strconv.ParseInt(pacienteIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	recetas, err := h.recetaUseCase.GetHistorialByPacienteID(c.Request.Context(), pacienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve historial",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Historial retrieved successfully",
		"data":    recetas,
	})
}

func (h *RecetaHandler) CheckDioptriasChange(c *gin.Context) {
	pacienteIDParam := c.Param("paciente_id")
	pacienteID, err := strconv.ParseInt(pacienteIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	change, err := h.recetaUseCase.CheckDioptriasChange(c.Request.Context(), pacienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check dioptrias change",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dioptrias change checked successfully",
		"data":    change,
	})
}

func (h *RecetaHandler) UpdateReceta(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid receta ID",
			"details": "Receta ID must be a number",
		})
		return
	}

	var req entities.UpdateRecetaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	receta, err := h.recetaUseCase.UpdateReceta(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update receta",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Receta updated successfully",
		"data":    receta,
	})
}

func (h *RecetaHandler) DeleteReceta(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid receta ID",
			"details": "Receta ID must be a number",
		})
		return
	}

	err = h.recetaUseCase.DeleteReceta(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to delete receta",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}