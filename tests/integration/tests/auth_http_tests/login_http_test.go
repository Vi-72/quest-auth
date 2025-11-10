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

// OpenAPI Validation Tests for Login (table-driven)
func (s *Suite) TestLoginHTTP_OpenAPIValidation() {
	ctx := context.Background()

	testCases := []struct {
		name        string
		email       string
		password    string
		expectError string
	}{
		{
			name:        "email_too_short",
			email:       "a@b",
			password:    "password123",
			expectError: "validation",
		},
		{
			name:        "email_too_long",
			email:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@example.com",
			password:    "password123",
			expectError: "validation",
		},
		{
			name:        "email_invalid_pattern",
			email:       "notanemail",
			password:    "password123",
			expectError: "validation",
		},
		{
			name:        "email_with_spaces",
			email:       "user name@example.com",
			password:    "password123",
			expectError: "validation",
		},
		{
			name:        "password_empty",
			email:       "user@example.com",
			password:    "",
			expectError: "validation",
		},
		{
			name:        "password_too_long",
			email:       "user@example.com",
			password:    string(make([]byte, 129)), // 129 chars, max 128
			expectError: "validation",
		},
		{
			name:        "both_missing",
			email:       "",
			password:    "",
			expectError: "",
		},
		{
			name:        "missing_password",
			email:       "user@example.com",
			password:    "",
			expectError: "validation",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := casesteps.LoginHTTPRequest(map[string]any{
				"email":    tc.email,
				"password": tc.password,
			})
			resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
			assertions.NewAuthHTTPAssertions(s.Assert()).HTTPErrorResponse(resp, err, 400, tc.expectError)
		})
	}
}
