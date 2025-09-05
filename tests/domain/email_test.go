package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"quest-auth/internal/core/domain/model/kernel"
)

func TestNewEmail_ValidInput(t *testing.T) {
	valid := []string{
		"user@example.com",
		"USER@EXAMPLE.COM",
		"user.name+tag@example.co.uk",
	}
	for _, s := range valid {
		t.Run(s, func(t *testing.T) {
			e, err := kernel.NewEmail(s)
			assert.NoError(t, err)
			assert.NotEmpty(t, e.String())
		})
	}
}

func TestNewEmail_InvalidInput(t *testing.T) {
	invalid := []string{
		"userexample.com",
		"user@",
		"@example.com",
		"user@.com",
	}
	for _, s := range invalid {
		t.Run(s, func(t *testing.T) {
			_, err := kernel.NewEmail(s)
			assert.Error(t, err)
		})
	}
}

func TestEmail_Empty(t *testing.T) {
	_, err := kernel.NewEmail("")
	assert.Error(t, err)
}

func TestEmail_EqualsAndString(t *testing.T) {
	e1, _ := kernel.NewEmail("User@Example.com")
	e2, _ := kernel.NewEmail("user@example.com")
	e3, _ := kernel.NewEmail("another@example.com")

	assert.True(t, e1.Equals(e2))
	assert.False(t, e1.Equals(e3))
	assert.Equal(t, "user@example.com", e1.String())
}
