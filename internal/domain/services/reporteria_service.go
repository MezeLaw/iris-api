package services

import (
	"context"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"
)

type ReporteriaService struct {
	reporteriaRepo repositories.ReporteriaRepository
}

func NewReporteriaService(reporteriaRepo repositories.ReporteriaRepository) *ReporteriaService {
	return &ReporteriaService{
		reporteriaRepo: reporteriaRepo,
	}
}

func (s *ReporteriaService) GetPacientesActivos(ctx context.Context, clientID int64) (*entities.ReportePacientesActivos, error) {
	pacientes, err := s.reporteriaRepo.GetPacientesActivos(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return &entities.ReportePacientesActivos{
		Pacientes: pacientes,
		Total:     len(pacientes),
	}, nil
}

func (s *ReporteriaService) GetPacientesInactivos(ctx context.Context, clientID int64) (*entities.ReportePacientesInactivos, error) {
	// 60 días = 2 meses aproximadamente
	diasInactividad := 60

	pacientes, err := s.reporteriaRepo.GetPacientesInactivos(ctx, clientID, diasInactividad)
	if err != nil {
		return nil, err
	}

	return &entities.ReportePacientesInactivos{
		Pacientes: pacientes,
		Total:     len(pacientes),
	}, nil
}
