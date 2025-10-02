package usecases

import (
	"context"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, req *entities.CreateUserRequest) (*entities.User, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
	GetUsers(ctx context.Context, limit, offset int) (*GetUsersResponse, error)
	UpdateUser(ctx context.Context, id int, req *entities.UpdateUserRequest) (*entities.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type GetUsersResponse struct {
	Users      []*entities.User `json:"users"`
	Total      int              `json:"total"`
	Limit      int              `json:"limit"`
	Offset     int              `json:"offset"`
	HasMore    bool             `json:"has_more"`
}

type userUseCase struct {
	userService services.UserService
}

func NewUserUseCase(userService services.UserService) UserUseCase {
	return &userUseCase{
		userService: userService,
	}
}

func (uc *userUseCase) CreateUser(ctx context.Context, req *entities.CreateUserRequest) (*entities.User, error) {
	return uc.userService.CreateUser(ctx, req)
}

func (uc *userUseCase) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	return uc.userService.GetUserByID(ctx, id)
}

func (uc *userUseCase) GetUsers(ctx context.Context, limit, offset int) (*GetUsersResponse, error) {
	users, total, err := uc.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	hasMore := (offset + limit) < total

	return &GetUsersResponse{
		Users:   users,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: hasMore,
	}, nil
}

func (uc *userUseCase) UpdateUser(ctx context.Context, id int, req *entities.UpdateUserRequest) (*entities.User, error) {
	return uc.userService.UpdateUser(ctx, id, req)
}

func (uc *userUseCase) DeleteUser(ctx context.Context, id int) error {
	return uc.userService.DeleteUser(ctx, id)
}