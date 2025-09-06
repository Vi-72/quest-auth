package testdatagenerators

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"quest-auth/internal/core/application/usecases/commands"
	"quest-auth/internal/generated/servers"

	openapi_types "github.com/oapi-codegen/runtime/types"
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

func (u UserTestData) ToRegisterRequest() servers.RegisterRequest {
	return servers.RegisterRequest{
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

func (u UserTestData) ToLoginRequest() servers.LoginRequest {
	return servers.LoginRequest{
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

func WithRand(rng *rand.Rand) Option {
	return func(_ *UserTestData, _ *rand.Rand) { defaultRng = rng }
}

func WithEmail(email string) Option { return func(u *UserTestData, _ *rand.Rand) { u.Email = email } }
func WithPhone(phone string) Option { return func(u *UserTestData, _ *rand.Rand) { u.Phone = phone } }
func WithName(name string) Option   { return func(u *UserTestData, _ *rand.Rand) { u.Name = name } }
func WithPassword(p string) Option  { return func(u *UserTestData, _ *rand.Rand) { u.Password = p } }

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
