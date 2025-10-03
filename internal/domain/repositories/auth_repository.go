package repositories

import (
	"context"

	"iris-api/internal/domain/entities"
)

type AuthRepository interface {
	Register(ctx context.Context, user *entities.AuthUser) (*entities.AuthUser, error)
	FindByEmail(ctx context.Context, email string) (*entities.AuthUser, error)
	FindByID(ctx context.Context, id int64) (*entities.AuthUser, error)
	UpdateLastLogin(ctx context.Context, userID int64) error
}
