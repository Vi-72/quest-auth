package http

import (
	"context"

	"quest-auth/internal/adapters/in/http/validations"
	"quest-auth/internal/core/application/usecases/auth"
	"quest-auth/internal/generated/servers"
	"quest-auth/internal/pkg/errs"
)

// Login implements POST /auth/login from OpenAPI.
func (a *APIHandler) Login(ctx context.Context, request servers.LoginRequestObject) (servers.LoginResponseObject, error) {
	// Validate login request body (includes nil check)
	validatedData, validationErr := validations.ValidateLoginUserRequestBody(request.Body)
	if validationErr != nil {
		// Use unified error converter
		return errs.ToLoginResponse(validationErr), nil
	}

	// Execute login command
	cmd := auth.LoginUserCommand{
		Email:    validatedData.Email,
		Password: validatedData.Password,
	}

	result, err := a.loginHandler.Handle(ctx, cmd)
	if err != nil {
		// Use unified error converter
		return errs.ToLoginResponse(err), nil
	}

	// Map result to response directly
	return servers.Login200JSONResponse(servers.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    int(result.ExpiresIn),
		User: servers.User{
			Id:    result.User.ID,
			Email: result.User.Email,
			Name:  result.User.Name,
		},
	}), nil
}
