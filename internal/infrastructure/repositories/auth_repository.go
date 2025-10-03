package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"iris-api/internal/domain/entities"
)

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *authRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Register(ctx context.Context, user *entities.AuthUser) (*entities.AuthUser, error) {
	query := `
		INSERT INTO auth_users (client_id, email, password, name, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.ClientID,
		user.Email,
		user.Password,
		user.Name,
		user.Role,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*entities.AuthUser, error) {
	query := `
		SELECT id, client_id, email, password, name, role, active, created_at, updated_at
		FROM auth_users
		WHERE email = $1 AND active = TRUE
	`

	var user entities.AuthUser
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) FindByID(ctx context.Context, id int64) (*entities.AuthUser, error) {
	query := `
		SELECT id, client_id, email, password, name, role, active, created_at, updated_at
		FROM auth_users
		WHERE id = $1 AND active = TRUE
	`

	var user entities.AuthUser
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	query := `
		UPDATE auth_users
		SET updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}
