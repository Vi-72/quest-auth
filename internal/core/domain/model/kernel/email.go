package kernel

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrEmailEmpty   = errors.New("email is empty")
	ErrEmailInvalid = errors.New("email is invalid")
)

var emailRe = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

type Email struct {
	value string
}

func NewEmail(s string) (Email, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return Email{}, ErrEmailEmpty
	}
	if !emailRe.MatchString(s) {
		return Email{}, fmt.Errorf("%w: %q", ErrEmailInvalid, s)
	}
	return Email{value: s}, nil
}

func (e Email) String() string { return e.value }

func (e Email) Equals(other Email) bool { return e.value == other.value }
