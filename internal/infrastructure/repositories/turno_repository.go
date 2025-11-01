package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"iris-api/internal/domain/entities"
)

type turnoRepository struct {
	db *sqlx.DB
}

func NewTurnoRepository(db *sqlx.DB) *turnoRepository {
	return &turnoRepository{db: db}
}

func (r *turnoRepository) Create(ctx context.Context, turno *entities.Turno) (*entities.Turno, error) {
	query := `
		INSERT INTO turnos (
			client_id, paciente_id, profesional_user_id, tipo_servicio,
			fecha_hora, duracion_minutos, hora_fin, estado, observaciones,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		turno.ClientID,
		turno.PacienteID,
		turno.ProfesionalUserID,
		turno.TipoServicio,
		turno.FechaHora,
		turno.DuracionMinutos,
		turno.HoraFin,
		turno.Estado,
		turno.Observaciones,
		turno.CreatedAt,
		turno.UpdatedAt,
	).Scan(&turno.ID, &turno.CreatedAt, &turno.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return turno, nil
}

func (r *turnoRepository) GetByID(ctx context.Context, id int64, clientID int64) (*entities.TurnoConDetalles, error) {
	query := `
		SELECT
			t.id, t.client_id, t.paciente_id, t.profesional_user_id,
			t.tipo_servicio, t.fecha_hora, t.duracion_minutos, t.hora_fin,
			t.estado, t.observaciones, t.created_at, t.updated_at, t.deleted_at,
			p.nombre as paciente_nombre,
			p.apellido as paciente_apellido,
			p.email as paciente_email,
			u.name as profesional_nombre,
			u.lastname as profesional_apellido,
			u.email as profesional_email
		FROM turnos t
		INNER JOIN pacientes p ON t.paciente_id = p.id
		INNER JOIN users u ON t.profesional_user_id = u.id
		WHERE t.id = $1 AND t.client_id = $2 AND t.deleted_at IS NULL
	`

	var turno entities.TurnoConDetalles
	err := r.db.GetContext(ctx, &turno, query, id, clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("turno not found")
		}
		return nil, err
	}

	return &turno, nil
}

func (r *turnoRepository) GetAll(ctx context.Context, filters *entities.TurnoFilters, clientID int64) ([]*entities.TurnoConDetalles, int, error) {
	// Build WHERE clause
	baseQuery := `
		FROM turnos t
		INNER JOIN pacientes p ON t.paciente_id = p.id
		INNER JOIN users u ON t.profesional_user_id = u.id
		WHERE t.client_id = $1 AND t.deleted_at IS NULL
	`

	args := []interface{}{clientID}
	argPos := 2

	whereConditions := ""

	if filters.FechaDesde != nil {
		whereConditions += fmt.Sprintf(" AND t.fecha_hora >= $%d", argPos)
		args = append(args, *filters.FechaDesde)
		argPos++
	}

	if filters.FechaHasta != nil {
		whereConditions += fmt.Sprintf(" AND t.fecha_hora <= $%d", argPos)
		args = append(args, *filters.FechaHasta)
		argPos++
	}

	if filters.ProfesionalUserID != nil {
		whereConditions += fmt.Sprintf(" AND t.profesional_user_id = $%d", argPos)
		args = append(args, *filters.ProfesionalUserID)
		argPos++
	}

	if filters.PacienteID != nil {
		whereConditions += fmt.Sprintf(" AND t.paciente_id = $%d", argPos)
		args = append(args, *filters.PacienteID)
		argPos++
	}

	if filters.Estado != nil {
		whereConditions += fmt.Sprintf(" AND t.estado = $%d", argPos)
		args = append(args, *filters.Estado)
		argPos++
	}

	// Count query
	countQuery := "SELECT COUNT(*) " + baseQuery + whereConditions
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Data query with pagination
	offset := (filters.Page - 1) * filters.PageSize
	dataQuery := `
		SELECT
			t.id, t.client_id, t.paciente_id, t.profesional_user_id,
			t.tipo_servicio, t.fecha_hora, t.duracion_minutos, t.hora_fin,
			t.estado, t.observaciones, t.created_at, t.updated_at, t.deleted_at,
			p.nombre as paciente_nombre,
			p.apellido as paciente_apellido,
			p.email as paciente_email,
			u.name as profesional_nombre,
			u.lastname as profesional_apellido,
			u.email as profesional_email
	` + baseQuery + whereConditions + `
		ORDER BY t.fecha_hora ASC
		LIMIT $` + fmt.Sprintf("%d", argPos) + ` OFFSET $` + fmt.Sprintf("%d", argPos+1)

	args = append(args, filters.PageSize, offset)

	var turnos []*entities.TurnoConDetalles
	err = r.db.SelectContext(ctx, &turnos, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return turnos, total, nil
}

func (r *turnoRepository) Update(ctx context.Context, id int64, turno *entities.Turno, clientID int64) (*entities.Turno, error) {
	query := `
		UPDATE turnos
		SET profesional_user_id = $1,
			tipo_servicio = $2,
			fecha_hora = $3,
			duracion_minutos = $4,
			hora_fin = $5,
			observaciones = $6,
			updated_at = $7
		WHERE id = $8 AND client_id = $9 AND deleted_at IS NULL
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		turno.ProfesionalUserID,
		turno.TipoServicio,
		turno.FechaHora,
		turno.DuracionMinutos,
		turno.HoraFin,
		turno.Observaciones,
		turno.UpdatedAt,
		id,
		clientID,
	).Scan(&turno.ID, &turno.CreatedAt, &turno.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("turno not found")
		}
		return nil, err
	}

	// Fill in the rest of the turno struct
	turno.ClientID = clientID
	turno.Estado = turno.Estado // Preserve estado

	return turno, nil
}

func (r *turnoRepository) Delete(ctx context.Context, id int64, clientID int64) error {
	query := `
		UPDATE turnos
		SET deleted_at = $1, updated_at = $2
		WHERE id = $3 AND client_id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), time.Now(), id, clientID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("turno not found")
	}

	return nil
}

func (r *turnoRepository) CambiarEstado(ctx context.Context, id int64, estado entities.EstadoTurno, clientID int64) error {
	query := `
		UPDATE turnos
		SET estado = $1, updated_at = $2
		WHERE id = $3 AND client_id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, estado, time.Now(), id, clientID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("turno not found")
	}

	return nil
}

// Vistas específicas
func (r *turnoRepository) GetByDia(ctx context.Context, fecha time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error) {
	inicioDay := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, fecha.Location())
	finDay := inicioDay.Add(24 * time.Hour)

	query := `
		SELECT
			t.id, t.client_id, t.paciente_id, t.profesional_user_id,
			t.tipo_servicio, t.fecha_hora, t.duracion_minutos, t.hora_fin,
			t.estado, t.observaciones, t.created_at, t.updated_at, t.deleted_at,
			p.nombre as paciente_nombre,
			p.apellido as paciente_apellido,
			p.email as paciente_email,
			u.name as profesional_nombre,
			u.lastname as profesional_apellido,
			u.email as profesional_email
		FROM turnos t
		INNER JOIN pacientes p ON t.paciente_id = p.id
		INNER JOIN users u ON t.profesional_user_id = u.id
		WHERE t.client_id = $1
			AND t.fecha_hora >= $2
			AND t.fecha_hora < $3
			AND t.deleted_at IS NULL
	`

	args := []interface{}{clientID, inicioDay, finDay}

	if profesionalID != nil {
		query += ` AND t.profesional_user_id = $4`
		args = append(args, *profesionalID)
	}

	query += ` ORDER BY t.fecha_hora ASC`

	var turnos []*entities.TurnoConDetalles
	err := r.db.SelectContext(ctx, &turnos, query, args...)
	if err != nil {
		return nil, err
	}

	return turnos, nil
}

func (r *turnoRepository) GetBySemana(ctx context.Context, fechaInicio time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error) {
	// Ajustar al inicio de la semana (lunes)
	diasHastaLunes := int(fechaInicio.Weekday()-time.Monday) % 7
	if diasHastaLunes < 0 {
		diasHastaLunes += 7
	}
	inicioSemana := fechaInicio.AddDate(0, 0, -diasHastaLunes)
	inicioSemana = time.Date(inicioSemana.Year(), inicioSemana.Month(), inicioSemana.Day(), 0, 0, 0, 0, inicioSemana.Location())
	finSemana := inicioSemana.AddDate(0, 0, 7)

	query := `
		SELECT
			t.id, t.client_id, t.paciente_id, t.profesional_user_id,
			t.tipo_servicio, t.fecha_hora, t.duracion_minutos, t.hora_fin,
			t.estado, t.observaciones, t.created_at, t.updated_at, t.deleted_at,
			p.nombre as paciente_nombre,
			p.apellido as paciente_apellido,
			p.email as paciente_email,
			u.name as profesional_nombre,
			u.lastname as profesional_apellido,
			u.email as profesional_email
		FROM turnos t
		INNER JOIN pacientes p ON t.paciente_id = p.id
		INNER JOIN users u ON t.profesional_user_id = u.id
		WHERE t.client_id = $1
			AND t.fecha_hora >= $2
			AND t.fecha_hora < $3
			AND t.deleted_at IS NULL
	`

	args := []interface{}{clientID, inicioSemana, finSemana}

	if profesionalID != nil {
		query += ` AND t.profesional_user_id = $4`
		args = append(args, *profesionalID)
	}

	query += ` ORDER BY t.fecha_hora ASC`

	var turnos []*entities.TurnoConDetalles
	err := r.db.SelectContext(ctx, &turnos, query, args...)
	if err != nil {
		return nil, err
	}

	return turnos, nil
}

func (r *turnoRepository) GetByProfesional(ctx context.Context, profesionalID int64, fechaDesde, fechaHasta *time.Time, clientID int64) ([]*entities.TurnoConDetalles, error) {
	query := `
		SELECT
			t.id, t.client_id, t.paciente_id, t.profesional_user_id,
			t.tipo_servicio, t.fecha_hora, t.duracion_minutos, t.hora_fin,
			t.estado, t.observaciones, t.created_at, t.updated_at, t.deleted_at,
			p.nombre as paciente_nombre,
			p.apellido as paciente_apellido,
			p.email as paciente_email,
			u.name as profesional_nombre,
			u.lastname as profesional_apellido,
			u.email as profesional_email
		FROM turnos t
		INNER JOIN pacientes p ON t.paciente_id = p.id
		INNER JOIN users u ON t.profesional_user_id = u.id
		WHERE t.client_id = $1
			AND t.profesional_user_id = $2
			AND t.deleted_at IS NULL
	`

	args := []interface{}{clientID, profesionalID}
	argPos := 3

	if fechaDesde != nil {
		query += fmt.Sprintf(" AND t.fecha_hora >= $%d", argPos)
		args = append(args, *fechaDesde)
		argPos++
	}

	if fechaHasta != nil {
		query += fmt.Sprintf(" AND t.fecha_hora <= $%d", argPos)
		args = append(args, *fechaHasta)
	}

	query += ` ORDER BY t.fecha_hora ASC`

	var turnos []*entities.TurnoConDetalles
	err := r.db.SelectContext(ctx, &turnos, query, args...)
	if err != nil {
		return nil, err
	}

	return turnos, nil
}

// Validaciones
func (r *turnoRepository) CheckDisponibilidad(ctx context.Context, req *entities.DisponibilidadRequest, clientID int64) (*entities.DisponibilidadResponse, error) {
	horaFin := req.FechaHora.Add(time.Duration(req.DuracionMinutos) * time.Minute)

	query := `
		SELECT
			id, fecha_hora, hora_fin, tipo_servicio, estado
		FROM turnos
		WHERE client_id = $1
			AND profesional_user_id = $2
			AND deleted_at IS NULL
			AND estado NOT IN ($3, $4)
			AND (
				(fecha_hora < $5 AND hora_fin > $6)
				OR (fecha_hora >= $6 AND fecha_hora < $5)
			)
	`

	args := []interface{}{
		clientID,
		req.ProfesionalUserID,
		entities.EstadoCancelado,
		entities.EstadoNoAsistio,
		horaFin,
		req.FechaHora,
	}

	if req.TurnoID != nil {
		query += ` AND id != $7`
		args = append(args, *req.TurnoID)
	}

	var turnosConflicto []entities.Turno
	err := r.db.SelectContext(ctx, &turnosConflicto, query, args...)
	if err != nil {
		return nil, err
	}

	if len(turnosConflicto) > 0 {
		return &entities.DisponibilidadResponse{
			Disponible:      false,
			Mensaje:         "El profesional tiene turnos agendados en ese horario",
			TurnosConflicto: turnosConflicto,
		}, nil
	}

	return &entities.DisponibilidadResponse{
		Disponible: true,
		Mensaje:    "El profesional está disponible en ese horario",
	}, nil
}
