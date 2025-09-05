package bcryptadapter

import (
	gobcrypt "golang.org/x/crypto/bcrypt"
	"quest-auth/internal/core/ports"
)

// Hasher implements PasswordHasher using bcrypt.
type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

// Hash generates bcrypt hash for the given password.
func (h *Hasher) Hash(raw string) (string, error) {
	b, err := gobcrypt.GenerateFromPassword([]byte(raw), gobcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compare checks whether given hash matches raw password.
func (h *Hasher) Compare(hash, raw string) bool {
	return gobcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}

var _ ports.PasswordHasher = (*Hasher)(nil)
