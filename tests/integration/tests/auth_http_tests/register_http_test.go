// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateRegisterRequest function

package auth_http_tests

import (
	"context"

	"github.com/Vi-72/quest-auth/tests/integration/core/assertions"
	casesteps "github.com/Vi-72/quest-auth/tests/integration/core/case_steps"
	testdatagenerators "github.com/Vi-72/quest-auth/tests/integration/core/test_data_generators"
)

// HTTP Register happy-path and validations
func (s *Suite) TestRegisterHTTP_Success() {
	ctx := context.Background()
	httpAsserts := assertions.NewAuthHTTPAssertions(s.Assert())
	fieldAsserts := assertions.NewUserFieldAssertions(s.Assert())
	tokenAsserts := assertions.NewAssignAssertions(s.Assert())

	// Pre-condition: prepare valid body
	body := testdatagenerators.RandomUserData().ToRegisterHTTPRequest()
	// Act: perform HTTP request
	req := casesteps.RegisterHTTPRequest(body)
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)

	// Assert
	reg := httpAsserts.RegisterHTTPCreatedSuccessfully(resp, err)
	fieldAsserts.VerifyHTTPResponseMatchesRegister(&reg, body["email"].(string), body["name"].(string), nil)
	tokenAsserts.VerifyTokensPresent(reg.TokenType, reg.AccessToken, reg.RefreshToken, reg.ExpiresIn)
}

func (s *Suite) TestRegisterHTTP_Validation_EmptyBody() {
	ctx := context.Background()
	// Pre-condition: empty body
	req := casesteps.RegisterHTTPRequest(map[string]any{})
	// Act
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
	// Assert
	assertions.NewAuthHTTPAssertions(s.Assert()).HTTPErrorResponse(resp, err, 400, "")
}

func (s *Suite) TestRegisterHTTP_Validation_InvalidEmail() {
	ctx := context.Background()
	// Pre-condition: invalid email
	body := testdatagenerators.DefaultUserData().ToRegisterHTTPRequest()
	body["email"] = "invalid-email"
	// Act
	req := casesteps.RegisterHTTPRequest(body)
	resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
	// Assert
	assertions.NewAuthHTTPAssertions(s.Assert()).HTTPErrorResponse(resp, err, 400, "")
}

func (s *Suite) TestRegisterHTTP_Validation_EmptyFields() {
	ctx := context.Background()
	// Pre-condition: different empty field cases
	cases := []map[string]any{
		{"email": "", "phone": "+1234567890", "name": "A", "password": "securepassword123"},
		{"email": "user@example.com", "phone": "", "name": "A", "password": "securepassword123"},
		{"email": "user@example.com", "phone": "+1234567890", "name": "", "password": "securepassword123"},
		{"email": "user@example.com", "phone": "+1234567890", "name": "A", "password": ""},
	}
	for _, b := range cases {
		// Act
		req := casesteps.RegisterHTTPRequest(b)
		resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)
		// Assert
		assertions.NewAuthHTTPAssertions(s.Assert()).HTTPErrorResponse(resp, err, 400, "")
	}
}
