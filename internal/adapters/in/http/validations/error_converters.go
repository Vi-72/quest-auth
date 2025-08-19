package validations

import (
	"quest-auth/internal/adapters/in/http/problems"
	"quest-auth/internal/pkg/errs"
)

// ConvertValidationErrorToProblem конвертирует ValidationError в Problem Details
func ConvertValidationErrorToProblem(err *ValidationError) *problems.ProblemDetails {
	return problems.NewBadRequest(err.Error())
}

// ConvertDomainValidationErrorToProblem конвертирует DomainValidationError в Problem Details
func ConvertDomainValidationErrorToProblem(err *errs.DomainValidationError) *problems.ProblemDetails {
	return problems.NewBadRequest(err.Error())
}

// ConvertNotFoundErrorToProblem конвертирует NotFoundError в Problem Details
func ConvertNotFoundErrorToProblem(err *errs.NotFoundError) *problems.ProblemDetails {
	return problems.NewNotFound(err.Error())
}
