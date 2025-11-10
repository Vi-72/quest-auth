package postgres

import (
	"context"

	"github.com/Vi-72/quest-auth/internal/adapters/out/postgres/eventrepo"
	"github.com/Vi-72/quest-auth/internal/adapters/out/postgres/userrepo"
	"github.com/Vi-72/quest-auth/internal/core/ports"

	"gorm.io/gorm"
)

// TransactionManager coordinates transactional work within the auth service.
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager backed by the given DB connection.
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// RunInTransaction executes the provided function within a database transaction.
func (tm *TransactionManager) RunInTransaction(
	ctx context.Context,
	fn func(ctx context.Context, repos ports.Repositories) error,
) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repos := ports.Repositories{
			User:  userrepo.NewRepository(tx),
			Event: eventrepo.NewRepository(tx),
		}
		return fn(ctx, repos)
	})
}
