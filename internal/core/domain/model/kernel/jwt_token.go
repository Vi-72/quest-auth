package kernel

import (
	"errors"
	"strings"
)

var (
	ErrJWTTokenEmpty = errors.New("jwt token is empty")
)

// JwtToken is a value object representing a raw JWT token string.
type JwtToken struct {
	value string
}

// NewJwtToken trims the input and validates that it is non-empty.
func NewJwtToken(s string) (JwtToken, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return JwtToken{}, ErrJWTTokenEmpty
	}
	return JwtToken{value: s}, nil
}

func (t JwtToken) String() string { return t.value }

func (t JwtToken) Equals(other JwtToken) bool { return t.value == other.value }
