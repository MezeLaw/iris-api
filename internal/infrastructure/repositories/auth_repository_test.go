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

func TestAuthRepository_Register(t *testing.T) {
	tests := []struct {
		name     string
		user     *entities.AuthUser
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, user *entities.AuthUser, err error)
	}{
		{
			name: "success - registers user",
			user: &entities.AuthUser{
				ClientID:  1,
				Email:     "test@example.com",
				Password:  "hashedpassword",
				Name:      "Test User",
				Role:      entities.RoleAdmin,
				Active:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow(1, time.Now(), time.Now())
				mock.ExpectQuery(`INSERT INTO auth_users`).
					WithArgs(int64(1), "test@example.com", "hashedpassword", "Test User", entities.RoleAdmin, true, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, int64(1), user.ID)
			},
		},
		{
			name: "error - database fails",
			user: &entities.AuthUser{
				ClientID: 1,
				Email:    "test@example.com",
			},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO auth_users`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
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

			repo := NewAuthRepository(sqlxDB)
			result, err := repo.Register(context.Background(), tt.user)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthRepository_FindByEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, user *entities.AuthUser, err error)
	}{
		{
			name:  "success - finds user by email",
			email: "test@example.com",
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "client_id", "email", "password", "name", "role", "active", "created_at", "updated_at"}).
					AddRow(1, 1, "test@example.com", "hashedpass", "Test", entities.RoleAdmin, true, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM auth_users WHERE email`).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "test@example.com", user.Email)
			},
		},
		{
			name:  "error - user not found",
			email: "notfound@example.com",
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM auth_users WHERE email`).
					WithArgs("notfound@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "user not found")
			},
		},
		{
			name:  "error - database fails",
			email: "test@example.com",
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM auth_users WHERE email`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
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

			repo := NewAuthRepository(sqlxDB)
			result, err := repo.FindByEmail(context.Background(), tt.email)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthRepository_FindByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, user *entities.AuthUser, err error)
	}{
		{
			name: "success - finds user by id",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "client_id", "email", "password", "name", "role", "active", "created_at", "updated_at"}).
					AddRow(1, 1, "test@example.com", "hashedpass", "Test", entities.RoleAdmin, true, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM auth_users WHERE id`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, int64(1), user.ID)
			},
		},
		{
			name: "error - user not found",
			id:   999,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM auth_users WHERE id`).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "user not found")
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

			repo := NewAuthRepository(sqlxDB)
			result, err := repo.FindByID(context.Background(), tt.id)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthRepository_UpdateLastLogin(t *testing.T) {
	tests := []struct {
		name     string
		userID   int64
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, err error)
	}{
		{
			name:   "success - updates last login",
			userID: 1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE auth_users`).
					WithArgs(sqlmock.AnyArg(), int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "error - database fails",
			userID: 1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE auth_users`).
					WillReturnError(errors.New("database error"))
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

			repo := NewAuthRepository(sqlxDB)
			err = repo.UpdateLastLogin(context.Background(), tt.userID)

			tt.asserts(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
