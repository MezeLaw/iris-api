package repositories

import (
	"context"

	"iris-api/internal/domain/entities"
)

type RecetaRepository interface {
	Create(ctx context.Context, receta *entities.Receta) (*entities.Receta, error)
	GetByID(ctx context.Context, id int64) (*entities.Receta, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Receta, error)
	GetByPacienteID(ctx context.Context, pacienteID int64, limit, offset int) ([]*entities.Receta, error)
	GetHistorialByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error)
	GetLastTwoByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error)
	Update(ctx context.Context, id int64, receta *entities.Receta) (*entities.Receta, error)
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
	CountByPacienteID(ctx context.Context, pacienteID int64) (int, error)
}
