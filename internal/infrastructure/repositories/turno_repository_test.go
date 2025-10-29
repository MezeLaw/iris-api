package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"iris-api/internal/domain/entities"
)

func TestTurnoRepository_Create(t *testing.T) {
	tests := []struct {
		name     string
		turno    *entities.Turno
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, turno *entities.Turno, err error)
	}{
		{
			name: "success - creates turno",
			turno: &entities.Turno{
				PacienteID:      1,
				ContactologoID:  2,
				FechaHora:       time.Now().Add(24 * time.Hour),
				DuracionMinutos: 30,
				TipoServicio:    entities.TipoServicio("consulta"),
				Estado:          entities.EstadoPendiente,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow(1, time.Now(), time.Now())
				mock.ExpectQuery(`INSERT INTO turnos`).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turno *entities.Turno, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, turno)
				assert.Equal(t, int64(1), turno.ID)
			},
		},
		{
			name:  "error - database fails",
			turno: &entities.Turno{PacienteID: 1},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO turnos`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, turno *entities.Turno, err error) {
				assert.Error(t, err)
				assert.Nil(t, turno)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.Create(context.Background(), tt.turno)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_GetByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, turno *entities.TurnoConDetalles, err error)
	}{
		{
			name: "success - gets turno by id",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient Name", "patient@example.com",
					"Doctor Name", "doctor@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN users p (.+) INNER JOIN users c (.+) WHERE t.id`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turno *entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, turno)
				assert.Equal(t, int64(1), turno.ID)
			},
		},
		{
			name: "error - turno not found",
			id:   999,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN users p (.+) INNER JOIN users c (.+) WHERE t.id`).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, turno *entities.TurnoConDetalles, err error) {
				assert.Error(t, err)
				assert.Nil(t, turno)
				assert.Contains(t, err.Error(), "turno not found")
			},
		},
		{
			name: "error - database generic error",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN users p (.+) INNER JOIN users c (.+) WHERE t.id`).
					WithArgs(int64(1)).
					WillReturnError(errors.New("database connection error"))
			},
			asserts: func(t *testing.T, turno *entities.TurnoConDetalles, err error) {
				assert.Error(t, err)
				assert.Nil(t, turno)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.GetByID(context.Background(), tt.id)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_Delete(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, err error)
	}{
		{
			name: "success - deletes turno",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM turnos WHERE id`).
					WithArgs(int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error - turno not found",
			id:   999,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM turnos WHERE id`).
					WithArgs(int64(999)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "turno not found")
			},
		},
		{
			name: "error - database exec fails",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM turnos WHERE id`).
					WithArgs(int64(1)).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "error - failed to get rows affected",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM turnos WHERE id`).
					WithArgs(int64(1)).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			err = repo.Delete(context.Background(), tt.id)

			tt.asserts(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_CancelTurno(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		motivo   string
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, err error)
	}{
		{
			name:   "success - cancels turno",
			id:     1,
			motivo: "Patient request",
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE turnos SET estado = (.+), observaciones = CONCAT`).
					WithArgs(entities.EstadoCancelado, "Patient request", sqlmock.AnyArg(), int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "error - database fails",
			id:     1,
			motivo: "Test",
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE turnos SET estado = (.+), observaciones = CONCAT`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:   "error - failed to get rows affected",
			id:     1,
			motivo: "Test",
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE turnos SET estado = (.+), observaciones = CONCAT`).
					WithArgs(entities.EstadoCancelado, "Test", sqlmock.AnyArg(), int64(1)).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:   "error - turno not found (0 rows affected)",
			id:     999,
			motivo: "Test",
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE turnos SET estado = (.+), observaciones = CONCAT`).
					WithArgs(entities.EstadoCancelado, "Test", sqlmock.AnyArg(), int64(999)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "turno not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			err = repo.CancelTurno(context.Background(), tt.id, tt.motivo)

			tt.asserts(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_CheckDisponibilidad(t *testing.T) {
	tests := []struct {
		name           string
		contactologoID int64
		fechaHora      time.Time
		duracion       int
		excludeID      *int64
		behavior       func(mock sqlmock.Sqlmock)
		asserts        func(t *testing.T, disponible bool, err error)
	}{
		{
			name:           "success - profesional disponible",
			contactologoID: 1,
			fechaHora:      time.Now().Add(24 * time.Hour),
			duracion:       30,
			excludeID:      nil,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM turnos`).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, disponible bool, err error) {
				assert.NoError(t, err)
				assert.True(t, disponible)
			},
		},
		{
			name:           "success - profesional no disponible",
			contactologoID: 1,
			fechaHora:      time.Now().Add(24 * time.Hour),
			duracion:       30,
			excludeID:      nil,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM turnos`).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, disponible bool, err error) {
				assert.NoError(t, err)
				assert.False(t, disponible)
			},
		},
		{
			name:           "success - disponible with excludeTurnoID",
			contactologoID: 1,
			fechaHora:      time.Now().Add(24 * time.Hour),
			duracion:       30,
			excludeID:      func() *int64 { v := int64(5); return &v }(),
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM turnos`).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, disponible bool, err error) {
				assert.NoError(t, err)
				assert.True(t, disponible)
			},
		},
		{
			name:           "error - database fails",
			contactologoID: 1,
			fechaHora:      time.Now().Add(24 * time.Hour),
			duracion:       30,
			excludeID:      nil,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM turnos`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, disponible bool, err error) {
				assert.Error(t, err)
				assert.False(t, disponible)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.CheckDisponibilidad(context.Background(), tt.contactologoID, tt.fechaHora, tt.duracion, tt.excludeID)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_Update(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		turno    *entities.Turno
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, turno *entities.Turno, err error)
	}{
		{
			name: "success - updates turno",
			id:   1,
			turno: &entities.Turno{
				PacienteID:      1,
				ContactologoID:  2,
				FechaHora:       time.Now(),
				DuracionMinutos: 45,
				TipoServicio:    entities.TipoServicio("seguimiento"),
				Estado:          entities.EstadoConfirmado,
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow(1, time.Now(), time.Now())
				mock.ExpectQuery(`UPDATE turnos SET (.+) WHERE id`).
					WithArgs(int64(1), int64(2), sqlmock.AnyArg(), 45, entities.TipoServicio("seguimiento"), entities.EstadoConfirmado, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turno *entities.Turno, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, turno)
				assert.Equal(t, int64(1), turno.ID)
			},
		},
		{
			name:  "error - turno not found",
			id:    999,
			turno: &entities.Turno{},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE turnos SET (.+) WHERE id`).
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, turno *entities.Turno, err error) {
				assert.Error(t, err)
				assert.Nil(t, turno)
				assert.Contains(t, err.Error(), "turno not found")
			},
		},
		{
			name:  "error - database generic error",
			id:    1,
			turno: &entities.Turno{},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE turnos SET (.+) WHERE id`).
					WillReturnError(errors.New("database connection error"))
			},
			asserts: func(t *testing.T, turno *entities.Turno, err error) {
				assert.Error(t, err)
				assert.Nil(t, turno)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.Update(context.Background(), tt.id, tt.turno)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_GetAll(t *testing.T) {
	tests := []struct {
		name     string
		filter   *entities.TurnoFilter
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, turnos []*entities.TurnoConDetalles, err error)
	}{
		{
			name: "success - gets all turnos",
			filter: &entities.TurnoFilter{
				Limit:  10,
				Offset: 0,
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				).AddRow(
					2, 2, 2, time.Now(), 30,
					"consulta", entities.EstadoConfirmado, "", "", false,
					time.Now(), time.Now(), "Patient 2", "patient2@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 2)
			},
		},
		{
			name: "success - filters by paciente",
			filter: &entities.TurnoFilter{
				PacienteID: func() *int64 { v := int64(1); return &v }(),
				Limit:      10,
				Offset:     0,
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
		{
			name: "success - filters by contactologo",
			filter: &entities.TurnoFilter{
				ContactologoID: func() *int64 { v := int64(2); return &v }(),
				Limit:          10,
				Offset:         0,
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
		{
			name: "success - filters by tipo_servicio",
			filter: &entities.TurnoFilter{
				TipoServicio: func() *entities.TipoServicio { v := entities.TipoServicio("consulta"); return &v }(),
				Limit:        10,
				Offset:       0,
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
		{
			name: "success - filters by estado",
			filter: &entities.TurnoFilter{
				Estado: func() *entities.EstadoTurno { v := entities.EstadoPendiente; return &v }(),
				Limit:  10,
				Offset: 0,
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
		{
			name: "success - filters by fecha range",
			filter: &entities.TurnoFilter{
				FechaDesde: func() *time.Time { v := time.Now().Add(-24 * time.Hour); return &v }(),
				FechaHasta: func() *time.Time { v := time.Now().Add(24 * time.Hour); return &v }(),
				Limit:      10,
				Offset:     0,
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
		{
			name: "error - database fails",
			filter: &entities.TurnoFilter{
				Limit:  10,
				Offset: 0,
			},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.Error(t, err)
				assert.Nil(t, turnos)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.GetAll(context.Background(), tt.filter)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_GetByDia(t *testing.T) {
	tests := []struct {
		name           string
		fecha          time.Time
		contactologoID *int64
		behavior       func(mock sqlmock.Sqlmock)
		asserts        func(t *testing.T, turnos []*entities.TurnoConDetalles, err error)
	}{
		{
			name:           "success - gets turnos by day",
			fecha:          time.Now(),
			contactologoID: nil,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, turnos)
			},
		},
		{
			name:           "success - gets turnos by day and profesional",
			fecha:          time.Now(),
			contactologoID: func() *int64 { v := int64(2); return &v }(),
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.GetByDia(context.Background(), tt.fecha, tt.contactologoID)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_GetBySemana(t *testing.T) {
	tests := []struct {
		name           string
		fechaInicio    time.Time
		contactologoID *int64
		behavior       func(mock sqlmock.Sqlmock)
		asserts        func(t *testing.T, turnos []*entities.TurnoConDetalles, err error)
	}{
		{
			name:           "success - gets turnos by week",
			fechaInicio:    time.Now(),
			contactologoID: nil,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, turnos)
			},
		},
		{
			name:           "error - database fails",
			fechaInicio:    time.Now(),
			contactologoID: nil,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.Error(t, err)
				assert.Nil(t, turnos)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.GetBySemana(context.Background(), tt.fechaInicio, tt.contactologoID)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_GetByProfesional(t *testing.T) {
	tests := []struct {
		name           string
		contactologoID int64
		fechaDesde     time.Time
		fechaHasta     time.Time
		behavior       func(mock sqlmock.Sqlmock)
		asserts        func(t *testing.T, turnos []*entities.TurnoConDetalles, err error)
	}{
		{
			name:           "success - gets turnos by profesional",
			contactologoID: 2,
			fechaDesde:     time.Now(),
			fechaHasta:     time.Now().Add(7 * 24 * time.Hour),
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now(), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
		{
			name:           "error - database fails",
			contactologoID: 2,
			fechaDesde:     time.Now(),
			fechaHasta:     time.Now().Add(7 * 24 * time.Hour),
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.Error(t, err)
				assert.Nil(t, turnos)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.GetByProfesional(context.Background(), tt.contactologoID, tt.fechaDesde, tt.fechaHasta)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_GetProximosTurnos(t *testing.T) {
	tests := []struct {
		name              string
		horasAnticipacion int
		behavior          func(mock sqlmock.Sqlmock)
		asserts           func(t *testing.T, turnos []*entities.TurnoConDetalles, err error)
	}{
		{
			name:              "success - gets upcoming turnos",
			horasAnticipacion: 24,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now().Add(2*time.Hour), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
		{
			name:              "error - database fails",
			horasAnticipacion: 24,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.Error(t, err)
				assert.Nil(t, turnos)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.GetProximosTurnos(context.Background(), tt.horasAnticipacion)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_GetTurnosSinConfirmar(t *testing.T) {
	tests := []struct {
		name     string
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, turnos []*entities.TurnoConDetalles, err error)
	}{
		{
			name: "success - gets unconfirmed turnos",
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "paciente_id", "contactologo_id", "fecha_hora", "duracion_minutos",
					"tipo_servicio", "estado", "motivo", "observaciones", "recordatorio_enviado",
					"created_at", "updated_at", "paciente_nombre", "paciente_email",
					"contactologo_nombre", "contactologo_email",
				}).AddRow(
					1, 1, 2, time.Now().Add(24*time.Hour), 30,
					"consulta", entities.EstadoPendiente, "", "", false,
					time.Now(), time.Now(), "Patient 1", "patient1@example.com",
					"Doctor 1", "doctor1@example.com",
				)
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.NoError(t, err)
				assert.Len(t, turnos, 1)
			},
		},
		{
			name: "error - database fails",
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM turnos t INNER JOIN`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, turnos []*entities.TurnoConDetalles, err error) {
				assert.Error(t, err)
				assert.Nil(t, turnos)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.GetTurnosSinConfirmar(context.Background())

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTurnoRepository_Count(t *testing.T) {
	tests := []struct {
		name     string
		filter   *entities.TurnoFilter
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, count int, err error)
	}{
		{
			name:   "success - counts all turnos",
			filter: &entities.TurnoFilter{},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(10)
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, count int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 10, count)
			},
		},
		{
			name: "success - counts with estado filter",
			filter: &entities.TurnoFilter{
				Estado: func() *entities.EstadoTurno { v := entities.EstadoPendiente; return &v }(),
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, count int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 5, count)
			},
		},
		{
			name: "success - counts with paciente_id filter",
			filter: &entities.TurnoFilter{
				PacienteID: func() *int64 { v := int64(1); return &v }(),
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(3)
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, count int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 3, count)
			},
		},
		{
			name: "success - counts with contactologo_id filter",
			filter: &entities.TurnoFilter{
				ContactologoID: func() *int64 { v := int64(2); return &v }(),
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(8)
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, count int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 8, count)
			},
		},
		{
			name: "success - counts with tipo_servicio filter",
			filter: &entities.TurnoFilter{
				TipoServicio: func() *entities.TipoServicio { v := entities.TipoServicio("consulta"); return &v }(),
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(12)
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, count int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 12, count)
			},
		},
		{
			name: "success - counts with fecha range filters",
			filter: &entities.TurnoFilter{
				FechaDesde: func() *time.Time { v := time.Now().Add(-24 * time.Hour); return &v }(),
				FechaHasta: func() *time.Time { v := time.Now().Add(24 * time.Hour); return &v }(),
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(6)
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(rows)
			},
			asserts: func(t *testing.T, count int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 6, count)
			},
		},
		{
			name:   "error - database fails",
			filter: &entities.TurnoFilter{},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, count int, err error) {
				assert.Error(t, err)
				assert.Equal(t, 0, count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			tt.behavior(mock)

			repo := NewTurnoRepository(sqlxDB)
			result, err := repo.Count(context.Background(), tt.filter)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
