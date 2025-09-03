package problems

import (
	"errors"

	"quest-auth/internal/generated/servers"
	"quest-auth/internal/pkg/errs"
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

	case isValidationError(err):
		return servers.Register400JSONResponse(servers.BadRequest{
			Type:   "bad-request",
			Title:  "Bad Request",
			Status: 400,
			Detail: err.Error(),
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
	var dvErr *errs.DomainValidationError
	if errors.As(err, &dvErr) {
		return dvErr.Field == "credentials"
	}
	return false
}

func isValidationError(err error) bool {
	var dvErr *errs.DomainValidationError
	if errors.As(err, &dvErr) {
		if dvErr.Field == "credentials" {
			return false
		}
		if dvErr.Field == "email" && dvErr.Message == "email already exists" {
			return false
		}
		if dvErr.Field == "phone" && dvErr.Message == "phone already exists" {
			return false
		}
		return true
	}
	return false
}

func isUserNotFoundError(err error) bool {
	var notFoundErr *errs.NotFoundError
	return errors.As(err, &notFoundErr)
}

func isEmailAlreadyExistsError(err error) bool {
	var dvErr *errs.DomainValidationError
	if errors.As(err, &dvErr) {
		return dvErr.Field == "email" && dvErr.Message == "email already exists"
	}
	return false
}

func isPhoneAlreadyExistsError(err error) bool {
	var dvErr *errs.DomainValidationError
	if errors.As(err, &dvErr) {
		return dvErr.Field == "phone" && dvErr.Message == "phone already exists"
	}
	return false
}
