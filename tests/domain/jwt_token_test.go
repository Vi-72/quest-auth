package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"quest-auth/internal/core/domain/model/kernel"
)

func TestNewJwtToken_Valid(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"plain", "abc.def.ghi", "abc.def.ghi"},
		{"trim spaces", "  token  ", "token"},
		{"tabs/newlines", "\t\nxxx.yyy.zzz\n", "xxx.yyy.zzz"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tok, err := kernel.NewJwtToken(c.input)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, tok.String())
		})
	}
}

func TestNewJwtToken_Empty(t *testing.T) {
	for _, s := range []string{"", " ", "\n\t  "} {
		t.Run("empty_"+s, func(t *testing.T) {
			_, err := kernel.NewJwtToken(s)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, kernel.ErrJWTTokenEmpty))
		})
	}
}

func TestJwtToken_Equals(t *testing.T) {
	a, _ := kernel.NewJwtToken("token")
	b, _ := kernel.NewJwtToken(" token ")
	c, _ := kernel.NewJwtToken("other")

	assert.True(t, a.Equals(b))
	assert.False(t, a.Equals(c))
}
