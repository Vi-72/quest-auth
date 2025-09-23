package jwt

import (
	"fmt"
	"time"

	"quest-auth/internal/core/ports"
	"quest-auth/internal/pkg/errs"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewService(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration) *Service {
	return &Service{
		secretKey:            []byte(secretKey),
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// Claims структура для JWT токена
type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt int64     `json:"created_at,omitempty"`
	Type      string    `json:"type"` // "access" или "refresh"
	jwt.RegisteredClaims
}

// GenerateTokenPair создает пару access и refresh токенов
func (s *Service) GenerateTokenPair(
	userID uuid.UUID,
	email, name, phone string,
	createdAt time.Time,
) (*ports.TokenPair, error) {
	now := time.Now()

	// Access token
	accessClaims := &Claims{
		UserID:    userID,
		Email:     email,
		Name:      name,
		Phone:     phone,
		CreatedAt: createdAt.Unix(),
		Type:      "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return nil, errs.WrapInfrastructureError("generating access token", err)
	}

	// Refresh token (PII не обязательно)
	refreshClaims := &Claims{
		UserID:    userID,
		Email:     email,
		Name:      name,
		Phone:     phone,
		CreatedAt: createdAt.Unix(),
		Type:      "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID.String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.secretKey)
	if err != nil {
		return nil, errs.WrapInfrastructureError("generating refresh token", err)
	}

	return &ports.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenDuration.Seconds()),
	}, nil
}

// parseToken разбирает и валидирует JWT токен
func (s *Service) parseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, errs.NewJWTValidationErrorWithCause("parsing token", err)
	}

	if !token.Valid {
		return nil, errs.NewJWTValidationError("token is invalid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errs.NewJWTValidationError("invalid token claims")
	}

	return claims, nil
}

// ValidateAccessToken проверяет валидность access токена
func (s *Service) ValidateAccessToken(tokenString string) (*ports.TokenClaims, error) {
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "access" {
		return nil, errs.NewJWTValidationError("token is not an access token")
	}

	return &ports.TokenClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		Name:      claims.Name,
		Phone:     claims.Phone,
		CreatedAt: time.Unix(claims.CreatedAt, 0),
		Exp:       claims.ExpiresAt.Time,
	}, nil
}

// RefreshTokens обновляет токены по refresh токену
func (s *Service) RefreshTokens(refreshTokenString string) (*ports.TokenPair, error) {
	claims, err := s.parseToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, errs.NewJWTValidationError("token is not a refresh token")
	}

	// Генерируем новую пару токенов
	return s.GenerateTokenPair(claims.UserID, claims.Email, claims.Name, claims.Phone, time.Unix(claims.CreatedAt, 0))
}

// Compile-time check that Service implements JWTService
var _ ports.JWTService = (*Service)(nil)
