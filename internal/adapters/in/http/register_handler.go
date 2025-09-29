package http

import (
	"context"

	v1 "github.com/Vi-72/quest-auth/api/http/auth/v1"
	"github.com/Vi-72/quest-auth/internal/core/application/usecases/commands"

	"github.com/Vi-72/quest-auth/internal/adapters/in/http/httperrs"
	"github.com/Vi-72/quest-auth/internal/adapters/in/http/validations"
)

// Register implements POST /auth/register from OpenAPI.
func (a *APIHandler) Register(
	ctx context.Context,
	request v1.RegisterRequestObject,
) (v1.RegisterResponseObject, error) {
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
	return v1.Register201JSONResponse(v1.RegisterResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    int(result.ExpiresIn),
		User: v1.User{
			Id:    result.User.ID,
			Email: result.User.Email,
			Name:  result.User.Name,
			Phone: &result.User.Phone,
		},
	}), nil
}
