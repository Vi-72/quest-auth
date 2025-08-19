package http

import (
	"context"

	"quest-auth/internal/adapters/in/http/problems"
	"quest-auth/internal/adapters/in/http/validations"
	"quest-auth/internal/core/application/usecases/auth"
	"quest-auth/internal/generated/servers"

	"github.com/google/uuid"
)

// Register implements POST /auth/register from OpenAPI.
func (h *ApiHandler) Register(ctx context.Context, request servers.RegisterRequestObject) (servers.RegisterResponseObject, error) {
	// Validate register request body (includes nil check)
	validatedData, validationErr := validations.ValidateRegisterUserRequestBody(request.Body)
	if validationErr != nil {
		// Use problems package to create structured error response
		return problems.ConvertToRegisterResponse(validationErr), nil
	}

	// Execute register command
	cmd := auth.RegisterUserCommand{
		Email:    validatedData.Email,
		Phone:    validatedData.Phone,
		Name:     validatedData.Name,
		Password: validatedData.Password,
	}

	result, err := h.registerHandler.Handle(cmd)
	if err != nil {
		// Use problems package to map domain errors to HTTP responses
		return problems.ConvertToRegisterResponse(err), nil
	}

	// Parse UserID
	userID, err := uuid.Parse(result.UserID)
	if err != nil {
		badRequest := problems.NewBadRequest("invalid user id format")
		return servers.Register400JSONResponse(servers.BadRequest{
			Type:   badRequest.Type,
			Title:  badRequest.Title,
			Status: badRequest.Status,
			Detail: badRequest.Detail,
		}), nil
	}

	// Success response
	apiResult := servers.RegisterResponse{
		UserID: userID,
		Email:  result.Email,
		Phone:  result.Phone,
		Name:   result.Name,
	}
	return servers.Register201JSONResponse(apiResult), nil
}
