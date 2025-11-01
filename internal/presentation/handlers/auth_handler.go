package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"iris-api/internal/domain/entities"
	"iris-api/internal/presentation/middleware"
)

// AuthUseCase interface defines the methods needed by AuthHandler
type AuthUseCase interface {
	Register(ctx context.Context, req *entities.RegisterRequest) (*entities.AuthUser, error)
	Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error)
	ValidateToken(token string) (*entities.TokenClaims, error)
	GetUserByID(ctx context.Context, userID int64) (*entities.AuthUser, error)
}

type AuthHandler struct {
	authUseCase AuthUseCase
}

func NewAuthHandler(authUseCase AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req entities.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	user, err := h.authUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to register user",
			"details": err.Error(),
		})
		return
	}

	// Clear password before returning
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req entities.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	loginResponse, err := h.authUseCase.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Login failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    loginResponse,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	claims, exists := middleware.GetUserClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	user, err := h.authUseCase.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    user,
	})
}
