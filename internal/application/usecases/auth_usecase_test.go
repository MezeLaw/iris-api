package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
	"iris-api/pkg/auth"
)

func TestAuthUseCase_Register(t *testing.T) {
	tests := []struct {
		name     string
		req      *entities.RegisterRequest
		behavior func(authRepo *MockAuthRepo, clientRepo *MockClientRepo)
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
			behavior: func(authRepo *MockAuthRepo, clientRepo *MockClientRepo) {
				clientRepo.On("Create", mock.Anything, mock.Anything).Return(&entities.Client{
					ID:   1,
					Name: "Test Clinic",
					Slug: "testclinic",
				}, nil)
				authRepo.On("Register", mock.Anything, mock.Anything).Return(&entities.AuthUser{
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
			},
		},
		{
			name: "error - service fails",
			req: &entities.RegisterRequest{
				ClientName:  "Test Clinic",
				ClientSlug:  "testclinic",
				ClientEmail: "clinic@example.com",
				Email:       "user@example.com",
				Password:    "securepassword",
				Name:        "Test User",
				Role:        entities.RoleAdmin,
			},
			behavior: func(authRepo *MockAuthRepo, clientRepo *MockClientRepo) {
				clientRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, user *entities.AuthUser, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo := new(MockAuthRepo)
			mockClientRepo := new(MockClientRepo)
			tt.behavior(mockAuthRepo, mockClientRepo)

			jwtManager := auth.NewJWTManager("test-secret", time.Hour)
			authService := services.NewAuthService(mockAuthRepo, mockClientRepo, jwtManager)
			useCase := NewAuthUseCase(authService)

			user, err := useCase.Register(context.Background(), tt.req)

			tt.asserts(t, user, err)
			mockAuthRepo.AssertExpectations(t)
			mockClientRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	hashedPassword, _ := auth.HashPassword("correctpassword")

	tests := []struct {
		name     string
		req      *entities.LoginRequest
		behavior func(authRepo *MockAuthRepo)
		asserts  func(t *testing.T, resp *entities.LoginResponse, err error)
	}{
		{
			name: "success - logs in with valid credentials",
			req: &entities.LoginRequest{
				Email:    "user@example.com",
				Password: "correctpassword",
			},
			behavior: func(authRepo *MockAuthRepo) {
				authRepo.On("FindByEmail", mock.Anything, "user@example.com").Return(&entities.AuthUser{
					ID:       1,
					ClientID: 100,
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
				assert.Empty(t, resp.User.Password)
			},
		},
		{
			name: "error - invalid credentials",
			req: &entities.LoginRequest{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			behavior: func(authRepo *MockAuthRepo) {
				authRepo.On("FindByEmail", mock.Anything, "user@example.com").Return(&entities.AuthUser{
					ID:       1,
					Email:    "user@example.com",
					Password: hashedPassword,
				}, nil)
			},
			asserts: func(t *testing.T, resp *entities.LoginResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
		{
			name: "error - user not found",
			req: &entities.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password",
			},
			behavior: func(authRepo *MockAuthRepo) {
				authRepo.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, resp *entities.LoginResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo := new(MockAuthRepo)
			mockClientRepo := new(MockClientRepo)
			tt.behavior(mockAuthRepo)

			jwtManager := auth.NewJWTManager("test-secret", time.Hour)
			authService := services.NewAuthService(mockAuthRepo, mockClientRepo, jwtManager)
			useCase := NewAuthUseCase(authService)

			resp, err := useCase.Login(context.Background(), tt.req)

			tt.asserts(t, resp, err)
			mockAuthRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_ValidateToken(t *testing.T) {
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
			},
		},
		{
			name:     "error - invalid token",
			token:    "invalid.token.string",
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

			mockAuthRepo := new(MockAuthRepo)
			mockClientRepo := new(MockClientRepo)
			authService := services.NewAuthService(mockAuthRepo, mockClientRepo, jwtManager)
			useCase := NewAuthUseCase(authService)

			claims, err := useCase.ValidateToken(tt.token)

			tt.asserts(t, claims, err)
		})
	}
}

func TestAuthUseCase_GetUserByID(t *testing.T) {
	tests := []struct {
		name     string
		userID   int64
		behavior func(authRepo *MockAuthRepo)
		asserts  func(t *testing.T, user *entities.AuthUser, err error)
	}{
		{
			name:   "success - returns user by id",
			userID: 1,
			behavior: func(authRepo *MockAuthRepo) {
				authRepo.On("FindByID", mock.Anything, int64(1)).Return(&entities.AuthUser{
					ID:       1,
					ClientID: 100,
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
			behavior: func(authRepo *MockAuthRepo) {
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
			mockAuthRepo := new(MockAuthRepo)
			mockClientRepo := new(MockClientRepo)
			tt.behavior(mockAuthRepo)

			jwtManager := auth.NewJWTManager("test-secret", time.Hour)
			authService := services.NewAuthService(mockAuthRepo, mockClientRepo, jwtManager)
			useCase := NewAuthUseCase(authService)

			user, err := useCase.GetUserByID(context.Background(), tt.userID)

			tt.asserts(t, user, err)
			mockAuthRepo.AssertExpectations(t)
		})
	}
}

// MockAuthRepo is a mock implementation of AuthRepository
type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) Register(ctx context.Context, user *entities.AuthUser) (*entities.AuthUser, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuthUser), args.Error(1)
}

func (m *MockAuthRepo) FindByEmail(ctx context.Context, email string) (*entities.AuthUser, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuthUser), args.Error(1)
}

func (m *MockAuthRepo) FindByID(ctx context.Context, id int64) (*entities.AuthUser, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuthUser), args.Error(1)
}

func (m *MockAuthRepo) UpdateLastLogin(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockClientRepo is a mock implementation of ClientRepository
type MockClientRepo struct {
	mock.Mock
}

func (m *MockClientRepo) Create(ctx context.Context, client *entities.Client) (*entities.Client, error) {
	args := m.Called(ctx, client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *MockClientRepo) FindByID(ctx context.Context, id int64) (*entities.Client, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *MockClientRepo) FindBySlug(ctx context.Context, slug string) (*entities.Client, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}
