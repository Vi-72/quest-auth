package testdatagenerators

import (
	"fmt"
	"math/rand"
	"os"
	"quest-auth/api/openapi"
	"strconv"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"quest-auth/internal/core/application/usecases/commands"
)

// ============================
// RNG (seed from ENV or default)
// ============================

var defaultRng *rand.Rand

func init() {
	seed := time.Now().UnixNano()
	if s, ok := os.LookupEnv("USER_GENERATOR_SEED"); ok {
		if parsed, err := strconv.ParseInt(s, 10, 64); err == nil {
			seed = parsed
		}
	}
	defaultRng = rand.New(rand.NewSource(seed))
}

// ============================
// Helpers
// ============================

func pick[T any](r *rand.Rand, xs []T) T { return xs[r.Intn(len(xs))] }

// ============================
// Test data model (User)
// ============================

type UserTestData struct {
	Email    string
	Phone    string
	Name     string
	Password string
}

// ============================
// Converters
// ============================

func (u UserTestData) ToRegisterCommand() commands.RegisterUserCommand {
	return commands.RegisterUserCommand{
		Email:    u.Email,
		Phone:    u.Phone,
		Name:     u.Name,
		Password: u.Password,
	}
}

func (u UserTestData) ToRegisterHTTPRequest() map[string]any {
	return map[string]any{
		"email":    u.Email,
		"phone":    u.Phone,
		"name":     u.Name,
		"password": u.Password,
	}
}

func (u UserTestData) ToRegisterRequest() openapi.RegisterRequest {
	return openapi.RegisterRequest{
		Email:    openapi_types.Email(u.Email),
		Phone:    u.Phone,
		Name:     u.Name,
		Password: u.Password,
	}
}

func (u UserTestData) ToLoginHTTPRequest() map[string]any {
	return map[string]any{
		"email":    u.Email,
		"password": u.Password,
	}
}

func (u UserTestData) ToLoginRequest() openapi.LoginRequest {
	return openapi.LoginRequest{
		Email:    openapi_types.Email(u.Email),
		Password: u.Password,
	}
}

// ============================
// Option pattern
// ============================

type Option func(*UserTestData, *rand.Rand)

func NewUser(opts ...Option) UserTestData {
	data := UserTestData{
		Email:    fmt.Sprintf("testuser_%d@example.com", defaultRng.Int63()),
		Phone:    "+1234567",
		Name:     "Test User",
		Password: "securepassword123",
	}
	for _, opt := range opts {
		opt(&data, defaultRng)
	}
	return data
}

func WithRandom() Option {
	names := []string{"Alice", "Bob", "Charlie", "Diana", "Eve", "Mallory", "Trent"}
	domains := []string{"example.com", "mail.com", "test.org"}
	phones := []string{"+1234567890", "+1234567891", "+1234567892", "+19991234567"}
	return func(u *UserTestData, r *rand.Rand) {
		u.Name = pick(r, names) + fmt.Sprintf(" %d", r.Intn(1000))
		u.Email = fmt.Sprintf("%s_%d@%s", pick(r, names), r.Int63(), pick(r, domains))
		u.Phone = pick(r, phones)
		// keep password constant but valid length
		u.Password = "securepassword123"
	}
}

// ============================
// Backward-compatible helpers
// ============================

func DefaultUserData() UserTestData { return NewUser() }
func RandomUserData() UserTestData  { return NewUser(WithRandom()) }
