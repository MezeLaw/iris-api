package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/domain/entities"
)

func TestTurnoHandler_CreateTurno(t *testing.T) {
	tests := []struct {
		name     string
		body     interface{}
		behavior func(m *MockTurnoUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "success - creates turno",
			body: map[string]interface{}{
				"paciente_id":      1,
				"contactologo_id":  2,
				"fecha_hora":       time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"duracion_minutos": 30,
				"tipo_servicio":    "consulta",
			},
			behavior: func(m *MockTurnoUseCase) {
				m.On("CreateTurno", mock.Anything, mock.Anything).Return(&entities.Turno{
					ID:              1,
					PacienteID:      1,
					ContactologoID:  2,
					DuracionMinutos: 30,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Turno created successfully", response["message"])
			},
		},
		{
			name:     "error - invalid request body",
			body:     "invalid json",
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid request body")
			},
		},
		{
			name: "error - usecase fails",
			body: map[string]interface{}{
				"paciente_id":     1,
				"contactologo_id": 2,
				"fecha_hora":      time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			},
			behavior: func(m *MockTurnoUseCase) {
				m.On("CreateTurno", mock.Anything, mock.Anything).Return(nil, errors.New("fecha en el pasado"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockTurnoUseCase)
			tt.behavior(mockUseCase)

			handler := NewTurnoHandler(mockUseCase)

			router := gin.New()
			router.POST("/turnos", handler.CreateTurno)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/turnos", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestTurnoHandler_GetTurnoByID(t *testing.T) {
	tests := []struct {
		name     string
		turnoID  string
		behavior func(m *MockTurnoUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:    "success - gets turno by id",
			turnoID: "1",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnoByID", mock.Anything, int64(1)).Return(&entities.TurnoConDetalles{
					Turno: entities.Turno{
						ID:              1,
						DuracionMinutos: 30,
					},
					PacienteNombre: "Test Patient",
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Turno retrieved successfully", response["message"])
			},
		},
		{
			name:     "error - invalid turno id",
			turnoID:  "invalid",
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid turno ID")
			},
		},
		{
			name:    "error - turno not found",
			turnoID: "999",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnoByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockTurnoUseCase)
			tt.behavior(mockUseCase)

			handler := NewTurnoHandler(mockUseCase)

			router := gin.New()
			router.GET("/turnos/:id", handler.GetTurnoByID)

			req := httptest.NewRequest(http.MethodGet, "/turnos/"+tt.turnoID, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestTurnoHandler_GetTurnos(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		behavior func(m *MockTurnoUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:  "success - gets turnos with default pagination",
			query: "",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnos", mock.Anything, mock.Anything).Return(map[string]interface{}{
					"turnos": []entities.TurnoConDetalles{},
					"total":  0,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - gets turnos with custom pagination",
			query: "?limit=20&offset=10",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnos", mock.Anything, mock.Anything).Return(map[string]interface{}{
					"turnos": []entities.TurnoConDetalles{},
					"total":  0,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - with paciente_id filter",
			query: "?paciente_id=1",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnos", mock.Anything, mock.Anything).Return(map[string]interface{}{
					"turnos": []entities.TurnoConDetalles{},
					"total":  0,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - with contactologo_id filter",
			query: "?contactologo_id=2",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnos", mock.Anything, mock.Anything).Return(map[string]interface{}{
					"turnos": []entities.TurnoConDetalles{},
					"total":  0,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - with tipo_servicio filter",
			query: "?tipo_servicio=consulta",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnos", mock.Anything, mock.Anything).Return(map[string]interface{}{
					"turnos": []entities.TurnoConDetalles{},
					"total":  0,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - with estado filter",
			query: "?estado=confirmado",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnos", mock.Anything, mock.Anything).Return(map[string]interface{}{
					"turnos": []entities.TurnoConDetalles{},
					"total":  0,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - with fecha_desde filter",
			query: "?fecha_desde=2025-01-01T00:00:00Z",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnos", mock.Anything, mock.Anything).Return(map[string]interface{}{
					"turnos": []entities.TurnoConDetalles{},
					"total":  0,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - with fecha_hasta filter",
			query: "?fecha_hasta=2025-12-31T23:59:59Z",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnos", mock.Anything, mock.Anything).Return(map[string]interface{}{
					"turnos": []entities.TurnoConDetalles{},
					"total":  0,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "error - usecase fails",
			query: "",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetTurnos", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockTurnoUseCase)
			tt.behavior(mockUseCase)

			handler := NewTurnoHandler(mockUseCase)

			router := gin.New()
			router.GET("/turnos", handler.GetTurnos)

			req := httptest.NewRequest(http.MethodGet, "/turnos"+tt.query, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestTurnoHandler_CancelTurno(t *testing.T) {
	tests := []struct {
		name     string
		turnoID  string
		body     interface{}
		behavior func(m *MockTurnoUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:    "success - cancels turno",
			turnoID: "1",
			body: map[string]interface{}{
				"motivo": "Paciente no puede asistir",
			},
			behavior: func(m *MockTurnoUseCase) {
				m.On("CancelTurno", mock.Anything, int64(1), "Paciente no puede asistir").Return(nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Turno cancelled successfully", response["message"])
			},
		},
		{
			name:     "error - invalid turno id",
			turnoID:  "invalid",
			body:     map[string]interface{}{},
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:     "error - invalid request body",
			turnoID:  "1",
			body:     "invalid json",
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:    "error - usecase fails",
			turnoID: "1",
			body: map[string]interface{}{
				"motivo": "Test motivo",
			},
			behavior: func(m *MockTurnoUseCase) {
				m.On("CancelTurno", mock.Anything, int64(1), "Test motivo").Return(errors.New("cancel failed"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockTurnoUseCase)
			tt.behavior(mockUseCase)

			handler := NewTurnoHandler(mockUseCase)

			router := gin.New()
			router.POST("/turnos/:id/cancel", handler.CancelTurno)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/turnos/"+tt.turnoID+"/cancel", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestTurnoHandler_DeleteTurno(t *testing.T) {
	tests := []struct {
		name     string
		turnoID  string
		behavior func(m *MockTurnoUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:    "success - deletes turno",
			turnoID: "1",
			behavior: func(m *MockTurnoUseCase) {
				m.On("DeleteTurno", mock.Anything, int64(1)).Return(nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, resp.Code)
			},
		},
		{
			name:     "error - invalid turno id",
			turnoID:  "invalid",
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:    "error - turno not found",
			turnoID: "999",
			behavior: func(m *MockTurnoUseCase) {
				m.On("DeleteTurno", mock.Anything, int64(999)).Return(errors.New("not found"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockTurnoUseCase)
			tt.behavior(mockUseCase)

			handler := NewTurnoHandler(mockUseCase)

			router := gin.New()
			router.DELETE("/turnos/:id", handler.DeleteTurno)

			req := httptest.NewRequest(http.MethodDelete, "/turnos/"+tt.turnoID, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestTurnoHandler_GetProximosTurnosAlert(t *testing.T) {
	tests := []struct {
		name     string
		behavior func(m *MockTurnoUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "success - gets alerts",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetProximosTurnosAlert", mock.Anything).Return(&entities.ProximosTurnosAlert{
					Count24h:          2,
					CountSinConfirmar: 1,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Alerts retrieved successfully", response["message"])
			},
		},
		{
			name: "error - usecase fails",
			behavior: func(m *MockTurnoUseCase) {
				m.On("GetProximosTurnosAlert", mock.Anything).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockTurnoUseCase)
			tt.behavior(mockUseCase)

			handler := NewTurnoHandler(mockUseCase)

			router := gin.New()
			router.GET("/turnos/alerts", handler.GetProximosTurnosAlert)

			req := httptest.NewRequest(http.MethodGet, "/turnos/alerts", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestTurnoHandler_UpdateTurno(t *testing.T) {
	tests := []struct {
		name     string
		turnoID  string
		body     interface{}
		behavior func(m *MockTurnoUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:    "success - updates turno",
			turnoID: "1",
			body: map[string]interface{}{
				"fecha_hora":       time.Now().Add(48 * time.Hour).Format(time.RFC3339),
				"duracion_minutos": 60,
			},
			behavior: func(m *MockTurnoUseCase) {
				m.On("UpdateTurno", mock.Anything, int64(1), mock.Anything).Return(&entities.Turno{
					ID:              1,
					DuracionMinutos: 60,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Turno updated successfully", response["message"])
			},
		},
		{
			name:     "error - invalid turno id",
			turnoID:  "invalid",
			body:     map[string]interface{}{},
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid turno ID")
			},
		},
		{
			name:     "error - invalid request body",
			turnoID:  "1",
			body:     "invalid json",
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid request body")
			},
		},
		{
			name:    "error - turno not found",
			turnoID: "999",
			body: map[string]interface{}{
				"fecha_hora": time.Now().Add(48 * time.Hour).Format(time.RFC3339),
			},
			behavior: func(m *MockTurnoUseCase) {
				m.On("UpdateTurno", mock.Anything, int64(999), mock.Anything).Return(nil, errors.New("turno not found"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:    "error - usecase validation fails",
			turnoID: "1",
			body: map[string]interface{}{
				"fecha_hora": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			},
			behavior: func(m *MockTurnoUseCase) {
				m.On("UpdateTurno", mock.Anything, int64(1), mock.Anything).Return(nil, errors.New("validation error: fecha en el pasado"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockTurnoUseCase)
			tt.behavior(mockUseCase)

			handler := NewTurnoHandler(mockUseCase)

			router := gin.New()
			router.PUT("/turnos/:id", handler.UpdateTurno)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/turnos/"+tt.turnoID, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestTurnoHandler_GetTurnosByDia(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		behavior func(m *MockTurnoUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:  "success - gets turnos by day",
			query: "?fecha=2025-12-25",
			behavior: func(m *MockTurnoUseCase) {
				fecha, _ := time.Parse("2006-01-02", "2025-12-25")
				m.On("GetTurnosByDia", mock.Anything, fecha, (*int64)(nil)).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 1}},
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Turnos retrieved successfully", response["message"])
			},
		},
		{
			name:  "success - with contactologo filter",
			query: "?fecha=2025-12-25&contactologo_id=1",
			behavior: func(m *MockTurnoUseCase) {
				fecha, _ := time.Parse("2006-01-02", "2025-12-25")
				contactologoID := int64(1)
				m.On("GetTurnosByDia", mock.Anything, fecha, &contactologoID).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 1}},
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:     "error - missing fecha parameter",
			query:    "",
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Missing fecha parameter")
			},
		},
		{
			name:     "error - invalid fecha format",
			query:    "?fecha=invalid-date",
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid fecha format")
			},
		},
		{
			name:  "error - usecase fails",
			query: "?fecha=2025-12-25",
			behavior: func(m *MockTurnoUseCase) {
				fecha, _ := time.Parse("2006-01-02", "2025-12-25")
				m.On("GetTurnosByDia", mock.Anything, fecha, (*int64)(nil)).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockTurnoUseCase)
			tt.behavior(mockUseCase)

			handler := NewTurnoHandler(mockUseCase)

			router := gin.New()
			router.GET("/turnos/dia", handler.GetTurnosByDia)

			req := httptest.NewRequest(http.MethodGet, "/turnos/dia"+tt.query, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestTurnoHandler_GetTurnosBySemana(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		behavior func(m *MockTurnoUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:  "success - gets turnos by week",
			query: "?fecha_inicio=2025-12-22",
			behavior: func(m *MockTurnoUseCase) {
				fecha, _ := time.Parse("2006-01-02", "2025-12-22")
				m.On("GetTurnosBySemana", mock.Anything, fecha, (*int64)(nil)).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 1}},
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Turnos retrieved successfully", response["message"])
			},
		},
		{
			name:  "success - with contactologo filter",
			query: "?fecha_inicio=2025-12-22&contactologo_id=1",
			behavior: func(m *MockTurnoUseCase) {
				fecha, _ := time.Parse("2006-01-02", "2025-12-22")
				contactologoID := int64(1)
				m.On("GetTurnosBySemana", mock.Anything, fecha, &contactologoID).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 1}},
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:     "error - missing fecha_inicio parameter",
			query:    "",
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Missing fecha_inicio parameter")
			},
		},
		{
			name:     "error - invalid fecha_inicio format",
			query:    "?fecha_inicio=invalid-date",
			behavior: func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid fecha_inicio format")
			},
		},
		{
			name:  "error - usecase fails",
			query: "?fecha_inicio=2025-12-22",
			behavior: func(m *MockTurnoUseCase) {
				fecha, _ := time.Parse("2006-01-02", "2025-12-22")
				m.On("GetTurnosBySemana", mock.Anything, fecha, (*int64)(nil)).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockTurnoUseCase)
			tt.behavior(mockUseCase)

			handler := NewTurnoHandler(mockUseCase)

			router := gin.New()
			router.GET("/turnos/semana", handler.GetTurnosBySemana)

			req := httptest.NewRequest(http.MethodGet, "/turnos/semana"+tt.query, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestTurnoHandler_GetTurnosByProfesional(t *testing.T) {
	tests := []struct {
		name           string
		contactologoID string
		query          string
		behavior       func(m *MockTurnoUseCase)
		asserts        func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:           "success - gets turnos by profesional",
			contactologoID: "1",
			query:          "?fecha_desde=2025-12-01&fecha_hasta=2025-12-31",
			behavior: func(m *MockTurnoUseCase) {
				fechaDesde, _ := time.Parse("2006-01-02", "2025-12-01")
				fechaHasta, _ := time.Parse("2006-01-02", "2025-12-31")
				m.On("GetTurnosByProfesional", mock.Anything, int64(1), fechaDesde, fechaHasta).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 1}},
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Turnos retrieved successfully", response["message"])
			},
		},
		{
			name:           "error - invalid contactologo id",
			contactologoID: "invalid",
			query:          "?fecha_desde=2025-12-01&fecha_hasta=2025-12-31",
			behavior:       func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid contactologo ID")
			},
		},
		{
			name:           "error - missing fecha_desde",
			contactologoID: "1",
			query:          "?fecha_hasta=2025-12-31",
			behavior:       func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Missing date parameters")
			},
		},
		{
			name:           "error - missing fecha_hasta",
			contactologoID: "1",
			query:          "?fecha_desde=2025-12-01",
			behavior:       func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Missing date parameters")
			},
		},
		{
			name:           "error - invalid fecha format",
			contactologoID: "1",
			query:          "?fecha_desde=invalid&fecha_hasta=2025-12-31",
			behavior:       func(m *MockTurnoUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid fecha_desde format")
			},
		},
		{
			name:           "error - usecase fails",
			contactologoID: "1",
			query:          "?fecha_desde=2025-12-01&fecha_hasta=2025-12-31",
			behavior: func(m *MockTurnoUseCase) {
				fechaDesde, _ := time.Parse("2006-01-02", "2025-12-01")
				fechaHasta, _ := time.Parse("2006-01-02", "2025-12-31")
				m.On("GetTurnosByProfesional", mock.Anything, int64(1), fechaDesde, fechaHasta).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockTurnoUseCase)
			tt.behavior(mockUseCase)

			handler := NewTurnoHandler(mockUseCase)

			router := gin.New()
			router.GET("/turnos/profesional/:contactologo_id", handler.GetTurnosByProfesional)

			req := httptest.NewRequest(http.MethodGet, "/turnos/profesional/"+tt.contactologoID+tt.query, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

// MockTurnoUseCase is a mock implementation of TurnoUseCase
type MockTurnoUseCase struct {
	mock.Mock
}

func (m *MockTurnoUseCase) CreateTurno(ctx context.Context, req *entities.CreateTurnoRequest) (*entities.Turno, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Turno), args.Error(1)
}

func (m *MockTurnoUseCase) GetTurnoByID(ctx context.Context, id int64) (*entities.TurnoConDetalles, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoUseCase) GetTurnos(ctx context.Context, filter *entities.TurnoFilter) (map[string]interface{}, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockTurnoUseCase) UpdateTurno(ctx context.Context, id int64, req *entities.UpdateTurnoRequest) (*entities.Turno, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Turno), args.Error(1)
}

func (m *MockTurnoUseCase) CancelTurno(ctx context.Context, id int64, motivo string) error {
	args := m.Called(ctx, id, motivo)
	return args.Error(0)
}

func (m *MockTurnoUseCase) DeleteTurno(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTurnoUseCase) GetTurnosByDia(ctx context.Context, fecha time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error) {
	args := m.Called(ctx, fecha, contactologoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoUseCase) GetTurnosBySemana(ctx context.Context, fechaInicio time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error) {
	args := m.Called(ctx, fechaInicio, contactologoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoUseCase) GetTurnosByProfesional(ctx context.Context, contactologoID int64, fechaDesde, fechaHasta time.Time) ([]*entities.TurnoConDetalles, error) {
	args := m.Called(ctx, contactologoID, fechaDesde, fechaHasta)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoUseCase) GetProximosTurnosAlert(ctx context.Context) (*entities.ProximosTurnosAlert, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProximosTurnosAlert), args.Error(1)
}
