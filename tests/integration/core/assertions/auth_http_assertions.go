package assertions

import (
	"encoding/json"
	"net/http"
	"quest-auth/api/openapi"

	"github.com/stretchr/testify/assert"

	casesteps "quest-auth/tests/integration/core/case_steps"
)

type AuthHTTPAssertions struct {
	assert *assert.Assertions
}

func NewAuthHTTPAssertions(a *assert.Assertions) *AuthHTTPAssertions {
	return &AuthHTTPAssertions{assert: a}
}

// RegisterHTTPCreatedSuccessfully verifies Register HTTP 201 response and parses it
func (a *AuthHTTPAssertions) RegisterHTTPCreatedSuccessfully(resp *casesteps.HTTPResponse, err error) openapi.RegisterResponse {
	a.assert.NoError(err)
	a.assert.Equal(http.StatusCreated, resp.StatusCode)
	var r openapi.RegisterResponse
	a.assert.NoError(json.Unmarshal([]byte(resp.Body), &r))
	a.assert.NotEmpty(r.User.Id)
	a.assert.NotEmpty(r.AccessToken)
	a.assert.NotEmpty(r.RefreshToken)
	a.assert.Equal("Bearer", r.TokenType)
	a.assert.Greater(r.ExpiresIn, 0)
	return r
}

// LoginHTTPSuccess verifies Login HTTP 200 response and parses it
func (a *AuthHTTPAssertions) LoginHTTPSuccess(resp *casesteps.HTTPResponse, err error) openapi.LoginResponse {
	a.assert.NoError(err)
	a.assert.Equal(http.StatusOK, resp.StatusCode)
	var r openapi.LoginResponse
	a.assert.NoError(json.Unmarshal([]byte(resp.Body), &r))
	a.assert.NotEmpty(r.User.Id)
	a.assert.NotEmpty(r.AccessToken)
	a.assert.NotEmpty(r.RefreshToken)
	a.assert.Equal("Bearer", r.TokenType)
	a.assert.Greater(r.ExpiresIn, 0)
	return r
}

// HTTPErrorResponse asserts generic error response code and optional message substring
func (a *AuthHTTPAssertions) HTTPErrorResponse(resp *casesteps.HTTPResponse, err error, expectedStatus int, contains string) {
	a.assert.NoError(err)
	a.assert.Equal(expectedStatus, resp.StatusCode)
	if contains != "" {
		a.assert.Contains(resp.Body, contains)
	}
}
