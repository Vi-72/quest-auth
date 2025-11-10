package http

import (
	"context"

	v1 "github.com/Vi-72/quest-auth/api/http/auth/v1"
	"github.com/Vi-72/quest-auth/internal/core/application/usecases/commands"

	"github.com/Vi-72/quest-auth/internal/adapters/in/http/httperrs"
)

// Login implements POST /auth/login from OpenAPI.
func (a *APIHandler) Login(ctx context.Context, request v1.LoginRequestObject) (v1.LoginResponseObject, error) {
	// OpenAPI validation middleware already validated the request
	// Just extract the data from request.Body
	body := request.Body

	// Execute login command
	cmd := commands.LoginUserCommand{
		Email:    string(body.Email),
		Password: body.Password,
	}

	result, err := a.loginHandler.Handle(ctx, cmd)
	if err != nil {
		// Use unified error converter
		return httperrs.ToLoginResponse(err), nil
	}

	// Map result to response directly
	return v1.Login200JSONResponse(v1.LoginResponse{
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
