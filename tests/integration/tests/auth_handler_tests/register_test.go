// HANDLER LAYER INTEGRATION TESTS
// Tests for RegisterUserHandler orchestration logic (no HTTP)

package auth_handler_tests

import (
	"context"
	"strings"

	testdatagenerators "quest-auth/tests/integration/core/test_data_generators"
)

func (s *Suite) TestRegisterHandler_Success() {
	ctx := context.Background()

	// Pre-condition: prepare random user data
	data := testdatagenerators.RandomUserData()
	// Act: call RegisterUserHandler directly
	res, err := s.TestDIContainer.RegisterUserHandler.Handle(ctx, data.ToRegisterCommand())
	// Assert
	s.Require().NoError(err)

	s.Assert().NotEmpty(res.User.ID)
	s.Assert().Equal(strings.ToLower(data.Email), res.User.Email)
	s.Assert().Equal(data.Name, res.User.Name)
	s.Assert().Equal(data.Phone, res.User.Phone)
	s.Assert().NotEmpty(res.AccessToken)
	s.Assert().NotEmpty(res.RefreshToken)
	s.Assert().Equal("Bearer", res.TokenType)
	s.Assert().Greater(int(res.ExpiresIn), 0)
}

func (s *Suite) TestRegisterHandler_Validation_InvalidEmail() {
	ctx := context.Background()
	data := testdatagenerators.RandomUserData()
	data.Email = "not-an-email"
	_, err := s.TestDIContainer.RegisterUserHandler.Handle(ctx, data.ToRegisterCommand())
	s.Require().Error(err)
}

func (s *Suite) TestRegisterHandler_Validation_InvalidPhone() {
	ctx := context.Background()
	data := testdatagenerators.RandomUserData()
	data.Phone = "12345" // invalid, missing + and too short
	_, err := s.TestDIContainer.RegisterUserHandler.Handle(ctx, data.ToRegisterCommand())
	s.Require().Error(err)
}

func (s *Suite) TestRegisterHandler_Validation_EmptyName() {
	ctx := context.Background()
	data := testdatagenerators.RandomUserData()
	data.Name = "  "
	_, err := s.TestDIContainer.RegisterUserHandler.Handle(ctx, data.ToRegisterCommand())
	s.Require().Error(err)
}

func (s *Suite) TestRegisterHandler_Validation_ShortPassword() {
	ctx := context.Background()
	data := testdatagenerators.RandomUserData()
	data.Password = "short"
	_, err := s.TestDIContainer.RegisterUserHandler.Handle(ctx, data.ToRegisterCommand())
	s.Require().Error(err)
}

func (s *Suite) TestRegisterHandler_Validation_EmailAlreadyExists() {
	ctx := context.Background()
	first := testdatagenerators.RandomUserData()
	_, err := s.TestDIContainer.RegisterUserHandler.Handle(ctx, first.ToRegisterCommand())
	s.Require().NoError(err)

	second := testdatagenerators.RandomUserData()
	second.Email = first.Email // duplicate email
	_, err = s.TestDIContainer.RegisterUserHandler.Handle(ctx, second.ToRegisterCommand())
	s.Require().Error(err)
}

func (s *Suite) TestRegisterHandler_Validation_PhoneAlreadyExists() {
	ctx := context.Background()
	first := testdatagenerators.RandomUserData()
	_, err := s.TestDIContainer.RegisterUserHandler.Handle(ctx, first.ToRegisterCommand())
	s.Require().NoError(err)

	second := testdatagenerators.RandomUserData()
	second.Phone = first.Phone // duplicate phone
	_, err = s.TestDIContainer.RegisterUserHandler.Handle(ctx, second.ToRegisterCommand())
	s.Require().Error(err)
}
