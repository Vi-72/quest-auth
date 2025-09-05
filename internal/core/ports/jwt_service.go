package ports

import (
	"time"

	"github.com/google/uuid"
)

// TokenPair представляет пару токенов
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64 // в секундах
}

// JWTService интерфейс для работы с JWT токенами
type JWTService interface {
	// GenerateTokenPair создает пару access и refresh токенов
	GenerateTokenPair(userID uuid.UUID, email, name, phone string, createdAt time.Time) (*TokenPair, error)

	// ValidateAccessToken проверяет валидность access токена
	ValidateAccessToken(token string) (*TokenClaims, error)

	// RefreshTokens обновляет токены по refresh токену
	RefreshTokens(refreshToken string) (*TokenPair, error)
}

// TokenClaims содержит данные из токена
type TokenClaims struct {
	UserID    uuid.UUID
	Email     string
	Name      string
	Phone     string
	Exp       time.Time
	CreatedAt time.Time
}
