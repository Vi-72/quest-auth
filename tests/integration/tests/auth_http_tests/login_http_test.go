// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateLoginRequest function

package auth_http_tests

import (
	"context"

	"quest-auth/tests/integration/core/assertions"
	casesteps "quest-auth/tests/integration/core/case_steps"
	testdatagenerators "quest-auth/tests/integration/core/test_data_generators"
)

// HTTP Login happy-path and validations
func (s *Suite) TestLoginHTTP_Success() {
	ctx := context.Background()
	httpAsserts := assertions.NewAuthHTTPAssertions(s.Assert())
	tokenAsserts := assertions.NewAssignAssertions(s.Assert())

	// Ensure user exists via use case
	data := testdatagenerators.DefaultUserData()
	_, err := casesteps.RegisterUserStepData(ctx, s.TestDIContainer.RegisterUserHandler, data)
	s.Require().NoError(err)

	req := casesteps.LoginHTTPRequest(data.ToLoginHTTPRequest())
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)

	login := httpAsserts.LoginHTTPSuccess(resp, err)
	tokenAsserts.VerifyTokensPresent(login.TokenType, login.AccessToken, login.RefreshToken, login.ExpiresIn)
}

func (s *Suite) TestLoginHTTP_Validation_EmptyBody() {
	ctx := context.Background()
	req := casesteps.LoginHTTPRequest(map[string]any{})
	resp, _ := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
	s.Assert().Equal(400, resp.StatusCode)
}

func (s *Suite) TestLoginHTTP_Validation_InvalidEmailFormat() {
	ctx := context.Background()
	req := casesteps.LoginHTTPRequest(map[string]any{"email": "invalid", "password": "password"})
	resp, _ := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
	s.Assert().Equal(400, resp.StatusCode)
}

func (s *Suite) TestLoginHTTP_Validation_MissingPassword() {
	ctx := context.Background()
	req := casesteps.LoginHTTPRequest(map[string]any{"email": "user@example.com"})
	resp, _ := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
	s.Assert().Equal(400, resp.StatusCode)
}
