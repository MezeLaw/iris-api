package services

import (
	"context"
	"errors"
	"time"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"
	"iris-api/pkg/auth"
)

type AuthService struct {
	authRepo   repositories.AuthRepository
	clientRepo repositories.ClientRepository
	jwtManager *auth.JWTManager
}

func NewAuthService(authRepo repositories.AuthRepository, clientRepo repositories.ClientRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		clientRepo: clientRepo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req *entities.RegisterRequest) (*entities.AuthUser, error) {
	// Create new client first
	client := &entities.Client{
		Name:      req.ClientName,
		Slug:      req.ClientSlug,
		Email:     req.ClientEmail,
		Phone:     req.ClientPhone,
		Address:   "",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdClient, err := s.clientRepo.Create(ctx, client)
	if err != nil {
		return nil, errors.New("failed to create client: " + err.Error())
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user for the new client
	user := &entities.AuthUser{
		ClientID:  createdClient.ID,
		Email:     req.Email,
		Password:  hashedPassword,
		Name:      req.Name,
		Role:      req.Role,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.authRepo.Register(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error) {
	// Find user by email (globally unique)
	user, err := s.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, expiresAt, err := s.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Update last login
	_ = s.authRepo.UpdateLastLogin(ctx, user.ID)

	// Clear password before returning
	user.Password = ""

	return &entities.LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*entities.TokenClaims, error) {
	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &entities.TokenClaims{
		UserID:   claims.UserID,
		ClientID: claims.ClientID,
		Email:    claims.Email,
		Role:     claims.Role,
	}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID int64) (*entities.AuthUser, error) {
	user, err := s.authRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Clear password
	user.Password = ""
	return user, nil
}
