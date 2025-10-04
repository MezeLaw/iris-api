package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"iris-api/internal/domain/entities"
)

type JWTManager interface {
	GenerateToken(user *entities.AuthUser) (string, time.Time, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type jwtManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type Claims struct {
	UserID   int64             `json:"user_id"`
	ClientID int64             `json:"client_id"`
	Email    string            `json:"email"`
	Role     entities.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) JWTManager {
	return &jwtManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (m *jwtManager) GenerateToken(user *entities.AuthUser) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.tokenDuration)

	claims := &Claims{
		UserID:   user.ID,
		ClientID: user.ClientID,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (m *jwtManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
