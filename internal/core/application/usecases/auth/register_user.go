package auth

import (
	"context"

	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"
	"quest-auth/internal/core/ports"
	"quest-auth/internal/pkg/errs"
)

// RegisterUserCommand — команда для регистрации пользователя
type RegisterUserCommand struct {
	Email    string
	Phone    string
	Name     string
	Password string
}

// RegisterUserResult — результат регистрации
type RegisterUserResult struct {
	User         UserInfo
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
}

// RegisterUserHandler — обработчик регистрации пользователя
type RegisterUserHandler struct {
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
	jwtService     ports.JWTService
}

func NewRegisterUserHandler(
	unitOfWork ports.UnitOfWork,
	eventPublisher ports.EventPublisher,
	jwtService ports.JWTService,
) *RegisterUserHandler {
	return &RegisterUserHandler{
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
		jwtService:     jwtService,
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

	userRepo := h.unitOfWork.UserRepository()

	// Проверка уникальности email
	emailExists, err := userRepo.EmailExists(email)
	if err != nil {
		return RegisterUserResult{}, err
	}
	if emailExists {
		return RegisterUserResult{}, errs.NewDomainValidationError("email", "email already exists")
	}

	// Проверка уникальности phone
	phoneExists, err := userRepo.PhoneExists(phone)
	if err != nil {
		return RegisterUserResult{}, err
	}
	if phoneExists {
		return RegisterUserResult{}, errs.NewDomainValidationError("phone", "phone already exists")
	}

	// Создание доменного объекта User
	user, err := auth.NewUser(email, phone, cmd.Name, cmd.Password)
	if err != nil {
		return RegisterUserResult{}, errs.NewDomainValidationError("user", err.Error())
	}

	// Сохранение в транзакции
	err = h.unitOfWork.Execute(ctx, func() error {
		if err := userRepo.Create(&user); err != nil {
			return err
		}

		// Публикация доменных событий
		return h.eventPublisher.PublishDomainEvents(ctx, user.GetDomainEvents())
	})

	if err != nil {
		return RegisterUserResult{}, err
	}

	// Генерация токенов
	tokenPair, err := h.jwtService.GenerateTokenPair(user.ID(), user.Email.String())
	if err != nil {
		return RegisterUserResult{}, err
	}

	return RegisterUserResult{
		User: UserInfo{
			ID:    user.ID(),
			Email: user.Email.String(),
			Name:  user.Name,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}
