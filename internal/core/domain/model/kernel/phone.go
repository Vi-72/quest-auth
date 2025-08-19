package kernel

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrPhoneEmpty   = errors.New("phone number is empty")
	ErrPhoneInvalid = errors.New("phone number is invalid")
)

// Простое регулярное выражение для международного формата (минимум 7 цифр после +)
var phoneRe = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)

type Phone struct {
	value string
}

func NewPhone(s string) (Phone, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Phone{}, ErrPhoneEmpty
	}
	if !phoneRe.MatchString(s) {
		return Phone{}, fmt.Errorf("%w: %q (expected format: +1234567890)", ErrPhoneInvalid, s)
	}
	return Phone{value: s}, nil
}

func (p Phone) String() string { return p.value }

func (p Phone) Equals(other Phone) bool { return p.value == other.value }
