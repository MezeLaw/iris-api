package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
)

func TestTurnoUseCase_CreateTurno(t *testing.T) {
	tests := []struct {
		name     string
		req      *entities.CreateTurnoRequest
		behavior func(m *MockTurnoRepo)
		asserts  func(t *testing.T, turno *entities.Turno, err error)
	}{
		{
			name: "success - creates turno",
			req: &entities.CreateTurnoRequest{
				PacienteID:      1,
				ContactologoID:  2,
				FechaHora:       time.Now().Add(24 * time.Hour),
				DuracionMinutos: 30,
				TipoServicio:    "consulta",
			},
			behavior: func(m *MockTurnoRepo) {
				m.On("CheckDisponibilidad", mock.Anything, int64(2), mock.Anything, 30, (*int64)(nil)).Return(true, nil)
				m.On("Create", mock.Anything, mock.Anything).Return(&entities.Turno{
					ID:              1,
					PacienteID:      1,
					ContactologoID:  2,
					DuracionMinutos: 30,
				}, nil)
			},
			asserts: func(t *testing.T, turno *entities.Turno, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, turno)
				assert.Equal(t, int64(1), turno.ID)
			},
		},
		{
			name: "error - fecha en el pasado",
			req: &entities.CreateTurnoRequest{
				PacienteID:      1,
				ContactologoID:  2,
				FechaHora:       time.Now().Add(-24 * time.Hour),
				DuracionMinutos: 30,
			},
			behavior: func(m *MockTurnoRepo) {},
			asserts: func(t *testing.T, turno *entities.Turno, err error) {
				assert.Error(t, err)
				assert.Nil(t, turno)
				assert.Contains(t, err.Error(), "pasado")
			},
		},
		{
			name: "error - profesional no disponible",
			req: &entities.CreateTurnoRequest{
				PacienteID:      1,
				ContactologoID:  2,
				FechaHora:       time.Now().Add(24 * time.Hour),
				DuracionMinutos: 30,
			},
			behavior: func(m *MockTurnoRepo) {
				m.On("CheckDisponibilidad", mock.Anything, int64(2), mock.Anything, 30, (*int64)(nil)).Return(false, nil)
			},
			asserts: func(t *testing.T, turno *entities.Turno, err error) {
				assert.Error(t, err)
				assert.Nil(t, turno)
				assert.Contains(t, err.Error(), "disponible")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepo)
			tt.behavior(mockRepo)

			turnoService := services.NewTurnoService(mockRepo)
			useCase := NewTurnoUseCase(turnoService)
			turno, err := useCase.CreateTurno(context.Background(), tt.req)

			tt.asserts(t, turno, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoUseCase_GetTurnoByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(m *MockTurnoRepo)
		asserts  func(t *testing.T, turno *entities.TurnoConDetalles, err error)
	}{
		{
			name: "success - gets turno by id",
			id:   1,
			behavior: func(m *MockTurnoRepo) {
				m.On("GetByID", mock.Anything, int64(1)).Return(&entities.TurnoConDetalles{
					Turno: entities.Turno{
						ID:              1,
						DuracionMinutos: 30,
					},
					PacienteNombre: "Test Patient",
				}, nil)
			},
			asserts: func(t *testing.T, turno *entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, turno)
				assert.Equal(t, int64(1), turno.ID)
			},
		},
		{
			name: "error - not found",
			id:   999,
			behavior: func(m *MockTurnoRepo) {
				m.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, turno *entities.TurnoConDetalles, err error) {
				assert.Error(t, err)
				assert.Nil(t, turno)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepo)
			tt.behavior(mockRepo)

			turnoService := services.NewTurnoService(mockRepo)
			useCase := NewTurnoUseCase(turnoService)
			turno, err := useCase.GetTurnoByID(context.Background(), tt.id)

			tt.asserts(t, turno, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoUseCase_GetTurnos(t *testing.T) {
	tests := []struct {
		name     string
		filter   *entities.TurnoFilter
		behavior func(m *MockTurnoRepo)
		asserts  func(t *testing.T, resp map[string]interface{}, err error)
	}{
		{
			name: "success - returns turnos",
			filter: &entities.TurnoFilter{
				Limit:  10,
				Offset: 0,
			},
			behavior: func(m *MockTurnoRepo) {
				m.On("GetAll", mock.Anything, mock.Anything).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 1}, PacienteNombre: "Patient 1"},
					{Turno: entities.Turno{ID: 2}, PacienteNombre: "Patient 2"},
				}, nil)
				m.On("Count", mock.Anything, mock.Anything).Return(2, nil)
			},
			asserts: func(t *testing.T, resp map[string]interface{}, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, 2, resp["total"])
			},
		},
		{
			name: "error - service fails",
			filter: &entities.TurnoFilter{
				Limit:  10,
				Offset: 0,
			},
			behavior: func(m *MockTurnoRepo) {
				m.On("GetAll", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp map[string]interface{}, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepo)
			tt.behavior(mockRepo)

			turnoService := services.NewTurnoService(mockRepo)
			useCase := NewTurnoUseCase(turnoService)
			resp, err := useCase.GetTurnos(context.Background(), tt.filter)

			tt.asserts(t, resp, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoUseCase_UpdateTurno(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		req      *entities.UpdateTurnoRequest
		behavior func(m *MockTurnoRepo)
		asserts  func(t *testing.T, turno *entities.Turno, err error)
	}{
		{
			name: "success - updates turno",
			id:   1,
			req: &entities.UpdateTurnoRequest{
				Estado: func() *entities.EstadoTurno { e := entities.EstadoConfirmado; return &e }(),
			},
			behavior: func(m *MockTurnoRepo) {
				m.On("GetByID", mock.Anything, int64(1)).Return(&entities.TurnoConDetalles{
					Turno: entities.Turno{
						ID:     1,
						Estado: entities.EstadoPendiente,
						FechaHora: time.Now().Add(24 * time.Hour),
					},
				}, nil)
				m.On("Update", mock.Anything, int64(1), mock.Anything).Return(&entities.Turno{
					ID:     1,
					Estado: entities.EstadoConfirmado,
				}, nil)
			},
			asserts: func(t *testing.T, turno *entities.Turno, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, turno)
			},
		},
		{
			name: "error - turno not found",
			id:   999,
			req:  &entities.UpdateTurnoRequest{},
			behavior: func(m *MockTurnoRepo) {
				m.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, turno *entities.Turno, err error) {
				assert.Error(t, err)
				assert.Nil(t, turno)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepo)
			tt.behavior(mockRepo)

			turnoService := services.NewTurnoService(mockRepo)
			useCase := NewTurnoUseCase(turnoService)
			turno, err := useCase.UpdateTurno(context.Background(), tt.id, tt.req)

			tt.asserts(t, turno, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoUseCase_DeleteTurno(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(m *MockTurnoRepo)
		asserts  func(t *testing.T, err error)
	}{
		{
			name: "success - deletes turno",
			id:   1,
			behavior: func(m *MockTurnoRepo) {
				m.On("Delete", mock.Anything, int64(1)).Return(nil)
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error - repository fails",
			id:   999,
			behavior: func(m *MockTurnoRepo) {
				m.On("Delete", mock.Anything, int64(999)).Return(errors.New("not found"))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepo)
			tt.behavior(mockRepo)

			turnoService := services.NewTurnoService(mockRepo)
			useCase := NewTurnoUseCase(turnoService)
			err := useCase.DeleteTurno(context.Background(), tt.id)

			tt.asserts(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoUseCase_GetTurnosByDia(t *testing.T) {
	tests := []struct {
		name           string
		fecha          time.Time
		contactologoID *int64
		behavior       func(m *MockTurnoRepo)
		asserts        func(t *testing.T, turnos []*entities.TurnoConDetalles, err error)
	}{
		{
			name:  "success - returns turnos by dia",
			fecha: time.Now(),
			behavior: func(m *MockTurnoRepo) {
				m.On("GetByDia", mock.Anything, mock.Anything, (*int64)(nil)).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 1}},
				}, nil)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepo)
			tt.behavior(mockRepo)

			turnoService := services.NewTurnoService(mockRepo)
			useCase := NewTurnoUseCase(turnoService)
			turnos, err := useCase.GetTurnosByDia(context.Background(), tt.fecha, tt.contactologoID)

			tt.asserts(t, turnos, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoUseCase_GetTurnosBySemana(t *testing.T) {
	tests := []struct {
		name           string
		fechaInicio    time.Time
		contactologoID *int64
		behavior       func(m *MockTurnoRepo)
		asserts        func(t *testing.T, turnos []*entities.TurnoConDetalles, err error)
	}{
		{
			name:        "success - returns turnos by semana",
			fechaInicio: time.Now(),
			behavior: func(m *MockTurnoRepo) {
				m.On("GetBySemana", mock.Anything, mock.Anything, (*int64)(nil)).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 1}},
				}, nil)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepo)
			tt.behavior(mockRepo)

			turnoService := services.NewTurnoService(mockRepo)
			useCase := NewTurnoUseCase(turnoService)
			turnos, err := useCase.GetTurnosBySemana(context.Background(), tt.fechaInicio, tt.contactologoID)

			tt.asserts(t, turnos, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoUseCase_GetTurnosByProfesional(t *testing.T) {
	tests := []struct {
		name           string
		contactologoID int64
		fechaDesde     time.Time
		fechaHasta     time.Time
		behavior       func(m *MockTurnoRepo)
		asserts        func(t *testing.T, turnos []*entities.TurnoConDetalles, err error)
	}{
		{
			name:           "success - returns turnos by profesional",
			contactologoID: 1,
			fechaDesde:     time.Now(),
			fechaHasta:     time.Now().Add(7 * 24 * time.Hour),
			behavior: func(m *MockTurnoRepo) {
				m.On("GetByProfesional", mock.Anything, int64(1), mock.Anything, mock.Anything).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 1}},
				}, nil)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepo)
			tt.behavior(mockRepo)

			turnoService := services.NewTurnoService(mockRepo)
			useCase := NewTurnoUseCase(turnoService)
			turnos, err := useCase.GetTurnosByProfesional(context.Background(), tt.contactologoID, tt.fechaDesde, tt.fechaHasta)

			tt.asserts(t, turnos, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoUseCase_GetProximosTurnosAlert(t *testing.T) {
	tests := []struct {
		name     string
		behavior func(m *MockTurnoRepo)
		asserts  func(t *testing.T, alert *entities.ProximosTurnosAlert, err error)
	}{
		{
			name: "success - returns alert",
			behavior: func(m *MockTurnoRepo) {
				m.On("GetProximosTurnos", mock.Anything, 24).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 1}},
				}, nil)
				m.On("GetTurnosSinConfirmar", mock.Anything).Return([]*entities.TurnoConDetalles{
					{Turno: entities.Turno{ID: 2}},
				}, nil)
			},
			asserts: func(t *testing.T, alert *entities.ProximosTurnosAlert, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, alert)
				assert.Equal(t, 1, alert.Count24h)
				assert.Equal(t, 1, alert.CountSinConfirmar)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepo)
			tt.behavior(mockRepo)

			turnoService := services.NewTurnoService(mockRepo)
			useCase := NewTurnoUseCase(turnoService)
			alert, err := useCase.GetProximosTurnosAlert(context.Background())

			tt.asserts(t, alert, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoUseCase_CancelTurno(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		motivo   string
		behavior func(m *MockTurnoRepo)
		asserts  func(t *testing.T, err error)
	}{
		{
			name:   "success - cancels turno",
			id:     1,
			motivo: "Paciente no puede asistir",
			behavior: func(m *MockTurnoRepo) {
				m.On("GetByID", mock.Anything, int64(1)).Return(&entities.TurnoConDetalles{
					Turno: entities.Turno{
						ID:     1,
						Estado: entities.EstadoPendiente,
					},
				}, nil)
				m.On("CancelTurno", mock.Anything, int64(1), "Paciente no puede asistir").Return(nil)
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "error - turno already canceled",
			id:     1,
			motivo: "Test",
			behavior: func(m *MockTurnoRepo) {
				m.On("GetByID", mock.Anything, int64(1)).Return(&entities.TurnoConDetalles{
					Turno: entities.Turno{
						ID:     1,
						Estado: entities.EstadoCancelado,
					},
				}, nil)
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cancelado")
			},
		},
		{
			name:   "error - turno not found",
			id:     999,
			motivo: "Test",
			behavior: func(m *MockTurnoRepo) {
				m.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepo)
			tt.behavior(mockRepo)

			turnoService := services.NewTurnoService(mockRepo)
			useCase := NewTurnoUseCase(turnoService)
			err := useCase.CancelTurno(context.Background(), tt.id, tt.motivo)

			tt.asserts(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// MockTurnoRepo is a mock implementation of TurnoRepository
type MockTurnoRepo struct {
	mock.Mock
}

func (m *MockTurnoRepo) Create(ctx context.Context, turno *entities.Turno) (*entities.Turno, error) {
	args := m.Called(ctx, turno)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Turno), args.Error(1)
}

func (m *MockTurnoRepo) GetByID(ctx context.Context, id int64) (*entities.TurnoConDetalles, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoRepo) GetAll(ctx context.Context, filter *entities.TurnoFilter) ([]*entities.TurnoConDetalles, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoRepo) Update(ctx context.Context, id int64, turno *entities.Turno) (*entities.Turno, error) {
	args := m.Called(ctx, id, turno)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Turno), args.Error(1)
}

func (m *MockTurnoRepo) CancelTurno(ctx context.Context, id int64, motivo string) error {
	args := m.Called(ctx, id, motivo)
	return args.Error(0)
}

func (m *MockTurnoRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTurnoRepo) GetByDia(ctx context.Context, fecha time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error) {
	args := m.Called(ctx, fecha, contactologoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoRepo) GetBySemana(ctx context.Context, fechaInicio time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error) {
	args := m.Called(ctx, fechaInicio, contactologoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoRepo) GetByProfesional(ctx context.Context, contactologoID int64, fechaDesde, fechaHasta time.Time) ([]*entities.TurnoConDetalles, error) {
	args := m.Called(ctx, contactologoID, fechaDesde, fechaHasta)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoRepo) GetProximosTurnos(ctx context.Context, horasAnticipacion int) ([]*entities.TurnoConDetalles, error) {
	args := m.Called(ctx, horasAnticipacion)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoRepo) GetTurnosSinConfirmar(ctx context.Context) ([]*entities.TurnoConDetalles, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.TurnoConDetalles), args.Error(1)
}

func (m *MockTurnoRepo) CheckDisponibilidad(ctx context.Context, contactologoID int64, fechaHora time.Time, duracionMinutos int, excludeTurnoID *int64) (bool, error) {
	args := m.Called(ctx, contactologoID, fechaHora, duracionMinutos, excludeTurnoID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTurnoRepo) Count(ctx context.Context, filter *entities.TurnoFilter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}
