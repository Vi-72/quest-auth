package problems

import (
	"strings"

	"quest-auth/internal/generated/servers"
)

// ConvertToLoginResponse конвертирует доменные ошибки в OpenAPI LoginResponseObject
func ConvertToLoginResponse(err error) servers.LoginResponseObject {
	switch {
	case isInvalidCredentialsError(err):
		return servers.Login401JSONResponse(servers.Unauthorized{
			Type:   "unauthorized",
			Title:  "Unauthorized",
			Status: 401,
			Detail: "Invalid email or password",
		})

	case isValidationError(err):
		return servers.Login400JSONResponse(servers.BadRequest{
			Type:   "bad-request",
			Title:  "Bad Request",
			Status: 400,
			Detail: err.Error(),
		})

	case isUserNotFoundError(err):
		// Маскируем "user not found" как invalid credentials для безопасности
		return servers.Login401JSONResponse(servers.Unauthorized{
			Type:   "unauthorized",
			Title:  "Unauthorized",
			Status: 401,
			Detail: "Invalid email or password",
		})

	default:
		// Общая ошибка сервера
		return servers.Login400JSONResponse(servers.BadRequest{
			Type:   "bad-request",
			Title:  "Bad Request",
			Status: 400,
			Detail: "Login failed",
		})
	}
}

// ConvertToRegisterResponse конвертирует доменные ошибки в OpenAPI RegisterResponseObject
func ConvertToRegisterResponse(err error) servers.RegisterResponseObject {
	switch {
	case isValidationError(err):
		return servers.Register400JSONResponse(servers.BadRequest{
			Type:   "bad-request",
			Title:  "Bad Request",
			Status: 400,
			Detail: err.Error(),
		})

	case isEmailAlreadyExistsError(err):
		return servers.Register400JSONResponse(servers.BadRequest{
			Type:   "conflict",
			Title:  "Conflict",
			Status: 400, // Note: OpenAPI schema uses 400 for all errors, not 409
			Detail: "Email already exists",
		})

	case isPhoneAlreadyExistsError(err):
		return servers.Register400JSONResponse(servers.BadRequest{
			Type:   "conflict",
			Title:  "Conflict",
			Status: 400, // Note: OpenAPI schema uses 400 for all errors, not 409
			Detail: "Phone number already exists",
		})

	default:
		// Общая ошибка сервера
		return servers.Register400JSONResponse(servers.BadRequest{
			Type:   "bad-request",
			Title:  "Bad Request",
			Status: 400,
			Detail: "Registration failed",
		})
	}
}

// Helper functions для определения типов ошибок

func isInvalidCredentialsError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "invalid email or password") ||
		strings.Contains(msg, "invalid credentials") ||
		strings.Contains(msg, "wrong password")
}

func isValidationError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "email") && strings.Contains(msg, "invalid") ||
		strings.Contains(msg, "phone") && strings.Contains(msg, "invalid") ||
		strings.Contains(msg, "password") && strings.Contains(msg, "too short") ||
		strings.Contains(msg, "name") && strings.Contains(msg, "empty")
}

func isUserNotFoundError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "user not found") ||
		strings.Contains(msg, "user does not exist")
}

func isEmailAlreadyExistsError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "email") &&
		(strings.Contains(msg, "already exists") ||
			strings.Contains(msg, "duplicate") ||
			strings.Contains(msg, "taken"))
}

func isPhoneAlreadyExistsError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "phone") &&
		(strings.Contains(msg, "already exists") ||
			strings.Contains(msg, "duplicate") ||
			strings.Contains(msg, "taken"))
}
