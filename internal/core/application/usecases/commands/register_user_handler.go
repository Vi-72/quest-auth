package commands

import (
	"context"

	"github.com/Vi-72/quest-auth/internal/core/domain/model/auth"
	"github.com/Vi-72/quest-auth/internal/core/domain/model/kernel"
	"github.com/Vi-72/quest-auth/internal/core/ports"
	"github.com/Vi-72/quest-auth/internal/pkg/errs"
)

// RegisterUserHandler — обработчик регистрации пользователя
type RegisterUserHandler struct {
	txManager      ports.TransactionManager
	jwtService     ports.JWTService
	passwordHasher ports.PasswordHasher
	clock          ports.Clock
}

func NewRegisterUserHandler(
	txManager ports.TransactionManager,
	jwtService ports.JWTService,
	passwordHasher ports.PasswordHasher,
	clock ports.Clock,
) *RegisterUserHandler {
	return &RegisterUserHandler{
		txManager:      txManager,
		jwtService:     jwtService,
		passwordHasher: passwordHasher,
		clock:          clock,
	}
}

// Handle выполняет регистрацию пользователя
func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUserCommand) (RegisterUserResult, error) {
	// Валидация email
	email, err := kernel.NewEmail(cmd.Email)
	if err != nil {
		return RegisterUserResult{}, errs.NewDomainValidationError("email", err.Error())
	}

	// Валидация phone
	phone, err := kernel.NewPhone(cmd.Phone)
	if err != nil {
		return RegisterUserResult{}, errs.NewDomainValidationError("phone", err.Error())
	}

	var createdUser auth.User
	err = h.txManager.RunInTransaction(ctx, func(ctx context.Context, repos ports.Repositories) error {
		userRepo := repos.User

		emailExists, err := userRepo.EmailExists(email)
		if err != nil {
			return err
		}
		if emailExists {
			return errs.NewDomainValidationError("email", "email already exists")
		}

		phoneExists, err := userRepo.PhoneExists(phone)
		if err != nil {
			return err
		}
		if phoneExists {
			return errs.NewDomainValidationError("phone", "phone already exists")
		}

		user, err := auth.NewUser(email, phone, cmd.Name, cmd.Password, h.passwordHasher, h.clock)
		if err != nil {
			return errs.NewDomainValidationError("user", err.Error())
		}

		if err := userRepo.Create(&user); err != nil {
			return err
		}

		if repos.Event != nil {
			if err := repos.Event.Publish(ctx, user.GetDomainEvents()...); err != nil {
				return err
			}
		}
		user.ClearDomainEvents()

		createdUser = user
		return nil
	})
	if err != nil {
		return RegisterUserResult{}, err
	}

	user := createdUser

	// Генерация токенов
	tokenPair, err := h.jwtService.GenerateTokenPair(
		user.ID(),
		user.Email.String(),
		user.Name,
		user.Phone.String(),
		user.CreatedAt,
	)
	if err != nil {
		return RegisterUserResult{}, err
	}

	return RegisterUserResult{
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
