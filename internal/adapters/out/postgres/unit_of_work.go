package postgres

import (
	"context"

	"quest-auth/internal/adapters/out/postgres/userrepo"
	"quest-auth/internal/core/ports"
	"quest-auth/internal/pkg/errs"

	"gorm.io/gorm"
)

var _ ports.UnitOfWork = &UnitOfWork{}

type UnitOfWork struct {
	tx             *gorm.DB
	db             *gorm.DB
	userRepository ports.UserRepository
}

func NewUnitOfWork(db *gorm.DB) (ports.UnitOfWork, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	uow := &UnitOfWork{db: db}

	// Создаем user repository
	userRepo := userrepo.NewRepository(uow.getDbInstance())
	uow.userRepository = userRepo

	return uow, nil
}

// getDbInstance возвращает активное соединение с БД (транзакция или основное)
func (u *UnitOfWork) getDbInstance() *gorm.DB {
	if u.tx != nil {
		return u.tx
	}
	return u.db
}

func (u *UnitOfWork) Tx() *gorm.DB {
	return u.tx
}

func (u *UnitOfWork) Db() *gorm.DB {
	return u.db
}

func (u *UnitOfWork) InTx() bool {
	return u.tx != nil
}

func (u *UnitOfWork) Begin(ctx context.Context) error {
	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	u.tx = tx

	// Обновляем repository для работы с транзакцией
	u.userRepository = userrepo.NewRepository(u.tx)

	return nil
}

func (u *UnitOfWork) Rollback() error {
	if u.tx != nil {
		err := u.tx.Rollback().Error
		u.tx = nil

		// Восстанавливаем repository для работы с основной БД
		u.userRepository = userrepo.NewRepository(u.db)

		return err
	}
	return nil
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}

	if err := u.tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}
	u.tx = nil

	// Восстанавливаем repository для работы с основной БД
	u.userRepository = userrepo.NewRepository(u.db)

	return nil
}

// Execute выполняет операцию в транзакции
func (u *UnitOfWork) Execute(fn func() error) error {
	ctx := context.Background()

	if err := u.Begin(ctx); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			u.Rollback()
			panic(r)
		}
	}()

	if err := fn(); err != nil {
		u.Rollback()
		return err
	}

	return u.Commit(ctx)
}

// Repository getters
func (u *UnitOfWork) UserRepository() ports.UserRepository {
	return u.userRepository
}
