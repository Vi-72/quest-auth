package grpc

import (
	"context"
	authv1 "github.com/Vi-72/quest-auth/api/grpc/sdk/go/proto/auth/v1"
	"strings"

	"quest-auth/internal/core/application/usecases/queries"
	"quest-auth/internal/pkg/errs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthHandler реализует gRPC сервис аутентификации
type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	authenticateByToken *queries.AuthenticateByTokenHandler
}

// NewAuthHandler создает новый gRPC handler для аутентификации
func NewAuthHandler(
	authenticateByToken *queries.AuthenticateByTokenHandler,
) *AuthHandler {
	return &AuthHandler{
		authenticateByToken: authenticateByToken,
	}
}

// Authenticate проверяет JWT токен и возвращает информацию о пользователе
func (h *AuthHandler) Authenticate(ctx context.Context, req *authv1.AuthenticateRequest) (*authv1.AuthenticateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if strings.TrimSpace(req.JwtToken) == "" {
		return nil, status.Error(codes.InvalidArgument, "jwt_token is required")
	}

	// Валидация JWT токена и извлечение данных из клеймов без обращения к БД
	info, err := h.authenticateByToken.Handle(ctx, queries.AuthenticateByTokenQuery{RawToken: req.JwtToken})
	if err != nil {
		return nil, h.convertErrorToGRPCStatus(err)
	}

	// Формируем ответ из данных клеймов
	response := &authv1.AuthenticateResponse{
		User: &authv1.User{
			Id:        info.ID.String(),
			Name:      info.Name,
			Email:     info.Email,
			Phone:     info.Phone,
			CreatedAt: timestamppb.New(info.CreatedAt),
		},
	}

	return response, nil
}

// convertErrorToGRPCStatus конвертирует доменные ошибки в gRPC статусы
func (h *AuthHandler) convertErrorToGRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	code := errs.ToGRPC(err)

	switch code {
	case codes.Unauthenticated:
		return status.Error(code, "invalid or expired JWT token")
	case codes.NotFound:
		return status.Error(code, "user not found")
	case codes.InvalidArgument:
		return status.Error(code, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
