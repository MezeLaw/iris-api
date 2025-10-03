package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"iris-api/internal/domain/entities"
)

type clientRepository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) *clientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(ctx context.Context, client *entities.Client) (*entities.Client, error) {
	query := `
		INSERT INTO clients (name, slug, email, phone, address, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		client.Name,
		client.Slug,
		client.Email,
		client.Phone,
		client.Address,
		client.Active,
		client.CreatedAt,
		client.UpdatedAt,
	).Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *clientRepository) FindBySlug(ctx context.Context, slug string) (*entities.Client, error) {
	query := `
		SELECT id, name, slug, email, phone, address, active, created_at, updated_at
		FROM clients
		WHERE slug = $1 AND active = TRUE
	`

	var client entities.Client
	err := r.db.GetContext(ctx, &client, query, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("client not found")
		}
		return nil, err
	}

	return &client, nil
}

func (r *clientRepository) FindByID(ctx context.Context, id int64) (*entities.Client, error) {
	query := `
		SELECT id, name, slug, email, phone, address, active, created_at, updated_at
		FROM clients
		WHERE id = $1 AND active = TRUE
	`

	var client entities.Client
	err := r.db.GetContext(ctx, &client, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("client not found")
		}
		return nil, err
	}

	return &client, nil
}
