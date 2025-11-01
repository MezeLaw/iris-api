package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"iris-api/internal/domain/entities"
)

const (
	AuthorizationHeader = "Authorization"
	UserContextKey      = "user"
	ClientContextKey    = "client_id"
)

type AuthUseCase interface {
	ValidateToken(tokenString string) (*entities.TokenClaims, error)
}

type AuthMiddleware struct {
	authUseCase AuthUseCase
}

func NewAuthMiddleware(authUseCase AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.authUseCase.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(UserContextKey, claims)
		c.Set(ClientContextKey, claims.ClientID)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(roles ...entities.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims, exists := c.Get(UserContextKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		claims := userClaims.(*entities.TokenClaims)

		// Check if user has required role
		hasRole := false
		for _, role := range roles {
			if claims.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Helper function to get user claims from context
func GetUserClaims(c *gin.Context) (*entities.TokenClaims, bool) {
	userClaims, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}
	claims, ok := userClaims.(*entities.TokenClaims)
	return claims, ok
}

// Helper function to get client ID from context
func GetClientID(c *gin.Context) (int64, bool) {
	clientID, exists := c.Get(ClientContextKey)
	if !exists {
		return 0, false
	}
	id, ok := clientID.(int64)
	return id, ok
}
