package http

import (
	"context"
	"quest-auth/internal/core/application/usecases/commands"

	"quest-auth/internal/adapters/in/http/httperrs"
	"quest-auth/internal/adapters/in/http/validations"
	"quest-auth/internal/generated/servers"
)

// Register implements POST /auth/register from OpenAPI.
func (a *APIHandler) Register(ctx context.Context, request servers.RegisterRequestObject) (servers.RegisterResponseObject, error) {
	// Validate register request body (includes nil check)
	validatedData, validationErr := validations.ValidateRegisterUserRequestBody(request.Body)
	if validationErr != nil {
		// Use unified error converter
		return httperrs.ToRegisterResponse(validationErr), nil
	}

	// Execute register command
	cmd := commands.RegisterUserCommand{
		Email:    validatedData.Email,
		Phone:    validatedData.Phone,
		Name:     validatedData.Name,
		Password: validatedData.Password,
	}

	result, err := a.registerHandler.Handle(ctx, cmd)
	if err != nil {
		// Use unified error converter
		return httperrs.ToRegisterResponse(err), nil
	}

	// Map result to response directly
	return servers.Register201JSONResponse(servers.RegisterResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    int(result.ExpiresIn),
		User: servers.User{
			Id:    result.User.ID,
			Email: result.User.Email,
			Name:  result.User.Name,
			Phone: &result.User.Phone,
		},
	}), nil
}
