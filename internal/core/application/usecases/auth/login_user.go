package auth

import (
	"quest-auth/internal/core/domain/model/kernel"
	"quest-auth/internal/core/ports"
	"quest-auth/internal/pkg/errs"
)

// LoginUserCommand — команда для входа пользователя
type LoginUserCommand struct {
	Email    string
	Password string
}

// LoginUserResult — результат входа
type LoginUserResult struct {
	UserID string
	Email  string
	Phone  string
	Name   string
}

// LoginUserHandler — обработчик входа пользователя
type LoginUserHandler struct {
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
}

func NewLoginUserHandler(
	unitOfWork ports.UnitOfWork,
	eventPublisher ports.EventPublisher,
) *LoginUserHandler {
	return &LoginUserHandler{
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
	}
}

// Handle выполняет вход пользователя
func (h *LoginUserHandler) Handle(cmd LoginUserCommand) (LoginUserResult, error) {
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
	user.MarkLoggedIn()

	// Публикация доменных событий
	err = h.eventPublisher.PublishDomainEvents(user.GetDomainEvents())
	if err != nil {
		return LoginUserResult{}, err
	}

	return LoginUserResult{
		UserID: user.ID().String(),
		Email:  user.Email.String(),
		Phone:  user.Phone.String(),
		Name:   user.Name,
	}, nil
}
