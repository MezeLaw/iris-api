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

func TestReporteriaService_GetPacientesActivos(t *testing.T) {
	tests := []struct {
		name     string
		clientID int64
		behavior func(m *MockReporteriaRepository)
		asserts  func(t *testing.T, reporte *entities.ReportePacientesActivos, err error)
	}{
		{
			name:     "success - returns active patients report",
			clientID: 1,
			behavior: func(m *MockReporteriaRepository) {
				m.On("GetPacientesActivos", mock.Anything, int64(1)).Return([]*entities.PacienteActivo{
					{
						ID:               1,
						Name:             "Patient 1",
						Email:            "patient1@example.com",
						UltimoTurno:      time.Now().Add(-10 * 24 * time.Hour),
						TotalTurnos:      5,
						TurnosPendientes: 2,
					},
					{
						ID:               2,
						Name:             "Patient 2",
						Email:            "patient2@example.com",
						UltimoTurno:      time.Now().Add(-30 * 24 * time.Hour),
						TotalTurnos:      3,
						TurnosPendientes: 1,
					},
				}, nil)
			},
			asserts: func(t *testing.T, reporte *entities.ReportePacientesActivos, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reporte)
				assert.Equal(t, 2, reporte.Total)
				assert.Len(t, reporte.Pacientes, 2)
				assert.Equal(t, "Patient 1", reporte.Pacientes[0].Name)
				assert.Equal(t, 5, reporte.Pacientes[0].TotalTurnos)
			},
		},
		{
			name:     "success - returns empty report",
			clientID: 1,
			behavior: func(m *MockReporteriaRepository) {
				m.On("GetPacientesActivos", mock.Anything, int64(1)).Return([]*entities.PacienteActivo{}, nil)
			},
			asserts: func(t *testing.T, reporte *entities.ReportePacientesActivos, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reporte)
				assert.Equal(t, 0, reporte.Total)
				assert.Empty(t, reporte.Pacientes)
			},
		},
		{
			name:     "error - repository fails",
			clientID: 1,
			behavior: func(m *MockReporteriaRepository) {
				m.On("GetPacientesActivos", mock.Anything, int64(1)).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, reporte *entities.ReportePacientesActivos, err error) {
				assert.Error(t, err)
				assert.Nil(t, reporte)
				assert.Contains(t, err.Error(), "database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockReporteriaRepository)
			tt.behavior(mockRepo)

			service := NewReporteriaService(mockRepo)
			reporte, err := service.GetPacientesActivos(context.Background(), tt.clientID)

			tt.asserts(t, reporte, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestReporteriaService_GetPacientesInactivos(t *testing.T) {
	tests := []struct {
		name     string
		clientID int64
		behavior func(m *MockReporteriaRepository)
		asserts  func(t *testing.T, reporte *entities.ReportePacientesInactivos, err error)
	}{
		{
			name:     "success - returns inactive patients report",
			clientID: 1,
			behavior: func(m *MockReporteriaRepository) {
				m.On("GetPacientesInactivos", mock.Anything, int64(1), 60).Return([]*entities.PacienteInactivo{
					{
						ID:           1,
						Name:         "Inactive Patient 1",
						Email:        "inactive1@example.com",
						UltimoTurno:  time.Now().Add(-90 * 24 * time.Hour),
						DiasInactivo: 90,
					},
					{
						ID:           2,
						Name:         "Inactive Patient 2",
						Email:        "inactive2@example.com",
						UltimoTurno:  time.Now().Add(-120 * 24 * time.Hour),
						DiasInactivo: 120,
					},
				}, nil)
			},
			asserts: func(t *testing.T, reporte *entities.ReportePacientesInactivos, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reporte)
				assert.Equal(t, 2, reporte.Total)
				assert.Len(t, reporte.Pacientes, 2)
				assert.Equal(t, "Inactive Patient 1", reporte.Pacientes[0].Name)
				assert.Equal(t, 90, reporte.Pacientes[0].DiasInactivo)
			},
		},
		{
			name:     "success - returns empty report when no inactive patients",
			clientID: 1,
			behavior: func(m *MockReporteriaRepository) {
				m.On("GetPacientesInactivos", mock.Anything, int64(1), 60).Return([]*entities.PacienteInactivo{}, nil)
			},
			asserts: func(t *testing.T, reporte *entities.ReportePacientesInactivos, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reporte)
				assert.Equal(t, 0, reporte.Total)
				assert.Empty(t, reporte.Pacientes)
			},
		},
		{
			name:     "error - repository fails",
			clientID: 1,
			behavior: func(m *MockReporteriaRepository) {
				m.On("GetPacientesInactivos", mock.Anything, int64(1), 60).Return(nil, errors.New("query error"))
			},
			asserts: func(t *testing.T, reporte *entities.ReportePacientesInactivos, err error) {
				assert.Error(t, err)
				assert.Nil(t, reporte)
				assert.Contains(t, err.Error(), "query error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockReporteriaRepository)
			tt.behavior(mockRepo)

			service := NewReporteriaService(mockRepo)
			reporte, err := service.GetPacientesInactivos(context.Background(), tt.clientID)

			tt.asserts(t, reporte, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// MockReporteriaRepository is a mock implementation of ReporteriaRepository
type MockReporteriaRepository struct {
	mock.Mock
}

func (m *MockReporteriaRepository) GetPacientesActivos(ctx context.Context, clientID int64) ([]*entities.PacienteActivo, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.PacienteActivo), args.Error(1)
}

func (m *MockReporteriaRepository) GetPacientesInactivos(ctx context.Context, clientID int64, diasInactividad int) ([]*entities.PacienteInactivo, error) {
	args := m.Called(ctx, clientID, diasInactividad)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.PacienteInactivo), args.Error(1)
}
