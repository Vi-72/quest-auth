package http

import (
	"context"

	"quest-auth/internal/adapters/in/http/problems"
	"quest-auth/internal/adapters/in/http/validations"
	"quest-auth/internal/core/application/usecases/auth"
	"quest-auth/internal/generated/servers"

	"github.com/google/uuid"
)

// Login implements POST /auth/login from OpenAPI.
func (a *ApiHandler) Login(ctx context.Context, request servers.LoginRequestObject) (servers.LoginResponseObject, error) {
	// Validate login request body (includes nil check)
	validatedData, validationErr := validations.ValidateLoginUserRequestBody(request.Body)
	if validationErr != nil {
		// Use problems package to create structured error response
		return problems.ConvertToLoginResponse(validationErr), nil
	}

	// Execute login command
	cmd := auth.LoginUserCommand{
		Email:    validatedData.Email,
		Password: validatedData.Password,
	}

	result, err := a.loginHandler.Handle(cmd)
	if err != nil {
		// Pass error to middleware for proper handling (400 for validation, 404 for not found, 500 for infrastructure)
		return nil, err
	}

	// Parse UserID
	userID, err := uuid.Parse(result.UserID)
	if err != nil {
		badRequest := problems.NewBadRequest("invalid user id format")
		return servers.Login400JSONResponse(servers.BadRequest{
			Type:   badRequest.Type,
			Title:  badRequest.Title,
			Status: badRequest.Status,
			Detail: badRequest.Detail,
		}), nil
	}

	// Success response
	apiResult := servers.LoginResponse{
		UserID: userID,
		Email:  result.Email,
		Phone:  result.Phone,
		Name:   result.Name,
	}
	return servers.Login200JSONResponse(apiResult), nil
}
