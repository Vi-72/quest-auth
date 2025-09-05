package mocks

import (
	"fmt"
	"sync"

	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"

	"github.com/google/uuid"
)

type MockUserRepository struct {
	mu      sync.RWMutex
	byID    map[uuid.UUID]*auth.User
	byEmail map[string]*auth.User
	byPhone map[string]*auth.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		byID:    make(map[uuid.UUID]*auth.User),
		byEmail: make(map[string]*auth.User),
		byPhone: make(map[string]*auth.User),
	}
}

func (m *MockUserRepository) Create(user *auth.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.byID[user.ID()]; ok {
		return fmt.Errorf("user already exists")
	}
	if _, ok := m.byEmail[user.Email.String()]; ok {
		return fmt.Errorf("email already exists")
	}
	if _, ok := m.byPhone[user.Phone.String()]; ok {
		return fmt.Errorf("phone already exists")
	}
	u := *user
	m.byID[u.ID()] = &u
	m.byEmail[u.Email.String()] = &u
	m.byPhone[u.Phone.String()] = &u
	return nil
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*auth.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.byID[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return u, nil
}

func (m *MockUserRepository) GetByEmail(email kernel.Email) (*auth.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.byEmail[email.String()]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return u, nil
}

func (m *MockUserRepository) GetByPhone(phone kernel.Phone) (*auth.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.byPhone[phone.String()]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return u, nil
}

func (m *MockUserRepository) Update(user *auth.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.byID[user.ID()]; !ok {
		return fmt.Errorf("not found")
	}
	u := *user
	m.byID[u.ID()] = &u
	m.byEmail[u.Email.String()] = &u
	m.byPhone[u.Phone.String()] = &u
	return nil
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.byID[id]
	if !ok {
		return fmt.Errorf("not found")
	}
	delete(m.byID, id)
	delete(m.byEmail, u.Email.String())
	delete(m.byPhone, u.Phone.String())
	return nil
}

func (m *MockUserRepository) EmailExists(email kernel.Email) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.byEmail[email.String()]
	return ok, nil
}

func (m *MockUserRepository) PhoneExists(phone kernel.Phone) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.byPhone[phone.String()]
	return ok, nil
}

// Helpers
func (m *MockUserRepository) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.byID = make(map[uuid.UUID]*auth.User)
	m.byEmail = make(map[string]*auth.User)
	m.byPhone = make(map[string]*auth.User)
}
