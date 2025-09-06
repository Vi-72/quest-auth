// DOMAIN LAYER UNIT TESTS
// Tests for domain model business rules and validation logic

package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"
)

func TestUser_NewUser_Success(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")

	u, err := auth.NewUser(email, phone, "John Doe", "password123", fakeHasher{}, fakeClock{})
	require.NoError(t, err)
	assert.Equal(t, "John Doe", u.Name)
	assert.Equal(t, email, u.Email)
	assert.Equal(t, phone, u.Phone)
	assert.True(t, u.VerifyPassword("password123", fakeHasher{}))
}

func TestUser_NewUser_Invalid(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")

	_, err := auth.NewUser(email, phone, "", "password123", fakeHasher{}, fakeClock{})
	require.Error(t, err)

	_, err = auth.NewUser(email, phone, "JD", "short", fakeHasher{}, fakeClock{})
	require.Error(t, err)
}

func TestUser_ChangeName(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", fakeHasher{}, fakeClock{})

	err := u.ChangeName("Jane Smith", fakeClock{})
	require.NoError(t, err)
	assert.Equal(t, "Jane Smith", u.Name)
}

func TestUser_ChangePhone(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", fakeHasher{}, fakeClock{})
	newPhone, _ := kernel.NewPhone("+1234567899")

	u.ChangePhone(newPhone, fakeClock{})
	assert.Equal(t, newPhone, u.Phone)
}

func TestUser_SetPassword(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", fakeHasher{}, fakeClock{})
	oldHash := u.PasswordHash

	err := u.SetPassword("newpassword456", fakeHasher{}, fakeClock{})
	require.NoError(t, err)
	assert.NotEqual(t, oldHash, u.PasswordHash)
	assert.True(t, u.VerifyPassword("newpassword456", fakeHasher{}))
}

func TestUser_DomainEvents(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", fakeHasher{}, fakeClock{})

	events := u.GetDomainEvents()
	require.NotEmpty(t, events)

	if e, ok := events[0].(interface{ GetAggregateID() uuid.UUID }); ok {
		assert.Equal(t, u.ID(), e.GetAggregateID())
	}
}
