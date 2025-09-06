package assertions

import (
	"github.com/stretchr/testify/assert"

	"quest-auth/internal/core/application/usecases/commands"
)

type HandlerAssertions struct{ assert *assert.Assertions }

func NewHandlerAssertions(a *assert.Assertions) *HandlerAssertions {
	return &HandlerAssertions{assert: a}
}

func (a *HandlerAssertions) VerifyRegisterResult(res commands.RegisterUserResult, err error) {
	a.assert.NoError(err)
	a.assert.NotEmpty(res.User.ID.String())
	a.assert.NotEmpty(res.User.Email)
	a.assert.NotEmpty(res.User.Name)
	a.assert.NotEmpty(res.AccessToken)
	a.assert.NotEmpty(res.RefreshToken)
	a.assert.Equal("Bearer", res.TokenType)
	a.assert.Greater(res.ExpiresIn, 0)
}

func (a *HandlerAssertions) VerifyLoginResult(res commands.LoginUserResult, err error) {
	a.assert.NoError(err)
	a.assert.NotEmpty(res.User.ID.String())
	a.assert.NotEmpty(res.User.Email)
	a.assert.NotEmpty(res.User.Name)
	a.assert.NotEmpty(res.AccessToken)
	a.assert.NotEmpty(res.RefreshToken)
	a.assert.Equal("Bearer", res.TokenType)
	a.assert.Greater(res.ExpiresIn, 0)
}
