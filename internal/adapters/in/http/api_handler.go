package http

import "quest-auth/internal/core/application/usecases/commands"

// APIHandler реализует StrictServerInterface для OpenAPI
type APIHandler struct {
	registerHandler *commands.RegisterUserHandler
	loginHandler    *commands.LoginUserHandler
}

func NewAPIHandler(
	registerHandler *commands.RegisterUserHandler,
	loginHandler *commands.LoginUserHandler,
) (*APIHandler, error) {
	return &APIHandler{
		registerHandler: registerHandler,
		loginHandler:    loginHandler,
	}, nil
}
