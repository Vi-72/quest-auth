package auth_grpc_tests

import (
	"testing"

	tests "quest-auth/tests/integration/tests"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	tests.DefaultSuite
}

func (s *Suite) SetupSuite()    { s.DefaultSuite = tests.NewDefault(s); s.DefaultSuite.SetupSuite() }
func (s *Suite) SetupTest()     { s.DefaultSuite.SetupTest() }
func (s *Suite) TearDownTest()  { s.DefaultSuite.TearDownTest() }
func (s *Suite) TearDownSuite() { s.DefaultSuite.TearDownSuite() }

func TestAuthGRPCOperations(t *testing.T) { suite.Run(t, new(Suite)) }
