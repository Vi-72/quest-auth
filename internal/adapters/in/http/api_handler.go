package http

import (
	"quest-auth/internal/core/application/usecases/auth"
)

// APIHandler реализует StrictServerInterface для OpenAPI
type APIHandler struct {
	registerHandler *auth.RegisterUserHandler
	loginHandler    *auth.LoginUserHandler
}

func NewAPIHandler(
	registerHandler *auth.RegisterUserHandler,
	loginHandler *auth.LoginUserHandler,
) (*APIHandler, error) {
	return &APIHandler{
		registerHandler: registerHandler,
		loginHandler:    loginHandler,
	}, nil
}
