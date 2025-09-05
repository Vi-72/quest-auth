package auth

import (
	"errors"
	"time"

	"quest-auth/internal/core/domain/model/kernel"
	clockpkg "quest-auth/internal/pkg/clock"
	"quest-auth/internal/pkg/ddd"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNameEmpty        = errors.New("name must not be empty")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)

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
func NewUser(clock clockpkg.Clock, email kernel.Email, phone kernel.Phone, name string, rawPassword string) (User, error) {
	if name = normalizeName(name); name == "" {
		return User{}, ErrNameEmpty
	}
	if len(rawPassword) < 8 {
		return User{}, ErrPasswordTooShort
	}

	hash, err := hashPassword(rawPassword)
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
func (u *User) ChangePhone(clock clockpkg.Clock, newPhone kernel.Phone) {
	old := u.Phone
	u.Phone = newPhone
	now := clock.Now()
	u.UpdatedAt = now
	u.RaiseDomainEvent(NewUserPhoneChanged(u.ID(), old.String(), newPhone.String(), now))
}

// ChangeName — обновление отображаемого имени.
func (u *User) ChangeName(clock clockpkg.Clock, newName string) error {
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
func (u *User) SetPassword(clock clockpkg.Clock, rawPassword string) error {
	if len(rawPassword) < 8 {
		return ErrPasswordTooShort
	}
	hash, err := hashPassword(rawPassword)
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
func (u *User) VerifyPassword(raw string) bool {
	if u.PasswordHash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(raw)) == nil
}

// MarkLoggedIn — доменное событие логина (можно вызывать после VerifyPassword).
func (u *User) MarkLoggedIn(clock clockpkg.Clock) {
	u.RaiseDomainEvent(NewUserLoggedIn(u.ID(), clock.Now()))
}

// Вспомогательные функции
func hashPassword(raw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

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
