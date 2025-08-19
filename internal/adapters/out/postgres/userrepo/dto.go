package userrepo

import (
	"time"

	"github.com/google/uuid"
)

// UserDTO — структура для работы с базой данных
type UserDTO struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Phone        string    `gorm:"uniqueIndex;not null"`
	Name         string    `gorm:"not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

// TableName определяет имя таблицы для GORM
func (UserDTO) TableName() string {
	return "users"
}
