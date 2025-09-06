// REPOSITORY LAYER INTEGRATION TESTS
// Tests for repository implementations and database interactions

//go:build integration

package repository

import (
	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"
	domainhelpers "quest-auth/tests/domain"
)

func (s *Suite) TestUserRepository_Create_And_GetByID() {
	// Pre-condition: build user aggregate
	email, _ := kernel.NewEmail("user.repo1@example.com")
	phone, _ := kernel.NewPhone("+1234567890")

	hasher := domainhelpers.NewMockPasswordHasher()
	clock := domainhelpers.NewMockClock()

	u, err := auth.NewUser(email, phone, "Repo User", "securepassword123", hasher, clock)
	s.Require().NoError(err)

	// Act: persist user
	err = s.TestDIContainer.UserRepository.Create(&u)
	s.Require().NoError(err)

	// Assert: fetched user matches
	found, err := s.TestDIContainer.UserRepository.GetByID(u.ID())
	s.Require().NoError(err)
	s.Equal(u.ID(), found.ID())
	s.Equal(u.Email.String(), found.Email.String())
	s.Equal(u.Phone.String(), found.Phone.String())
	s.Equal(u.Name, found.Name)
}

func (s *Suite) TestUserRepository_GetByEmail_And_Phone() {
	// Pre-condition: existing user
	email, _ := kernel.NewEmail("user.repo2@example.com")
	phone, _ := kernel.NewPhone("+1234567891")
	hasher := domainhelpers.NewMockPasswordHasher()
	clock := domainhelpers.NewMockClock()
	u, err := auth.NewUser(email, phone, "Repo User 2", "securepassword123", hasher, clock)
	s.Require().NoError(err)
	s.Require().NoError(s.TestDIContainer.UserRepository.Create(&u))

	// Act & Assert: get by email
	byEmail, err := s.TestDIContainer.UserRepository.GetByEmail(email)
	s.Require().NoError(err)
	s.Equal(u.ID(), byEmail.ID())

	// Act & Assert: get by phone
	byPhone, err := s.TestDIContainer.UserRepository.GetByPhone(phone)
	s.Require().NoError(err)
	s.Equal(u.ID(), byPhone.ID())
}

func (s *Suite) TestUserRepository_Update() {
	// Pre-condition: existing user
	email, _ := kernel.NewEmail("user.repo3@example.com")
	phone, _ := kernel.NewPhone("+1234567892")
	hasher := domainhelpers.NewMockPasswordHasher()
	clock := domainhelpers.NewMockClock()
	u, err := auth.NewUser(email, phone, "Repo User 3", "securepassword123", hasher, clock)
	s.Require().NoError(err)
	s.Require().NoError(s.TestDIContainer.UserRepository.Create(&u))

	// Act: update fields
	_ = u.ChangeName("Updated Name", clock)
	newPhone, _ := kernel.NewPhone("+19991234567")
	u.ChangePhone(newPhone, clock)

	s.Require().NoError(s.TestDIContainer.UserRepository.Update(&u))

	// Assert: persisted changes
	found, err := s.TestDIContainer.UserRepository.GetByID(u.ID())
	s.Require().NoError(err)
	s.Equal("Updated Name", found.Name)
	s.Equal(newPhone.String(), found.Phone.String())
}

func (s *Suite) TestUserRepository_Delete() {
	// Pre-condition: existing user
	email, _ := kernel.NewEmail("user.repo4@example.com")
	phone, _ := kernel.NewPhone("+1234567893")
	hasher := domainhelpers.NewMockPasswordHasher()
	clock := domainhelpers.NewMockClock()
	u, err := auth.NewUser(email, phone, "Repo User 4", "securepassword123", hasher, clock)
	s.Require().NoError(err)
	s.Require().NoError(s.TestDIContainer.UserRepository.Create(&u))

	// Act: delete
	s.Require().NoError(s.TestDIContainer.UserRepository.Delete(u.ID()))

	// Assert: not found
	_, err = s.TestDIContainer.UserRepository.GetByID(u.ID())
	s.Require().Error(err)
}

func (s *Suite) TestUserRepository_EmailAndPhoneExists() {
	// Pre-condition: existing user
	email, _ := kernel.NewEmail("user.repo5@example.com")
	phone, _ := kernel.NewPhone("+1234567894")
	hasher := domainhelpers.NewMockPasswordHasher()
	clock := domainhelpers.NewMockClock()
	u, err := auth.NewUser(email, phone, "Repo User 5", "securepassword123", hasher, clock)
	s.Require().NoError(err)
	s.Require().NoError(s.TestDIContainer.UserRepository.Create(&u))

	// Act & Assert: existence checks
	existsEmail, err := s.TestDIContainer.UserRepository.EmailExists(email)
	s.Require().NoError(err)
	s.True(existsEmail)

	existsPhone, err := s.TestDIContainer.UserRepository.PhoneExists(phone)
	s.Require().NoError(err)
	s.True(existsPhone)

	otherEmail, _ := kernel.NewEmail("other@example.com")
	otherPhone, _ := kernel.NewPhone("+15551234567")
	notEmail, _ := s.TestDIContainer.UserRepository.EmailExists(otherEmail)
	notPhone, _ := s.TestDIContainer.UserRepository.PhoneExists(otherPhone)
	s.False(notEmail)
	s.False(notPhone)
}
