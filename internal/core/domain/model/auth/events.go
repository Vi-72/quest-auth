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

// NewUserRegistered creates UserRegistered domain event at the specified time.
func NewUserRegistered(userID uuid.UUID, email, phone string, at time.Time) UserRegistered {
	return UserRegistered{
		ID:     uuid.New(),
		UserID: userID,
		Email:  email,
		Phone:  phone,
		At:     at,
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

// NewUserPhoneChanged creates UserPhoneChanged event at the given time.
func NewUserPhoneChanged(userID uuid.UUID, old, new string, at time.Time) UserPhoneChanged {
	return UserPhoneChanged{
		ID:     uuid.New(),
		UserID: userID,
		Old:    old,
		New:    new,
		At:     at,
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

// NewUserNameChanged creates UserNameChanged event at the specified time.
func NewUserNameChanged(userID uuid.UUID, old, new string, at time.Time) UserNameChanged {
	return UserNameChanged{
		ID:     uuid.New(),
		UserID: userID,
		Old:    old,
		New:    new,
		At:     at,
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

// NewUserPasswordChanged creates UserPasswordChanged event at the given time.
func NewUserPasswordChanged(userID uuid.UUID, at time.Time) UserPasswordChanged {
	return UserPasswordChanged{
		ID:     uuid.New(),
		UserID: userID,
		At:     at,
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
