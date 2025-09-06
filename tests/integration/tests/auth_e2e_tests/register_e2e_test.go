package auth_e2e_tests

import (
	"context"

	"quest-auth/tests/integration/core/assertions"
	casesteps "quest-auth/tests/integration/core/case_steps"
	testdatagenerators "quest-auth/tests/integration/core/test_data_generators"
)

// E2E 1: Create user via HTTP (register) and verify token fields
func (s *E2ESuite) TestCreateUserThroughAPI() {
	ctx := context.Background()
	httpAsserts := assertions.NewAuthHTTPAssertions(s.Assert())
	fieldAsserts := assertions.NewUserFieldAssertions(s.Assert())
	tokendAsserts := assertions.NewAssignAssertions(s.Assert())

	// Prepare request
	reqBody := testdatagenerators.DefaultUserData().ToRegisterHTTPRequest()

	// Execute HTTP request
	req := casesteps.RegisterHTTPRequest(reqBody)
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)

	// Assert HTTP response and parse
	reg := httpAsserts.RegisterHTTPCreatedSuccessfully(resp, err)
	fieldAsserts.VerifyHTTPResponseMatchesRegister(&reg, reqBody["email"].(string), reqBody["name"].(string), (*string)(nil))
	tokendAsserts.VerifyTokensPresent(reg.TokenType, reg.AccessToken, reg.RefreshToken, reg.ExpiresIn)

	// DB: verify user exists
	userID := reg.User.Id
	found, err := s.TestDIContainer.UserRepository.GetByID(userID)
	s.Require().NoError(err)
	s.Equal(userID, found.ID())

	// Event row exists for registration with expected name
	events, err := s.TestDIContainer.EventStorage.GetEventsByType(context.Background(), "user.registered")
	s.Require().NoError(err)
	s.Assert().GreaterOrEqual(len(events), 1)
}

// добавить тест на проверку валидации с самого нижнего лсоя доменной модели
