package contracts

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"
	"quest-auth/internal/core/ports"
	"quest-auth/tests/contracts/mocks"
)

type UserRepositoryContractSuite struct {
	suite.Suite
	repo ports.UserRepository
}

func (s *UserRepositoryContractSuite) SetupSuite() {
	s.repo = mocks.NewMockUserRepository()
}

func (s *UserRepositoryContractSuite) SetupTest() {
	if mr, ok := s.repo.(*mocks.MockUserRepository); ok {
		mr.Clear()
	}
}

func TestUserRepositoryContract(t *testing.T) { suite.Run(t, new(UserRepositoryContractSuite)) }

func (s *UserRepositoryContractSuite) TestCreateAndGetByID() {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", testHasher{}, testClock{})

	err := s.repo.Create(&u)
	s.Require().NoError(err)

	got, err := s.repo.GetByID(u.ID())
	s.Require().NoError(err)
	s.Equal(u.ID(), got.ID())
}

func (s *UserRepositoryContractSuite) TestGetByEmailAndPhone() {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", testHasher{}, testClock{})
	_ = s.repo.Create(&u)

	ge, err := s.repo.GetByEmail(email)
	s.Require().NoError(err)
	s.Equal(u.ID(), ge.ID())

	gp, err := s.repo.GetByPhone(phone)
	s.Require().NoError(err)
	s.Equal(u.ID(), gp.ID())
}

func (s *UserRepositoryContractSuite) TestUpdateAndDelete() {
	email, _ := kernel.NewEmail("user@example.com")
	phone, _ := kernel.NewPhone("+1234567890")
	u, _ := auth.NewUser(email, phone, "John Doe", "password123", testHasher{}, testClock{})
	_ = s.repo.Create(&u)

	// Change name and update
	_ = u.ChangeName("Jane", testClock{})
	err := s.repo.Update(&u)
	s.Require().NoError(err)

	// Delete
	err = s.repo.Delete(u.ID())
	s.Require().NoError(err)

	_, err = s.repo.GetByID(uuid.New())
	s.Error(err)
}

type testHasher struct{}

func (testHasher) Hash(raw string) (string, error) { return "h:" + raw, nil }
func (testHasher) Compare(hash, raw string) bool   { return hash == "h:"+raw }

type testClock struct{}

func (testClock) Now() time.Time { return time.Now() }
