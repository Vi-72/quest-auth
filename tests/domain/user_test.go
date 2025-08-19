package domain

import (
	"testing"

	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser_Success(t *testing.T) {
	// Arrange
	email, err := kernel.NewEmail("test@example.com")
	require.NoError(t, err)

	phone, err := kernel.NewPhone("+1234567890")
	require.NoError(t, err)

	name := "John Doe"
	password := "securepassword123"

	// Act
	user, err := auth.NewUser(email, phone, name, password)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, phone, user.Phone)
	assert.Equal(t, name, user.Name)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEqual(t, password, user.PasswordHash) // Пароль должен быть хеширован
	assert.NotZero(t, user.ID())
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())

	// Проверяем, что создалось доменное событие
	events := user.GetDomainEvents()
	assert.Len(t, events, 1)

	userRegistered, ok := events[0].(auth.UserRegistered)
	assert.True(t, ok)
	assert.Equal(t, user.ID(), userRegistered.UserID)
	assert.Equal(t, email.String(), userRegistered.Email)
	assert.Equal(t, phone.String(), userRegistered.Phone)
}

func TestNewUser_EmptyName_ReturnsError(t *testing.T) {
	// Arrange
	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")

	// Act
	_, err := auth.NewUser(email, phone, "", "password123")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name must not be empty")
}

func TestNewUser_ShortPassword_ReturnsError(t *testing.T) {
	// Arrange
	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")

	// Act
	_, err := auth.NewUser(email, phone, "John Doe", "123")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password must be at least 8 characters")
}

func TestUser_VerifyPassword_Success(t *testing.T) {
	// Arrange
	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	password := "securepassword123"

	user, err := auth.NewUser(email, phone, "John Doe", password)
	require.NoError(t, err)

	// Act & Assert
	assert.True(t, user.VerifyPassword(password))
	assert.False(t, user.VerifyPassword("wrongpassword"))
}

func TestUser_ChangePhone_Success(t *testing.T) {
	// Arrange
	email, _ := kernel.NewEmail("test@example.com")
	oldPhone, _ := kernel.NewPhone("+1234567890")
	newPhone, _ := kernel.NewPhone("+9876543210")

	user, err := auth.NewUser(email, oldPhone, "John Doe", "password123")
	require.NoError(t, err)

	user.ClearDomainEvents() // Очищаем события от создания

	// Act
	user.ChangePhone(newPhone)

	// Assert
	assert.Equal(t, newPhone, user.Phone)

	// Проверяем событие
	events := user.GetDomainEvents()
	assert.Len(t, events, 1)

	phoneChanged, ok := events[0].(auth.UserPhoneChanged)
	assert.True(t, ok)
	assert.Equal(t, user.ID(), phoneChanged.UserID)
	assert.Equal(t, oldPhone.String(), phoneChanged.Old)
	assert.Equal(t, newPhone.String(), phoneChanged.New)
}

func TestUser_ChangeName_Success(t *testing.T) {
	// Arrange
	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")

	user, err := auth.NewUser(email, phone, "John Doe", "password123")
	require.NoError(t, err)

	user.ClearDomainEvents() // Очищаем события от создания

	// Act
	err = user.ChangeName("Jane Smith")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Jane Smith", user.Name)

	// Проверяем событие
	events := user.GetDomainEvents()
	assert.Len(t, events, 1)

	nameChanged, ok := events[0].(auth.UserNameChanged)
	assert.True(t, ok)
	assert.Equal(t, user.ID(), nameChanged.UserID)
	assert.Equal(t, "John Doe", nameChanged.Old)
	assert.Equal(t, "Jane Smith", nameChanged.New)
}

func TestUser_ChangeName_EmptyName_ReturnsError(t *testing.T) {
	// Arrange
	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")

	user, err := auth.NewUser(email, phone, "John Doe", "password123")
	require.NoError(t, err)

	// Act
	err = user.ChangeName("")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name must not be empty")
}

func TestUser_SetPassword_Success(t *testing.T) {
	// Arrange
	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")

	user, err := auth.NewUser(email, phone, "John Doe", "oldpassword123")
	require.NoError(t, err)

	oldHash := user.PasswordHash
	user.ClearDomainEvents()

	// Act
	err = user.SetPassword("newpassword123")

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, oldHash, user.PasswordHash)
	assert.True(t, user.VerifyPassword("newpassword123"))
	assert.False(t, user.VerifyPassword("oldpassword123"))

	// Проверяем событие
	events := user.GetDomainEvents()
	assert.Len(t, events, 1)

	passwordChanged, ok := events[0].(auth.UserPasswordChanged)
	assert.True(t, ok)
	assert.Equal(t, user.ID(), passwordChanged.UserID)
}

func TestUser_MarkLoggedIn_CreatesEvent(t *testing.T) {
	// Arrange
	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")

	user, err := auth.NewUser(email, phone, "John Doe", "password123")
	require.NoError(t, err)

	user.ClearDomainEvents()

	// Act
	user.MarkLoggedIn()

	// Assert
	events := user.GetDomainEvents()
	assert.Len(t, events, 1)

	loggedIn, ok := events[0].(auth.UserLoggedIn)
	assert.True(t, ok)
	assert.Equal(t, user.ID(), loggedIn.UserID)
	assert.NotZero(t, loggedIn.At) // Проверяем что время установлено
}
