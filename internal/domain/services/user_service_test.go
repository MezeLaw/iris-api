package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iris-api/internal/domain/entities"
)

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name     string
		req      *entities.CreateUserRequest
		behavior func(m *MockUserRepository)
		asserts  func(t *testing.T, user *entities.User, err error)
	}{
		{
			name: "success - creates user with valid data",
			req: &entities.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			behavior: func(m *MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("not found"))
				m.On("Create", mock.Anything, mock.MatchedBy(func(u *entities.User) bool {
					return u.Email == "test@example.com" && u.Name == "Test User"
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
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, "Test User", user.Name)
			},
		},
		{
			name: "success - trims and lowercases email",
			req: &entities.CreateUserRequest{
				Email: "  TEST@EXAMPLE.COM  ",
				Name:  "  Test User  ",
			},
			behavior: func(m *MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "  TEST@EXAMPLE.COM  ").Return(nil, errors.New("not found"))
				m.On("Create", mock.Anything, mock.MatchedBy(func(u *entities.User) bool {
					return u.Email == "test@example.com" && u.Name == "Test User"
				})).Return(&entities.User{
					ID:    1,
					Email: "test@example.com",
					Name:  "Test User",
				}, nil)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, "Test User", user.Name)
			},
		},
		{
			name: "error - email is required",
			req: &entities.CreateUserRequest{
				Email: "",
				Name:  "Test User",
			},
			behavior: func(m *MockUserRepository) {},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "email is required")
			},
		},
		{
			name: "error - name is required",
			req: &entities.CreateUserRequest{
				Email: "test@example.com",
				Name:  "",
			},
			behavior: func(m *MockUserRepository) {},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "name is required")
			},
		},
		{
			name: "error - name too short",
			req: &entities.CreateUserRequest{
				Email: "test@example.com",
				Name:  "A",
			},
			behavior: func(m *MockUserRepository) {},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "at least 2 characters")
			},
		},
		{
			name: "error - email already exists",
			req: &entities.CreateUserRequest{
				Email: "existing@example.com",
				Name:  "Test User",
			},
			behavior: func(m *MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "existing@example.com").Return(&entities.User{
					ID:    1,
					Email: "existing@example.com",
				}, nil)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "already exists")
			},
		},
		{
			name: "error - repository create fails",
			req: &entities.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			behavior: func(m *MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("not found"))
				m.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.behavior(mockRepo)

			service := NewUserService(mockRepo)
			user, err := service.CreateUser(context.Background(), tt.req)

			tt.asserts(t, user, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		behavior func(m *MockUserRepository)
		asserts  func(t *testing.T, user *entities.User, err error)
	}{
		{
			name: "success - returns user by id",
			id:   1,
			behavior: func(m *MockUserRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&entities.User{
					ID:    1,
					Email: "test@example.com",
					Name:  "Test User",
				}, nil)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, 1, user.ID)
				assert.Equal(t, "test@example.com", user.Email)
			},
		},
		{
			name:     "error - invalid id (zero)",
			id:       0,
			behavior: func(m *MockUserRepository) {},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "invalid user ID")
			},
		},
		{
			name:     "error - invalid id (negative)",
			id:       -1,
			behavior: func(m *MockUserRepository) {},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "invalid user ID")
			},
		},
		{
			name: "error - user not found",
			id:   999,
			behavior: func(m *MockUserRepository) {
				m.On("GetByID", mock.Anything, 999).Return(nil, errors.New("user not found"))
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.behavior(mockRepo)

			service := NewUserService(mockRepo)
			user, err := service.GetUserByID(context.Background(), tt.id)

			tt.asserts(t, user, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUsers(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		offset   int
		behavior func(m *MockUserRepository)
		asserts  func(t *testing.T, users []*entities.User, total int, err error)
	}{
		{
			name:   "success - returns users with default limit",
			limit:  10,
			offset: 0,
			behavior: func(m *MockUserRepository) {
				m.On("GetAll", mock.Anything, 10, 0).Return([]*entities.User{
					{ID: 1, Email: "user1@example.com", Name: "User 1"},
					{ID: 2, Email: "user2@example.com", Name: "User 2"},
				}, nil)
				m.On("Count", mock.Anything).Return(2, nil)
			},
			asserts: func(t *testing.T, users []*entities.User, total int, err error) {
				assert.NoError(t, err)
				assert.Len(t, users, 2)
				assert.Equal(t, 2, total)
			},
		},
		{
			name:   "success - applies default limit when zero",
			limit:  0,
			offset: 0,
			behavior: func(m *MockUserRepository) {
				m.On("GetAll", mock.Anything, 10, 0).Return([]*entities.User{}, nil)
				m.On("Count", mock.Anything).Return(0, nil)
			},
			asserts: func(t *testing.T, users []*entities.User, total int, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, users)
			},
		},
		{
			name:   "success - caps limit at 100",
			limit:  200,
			offset: 0,
			behavior: func(m *MockUserRepository) {
				m.On("GetAll", mock.Anything, 100, 0).Return([]*entities.User{}, nil)
				m.On("Count", mock.Anything).Return(0, nil)
			},
			asserts: func(t *testing.T, users []*entities.User, total int, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "success - sets negative offset to zero",
			limit:  10,
			offset: -5,
			behavior: func(m *MockUserRepository) {
				m.On("GetAll", mock.Anything, 10, 0).Return([]*entities.User{}, nil)
				m.On("Count", mock.Anything).Return(0, nil)
			},
			asserts: func(t *testing.T, users []*entities.User, total int, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "error - repository GetAll fails",
			limit:  10,
			offset: 0,
			behavior: func(m *MockUserRepository) {
				m.On("GetAll", mock.Anything, 10, 0).Return(nil, errors.New("database error"))
			},
			asserts: func(t *testing.T, users []*entities.User, total int, err error) {
				assert.Error(t, err)
				assert.Nil(t, users)
				assert.Equal(t, 0, total)
			},
		},
		{
			name:   "error - repository Count fails",
			limit:  10,
			offset: 0,
			behavior: func(m *MockUserRepository) {
				m.On("GetAll", mock.Anything, 10, 0).Return([]*entities.User{}, nil)
				m.On("Count", mock.Anything).Return(0, errors.New("count error"))
			},
			asserts: func(t *testing.T, users []*entities.User, total int, err error) {
				assert.Error(t, err)
				assert.Nil(t, users)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.behavior(mockRepo)

			service := NewUserService(mockRepo)
			users, total, err := service.GetUsers(context.Background(), tt.limit, tt.offset)

			tt.asserts(t, users, total, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		req      *entities.UpdateUserRequest
		behavior func(m *MockUserRepository)
		asserts  func(t *testing.T, user *entities.User, err error)
	}{
		{
			name: "success - updates user name",
			id:   1,
			req: &entities.UpdateUserRequest{
				Name: "Updated Name",
			},
			behavior: func(m *MockUserRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&entities.User{
					ID:    1,
					Email: "test@example.com",
					Name:  "Old Name",
				}, nil)
				m.On("Update", mock.Anything, 1, mock.MatchedBy(func(u *entities.User) bool {
					return u.Name == "Updated Name"
				})).Return(&entities.User{
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
			name: "success - updates user email",
			id:   1,
			req: &entities.UpdateUserRequest{
				Email: "new@example.com",
			},
			behavior: func(m *MockUserRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&entities.User{
					ID:    1,
					Email: "old@example.com",
					Name:  "Test User",
				}, nil)
				m.On("GetByEmail", mock.Anything, "new@example.com").Return(nil, errors.New("not found"))
				m.On("Update", mock.Anything, 1, mock.MatchedBy(func(u *entities.User) bool {
					return u.Email == "new@example.com"
				})).Return(&entities.User{
					ID:    1,
					Email: "new@example.com",
					Name:  "Test User",
				}, nil)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "new@example.com", user.Email)
			},
		},
		{
			name:     "error - invalid user id",
			id:       0,
			req:      &entities.UpdateUserRequest{Name: "Test"},
			behavior: func(m *MockUserRepository) {},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "invalid user ID")
			},
		},
		{
			name: "error - name too short",
			id:   1,
			req: &entities.UpdateUserRequest{
				Name: "A",
			},
			behavior: func(m *MockUserRepository) {},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "at least 2 characters")
			},
		},
		{
			name: "error - user not found",
			id:   999,
			req: &entities.UpdateUserRequest{
				Name: "Test User",
			},
			behavior: func(m *MockUserRepository) {
				m.On("GetByID", mock.Anything, 999).Return(nil, errors.New("user not found"))
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
		{
			name: "error - email already exists for different user",
			id:   1,
			req: &entities.UpdateUserRequest{
				Email: "taken@example.com",
			},
			behavior: func(m *MockUserRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&entities.User{
					ID:    1,
					Email: "old@example.com",
					Name:  "Test User",
				}, nil)
				m.On("GetByEmail", mock.Anything, "taken@example.com").Return(&entities.User{
					ID:    2,
					Email: "taken@example.com",
				}, nil)
			},
			asserts: func(t *testing.T, user *entities.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), "already exists")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.behavior(mockRepo)

			service := NewUserService(mockRepo)
			user, err := service.UpdateUser(context.Background(), tt.id, tt.req)

			tt.asserts(t, user, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		behavior func(m *MockUserRepository)
		asserts  func(t *testing.T, err error)
	}{
		{
			name: "success - deletes user",
			id:   1,
			behavior: func(m *MockUserRepository) {
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			asserts: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:     "error - invalid user id",
			id:       0,
			behavior: func(m *MockUserRepository) {},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid user ID")
			},
		},
		{
			name: "error - repository delete fails",
			id:   1,
			behavior: func(m *MockUserRepository) {
				m.On("Delete", mock.Anything, 1).Return(errors.New("database error"))
			},
			asserts: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.behavior(mockRepo)

			service := NewUserService(mockRepo)
			err := service.DeleteUser(context.Background(), tt.id)

			tt.asserts(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id int, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}
