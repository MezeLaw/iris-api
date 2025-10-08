package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/domain/entities"
)

func TestRecetaService_CreateReceta(t *testing.T) {
	tests := []struct {
		name     string
		req      *entities.CreateRecetaRequest
		behavior func(m *MockRecetaRepo)
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
			behavior: func(m *MockRecetaRepo) {
				m.On("Create", mock.Anything, mock.Anything).Return(&entities.Receta{
					ID:         1,
					PacienteID: 1,
					ODEsfera:   -2.5,
					OIEsfera:   -3.0,
					TipoLente:  "monofocal",
				}, nil)
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, receta)
				assert.Equal(t, int64(1), receta.ID)
			},
		},
		{
			name: "error - invalid paciente_id",
			req: &entities.CreateRecetaRequest{
				PacienteID: 0,
				Fecha:      time.Now(),
				ODEsfera:   -2.5,
				OIEsfera:   -3.0,
				TipoLente:  "monofocal",
			},
			behavior: func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
				assert.Contains(t, err.Error(), "paciente_id")
			},
		},
		{
			name: "error - invalid od_esfera range",
			req: &entities.CreateRecetaRequest{
				PacienteID: 1,
				Fecha:      time.Now(),
				ODEsfera:   -25.0,
				OIEsfera:   -3.0,
				TipoLente:  "monofocal",
			},
			behavior: func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
				assert.Contains(t, err.Error(), "od_esfera")
			},
		},
		{
			name: "error - invalid tipo_lente",
			req: &entities.CreateRecetaRequest{
				PacienteID: 1,
				Fecha:      time.Now(),
				ODEsfera:   -2.5,
				OIEsfera:   -3.0,
				TipoLente:  "invalid",
			},
			behavior: func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
				assert.Contains(t, err.Error(), "tipo_lente")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRecetaRepo)
			tt.behavior(mockRepo)

			service := NewRecetaService(mockRepo)
			receta, err := service.CreateReceta(context.Background(), tt.req)

			tt.asserts(t, receta, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRecetaService_GetRecetaByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(m *MockRecetaRepo)
		asserts  func(t *testing.T, receta *entities.Receta, err error)
	}{
		{
			name: "success - gets receta by id",
			id:   1,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetByID", mock.Anything, int64(1)).Return(&entities.Receta{
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
			name:     "error - invalid id",
			id:       0,
			behavior: func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
				assert.Contains(t, err.Error(), "invalid")
			},
		},
		{
			name: "error - not found",
			id:   999,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRecetaRepo)
			tt.behavior(mockRepo)

			service := NewRecetaService(mockRepo)
			receta, err := service.GetRecetaByID(context.Background(), tt.id)

			tt.asserts(t, receta, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRecetaService_GetRecetas(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		offset   int
		behavior func(m *MockRecetaRepo)
		asserts  func(t *testing.T, recetas []*entities.Receta, total int, err error)
	}{
		{
			name:   "success - returns recetas",
			limit:  10,
			offset: 0,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetAll", mock.Anything, 10, 0).Return([]*entities.Receta{
					{ID: 1, PacienteID: 1},
					{ID: 2, PacienteID: 2},
				}, nil)
				m.On("Count", mock.Anything).Return(2, nil)
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, total int, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, recetas)
				assert.Len(t, recetas, 2)
				assert.Equal(t, 2, total)
			},
		},
		{
			name:   "success - applies default limit",
			limit:  0,
			offset: 0,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetAll", mock.Anything, 10, 0).Return([]*entities.Receta{}, nil)
				m.On("Count", mock.Anything).Return(0, nil)
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, total int, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, recetas)
			},
		},
		{
			name:   "success - caps limit at 100",
			limit:  200,
			offset: 0,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetAll", mock.Anything, 100, 0).Return([]*entities.Receta{}, nil)
				m.On("Count", mock.Anything).Return(0, nil)
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, total int, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, recetas)
			},
		},
		{
			name:   "error - repository fails",
			limit:  10,
			offset: 0,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetAll", mock.Anything, 10, 0).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, total int, err error) {
				assert.Error(t, err)
				assert.Nil(t, recetas)
				assert.Equal(t, 0, total)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRecetaRepo)
			tt.behavior(mockRepo)

			service := NewRecetaService(mockRepo)
			recetas, total, err := service.GetRecetas(context.Background(), tt.limit, tt.offset)

			tt.asserts(t, recetas, total, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRecetaService_CheckDioptriasChange(t *testing.T) {
	tests := []struct {
		name       string
		pacienteID int64
		behavior   func(m *MockRecetaRepo)
		asserts    func(t *testing.T, change *entities.DioptriasChange, err error)
	}{
		{
			name:       "success - detects significant change",
			pacienteID: 1,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetLastTwoByPacienteID", mock.Anything, int64(1)).Return([]*entities.Receta{
					{ID: 2, PacienteID: 1, ODEsfera: -3.0, OIEsfera: -2.5},
					{ID: 1, PacienteID: 1, ODEsfera: -2.5, OIEsfera: -2.0},
				}, nil)
			},
			asserts: func(t *testing.T, change *entities.DioptriasChange, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, change)
				assert.True(t, change.HasAlert)
				assert.Equal(t, 0.5, change.ODChange)
				assert.Equal(t, 0.5, change.OIChange)
			},
		},
		{
			name:       "success - no significant change",
			pacienteID: 1,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetLastTwoByPacienteID", mock.Anything, int64(1)).Return([]*entities.Receta{
					{ID: 2, PacienteID: 1, ODEsfera: -2.5, OIEsfera: -2.0},
					{ID: 1, PacienteID: 1, ODEsfera: -2.4, OIEsfera: -1.9},
				}, nil)
			},
			asserts: func(t *testing.T, change *entities.DioptriasChange, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, change)
				assert.False(t, change.HasAlert)
			},
		},
		{
			name:       "success - insufficient recetas",
			pacienteID: 1,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetLastTwoByPacienteID", mock.Anything, int64(1)).Return([]*entities.Receta{
					{ID: 1, PacienteID: 1, ODEsfera: -2.5, OIEsfera: -2.0},
				}, nil)
			},
			asserts: func(t *testing.T, change *entities.DioptriasChange, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, change)
				assert.False(t, change.HasAlert)
				assert.Contains(t, change.Message, "suficientes")
			},
		},
		{
			name:       "error - invalid paciente id",
			pacienteID: 0,
			behavior:   func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, change *entities.DioptriasChange, err error) {
				assert.Error(t, err)
				assert.Nil(t, change)
				assert.Contains(t, err.Error(), "invalid")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRecetaRepo)
			tt.behavior(mockRepo)

			service := NewRecetaService(mockRepo)
			change, err := service.CheckDioptriasChange(context.Background(), tt.pacienteID)

			tt.asserts(t, change, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRecetaService_GetRecetasByPacienteID(t *testing.T) {
	tests := []struct {
		name       string
		pacienteID int64
		limit      int
		offset     int
		behavior   func(m *MockRecetaRepo)
		asserts    func(t *testing.T, recetas []*entities.Receta, total int, err error)
	}{
		{
			name:       "success - returns recetas by paciente",
			pacienteID: 1,
			limit:      10,
			offset:     0,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetByPacienteID", mock.Anything, int64(1), 10, 0).Return([]*entities.Receta{
					{ID: 1, PacienteID: 1},
					{ID: 2, PacienteID: 1},
				}, nil)
				m.On("CountByPacienteID", mock.Anything, int64(1)).Return(2, nil)
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, total int, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, recetas)
				assert.Len(t, recetas, 2)
				assert.Equal(t, 2, total)
			},
		},
		{
			name:       "error - invalid paciente id",
			pacienteID: 0,
			limit:      10,
			offset:     0,
			behavior:   func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, recetas []*entities.Receta, total int, err error) {
				assert.Error(t, err)
				assert.Nil(t, recetas)
				assert.Contains(t, err.Error(), "invalid")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRecetaRepo)
			tt.behavior(mockRepo)

			service := NewRecetaService(mockRepo)
			recetas, total, err := service.GetRecetasByPacienteID(context.Background(), tt.pacienteID, tt.limit, tt.offset)

			tt.asserts(t, recetas, total, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRecetaService_GetHistorialByPacienteID(t *testing.T) {
	tests := []struct {
		name       string
		pacienteID int64
		behavior   func(m *MockRecetaRepo)
		asserts    func(t *testing.T, recetas []*entities.Receta, err error)
	}{
		{
			name:       "success - returns historial",
			pacienteID: 1,
			behavior: func(m *MockRecetaRepo) {
				m.On("GetHistorialByPacienteID", mock.Anything, int64(1)).Return([]*entities.Receta{
					{ID: 1, PacienteID: 1},
					{ID: 2, PacienteID: 1},
				}, nil)
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, err error) {
				assert.NoError(t, err)
				assert.Len(t, recetas, 2)
			},
		},
		{
			name:       "error - invalid paciente id",
			pacienteID: 0,
			behavior:   func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, recetas []*entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, recetas)
				assert.Contains(t, err.Error(), "invalid")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRecetaRepo)
			tt.behavior(mockRepo)

			service := NewRecetaService(mockRepo)
			recetas, err := service.GetHistorialByPacienteID(context.Background(), tt.pacienteID)

			tt.asserts(t, recetas, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRecetaService_UpdateReceta(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		req      *entities.UpdateRecetaRequest
		behavior func(m *MockRecetaRepo)
		asserts  func(t *testing.T, receta *entities.Receta, err error)
	}{
		{
			name: "success - updates receta",
			id:   1,
			req: &entities.UpdateRecetaRequest{
				ODEsfera: func() *float64 { v := -3.0; return &v }(),
			},
			behavior: func(m *MockRecetaRepo) {
				m.On("GetByID", mock.Anything, int64(1)).Return(&entities.Receta{
					ID:       1,
					ODEsfera: -2.5,
				}, nil)
				m.On("Update", mock.Anything, int64(1), mock.Anything).Return(&entities.Receta{
					ID:       1,
					ODEsfera: -3.0,
				}, nil)
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, receta)
			},
		},
		{
			name: "error - invalid id",
			id:   0,
			req:  &entities.UpdateRecetaRequest{},
			behavior: func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
				assert.Contains(t, err.Error(), "invalid")
			},
		},
		{
			name: "error - invalid od_esfera",
			id:   1,
			req: &entities.UpdateRecetaRequest{
				ODEsfera: func() *float64 { v := -25.0; return &v }(),
			},
			behavior: func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
				assert.Contains(t, err.Error(), "od_esfera")
			},
		},
		{
			name: "error - invalid tipo_lente",
			id:   1,
			req: &entities.UpdateRecetaRequest{
				TipoLente: func() *string { v := "invalid"; return &v }(),
			},
			behavior: func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
				assert.Contains(t, err.Error(), "tipo_lente")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRecetaRepo)
			tt.behavior(mockRepo)

			service := NewRecetaService(mockRepo)
			receta, err := service.UpdateReceta(context.Background(), tt.id, tt.req)

			tt.asserts(t, receta, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRecetaService_DeleteReceta(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(m *MockRecetaRepo)
		asserts  func(t *testing.T, err error)
	}{
		{
			name: "success - deletes receta",
			id:   1,
			behavior: func(m *MockRecetaRepo) {
				m.On("Delete", mock.Anything, int64(1)).Return(nil)
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:     "error - invalid id",
			id:       0,
			behavior: func(m *MockRecetaRepo) {},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid")
			},
		},
		{
			name: "error - repository fails",
			id:   999,
			behavior: func(m *MockRecetaRepo) {
				m.On("Delete", mock.Anything, int64(999)).Return(errors.New("not found"))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRecetaRepo)
			tt.behavior(mockRepo)

			service := NewRecetaService(mockRepo)
			err := service.DeleteReceta(context.Background(), tt.id)

			tt.asserts(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// MockRecetaRepo is a mock implementation of RecetaRepository
type MockRecetaRepo struct {
	mock.Mock
}

func (m *MockRecetaRepo) Create(ctx context.Context, receta *entities.Receta) (*entities.Receta, error) {
	args := m.Called(ctx, receta)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Receta), args.Error(1)
}

func (m *MockRecetaRepo) GetByID(ctx context.Context, id int64) (*entities.Receta, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Receta), args.Error(1)
}

func (m *MockRecetaRepo) GetAll(ctx context.Context, limit, offset int) ([]*entities.Receta, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Receta), args.Error(1)
}

func (m *MockRecetaRepo) GetByPacienteID(ctx context.Context, pacienteID int64, limit, offset int) ([]*entities.Receta, error) {
	args := m.Called(ctx, pacienteID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Receta), args.Error(1)
}

func (m *MockRecetaRepo) GetHistorialByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error) {
	args := m.Called(ctx, pacienteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Receta), args.Error(1)
}

func (m *MockRecetaRepo) GetLastTwoByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error) {
	args := m.Called(ctx, pacienteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Receta), args.Error(1)
}

func (m *MockRecetaRepo) Update(ctx context.Context, id int64, receta *entities.Receta) (*entities.Receta, error) {
	args := m.Called(ctx, id, receta)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Receta), args.Error(1)
}

func (m *MockRecetaRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRecetaRepo) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockRecetaRepo) CountByPacienteID(ctx context.Context, pacienteID int64) (int, error) {
	args := m.Called(ctx, pacienteID)
	return args.Int(0), args.Error(1)
}
