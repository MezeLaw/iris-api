package usecases

import (
	"context"
	"time"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
)

type TurnoUseCase interface {
	CreateTurno(ctx context.Context, req *entities.CreateTurnoRequest) (*entities.Turno, error)
	GetTurnoByID(ctx context.Context, id int64) (*entities.TurnoConDetalles, error)
	GetTurnos(ctx context.Context, filter *entities.TurnoFilter) (map[string]interface{}, error)
	UpdateTurno(ctx context.Context, id int64, req *entities.UpdateTurnoRequest) (*entities.Turno, error)
	CancelTurno(ctx context.Context, id int64, motivo string) error
	DeleteTurno(ctx context.Context, id int64) error
	GetTurnosByDia(ctx context.Context, fecha time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error)
	GetTurnosBySemana(ctx context.Context, fechaInicio time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error)
	GetTurnosByProfesional(ctx context.Context, contactologoID int64, fechaDesde, fechaHasta time.Time) ([]*entities.TurnoConDetalles, error)
	GetProximosTurnosAlert(ctx context.Context) (*entities.ProximosTurnosAlert, error)
}

type turnoUseCase struct {
	turnoService services.TurnoService
}

func NewTurnoUseCase(turnoService services.TurnoService) TurnoUseCase {
	return &turnoUseCase{
		turnoService: turnoService,
	}
}

func (uc *turnoUseCase) CreateTurno(ctx context.Context, req *entities.CreateTurnoRequest) (*entities.Turno, error) {
	return uc.turnoService.CreateTurno(ctx, req)
}

func (uc *turnoUseCase) GetTurnoByID(ctx context.Context, id int64) (*entities.TurnoConDetalles, error) {
	return uc.turnoService.GetTurnoByID(ctx, id)
}

func (uc *turnoUseCase) GetTurnos(ctx context.Context, filter *entities.TurnoFilter) (map[string]interface{}, error) {
	turnos, err := uc.turnoService.GetTurnos(ctx, filter)
	if err != nil {
		return nil, err
	}

	total, err := uc.turnoService.CountTurnos(ctx, filter)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"turnos": turnos,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}, nil
}

func (uc *turnoUseCase) UpdateTurno(ctx context.Context, id int64, req *entities.UpdateTurnoRequest) (*entities.Turno, error) {
	return uc.turnoService.UpdateTurno(ctx, id, req)
}

func (uc *turnoUseCase) CancelTurno(ctx context.Context, id int64, motivo string) error {
	return uc.turnoService.CancelTurno(ctx, id, motivo)
}

func (uc *turnoUseCase) DeleteTurno(ctx context.Context, id int64) error {
	return uc.turnoService.DeleteTurno(ctx, id)
}

func (uc *turnoUseCase) GetTurnosByDia(ctx context.Context, fecha time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error) {
	return uc.turnoService.GetTurnosByDia(ctx, fecha, contactologoID)
}

func (uc *turnoUseCase) GetTurnosBySemana(ctx context.Context, fechaInicio time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error) {
	return uc.turnoService.GetTurnosBySemana(ctx, fechaInicio, contactologoID)
}

func (uc *turnoUseCase) GetTurnosByProfesional(ctx context.Context, contactologoID int64, fechaDesde, fechaHasta time.Time) ([]*entities.TurnoConDetalles, error) {
	return uc.turnoService.GetTurnosByProfesional(ctx, contactologoID, fechaDesde, fechaHasta)
}

func (uc *turnoUseCase) GetProximosTurnosAlert(ctx context.Context) (*entities.ProximosTurnosAlert, error) {
	return uc.turnoService.GetProximosTurnosAlert(ctx)
}
