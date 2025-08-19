package ports

import (
	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"

	"github.com/google/uuid"
)

type UserRepository interface {
	// Create — сохранение нового пользователя
	Create(user auth.User) error

	// GetByID — поиск пользователя по ID
	GetByID(id uuid.UUID) (auth.User, error)

	// GetByEmail — поиск пользователя по email
	GetByEmail(email kernel.Email) (auth.User, error)

	// GetByPhone — поиск пользователя по телефону
	GetByPhone(phone kernel.Phone) (auth.User, error)

	// Update — обновление существующего пользователя
	Update(user auth.User) error

	// Delete — удаление пользователя (мягкое или жёсткое)
	Delete(id uuid.UUID) error

	// EmailExists — проверка существования email
	EmailExists(email kernel.Email) (bool, error)

	// PhoneExists — проверка существования телефона
	PhoneExists(phone kernel.Phone) (bool, error)
}
