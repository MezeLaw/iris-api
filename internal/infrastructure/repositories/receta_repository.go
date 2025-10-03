package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"
)

type recetaRepository struct {
	db *sqlx.DB
}

func NewRecetaRepository(db *sqlx.DB) repositories.RecetaRepository {
	return &recetaRepository{db: db}
}

func (r *recetaRepository) Create(ctx context.Context, receta *entities.Receta) (*entities.Receta, error) {
	query := `
		INSERT INTO recetas (paciente_id, fecha, od_esfera, od_cilindro, od_eje, oi_esfera, oi_cilindro, oi_eje, tipo_lente, observaciones, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, paciente_id, fecha, od_esfera, od_cilindro, od_eje, oi_esfera, oi_cilindro, oi_eje, tipo_lente, observaciones, created_at, updated_at
	`

	now := time.Now()
	receta.CreatedAt = now
	receta.UpdatedAt = now

	row := r.db.QueryRowContext(ctx, query,
		receta.PacienteID,
		receta.Fecha,
		receta.ODEsfera,
		receta.ODCilindro,
		receta.ODEje,
		receta.OIEsfera,
		receta.OICilindro,
		receta.OIEje,
		receta.TipoLente,
		receta.Observaciones,
		receta.CreatedAt,
		receta.UpdatedAt,
	)

	var createdReceta entities.Receta
	err := row.Scan(
		&createdReceta.ID,
		&createdReceta.PacienteID,
		&createdReceta.Fecha,
		&createdReceta.ODEsfera,
		&createdReceta.ODCilindro,
		&createdReceta.ODEje,
		&createdReceta.OIEsfera,
		&createdReceta.OICilindro,
		&createdReceta.OIEje,
		&createdReceta.TipoLente,
		&createdReceta.Observaciones,
		&createdReceta.CreatedAt,
		&createdReceta.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create receta: %w", err)
	}

	return &createdReceta, nil
}

func (r *recetaRepository) GetByID(ctx context.Context, id int64) (*entities.Receta, error) {
	query := `
		SELECT id, paciente_id, fecha, od_esfera, od_cilindro, od_eje, oi_esfera, oi_cilindro, oi_eje, tipo_lente, observaciones, created_at, updated_at
		FROM recetas
		WHERE id = $1
	`

	var receta entities.Receta
	err := r.db.GetContext(ctx, &receta, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("receta not found")
		}
		return nil, fmt.Errorf("failed to get receta: %w", err)
	}

	return &receta, nil
}

func (r *recetaRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Receta, error) {
	query := `
		SELECT id, paciente_id, fecha, od_esfera, od_cilindro, od_eje, oi_esfera, oi_cilindro, oi_eje, tipo_lente, observaciones, created_at, updated_at
		FROM recetas
		ORDER BY fecha DESC, created_at DESC
		LIMIT $1 OFFSET $2
	`

	var recetas []*entities.Receta
	err := r.db.SelectContext(ctx, &recetas, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get recetas: %w", err)
	}

	return recetas, nil
}

func (r *recetaRepository) GetByPacienteID(ctx context.Context, pacienteID int64, limit, offset int) ([]*entities.Receta, error) {
	query := `
		SELECT id, paciente_id, fecha, od_esfera, od_cilindro, od_eje, oi_esfera, oi_cilindro, oi_eje, tipo_lente, observaciones, created_at, updated_at
		FROM recetas
		WHERE paciente_id = $1
		ORDER BY fecha DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	var recetas []*entities.Receta
	err := r.db.SelectContext(ctx, &recetas, query, pacienteID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get recetas for paciente: %w", err)
	}

	return recetas, nil
}

func (r *recetaRepository) GetHistorialByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error) {
	query := `
		SELECT id, paciente_id, fecha, od_esfera, od_cilindro, od_eje, oi_esfera, oi_cilindro, oi_eje, tipo_lente, observaciones, created_at, updated_at
		FROM recetas
		WHERE paciente_id = $1
		ORDER BY fecha DESC, created_at DESC
	`

	var recetas []*entities.Receta
	err := r.db.SelectContext(ctx, &recetas, query, pacienteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get historial for paciente: %w", err)
	}

	return recetas, nil
}

func (r *recetaRepository) GetLastTwoByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error) {
	query := `
		SELECT id, paciente_id, fecha, od_esfera, od_cilindro, od_eje, oi_esfera, oi_cilindro, oi_eje, tipo_lente, observaciones, created_at, updated_at
		FROM recetas
		WHERE paciente_id = $1
		ORDER BY fecha DESC, created_at DESC
		LIMIT 2
	`

	var recetas []*entities.Receta
	err := r.db.SelectContext(ctx, &recetas, query, pacienteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last two recetas for paciente: %w", err)
	}

	return recetas, nil
}

func (r *recetaRepository) Update(ctx context.Context, id int64, receta *entities.Receta) (*entities.Receta, error) {
	query := `
		UPDATE recetas
		SET paciente_id = $1, fecha = $2, od_esfera = $3, od_cilindro = $4, od_eje = $5, oi_esfera = $6, oi_cilindro = $7, oi_eje = $8, tipo_lente = $9, observaciones = $10, updated_at = $11
		WHERE id = $12
		RETURNING id, paciente_id, fecha, od_esfera, od_cilindro, od_eje, oi_esfera, oi_cilindro, oi_eje, tipo_lente, observaciones, created_at, updated_at
	`

	receta.UpdatedAt = time.Now()

	row := r.db.QueryRowContext(ctx, query,
		receta.PacienteID,
		receta.Fecha,
		receta.ODEsfera,
		receta.ODCilindro,
		receta.ODEje,
		receta.OIEsfera,
		receta.OICilindro,
		receta.OIEje,
		receta.TipoLente,
		receta.Observaciones,
		receta.UpdatedAt,
		id,
	)

	var updatedReceta entities.Receta
	err := row.Scan(
		&updatedReceta.ID,
		&updatedReceta.PacienteID,
		&updatedReceta.Fecha,
		&updatedReceta.ODEsfera,
		&updatedReceta.ODCilindro,
		&updatedReceta.ODEje,
		&updatedReceta.OIEsfera,
		&updatedReceta.OICilindro,
		&updatedReceta.OIEje,
		&updatedReceta.TipoLente,
		&updatedReceta.Observaciones,
		&updatedReceta.CreatedAt,
		&updatedReceta.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("receta not found")
		}
		return nil, fmt.Errorf("failed to update receta: %w", err)
	}

	return &updatedReceta, nil
}

func (r *recetaRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM recetas WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete receta: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("receta not found")
	}

	return nil
}

func (r *recetaRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM recetas`

	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count recetas: %w", err)
	}

	return count, nil
}

func (r *recetaRepository) CountByPacienteID(ctx context.Context, pacienteID int64) (int, error) {
	query := `SELECT COUNT(*) FROM recetas WHERE paciente_id = $1`

	var count int
	err := r.db.GetContext(ctx, &count, query, pacienteID)
	if err != nil {
		return 0, fmt.Errorf("failed to count recetas for paciente: %w", err)
	}

	return count, nil
}
