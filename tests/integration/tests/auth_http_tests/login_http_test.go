// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateLoginRequest function

package auth_http_tests

import (
	"context"

	"github.com/Vi-72/quest-auth/tests/integration/core/assertions"
	casesteps "github.com/Vi-72/quest-auth/tests/integration/core/case_steps"
	testdatagenerators "github.com/Vi-72/quest-auth/tests/integration/core/test_data_generators"
)

// HTTP Login happy-path and validations
func (s *Suite) TestLoginHTTP_Success() {
	ctx := context.Background()
	httpAsserts := assertions.NewAuthHTTPAssertions(s.Assert())
	tokenAsserts := assertions.NewAssignAssertions(s.Assert())

	// Pre-condition: ensure user exists via use case
	data := testdatagenerators.DefaultUserData()
	_, err := casesteps.RegisterUserStepData(ctx, s.TestDIContainer.RegisterUserHandler, data)
	s.Require().NoError(err)

	// Act: perform HTTP login
	req := casesteps.LoginHTTPRequest(data.ToLoginHTTPRequest())
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)

	// Assert
	login := httpAsserts.LoginHTTPSuccess(resp, err)
	tokenAsserts.VerifyTokensPresent(login.TokenType, login.AccessToken, login.RefreshToken, login.ExpiresIn)
}

func (s *Suite) TestLoginHTTP_Validation_EmptyBody() {
	ctx := context.Background()
	// Pre-condition: empty body
	req := casesteps.LoginHTTPRequest(map[string]any{})
	// Act
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
	// Assert
	assertions.NewAuthHTTPAssertions(s.Assert()).HTTPErrorResponse(resp, err, 400, "")
}

func (s *Suite) TestLoginHTTP_Validation_InvalidEmailFormat() {
	ctx := context.Background()
	// Pre-condition: invalid email format
	req := casesteps.LoginHTTPRequest(map[string]any{"email": "invalid", "password": "password"})
	// Act
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
	// Assert
	assertions.NewAuthHTTPAssertions(s.Assert()).HTTPErrorResponse(resp, err, 400, "")
}

func (s *Suite) TestLoginHTTP_Validation_MissingPassword() {
	ctx := context.Background()
	// Pre-condition: missing password
	req := casesteps.LoginHTTPRequest(map[string]any{"email": "user@example.com"})
	// Act
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
	// Assert
	assertions.NewAuthHTTPAssertions(s.Assert()).HTTPErrorResponse(resp, err, 400, "")
}
