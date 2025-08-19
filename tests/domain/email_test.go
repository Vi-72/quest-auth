package domain

import (
	"testing"

	"quest-auth/internal/core/domain/model/kernel"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmail_ValidEmails_Success(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"user+tag@example.org",
		"user123@test-domain.com",
		"simple@test.io",
	}

	for _, emailStr := range validEmails {
		t.Run(emailStr, func(t *testing.T) {
			// Act
			email, err := kernel.NewEmail(emailStr)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, emailStr, email.String())
		})
	}
}

func TestNewEmail_InvalidEmails_ReturnsError(t *testing.T) {
	invalidEmails := []string{
		"",
		"invalid",
		"@domain.com",
		"user@",
		"user@domain",
		"user.domain.com",
		"user @domain.com",
		"user@domain .com",
	}

	for _, emailStr := range invalidEmails {
		t.Run(emailStr, func(t *testing.T) {
			// Act
			_, err := kernel.NewEmail(emailStr)

			// Assert
			assert.Error(t, err)
		})
	}
}

func TestNewEmail_EmptyString_ReturnsEmptyError(t *testing.T) {
	// Act
	_, err := kernel.NewEmail("")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email is empty")
}

func TestNewEmail_InvalidFormat_ReturnsInvalidError(t *testing.T) {
	// Act
	_, err := kernel.NewEmail("invalid-email")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email is invalid")
}

func TestNewEmail_WhitespaceHandling(t *testing.T) {
	// Act
	email, err := kernel.NewEmail("  TEST@EXAMPLE.COM  ")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", email.String()) // Должно быть lowercase и trimmed
}

func TestEmail_Equals(t *testing.T) {
	// Arrange
	email1, _ := kernel.NewEmail("test@example.com")
	email2, _ := kernel.NewEmail("test@example.com")
	email3, _ := kernel.NewEmail("other@example.com")

	// Act & Assert
	assert.True(t, email1.Equals(email2))
	assert.False(t, email1.Equals(email3))
}

func TestEmail_EqualsIgnoreCase(t *testing.T) {
	// Arrange
	email1, _ := kernel.NewEmail("Test@Example.Com")
	email2, _ := kernel.NewEmail("test@example.com")

	// Act & Assert
	assert.True(t, email1.Equals(email2)) // Должны быть равны после нормализации
}
