package assertions

import (
	"github.com/stretchr/testify/assert"

	"quest-auth/internal/generated/servers"
)

type AssignAssertions struct{ assert *assert.Assertions }

func NewAssignAssertions(a *assert.Assertions) *AssignAssertions { return &AssignAssertions{assert: a} }

// VerifyTokensPresent asserts token fields are present in response
func (a *AssignAssertions) VerifyTokensPresent(tokenType string, access string, refresh string, expiresIn int) {
	a.assert.Equal("Bearer", tokenType)
	a.assert.NotEmpty(access)
	a.assert.NotEmpty(refresh)
	a.assert.Greater(expiresIn, 0)
}

// VerifyHTTPUser ensures HTTP user fields are filled
func (a *AssignAssertions) VerifyHTTPUser(u servers.User) {
	a.assert.NotEmpty(u.Id)
	a.assert.NotEmpty(u.Email)
	a.assert.NotEmpty(u.Name)
}
