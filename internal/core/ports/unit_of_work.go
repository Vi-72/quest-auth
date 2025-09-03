package ports

import (
	"context"
)

type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback() error
	Execute(ctx context.Context, fn func() error) error
	UserRepository() UserRepository
}
