package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/domain/entities"
)

func TestRecetaUseCase_CreateReceta(t *testing.T) {
	tests := []struct {
		name     string
		req      *entities.CreateRecetaRequest
		behavior func(m *MockRecetaService)
		asserts  func(t *testing.T, receta *entities.Receta, err error)
	}{
		{
			name: "success - creates receta",
			req: &entities.CreateRecetaRequest{
				PacienteID: 1,
				Fecha:      time.Now(),
				ODEsfera:   -2.5,
				OIEsfera:   -3.0,
				TipoLente:  "monofocal",
			},
			behavior: func(m *MockRecetaService) {
				m.On("CreateReceta", mock.Anything, mock.Anything).Return(&entities.Receta{
					ID:         1,
					PacienteID: 1,
					ODEsfera:   -2.5,
					OIEsfera:   -3.0,
				}, nil)
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, receta)
				assert.Equal(t, int64(1), receta.ID)
			},
		},
		{
			name: "error - service fails",
			req:  &entities.CreateRecetaRequest{},
			behavior: func(m *MockRecetaService) {
				m.On("CreateReceta", mock.Anything, mock.Anything).Return(nil, errors.New("validation error"))
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockRecetaService)
			tt.behavior(mockService)

			useCase := NewRecetaUseCase(mockService)
			receta, err := useCase.CreateReceta(context.Background(), tt.req)

			tt.asserts(t, receta, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestRecetaUseCase_GetRecetaByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(m *MockRecetaService)
		asserts  func(t *testing.T, receta *entities.Receta, err error)
	}{
		{
			name: "success - gets receta by id",
			id:   1,
			behavior: func(m *MockRecetaService) {
				m.On("GetRecetaByID", mock.Anything, int64(1)).Return(&entities.Receta{
					ID:         1,
					PacienteID: 1,
					ODEsfera:   -2.5,
				}, nil)
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, receta)
				assert.Equal(t, int64(1), receta.ID)
			},
		},
		{
			name: "error - not found",
			id:   999,
			behavior: func(m *MockRecetaService) {
				m.On("GetRecetaByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockRecetaService)
			tt.behavior(mockService)

			useCase := NewRecetaUseCase(mockService)
			receta, err := useCase.GetRecetaByID(context.Background(), tt.id)

			tt.asserts(t, receta, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestRecetaUseCase_GetRecetas(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		offset   int
		behavior func(m *MockRecetaService)
		asserts  func(t *testing.T, resp *GetRecetasResponse, err error)
	}{
		{
			name:   "success - returns recetas with has_more true",
			limit:  10,
			offset: 0,
			behavior: func(m *MockRecetaService) {
				m.On("GetRecetas", mock.Anything, 10, 0).Return([]*entities.Receta{
					{ID: 1, PacienteID: 1},
					{ID: 2, PacienteID: 2},
				}, 20, nil)
			},
			asserts: func(t *testing.T, resp *GetRecetasResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Recetas, 2)
				assert.Equal(t, 20, resp.Total)
				assert.True(t, resp.HasMore)
			},
		},
		{
			name:   "success - returns recetas with has_more false",
			limit:  10,
			offset: 15,
			behavior: func(m *MockRecetaService) {
				m.On("GetRecetas", mock.Anything, 10, 15).Return([]*entities.Receta{
					{ID: 16, PacienteID: 1},
				}, 20, nil)
			},
			asserts: func(t *testing.T, resp *GetRecetasResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.False(t, resp.HasMore)
			},
		},
		{
			name:   "error - service fails",
			limit:  10,
			offset: 0,
			behavior: func(m *MockRecetaService) {
				m.On("GetRecetas", mock.Anything, 10, 0).Return(nil, 0, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *GetRecetasResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockRecetaService)
			tt.behavior(mockService)

			useCase := NewRecetaUseCase(mockService)
			resp, err := useCase.GetRecetas(context.Background(), tt.limit, tt.offset)

			tt.asserts(t, resp, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestRecetaUseCase_CheckDioptriasChange(t *testing.T) {
	tests := []struct {
		name       string
		pacienteID int64
		behavior   func(m *MockRecetaService)
		asserts    func(t *testing.T, change *entities.DioptriasChange, err error)
	}{
		{
			name:       "success - returns dioptrías change",
			pacienteID: 1,
			behavior: func(m *MockRecetaService) {
				m.On("CheckDioptriasChange", mock.Anything, int64(1)).Return(&entities.DioptriasChange{
					ODChange: 0.5,
					OIChange: 0.25,
					HasAlert: true,
					Message:  "Cambio significativo detectado",
				}, nil)
			},
			asserts: func(t *testing.T, change *entities.DioptriasChange, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, change)
				assert.True(t, change.HasAlert)
			},
		},
		{
			name:       "error - service fails",
			pacienteID: 999,
			behavior: func(m *MockRecetaService) {
				m.On("CheckDioptriasChange", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, change *entities.DioptriasChange, err error) {
				assert.Error(t, err)
				assert.Nil(t, change)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockRecetaService)
			tt.behavior(mockService)

			useCase := NewRecetaUseCase(mockService)
			change, err := useCase.CheckDioptriasChange(context.Background(), tt.pacienteID)

			tt.asserts(t, change, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestRecetaUseCase_DeleteReceta(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(m *MockRecetaService)
		asserts  func(t *testing.T, err error)
	}{
		{
			name: "success - deletes receta",
			id:   1,
			behavior: func(m *MockRecetaService) {
				m.On("DeleteReceta", mock.Anything, int64(1)).Return(nil)
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error - service fails",
			id:   999,
			behavior: func(m *MockRecetaService) {
				m.On("DeleteReceta", mock.Anything, int64(999)).Return(errors.New("not found"))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockRecetaService)
			tt.behavior(mockService)

			useCase := NewRecetaUseCase(mockService)
			err := useCase.DeleteReceta(context.Background(), tt.id)

			tt.asserts(t, err)
			mockService.AssertExpectations(t)
		})
	}
}

// MockRecetaService is a mock implementation of RecetaService
type MockRecetaService struct {
	mock.Mock
}

func (m *MockRecetaService) CreateReceta(ctx context.Context, req *entities.CreateRecetaRequest) (*entities.Receta, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Receta), args.Error(1)
}

func (m *MockRecetaService) GetRecetaByID(ctx context.Context, id int64) (*entities.Receta, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Receta), args.Error(1)
}

func (m *MockRecetaService) GetRecetas(ctx context.Context, limit, offset int) ([]*entities.Receta, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.Receta), args.Int(1), args.Error(2)
}

func (m *MockRecetaService) GetRecetasByPacienteID(ctx context.Context, pacienteID int64, limit, offset int) ([]*entities.Receta, int, error) {
	args := m.Called(ctx, pacienteID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.Receta), args.Int(1), args.Error(2)
}

func (m *MockRecetaService) GetHistorialByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error) {
	args := m.Called(ctx, pacienteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Receta), args.Error(1)
}

func (m *MockRecetaService) CheckDioptriasChange(ctx context.Context, pacienteID int64) (*entities.DioptriasChange, error) {
	args := m.Called(ctx, pacienteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.DioptriasChange), args.Error(1)
}

func (m *MockRecetaService) UpdateReceta(ctx context.Context, id int64, req *entities.UpdateRecetaRequest) (*entities.Receta, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Receta), args.Error(1)
}

func (m *MockRecetaService) DeleteReceta(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
