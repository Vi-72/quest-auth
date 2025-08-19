package domain

import (
	"testing"

	"quest-auth/internal/core/domain/model/kernel"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPhone_ValidPhones_Success(t *testing.T) {
	validPhones := []string{
		"+1234567890",
		"+123456789012345", // 15 digits max (после +)
		"+19876543210",
		"+447700900123",  // UK
		"+33123456789",   // France
		"+8612345678901", // China
	}

	for _, phoneStr := range validPhones {
		t.Run(phoneStr, func(t *testing.T) {
			// Act
			phone, err := kernel.NewPhone(phoneStr)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, phoneStr, phone.String())
		})
	}
}

func TestNewPhone_InvalidPhones_ReturnsError(t *testing.T) {
	invalidPhones := []string{
		"",
		"1234567890",        // No +
		"+0123456789",       // Starts with 0 after +
		"+1234567890123456", // Too long (>15 digits after +)
		"+123456",           // Too short (need at least 7 digits after +)
		"+123abc456",        // Contains letters
		"+ 1234567890",      // Space after +
		"+1234 567890",      // Space in middle
		"+12345-67890",      // Contains dash
	}

	for _, phoneStr := range invalidPhones {
		t.Run(phoneStr, func(t *testing.T) {
			// Act
			_, err := kernel.NewPhone(phoneStr)

			// Assert
			assert.Error(t, err)
		})
	}
}

func TestNewPhone_EmptyString_ReturnsEmptyError(t *testing.T) {
	// Act
	_, err := kernel.NewPhone("")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "phone number is empty")
}

func TestNewPhone_InvalidFormat_ReturnsInvalidError(t *testing.T) {
	// Act
	_, err := kernel.NewPhone("1234567890")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "phone number is invalid")
	assert.Contains(t, err.Error(), "expected format: +1234567890")
}

func TestNewPhone_WhitespaceHandling(t *testing.T) {
	// Act
	phone, err := kernel.NewPhone("  +1234567890  ")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "+1234567890", phone.String()) // Должно быть trimmed
}

func TestPhone_Equals(t *testing.T) {
	// Arrange
	phone1, _ := kernel.NewPhone("+1234567890")
	phone2, _ := kernel.NewPhone("+1234567890")
	phone3, _ := kernel.NewPhone("+9876543210")

	// Act & Assert
	assert.True(t, phone1.Equals(phone2))
	assert.False(t, phone1.Equals(phone3))
}

func TestPhone_String(t *testing.T) {
	// Arrange
	phoneStr := "+1234567890"
	phone, _ := kernel.NewPhone(phoneStr)

	// Act
	result := phone.String()

	// Assert
	assert.Equal(t, phoneStr, result)
}
