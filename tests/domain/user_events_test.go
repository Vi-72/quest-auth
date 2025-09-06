package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"
)

func TestUser_Events_OnNewUser(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")

	u, err := auth.NewUser(email, phone, "John Doe", "password123", FakeHasher{}, FakeClock{})
	require.NoError(t, err)

	events := u.GetDomainEvents()
	require.Len(t, events, 1)
	// First event should be user.registered
	assert.Equal(t, "user.registered", events[0].GetName())
}

func TestUser_Events_OnChangePhone(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", FakeHasher{}, FakeClock{})
	u.ClearDomainEvents()

	newPhone, _ := kernel.NewPhone("+1234567899")
	u.ChangePhone(newPhone, FakeClock{})

	events := u.GetDomainEvents()
	require.Len(t, events, 1)
	assert.Equal(t, "UserPhoneChanged", events[0].GetName())
}

func TestUser_Events_OnChangeName(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", FakeHasher{}, FakeClock{})
	u.ClearDomainEvents()

	err := u.ChangeName("Jane Smith", FakeClock{})
	require.NoError(t, err)

	events := u.GetDomainEvents()
	require.Len(t, events, 1)
	assert.Equal(t, "UserNameChanged", events[0].GetName())
}

func TestUser_Events_OnSetPassword(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", FakeHasher{}, FakeClock{})
	u.ClearDomainEvents()

	err := u.SetPassword("newpassword456", FakeHasher{}, FakeClock{})
	require.NoError(t, err)

	events := u.GetDomainEvents()
	require.Len(t, events, 1)
	assert.Equal(t, "UserPasswordChanged", events[0].GetName())
}

func TestUser_Events_OnLoggedIn(t *testing.T) {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", FakeHasher{}, FakeClock{})
	u.ClearDomainEvents()

	// simulate login at specific time for determinism of At ordering if needed
	time.Sleep(1 * time.Millisecond)
	u.MarkLoggedIn(FakeClock{})

	events := u.GetDomainEvents()
	require.Len(t, events, 1)
	assert.Equal(t, "user.login", events[0].GetName())
}
