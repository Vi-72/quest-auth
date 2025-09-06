package assertions

import (
	"github.com/stretchr/testify/assert"
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
