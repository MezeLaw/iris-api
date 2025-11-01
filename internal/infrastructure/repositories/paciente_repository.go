package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"
	"math"
	"strings"

	"github.com/jmoiron/sqlx"
)

type pacienteRepository struct {
	db *sqlx.DB
}

// NewPacienteRepository crea una nueva instancia del repositorio de pacientes
func NewPacienteRepository(db *sqlx.DB) repositories.PacienteRepository {
	return &pacienteRepository{db: db}
}

func (r *pacienteRepository) Create(ctx context.Context, paciente *entities.Paciente) error {
	query := `
		INSERT INTO pacientes (
			client_id, nombre_completo, dni, fecha_nacimiento, genero,
			telefono, email, direccion, ocupacion, motivo_consulta,
			fecha_primera_visita, observaciones
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) RETURNING id, edad, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query,
		paciente.ClientID,
		paciente.NombreCompleto,
		paciente.DNI,
		paciente.FechaNacimiento,
		paciente.Genero,
		paciente.Telefono,
		paciente.Email,
		paciente.Direccion,
		paciente.Ocupacion,
		paciente.MotivoConsulta,
		paciente.FechaPrimeraVisita,
		paciente.Observaciones,
	).Scan(&paciente.ID, &paciente.Edad, &paciente.CreatedAt, &paciente.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating paciente: %w", err)
	}

	return nil
}

func (r *pacienteRepository) GetByID(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error) {
	var paciente entities.Paciente
	query := `
		SELECT * FROM pacientes
		WHERE id = $1 AND client_id = $2 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &paciente, query, id, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("paciente not found")
		}
		return nil, fmt.Errorf("error getting paciente: %w", err)
	}

	return &paciente, nil
}

func (r *pacienteRepository) Update(ctx context.Context, paciente *entities.Paciente) error {
	query := `
		UPDATE pacientes SET
			nombre_completo = $1,
			dni = $2,
			fecha_nacimiento = $3,
			genero = $4,
			telefono = $5,
			email = $6,
			direccion = $7,
			ocupacion = $8,
			motivo_consulta = $9,
			fecha_primera_visita = $10,
			observaciones = $11,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $12 AND client_id = $13 AND deleted_at IS NULL
		RETURNING edad, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query,
		paciente.NombreCompleto,
		paciente.DNI,
		paciente.FechaNacimiento,
		paciente.Genero,
		paciente.Telefono,
		paciente.Email,
		paciente.Direccion,
		paciente.Ocupacion,
		paciente.MotivoConsulta,
		paciente.FechaPrimeraVisita,
		paciente.Observaciones,
		paciente.ID,
		paciente.ClientID,
	).Scan(&paciente.Edad, &paciente.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("paciente not found or already deleted")
		}
		return fmt.Errorf("error updating paciente: %w", err)
	}

	return nil
}

func (r *pacienteRepository) Delete(ctx context.Context, id int64, clientID int64) error {
	query := `
		UPDATE pacientes SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND client_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, clientID)
	if err != nil {
		return fmt.Errorf("error deleting paciente: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("paciente not found or already deleted")
	}

	return nil
}

func (r *pacienteRepository) List(ctx context.Context, clientID int64, page, pageSize int) ([]entities.Paciente, int64, error) {
	offset := (page - 1) * pageSize

	// Obtener total de registros
	var total int64
	countQuery := `SELECT COUNT(*) FROM pacientes WHERE client_id = $1 AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &total, countQuery, clientID)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting pacientes: %w", err)
	}

	// Obtener pacientes paginados
	var pacientes []entities.Paciente
	query := `
		SELECT * FROM pacientes
		WHERE client_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &pacientes, query, clientID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing pacientes: %w", err)
	}

	return pacientes, total, nil
}

func (r *pacienteRepository) Search(ctx context.Context, clientID int64, query string, page, pageSize int) ([]entities.Paciente, int64, error) {
	offset := (page - 1) * pageSize
	searchPattern := "%" + strings.ToLower(query) + "%"

	// Obtener total de registros
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM pacientes
		WHERE client_id = $1 AND deleted_at IS NULL
		AND (
			LOWER(nombre_completo) LIKE $2 OR
			LOWER(dni) LIKE $2 OR
			LOWER(email) LIKE $2 OR
			LOWER(telefono) LIKE $2
		)
	`
	err := r.db.GetContext(ctx, &total, countQuery, clientID, searchPattern)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting pacientes: %w", err)
	}

	// Obtener pacientes paginados
	var pacientes []entities.Paciente
	searchQuery := `
		SELECT * FROM pacientes
		WHERE client_id = $1 AND deleted_at IS NULL
		AND (
			LOWER(nombre_completo) LIKE $2 OR
			LOWER(dni) LIKE $2 OR
			LOWER(email) LIKE $2 OR
			LOWER(telefono) LIKE $2
		)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	err = r.db.SelectContext(ctx, &pacientes, searchQuery, clientID, searchPattern, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error searching pacientes: %w", err)
	}

	return pacientes, total, nil
}

func (r *pacienteRepository) GetByDNI(ctx context.Context, dni string, clientID int64) (*entities.Paciente, error) {
	var paciente entities.Paciente
	query := `
		SELECT * FROM pacientes
		WHERE dni = $1 AND client_id = $2 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &paciente, query, dni, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("paciente not found")
		}
		return nil, fmt.Errorf("error getting paciente by DNI: %w", err)
	}

	return &paciente, nil
}

func (r *pacienteRepository) GetWithAntecedentes(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error) {
	paciente, err := r.GetByID(ctx, id, clientID)
	if err != nil {
		return nil, err
	}

	// Cargar antecedentes médicos
	var antecedentesMedicos entities.AntecedentesMedicos
	queryMedicos := `SELECT * FROM antecedentes_medicos WHERE paciente_id = $1`
	err = r.db.GetContext(ctx, &antecedentesMedicos, queryMedicos, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting antecedentes medicos: %w", err)
	}
	if err == nil {
		paciente.AntecedentesMedicos = &antecedentesMedicos
	}

	// Cargar antecedentes visuales
	var antecedentesVisuales entities.AntecedentesVisuales
	queryVisuales := `SELECT * FROM antecedentes_visuales WHERE paciente_id = $1`
	err = r.db.GetContext(ctx, &antecedentesVisuales, queryVisuales, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting antecedentes visuales: %w", err)
	}
	if err == nil {
		paciente.AntecedentesVisuales = &antecedentesVisuales
	}

	return paciente, nil
}

func (r *pacienteRepository) GetWithExamenes(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error) {
	paciente, err := r.GetByID(ctx, id, clientID)
	if err != nil {
		return nil, err
	}

	// Cargar exámenes visuales
	var examenes []entities.ExamenVisual
	queryExamenes := `
		SELECT * FROM examenes_visuales
		WHERE paciente_id = $1
		ORDER BY fecha_examen DESC
	`
	err = r.db.SelectContext(ctx, &examenes, queryExamenes, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting examenes visuales: %w", err)
	}
	if err == nil {
		paciente.ExamenesVisuales = examenes
	}

	return paciente, nil
}

func (r *pacienteRepository) GetComplete(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error) {
	paciente, err := r.GetByID(ctx, id, clientID)
	if err != nil {
		return nil, err
	}

	// Cargar antecedentes médicos
	var antecedentesMedicos entities.AntecedentesMedicos
	queryMedicos := `SELECT * FROM antecedentes_medicos WHERE paciente_id = $1`
	err = r.db.GetContext(ctx, &antecedentesMedicos, queryMedicos, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting antecedentes medicos: %w", err)
	}
	if err == nil {
		paciente.AntecedentesMedicos = &antecedentesMedicos
	}

	// Cargar antecedentes visuales
	var antecedentesVisuales entities.AntecedentesVisuales
	queryVisuales := `SELECT * FROM antecedentes_visuales WHERE paciente_id = $1`
	err = r.db.GetContext(ctx, &antecedentesVisuales, queryVisuales, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting antecedentes visuales: %w", err)
	}
	if err == nil {
		paciente.AntecedentesVisuales = &antecedentesVisuales
	}

	// Cargar exámenes visuales
	var examenes []entities.ExamenVisual
	queryExamenes := `
		SELECT * FROM examenes_visuales
		WHERE paciente_id = $1
		ORDER BY fecha_examen DESC
	`
	err = r.db.SelectContext(ctx, &examenes, queryExamenes, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting examenes visuales: %w", err)
	}
	if err == nil {
		paciente.ExamenesVisuales = examenes
	}

	return paciente, nil
}

func (r *pacienteRepository) CountByClientID(ctx context.Context, clientID int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM pacientes WHERE client_id = $1 AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &count, query, clientID)
	if err != nil {
		return 0, fmt.Errorf("error counting pacientes: %w", err)
	}

	return count, nil
}

// Helper function para calcular total de páginas
func calculateTotalPages(total int64, pageSize int) int {
	return int(math.Ceil(float64(total) / float64(pageSize)))
}
