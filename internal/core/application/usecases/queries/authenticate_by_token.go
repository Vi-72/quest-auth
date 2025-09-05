package queries

import (
	"context"

	"quest-auth/internal/core/domain/model/kernel"
	"quest-auth/internal/core/ports"

	"github.com/google/uuid"
)

type AuthenticatedInfo struct {
	ID       uuid.UUID
	FullName string
	Email    string
	Phone    string
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
		return AuthenticatedInfo{}, err
	}

	claims, err := h.jwt.ValidateAccessToken(token.String())
	if err != nil {
		return AuthenticatedInfo{}, err
	}

	// Собираем ответ из доступных клеймов (без похода в БД)
	return AuthenticatedInfo{
		ID:    claims.UserID,
		Email: claims.Email,
		// FullName и Phone останутся пустыми, если не добавим их в клеймы позже
	}, nil
}
