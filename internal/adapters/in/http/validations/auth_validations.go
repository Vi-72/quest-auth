package validations

import (
	"fmt"
	"quest-auth/api/http"
	"strings"
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

// requireNotEmpty проверяет, что строковое значение не пустое после TrimSpace.
// Если значение пустое, возвращает ValidationError для указанного поля.
func requireNotEmpty(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return ValidationError{
			Field:   field,
			Message: "is required",
		}
	}
	return nil
}

// requireNotZeroLen проверяет, что строка не пустая (без тримминга).
// Используется, например, для проверки пароля.
func requireNotZeroLen(field, value string) error {
	if len(value) == 0 {
		return ValidationError{
			Field:   field,
			Message: "is required",
		}
	}
	return nil
}

// ValidateRegisterUserRequestBody валидирует OpenAPI тип для регистрации
func ValidateRegisterUserRequestBody(body *http.RegisterJSONRequestBody) (RegisterUserRequest, error) {
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
	if err := requireNotEmpty("email", req.Email); err != nil {
		return RegisterUserRequest{}, err
	}

	// Валидация phone
	if err := requireNotEmpty("phone", req.Phone); err != nil {
		return RegisterUserRequest{}, err
	}

	// Валидация name
	if err := requireNotEmpty("name", req.Name); err != nil {
		return RegisterUserRequest{}, err
	}

	// Валидация password
	if err := requireNotZeroLen("password", req.Password); err != nil {
		return RegisterUserRequest{}, err
	}

	return req, nil
}

// ValidateLoginUserRequestBody валидирует OpenAPI тип для входа
func ValidateLoginUserRequestBody(body *http.LoginJSONRequestBody) (LoginUserRequest, error) {
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
	if err := requireNotEmpty("email", req.Email); err != nil {
		return LoginUserRequest{}, err
	}

	// Валидация password
	if err := requireNotZeroLen("password", req.Password); err != nil {
		return LoginUserRequest{}, err
	}

	return req, nil
}
