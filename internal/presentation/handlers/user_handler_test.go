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

	"iris-api/internal/application/usecases"
	"iris-api/internal/domain/entities"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name     string
		body     interface{}
		behavior func(m *MockUserUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "success - creates user",
			body: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
			behavior: func(m *MockUserUseCase) {
				m.On("CreateUser", mock.Anything, mock.Anything).Return(&entities.User{
					ID:    1,
					Email: "test@example.com",
					Name:  "Test User",
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "User created successfully", response["message"])
			},
		},
		{
			name:     "error - invalid request body",
			body:     "invalid json",
			behavior: func(m *MockUserUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid request body")
			},
		},
		{
			name: "error - usecase fails",
			body: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
			behavior: func(m *MockUserUseCase) {
				m.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("usecase error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockUserUseCase)
			tt.behavior(mockUseCase)

			handler := NewUserHandler(mockUseCase)

			router := gin.New()
			router.POST("/users", handler.CreateUser)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		behavior func(m *MockUserUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:   "success - gets user by id",
			userID: "1",
			behavior: func(m *MockUserUseCase) {
				m.On("GetUserByID", mock.Anything, 1).Return(&entities.User{
					ID:    1,
					Email: "test@example.com",
					Name:  "Test User",
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, "User retrieved successfully", response["message"])
			},
		},
		{
			name:     "error - invalid user id",
			userID:   "invalid",
			behavior: func(m *MockUserUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid user ID")
			},
		},
		{
			name:   "error - user not found",
			userID: "999",
			behavior: func(m *MockUserUseCase) {
				m.On("GetUserByID", mock.Anything, 999).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockUserUseCase)
			tt.behavior(mockUseCase)

			handler := NewUserHandler(mockUseCase)

			router := gin.New()
			router.GET("/users/:id", handler.GetUserByID)

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		behavior func(m *MockUserUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:  "success - gets users with default pagination",
			query: "",
			behavior: func(m *MockUserUseCase) {
				m.On("GetUsers", mock.Anything, 10, 0).Return(&usecases.GetUsersResponse{
					Users:   []*entities.User{{ID: 1, Email: "test@example.com", Name: "Test"}},
					Total:   1,
					Limit:   10,
					Offset:  0,
					HasMore: false,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - gets users with custom pagination",
			query: "?limit=20&offset=10",
			behavior: func(m *MockUserUseCase) {
				m.On("GetUsers", mock.Anything, 20, 10).Return(&usecases.GetUsersResponse{
					Users:   []*entities.User{},
					Total:   0,
					Limit:   20,
					Offset:  10,
					HasMore: false,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - limit zero uses default",
			query: "?limit=0",
			behavior: func(m *MockUserUseCase) {
				m.On("GetUsers", mock.Anything, 10, 0).Return(&usecases.GetUsersResponse{
					Users:   []*entities.User{},
					Total:   0,
					Limit:   10,
					Offset:  0,
					HasMore: false,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - negative limit uses default",
			query: "?limit=-5",
			behavior: func(m *MockUserUseCase) {
				m.On("GetUsers", mock.Anything, 10, 0).Return(&usecases.GetUsersResponse{
					Users:   []*entities.User{},
					Total:   0,
					Limit:   10,
					Offset:  0,
					HasMore: false,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "success - negative offset uses default",
			query: "?offset=-10",
			behavior: func(m *MockUserUseCase) {
				m.On("GetUsers", mock.Anything, 10, 0).Return(&usecases.GetUsersResponse{
					Users:   []*entities.User{},
					Total:   0,
					Limit:   10,
					Offset:  0,
					HasMore: false,
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:  "error - usecase fails",
			query: "",
			behavior: func(m *MockUserUseCase) {
				m.On("GetUsers", mock.Anything, 10, 0).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockUserUseCase)
			tt.behavior(mockUseCase)

			handler := NewUserHandler(mockUseCase)

			router := gin.New()
			router.GET("/users", handler.GetUsers)

			req := httptest.NewRequest(http.MethodGet, "/users"+tt.query, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		body     interface{}
		behavior func(m *MockUserUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:   "success - updates user",
			userID: "1",
			body: map[string]interface{}{
				"name": "Updated Name",
			},
			behavior: func(m *MockUserUseCase) {
				m.On("UpdateUser", mock.Anything, 1, mock.Anything).Return(&entities.User{
					ID:    1,
					Email: "test@example.com",
					Name:  "Updated Name",
				}, nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:     "error - invalid user id",
			userID:   "invalid",
			body:     map[string]interface{}{},
			behavior: func(m *MockUserUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:     "error - invalid body",
			userID:   "1",
			body:     "invalid",
			behavior: func(m *MockUserUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:   "error - user not found",
			userID: "999",
			body: map[string]interface{}{
				"name": "Updated Name",
			},
			behavior: func(m *MockUserUseCase) {
				m.On("UpdateUser", mock.Anything, 999, mock.Anything).Return(nil, errors.New("user not found"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockUserUseCase)
			tt.behavior(mockUseCase)

			handler := NewUserHandler(mockUseCase)

			router := gin.New()
			router.PUT("/users/:id", handler.UpdateUser)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.userID, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		behavior func(m *MockUserUseCase)
		asserts  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:   "success - deletes user",
			userID: "1",
			behavior: func(m *MockUserUseCase) {
				m.On("DeleteUser", mock.Anything, 1).Return(nil)
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, resp.Code)
			},
		},
		{
			name:     "error - invalid user id",
			userID:   "invalid",
			behavior: func(m *MockUserUseCase) {},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:   "error - user not found",
			userID: "999",
			behavior: func(m *MockUserUseCase) {
				m.On("DeleteUser", mock.Anything, 999).Return(errors.New("not found"))
			},
			asserts: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockUserUseCase)
			tt.behavior(mockUseCase)

			handler := NewUserHandler(mockUseCase)

			router := gin.New()
			router.DELETE("/users/:id", handler.DeleteUser)

			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tt.asserts(t, resp)
			mockUseCase.AssertExpectations(t)
		})
	}
}

// MockUserUseCase is a mock implementation of UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) CreateUser(ctx context.Context, req *entities.CreateUserRequest) (*entities.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) GetUsers(ctx context.Context, limit, offset int) (*usecases.GetUsersResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.GetUsersResponse), args.Error(1)
}

func (m *MockUserUseCase) UpdateUser(ctx context.Context, id int, req *entities.UpdateUserRequest) (*entities.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUseCase) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
