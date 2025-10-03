package repositories

import (
	"context"
	"time"

	"iris-api/internal/domain/entities"
)

type TurnoRepository interface {
	Create(ctx context.Context, turno *entities.Turno) (*entities.Turno, error)
	GetByID(ctx context.Context, id int64) (*entities.TurnoConDetalles, error)
	GetAll(ctx context.Context, filter *entities.TurnoFilter) ([]*entities.TurnoConDetalles, error)
	Update(ctx context.Context, id int64, turno *entities.Turno) (*entities.Turno, error)
	CancelTurno(ctx context.Context, id int64, motivo string) error
	Delete(ctx context.Context, id int64) error

	// Vistas específicas
	GetByDia(ctx context.Context, fecha time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error)
	GetBySemana(ctx context.Context, fechaInicio time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error)
	GetByProfesional(ctx context.Context, contactologoID int64, fechaDesde, fechaHasta time.Time) ([]*entities.TurnoConDetalles, error)

	// Alertas y notificaciones
	GetProximosTurnos(ctx context.Context, horasAnticipacion int) ([]*entities.TurnoConDetalles, error)
	GetTurnosSinConfirmar(ctx context.Context) ([]*entities.TurnoConDetalles, error)

	// Validaciones
	CheckDisponibilidad(ctx context.Context, contactologoID int64, fechaHora time.Time, duracionMinutos int, excludeTurnoID *int64) (bool, error)

	Count(ctx context.Context, filter *entities.TurnoFilter) (int, error)
}
