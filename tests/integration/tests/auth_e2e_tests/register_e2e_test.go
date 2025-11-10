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

// OpenAPI Validation Tests for Register (table-driven)
func (s *E2ESuite) TestRegisterHTTP_OpenAPIValidation() {
	ctx := context.Background()

	testCases := []struct {
		name        string
		modifyBody  func(body map[string]any)
		expectError string
	}{
		{
			name:        "email_too_short",
			modifyBody:  func(b map[string]any) { b["email"] = "a@b" },
			expectError: "validation",
		},
		{
			name: "email_too_long",
			modifyBody: func(b map[string]any) {
				b["email"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@example.com"
			},
			expectError: "validation",
		},
		{
			name:        "email_invalid_pattern",
			modifyBody:  func(b map[string]any) { b["email"] = "notanemail" },
			expectError: "validation",
		},
		{
			name:        "phone_too_short",
			modifyBody:  func(b map[string]any) { b["phone"] = "+123" },
			expectError: "validation",
		},
		{
			name:        "phone_invalid_pattern",
			modifyBody:  func(b map[string]any) { b["phone"] = "1234567890" },
			expectError: "validation",
		},
		{
			name:        "phone_starts_with_zero",
			modifyBody:  func(b map[string]any) { b["phone"] = "+0123456789" },
			expectError: "validation",
		},
		{
			name:        "name_empty",
			modifyBody:  func(b map[string]any) { b["name"] = "" },
			expectError: "validation",
		},
		{
			name:        "name_only_whitespace",
			modifyBody:  func(b map[string]any) { b["name"] = "   " },
			expectError: "validation",
		},
		{
			name:        "password_too_short",
			modifyBody:  func(b map[string]any) { b["password"] = "short12" },
			expectError: "validation",
		},
		{
			name:        "missing_email",
			modifyBody:  func(b map[string]any) { delete(b, "email") },
			expectError: "",
		},
		{
			name:        "missing_phone",
			modifyBody:  func(b map[string]any) { delete(b, "phone") },
			expectError: "",
		},
		{
			name:        "missing_password",
			modifyBody:  func(b map[string]any) { delete(b, "password") },
			expectError: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			body := testdatagenerators.DefaultUserData().ToRegisterHTTPRequest()
			tc.modifyBody(body)

			req := casesteps.RegisterHTTPRequest(body)
			resp, _ := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)

			s.Assert().Equal(400, resp.StatusCode, "Expected 400 status code")
			if tc.expectError != "" {
				s.Assert().Contains(resp.Body, tc.expectError)
			}
		})
	}
}
