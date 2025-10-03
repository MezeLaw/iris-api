package usecases

import (
	"context"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
)

type AuthUseCase struct {
	authService *services.AuthService
}

func NewAuthUseCase(authService *services.AuthService) *AuthUseCase {
	return &AuthUseCase{
		authService: authService,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, req *entities.RegisterRequest) (*entities.AuthUser, error) {
	return uc.authService.Register(ctx, req)
}

func (uc *AuthUseCase) Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error) {
	return uc.authService.Login(ctx, req)
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (*entities.TokenClaims, error) {
	return uc.authService.ValidateToken(tokenString)
}

func (uc *AuthUseCase) GetUserByID(ctx context.Context, userID int64) (*entities.AuthUser, error) {
	return uc.authService.GetUserByID(ctx, userID)
}
