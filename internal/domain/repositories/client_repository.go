package repositories

import (
	"context"

	"iris-api/internal/domain/entities"
)

type ClientRepository interface {
	Create(ctx context.Context, client *entities.Client) (*entities.Client, error)
	FindBySlug(ctx context.Context, slug string) (*entities.Client, error)
	FindByID(ctx context.Context, id int64) (*entities.Client, error)
}
