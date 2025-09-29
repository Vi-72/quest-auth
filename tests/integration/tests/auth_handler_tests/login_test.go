// HANDLER LAYER INTEGRATION TESTS
// Tests for LoginUserHandler orchestration logic (no HTTP)

package auth_handler_tests

import (
	"context"
	"strings"

	casesteps "github.com/Vi-72/quest-auth/tests/integration/core/case_steps"
	testdatagenerators "github.com/Vi-72/quest-auth/tests/integration/core/test_data_generators"
)

func (s *Suite) TestLoginHandler_Success() {
	ctx := context.Background()

	// Pre-condition: register a user via RegisterUserHandler
	data := testdatagenerators.RandomUserData()
	regRes, err := casesteps.RegisterUserStepData(ctx, s.TestDIContainer.RegisterUserHandler, data)
	// Assert pre-condition
	s.Require().NoError(err)

	// Act: perform login via LoginUserHandler
	loginRes, err := casesteps.LoginUserStep(ctx, s.TestDIContainer.LoginUserHandler, data.Email, data.Password)

	// Assert
	s.Require().NoError(err)
	s.Assert().NotEmpty(loginRes.AccessToken)
	s.Assert().NotEmpty(loginRes.RefreshToken)
	s.Assert().Equal("Bearer", loginRes.TokenType)
	s.Assert().Greater(int(loginRes.ExpiresIn), 0)

	s.Assert().Equal(regRes.User.ID, loginRes.User.ID)
	s.Assert().Equal(strings.ToLower(data.Email), loginRes.User.Email)
	s.Assert().Equal(data.Name, loginRes.User.Name)
	s.Assert().Equal(data.Phone, loginRes.User.Phone)
}

func (s *Suite) TestLoginHandler_Unauthorized_InvalidCredentials() {
	ctx := context.Background()

	// Pre-condition: no user registered

	// Act: call login with wrong credentials
	_, err := casesteps.LoginUserStep(ctx, s.TestDIContainer.LoginUserHandler, "nouser@example.com", "wrongpass")

	// Assert: expect error
	s.Require().Error(err)
}

func (s *Suite) TestLoginHandler_Validation_InvalidEmail() {
	ctx := context.Background()
	// Pre-condition: no need to register for email validation
	_, err := casesteps.LoginUserStep(ctx, s.TestDIContainer.LoginUserHandler, "bad-email", "somepassword")
	// Assert
	s.Require().Error(err)
}

func (s *Suite) TestLoginHandler_Validation_WrongPassword() {
	ctx := context.Background()
	// Pre-condition: register valid user
	data := testdatagenerators.RandomUserData()
	_, err := casesteps.RegisterUserStepData(ctx, s.TestDIContainer.RegisterUserHandler, data)
	s.Require().NoError(err)

	// Act: login with wrong password
	_, err = casesteps.LoginUserStep(ctx, s.TestDIContainer.LoginUserHandler, data.Email, "wrongpass")
	// Assert
	s.Require().Error(err)
}
