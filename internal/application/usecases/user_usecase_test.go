package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/domain/entities"
)

func TestUserUseCase_CreateUser(t *testing.T) {
	tests := []struct {
		name     string
		req      *entities.CreateUserRequest
		behavior func(m *MockUserService)
		asserts  func(t *testing.T, user *entities.User, err error)
	}{
		{
			name: "success - creates user",
			req: &entities.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			behavior: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.MatchedBy(func(req *entities.CreateUserRequest) bool {
					return req.Email == "test@example.com"
				})).Return(&entities.User{
					ID:    1,
					Email: "test@example.com",
					Name:  "Test User",
				}, nil)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, 1, user.ID)
			},
		},
		{
			name: "error - service fails",
			req: &entities.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			behavior: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("service error"))
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.behavior(mockService)

			useCase := NewUserUseCase(mockService)
			user, err := useCase.CreateUser(context.Background(), tt.req)

			tt.asserts(t, user, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		behavior func(m *MockUserService)
		asserts  func(t *testing.T, user *entities.User, err error)
	}{
		{
			name: "success - gets user by id",
			id:   1,
			behavior: func(m *MockUserService) {
				m.On("GetUserByID", mock.Anything, 1).Return(&entities.User{
					ID:    1,
					Email: "test@example.com",
					Name:  "Test User",
				}, nil)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, 1, user.ID)
			},
		},
		{
			name: "error - user not found",
			id:   999,
			behavior: func(m *MockUserService) {
				m.On("GetUserByID", mock.Anything, 999).Return(nil, errors.New("not found"))
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.behavior(mockService)

			useCase := NewUserUseCase(mockService)
			user, err := useCase.GetUserByID(context.Background(), tt.id)

			tt.asserts(t, user, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_GetUsers(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		offset   int
		behavior func(m *MockUserService)
		asserts  func(t *testing.T, resp *GetUsersResponse, err error)
	}{
		{
			name:   "success - returns users with has_more true",
			limit:  10,
			offset: 0,
			behavior: func(m *MockUserService) {
				m.On("GetUsers", mock.Anything, 10, 0).Return([]*entities.User{
					{ID: 1, Email: "user1@example.com", Name: "User 1"},
					{ID: 2, Email: "user2@example.com", Name: "User 2"},
				}, 20, nil)
			},
			asserts: func(t *testing.T, resp *GetUsersResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Users, 2)
				assert.Equal(t, 20, resp.Total)
				assert.Equal(t, 10, resp.Limit)
				assert.Equal(t, 0, resp.Offset)
				assert.True(t, resp.HasMore)
			},
		},
		{
			name:   "success - returns users with has_more false",
			limit:  10,
			offset: 15,
			behavior: func(m *MockUserService) {
				m.On("GetUsers", mock.Anything, 10, 15).Return([]*entities.User{
					{ID: 16, Email: "user16@example.com", Name: "User 16"},
				}, 20, nil)
			},
			asserts: func(t *testing.T, resp *GetUsersResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, 20, resp.Total)
				assert.False(t, resp.HasMore)
			},
		},
		{
			name:   "error - service fails",
			limit:  10,
			offset: 0,
			behavior: func(m *MockUserService) {
				m.On("GetUsers", mock.Anything, 10, 0).Return(nil, 0, errors.New("database error"))
			},
			asserts: func(t *testing.T, resp *GetUsersResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.behavior(mockService)

			useCase := NewUserUseCase(mockService)
			resp, err := useCase.GetUsers(context.Background(), tt.limit, tt.offset)

			tt.asserts(t, resp, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		req      *entities.UpdateUserRequest
		behavior func(m *MockUserService)
		asserts  func(t *testing.T, user *entities.User, err error)
	}{
		{
			name: "success - updates user",
			id:   1,
			req: &entities.UpdateUserRequest{
				Name: "Updated Name",
			},
			behavior: func(m *MockUserService) {
				m.On("UpdateUser", mock.Anything, 1, mock.Anything).Return(&entities.User{
					ID:    1,
					Email: "test@example.com",
					Name:  "Updated Name",
				}, nil)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "Updated Name", user.Name)
			},
		},
		{
			name: "error - service fails",
			id:   1,
			req:  &entities.UpdateUserRequest{Name: "Test"},
			behavior: func(m *MockUserService) {
				m.On("UpdateUser", mock.Anything, 1, mock.Anything).Return(nil, errors.New("update failed"))
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.behavior(mockService)

			useCase := NewUserUseCase(mockService)
			user, err := useCase.UpdateUser(context.Background(), tt.id, tt.req)

			tt.asserts(t, user, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		behavior func(m *MockUserService)
		asserts  func(t *testing.T, err error)
	}{
		{
			name: "success - deletes user",
			id:   1,
			behavior: func(m *MockUserService) {
				m.On("DeleteUser", mock.Anything, 1).Return(nil)
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error - service fails",
			id:   999,
			behavior: func(m *MockUserService) {
				m.On("DeleteUser", mock.Anything, 999).Return(errors.New("delete failed"))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.behavior(mockService)

			useCase := NewUserUseCase(mockService)
			err := useCase.DeleteUser(context.Background(), tt.id)

			tt.asserts(t, err)
			mockService.AssertExpectations(t)
		})
	}
}

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req *entities.CreateUserRequest) (*entities.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserService) GetUsers(ctx context.Context, limit, offset int) ([]*entities.User, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.User), args.Int(1), args.Error(2)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id int, req *entities.UpdateUserRequest) (*entities.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
