package auth_e2e_tests

import (
	"context"

	"github.com/Vi-72/quest-auth/tests/integration/core/assertions"
	casesteps "github.com/Vi-72/quest-auth/tests/integration/core/case_steps"
	testdatagenerators "github.com/Vi-72/quest-auth/tests/integration/core/test_data_generators"
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

// Negative E2E: Domain-level validation error during registration (short password)
func (s *E2ESuite) TestCreateUserThroughAPI_DomainValidation_ShortPassword() {
	ctx := context.Background()

	// Prepare invalid request (domain violation: password too short)
	invalid := testdatagenerators.DefaultUserData().ToRegisterHTTPRequest()
	invalid["password"] = "short"

	req := casesteps.RegisterHTTPRequest(invalid)
	resp, _ := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)

	// Should surface as HTTP 400 with problem details
	s.Assert().Equal(400, resp.StatusCode)
	s.Assert().Contains(resp.Body, "Bad Request")
}
