//go:build integration

package auth

import (
	"testing"

	"quest-auth/internal/adapters/out/postgres/userrepo"
	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"
	"quest-auth/tests/integration/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create_Success(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := userrepo.NewRepository(testDB.DB)

	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	user, err := auth.NewUser(email, phone, "John Doe", "password123")
	require.NoError(t, err)

	// Act
	err = repo.Create(&user)

	// Assert
	require.NoError(t, err)

	// Проверяем, что пользователь сохранился
	savedUser, err := repo.GetByID(user.ID())
	require.NoError(t, err)
	assert.Equal(t, user.Email.String(), savedUser.Email.String())
	assert.Equal(t, user.Phone.String(), savedUser.Phone.String())
	assert.Equal(t, user.Name, savedUser.Name)
}

func TestUserRepository_GetByEmail_Success(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := userrepo.NewRepository(testDB.DB)

	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	user, err := auth.NewUser(email, phone, "John Doe", "password123")
	require.NoError(t, err)

	err = repo.Create(&user)
	require.NoError(t, err)

	// Act
	foundUser, err := repo.GetByEmail(email)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, user.ID(), foundUser.ID())
	assert.Equal(t, user.Email.String(), foundUser.Email.String())
}

func TestUserRepository_GetByPhone_Success(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := userrepo.NewRepository(testDB.DB)

	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	user, err := auth.NewUser(email, phone, "John Doe", "password123")
	require.NoError(t, err)

	err = repo.Create(&user)
	require.NoError(t, err)

	// Act
	foundUser, err := repo.GetByPhone(phone)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, user.ID(), foundUser.ID())
	assert.Equal(t, user.Phone.String(), foundUser.Phone.String())
}

func TestUserRepository_EmailExists_True(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := userrepo.NewRepository(testDB.DB)

	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	user, err := auth.NewUser(email, phone, "John Doe", "password123")
	require.NoError(t, err)

	err = repo.Create(&user)
	require.NoError(t, err)

	// Act
	exists, err := repo.EmailExists(email)

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestUserRepository_EmailExists_False(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := userrepo.NewRepository(testDB.DB)

	email, _ := kernel.NewEmail("nonexistent@example.com")

	// Act
	exists, err := repo.EmailExists(email)

	// Assert
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestUserRepository_PhoneExists_True(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := userrepo.NewRepository(testDB.DB)

	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	user, err := auth.NewUser(email, phone, "John Doe", "password123")
	require.NoError(t, err)

	err = repo.Create(&user)
	require.NoError(t, err)

	// Act
	exists, err := repo.PhoneExists(phone)

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestUserRepository_Update_Success(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := userrepo.NewRepository(testDB.DB)

	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	user, err := auth.NewUser(email, phone, "John Doe", "password123")
	require.NoError(t, err)

	err = repo.Create(&user)
	require.NoError(t, err)

	// Изменяем имя
	err = user.ChangeName("Jane Smith")
	require.NoError(t, err)

	// Act
	err = repo.Update(&user)

	// Assert
	require.NoError(t, err)

	// Проверяем, что изменения сохранились
	updatedUser, err := repo.GetByID(user.ID())
	require.NoError(t, err)
	assert.Equal(t, "Jane Smith", updatedUser.Name)
}

func TestUserRepository_Delete_Success(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := userrepo.NewRepository(testDB.DB)

	email, _ := kernel.NewEmail("test@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	user, err := auth.NewUser(email, phone, "John Doe", "password123")
	require.NoError(t, err)

	err = repo.Create(&user)
	require.NoError(t, err)

	// Act
	err = repo.Delete(user.ID())

	// Assert
	require.NoError(t, err)

	// Проверяем, что пользователь удален
	_, err = repo.GetByID(user.ID())
	assert.Error(t, err) // Должна быть ошибка NotFound
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := userrepo.NewRepository(testDB.DB)
	nonExistentID := uuid.New()

	// Act
	_, err := repo.GetByID(nonExistentID)

	// Assert
	assert.Error(t, err)
	// Можно добавить проверку на конкретный тип ошибки NotFoundError
}
