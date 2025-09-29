package casesteps

import (
	"context"

	"github.com/Vi-72/quest-auth/internal/core/application/usecases/commands"
)

// LoginUserStep logs user in through the command handler
func LoginUserStep(ctx context.Context, handler *commands.LoginUserHandler, email, password string) (commands.LoginUserResult, error) {
	cmd := commands.LoginUserCommand{
		Email:    email,
		Password: password,
	}
	return handler.Handle(ctx, cmd)
}
