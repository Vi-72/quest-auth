package assertions

import (
	"strings"

	"github.com/stretchr/testify/assert"

	"quest-auth/internal/generated/servers"
)

type UserFieldAssertions struct{ assert *assert.Assertions }

func NewUserFieldAssertions(a *assert.Assertions) *UserFieldAssertions {
	return &UserFieldAssertions{assert: a}
}

// VerifyHTTPResponseMatchesRegister verifies that register HTTP response matches request
func (a *UserFieldAssertions) VerifyHTTPResponseMatchesRegister(resp *servers.RegisterResponse, reqEmail, reqName string, reqPhone *string) {
	a.assert.Equal(strings.ToLower(reqEmail), strings.ToLower(resp.User.Email))
	a.assert.Equal(reqName, resp.User.Name)
	if reqPhone != nil {
		a.assert.NotNil(resp.User.Phone)
		a.assert.Equal(*reqPhone, *resp.User.Phone)
	}
}

// VerifyHTTPResponseMatchesLogin verifies that login HTTP response contains expected user fields
func (a *UserFieldAssertions) VerifyHTTPResponseMatchesLogin(resp *servers.LoginResponse, expectedEmail, expectedName string, expectedPhone *string) {
	a.assert.Equal(strings.ToLower(expectedEmail), strings.ToLower(resp.User.Email))
	a.assert.Equal(expectedName, resp.User.Name)
	if expectedPhone != nil {
		a.assert.NotNil(resp.User.Phone)
		a.assert.Equal(*expectedPhone, *resp.User.Phone)
	}
}
