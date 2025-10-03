package entities

import "time"

type UserRole string

const (
	RoleAdmin         UserRole = "administrador"
	RoleOptometrista  UserRole = "optometrista"
	RoleRecepcionista UserRole = "recepcionista"
)

type AuthUser struct {
	ID        int64     `json:"id" db:"id"`
	ClientID  int64     `json:"client_id" db:"client_id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Name      string    `json:"name" db:"name"`
	Role      UserRole  `json:"role" db:"role"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	// Client info (for new client creation)
	ClientName  string `json:"client_name" validate:"required,min=2,max=200"`
	ClientSlug  string `json:"client_slug" validate:"required,min=2,max=100,alphanum"`
	ClientEmail string `json:"client_email" validate:"required,email"`
	ClientPhone string `json:"client_phone,omitempty" validate:"omitempty,max=20"`

	// User info
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Name     string   `json:"name" validate:"required,min=2,max=200"`
	Role     UserRole `json:"role" validate:"required,oneof=administrador optometrista recepcionista"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	User      *AuthUser `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TokenClaims struct {
	UserID   int64    `json:"user_id"`
	ClientID int64    `json:"client_id"`
	Email    string   `json:"email"`
	Role     UserRole `json:"role"`
}
