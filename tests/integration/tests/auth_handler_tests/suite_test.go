package auth_handler_tests

import (
	"github.com/Vi-72/quest-auth/tests/integration/tests"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	tests.DefaultSuite
}

func TestAuthHandlerOperations(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	s.DefaultSuite = tests.NewDefault(s)
	s.DefaultSuite.SetupSuite()
}

func (s *Suite) SetupTest() {
	s.DefaultSuite.SetupTest()
}

func (s *Suite) TearDownTest() {
	s.DefaultSuite.TearDownTest()
}

func (s *Suite) TearDownSuite() {
	s.DefaultSuite.TearDownSuite()
}
