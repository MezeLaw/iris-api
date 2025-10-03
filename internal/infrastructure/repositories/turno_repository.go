package repositories

import (
	"context"
	"database/sql"
	"errors"
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
		INSERT INTO turnos (paciente_id, contactologo_id, fecha_hora, duracion_minutos, tipo_servicio, estado, motivo, observaciones, recordatorio_enviado, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		turno.PacienteID,
		turno.ContactologoID,
		turno.FechaHora,
		turno.DuracionMinutos,
		turno.TipoServicio,
		turno.Estado,
		turno.Motivo,
		turno.Observaciones,
		turno.RecordatorioEnviado,
		turno.CreatedAt,
		turno.UpdatedAt,
	).Scan(&turno.ID, &turno.CreatedAt, &turno.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return turno, nil
}

func (r *turnoRepository) GetByID(ctx context.Context, id int64) (*entities.TurnoConDetalles, error) {
	query := `
		SELECT
			t.id, t.paciente_id, t.contactologo_id, t.fecha_hora, t.duracion_minutos,
			t.tipo_servicio, t.estado, t.motivo, t.observaciones, t.recordatorio_enviado,
			t.created_at, t.updated_at,
			p.name as paciente_nombre, p.email as paciente_email,
			c.name as contactologo_nombre, c.email as contactologo_email
		FROM turnos t
		INNER JOIN users p ON t.paciente_id = p.id
		INNER JOIN users c ON t.contactologo_id = c.id
		WHERE t.id = $1
	`

	var turno entities.TurnoConDetalles
	err := r.db.GetContext(ctx, &turno, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("turno not found")
		}
		return nil, err
	}

	return &turno, nil
}

func (r *turnoRepository) GetAll(ctx context.Context, filter *entities.TurnoFilter) ([]*entities.TurnoConDetalles, error) {
	query := `
		SELECT
			t.id, t.paciente_id, t.contactologo_id, t.fecha_hora, t.duracion_minutos,
			t.tipo_servicio, t.estado, t.motivo, t.observaciones, t.recordatorio_enviado,
			t.created_at, t.updated_at,
			p.name as paciente_nombre, p.email as paciente_email,
			c.name as contactologo_nombre, c.email as contactologo_email
		FROM turnos t
		INNER JOIN users p ON t.paciente_id = p.id
		INNER JOIN users c ON t.contactologo_id = c.id
		WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if filter.PacienteID != nil {
		query += ` AND t.paciente_id = $` + string(rune(argPos+'0'))
		args = append(args, *filter.PacienteID)
		argPos++
	}

	if filter.ContactologoID != nil {
		query += ` AND t.contactologo_id = $` + string(rune(argPos+'0'))
		args = append(args, *filter.ContactologoID)
		argPos++
	}

	if filter.TipoServicio != nil {
		query += ` AND t.tipo_servicio = $` + string(rune(argPos+'0'))
		args = append(args, *filter.TipoServicio)
		argPos++
	}

	if filter.Estado != nil {
		query += ` AND t.estado = $` + string(rune(argPos+'0'))
		args = append(args, *filter.Estado)
		argPos++
	}

	if filter.FechaDesde != nil {
		query += ` AND t.fecha_hora >= $` + string(rune(argPos+'0'))
		args = append(args, *filter.FechaDesde)
		argPos++
	}

	if filter.FechaHasta != nil {
		query += ` AND t.fecha_hora <= $` + string(rune(argPos+'0'))
		args = append(args, *filter.FechaHasta)
		argPos++
	}

	query += ` ORDER BY t.fecha_hora ASC LIMIT $` + string(rune(argPos+'0')) + ` OFFSET $` + string(rune(argPos+1+'0'))
	args = append(args, filter.Limit, filter.Offset)

	var turnos []*entities.TurnoConDetalles
	err := r.db.SelectContext(ctx, &turnos, query, args...)
	if err != nil {
		return nil, err
	}

	return turnos, nil
}

func (r *turnoRepository) Update(ctx context.Context, id int64, turno *entities.Turno) (*entities.Turno, error) {
	query := `
		UPDATE turnos
		SET paciente_id = $1, contactologo_id = $2, fecha_hora = $3, duracion_minutos = $4,
		    tipo_servicio = $5, estado = $6, motivo = $7, observaciones = $8,
		    recordatorio_enviado = $9, updated_at = $10
		WHERE id = $11
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		turno.PacienteID,
		turno.ContactologoID,
		turno.FechaHora,
		turno.DuracionMinutos,
		turno.TipoServicio,
		turno.Estado,
		turno.Motivo,
		turno.Observaciones,
		turno.RecordatorioEnviado,
		turno.UpdatedAt,
		id,
	).Scan(&turno.ID, &turno.CreatedAt, &turno.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("turno not found")
		}
		return nil, err
	}

	return turno, nil
}

func (r *turnoRepository) CancelTurno(ctx context.Context, id int64, motivo string) error {
	query := `
		UPDATE turnos
		SET estado = $1, observaciones = CONCAT(observaciones, ' | Cancelación: ', $2), updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, entities.EstadoCancelado, motivo, time.Now(), id)
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

func (r *turnoRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM turnos WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *turnoRepository) GetByDia(ctx context.Context, fecha time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error) {
	inicioDay := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, fecha.Location())
	finDay := inicioDay.Add(24 * time.Hour)

	query := `
		SELECT
			t.id, t.paciente_id, t.contactologo_id, t.fecha_hora, t.duracion_minutos,
			t.tipo_servicio, t.estado, t.motivo, t.observaciones, t.recordatorio_enviado,
			t.created_at, t.updated_at,
			p.name as paciente_nombre, p.email as paciente_email,
			c.name as contactologo_nombre, c.email as contactologo_email
		FROM turnos t
		INNER JOIN users p ON t.paciente_id = p.id
		INNER JOIN users c ON t.contactologo_id = c.id
		WHERE t.fecha_hora >= $1 AND t.fecha_hora < $2
	`

	args := []interface{}{inicioDay, finDay}

	if contactologoID != nil {
		query += ` AND t.contactologo_id = $3`
		args = append(args, *contactologoID)
	}

	query += ` ORDER BY t.fecha_hora ASC`

	var turnos []*entities.TurnoConDetalles
	err := r.db.SelectContext(ctx, &turnos, query, args...)
	if err != nil {
		return nil, err
	}

	return turnos, nil
}

func (r *turnoRepository) GetBySemana(ctx context.Context, fechaInicio time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error) {
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
			t.id, t.paciente_id, t.contactologo_id, t.fecha_hora, t.duracion_minutos,
			t.tipo_servicio, t.estado, t.motivo, t.observaciones, t.recordatorio_enviado,
			t.created_at, t.updated_at,
			p.name as paciente_nombre, p.email as paciente_email,
			c.name as contactologo_nombre, c.email as contactologo_email
		FROM turnos t
		INNER JOIN users p ON t.paciente_id = p.id
		INNER JOIN users c ON t.contactologo_id = c.id
		WHERE t.fecha_hora >= $1 AND t.fecha_hora < $2
	`

	args := []interface{}{inicioSemana, finSemana}

	if contactologoID != nil {
		query += ` AND t.contactologo_id = $3`
		args = append(args, *contactologoID)
	}

	query += ` ORDER BY t.fecha_hora ASC`

	var turnos []*entities.TurnoConDetalles
	err := r.db.SelectContext(ctx, &turnos, query, args...)
	if err != nil {
		return nil, err
	}

	return turnos, nil
}

func (r *turnoRepository) GetByProfesional(ctx context.Context, contactologoID int64, fechaDesde, fechaHasta time.Time) ([]*entities.TurnoConDetalles, error) {
	query := `
		SELECT
			t.id, t.paciente_id, t.contactologo_id, t.fecha_hora, t.duracion_minutos,
			t.tipo_servicio, t.estado, t.motivo, t.observaciones, t.recordatorio_enviado,
			t.created_at, t.updated_at,
			p.name as paciente_nombre, p.email as paciente_email,
			c.name as contactologo_nombre, c.email as contactologo_email
		FROM turnos t
		INNER JOIN users p ON t.paciente_id = p.id
		INNER JOIN users c ON t.contactologo_id = c.id
		WHERE t.contactologo_id = $1 AND t.fecha_hora >= $2 AND t.fecha_hora <= $3
		ORDER BY t.fecha_hora ASC
	`

	var turnos []*entities.TurnoConDetalles
	err := r.db.SelectContext(ctx, &turnos, query, contactologoID, fechaDesde, fechaHasta)
	if err != nil {
		return nil, err
	}

	return turnos, nil
}

func (r *turnoRepository) GetProximosTurnos(ctx context.Context, horasAnticipacion int) ([]*entities.TurnoConDetalles, error) {
	ahora := time.Now()
	limite := ahora.Add(time.Duration(horasAnticipacion) * time.Hour)

	query := `
		SELECT
			t.id, t.paciente_id, t.contactologo_id, t.fecha_hora, t.duracion_minutos,
			t.tipo_servicio, t.estado, t.motivo, t.observaciones, t.recordatorio_enviado,
			t.created_at, t.updated_at,
			p.name as paciente_nombre, p.email as paciente_email,
			c.name as contactologo_nombre, c.email as contactologo_email
		FROM turnos t
		INNER JOIN users p ON t.paciente_id = p.id
		INNER JOIN users c ON t.contactologo_id = c.id
		WHERE t.fecha_hora >= $1 AND t.fecha_hora <= $2 AND t.estado != $3
		ORDER BY t.fecha_hora ASC
	`

	var turnos []*entities.TurnoConDetalles
	err := r.db.SelectContext(ctx, &turnos, query, ahora, limite, entities.EstadoCancelado)
	if err != nil {
		return nil, err
	}

	return turnos, nil
}

func (r *turnoRepository) GetTurnosSinConfirmar(ctx context.Context) ([]*entities.TurnoConDetalles, error) {
	query := `
		SELECT
			t.id, t.paciente_id, t.contactologo_id, t.fecha_hora, t.duracion_minutos,
			t.tipo_servicio, t.estado, t.motivo, t.observaciones, t.recordatorio_enviado,
			t.created_at, t.updated_at,
			p.name as paciente_nombre, p.email as paciente_email,
			c.name as contactologo_nombre, c.email as contactologo_email
		FROM turnos t
		INNER JOIN users p ON t.paciente_id = p.id
		INNER JOIN users c ON t.contactologo_id = c.id
		WHERE t.estado = $1 AND t.fecha_hora > $2
		ORDER BY t.fecha_hora ASC
	`

	var turnos []*entities.TurnoConDetalles
	err := r.db.SelectContext(ctx, &turnos, query, entities.EstadoPendiente, time.Now())
	if err != nil {
		return nil, err
	}

	return turnos, nil
}

func (r *turnoRepository) CheckDisponibilidad(ctx context.Context, contactologoID int64, fechaHora time.Time, duracionMinutos int, excludeTurnoID *int64) (bool, error) {
	finTurno := fechaHora.Add(time.Duration(duracionMinutos) * time.Minute)

	query := `
		SELECT COUNT(*)
		FROM turnos
		WHERE contactologo_id = $1
		  AND estado NOT IN ($2, $3)
		  AND (
		    (fecha_hora < $4 AND fecha_hora + (duracion_minutos || ' minutes')::INTERVAL > $5)
		    OR
		    (fecha_hora >= $5 AND fecha_hora < $4)
		  )
	`

	args := []interface{}{contactologoID, entities.EstadoCancelado, entities.EstadoNoAsistio, finTurno, fechaHora}

	if excludeTurnoID != nil {
		query += ` AND id != $6`
		args = append(args, *excludeTurnoID)
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *turnoRepository) Count(ctx context.Context, filter *entities.TurnoFilter) (int, error) {
	query := `SELECT COUNT(*) FROM turnos t WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if filter.PacienteID != nil {
		query += ` AND t.paciente_id = $` + string(rune(argPos+'0'))
		args = append(args, *filter.PacienteID)
		argPos++
	}

	if filter.ContactologoID != nil {
		query += ` AND t.contactologo_id = $` + string(rune(argPos+'0'))
		args = append(args, *filter.ContactologoID)
		argPos++
	}

	if filter.TipoServicio != nil {
		query += ` AND t.tipo_servicio = $` + string(rune(argPos+'0'))
		args = append(args, *filter.TipoServicio)
		argPos++
	}

	if filter.Estado != nil {
		query += ` AND t.estado = $` + string(rune(argPos+'0'))
		args = append(args, *filter.Estado)
		argPos++
	}

	if filter.FechaDesde != nil {
		query += ` AND t.fecha_hora >= $` + string(rune(argPos+'0'))
		args = append(args, *filter.FechaDesde)
		argPos++
	}

	if filter.FechaHasta != nil {
		query += ` AND t.fecha_hora <= $` + string(rune(argPos+'0'))
		args = append(args, *filter.FechaHasta)
		argPos++
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}
