package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		behavior func()
		asserts  func(t *testing.T, hash string, err error)
	}{
		{
			name:     "success - hashes password",
			password: "securepassword123",
			behavior: func() {},
			asserts: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, "securepassword123", hash)
				assert.True(t, len(hash) > 0)
			},
		},
		{
			name:     "success - different passwords produce different hashes",
			password: "password1",
			behavior: func() {},
			asserts: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err)
				hash2, err2 := HashPassword("password2")
				assert.NoError(t, err2)
				assert.NotEqual(t, hash, hash2)
			},
		},
		{
			name:     "success - same password produces different hashes due to salt",
			password: "samepassword",
			behavior: func() {},
			asserts: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err)
				hash2, err2 := HashPassword("samepassword")
				assert.NoError(t, err2)
				assert.NotEqual(t, hash, hash2)
			},
		},
		{
			name:     "success - hashes empty string",
			password: "",
			behavior: func() {},
			asserts: func(t *testing.T, hash string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behavior()
			hash, err := HashPassword(tt.password)
			tt.asserts(t, hash, err)
		})
	}
}

func TestCheckPassword(t *testing.T) {
	validHash, _ := HashPassword("correctpassword")

	tests := []struct {
		name     string
		password string
		hash     string
		behavior func()
		asserts  func(t *testing.T, result bool)
	}{
		{
			name:     "success - correct password matches hash",
			password: "correctpassword",
			hash:     validHash,
			behavior: func() {},
			asserts: func(t *testing.T, result bool) {
				assert.True(t, result)
			},
		},
		{
			name:     "failure - incorrect password does not match",
			password: "wrongpassword",
			hash:     validHash,
			behavior: func() {},
			asserts: func(t *testing.T, result bool) {
				assert.False(t, result)
			},
		},
		{
			name:     "failure - empty password does not match",
			password: "",
			hash:     validHash,
			behavior: func() {},
			asserts: func(t *testing.T, result bool) {
				assert.False(t, result)
			},
		},
		{
			name:     "failure - invalid hash format",
			password: "password",
			hash:     "invalid-hash",
			behavior: func() {},
			asserts: func(t *testing.T, result bool) {
				assert.False(t, result)
			},
		},
		{
			name:     "failure - empty hash",
			password: "password",
			hash:     "",
			behavior: func() {},
			asserts: func(t *testing.T, result bool) {
				assert.False(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behavior()
			result := CheckPassword(tt.password, tt.hash)
			tt.asserts(t, result)
		})
	}
}
