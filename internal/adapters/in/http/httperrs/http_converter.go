package httperrs

import (
	"errors"
	stdhttp "net/http"
	"quest-auth/api/http/auth/v1"
	"quest-auth/internal/pkg/errs"
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
			StatusCode: stdhttp.StatusInternalServerError,
		}
	}

	// Check for domain validation errors
	var domainValidationErr *errs.DomainValidationError
	if errors.As(err, &domainValidationErr) {
		switch domainValidationErr.Field {
		case "email":
			if domainValidationErr.Message == "email already exists" {
				return HTTPError{
					Type:       "conflict",
					Title:      "Conflict",
					Status:     409,
					Detail:     "Email already exists",
					StatusCode: stdhttp.StatusConflict,
				}
			}
			return HTTPError{
				Type:       "bad-request",
				Title:      "Bad Request",
				Status:     400,
				Detail:     domainValidationErr.Error(),
				StatusCode: stdhttp.StatusBadRequest,
			}
		case "phone":
			if domainValidationErr.Message == "phone already exists" {
				return HTTPError{
					Type:       "conflict",
					Title:      "Conflict",
					Status:     409,
					Detail:     "Phone number already exists",
					StatusCode: stdhttp.StatusConflict,
				}
			}
			return HTTPError{
				Type:       "bad-request",
				Title:      "Bad Request",
				Status:     400,
				Detail:     domainValidationErr.Error(),
				StatusCode: stdhttp.StatusBadRequest,
			}
		case "credentials":
			return HTTPError{
				Type:       "unauthorized",
				Title:      "Unauthorized",
				Status:     401,
				Detail:     "Invalid credentials",
				StatusCode: stdhttp.StatusUnauthorized,
			}
		default:
			return HTTPError{
				Type:       "bad-request",
				Title:      "Bad Request",
				Status:     400,
				Detail:     domainValidationErr.Error(),
				StatusCode: stdhttp.StatusBadRequest,
			}
		}
	}

	// Check for not found errors
	var notFoundErr *errs.NotFoundError
	if errors.As(err, &notFoundErr) {
		return HTTPError{
			Type:       "not-found",
			Title:      "Not Found",
			Status:     404,
			Detail:     notFoundErr.Error(),
			StatusCode: stdhttp.StatusNotFound,
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

	// Default to internal server error
	return HTTPError{
		Type:       "internal-server-error",
		Title:      "Internal Server Error",
		Status:     500,
		Detail:     "An unexpected error occurred",
		StatusCode: stdhttp.StatusInternalServerError,
	}
}

// ToRegisterResponse converts error to Register strict response wrapper
func ToRegisterResponse(err error) v1.RegisterResponseObject {
	httpErr := ToHTTP(err)

	// OpenAPI schema only supports 400 and 500 for register endpoint
	if httpErr.StatusCode == stdhttp.StatusConflict {
		httpErr.Status = 400
		httpErr.StatusCode = stdhttp.StatusBadRequest
	}

	return v1.Register400JSONResponse(v1.BadRequest{
		Type:   httpErr.Type,
		Title:  httpErr.Title,
		Status: httpErr.Status,
		Detail: httpErr.Detail,
	})
}

// ToLoginResponse converts error to Login strict response wrapper
func ToLoginResponse(err error) v1.LoginResponseObject {
	httpErr := ToHTTP(err)

	switch httpErr.StatusCode {
	case stdhttp.StatusUnauthorized:
		return v1.Login401JSONResponse(v1.Unauthorized{
			Type:   httpErr.Type,
			Title:  httpErr.Title,
			Status: httpErr.Status,
			Detail: httpErr.Detail,
		})
	case stdhttp.StatusBadRequest:
		return v1.Login400JSONResponse(v1.BadRequest{
			Type:   httpErr.Type,
			Title:  httpErr.Title,
			Status: httpErr.Status,
			Detail: httpErr.Detail,
		})
	default:
		// Default to bad request for other errors
		return v1.Login400JSONResponse(v1.BadRequest{
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
	case stdhttp.StatusBadRequest:
		return "bad-request"
	case stdhttp.StatusUnauthorized:
		return "unauthorized"
	case stdhttp.StatusNotFound:
		return "not-found"
	case stdhttp.StatusConflict:
		return "conflict"
	case stdhttp.StatusInternalServerError:
		return "internal-server-error"
	default:
		return "error"
	}
}

func getTitleFromStatus(status int) string {
	switch status {
	case stdhttp.StatusBadRequest:
		return "Bad Request"
	case stdhttp.StatusUnauthorized:
		return "Unauthorized"
	case stdhttp.StatusNotFound:
		return "Not Found"
	case stdhttp.StatusConflict:
		return "Conflict"
	case stdhttp.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Error"
	}
}
