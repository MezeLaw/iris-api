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
	turnoUseCase *usecases.TurnoUseCase
}

func NewTurnoHandler(turnoUseCase *usecases.TurnoUseCase) *TurnoHandler {
	return &TurnoHandler{
		turnoUseCase: turnoUseCase,
	}
}

func (h *TurnoHandler) CreateTurno(c *gin.Context) {
	var req entities.CreateTurnoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	turno, err := h.turnoUseCase.CreateTurno(c.Request.Context(), &req)
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

	turno, err := h.turnoUseCase.GetTurnoByID(c.Request.Context(), id)
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

func (h *TurnoHandler) GetTurnos(c *gin.Context) {
	filter := &entities.TurnoFilter{
		Limit:  10,
		Offset: 0,
	}

	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil && val > 0 {
			filter.Limit = val
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil && val >= 0 {
			filter.Offset = val
		}
	}

	if pacienteID := c.Query("paciente_id"); pacienteID != "" {
		if val, err := strconv.ParseInt(pacienteID, 10, 64); err == nil {
			filter.PacienteID = &val
		}
	}

	if contactologoID := c.Query("contactologo_id"); contactologoID != "" {
		if val, err := strconv.ParseInt(contactologoID, 10, 64); err == nil {
			filter.ContactologoID = &val
		}
	}

	if tipoServicio := c.Query("tipo_servicio"); tipoServicio != "" {
		ts := entities.TipoServicio(tipoServicio)
		filter.TipoServicio = &ts
	}

	if estado := c.Query("estado"); estado != "" {
		est := entities.EstadoTurno(estado)
		filter.Estado = &est
	}

	if fechaDesde := c.Query("fecha_desde"); fechaDesde != "" {
		if val, err := time.Parse(time.RFC3339, fechaDesde); err == nil {
			filter.FechaDesde = &val
		}
	}

	if fechaHasta := c.Query("fecha_hasta"); fechaHasta != "" {
		if val, err := time.Parse(time.RFC3339, fechaHasta); err == nil {
			filter.FechaHasta = &val
		}
	}

	response, err := h.turnoUseCase.GetTurnos(c.Request.Context(), filter)
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

	turno, err := h.turnoUseCase.UpdateTurno(c.Request.Context(), id, &req)
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

func (h *TurnoHandler) CancelTurno(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid turno ID",
			"details": "Turno ID must be a number",
		})
		return
	}

	var req struct {
		Motivo string `json:"motivo" validate:"required,max=500"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	err = h.turnoUseCase.CancelTurno(c.Request.Context(), id, req.Motivo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to cancel turno",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Turno cancelled successfully",
	})
}

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

	err = h.turnoUseCase.DeleteTurno(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to delete turno",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// Vista por día
func (h *TurnoHandler) GetTurnosByDia(c *gin.Context) {
	fechaParam := c.Query("fecha")
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

	var contactologoID *int64
	if contactologoParam := c.Query("contactologo_id"); contactologoParam != "" {
		if val, err := strconv.ParseInt(contactologoParam, 10, 64); err == nil {
			contactologoID = &val
		}
	}

	turnos, err := h.turnoUseCase.GetTurnosByDia(c.Request.Context(), fecha, contactologoID)
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

// Vista por semana
func (h *TurnoHandler) GetTurnosBySemana(c *gin.Context) {
	fechaParam := c.Query("fecha_inicio")
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

	var contactologoID *int64
	if contactologoParam := c.Query("contactologo_id"); contactologoParam != "" {
		if val, err := strconv.ParseInt(contactologoParam, 10, 64); err == nil {
			contactologoID = &val
		}
	}

	turnos, err := h.turnoUseCase.GetTurnosBySemana(c.Request.Context(), fechaInicio, contactologoID)
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

// Vista por profesional
func (h *TurnoHandler) GetTurnosByProfesional(c *gin.Context) {
	contactologoIDParam := c.Param("contactologo_id")
	contactologoID, err := strconv.ParseInt(contactologoIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid contactologo ID",
			"details": "Contactologo ID must be a number",
		})
		return
	}

	fechaDesdeParam := c.Query("fecha_desde")
	fechaHastaParam := c.Query("fecha_hasta")

	if fechaDesdeParam == "" || fechaHastaParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing date parameters",
			"details": "fecha_desde and fecha_hasta are required (format: YYYY-MM-DD)",
		})
		return
	}

	fechaDesde, err := time.Parse("2006-01-02", fechaDesdeParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid fecha_desde format",
			"details": "fecha_desde must be in format YYYY-MM-DD",
		})
		return
	}

	fechaHasta, err := time.Parse("2006-01-02", fechaHastaParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid fecha_hasta format",
			"details": "fecha_hasta must be in format YYYY-MM-DD",
		})
		return
	}

	turnos, err := h.turnoUseCase.GetTurnosByProfesional(c.Request.Context(), contactologoID, fechaDesde, fechaHasta)
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

// Alertas de turnos próximos
func (h *TurnoHandler) GetProximosTurnosAlert(c *gin.Context) {
	alert, err := h.turnoUseCase.GetProximosTurnosAlert(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve alerts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Alerts retrieved successfully",
		"data":    alert,
	})
}
