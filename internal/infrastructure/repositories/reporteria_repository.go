package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"

	"iris-api/internal/domain/entities"
)

type reporteriaRepository struct {
	db *sqlx.DB
}

func NewReporteriaRepository(db *sqlx.DB) *reporteriaRepository {
	return &reporteriaRepository{db: db}
}

func (r *reporteriaRepository) GetPacientesActivos(ctx context.Context, clientID int64) ([]*entities.PacienteActivo, error) {
	query := `
		SELECT
			u.id,
			u.name,
			u.email,
			MAX(t.fecha_hora) as ultimo_turno,
			COUNT(t.id) as total_turnos,
			COUNT(CASE WHEN t.estado = 'pendiente' OR t.estado = 'confirmado' THEN 1 END) as turnos_pendientes
		FROM users u
		INNER JOIN turnos t ON u.id = t.paciente_id
		WHERE u.client_id = $1
		  AND t.client_id = $1
		  AND t.fecha_hora >= NOW() - INTERVAL '60 days'
		GROUP BY u.id, u.name, u.email
		ORDER BY ultimo_turno DESC
	`

	var pacientes []*entities.PacienteActivo
	err := r.db.SelectContext(ctx, &pacientes, query, clientID)
	if err != nil {
		return nil, err
	}

	return pacientes, nil
}

func (r *reporteriaRepository) GetPacientesInactivos(ctx context.Context, clientID int64, diasInactividad int) ([]*entities.PacienteInactivo, error) {
	query := `
		SELECT
			u.id,
			u.name,
			u.email,
			MAX(t.fecha_hora) as ultimo_turno,
			EXTRACT(DAY FROM (NOW() - MAX(t.fecha_hora)))::int as dias_inactivo
		FROM users u
		INNER JOIN turnos t ON u.id = t.paciente_id
		WHERE u.client_id = $1
		  AND t.client_id = $1
		GROUP BY u.id, u.name, u.email
		HAVING MAX(t.fecha_hora) < NOW() - INTERVAL '1 day' * $2
		ORDER BY dias_inactivo DESC
	`

	var pacientes []*entities.PacienteInactivo
	err := r.db.SelectContext(ctx, &pacientes, query, clientID, diasInactividad)
	if err != nil {
		return nil, err
	}

	return pacientes, nil
}
