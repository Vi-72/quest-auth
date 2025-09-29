package contracts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Vi-72/quest-auth/internal/core/ports"
	"github.com/Vi-72/quest-auth/tests/contracts/mocks"
)

type UnitOfWorkContractSuite struct {
	suite.Suite
	uow ports.UnitOfWork
	ctx context.Context
}

func (s *UnitOfWorkContractSuite) SetupSuite() {
	s.uow = mocks.NewMockUnitOfWork(mocks.NewMockUserRepository())
	s.ctx = context.Background()
}

func TestUnitOfWorkContract(t *testing.T) { suite.Run(t, new(UnitOfWorkContractSuite)) }

func (s *UnitOfWorkContractSuite) TestBeginCommitRollback() {
	s.Require().NoError(s.uow.Begin(s.ctx))
	s.Require().NoError(s.uow.Commit(s.ctx))
	s.Require().NoError(s.uow.Begin(s.ctx))
	s.Require().NoError(s.uow.Rollback())
}

func (s *UnitOfWorkContractSuite) TestExecute() {
	err := s.uow.Execute(s.ctx, func() error { return nil })
	s.Require().NoError(err)
}

func (s *UnitOfWorkContractSuite) TestUserRepositoryAccessible() {
	repo := s.uow.UserRepository()
	s.Assert().NotNil(repo)
}
