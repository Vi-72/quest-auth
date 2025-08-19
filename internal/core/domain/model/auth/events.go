package auth

import (
	"time"

	"github.com/google/uuid"
)

// Ниже — простые фабрики доменных событий. Структуры/интерфейсы ddd.Event
// берём из вашего пакета ddd. Если требуется другой формат — подстрою.

type UserRegistered struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Email  string
	Phone  string
	At     time.Time
}

func NewUserRegistered(userID uuid.UUID, email, phone string) UserRegistered {
	return UserRegistered{
		ID:     uuid.New(),
		UserID: userID,
		Email:  email,
		Phone:  phone,
		At:     time.Now(),
	}
}

func (e UserRegistered) GetID() uuid.UUID          { return e.ID }
func (e UserRegistered) GetName() string           { return "UserRegistered" }
func (e UserRegistered) GetAggregateID() uuid.UUID { return e.UserID }

type UserPhoneChanged struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Old    string
	New    string
	At     time.Time
}

func NewUserPhoneChanged(userID uuid.UUID, old, new string) UserPhoneChanged {
	return UserPhoneChanged{
		ID:     uuid.New(),
		UserID: userID,
		Old:    old,
		New:    new,
		At:     time.Now(),
	}
}

func (e UserPhoneChanged) GetID() uuid.UUID          { return e.ID }
func (e UserPhoneChanged) GetName() string           { return "UserPhoneChanged" }
func (e UserPhoneChanged) GetAggregateID() uuid.UUID { return e.UserID }

type UserNameChanged struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Old    string
	New    string
	At     time.Time
}

func NewUserNameChanged(userID uuid.UUID, old, new string) UserNameChanged {
	return UserNameChanged{
		ID:     uuid.New(),
		UserID: userID,
		Old:    old,
		New:    new,
		At:     time.Now(),
	}
}

func (e UserNameChanged) GetID() uuid.UUID          { return e.ID }
func (e UserNameChanged) GetName() string           { return "UserNameChanged" }
func (e UserNameChanged) GetAggregateID() uuid.UUID { return e.UserID }

type UserPasswordChanged struct {
	ID     uuid.UUID
	UserID uuid.UUID
	At     time.Time
}

func NewUserPasswordChanged(userID uuid.UUID) UserPasswordChanged {
	return UserPasswordChanged{
		ID:     uuid.New(),
		UserID: userID,
		At:     time.Now(),
	}
}

func (e UserPasswordChanged) GetID() uuid.UUID          { return e.ID }
func (e UserPasswordChanged) GetName() string           { return "UserPasswordChanged" }
func (e UserPasswordChanged) GetAggregateID() uuid.UUID { return e.UserID }

type UserLoggedIn struct {
	ID     uuid.UUID
	UserID uuid.UUID
	At     time.Time
}

func NewUserLoggedIn(userID uuid.UUID, at time.Time) UserLoggedIn {
	return UserLoggedIn{
		ID:     uuid.New(),
		UserID: userID,
		At:     at,
	}
}

func (e UserLoggedIn) GetID() uuid.UUID          { return e.ID }
func (e UserLoggedIn) GetName() string           { return "UserLoggedIn" }
func (e UserLoggedIn) GetAggregateID() uuid.UUID { return e.UserID }
