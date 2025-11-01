package services

import (
	"context"
	"fmt"
	"strings"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, req *entities.CreateUserRequest) (*entities.User, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*entities.User, int, error)
	UpdateUser(ctx context.Context, id int, req *entities.UpdateUserRequest) (*entities.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *entities.CreateUserRequest) (*entities.User, error) {
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &entities.User{
		Email: strings.ToLower(strings.TrimSpace(req.Email)),
		Name:  strings.TrimSpace(req.Name),
	}

	return s.userRepo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetUsers(ctx context.Context, limit, offset int) ([]*entities.User, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int, req *entities.UpdateUserRequest) (*entities.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if err := s.validateUpdateUserRequest(req); err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Email != "" && req.Email != existingUser.Email {
		emailUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
		if emailUser != nil && emailUser.ID != id {
			return nil, fmt.Errorf("user with email %s already exists", req.Email)
		}
		existingUser.Email = strings.ToLower(strings.TrimSpace(req.Email))
	}

	if req.Name != "" {
		existingUser.Name = strings.TrimSpace(req.Name)
	}

	return s.userRepo.Update(ctx, id, existingUser)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *userService) validateCreateUserRequest(req *entities.CreateUserRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}
	if len(req.Name) > 100 {
		return fmt.Errorf("name must be less than 100 characters long")
	}
	return nil
}

func (s *userService) validateUpdateUserRequest(req *entities.UpdateUserRequest) error {
	if req.Name != "" && len(req.Name) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}
	if req.Name != "" && len(req.Name) > 100 {
		return fmt.Errorf("name must be less than 100 characters long")
	}
	return nil
}
