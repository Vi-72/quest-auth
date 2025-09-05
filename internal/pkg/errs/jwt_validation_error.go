package errs

import "fmt"

// JWTValidationError represents errors that occur during JWT validation.
type JWTValidationError struct {
	Message string
	Cause   error
}

func (e *JWTValidationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("jwt validation error: %s (cause: %v)", e.Message, e.Cause)
	}
	return fmt.Sprintf("jwt validation error: %s", e.Message)
}

func (e *JWTValidationError) Unwrap() error { return e.Cause }

func NewJWTValidationError(message string) *JWTValidationError {
	return &JWTValidationError{Message: message}
}

func NewJWTValidationErrorWithCause(message string, cause error) *JWTValidationError {
	return &JWTValidationError{Message: message, Cause: cause}
}
