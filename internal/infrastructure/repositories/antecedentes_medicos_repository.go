package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"

	"github.com/jmoiron/sqlx"
)

type antecedentesMedicosRepository struct {
	db *sqlx.DB
}

// NewAntecedentesMedicosRepository crea una nueva instancia del repositorio
func NewAntecedentesMedicosRepository(db *sqlx.DB) repositories.AntecedentesMedicosRepository {
	return &antecedentesMedicosRepository{db: db}
}

func (r *antecedentesMedicosRepository) Create(ctx context.Context, antecedentes *entities.AntecedentesMedicos) error {
	query := `
		INSERT INTO antecedentes_medicos (
			paciente_id, tiene_diabetes, tiene_hipertension, tiene_alergias,
			detalle_alergias, otras_enfermedades, medicacion_habitual,
			cirugias_previas, cirugias_oculares
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query,
		antecedentes.PacienteID,
		antecedentes.TieneDiabetes,
		antecedentes.TieneHipertension,
		antecedentes.TieneAlergias,
		antecedentes.DetalleAlergias,
		antecedentes.OtrasEnfermedades,
		antecedentes.MedicacionHabitual,
		antecedentes.CirugiasPrevias,
		antecedentes.CirugiasOculares,
	).Scan(&antecedentes.ID, &antecedentes.CreatedAt, &antecedentes.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating antecedentes medicos: %w", err)
	}

	return nil
}

func (r *antecedentesMedicosRepository) GetByPacienteID(ctx context.Context, pacienteID int64) (*entities.AntecedentesMedicos, error) {
	var antecedentes entities.AntecedentesMedicos
	query := `SELECT * FROM antecedentes_medicos WHERE paciente_id = $1`

	err := r.db.GetContext(ctx, &antecedentes, query, pacienteID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No hay antecedentes aún
		}
		return nil, fmt.Errorf("error getting antecedentes medicos: %w", err)
	}

	return &antecedentes, nil
}

func (r *antecedentesMedicosRepository) Update(ctx context.Context, antecedentes *entities.AntecedentesMedicos) error {
	query := `
		UPDATE antecedentes_medicos SET
			tiene_diabetes = $1,
			tiene_hipertension = $2,
			tiene_alergias = $3,
			detalle_alergias = $4,
			otras_enfermedades = $5,
			medicacion_habitual = $6,
			cirugias_previas = $7,
			cirugias_oculares = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE paciente_id = $9
		RETURNING updated_at
	`

	err := r.db.QueryRowxContext(ctx, query,
		antecedentes.TieneDiabetes,
		antecedentes.TieneHipertension,
		antecedentes.TieneAlergias,
		antecedentes.DetalleAlergias,
		antecedentes.OtrasEnfermedades,
		antecedentes.MedicacionHabitual,
		antecedentes.CirugiasPrevias,
		antecedentes.CirugiasOculares,
		antecedentes.PacienteID,
	).Scan(&antecedentes.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("antecedentes medicos not found")
		}
		return fmt.Errorf("error updating antecedentes medicos: %w", err)
	}

	return nil
}

func (r *antecedentesMedicosRepository) Delete(ctx context.Context, pacienteID int64) error {
	query := `DELETE FROM antecedentes_medicos WHERE paciente_id = $1`

	result, err := r.db.ExecContext(ctx, query, pacienteID)
	if err != nil {
		return fmt.Errorf("error deleting antecedentes medicos: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("antecedentes medicos not found")
	}

	return nil
}
