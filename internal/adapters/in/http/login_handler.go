package http

import (
	"context"
	"quest-auth/internal/core/application/usecases/commands"

	"quest-auth/internal/adapters/in/http/httperrs"
	"quest-auth/internal/adapters/in/http/validations"
)

// Login implements POST /auth/login from OpenAPI.
func (a *APIHandler) Login(ctx context.Context, request http.LoginRequestObject) (http.LoginResponseObject, error) {
	// Validate login request body (includes nil check)
	validatedData, validationErr := validations.ValidateLoginUserRequestBody(request.Body)
	if validationErr != nil {
		// Use unified error converter
		return httperrs.ToLoginResponse(validationErr), nil
	}

	// Execute login command
	cmd := commands.LoginUserCommand{
		Email:    validatedData.Email,
		Password: validatedData.Password,
	}

	result, err := a.loginHandler.Handle(ctx, cmd)
	if err != nil {
		// Use unified error converter
		return httperrs.ToLoginResponse(err), nil
	}

	// Map result to response directly
	return http.Login200JSONResponse(http.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    int(result.ExpiresIn),
		User: http.User{
			Id:    result.User.ID,
			Email: result.User.Email,
			Name:  result.User.Name,
			Phone: &result.User.Phone,
		},
	}), nil
}
