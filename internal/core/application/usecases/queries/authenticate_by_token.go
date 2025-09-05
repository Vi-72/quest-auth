package queries

import (
	"context"
	"time"

	"quest-auth/internal/core/domain/model/kernel"
	"quest-auth/internal/core/ports"
	"quest-auth/internal/pkg/errs"

	"github.com/google/uuid"
)

type AuthenticatedInfo struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Phone     string
	CreatedAt time.Time
}

type AuthenticateByTokenQuery struct {
	RawToken string
}

type AuthenticateByTokenHandler struct {
	jwt ports.JWTService
}

func NewAuthenticateByTokenHandler(jwt ports.JWTService) *AuthenticateByTokenHandler {
	return &AuthenticateByTokenHandler{jwt: jwt}
}

func (h *AuthenticateByTokenHandler) Handle(ctx context.Context, q AuthenticateByTokenQuery) (AuthenticatedInfo, error) {
	// Создаём доменную модель токена (валидирует пустые значения и пробелы)
	token, err := kernel.NewJwtToken(q.RawToken)
	if err != nil {
		return AuthenticatedInfo{}, errs.NewDomainValidationError("jwt_token", "value is required")
	}

	claims, err := h.jwt.ValidateAccessToken(token.String())
	if err != nil {
		return AuthenticatedInfo{}, err
	}

	// Собираем ответ из доступных клеймов (без похода в БД)
	return AuthenticatedInfo{
		ID:        claims.UserID,
		Name:      claims.Name,
		Email:     claims.Email,
		Phone:     claims.Phone,
		CreatedAt: claims.CreatedAt,
	}, nil
}
