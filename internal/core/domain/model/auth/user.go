package auth

import (
	"errors"
	"time"

	"quest-auth/internal/core/domain/model/kernel"
	"quest-auth/internal/pkg/ddd"

	"github.com/google/uuid"
)

var (
	ErrNameEmpty        = errors.New("name must not be empty")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)

// PasswordHasher provides methods to hash and compare passwords.
type PasswordHasher interface {
	Hash(raw string) (string, error)
	Compare(hash, raw string) bool
}

// Clock provides current time.
type Clock interface {
	Now() time.Time
}

// User — агрегат домена аутентификации.
type User struct {
	*ddd.BaseAggregate[uuid.UUID]

	Email        kernel.Email
	Phone        kernel.Phone
	Name         string
	PasswordHash string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser — регистрация пользователя (создание аккаунта).
// Сразу валидирует email/phone/name и хеширует пароль.
func NewUser(email kernel.Email, phone kernel.Phone, name string, rawPassword string, hasher PasswordHasher, clock Clock) (User, error) {
	if name = normalizeName(name); name == "" {
		return User{}, ErrNameEmpty
	}
	if len(rawPassword) < 8 {
		return User{}, ErrPasswordTooShort
	}

	hash, err := hasher.Hash(rawPassword)
	if err != nil {
		return User{}, err
	}

	id := uuid.New()
	now := clock.Now()

	u := User{
		BaseAggregate: ddd.NewBaseAggregate(id),
		Email:         email,
		Phone:         phone,
		Name:          name,
		PasswordHash:  hash,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	u.RaiseDomainEvent(NewUserRegistered(id, email.String(), phone.String(), now))
	return u, nil
}

// ChangePhone — смена телефона (например, после подтверждения OTP).
func (u *User) ChangePhone(newPhone kernel.Phone, clock Clock) {
	old := u.Phone
	u.Phone = newPhone
	now := clock.Now()
	u.UpdatedAt = now
	u.RaiseDomainEvent(NewUserPhoneChanged(u.ID(), old.String(), newPhone.String(), now))
}

// ChangeName — обновление отображаемого имени.
func (u *User) ChangeName(newName string, clock Clock) error {
	newName = normalizeName(newName)
	if newName == "" {
		return ErrNameEmpty
	}
	if newName == u.Name {
		return nil
	}
	old := u.Name
	u.Name = newName
	now := clock.Now()
	u.UpdatedAt = now
	u.RaiseDomainEvent(NewUserNameChanged(u.ID(), old, newName, now))
	return nil
}

// SetPassword — смена пароля (с валидацией и перезаписью хеша).
func (u *User) SetPassword(rawPassword string, hasher PasswordHasher, clock Clock) error {
	if len(rawPassword) < 8 {
		return ErrPasswordTooShort
	}
	hash, err := hasher.Hash(rawPassword)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	now := clock.Now()
	u.UpdatedAt = now
	u.RaiseDomainEvent(NewUserPasswordChanged(u.ID(), now))
	return nil
}

// VerifyPassword — проверка пароля при логине.
func (u *User) VerifyPassword(raw string, hasher PasswordHasher) bool {
	if u.PasswordHash == "" {
		return false
	}
	return hasher.Compare(u.PasswordHash, raw)
}

// MarkLoggedIn — доменное событие логина (можно вызывать после VerifyPassword).
func (u *User) MarkLoggedIn(clock Clock) {
	u.RaiseDomainEvent(NewUserLoggedIn(u.ID(), clock.Now()))
}

// Вспомогательные функции
func normalizeName(s string) string {
	// лёгкая нормализация; можно добавить unicode.TrimSpace/Title
	if len(s) == 0 {
		return s
	}
	// обрезка пробелов
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\n') {
		s = s[:len(s)-1]
	}
	return s
}
