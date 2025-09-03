package errs

import (
	"errors"
	"net/http"
	"strings"

	"quest-auth/internal/generated/servers"
)

// HTTPError represents a structured HTTP error response
type HTTPError struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Status     int    `json:"status"`
	Detail     string `json:"detail"`
	StatusCode int    `json:"-"` // For HTTP response status
}

func (e HTTPError) Error() string {
	return e.Detail
}

// ToHTTP converts any error to a structured HTTP error
func ToHTTP(err error) HTTPError {
	if err == nil {
		return HTTPError{
			Type:       "internal-server-error",
			Title:      "Internal Server Error",
			Status:     500,
			Detail:     "An unexpected error occurred",
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Check for domain validation errors
	var domainValidationErr *DomainValidationError
	if errors.As(err, &domainValidationErr) {
		return HTTPError{
			Type:       "bad-request",
			Title:      "Bad Request",
			Status:     400,
			Detail:     domainValidationErr.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	// Check for not found errors
	var notFoundErr *NotFoundError
	if errors.As(err, &notFoundErr) {
		return HTTPError{
			Type:       "not-found",
			Title:      "Not Found",
			Status:     404,
			Detail:     notFoundErr.Error(),
			StatusCode: http.StatusNotFound,
		}
	}

	// Check for errors with status
	var statusErr *ErrorWithStatus
	if errors.As(err, &statusErr) {
		return HTTPError{
			Type:       getTypeFromStatus(statusErr.StatusCode),
			Title:      getTitleFromStatus(statusErr.StatusCode),
			Status:     statusErr.StatusCode,
			Detail:     statusErr.Error(),
			StatusCode: statusErr.StatusCode,
		}
	}

	// Check for specific validation patterns
	errMsg := strings.ToLower(err.Error())

	// Email/Phone already exists
	if strings.Contains(errMsg, "email already exists") {
		return HTTPError{
			Type:       "conflict",
			Title:      "Conflict",
			Status:     409,
			Detail:     "Email already exists",
			StatusCode: http.StatusConflict,
		}
	}

	if strings.Contains(errMsg, "phone already exists") {
		return HTTPError{
			Type:       "conflict",
			Title:      "Conflict",
			Status:     409,
			Detail:     "Phone number already exists",
			StatusCode: http.StatusConflict,
		}
	}

	// Invalid credentials
	if strings.Contains(errMsg, "invalid email or password") ||
		strings.Contains(errMsg, "invalid credentials") {
		return HTTPError{
			Type:       "unauthorized",
			Title:      "Unauthorized",
			Status:     401,
			Detail:     "Invalid credentials",
			StatusCode: http.StatusUnauthorized,
		}
	}

	// Validation errors
	if isValidationError(errMsg) {
		return HTTPError{
			Type:       "bad-request",
			Title:      "Bad Request",
			Status:     400,
			Detail:     err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	// Default to internal server error
	return HTTPError{
		Type:       "internal-server-error",
		Title:      "Internal Server Error",
		Status:     500,
		Detail:     "An unexpected error occurred",
		StatusCode: http.StatusInternalServerError,
	}
}

// ToRegisterResponse converts error to OpenAPI RegisterResponseObject
func ToRegisterResponse(err error) servers.RegisterResponseObject {
	httpErr := ToHTTP(err)

	// OpenAPI schema only supports 400 and 500 for register endpoint
	if httpErr.StatusCode == http.StatusConflict {
		httpErr.Status = 400
		httpErr.StatusCode = http.StatusBadRequest
	}

	return servers.Register400JSONResponse(servers.BadRequest{
		Type:   httpErr.Type,
		Title:  httpErr.Title,
		Status: httpErr.Status,
		Detail: httpErr.Detail,
	})
}

// ToLoginResponse converts error to OpenAPI LoginResponseObject
func ToLoginResponse(err error) servers.LoginResponseObject {
	httpErr := ToHTTP(err)

	switch httpErr.StatusCode {
	case http.StatusUnauthorized:
		return servers.Login401JSONResponse(servers.Unauthorized{
			Type:   httpErr.Type,
			Title:  httpErr.Title,
			Status: httpErr.Status,
			Detail: httpErr.Detail,
		})
	case http.StatusBadRequest:
		return servers.Login400JSONResponse(servers.BadRequest{
			Type:   httpErr.Type,
			Title:  httpErr.Title,
			Status: httpErr.Status,
			Detail: httpErr.Detail,
		})
	default:
		// Default to bad request for other errors
		return servers.Login400JSONResponse(servers.BadRequest{
			Type:   "bad-request",
			Title:  "Bad Request",
			Status: 400,
			Detail: httpErr.Detail,
		})
	}
}

// Helper functions
func getTypeFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad-request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusNotFound:
		return "not-found"
	case http.StatusConflict:
		return "conflict"
	case http.StatusInternalServerError:
		return "internal-server-error"
	default:
		return "error"
	}
}

func getTitleFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusConflict:
		return "Conflict"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Error"
	}
}

func isValidationError(errMsg string) bool {
	validationKeywords := []string{
		"email",
		"phone",
		"password",
		"name",
		"invalid",
		"required",
		"empty",
		"too short",
		"validation",
	}

	for _, keyword := range validationKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}
	return false
}
