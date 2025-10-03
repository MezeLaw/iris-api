package usecases

import (
	"context"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
)

type ReporteriaUseCase struct {
	reporteriaService *services.ReporteriaService
}

func NewReporteriaUseCase(reporteriaService *services.ReporteriaService) *ReporteriaUseCase {
	return &ReporteriaUseCase{
		reporteriaService: reporteriaService,
	}
}

func (uc *ReporteriaUseCase) GetPacientesActivos(ctx context.Context, clientID int64) (*entities.ReportePacientesActivos, error) {
	return uc.reporteriaService.GetPacientesActivos(ctx, clientID)
}

func (uc *ReporteriaUseCase) GetPacientesInactivos(ctx context.Context, clientID int64) (*entities.ReportePacientesInactivos, error) {
	return uc.reporteriaService.GetPacientesInactivos(ctx, clientID)
}
