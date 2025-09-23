package commands

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"
	"quest-auth/internal/core/ports"
	"quest-auth/internal/pkg/ddd"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRegisterUserHandler_RollsBackWhenPublishDomainEventsFails(t *testing.T) {
	ctx := context.Background()

	mainRepo := newFakeUserRepository()
	uow := newFakeUnitOfWork(mainRepo)

	publishErr := errors.New("publish failed")
	eventPublisher := &failingEventPublisher{err: publishErr}
	jwtService := &spyJWTService{}
	passwordHasher := fakePasswordHasher{}
	clock := fakeClock{now: time.Unix(0, 0)}

	handler := NewRegisterUserHandler(uow, eventPublisher, jwtService, passwordHasher, clock)

	cmd := RegisterUserCommand{
		Email:    "user@example.com",
		Phone:    "+12345678901",
		Name:     "Test User",
		Password: "password123",
	}

	_, err := handler.Handle(ctx, cmd)
	require.Error(t, err)
	require.ErrorIs(t, err, publishErr)

	require.False(t, mainRepo.HasUserByEmail("user@example.com"), "user must not be persisted when events publishing fails")
	require.Zero(t, jwtService.generateCalls, "jwt tokens must not be generated when transaction fails")
	require.Equal(t, 1, eventPublisher.callCount, "domain events must be attempted to publish once")
}

type fakeUnitOfWork struct {
	mainRepo *fakeUserRepository
	txRepo   *fakeUserRepository
	current  ports.UserRepository
	inTx     bool
}

func newFakeUnitOfWork(repo *fakeUserRepository) *fakeUnitOfWork {
	return &fakeUnitOfWork{mainRepo: repo, current: repo}
}

func (u *fakeUnitOfWork) Begin(ctx context.Context) error {
	if u.inTx {
		return fmt.Errorf("transaction already started")
	}
	u.txRepo = u.mainRepo.Clone()
	u.current = u.txRepo
	u.inTx = true
	return nil
}

func (u *fakeUnitOfWork) Commit(ctx context.Context) error {
	if !u.inTx {
		return fmt.Errorf("transaction not started")
	}
	u.mainRepo = u.txRepo.Clone()
	u.current = u.mainRepo
	u.txRepo = nil
	u.inTx = false
	return nil
}

func (u *fakeUnitOfWork) Rollback() error {
	if u.inTx {
		u.current = u.mainRepo
		u.txRepo = nil
		u.inTx = false
	}
	return nil
}

func (u *fakeUnitOfWork) Execute(ctx context.Context, fn func() error) error {
	if err := u.Begin(ctx); err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			_ = u.Rollback()
			panic(r)
		}
	}()
	if err := fn(); err != nil {
		_ = u.Rollback()
		return err
	}
	return u.Commit(ctx)
}

func (u *fakeUnitOfWork) UserRepository() ports.UserRepository {
	return u.current
}

type fakeUserRepository struct {
	mu      sync.RWMutex
	byID    map[uuid.UUID]*auth.User
	byEmail map[string]*auth.User
	byPhone map[string]*auth.User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		byID:    make(map[uuid.UUID]*auth.User),
		byEmail: make(map[string]*auth.User),
		byPhone: make(map[string]*auth.User),
	}
}

func (r *fakeUserRepository) Clone() *fakeUserRepository {
	r.mu.RLock()
	defer r.mu.RUnlock()
	clone := newFakeUserRepository()
	for id, user := range r.byID {
		copy := *user
		clone.byID[id] = &copy
		clone.byEmail[copy.Email.String()] = &copy
		clone.byPhone[copy.Phone.String()] = &copy
	}
	return clone
}

func (r *fakeUserRepository) HasUserByEmail(email string) bool {
	normalized, err := kernel.NewEmail(email)
	if err != nil {
		panic(err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.byEmail[normalized.String()]
	return ok
}

func (r *fakeUserRepository) Create(user *auth.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[user.ID()]; ok {
		return fmt.Errorf("user already exists")
	}
	if _, ok := r.byEmail[user.Email.String()]; ok {
		return fmt.Errorf("email already exists")
	}
	if _, ok := r.byPhone[user.Phone.String()]; ok {
		return fmt.Errorf("phone already exists")
	}
	copy := *user
	r.byID[copy.ID()] = &copy
	r.byEmail[copy.Email.String()] = &copy
	r.byPhone[copy.Phone.String()] = &copy
	return nil
}

func (r *fakeUserRepository) GetByID(id uuid.UUID) (*auth.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *fakeUserRepository) GetByEmail(email kernel.Email) (*auth.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *fakeUserRepository) GetByPhone(phone kernel.Phone) (*auth.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *fakeUserRepository) Update(user *auth.User) error {
	return fmt.Errorf("not implemented")
}

func (r *fakeUserRepository) Delete(id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (r *fakeUserRepository) EmailExists(email kernel.Email) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.byEmail[email.String()]
	return ok, nil
}

func (r *fakeUserRepository) PhoneExists(phone kernel.Phone) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.byPhone[phone.String()]
	return ok, nil
}

type failingEventPublisher struct {
	err       error
	callCount int
}

func (p *failingEventPublisher) PublishDomainEvents(ctx context.Context, events []ddd.DomainEvent) error {
	p.callCount++
	return p.err
}

type spyJWTService struct {
	generateCalls int
}

func (s *spyJWTService) GenerateTokenPair(userID uuid.UUID, email, name, phone string, createdAt time.Time) (*ports.TokenPair, error) {
	s.generateCalls++
	return &ports.TokenPair{}, nil
}

func (s *spyJWTService) ValidateAccessToken(token string) (*ports.TokenClaims, error) {
	panic("not implemented")
}

func (s *spyJWTService) RefreshTokens(refreshToken string) (*ports.TokenPair, error) {
	panic("not implemented")
}

type fakePasswordHasher struct{}

func (fakePasswordHasher) Hash(raw string) (string, error) {
	return "hashed:" + raw, nil
}

func (fakePasswordHasher) Compare(hash, raw string) bool {
	return hash == "hashed:"+raw
}

type fakeClock struct {
	now time.Time
}

func (c fakeClock) Now() time.Time {
	return c.now
}
