package ports

// PasswordHasher provides methods for hashing and comparing passwords.
type PasswordHasher interface {
	// Hash generates a hash for the given raw password.
	Hash(raw string) (string, error)
	// Compare checks whether the provided raw password matches the given hash.
	Compare(hash, raw string) bool
}
