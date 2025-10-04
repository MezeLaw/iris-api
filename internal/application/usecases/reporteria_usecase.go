package usecases

import (
	"context"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
)

type ReporteriaUseCase interface {
	GetPacientesActivos(ctx context.Context, clientID int64) (*entities.ReportePacientesActivos, error)
	GetPacientesInactivos(ctx context.Context, clientID int64) (*entities.ReportePacientesInactivos, error)
}

type reporteriaUseCase struct {
	reporteriaService services.ReporteriaService
}

func NewReporteriaUseCase(reporteriaService services.ReporteriaService) ReporteriaUseCase {
	return &reporteriaUseCase{
		reporteriaService: reporteriaService,
	}
}

func (uc *reporteriaUseCase) GetPacientesActivos(ctx context.Context, clientID int64) (*entities.ReportePacientesActivos, error) {
	return uc.reporteriaService.GetPacientesActivos(ctx, clientID)
}

func (uc *reporteriaUseCase) GetPacientesInactivos(ctx context.Context, clientID int64) (*entities.ReportePacientesInactivos, error) {
	return uc.reporteriaService.GetPacientesInactivos(ctx, clientID)
}
