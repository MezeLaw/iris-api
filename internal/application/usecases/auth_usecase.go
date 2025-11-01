package usecases

import (
	"context"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
)

type AuthUseCase interface {
	Register(ctx context.Context, req *entities.RegisterRequest) (*entities.AuthUser, error)
	Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error)
	ValidateToken(tokenString string) (*entities.TokenClaims, error)
	GetUserByID(ctx context.Context, userID int64) (*entities.AuthUser, error)
}

type authUseCase struct {
	authService services.AuthService
}

func NewAuthUseCase(authService services.AuthService) AuthUseCase {
	return &authUseCase{
		authService: authService,
	}
}

func (uc *authUseCase) Register(ctx context.Context, req *entities.RegisterRequest) (*entities.AuthUser, error) {
	return uc.authService.Register(ctx, req)
}

func (uc *authUseCase) Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error) {
	return uc.authService.Login(ctx, req)
}

func (uc *authUseCase) ValidateToken(tokenString string) (*entities.TokenClaims, error) {
	return uc.authService.ValidateToken(tokenString)
}

func (uc *authUseCase) GetUserByID(ctx context.Context, userID int64) (*entities.AuthUser, error) {
	return uc.authService.GetUserByID(ctx, userID)
}
