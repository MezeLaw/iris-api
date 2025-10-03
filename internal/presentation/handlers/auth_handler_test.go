package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/domain/entities"
	"iris-api/internal/presentation/middleware"
)

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name     string
		body     interface{}
		behavior func(m *MockAuthUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "success - registers new user",
			body: map[string]interface{}{
				"client_name":  "Test Clinic",
				"client_slug":  "testclinic",
				"client_email": "clinic@example.com",
				"client_phone": "1234567890",
				"email":        "user@example.com",
				"password":     "securepassword",
				"name":         "Test User",
				"role":         "admin",
			},
			behavior: func(m *MockAuthUseCase) {
				m.On("Register", mock.Anything, mock.Anything).Return(&entities.AuthUser{
					ID:       1,
					ClientID: 1,
					Email:    "user@example.com",
					Name:     "Test User",
					Role:     entities.RoleAdmin,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "User registered successfully", response["message"])
			},
		},
		{
			name: "error - invalid request body",
			body: "invalid json",
			behavior: func(m *MockAuthUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid request body")
			},
		},
		{
			name: "error - registration fails",
			body: map[string]interface{}{
				"client_name": "Test Clinic",
				"client_slug": "testclinic",
				"email":       "user@example.com",
				"password":    "securepassword",
				"name":        "Test User",
				"role":        "admin",
			},
			behavior: func(m *MockAuthUseCase) {
				m.On("Register", mock.Anything, mock.Anything).Return(nil, errors.New("email already exists"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Failed to register user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockAuthUseCase)
			tt.behavior(mockUseCase)

			handler := NewAuthHandler(mockUseCase)

			router := gin.New()
			router.POST("/register", handler.Register)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name     string
		body     interface{}
		behavior func(m *MockAuthUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "success - logs in user",
			body: map[string]interface{}{
				"email":    "user@example.com",
				"password": "correctpassword",
			},
			behavior: func(m *MockAuthUseCase) {
				m.On("Login", mock.Anything, mock.Anything).Return(&entities.LoginResponse{
					Token: "valid.jwt.token",
					User: &entities.AuthUser{
						ID:       1,
						ClientID: 100,
						Email:    "user@example.com",
						Name:     "Test User",
						Role:     entities.RoleAdmin,
					},
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Login successful", response["message"])
				data := response["data"].(map[string]interface{})
				assert.NotEmpty(t, data["token"])
			},
		},
		{
			name: "error - invalid request body",
			body: "invalid json",
			behavior: func(m *MockAuthUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid request body")
			},
		},
		{
			name: "error - invalid credentials",
			body: map[string]interface{}{
				"email":    "user@example.com",
				"password": "wrongpassword",
			},
			behavior: func(m *MockAuthUseCase) {
				m.On("Login", mock.Anything, mock.Anything).Return(nil, errors.New("invalid credentials"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Login failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockAuthUseCase)
			tt.behavior(mockUseCase)

			handler := NewAuthHandler(mockUseCase)

			router := gin.New()
			router.POST("/login", handler.Login)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_GetProfile(t *testing.T) {
	tests := []struct {
		name     string
		setClaim bool
		userID   int64
		behavior func(m *MockAuthUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:     "success - gets user profile",
			setClaim: true,
			userID:   1,
			behavior: func(m *MockAuthUseCase) {
				m.On("GetUserByID", mock.Anything, int64(1)).Return(&entities.AuthUser{
					ID:       1,
					ClientID: 100,
					Email:    "user@example.com",
					Name:     "Test User",
					Role:     entities.RoleAdmin,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "Profile retrieved successfully", response["message"])
			},
		},
		{
			name:     "error - user not authenticated",
			setClaim: false,
			behavior: func(m *MockAuthUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "not authenticated")
			},
		},
		{
			name:     "error - user not found",
			setClaim: true,
			userID:   999,
			behavior: func(m *MockAuthUseCase) {
				m.On("GetUserByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockAuthUseCase)
			tt.behavior(mockUseCase)

			handler := NewAuthHandler(mockUseCase)

			router := gin.New()
			router.GET("/profile", func(c *gin.Context) {
				if tt.setClaim {
					c.Set(middleware.UserContextKey, &entities.TokenClaims{
						UserID:   tt.userID,
						ClientID: 100,
						Email:    "user@example.com",
						Role:     entities.RoleAdmin,
					})
				}
				c.Next()
			}, handler.GetProfile)

			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

// MockAuthUseCase is a mock implementation of AuthUseCase
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Register(ctx context.Context, req *entities.RegisterRequest) (*entities.AuthUser, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuthUser), args.Error(1)
}

func (m *MockAuthUseCase) Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.LoginResponse), args.Error(1)
}

func (m *MockAuthUseCase) ValidateToken(token string) (*entities.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TokenClaims), args.Error(1)
}

func (m *MockAuthUseCase) GetUserByID(ctx context.Context, userID int64) (*entities.AuthUser, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuthUser), args.Error(1)
}
