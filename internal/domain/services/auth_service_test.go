package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/domain/entities"
	"iris-api/pkg/auth"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name     string
		req      *entities.RegisterRequest
		behavior func(clientRepo *MockClientRepository, authRepo *MockAuthRepository)
		asserts  func(t *testing.T, user *entities.AuthUser, err error)
	}{
		{
			name: "success - registers new user and client",
			req: &entities.RegisterRequest{
				ClientName:  "Test Clinic",
				ClientSlug:  "testclinic",
				ClientEmail: "clinic@example.com",
				ClientPhone: "1234567890",
				Email:       "user@example.com",
				Password:    "securepassword",
				Name:        "Test User",
				Role:        entities.RoleAdmin,
			},
			behavior: func(clientRepo *MockClientRepository, authRepo *MockAuthRepository) {
				clientRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *entities.Client) bool {
					return c.Name == "Test Clinic" && c.Slug == "testclinic"
				})).Return(&entities.Client{
					ID:    1,
					Name:  "Test Clinic",
					Slug:  "testclinic",
					Email: "clinic@example.com",
					Phone: "1234567890",
				}, nil)

				authRepo.On("Register", mock.Anything, mock.MatchedBy(func(u *entities.AuthUser) bool {
					return u.ClientID == 1 && u.Email == "user@example.com" && u.Name == "Test User"
				})).Return(&entities.AuthUser{
					ID:       1,
					ClientID: 1,
					Email:    "user@example.com",
					Name:     "Test User",
					Role:     entities.RoleAdmin,
				}, nil)
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, int64(1), user.ClientID)
				assert.Equal(t, "user@example.com", user.Email)
				assert.Equal(t, entities.RoleAdmin, user.Role)
			},
		},
		{
			name: "error - client creation fails",
			req: &entities.RegisterRequest{
				ClientName:  "Test Clinic",
				ClientSlug:  "testclinic",
				ClientEmail: "clinic@example.com",
				Email:       "user@example.com",
				Password:    "securepassword",
				Name:        "Test User",
				Role:        entities.RoleAdmin,
			},
			behavior: func(clientRepo *MockClientRepository, authRepo *MockAuthRepository) {
				clientRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "failed to create client")
			},
		},
		{
			name: "error - user registration fails",
			req: &entities.RegisterRequest{
				ClientName:  "Test Clinic",
				ClientSlug:  "testclinic",
				ClientEmail: "clinic@example.com",
				Email:       "user@example.com",
				Password:    "securepassword",
				Name:        "Test User",
				Role:        entities.RoleAdmin,
			},
			behavior: func(clientRepo *MockClientRepository, authRepo *MockAuthRepository) {
				clientRepo.On("Create", mock.Anything, mock.Anything).Return(&entities.Client{
					ID: 1,
				}, nil)
				authRepo.On("Register", mock.Anything, mock.Anything).Return(nil, errors.New("email already exists"))
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClientRepo := new(MockClientRepository)
			mockAuthRepo := new(MockAuthRepository)
			tt.behavior(mockClientRepo, mockAuthRepo)

			jwtManager := auth.NewJWTManager("test-secret", time.Hour)
			service := NewAuthService(mockAuthRepo, mockClientRepo, jwtManager)

			user, err := service.Register(context.Background(), tt.req)

			tt.asserts(t, user, err)
			mockClientRepo.AssertExpectations(t)
			mockAuthRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	hashedPassword, _ := auth.HashPassword("correctpassword")

	tests := []struct {
		name     string
		req      *entities.LoginRequest
		behavior func(authRepo *MockAuthRepository)
		asserts  func(t *testing.T, resp *entities.LoginResponse, err error)
	}{
		{
			name: "success - logs in with valid credentials",
			req: &entities.LoginRequest{
				Email:    "user@example.com",
				Password: "correctpassword",
			},
			behavior: func(authRepo *MockAuthRepository) {
				authRepo.On("FindByEmail", mock.Anything, "user@example.com").Return(&entities.AuthUser{
					ID:       1,
					ClientID: 1,
					Email:    "user@example.com",
					Password: hashedPassword,
					Name:     "Test User",
					Role:     entities.RoleAdmin,
				}, nil)
				authRepo.On("UpdateLastLogin", mock.Anything, int64(1)).Return(nil)
			},
			asserts: func(t *testing.T, resp *entities.LoginResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Token)
				assert.NotNil(t, resp.User)
				assert.Equal(t, "user@example.com", resp.User.Email)
				assert.Empty(t, resp.User.Password)
				assert.NotZero(t, resp.ExpiresAt)
			},
		},
		{
			name: "error - user not found",
			req: &entities.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password",
			},
			behavior: func(authRepo *MockAuthRepository) {
				authRepo.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, resp *entities.LoginResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), "invalid credentials")
			},
		},
		{
			name: "error - wrong password",
			req: &entities.LoginRequest{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			behavior: func(authRepo *MockAuthRepository) {
				authRepo.On("FindByEmail", mock.Anything, "user@example.com").Return(&entities.AuthUser{
					ID:       1,
					Email:    "user@example.com",
					Password: hashedPassword,
				}, nil)
			},
			asserts: func(t *testing.T, resp *entities.LoginResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), "invalid credentials")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo := new(MockAuthRepository)
			mockClientRepo := new(MockClientRepository)
			tt.behavior(mockAuthRepo)

			jwtManager := auth.NewJWTManager("test-secret", time.Hour)
			service := NewAuthService(mockAuthRepo, mockClientRepo, jwtManager)

			resp, err := service.Login(context.Background(), tt.req)

			tt.asserts(t, resp, err)
			mockAuthRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", time.Hour)

	validUser := &entities.AuthUser{
		ID:       1,
		ClientID: 100,
		Email:    "test@example.com",
		Role:     entities.RoleAdmin,
	}
	validToken, _, _ := jwtManager.GenerateToken(validUser)

	tests := []struct {
		name     string
		token    string
		behavior func()
		asserts  func(t *testing.T, claims *entities.TokenClaims, err error)
	}{
		{
			name:     "success - validates valid token",
			token:    validToken,
			behavior: func() {},
			asserts: func(t *testing.T, claims *entities.TokenClaims, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, int64(1), claims.UserID)
				assert.Equal(t, int64(100), claims.ClientID)
				assert.Equal(t, "test@example.com", claims.Email)
				assert.Equal(t, entities.RoleAdmin, claims.Role)
			},
		},
		{
			name:     "error - invalid token format",
			token:    "invalid.token.string",
			behavior: func() {},
			asserts: func(t *testing.T, claims *entities.TokenClaims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
		{
			name:     "error - empty token",
			token:    "",
			behavior: func() {},
			asserts: func(t *testing.T, claims *entities.TokenClaims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behavior()

			mockAuthRepo := new(MockAuthRepository)
			mockClientRepo := new(MockClientRepository)
			service := NewAuthService(mockAuthRepo, mockClientRepo, jwtManager)

			claims, err := service.ValidateToken(tt.token)

			tt.asserts(t, claims, err)
		})
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	tests := []struct {
		name     string
		userID   int64
		behavior func(authRepo *MockAuthRepository)
		asserts  func(t *testing.T, user *entities.AuthUser, err error)
	}{
		{
			name:   "success - returns user by id",
			userID: 1,
			behavior: func(authRepo *MockAuthRepository) {
				authRepo.On("FindByID", mock.Anything, int64(1)).Return(&entities.AuthUser{
					ID:       1,
					ClientID: 1,
					Email:    "test@example.com",
					Password: "hashedpassword",
					Name:     "Test User",
				}, nil)
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, int64(1), user.ID)
				assert.Empty(t, user.Password)
			},
		},
		{
			name:   "error - user not found",
			userID: 999,
			behavior: func(authRepo *MockAuthRepository) {
				authRepo.On("FindByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo := new(MockAuthRepository)
			mockClientRepo := new(MockClientRepository)
			tt.behavior(mockAuthRepo)

			jwtManager := auth.NewJWTManager("test-secret", time.Hour)
			service := NewAuthService(mockAuthRepo, mockClientRepo, jwtManager)

			user, err := service.GetUserByID(context.Background(), tt.userID)

			tt.asserts(t, user, err)
			mockAuthRepo.AssertExpectations(t)
		})
	}
}

// MockAuthRepository is a mock implementation of AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Register(ctx context.Context, user *entities.AuthUser) (*entities.AuthUser, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuthUser), args.Error(1)
}

func (m *MockAuthRepository) FindByEmail(ctx context.Context, email string) (*entities.AuthUser, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuthUser), args.Error(1)
}

func (m *MockAuthRepository) FindByID(ctx context.Context, id int64) (*entities.AuthUser, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuthUser), args.Error(1)
}

func (m *MockAuthRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockClientRepository is a mock implementation of ClientRepository
type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) Create(ctx context.Context, client *entities.Client) (*entities.Client, error) {
	args := m.Called(ctx, client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *MockClientRepository) FindByID(ctx context.Context, id int64) (*entities.Client, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *MockClientRepository) FindBySlug(ctx context.Context, slug string) (*entities.Client, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}
