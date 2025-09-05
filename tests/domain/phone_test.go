package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"quest-auth/internal/core/domain/model/kernel"
)

func TestNewPhone_ValidInput(t *testing.T) {
	valid := []string{
		"+1234567",
		"+123456789012345",
		"+79991234567",
	}
	for _, s := range valid {
		t.Run(s, func(t *testing.T) {
			p, err := kernel.NewPhone(s)
			assert.NoError(t, err)
			assert.Equal(t, s, p.String())
		})
	}
}

func TestNewPhone_InvalidInput(t *testing.T) {
	invalid := []string{
		"",
		"+12345",
		"1234567890",
		"+0123456789",
	}
	for _, s := range invalid {
		t.Run(s, func(t *testing.T) {
			_, err := kernel.NewPhone(s)
			assert.Error(t, err)
		})
	}
}

func TestPhone_Equals(t *testing.T) {
	p1, _ := kernel.NewPhone("+1234567890")
	p2, _ := kernel.NewPhone("+1234567890")
	p3, _ := kernel.NewPhone("+1234567899")
	assert.True(t, p1.Equals(p2))
	assert.False(t, p1.Equals(p3))
}
