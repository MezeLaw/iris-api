package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"iris-api/internal/application/usecases"
	"iris-api/internal/domain/entities"
)

type TurnoHandler struct {
	turnoUseCase usecases.TurnoUseCase
}

func NewTurnoHandler(turnoUseCase usecases.TurnoUseCase) *TurnoHandler {
	return &TurnoHandler{
		turnoUseCase: turnoUseCase,
	}
}

// CreateTurno creates a new turno
// POST /api/v1/turnos
func (h *TurnoHandler) CreateTurno(c *gin.Context) {
	var req entities.CreateTurnoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get clientID from context (set by auth middleware)
	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	turno, err := h.turnoUseCase.CreateTurno(c.Request.Context(), &req, clientID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create turno",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Turno created successfully",
		"data":    turno,
	})
}

// GetTurnoByID retrieves a turno by ID
// GET /api/v1/turnos/:id
func (h *TurnoHandler) GetTurnoByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid turno ID",
			"details": "Turno ID must be a number",
		})
		return
	}

	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	turno, err := h.turnoUseCase.GetTurnoByID(c.Request.Context(), id, clientID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Turno not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Turno retrieved successfully",
		"data":    turno,
	})
}

// GetTurnos retrieves all turnos with filters and pagination
// GET /api/v1/turnos
func (h *TurnoHandler) GetTurnos(c *gin.Context) {
	filters := &entities.TurnoFilters{
		Page:     1,
		PageSize: 10,
	}

	if page := c.Query("page"); page != "" {
		if val, err := strconv.Atoi(page); err == nil && val > 0 {
			filters.Page = val
		}
	}

	if pageSize := c.Query("page_size"); pageSize != "" {
		if val, err := strconv.Atoi(pageSize); err == nil && val > 0 {
			filters.PageSize = val
		}
	}

	if pacienteID := c.Query("paciente_id"); pacienteID != "" {
		if val, err := strconv.ParseInt(pacienteID, 10, 64); err == nil {
			filters.PacienteID = &val
		}
	}

	if profesionalID := c.Query("profesional_user_id"); profesionalID != "" {
		if val, err := strconv.ParseInt(profesionalID, 10, 64); err == nil {
			filters.ProfesionalUserID = &val
		}
	}

	if estado := c.Query("estado"); estado != "" {
		est := entities.EstadoTurno(estado)
		filters.Estado = &est
	}

	if fechaDesde := c.Query("fecha_desde"); fechaDesde != "" {
		if val, err := time.Parse(time.RFC3339, fechaDesde); err == nil {
			filters.FechaDesde = &val
		}
	}

	if fechaHasta := c.Query("fecha_hasta"); fechaHasta != "" {
		if val, err := time.Parse(time.RFC3339, fechaHasta); err == nil {
			filters.FechaHasta = &val
		}
	}

	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	response, err := h.turnoUseCase.GetTurnos(c.Request.Context(), filters, clientID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve turnos",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Turnos retrieved successfully",
		"data":    response,
	})
}

// UpdateTurno updates a turno
// PUT /api/v1/turnos/:id
func (h *TurnoHandler) UpdateTurno(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid turno ID",
			"details": "Turno ID must be a number",
		})
		return
	}

	var req entities.UpdateTurnoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	turno, err := h.turnoUseCase.UpdateTurno(c.Request.Context(), id, &req, clientID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update turno",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Turno updated successfully",
		"data":    turno,
	})
}

// CambiarEstado changes the estado of a turno
// PATCH /api/v1/turnos/:id/estado
func (h *TurnoHandler) CambiarEstado(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid turno ID",
			"details": "Turno ID must be a number",
		})
		return
	}

	var req entities.CambiarEstadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	err = h.turnoUseCase.CambiarEstado(c.Request.Context(), id, &req, clientID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to change turno estado",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Turno estado changed successfully",
	})
}

// DeleteTurno soft deletes a turno
// DELETE /api/v1/turnos/:id
func (h *TurnoHandler) DeleteTurno(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid turno ID",
			"details": "Turno ID must be a number",
		})
		return
	}

	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	err = h.turnoUseCase.DeleteTurno(c.Request.Context(), id, clientID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to delete turno",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Turno deleted successfully",
	})
}

// GetTurnosByDia retrieves turnos for a specific day
// GET /api/v1/turnos/por-dia/:fecha
func (h *TurnoHandler) GetTurnosByDia(c *gin.Context) {
	fechaParam := c.Param("fecha")
	if fechaParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing fecha parameter",
			"details": "fecha parameter is required (format: YYYY-MM-DD)",
		})
		return
	}

	fecha, err := time.Parse("2006-01-02", fechaParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid fecha format",
			"details": "fecha must be in format YYYY-MM-DD",
		})
		return
	}

	var profesionalID *int64
	if profesionalParam := c.Query("profesional_user_id"); profesionalParam != "" {
		if val, err := strconv.ParseInt(profesionalParam, 10, 64); err == nil {
			profesionalID = &val
		}
	}

	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	turnos, err := h.turnoUseCase.GetTurnosByDia(c.Request.Context(), fecha, profesionalID, clientID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve turnos",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Turnos retrieved successfully",
		"data":    turnos,
	})
}

// GetTurnosBySemana retrieves turnos for a specific week
// GET /api/v1/turnos/por-semana/:fecha_inicio
func (h *TurnoHandler) GetTurnosBySemana(c *gin.Context) {
	fechaParam := c.Param("fecha_inicio")
	if fechaParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing fecha_inicio parameter",
			"details": "fecha_inicio parameter is required (format: YYYY-MM-DD)",
		})
		return
	}

	fechaInicio, err := time.Parse("2006-01-02", fechaParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid fecha_inicio format",
			"details": "fecha_inicio must be in format YYYY-MM-DD",
		})
		return
	}

	var profesionalID *int64
	if profesionalParam := c.Query("profesional_user_id"); profesionalParam != "" {
		if val, err := strconv.ParseInt(profesionalParam, 10, 64); err == nil {
			profesionalID = &val
		}
	}

	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	turnos, err := h.turnoUseCase.GetTurnosBySemana(c.Request.Context(), fechaInicio, profesionalID, clientID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve turnos",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Turnos retrieved successfully",
		"data":    turnos,
	})
}

// GetTurnosByProfesional retrieves turnos for a specific professional
// GET /api/v1/turnos/por-profesional/:user_id
func (h *TurnoHandler) GetTurnosByProfesional(c *gin.Context) {
	profesionalIDParam := c.Param("user_id")
	profesionalID, err := strconv.ParseInt(profesionalIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid profesional ID",
			"details": "Profesional ID must be a number",
		})
		return
	}

	var fechaDesde, fechaHasta *time.Time
	if fechaDesdeParam := c.Query("fecha_desde"); fechaDesdeParam != "" {
		if val, err := time.Parse("2006-01-02", fechaDesdeParam); err == nil {
			fechaDesde = &val
		}
	}

	if fechaHastaParam := c.Query("fecha_hasta"); fechaHastaParam != "" {
		if val, err := time.Parse("2006-01-02", fechaHastaParam); err == nil {
			fechaHasta = &val
		}
	}

	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	turnos, err := h.turnoUseCase.GetTurnosByProfesional(c.Request.Context(), profesionalID, fechaDesde, fechaHasta, clientID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve turnos",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Turnos retrieved successfully",
		"data":    turnos,
	})
}

// CheckDisponibilidad checks if a professional is available at a specific time
// POST /api/v1/turnos/disponibilidad
func (h *TurnoHandler) CheckDisponibilidad(c *gin.Context) {
	var req entities.DisponibilidadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Client ID not found in context",
		})
		return
	}

	response, err := h.turnoUseCase.CheckDisponibilidad(c.Request.Context(), &req, clientID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check disponibilidad",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Disponibilidad checked successfully",
		"data":    response,
	})
}
