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

	result, err := a.loginHandler.Handle(ctx, cmd)
	if err != nil {
		// Pass error to middleware for proper handling (400 for validation, 404 for not found, 500 for infrastructure)
		return nil, err
	}

	// Success response
	apiResult := AuthResponse{
		User: UserResponse{
			ID:    result.User.ID,
			Email: result.User.Email,
			Name:  result.User.Name,
		},
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
	}
	// Parse user ID
	userID, err := uuid.Parse(apiResult.User.ID)
	if err != nil {
		return nil, err
	}

	return servers.Login200JSONResponse(servers.LoginResponse{
		AccessToken:  apiResult.AccessToken,
		RefreshToken: apiResult.RefreshToken,
		TokenType:    apiResult.TokenType,
		ExpiresIn:    int(apiResult.ExpiresIn),
		User: servers.User{
			Id:    userID,
			Email: apiResult.User.Email,
			Name:  apiResult.User.Name,
		},
	}), nil
}
