package ports

import "context"

// Repositories groups repositories available within a transactional boundary.
type Repositories struct {
	User  UserRepository
	Event EventPublisher
}

// TransactionManager defines transactional coordination for use cases.
type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context, repos Repositories) error) error
}
