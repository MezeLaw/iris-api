package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/application/usecases"
	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
	"iris-api/pkg/auth"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthMiddleware_RequireAuth(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		behavior   func()
		asserts    func(t *testing.T, resp *httptest.ResponseRecorder, contextCalled bool)
	}{
		{
			name:       "success - valid token",
			authHeader: "Bearer valid.token.here",
			behavior:   func() {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder, contextCalled bool) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.True(t, contextCalled)
			},
		},
		{
			name:       "error - missing authorization header",
			authHeader: "",
			behavior:   func() {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder, contextCalled bool) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
				assert.False(t, contextCalled)
			},
		},
		{
			name:       "error - invalid authorization header format (no Bearer)",
			authHeader: "token.without.bearer",
			behavior:   func() {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder, contextCalled bool) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
				assert.False(t, contextCalled)
			},
		},
		{
			name:       "error - invalid authorization header format (wrong prefix)",
			authHeader: "Basic sometoken",
			behavior:   func() {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder, contextCalled bool) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
				assert.False(t, contextCalled)
			},
		},
		{
			name:       "error - invalid token",
			authHeader: "Bearer invalid.token",
			behavior:   func() {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder, contextCalled bool) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
				assert.False(t, contextCalled)
			},
		},
		{
			name:       "error - expired token",
			authHeader: "Bearer expired.token",
			behavior:   func() {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder, contextCalled bool) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
				assert.False(t, contextCalled)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behavior()

			mockAuthRepo := new(MockAuthRepo)
			mockClientRepo := new(MockClientRepo)
			jwtManager := auth.NewJWTManager("test-secret", time.Hour)

			authService := services.NewAuthService(mockAuthRepo, mockClientRepo, jwtManager)
			authUseCase := usecases.NewAuthUseCase(authService)

			// Generate a real token for valid test cases
			if tt.authHeader == "Bearer valid.token.here" {
				user := &entities.AuthUser{
					ID:       1,
					ClientID: 100,
					Email:    "test@example.com",
					Role:     entities.RoleAdmin,
				}
				token, _, _ := jwtManager.GenerateToken(user)
				tt.authHeader = "Bearer " + token
			}

			middleware := NewAuthMiddleware(authUseCase)

			router := gin.New()
			contextCalled := false
			router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
				contextCalled = true
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp, contextCalled)
		})
	}
}

func TestAuthMiddleware_RequireRole(t *testing.T) {
	tests := []struct {
		name          string
		userClaims    *entities.TokenClaims
		requiredRoles []entities.UserRole
		behavior      func()
		asserts       func(t *testing.T, resp *httptest.ResponseRecorder, handlerCalled bool)
	}{
		{
			name: "success - user has required role (admin)",
			userClaims: &entities.TokenClaims{
				UserID:   1,
				ClientID: 100,
				Email:    "admin@example.com",
				Role:     entities.RoleAdmin,
			},
			requiredRoles: []entities.UserRole{entities.RoleAdmin},
			behavior:      func() {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder, handlerCalled bool) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.True(t, handlerCalled)
			},
		},
		{
			name: "success - user has one of multiple required roles",
			userClaims: &entities.TokenClaims{
				UserID:   2,
				ClientID: 100,
				Email:    "opto@example.com",
				Role:     entities.RoleOptometrista,
			},
			requiredRoles: []entities.UserRole{entities.RoleAdmin, entities.RoleOptometrista},
			behavior:      func() {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder, handlerCalled bool) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.True(t, handlerCalled)
			},
		},
		{
			name: "error - user does not have required role",
			userClaims: &entities.TokenClaims{
				UserID:   3,
				ClientID: 100,
				Email:    "recep@example.com",
				Role:     entities.RoleRecepcionista,
			},
			requiredRoles: []entities.UserRole{entities.RoleAdmin},
			behavior:      func() {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder, handlerCalled bool) {
				assert.Equal(t, http.StatusForbidden, resp.Code)
				assert.False(t, handlerCalled)
			},
		},
		{
			name:          "error - user not authenticated (no claims in context)",
			userClaims:    nil,
			requiredRoles: []entities.UserRole{entities.RoleAdmin},
			behavior:      func() {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder, handlerCalled bool) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
				assert.False(t, handlerCalled)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behavior()

			mockAuthRepo := new(MockAuthRepo)
			mockClientRepo := new(MockClientRepo)
			jwtManager := auth.NewJWTManager("test-secret", time.Hour)

			authService := services.NewAuthService(mockAuthRepo, mockClientRepo, jwtManager)
			authUseCase := usecases.NewAuthUseCase(authService)
			middleware := NewAuthMiddleware(authUseCase)

			router := gin.New()
			handlerCalled := false

			router.GET("/test", func(c *gin.Context) {
				// Simulate RequireAuth setting user in context
				if tt.userClaims != nil {
					c.Set(UserContextKey, tt.userClaims)
					c.Set(ClientContextKey, tt.userClaims.ClientID)
				}
				c.Next()
			}, middleware.RequireRole(tt.requiredRoles...), func(c *gin.Context) {
				handlerCalled = true
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp, handlerCalled)
		})
	}
}

func TestGetUserClaims(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(c *gin.Context)
		behavior func()
		asserts  func(t *testing.T, claims *entities.TokenClaims, exists bool)
	}{
		{
			name: "success - returns user claims",
			setup: func(c *gin.Context) {
				c.Set(UserContextKey, &entities.TokenClaims{
					UserID:   1,
					ClientID: 100,
					Email:    "test@example.com",
					Role:     entities.RoleAdmin,
				})
			},
			behavior: func() {},
			asserts: func(t *testing.T, claims *entities.TokenClaims, exists bool) {
				assert.True(t, exists)
				assert.NotNil(t, claims)
				assert.Equal(t, int64(1), claims.UserID)
				assert.Equal(t, int64(100), claims.ClientID)
			},
		},
		{
			name:     "failure - user claims not in context",
			setup:    func(c *gin.Context) {},
			behavior: func() {},
			asserts: func(t *testing.T, claims *entities.TokenClaims, exists bool) {
				assert.False(t, exists)
				assert.Nil(t, claims)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behavior()

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			tt.setup(c)

			claims, exists := GetUserClaims(c)

			tt.asserts(t, claims, exists)
		})
	}
}

func TestGetClientID(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(c *gin.Context)
		behavior func()
		asserts  func(t *testing.T, clientID int64, exists bool)
	}{
		{
			name: "success - returns client id",
			setup: func(c *gin.Context) {
				c.Set(ClientContextKey, int64(100))
			},
			behavior: func() {},
			asserts: func(t *testing.T, clientID int64, exists bool) {
				assert.True(t, exists)
				assert.Equal(t, int64(100), clientID)
			},
		},
		{
			name:     "failure - client id not in context",
			setup:    func(c *gin.Context) {},
			behavior: func() {},
			asserts: func(t *testing.T, clientID int64, exists bool) {
				assert.False(t, exists)
				assert.Equal(t, int64(0), clientID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behavior()

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			tt.setup(c)

			clientID, exists := GetClientID(c)

			tt.asserts(t, clientID, exists)
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
