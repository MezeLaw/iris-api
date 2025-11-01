package repositories

import (
	"context"
	"time"

	"iris-api/internal/domain/entities"
)

type TurnoRepository interface {
	Create(ctx context.Context, turno *entities.Turno) (*entities.Turno, error)
	GetByID(ctx context.Context, id int64, clientID int64) (*entities.TurnoConDetalles, error)
	GetAll(ctx context.Context, filters *entities.TurnoFilters, clientID int64) ([]*entities.TurnoConDetalles, int, error)
	Update(ctx context.Context, id int64, turno *entities.Turno, clientID int64) (*entities.Turno, error)
	Delete(ctx context.Context, id int64, clientID int64) error
	CambiarEstado(ctx context.Context, id int64, estado entities.EstadoTurno, clientID int64) error

	// Vistas específicas
	GetByDia(ctx context.Context, fecha time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error)
	GetBySemana(ctx context.Context, fechaInicio time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error)
	GetByProfesional(ctx context.Context, profesionalID int64, fechaDesde, fechaHasta *time.Time, clientID int64) ([]*entities.TurnoConDetalles, error)

	// Validaciones
	CheckDisponibilidad(ctx context.Context, req *entities.DisponibilidadRequest, clientID int64) (*entities.DisponibilidadResponse, error)
}
