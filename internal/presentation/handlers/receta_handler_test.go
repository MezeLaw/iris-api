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

	"iris-api/internal/application/usecases"
	"iris-api/internal/domain/entities"
)

func TestRecetaHandler_CreateReceta(t *testing.T) {
	tests := []struct {
		name     string
		body     interface{}
		behavior func(m *MockRecetaUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "success - creates receta",
			body: map[string]interface{}{
				"paciente_id": 1,
				"fecha":       time.Now().Format(time.RFC3339),
				"od_esfera":   -2.5,
				"oi_esfera":   -3.0,
				"tipo_lente":  "monofocal",
			},
			behavior: func(m *MockRecetaUseCase) {
				m.On("CreateReceta", mock.Anything, mock.Anything).Return(&entities.Receta{
					ID:         1,
					PacienteID: 1,
					ODEsfera:   -2.5,
					OIEsfera:   -3.0,
					TipoLente:  "monofocal",
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Receta created successfully", response["message"])
			},
		},
		{
			name: "error - invalid request body",
			body: "invalid json",
			behavior: func(m *MockRecetaUseCase) {},
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
				"paciente_id": 1,
				"fecha":       time.Now().Format(time.RFC3339),
				"od_esfera":   -2.5,
				"tipo_lente":  "invalid_type",
			},
			behavior: func(m *MockRecetaUseCase) {
				m.On("CreateReceta", mock.Anything, mock.Anything).Return(nil, errors.New("validation error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockRecetaUseCase)
			tt.behavior(mockUseCase)

			handler := NewRecetaHandler(mockUseCase)

			router := gin.New()
			router.POST("/recetas", handler.CreateReceta)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/recetas", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestRecetaHandler_GetRecetaByID(t *testing.T) {
	tests := []struct {
		name      string
		recetaID  string
		behavior  func(m *MockRecetaUseCase)
		asserts   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:     "success - gets receta by id",
			recetaID: "1",
			behavior: func(m *MockRecetaUseCase) {
				m.On("GetRecetaByID", mock.Anything, int64(1)).Return(&entities.Receta{
					ID:         1,
					PacienteID: 1,
					ODEsfera:   -2.5,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Receta retrieved successfully", response["message"])
			},
		},
		{
			name:     "error - invalid receta id",
			recetaID: "invalid",
			behavior: func(m *MockRecetaUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid receta ID")
			},
		},
		{
			name:     "error - receta not found",
			recetaID: "999",
			behavior: func(m *MockRecetaUseCase) {
				m.On("GetRecetaByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockRecetaUseCase)
			tt.behavior(mockUseCase)

			handler := NewRecetaHandler(mockUseCase)

			router := gin.New()
			router.GET("/recetas/:id", handler.GetRecetaByID)

			req := httptest.NewRequest(http.MethodGet, "/recetas/"+tt.recetaID, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestRecetaHandler_GetRecetas(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		behavior func(m *MockRecetaUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:  "success - gets recetas with default pagination",
			query: "",
			behavior: func(m *MockRecetaUseCase) {
				m.On("GetRecetas", mock.Anything, 10, 0).Return(&usecases.GetRecetasResponse{
					Recetas: []*entities.Receta{{ID: 1, PacienteID: 1}},
					Total:   1,
					HasMore: false,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - gets recetas with custom pagination",
			query: "?limit=20&offset=10",
			behavior: func(m *MockRecetaUseCase) {
				m.On("GetRecetas", mock.Anything, 20, 10).Return(&usecases.GetRecetasResponse{
					Recetas: []*entities.Receta{},
					Total:   0,
					HasMore: false,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "error - usecase fails",
			query: "",
			behavior: func(m *MockRecetaUseCase) {
				m.On("GetRecetas", mock.Anything, 10, 0).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockRecetaUseCase)
			tt.behavior(mockUseCase)

			handler := NewRecetaHandler(mockUseCase)

			router := gin.New()
			router.GET("/recetas", handler.GetRecetas)

			req := httptest.NewRequest(http.MethodGet, "/recetas"+tt.query, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestRecetaHandler_DeleteReceta(t *testing.T) {
	tests := []struct {
		name     string
		recetaID string
		behavior func(m *MockRecetaUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:     "success - deletes receta",
			recetaID: "1",
			behavior: func(m *MockRecetaUseCase) {
				m.On("DeleteReceta", mock.Anything, int64(1)).Return(nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, resp.Code)
			},
		},
		{
			name:     "error - invalid receta id",
			recetaID: "invalid",
			behavior: func(m *MockRecetaUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:     "error - receta not found",
			recetaID: "999",
			behavior: func(m *MockRecetaUseCase) {
				m.On("DeleteReceta", mock.Anything, int64(999)).Return(errors.New("not found"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockRecetaUseCase)
			tt.behavior(mockUseCase)

			handler := NewRecetaHandler(mockUseCase)

			router := gin.New()
			router.DELETE("/recetas/:id", handler.DeleteReceta)

			req := httptest.NewRequest(http.MethodDelete, "/recetas/"+tt.recetaID, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestRecetaHandler_CheckDioptriasChange(t *testing.T) {
	tests := []struct {
		name        string
		pacienteID  string
		behavior    func(m *MockRecetaUseCase)
		asserts     func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:       "success - checks dioptrías change",
			pacienteID: "1",
			behavior: func(m *MockRecetaUseCase) {
				m.On("CheckDioptriasChange", mock.Anything, int64(1)).Return(&entities.DioptriasChange{
					ODChange: 0.5,
					OIChange: 0.25,
					HasAlert: true,
					Message:  "Cambio significativo detectado",
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Dioptrias change checked successfully", response["message"])
			},
		},
		{
			name:       "error - invalid paciente id",
			pacienteID: "invalid",
			behavior:   func(m *MockRecetaUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockRecetaUseCase)
			tt.behavior(mockUseCase)

			handler := NewRecetaHandler(mockUseCase)

			router := gin.New()
			router.GET("/recetas/paciente/:paciente_id/check-change", handler.CheckDioptriasChange)

			req := httptest.NewRequest(http.MethodGet, "/recetas/paciente/"+tt.pacienteID+"/check-change", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

// MockRecetaUseCase is a mock implementation of RecetaUseCase
type MockRecetaUseCase struct {
	mock.Mock
}

func (m *MockRecetaUseCase) CreateReceta(ctx context.Context, req *entities.CreateRecetaRequest) (*entities.Receta, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Receta), args.Error(1)
}

func (m *MockRecetaUseCase) GetRecetaByID(ctx context.Context, id int64) (*entities.Receta, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Receta), args.Error(1)
}

func (m *MockRecetaUseCase) GetRecetas(ctx context.Context, limit, offset int) (*usecases.GetRecetasResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.GetRecetasResponse), args.Error(1)
}

func (m *MockRecetaUseCase) GetRecetasByPacienteID(ctx context.Context, pacienteID int64, limit, offset int) (*usecases.GetRecetasResponse, error) {
	args := m.Called(ctx, pacienteID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.GetRecetasResponse), args.Error(1)
}

func (m *MockRecetaUseCase) GetHistorialByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error) {
	args := m.Called(ctx, pacienteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Receta), args.Error(1)
}

func (m *MockRecetaUseCase) CheckDioptriasChange(ctx context.Context, pacienteID int64) (*entities.DioptriasChange, error) {
	args := m.Called(ctx, pacienteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.DioptriasChange), args.Error(1)
}

func (m *MockRecetaUseCase) UpdateReceta(ctx context.Context, id int64, req *entities.UpdateRecetaRequest) (*entities.Receta, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Receta), args.Error(1)
}

func (m *MockRecetaUseCase) DeleteReceta(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
