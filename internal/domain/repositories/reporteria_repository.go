package repositories

import (
	"context"

	"iris-api/internal/domain/entities"
)

type ReporteriaRepository interface {
	GetPacientesActivos(ctx context.Context, clientID int64) ([]*entities.PacienteActivo, error)
	GetPacientesInactivos(ctx context.Context, clientID int64, diasInactividad int) ([]*entities.PacienteInactivo, error)
}
