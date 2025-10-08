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

func TestRecetaRepository_Create(t *testing.T) {
	tests := []struct {
		name     string
		receta   *entities.Receta
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, receta *entities.Receta, err error)
	}{
		{
			name: "success - creates receta",
			receta: &entities.Receta{
				PacienteID:  1,
				Fecha:       time.Now(),
				ODEsfera:    -2.5,
				ODCilindro:  -1.0,
				ODEje:       90,
				OIEsfera:    -3.0,
				OICilindro:  -0.5,
				OIEje:       180,
				TipoLente:   "monofocal",
				Observaciones: "Test observation",
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "paciente_id", "fecha", "od_esfera", "od_cilindro", "od_eje", "oi_esfera", "oi_cilindro", "oi_eje", "tipo_lente", "observaciones", "created_at", "updated_at"}).
					AddRow(1, 1, time.Now(), -2.5, -1.0, 90, -3.0, -0.5, 180, "monofocal", "Test observation", time.Now(), time.Now())
				mock.ExpectQuery(`INSERT INTO recetas`).
					WithArgs(int64(1), sqlmock.AnyArg(), -2.5, -1.0, 90, -3.0, -0.5, 180, "monofocal", "Test observation", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, receta)
				assert.Equal(t, int64(1), receta.ID)
			},
		},
		{
			name:   "error - database fails",
			receta: &entities.Receta{PacienteID: 1},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO recetas`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
				assert.Contains(t, err.Error(), "failed to create receta")
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

			repo := NewRecetaRepository(sqlxDB)
			result, err := repo.Create(context.Background(), tt.receta)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecetaRepository_GetByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, receta *entities.Receta, err error)
	}{
		{
			name: "success - gets receta by id",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "paciente_id", "fecha", "od_esfera", "od_cilindro", "od_eje", "oi_esfera", "oi_cilindro", "oi_eje", "tipo_lente", "observaciones", "created_at", "updated_at"}).
					AddRow(1, 1, time.Now(), -2.5, -1.0, 90, -3.0, -0.5, 180, "monofocal", "Test", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM recetas WHERE id`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, receta)
				assert.Equal(t, int64(1), receta.ID)
			},
		},
		{
			name: "error - receta not found",
			id:   999,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM recetas WHERE id`).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
				assert.Contains(t, err.Error(), "receta not found")
			},
		},
		{
			name: "error - database fails",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM recetas WHERE id`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
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

			repo := NewRecetaRepository(sqlxDB)
			result, err := repo.GetByID(context.Background(), tt.id)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecetaRepository_GetAll(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		offset   int
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, recetas []*entities.Receta, err error)
	}{
		{
			name:   "success - gets all recetas",
			limit:  10,
			offset: 0,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "paciente_id", "fecha", "od_esfera", "od_cilindro", "od_eje", "oi_esfera", "oi_cilindro", "oi_eje", "tipo_lente", "observaciones", "created_at", "updated_at"}).
					AddRow(1, 1, time.Now(), -2.5, -1.0, 90, -3.0, -0.5, 180, "monofocal", "Test", time.Now(), time.Now()).
					AddRow(2, 2, time.Now(), -1.5, -0.5, 90, -2.0, -0.25, 180, "bifocal", "Test 2", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM recetas ORDER BY fecha DESC, created_at DESC`).
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, err error) {
				assert.NoError(t, err)
				assert.Len(t, recetas, 2)
			},
		},
		{
			name:   "error - database fails",
			limit:  10,
			offset: 0,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM recetas`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, recetas)
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

			repo := NewRecetaRepository(sqlxDB)
			result, err := repo.GetAll(context.Background(), tt.limit, tt.offset)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecetaRepository_Delete(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, err error)
	}{
		{
			name: "success - deletes receta",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM recetas WHERE id`).
					WithArgs(int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error - receta not found",
			id:   999,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM recetas WHERE id`).
					WithArgs(int64(999)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "receta not found")
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

			repo := NewRecetaRepository(sqlxDB)
			err = repo.Delete(context.Background(), tt.id)

			tt.asserts(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecetaRepository_Count(t *testing.T) {
	tests := []struct {
		name     string
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, count int, err error)
	}{
		{
			name: "success - counts recetas",
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(10)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM recetas`).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, count int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 10, count)
			},
		},
		{
			name: "error - database fails",
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

			repo := NewRecetaRepository(sqlxDB)
			result, err := repo.Count(context.Background())

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecetaRepository_GetByPacienteID(t *testing.T) {
	tests := []struct {
		name       string
		pacienteID int64
		limit      int
		offset     int
		behavior   func(mock sqlmock.Sqlmock)
		asserts    func(t *testing.T, recetas []*entities.Receta, err error)
	}{
		{
			name:       "success - gets recetas by paciente",
			pacienteID: 1,
			limit:      10,
			offset:     0,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "paciente_id", "fecha", "od_esfera", "od_cilindro", "od_eje", "oi_esfera", "oi_cilindro", "oi_eje", "tipo_lente", "observaciones", "created_at", "updated_at"}).
					AddRow(1, 1, time.Now(), -2.5, -1.0, 90, -3.0, -0.5, 180, "monofocal", "Test", time.Now(), time.Now()).
					AddRow(2, 1, time.Now(), -2.0, -0.5, 90, -2.5, -0.25, 180, "bifocal", "Test 2", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM recetas WHERE paciente_id`).
					WithArgs(int64(1), 10, 0).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, err error) {
				assert.NoError(t, err)
				assert.Len(t, recetas, 2)
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

			repo := NewRecetaRepository(sqlxDB)
			result, err := repo.GetByPacienteID(context.Background(), tt.pacienteID, tt.limit, tt.offset)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecetaRepository_GetHistorialByPacienteID(t *testing.T) {
	tests := []struct {
		name       string
		pacienteID int64
		behavior   func(mock sqlmock.Sqlmock)
		asserts    func(t *testing.T, recetas []*entities.Receta, err error)
	}{
		{
			name:       "success - gets historial by paciente",
			pacienteID: 1,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "paciente_id", "fecha", "od_esfera", "od_cilindro", "od_eje", "oi_esfera", "oi_cilindro", "oi_eje", "tipo_lente", "observaciones", "created_at", "updated_at"}).
					AddRow(1, 1, time.Now(), -2.5, -1.0, 90, -3.0, -0.5, 180, "monofocal", "Test", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM recetas WHERE paciente_id`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, err error) {
				assert.NoError(t, err)
				assert.Len(t, recetas, 1)
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

			repo := NewRecetaRepository(sqlxDB)
			result, err := repo.GetHistorialByPacienteID(context.Background(), tt.pacienteID)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecetaRepository_GetLastTwoByPacienteID(t *testing.T) {
	tests := []struct {
		name       string
		pacienteID int64
		behavior   func(mock sqlmock.Sqlmock)
		asserts    func(t *testing.T, recetas []*entities.Receta, err error)
	}{
		{
			name:       "success - gets last two recetas",
			pacienteID: 1,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "paciente_id", "fecha", "od_esfera", "od_cilindro", "od_eje", "oi_esfera", "oi_cilindro", "oi_eje", "tipo_lente", "observaciones", "created_at", "updated_at"}).
					AddRow(1, 1, time.Now(), -2.5, -1.0, 90, -3.0, -0.5, 180, "monofocal", "Test", time.Now(), time.Now()).
					AddRow(2, 1, time.Now(), -2.0, -0.5, 90, -2.5, -0.25, 180, "bifocal", "Test 2", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM recetas WHERE paciente_id (.+) LIMIT 2`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, recetas []*entities.Receta, err error) {
				assert.NoError(t, err)
				assert.Len(t, recetas, 2)
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

			repo := NewRecetaRepository(sqlxDB)
			result, err := repo.GetLastTwoByPacienteID(context.Background(), tt.pacienteID)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecetaRepository_Update(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		receta   *entities.Receta
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, receta *entities.Receta, err error)
	}{
		{
			name: "success - updates receta",
			id:   1,
			receta: &entities.Receta{
				PacienteID:     1,
				Fecha:          time.Now(),
				ODEsfera:       -3.0,
				ODCilindro:     -1.5,
				ODEje:          90,
				OIEsfera:       -3.5,
				OICilindro:     -1.0,
				OIEje:          180,
				TipoLente:      "progresivo",
				Observaciones:  "Updated",
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "paciente_id", "fecha", "od_esfera", "od_cilindro", "od_eje", "oi_esfera", "oi_cilindro", "oi_eje", "tipo_lente", "observaciones", "created_at", "updated_at"}).
					AddRow(1, 1, time.Now(), -3.0, -1.5, 90, -3.5, -1.0, 180, "progresivo", "Updated", time.Now(), time.Now())
				mock.ExpectQuery(`UPDATE recetas SET (.+) WHERE id`).
					WithArgs(int64(1), sqlmock.AnyArg(), -3.0, -1.5, 90, -3.5, -1.0, 180, "progresivo", "Updated", sqlmock.AnyArg(), int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, receta)
				assert.Equal(t, int64(1), receta.ID)
			},
		},
		{
			name:   "error - receta not found",
			id:     999,
			receta: &entities.Receta{},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE recetas SET (.+) WHERE id`).
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, receta *entities.Receta, err error) {
				assert.Error(t, err)
				assert.Nil(t, receta)
				assert.Contains(t, err.Error(), "receta not found")
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

			repo := NewRecetaRepository(sqlxDB)
			result, err := repo.Update(context.Background(), tt.id, tt.receta)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecetaRepository_CountByPacienteID(t *testing.T) {
	tests := []struct {
		name       string
		pacienteID int64
		behavior   func(mock sqlmock.Sqlmock)
		asserts    func(t *testing.T, count int, err error)
	}{
		{
			name:       "success - counts recetas by paciente",
			pacienteID: 1,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM recetas WHERE paciente_id`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, count int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 5, count)
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

			repo := NewRecetaRepository(sqlxDB)
			result, err := repo.CountByPacienteID(context.Background(), tt.pacienteID)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
