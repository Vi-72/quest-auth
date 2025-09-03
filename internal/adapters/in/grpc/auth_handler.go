package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	authpb "quest-auth/api/proto"
	"quest-auth/internal/core/ports"
	"quest-auth/internal/pkg/errs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthHandler реализует gRPC сервис аутентификации
type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	jwtService     ports.JWTService
	userRepository ports.UserRepository
}

// NewAuthHandler создает новый gRPC handler для аутентификации
func NewAuthHandler(
	jwtService ports.JWTService,
	userRepository ports.UserRepository,
) *AuthHandler {
	return &AuthHandler{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
}

// Authenticate проверяет JWT токен и возвращает информацию о пользователе
func (h *AuthHandler) Authenticate(ctx context.Context, req *authpb.AuthenticateRequest) (*authpb.AuthenticateResponse, error) {
	// Валидация входящего запроса
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if strings.TrimSpace(req.JwtToken) == "" {
		return nil, status.Error(codes.InvalidArgument, "jwt_token is required")
	}

	// Валидация JWT токена
	claims, err := h.jwtService.ValidateAccessToken(req.JwtToken)
	if err != nil {
		// Конвертируем доменные ошибки в gRPC статусы
		return nil, h.convertErrorToGRPCStatus(err)
	}

	// Получение пользователя из базы данных по ID из токена
	user, err := h.userRepository.GetByID(claims.UserID)
	if err != nil {
		return nil, h.convertErrorToGRPCStatus(err)
	}

	// Преобразование доменной модели в gRPC ответ
	response := &authpb.AuthenticateResponse{
		User: &authpb.User{
			Id:        user.ID().String(),
			FullName:  user.Name,
			Email:     user.Email.String(),
			Phone:     user.Phone.String(),
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}

	return response, nil
}

// convertErrorToGRPCStatus конвертирует доменные ошибки в gRPC статусы
func (h *AuthHandler) convertErrorToGRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	// Проверяем тип ошибки и конвертируем в соответствующий gRPC код
	switch {
	case isJWTValidationError(err):
		return status.Error(codes.Unauthenticated, "invalid or expired JWT token")
	case isUserNotFoundError(err):
		return status.Error(codes.NotFound, "user not found")
	case isInfrastructureError(err):
		return status.Error(codes.Internal, "internal server error")
	default:
		return status.Error(codes.Internal, fmt.Sprintf("unexpected error: %v", err))
	}
}

// Helper functions для определения типов ошибок

func isJWTValidationError(err error) bool {
	// Проверяем, является ли ошибка связанной с JWT валидацией
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "token") &&
		(strings.Contains(errMsg, "invalid") ||
			strings.Contains(errMsg, "expired") ||
			strings.Contains(errMsg, "parsing"))
}

func isUserNotFoundError(err error) bool {
	var notFoundErr *errs.NotFoundError
	return errors.As(err, &notFoundErr)
}

func isInfrastructureError(err error) bool {
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "infrastructure") ||
		strings.Contains(errMsg, "database") ||
		strings.Contains(errMsg, "connection")
}
