package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"iris-api/internal/domain/entities"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	tests := []struct {
		name     string
		user     *entities.AuthUser
		behavior func()
		asserts  func(t *testing.T, token string, expiresAt time.Time, err error)
	}{
		{
			name: "success - generates valid token",
			user: &entities.AuthUser{
				ID:       1,
				ClientID: 100,
				Email:    "test@example.com",
				Role:     entities.RoleAdmin,
			},
			behavior: func() {},
			asserts: func(t *testing.T, token string, expiresAt time.Time, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.True(t, expiresAt.After(time.Now()))
			},
		},
		{
			name: "success - generates token for optometrista",
			user: &entities.AuthUser{
				ID:       2,
				ClientID: 200,
				Email:    "optometrista@example.com",
				Role:     entities.RoleOptometrista,
			},
			behavior: func() {},
			asserts: func(t *testing.T, token string, expiresAt time.Time, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			},
		},
		{
			name: "success - generates token for recepcionista",
			user: &entities.AuthUser{
				ID:       3,
				ClientID: 300,
				Email:    "recepcionista@example.com",
				Role:     entities.RoleRecepcionista,
			},
			behavior: func() {},
			asserts: func(t *testing.T, token string, expiresAt time.Time, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behavior()

			manager := NewJWTManager("test-secret-key", time.Hour)
			token, expiresAt, err := manager.GenerateToken(tt.user)

			tt.asserts(t, token, expiresAt, err)
		})
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", time.Hour)
	shortManager := NewJWTManager("test-secret-key", 1*time.Millisecond)
	differentSecretManager := NewJWTManager("different-secret", time.Hour)

	validUser := &entities.AuthUser{
		ID:       1,
		ClientID: 100,
		Email:    "test@example.com",
		Role:     entities.RoleAdmin,
	}
	validToken, _, _ := manager.GenerateToken(validUser)

	// Generate expired token
	expiredToken, _, _ := shortManager.GenerateToken(validUser)
	time.Sleep(10 * time.Millisecond)

	// Generate token with different secret
	differentSecretToken, _, _ := differentSecretManager.GenerateToken(validUser)

	tests := []struct {
		name        string
		token       string
		useManager  JWTManager
		behavior    func()
		asserts     func(t *testing.T, claims *Claims, err error)
	}{
		{
			name:       "success - validates valid token",
			token:      validToken,
			useManager: manager,
			behavior:   func() {},
			asserts: func(t *testing.T, claims *Claims, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, int64(1), claims.UserID)
				assert.Equal(t, int64(100), claims.ClientID)
				assert.Equal(t, "test@example.com", claims.Email)
				assert.Equal(t, entities.RoleAdmin, claims.Role)
			},
		},
		{
			name:       "error - expired token",
			token:      expiredToken,
			useManager: shortManager,
			behavior:   func() {},
			asserts: func(t *testing.T, claims *Claims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
		{
			name:       "error - invalid signature (different secret)",
			token:      differentSecretToken,
			useManager: manager,
			behavior:   func() {},
			asserts: func(t *testing.T, claims *Claims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
		{
			name:       "error - malformed token",
			token:      "malformed.token.here",
			useManager: manager,
			behavior:   func() {},
			asserts: func(t *testing.T, claims *Claims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
		{
			name:       "error - empty token",
			token:      "",
			useManager: manager,
			behavior:   func() {},
			asserts: func(t *testing.T, claims *Claims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
		{
			name:       "error - random string",
			token:      "notavalidtoken",
			useManager: manager,
			behavior:   func() {},
			asserts: func(t *testing.T, claims *Claims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behavior()
			claims, err := tt.useManager.ValidateToken(tt.token)
			tt.asserts(t, claims, err)
		})
	}
}

func TestJWTManager_TokenExpiration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		behavior func()
		asserts  func(t *testing.T, expiresAt time.Time)
	}{
		{
			name:     "token expires in 1 hour",
			duration: time.Hour,
			behavior: func() {},
			asserts: func(t *testing.T, expiresAt time.Time) {
				expectedExpiry := time.Now().Add(time.Hour)
				assert.WithinDuration(t, expectedExpiry, expiresAt, 5*time.Second)
			},
		},
		{
			name:     "token expires in 24 hours",
			duration: 24 * time.Hour,
			behavior: func() {},
			asserts: func(t *testing.T, expiresAt time.Time) {
				expectedExpiry := time.Now().Add(24 * time.Hour)
				assert.WithinDuration(t, expectedExpiry, expiresAt, 5*time.Second)
			},
		},
		{
			name:     "token expires in 30 minutes",
			duration: 30 * time.Minute,
			behavior: func() {},
			asserts: func(t *testing.T, expiresAt time.Time) {
				expectedExpiry := time.Now().Add(30 * time.Minute)
				assert.WithinDuration(t, expectedExpiry, expiresAt, 5*time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behavior()

			manager := NewJWTManager("test-secret", tt.duration)
			user := &entities.AuthUser{
				ID:       1,
				ClientID: 1,
				Email:    "test@example.com",
				Role:     entities.RoleAdmin,
			}

			_, expiresAt, err := manager.GenerateToken(user)
			assert.NoError(t, err)
			tt.asserts(t, expiresAt)
		})
	}
}
