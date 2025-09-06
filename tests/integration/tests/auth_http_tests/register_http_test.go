// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateRegisterRequest function

package auth_http_tests

import (
	"context"

	"quest-auth/tests/integration/core/assertions"
	casesteps "quest-auth/tests/integration/core/case_steps"
	testdatagenerators "quest-auth/tests/integration/core/test_data_generators"
)

// HTTP Register happy-path and validations
func (s *Suite) TestRegisterHTTP_Success() {
	ctx := context.Background()
	httpAsserts := assertions.NewAuthHTTPAssertions(s.Assert())
	fieldAsserts := assertions.NewUserFieldAssertions(s.Assert())
	tokenAsserts := assertions.NewAssignAssertions(s.Assert())

	body := testdatagenerators.RandomUserData().ToRegisterHTTPRequest()
	req := casesteps.RegisterHTTPRequest(body)
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)

	reg := httpAsserts.RegisterHTTPCreatedSuccessfully(resp, err)
	fieldAsserts.VerifyHTTPResponseMatchesRegister(&reg, body["email"].(string), body["name"].(string), nil)
	tokenAsserts.VerifyTokensPresent(reg.TokenType, reg.AccessToken, reg.RefreshToken, reg.ExpiresIn)
}

func (s *Suite) TestRegisterHTTP_Validation_EmptyBody() {
	ctx := context.Background()
	req := casesteps.RegisterHTTPRequest(map[string]any{})
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
	s.Assert().Equal(400, resp.StatusCode)
	s.Assert().NoError(err)
}

func (s *Suite) TestRegisterHTTP_Validation_InvalidEmail() {
	ctx := context.Background()
	body := testdatagenerators.DefaultUserData().ToRegisterHTTPRequest()
	body["email"] = "invalid-email"
	req := casesteps.RegisterHTTPRequest(body)
	resp, _ := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
	s.Assert().Equal(400, resp.StatusCode)
}

func (s *Suite) TestRegisterHTTP_Validation_EmptyFields() {
	ctx := context.Background()
	cases := []map[string]any{
		{"email": "", "phone": "+1234567890", "name": "A", "password": "securepassword123"},
		{"email": "user@example.com", "phone": "", "name": "A", "password": "securepassword123"},
		{"email": "user@example.com", "phone": "+1234567890", "name": "", "password": "securepassword123"},
		{"email": "user@example.com", "phone": "+1234567890", "name": "A", "password": ""},
	}
	for _, b := range cases {
		req := casesteps.RegisterHTTPRequest(b)
		resp, _ := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
		s.Assert().Equal(400, resp.StatusCode)
	}
}
