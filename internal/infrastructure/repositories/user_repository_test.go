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

func TestUserRepository_Create(t *testing.T) {
	tests := []struct {
		name     string
		user     *entities.User
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, user *entities.User, err error)
	}{
		{
			name: "success - creates user",
			user: &entities.User{
				Email: "test@example.com",
				Name:  "Test User",
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
					AddRow(1, "test@example.com", "Test User", time.Now(), time.Now())
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs("test@example.com", "Test User", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, 1, user.ID)
				assert.Equal(t, "test@example.com", user.Email)
			},
		},
		{
			name: "error - database fails",
			user: &entities.User{
				Email: "test@example.com",
				Name:  "Test User",
			},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "failed to create user")
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

			repo := NewUserRepository(sqlxDB)
			result, err := repo.Create(context.Background(), tt.user)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, user *entities.User, err error)
	}{
		{
			name: "success - gets user by id",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
					AddRow(1, "test@example.com", "Test User", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM users WHERE id`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, 1, user.ID)
			},
		},
		{
			name: "error - user not found",
			id:   999,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM users WHERE id`).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "user not found")
			},
		},
		{
			name: "error - database fails",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM users WHERE id`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "failed to get user")
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

			repo := NewUserRepository(sqlxDB)
			result, err := repo.GetByID(context.Background(), tt.id)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, user *entities.User, err error)
	}{
		{
			name:  "success - gets user by email",
			email: "test@example.com",
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
					AddRow(1, "test@example.com", "Test User", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM users WHERE email`).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "test@example.com", user.Email)
			},
		},
		{
			name:  "error - user not found",
			email: "notfound@example.com",
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM users WHERE email`).
					WithArgs("notfound@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
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

			repo := NewUserRepository(sqlxDB)
			result, err := repo.GetByEmail(context.Background(), tt.email)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetAll(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		offset   int
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, users []*entities.User, err error)
	}{
		{
			name:   "success - gets all users",
			limit:  10,
			offset: 0,
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
					AddRow(1, "user1@example.com", "User 1", time.Now(), time.Now()).
					AddRow(2, "user2@example.com", "User 2", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM users ORDER BY created_at`).
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, users []*entities.User, err error) {
				assert.NoError(t, err)
				assert.Len(t, users, 2)
			},
		},
		{
			name:   "error - database fails",
			limit:  10,
			offset: 0,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM users`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, users []*entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, users)
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

			repo := NewUserRepository(sqlxDB)
			result, err := repo.GetAll(context.Background(), tt.limit, tt.offset)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		user     *entities.User
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, user *entities.User, err error)
	}{
		{
			name: "success - updates user",
			id:   1,
			user: &entities.User{
				Email: "updated@example.com",
				Name:  "Updated User",
			},
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
					AddRow(1, "updated@example.com", "Updated User", time.Now(), time.Now())
				mock.ExpectQuery(`UPDATE users`).
					WithArgs("updated@example.com", "Updated User", sqlmock.AnyArg(), 1).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "updated@example.com", user.Email)
			},
		},
		{
			name: "error - user not found",
			id:   999,
			user: &entities.User{Email: "test@example.com", Name: "Test"},
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE users`).
					WillReturnError(sql.ErrNoRows)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
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

			repo := NewUserRepository(sqlxDB)
			result, err := repo.Update(context.Background(), tt.id, tt.user)

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, err error)
	}{
		{
			name: "success - deletes user",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM users WHERE id`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error - user not found",
			id:   999,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM users WHERE id`).
					WithArgs(999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user not found")
			},
		},
		{
			name: "error - database fails",
			id:   1,
			behavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM users`).
					WillReturnError(errors.New("database error"))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete user")
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

			repo := NewUserRepository(sqlxDB)
			err = repo.Delete(context.Background(), tt.id)

			tt.asserts(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Count(t *testing.T) {
	tests := []struct {
		name     string
		behavior func(mock sqlmock.Sqlmock)
		asserts  func(t *testing.T, count int, err error)
	}{
		{
			name: "success - counts users",
			behavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, count int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 5, count)
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

			repo := NewUserRepository(sqlxDB)
			result, err := repo.Count(context.Background())

			tt.asserts(t, result, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
