package commands

import (
	"context"

	"github.com/Vi-72/quest-auth/internal/core/domain/model/auth"
	"github.com/Vi-72/quest-auth/internal/core/domain/model/kernel"
	"github.com/Vi-72/quest-auth/internal/core/ports"
	"github.com/Vi-72/quest-auth/internal/pkg/errs"
)

// LoginUserHandler — обработчик входа пользователя
type LoginUserHandler struct {
	txManager      ports.TransactionManager
	jwtService     ports.JWTService
	passwordHasher ports.PasswordHasher
	clock          ports.Clock
}

func NewLoginUserHandler(
	txManager ports.TransactionManager,
	jwtService ports.JWTService,
	passwordHasher ports.PasswordHasher,
	clock ports.Clock,
) *LoginUserHandler {
	return &LoginUserHandler{
		txManager:      txManager,
		jwtService:     jwtService,
		passwordHasher: passwordHasher,
		clock:          clock,
	}
}

// Handle выполняет вход пользователя
func (h *LoginUserHandler) Handle(ctx context.Context, cmd LoginUserCommand) (LoginUserResult, error) {
	// Валидация email
	email, err := kernel.NewEmail(cmd.Email)
	if err != nil {
		return LoginUserResult{}, errs.NewDomainValidationError("email", err.Error())
	}

	var loggedInUser *auth.User
	err = h.txManager.RunInTransaction(ctx, func(ctx context.Context, repos ports.Repositories) error {
		userRepo := repos.User

		user, err := userRepo.GetByEmail(email)
		if err != nil {
			return errs.NewDomainValidationError("credentials", "invalid email or password")
		}

		if !user.VerifyPassword(cmd.Password, h.passwordHasher) {
			return errs.NewDomainValidationError("credentials", "invalid email or password")
		}

		user.MarkLoggedIn(h.clock)

		if repos.Event != nil {
			if err := repos.Event.Publish(ctx, user.GetDomainEvents()...); err != nil {
				return err
			}
		}
		user.ClearDomainEvents()

		loggedInUser = user
		return nil
	})
	if err != nil {
		return LoginUserResult{}, err
	}

	user := loggedInUser

	// Генерация токенов
	tokenPair, err := h.jwtService.GenerateTokenPair(
		user.ID(),
		user.Email.String(),
		user.Name,
		user.Phone.String(),
		user.CreatedAt,
	)
	if err != nil {
		return LoginUserResult{}, err
	}

	return LoginUserResult{
		User: UserInfo{
			ID:    user.ID(),
			Email: user.Email.String(),
			Name:  user.Name,
			Phone: user.Phone.String(),
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}
