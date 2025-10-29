package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"

	"github.com/jmoiron/sqlx"
)

type antecedentesVisualesRepository struct {
	db *sqlx.DB
}

// NewAntecedentesVisualesRepository crea una nueva instancia del repositorio
func NewAntecedentesVisualesRepository(db *sqlx.DB) repositories.AntecedentesVisualesRepository {
	return &antecedentesVisualesRepository{db: db}
}

func (r *antecedentesVisualesRepository) Create(ctx context.Context, antecedentes *entities.AntecedentesVisuales) error {
	query := `
		INSERT INTO antecedentes_visuales (
			paciente_id, usa_lentes_contacto, tipo_lente_actual,
			tiempo_uso_diario, marca_modelo, fecha_ultima_adaptacion,
			molestias_complicaciones
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query,
		antecedentes.PacienteID,
		antecedentes.UsaLentesContacto,
		antecedentes.TipoLenteActual,
		antecedentes.TiempoUsoDiario,
		antecedentes.MarcaModelo,
		antecedentes.FechaUltimaAdaptacion,
		antecedentes.MolestiaComplicaciones,
	).Scan(&antecedentes.ID, &antecedentes.CreatedAt, &antecedentes.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating antecedentes visuales: %w", err)
	}

	return nil
}

func (r *antecedentesVisualesRepository) GetByPacienteID(ctx context.Context, pacienteID int64) (*entities.AntecedentesVisuales, error) {
	var antecedentes entities.AntecedentesVisuales
	query := `SELECT * FROM antecedentes_visuales WHERE paciente_id = $1`

	err := r.db.GetContext(ctx, &antecedentes, query, pacienteID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No hay antecedentes aún
		}
		return nil, fmt.Errorf("error getting antecedentes visuales: %w", err)
	}

	return &antecedentes, nil
}

func (r *antecedentesVisualesRepository) Update(ctx context.Context, antecedentes *entities.AntecedentesVisuales) error {
	query := `
		UPDATE antecedentes_visuales SET
			usa_lentes_contacto = $1,
			tipo_lente_actual = $2,
			tiempo_uso_diario = $3,
			marca_modelo = $4,
			fecha_ultima_adaptacion = $5,
			molestias_complicaciones = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE paciente_id = $7
		RETURNING updated_at
	`

	err := r.db.QueryRowxContext(ctx, query,
		antecedentes.UsaLentesContacto,
		antecedentes.TipoLenteActual,
		antecedentes.TiempoUsoDiario,
		antecedentes.MarcaModelo,
		antecedentes.FechaUltimaAdaptacion,
		antecedentes.MolestiaComplicaciones,
		antecedentes.PacienteID,
	).Scan(&antecedentes.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("antecedentes visuales not found")
		}
		return fmt.Errorf("error updating antecedentes visuales: %w", err)
	}

	return nil
}

func (r *antecedentesVisualesRepository) Delete(ctx context.Context, pacienteID int64) error {
	query := `DELETE FROM antecedentes_visuales WHERE paciente_id = $1`

	result, err := r.db.ExecContext(ctx, query, pacienteID)
	if err != nil {
		return fmt.Errorf("error deleting antecedentes visuales: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("antecedentes visuales not found")
	}

	return nil
}
