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

func TestClientRepository_Create(t *testing.T) {
	tests := []struct {
		name     string
		client   *entities.Client
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, client *entities.Client, err error)
	}{
		{
			name: "success - creates client",
			client: &entities.Client{
				Name:      "Test Clinic",
				Slug:      "test-clinic",
				Email:     "clinic@example.com",
				Phone:     "1234567890",
				Address:   "Test Address",
				Active:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow(1, time.Now(), time.Now())
				mock.ExpectQuery(`INSERT INTO clients`).
					WithArgs("Test Clinic", "test-clinic", "clinic@example.com", "1234567890", "Test Address", true, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, client *entities.Client, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, int64(1), client.ID)
			},
		},
		{
			name:   "error - database fails",
			client: &entities.Client{Name: "Test"},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO clients`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, client *entities.Client, err error) {
				assert.Error(t, err)
				assert.Nil(t, client)
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

			repo := NewClientRepository(sqlxDB)
			result, err := repo.Create(context.Background(), tt.client)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestClientRepository_FindBySlug(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, client *entities.Client, err error)
	}{
		{
			name: "success - finds client by slug",
			slug: "test-clinic",
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "slug", "email", "phone", "address", "active", "created_at", "updated_at"}).
					AddRow(1, "Test Clinic", "test-clinic", "clinic@example.com", "123", "Address", true, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM clients WHERE slug`).
					WithArgs("test-clinic").
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, client *entities.Client, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, "test-clinic", client.Slug)
			},
		},
		{
			name: "error - client not found",
			slug: "notfound",
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM clients WHERE slug`).
					WithArgs("notfound").
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, client *entities.Client, err error) {
				assert.Error(t, err)
				assert.Nil(t, client)
				assert.Contains(t, err.Error(), "client not found")
			},
		},
		{
			name: "error - database fails",
			slug: "test",
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM clients WHERE slug`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, client *entities.Client, err error) {
				assert.Error(t, err)
				assert.Nil(t, client)
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

			repo := NewClientRepository(sqlxDB)
			result, err := repo.FindBySlug(context.Background(), tt.slug)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestClientRepository_FindByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, client *entities.Client, err error)
	}{
		{
			name: "success - finds client by id",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "slug", "email", "phone", "address", "active", "created_at", "updated_at"}).
					AddRow(1, "Test Clinic", "test-clinic", "clinic@example.com", "123", "Address", true, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM clients WHERE id`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, client *entities.Client, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, int64(1), client.ID)
			},
		},
		{
			name: "error - client not found",
			id:   999,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM clients WHERE id`).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, client *entities.Client, err error) {
				assert.Error(t, err)
				assert.Nil(t, client)
				assert.Contains(t, err.Error(), "client not found")
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

			repo := NewClientRepository(sqlxDB)
			result, err := repo.FindByID(context.Background(), tt.id)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
