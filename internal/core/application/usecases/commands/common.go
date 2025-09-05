package commands

import "github.com/google/uuid"

// UserInfo — информация о пользователе
type UserInfo struct {
	ID    uuid.UUID
	Email string
	Name  string
	Phone string
}
