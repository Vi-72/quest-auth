package http

import (
	"quest-auth/internal/core/application/usecases/auth"
)

// ApiHandler реализует StrictServerInterface для OpenAPI
type ApiHandler struct {
	registerHandler *auth.RegisterUserHandler
	loginHandler    *auth.LoginUserHandler
}

func NewApiHandler(
	registerHandler *auth.RegisterUserHandler,
	loginHandler *auth.LoginUserHandler,
) (*ApiHandler, error) {
	return &ApiHandler{
		registerHandler: registerHandler,
		loginHandler:    loginHandler,
	}, nil
}
