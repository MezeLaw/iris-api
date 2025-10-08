package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"iris-api/internal/domain/entities"
)

func TestReporteriaRepository_GetPacientesActivos(t *testing.T) {
	tests := []struct {
		name     string
		clientID int64
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, pacientes []*entities.PacienteActivo, err error)
	}{
		{
			name:     "success - returns active patients",
			clientID: 1,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "ultimo_turno", "total_turnos", "turnos_pendientes"}).
					AddRow(1, "Patient 1", "patient1@example.com", time.Now(), 5, 2).
					AddRow(2, "Patient 2", "patient2@example.com", time.Now(), 3, 1)
				mock.ExpectQuery(`SELECT (.+) FROM users u INNER JOIN turnos t`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, pacientes []*entities.PacienteActivo, err error) {
				assert.NoError(t, err)
				assert.Len(t, pacientes, 2)
				assert.Equal(t, int64(1), pacientes[0].ID)
			},
		},
		{
			name:     "success - no active patients",
			clientID: 1,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "ultimo_turno", "total_turnos", "turnos_pendientes"})
				mock.ExpectQuery(`SELECT (.+) FROM users u INNER JOIN turnos t`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, pacientes []*entities.PacienteActivo, err error) {
				assert.NoError(t, err)
				assert.Len(t, pacientes, 0)
			},
		},
		{
			name:     "error - database fails",
			clientID: 1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM users u INNER JOIN turnos t`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, pacientes []*entities.PacienteActivo, err error) {
				assert.Error(t, err)
				assert.Nil(t, pacientes)
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

			repo := NewReporteriaRepository(sqlxDB)
			result, err := repo.GetPacientesActivos(context.Background(), tt.clientID)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReporteriaRepository_GetPacientesInactivos(t *testing.T) {
	tests := []struct {
		name            string
		clientID        int64
		diasInactividad int
		behavior        func(mock sqlmock.Sqlmock)
		asserts         func(t *testing.T, pacientes []*entities.PacienteInactivo, err error)
	}{
		{
			name:            "success - returns inactive patients",
			clientID:        1,
			diasInactividad: 60,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "ultimo_turno", "dias_inactivo"}).
					AddRow(1, "Inactive Patient 1", "inactive1@example.com", time.Now().Add(-90*24*time.Hour), 90).
					AddRow(2, "Inactive Patient 2", "inactive2@example.com", time.Now().Add(-70*24*time.Hour), 70)
				mock.ExpectQuery(`SELECT (.+) FROM users u INNER JOIN turnos t`).
					WithArgs(int64(1), 60).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, pacientes []*entities.PacienteInactivo, err error) {
				assert.NoError(t, err)
				assert.Len(t, pacientes, 2)
				assert.Equal(t, 90, pacientes[0].DiasInactivo)
			},
		},
		{
			name:            "success - no inactive patients",
			clientID:        1,
			diasInactividad: 60,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "ultimo_turno", "dias_inactivo"})
				mock.ExpectQuery(`SELECT (.+) FROM users u INNER JOIN turnos t`).
					WithArgs(int64(1), 60).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, pacientes []*entities.PacienteInactivo, err error) {
				assert.NoError(t, err)
				assert.Len(t, pacientes, 0)
			},
		},
		{
			name:            "error - database fails",
			clientID:        1,
			diasInactividad: 60,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM users u INNER JOIN turnos t`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, pacientes []*entities.PacienteInactivo, err error) {
				assert.Error(t, err)
				assert.Nil(t, pacientes)
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

			repo := NewReporteriaRepository(sqlxDB)
			result, err := repo.GetPacientesInactivos(context.Background(), tt.clientID, tt.diasInactividad)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
