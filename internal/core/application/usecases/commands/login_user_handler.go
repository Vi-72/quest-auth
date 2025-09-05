package commands

import (
	"context"

	"quest-auth/internal/core/domain/model/kernel"
	"quest-auth/internal/core/ports"
	clockpkg "quest-auth/internal/pkg/clock"
	"quest-auth/internal/pkg/errs"
)

// LoginUserHandler — обработчик входа пользователя
type LoginUserHandler struct {
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
	jwtService     ports.JWTService
	clock          clockpkg.Clock
}

func NewLoginUserHandler(
	unitOfWork ports.UnitOfWork,
	eventPublisher ports.EventPublisher,
	jwtService ports.JWTService,
	clock clockpkg.Clock,
) *LoginUserHandler {
	return &LoginUserHandler{
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
		jwtService:     jwtService,
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

	userRepo := h.unitOfWork.UserRepository()
	// Поиск пользователя по email
	user, err := userRepo.GetByEmail(email)
	if err != nil {
		return LoginUserResult{}, errs.NewDomainValidationError("credentials", "invalid email or password")
	}

	// Проверка пароля
	if !user.VerifyPassword(cmd.Password) {
		return LoginUserResult{}, errs.NewDomainValidationError("credentials", "invalid email or password")
	}

	// Отметка о входе (создание доменного события)
	user.MarkLoggedIn(h.clock)

	// Публикация доменных событий
	err = h.eventPublisher.PublishDomainEvents(ctx, user.GetDomainEvents())
	if err != nil {
		return LoginUserResult{}, err
	}

	// Генерация токенов
	tokenPair, err := h.jwtService.GenerateTokenPair(user.ID(), user.Email.String(), user.Name, user.Phone.String(), user.CreatedAt)
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
