package http

import (
	"context"

	v1 "github.com/Vi-72/quest-auth/api/http/auth/v1"
	"github.com/Vi-72/quest-auth/internal/core/application/usecases/commands"

	"github.com/Vi-72/quest-auth/internal/adapters/in/http/httperrs"
)

// Register implements POST /auth/register from OpenAPI.
func (a *APIHandler) Register(
	ctx context.Context,
	request v1.RegisterRequestObject,
) (v1.RegisterResponseObject, error) {
	// OpenAPI validation middleware already validated the request
	// Just extract the data from request.Body
	body := request.Body

	// Execute register command
	cmd := commands.RegisterUserCommand{
		Email:    string(body.Email),
		Phone:    body.Phone,
		Name:     body.Name,
		Password: body.Password,
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
