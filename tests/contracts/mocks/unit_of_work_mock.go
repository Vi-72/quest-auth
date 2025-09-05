package mocks

import (
	"context"

	"quest-auth/internal/core/ports"
)

type MockUnitOfWork struct {
	userRepo ports.UserRepository
	inTx     bool
}

func NewMockUnitOfWork(userRepo ports.UserRepository) *MockUnitOfWork {
	return &MockUnitOfWork{userRepo: userRepo}
}

func (m *MockUnitOfWork) Begin(ctx context.Context) error  { m.inTx = true; return nil }
func (m *MockUnitOfWork) Commit(ctx context.Context) error { m.inTx = false; return nil }
func (m *MockUnitOfWork) Rollback() error                  { m.inTx = false; return nil }
func (m *MockUnitOfWork) Execute(ctx context.Context, fn func() error) error {
	if err := m.Begin(ctx); err != nil {
		return err
	}
	defer m.Rollback()
	if err := fn(); err != nil {
		return err
	}
	return m.Commit(ctx)
}
func (m *MockUnitOfWork) UserRepository() ports.UserRepository { return m.userRepo }
