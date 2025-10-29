package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"
	"math"

	"github.com/jmoiron/sqlx"
)

type examenVisualRepository struct {
	db *sqlx.DB
}

// NewExamenVisualRepository crea una nueva instancia del repositorio
func NewExamenVisualRepository(db *sqlx.DB) repositories.ExamenVisualRepository {
	return &examenVisualRepository{db: db}
}

func (r *examenVisualRepository) Create(ctx context.Context, examen *entities.ExamenVisual) error {
	query := `
		INSERT INTO examenes_visuales (
			paciente_id, fecha_examen, av_sc_od, av_sc_oi, av_cc_od, av_cc_oi,
			od_esfera, od_cilindro, od_eje, od_add,
			oi_esfera, oi_cilindro, oi_eje, oi_add,
			tipo_lente, tipo_lente_otro, observaciones, realizado_por_user_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query,
		examen.PacienteID,
		examen.FechaExamen,
		examen.AvScOD,
		examen.AvScOI,
		examen.AvCcOD,
		examen.AvCcOI,
		examen.ODEsfera,
		examen.ODCilindro,
		examen.ODEje,
		examen.ODAdd,
		examen.OIEsfera,
		examen.OICilindro,
		examen.OIEje,
		examen.OIAdd,
		examen.TipoLente,
		examen.TipoLenteOtro,
		examen.Observaciones,
		examen.RealizadoPorUserID,
	).Scan(&examen.ID, &examen.CreatedAt, &examen.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating examen visual: %w", err)
	}

	return nil
}

func (r *examenVisualRepository) GetByID(ctx context.Context, id int64) (*entities.ExamenVisual, error) {
	var examen entities.ExamenVisual
	query := `SELECT * FROM examenes_visuales WHERE id = $1`

	err := r.db.GetContext(ctx, &examen, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("examen visual not found")
		}
		return nil, fmt.Errorf("error getting examen visual: %w", err)
	}

	return &examen, nil
}

func (r *examenVisualRepository) Update(ctx context.Context, examen *entities.ExamenVisual) error {
	query := `
		UPDATE examenes_visuales SET
			fecha_examen = $1,
			av_sc_od = $2,
			av_sc_oi = $3,
			av_cc_od = $4,
			av_cc_oi = $5,
			od_esfera = $6,
			od_cilindro = $7,
			od_eje = $8,
			od_add = $9,
			oi_esfera = $10,
			oi_cilindro = $11,
			oi_eje = $12,
			oi_add = $13,
			tipo_lente = $14,
			tipo_lente_otro = $15,
			observaciones = $16,
			realizado_por_user_id = $17,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $18
		RETURNING updated_at
	`

	err := r.db.QueryRowxContext(ctx, query,
		examen.FechaExamen,
		examen.AvScOD,
		examen.AvScOI,
		examen.AvCcOD,
		examen.AvCcOI,
		examen.ODEsfera,
		examen.ODCilindro,
		examen.ODEje,
		examen.ODAdd,
		examen.OIEsfera,
		examen.OICilindro,
		examen.OIEje,
		examen.OIAdd,
		examen.TipoLente,
		examen.TipoLenteOtro,
		examen.Observaciones,
		examen.RealizadoPorUserID,
		examen.ID,
	).Scan(&examen.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("examen visual not found")
		}
		return fmt.Errorf("error updating examen visual: %w", err)
	}

	return nil
}

func (r *examenVisualRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM examenes_visuales WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting examen visual: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("examen visual not found")
	}

	return nil
}

func (r *examenVisualRepository) ListByPacienteID(ctx context.Context, pacienteID int64) ([]entities.ExamenVisual, error) {
	var examenes []entities.ExamenVisual
	query := `
		SELECT * FROM examenes_visuales
		WHERE paciente_id = $1
		ORDER BY fecha_examen DESC
	`

	err := r.db.SelectContext(ctx, &examenes, query, pacienteID)
	if err != nil {
		return nil, fmt.Errorf("error listing examenes visuales: %w", err)
	}

	return examenes, nil
}

func (r *examenVisualRepository) GetLatestByPacienteID(ctx context.Context, pacienteID int64) (*entities.ExamenVisual, error) {
	var examen entities.ExamenVisual
	query := `
		SELECT * FROM examenes_visuales
		WHERE paciente_id = $1
		ORDER BY fecha_examen DESC
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &examen, query, pacienteID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No hay exámenes aún
		}
		return nil, fmt.Errorf("error getting latest examen visual: %w", err)
	}

	return &examen, nil
}

func (r *examenVisualRepository) GetLastNExamenes(ctx context.Context, pacienteID int64, n int) ([]entities.ExamenVisual, error) {
	var examenes []entities.ExamenVisual
	query := `
		SELECT * FROM examenes_visuales
		WHERE paciente_id = $1
		ORDER BY fecha_examen DESC
		LIMIT $2
	`

	err := r.db.SelectContext(ctx, &examenes, query, pacienteID, n)
	if err != nil {
		return nil, fmt.Errorf("error getting last N examenes: %w", err)
	}

	return examenes, nil
}

func (r *examenVisualRepository) GetDiferencia(ctx context.Context, examenAnteriorID, examenActualID int64) (*entities.DiferenciaRefraccion, error) {
	// Obtener ambos exámenes
	examenAnterior, err := r.GetByID(ctx, examenAnteriorID)
	if err != nil {
		return nil, fmt.Errorf("error getting examen anterior: %w", err)
	}

	examenActual, err := r.GetByID(ctx, examenActualID)
	if err != nil {
		return nil, fmt.Errorf("error getting examen actual: %w", err)
	}

	// Verificar que sean del mismo paciente
	if examenAnterior.PacienteID != examenActual.PacienteID {
		return nil, fmt.Errorf("examenes are not from the same patient")
	}

	diferencia := &entities.DiferenciaRefraccion{
		PacienteID:     examenActual.PacienteID,
		ExamenAnterior: examenAnteriorID,
		ExamenActual:   examenActualID,
		FechaAnterior:  examenAnterior.FechaExamen,
		FechaActual:    examenActual.FechaExamen,
	}

	// Calcular diferencias OD
	if examenAnterior.ODEsfera.Valid && examenActual.ODEsfera.Valid {
		diff := examenActual.ODEsfera.Float64 - examenAnterior.ODEsfera.Float64
		diferencia.DiferenciaODEsfera = &diff
	}

	if examenAnterior.ODCilindro.Valid && examenActual.ODCilindro.Valid {
		diff := examenActual.ODCilindro.Float64 - examenAnterior.ODCilindro.Float64
		diferencia.DiferenciaODCilindro = &diff
	}

	// Calcular diferencias OI
	if examenAnterior.OIEsfera.Valid && examenActual.OIEsfera.Valid {
		diff := examenActual.OIEsfera.Float64 - examenAnterior.OIEsfera.Float64
		diferencia.DiferenciaOIEsfera = &diff
	}

	if examenAnterior.OICilindro.Valid && examenActual.OICilindro.Valid {
		diff := examenActual.OICilindro.Float64 - examenAnterior.OICilindro.Float64
		diferencia.DiferenciaOICilindro = &diff
	}

	// Detectar cambios significativos (> 0.25D)
	const umbralSignificativo = 0.25
	cambioSignificativo := false
	var mensajes []string

	if diferencia.DiferenciaODEsfera != nil && math.Abs(*diferencia.DiferenciaODEsfera) > umbralSignificativo {
		cambioSignificativo = true
		mensajes = append(mensajes, fmt.Sprintf("Cambio significativo en esfera OD: %.2fD", *diferencia.DiferenciaODEsfera))
	}

	if diferencia.DiferenciaODCilindro != nil && math.Abs(*diferencia.DiferenciaODCilindro) > umbralSignificativo {
		cambioSignificativo = true
		mensajes = append(mensajes, fmt.Sprintf("Cambio significativo en cilindro OD: %.2fD", *diferencia.DiferenciaODCilindro))
	}

	if diferencia.DiferenciaOIEsfera != nil && math.Abs(*diferencia.DiferenciaOIEsfera) > umbralSignificativo {
		cambioSignificativo = true
		mensajes = append(mensajes, fmt.Sprintf("Cambio significativo en esfera OI: %.2fD", *diferencia.DiferenciaOIEsfera))
	}

	if diferencia.DiferenciaOICilindro != nil && math.Abs(*diferencia.DiferenciaOICilindro) > umbralSignificativo {
		cambioSignificativo = true
		mensajes = append(mensajes, fmt.Sprintf("Cambio significativo en cilindro OI: %.2fD", *diferencia.DiferenciaOICilindro))
	}

	diferencia.AlertaCambioSignificativo = cambioSignificativo

	if cambioSignificativo {
		diferencia.Mensaje = "ALERTA: Se detectaron los siguientes cambios significativos (>0.25D): "
		for i, msg := range mensajes {
			if i > 0 {
				diferencia.Mensaje += ", "
			}
			diferencia.Mensaje += msg
		}
	} else {
		diferencia.Mensaje = "No se detectaron cambios significativos en la refracción"
	}

	return diferencia, nil
}
