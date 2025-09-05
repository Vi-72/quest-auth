package mocks

import (
	"quest-auth/internal/core/ports"
)

type ContractDIContainer struct {
	UserRepository ports.UserRepository
	UnitOfWork     ports.UnitOfWork
	EventPublisher ports.EventPublisher
}

func NewContractDIContainer() *ContractDIContainer {
	userRepo := NewMockUserRepository()
	uow := NewMockUnitOfWork(userRepo)
	return &ContractDIContainer{
		UserRepository: userRepo,
		UnitOfWork:     uow,
		EventPublisher: &ports.NullEventPublisher{},
	}
}

func (c *ContractDIContainer) CleanupAll() {
	if mr, ok := c.UserRepository.(*MockUserRepository); ok {
		mr.Clear()
	}
}
