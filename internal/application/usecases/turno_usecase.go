package usecases

import (
	"context"
	"time"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
)

type TurnoUseCase interface {
	CreateTurno(ctx context.Context, req *entities.CreateTurnoRequest, clientID int64) (*entities.Turno, error)
	GetTurnoByID(ctx context.Context, id int64, clientID int64) (*entities.TurnoConDetalles, error)
	GetTurnos(ctx context.Context, filters *entities.TurnoFilters, clientID int64) (map[string]interface{}, error)
	UpdateTurno(ctx context.Context, id int64, req *entities.UpdateTurnoRequest, clientID int64) (*entities.Turno, error)
	DeleteTurno(ctx context.Context, id int64, clientID int64) error
	CambiarEstado(ctx context.Context, id int64, req *entities.CambiarEstadoRequest, clientID int64) error
	GetTurnosByDia(ctx context.Context, fecha time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error)
	GetTurnosBySemana(ctx context.Context, fechaInicio time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error)
	GetTurnosByProfesional(ctx context.Context, profesionalID int64, fechaDesde, fechaHasta *time.Time, clientID int64) ([]*entities.TurnoConDetalles, error)
	CheckDisponibilidad(ctx context.Context, req *entities.DisponibilidadRequest, clientID int64) (*entities.DisponibilidadResponse, error)
}

type turnoUseCase struct {
	turnoService services.TurnoService
}

func NewTurnoUseCase(turnoService services.TurnoService) TurnoUseCase {
	return &turnoUseCase{
		turnoService: turnoService,
	}
}

func (uc *turnoUseCase) CreateTurno(ctx context.Context, req *entities.CreateTurnoRequest, clientID int64) (*entities.Turno, error) {
	return uc.turnoService.CreateTurno(ctx, req, clientID)
}

func (uc *turnoUseCase) GetTurnoByID(ctx context.Context, id int64, clientID int64) (*entities.TurnoConDetalles, error) {
	return uc.turnoService.GetTurnoByID(ctx, id, clientID)
}

func (uc *turnoUseCase) GetTurnos(ctx context.Context, filters *entities.TurnoFilters, clientID int64) (map[string]interface{}, error) {
	turnos, total, err := uc.turnoService.GetTurnos(ctx, filters, clientID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"turnos":    turnos,
		"total":     total,
		"page":      filters.Page,
		"page_size": filters.PageSize,
	}, nil
}

func (uc *turnoUseCase) UpdateTurno(ctx context.Context, id int64, req *entities.UpdateTurnoRequest, clientID int64) (*entities.Turno, error) {
	return uc.turnoService.UpdateTurno(ctx, id, req, clientID)
}

func (uc *turnoUseCase) CambiarEstado(ctx context.Context, id int64, req *entities.CambiarEstadoRequest, clientID int64) error {
	return uc.turnoService.CambiarEstado(ctx, id, req, clientID)
}

func (uc *turnoUseCase) DeleteTurno(ctx context.Context, id int64, clientID int64) error {
	return uc.turnoService.DeleteTurno(ctx, id, clientID)
}

func (uc *turnoUseCase) GetTurnosByDia(ctx context.Context, fecha time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error) {
	return uc.turnoService.GetTurnosByDia(ctx, fecha, profesionalID, clientID)
}

func (uc *turnoUseCase) GetTurnosBySemana(ctx context.Context, fechaInicio time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error) {
	return uc.turnoService.GetTurnosBySemana(ctx, fechaInicio, profesionalID, clientID)
}

func (uc *turnoUseCase) GetTurnosByProfesional(ctx context.Context, profesionalID int64, fechaDesde, fechaHasta *time.Time, clientID int64) ([]*entities.TurnoConDetalles, error) {
	return uc.turnoService.GetTurnosByProfesional(ctx, profesionalID, fechaDesde, fechaHasta, clientID)
}

func (uc *turnoUseCase) CheckDisponibilidad(ctx context.Context, req *entities.DisponibilidadRequest, clientID int64) (*entities.DisponibilidadResponse, error) {
	return uc.turnoService.CheckDisponibilidad(ctx, req, clientID)
}
