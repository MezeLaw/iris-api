package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"iris-api/internal/application/usecases"
	"iris-api/internal/domain/entities"
)

type PacienteHandler struct {
	pacienteUseCase usecases.PacienteUseCase
}

func NewPacienteHandler(pacienteUseCase usecases.PacienteUseCase) *PacienteHandler {
	return &PacienteHandler{
		pacienteUseCase: pacienteUseCase,
	}
}

// Helper para obtener client_id del contexto (seteado por middleware de autenticación)
func getClientID(c *gin.Context) (int64, error) {
	clientID, exists := c.Get("client_id")
	if !exists {
		return 0, gin.Error{Err: gin.Error{}.Err, Type: gin.ErrorTypePrivate}
	}
	return clientID.(int64), nil
}

// CRUD Pacientes

func (h *PacienteHandler) CreatePaciente(c *gin.Context) {
	var req entities.CreatePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Obtener client_id del contexto (usuario autenticado)
	clientID, err := getClientID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized - client_id not found",
		})
		return
	}
	req.ClientID = clientID

	paciente, err := h.pacienteUseCase.CreatePaciente(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create paciente",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Paciente created successfully",
		"data":    paciente,
	})
}

func (h *PacienteHandler) GetPacienteByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	clientID, err := getClientID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized - client_id not found",
		})
		return
	}

	// Verificar si se solicita información completa
	complete := c.Query("complete")
	var paciente *entities.Paciente

	if complete == "true" {
		paciente, err = h.pacienteUseCase.GetPacienteComplete(c.Request.Context(), id, clientID)
	} else {
		paciente, err = h.pacienteUseCase.GetPacienteByID(c.Request.Context(), id, clientID)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Paciente not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Paciente retrieved successfully",
		"data":    paciente,
	})
}

func (h *PacienteHandler) ListPacientes(c *gin.Context) {
	pageParam := c.DefaultQuery("page", "1")
	pageSizeParam := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	clientID, err := getClientID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized - client_id not found",
		})
		return
	}

	response, err := h.pacienteUseCase.ListPacientes(c.Request.Context(), clientID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve pacientes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pacientes retrieved successfully",
		"data":    response,
	})
}

func (h *PacienteHandler) SearchPacientes(c *gin.Context) {
	query := c.Query("q")
	pageParam := c.DefaultQuery("page", "1")
	pageSizeParam := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	clientID, err := getClientID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized - client_id not found",
		})
		return
	}

	response, err := h.pacienteUseCase.SearchPacientes(c.Request.Context(), clientID, query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search pacientes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Search completed successfully",
		"data":    response,
	})
}

func (h *PacienteHandler) UpdatePaciente(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	clientID, err := getClientID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized - client_id not found",
		})
		return
	}

	var req entities.UpdatePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	paciente, err := h.pacienteUseCase.UpdatePaciente(c.Request.Context(), id, clientID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update paciente",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Paciente updated successfully",
		"data":    paciente,
	})
}

func (h *PacienteHandler) DeletePaciente(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	clientID, err := getClientID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized - client_id not found",
		})
		return
	}

	err = h.pacienteUseCase.DeletePaciente(c.Request.Context(), id, clientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to delete paciente",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// Antecedentes Médicos

func (h *PacienteHandler) CreateOrUpdateAntecedentesMedicos(c *gin.Context) {
	idParam := c.Param("id")
	pacienteID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	var req entities.CreateAntecedentesMedicosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	antecedentes, err := h.pacienteUseCase.CreateOrUpdateAntecedentesMedicos(c.Request.Context(), pacienteID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create/update antecedentes medicos",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Antecedentes medicos saved successfully",
		"data":    antecedentes,
	})
}

func (h *PacienteHandler) GetAntecedentesMedicos(c *gin.Context) {
	idParam := c.Param("id")
	pacienteID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	antecedentes, err := h.pacienteUseCase.GetAntecedentesMedicos(c.Request.Context(), pacienteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Antecedentes medicos not found",
			"details": err.Error(),
		})
		return
	}

	if antecedentes == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Antecedentes medicos not found for this paciente",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Antecedentes medicos retrieved successfully",
		"data":    antecedentes,
	})
}

// Antecedentes Visuales

func (h *PacienteHandler) CreateOrUpdateAntecedentesVisuales(c *gin.Context) {
	idParam := c.Param("id")
	pacienteID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	var req entities.CreateAntecedentesVisualesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	antecedentes, err := h.pacienteUseCase.CreateOrUpdateAntecedentesVisuales(c.Request.Context(), pacienteID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create/update antecedentes visuales",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Antecedentes visuales saved successfully",
		"data":    antecedentes,
	})
}

func (h *PacienteHandler) GetAntecedentesVisuales(c *gin.Context) {
	idParam := c.Param("id")
	pacienteID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	antecedentes, err := h.pacienteUseCase.GetAntecedentesVisuales(c.Request.Context(), pacienteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Antecedentes visuales not found",
			"details": err.Error(),
		})
		return
	}

	if antecedentes == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Antecedentes visuales not found for this paciente",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Antecedentes visuales retrieved successfully",
		"data":    antecedentes,
	})
}

// Exámenes Visuales

func (h *PacienteHandler) CreateExamenVisual(c *gin.Context) {
	var req entities.CreateExamenVisualRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	examen, err := h.pacienteUseCase.CreateExamenVisual(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create examen visual",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Examen visual created successfully",
		"data":    examen,
	})
}

func (h *PacienteHandler) GetExamenVisual(c *gin.Context) {
	idParam := c.Param("examenId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid examen visual ID",
			"details": "Examen visual ID must be a number",
		})
		return
	}

	examen, err := h.pacienteUseCase.GetExamenVisual(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Examen visual not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Examen visual retrieved successfully",
		"data":    examen,
	})
}

func (h *PacienteHandler) GetExamenesVisualesByPaciente(c *gin.Context) {
	idParam := c.Param("id")
	pacienteID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid paciente ID",
			"details": "Paciente ID must be a number",
		})
		return
	}

	examenes, err := h.pacienteUseCase.GetExamenesVisualesByPaciente(c.Request.Context(), pacienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve examenes visuales",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Examenes visuales retrieved successfully",
		"data":    examenes,
	})
}

func (h *PacienteHandler) UpdateExamenVisual(c *gin.Context) {
	idParam := c.Param("examenId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid examen visual ID",
			"details": "Examen visual ID must be a number",
		})
		return
	}

	var req entities.UpdateExamenVisualRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	examen, err := h.pacienteUseCase.UpdateExamenVisual(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update examen visual",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Examen visual updated successfully",
		"data":    examen,
	})
}

func (h *PacienteHandler) DeleteExamenVisual(c *gin.Context) {
	idParam := c.Param("examenId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid examen visual ID",
			"details": "Examen visual ID must be a number",
		})
		return
	}

	err = h.pacienteUseCase.DeleteExamenVisual(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to delete examen visual",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PacienteHandler) CompararExamenes(c *gin.Context) {
	anteriorParam := c.Query("anterior")
	actualParam := c.Query("actual")

	anteriorID, err := strconv.ParseInt(anteriorParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid anterior examen ID",
			"details": "Anterior ID must be a number",
		})
		return
	}

	actualID, err := strconv.ParseInt(actualParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid actual examen ID",
			"details": "Actual ID must be a number",
		})
		return
	}

	diferencia, err := h.pacienteUseCase.CompararExamenes(c.Request.Context(), anteriorID, actualID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to compare examenes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comparison completed successfully",
		"data":    diferencia,
	})
}
