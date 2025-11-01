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

func TestReporteriaUseCase_GetPacientesActivos(t *testing.T) {
	tests := []struct {
		name     string
		clientID int64
		behavior func(m *MockReporteriaRepo)
		asserts  func(t *testing.T, reporte *entities.ReportePacientesActivos, err error)
	}{
		{
			name:     "success - returns active patients report",
			clientID: 1,
			behavior: func(m *MockReporteriaRepo) {
				m.On("GetPacientesActivos", mock.Anything, int64(1)).Return([]*entities.PacienteActivo{
					{
						ID:               1,
						Name:             "Patient 1",
						Email:            "patient1@example.com",
						UltimoTurno:      time.Now(),
						TotalTurnos:      5,
						TurnosPendientes: 2,
					},
				}, nil)
			},
			asserts: func(t *testing.T, reporte *entities.ReportePacientesActivos, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reporte)
				assert.Equal(t, 1, reporte.Total)
				assert.Len(t, reporte.Pacientes, 1)
			},
		},
		{
			name:     "error - service fails",
			clientID: 1,
			behavior: func(m *MockReporteriaRepo) {
				m.On("GetPacientesActivos", mock.Anything, int64(1)).Return(nil, errors.New("service error"))
			},
			asserts: func(t *testing.T, reporte *entities.ReportePacientesActivos, err error) {
				assert.Error(t, err)
				assert.Nil(t, reporte)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockReporteriaRepo)
			tt.behavior(mockRepo)

			service := services.NewReporteriaService(mockRepo)
			useCase := NewReporteriaUseCase(service)
			reporte, err := useCase.GetPacientesActivos(context.Background(), tt.clientID)

			tt.asserts(t, reporte, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestReporteriaUseCase_GetPacientesInactivos(t *testing.T) {
	tests := []struct {
		name     string
		clientID int64
		behavior func(m *MockReporteriaRepo)
		asserts  func(t *testing.T, reporte *entities.ReportePacientesInactivos, err error)
	}{
		{
			name:     "success - returns inactive patients report",
			clientID: 1,
			behavior: func(m *MockReporteriaRepo) {
				m.On("GetPacientesInactivos", mock.Anything, int64(1), 60).Return([]*entities.PacienteInactivo{
					{
						ID:           1,
						Name:         "Inactive Patient",
						Email:        "inactive@example.com",
						UltimoTurno:  time.Now().Add(-90 * 24 * time.Hour),
						DiasInactivo: 90,
					},
				}, nil)
			},
			asserts: func(t *testing.T, reporte *entities.ReportePacientesInactivos, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reporte)
				assert.Equal(t, 1, reporte.Total)
				assert.Len(t, reporte.Pacientes, 1)
			},
		},
		{
			name:     "error - service fails",
			clientID: 1,
			behavior: func(m *MockReporteriaRepo) {
				m.On("GetPacientesInactivos", mock.Anything, int64(1), 60).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, reporte *entities.ReportePacientesInactivos, err error) {
				assert.Error(t, err)
				assert.Nil(t, reporte)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockReporteriaRepo)
			tt.behavior(mockRepo)

			service := services.NewReporteriaService(mockRepo)
			useCase := NewReporteriaUseCase(service)
			reporte, err := useCase.GetPacientesInactivos(context.Background(), tt.clientID)

			tt.asserts(t, reporte, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// MockReporteriaRepo is a mock implementation of ReporteriaRepository
type MockReporteriaRepo struct {
	mock.Mock
}

func (m *MockReporteriaRepo) GetPacientesActivos(ctx context.Context, clientID int64) ([]*entities.PacienteActivo, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.PacienteActivo), args.Error(1)
}

func (m *MockReporteriaRepo) GetPacientesInactivos(ctx context.Context, clientID int64, diasInactividad int) ([]*entities.PacienteInactivo, error) {
	args := m.Called(ctx, clientID, diasInactividad)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.PacienteInactivo), args.Error(1)
}
