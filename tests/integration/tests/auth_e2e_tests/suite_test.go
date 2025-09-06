package auth_e2e_tests

import (
	"testing"

	tests "quest-auth/tests/integration/tests"

	"github.com/stretchr/testify/suite"
)

type E2ESuite struct {
	suite.Suite
	tests.DefaultSuite
}

func (s *E2ESuite) SetupSuite() {
	s.DefaultSuite = tests.NewDefault(s)
	s.DefaultSuite.SetupSuite()
}

func (s *E2ESuite) SetupTest() {
	s.DefaultSuite.SetupTest()
}

func (s *E2ESuite) TearDownTest() {
	s.DefaultSuite.TearDownTest()
}

func (s *E2ESuite) TearDownSuite() {
	s.DefaultSuite.TearDownSuite()
}

func TestE2EAuthOperations(t *testing.T) {
	suite.Run(t, new(E2ESuite))
}
