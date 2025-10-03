package entities

import "time"

type Client struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone,omitempty" db:"phone"`
	Address   string    `json:"address,omitempty" db:"address"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateClientRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=200"`
	Slug    string `json:"slug" validate:"required,min=2,max=100,alphanum"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Address string `json:"address,omitempty" validate:"max=500"`
}

type UpdateClientRequest struct {
	Name    *string `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Email   *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone   *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Address *string `json:"address,omitempty" validate:"omitempty,max=500"`
	Active  *bool   `json:"active,omitempty"`
}
