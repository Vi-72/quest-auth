package auth

import (
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
	UserID string
	Email  string
	Phone  string
	Name   string
}

// RegisterUserHandler — обработчик регистрации пользователя
type RegisterUserHandler struct {
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
}

func NewRegisterUserHandler(
	unitOfWork ports.UnitOfWork,
	eventPublisher ports.EventPublisher,
) *RegisterUserHandler {
	return &RegisterUserHandler{
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
	}
}

// Handle выполняет регистрацию пользователя
func (h *RegisterUserHandler) Handle(cmd RegisterUserCommand) (RegisterUserResult, error) {
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
	err = h.unitOfWork.Execute(func() error {
		if err := userRepo.Create(user); err != nil {
			return err
		}

		// Публикация доменных событий
		return h.eventPublisher.PublishDomainEvents(user.GetDomainEvents())
	})

	if err != nil {
		return RegisterUserResult{}, err
	}

	return RegisterUserResult{
		UserID: user.ID().String(),
		Email:  user.Email.String(),
		Phone:  user.Phone.String(),
		Name:   user.Name,
	}, nil
}
