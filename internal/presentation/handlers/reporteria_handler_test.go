package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/domain/entities"
	"iris-api/internal/presentation/middleware"
)

func TestReporteriaHandler_GetPacientesActivos(t *testing.T) {
	tests := []struct {
		name       string
		setContext bool
		clientID   int64
		behavior   func(m *MockReporteriaUseCase)
		asserts    func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:       "success - gets active patients report",
			setContext: true,
			clientID:   100,
			behavior: func(m *MockReporteriaUseCase) {
				m.On("GetPacientesActivos", mock.Anything, int64(100)).Return(&entities.ReportePacientesActivos{
					Total: 10,
					Pacientes: []*entities.PacienteActivo{
						{ID: 1, Name: "Patient 1", TotalTurnos: 5},
						{ID: 2, Name: "Patient 2", TotalTurnos: 3},
					},
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Active patients report retrieved successfully", response["message"])
			},
		},
		{
			name:       "error - client id not in context",
			setContext: false,
			behavior:   func(m *MockReporteriaUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Client ID not found")
			},
		},
		{
			name:       "error - usecase fails",
			setContext: true,
			clientID:   100,
			behavior: func(m *MockReporteriaUseCase) {
				m.On("GetPacientesActivos", mock.Anything, int64(100)).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Failed to retrieve")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockReporteriaUseCase)
			tt.behavior(mockUseCase)

			handler := NewReporteriaHandler(mockUseCase)

			router := gin.New()
			router.GET("/reportes/pacientes-activos", func(c *gin.Context) {
				if tt.setContext {
					c.Set(middleware.ClientContextKey, tt.clientID)
				}
				c.Next()
			}, handler.GetPacientesActivos)

			req := httptest.NewRequest(http.MethodGet, "/reportes/pacientes-activos", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestReporteriaHandler_GetPacientesInactivos(t *testing.T) {
	tests := []struct {
		name       string
		setContext bool
		clientID   int64
		behavior   func(m *MockReporteriaUseCase)
		asserts    func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:       "success - gets inactive patients report",
			setContext: true,
			clientID:   100,
			behavior: func(m *MockReporteriaUseCase) {
				m.On("GetPacientesInactivos", mock.Anything, int64(100)).Return(&entities.ReportePacientesInactivos{
					Total: 5,
					Pacientes: []*entities.PacienteInactivo{
						{ID: 3, Name: "Patient 3", DiasInactivo: 90},
						{ID: 4, Name: "Patient 4", DiasInactivo: 120},
					},
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Inactive patients report retrieved successfully", response["message"])
			},
		},
		{
			name:       "error - client id not in context",
			setContext: false,
			behavior:   func(m *MockReporteriaUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Client ID not found")
			},
		},
		{
			name:       "error - usecase fails",
			setContext: true,
			clientID:   100,
			behavior: func(m *MockReporteriaUseCase) {
				m.On("GetPacientesInactivos", mock.Anything, int64(100)).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Failed to retrieve")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockReporteriaUseCase)
			tt.behavior(mockUseCase)

			handler := NewReporteriaHandler(mockUseCase)

			router := gin.New()
			router.GET("/reportes/pacientes-inactivos", func(c *gin.Context) {
				if tt.setContext {
					c.Set(middleware.ClientContextKey, tt.clientID)
				}
				c.Next()
			}, handler.GetPacientesInactivos)

			req := httptest.NewRequest(http.MethodGet, "/reportes/pacientes-inactivos", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

// MockReporteriaUseCase is a mock implementation of ReporteriaUseCase
type MockReporteriaUseCase struct {
	mock.Mock
}

func (m *MockReporteriaUseCase) GetPacientesActivos(ctx context.Context, clientID int64) (*entities.ReportePacientesActivos, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ReportePacientesActivos), args.Error(1)
}

func (m *MockReporteriaUseCase) GetPacientesInactivos(ctx context.Context, clientID int64) (*entities.ReportePacientesInactivos, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ReportePacientesInactivos), args.Error(1)
}
