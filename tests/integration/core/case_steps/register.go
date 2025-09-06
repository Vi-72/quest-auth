package casesteps

import (
	"context"

	"quest-auth/internal/core/application/usecases/commands"
	testdatagenerators "quest-auth/tests/integration/core/test_data_generators"
)

// RegisterUserStep registers a user through the command handler
func RegisterUserStep(ctx context.Context, handler *commands.RegisterUserHandler, email, phone, name, password string) (commands.RegisterUserResult, error) {
	cmd := commands.RegisterUserCommand{
		Email:    email,
		Phone:    phone,
		Name:     name,
		Password: password,
	}
	return handler.Handle(ctx, cmd)
}

// RegisterUserStepData registers user using a generated UserTestData
func RegisterUserStepData(ctx context.Context, handler *commands.RegisterUserHandler, data testdatagenerators.UserTestData) (commands.RegisterUserResult, error) {
	return RegisterUserStep(ctx, handler, data.Email, data.Phone, data.Name, data.Password)
}
