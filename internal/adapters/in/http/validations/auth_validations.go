package validations

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"quest-auth/internal/generated/servers"
)

// RegisterUserRequest — запрос на регистрацию пользователя
type RegisterUserRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LoginUserRequest — запрос на вход пользователя
type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ValidationError — ошибка валидации
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed: field '%s' %s", e.Field, e.Message)
}

// ValidateRegisterUserRequest валидирует запрос на регистрацию
func ValidateRegisterUserRequest(body io.Reader) (RegisterUserRequest, error) {
	var req RegisterUserRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		return RegisterUserRequest{}, ValidationError{
			Field:   "body",
			Message: "invalid JSON format",
		}
	}

	// Валидация email
	if strings.TrimSpace(req.Email) == "" {
		return RegisterUserRequest{}, ValidationError{
			Field:   "email",
			Message: "is required",
		}
	}

	// Валидация phone
	if strings.TrimSpace(req.Phone) == "" {
		return RegisterUserRequest{}, ValidationError{
			Field:   "phone",
			Message: "is required",
		}
	}

	// Валидация name
	if strings.TrimSpace(req.Name) == "" {
		return RegisterUserRequest{}, ValidationError{
			Field:   "name",
			Message: "is required",
		}
	}

	// Валидация password
	if len(req.Password) == 0 {
		return RegisterUserRequest{}, ValidationError{
			Field:   "password",
			Message: "is required",
		}
	}

	return req, nil
}

// ValidateLoginUserRequest валидирует запрос на вход
func ValidateLoginUserRequest(body io.Reader) (LoginUserRequest, error) {
	var req LoginUserRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		return LoginUserRequest{}, ValidationError{
			Field:   "body",
			Message: "invalid JSON format",
		}
	}

	// Валидация email
	if strings.TrimSpace(req.Email) == "" {
		return LoginUserRequest{}, ValidationError{
			Field:   "email",
			Message: "is required",
		}
	}

	// Валидация password
	if len(req.Password) == 0 {
		return LoginUserRequest{}, ValidationError{
			Field:   "password",
			Message: "is required",
		}
	}

	return req, nil
}

// ValidateRegisterUserRequestBody валидирует OpenAPI тип для регистрации
func ValidateRegisterUserRequestBody(body *servers.RegisterJSONRequestBody) (RegisterUserRequest, error) {
	if body == nil {
		return RegisterUserRequest{}, ValidationError{
			Field:   "body",
			Message: "request body is required",
		}
	}

	// Конвертируем OpenAPI типы в наши внутренние типы
	req := RegisterUserRequest{
		Email:    string(body.Email),
		Phone:    body.Phone,
		Name:     body.Name,
		Password: body.Password,
	}

	// Валидация email
	if strings.TrimSpace(req.Email) == "" {
		return RegisterUserRequest{}, ValidationError{
			Field:   "email",
			Message: "is required",
		}
	}

	// Валидация phone
	if strings.TrimSpace(req.Phone) == "" {
		return RegisterUserRequest{}, ValidationError{
			Field:   "phone",
			Message: "is required",
		}
	}

	// Валидация name
	if strings.TrimSpace(req.Name) == "" {
		return RegisterUserRequest{}, ValidationError{
			Field:   "name",
			Message: "is required",
		}
	}

	// Валидация password
	if len(req.Password) == 0 {
		return RegisterUserRequest{}, ValidationError{
			Field:   "password",
			Message: "is required",
		}
	}

	return req, nil
}

// ValidateLoginUserRequestBody валидирует OpenAPI тип для входа
func ValidateLoginUserRequestBody(body *servers.LoginJSONRequestBody) (LoginUserRequest, error) {
	if body == nil {
		return LoginUserRequest{}, ValidationError{
			Field:   "body",
			Message: "request body is required",
		}
	}

	// Конвертируем OpenAPI типы в наши внутренние типы
	req := LoginUserRequest{
		Email:    string(body.Email),
		Password: body.Password,
	}

	// Валидация email
	if strings.TrimSpace(req.Email) == "" {
		return LoginUserRequest{}, ValidationError{
			Field:   "email",
			Message: "is required",
		}
	}

	// Валидация password
	if len(req.Password) == 0 {
		return LoginUserRequest{}, ValidationError{
			Field:   "password",
			Message: "is required",
		}
	}

	return req, nil
}
