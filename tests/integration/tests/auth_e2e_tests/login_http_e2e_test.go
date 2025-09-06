package auth_e2e_tests

import (
	"context"
	"net/http"

	casesteps "quest-auth/tests/integration/core/case_steps"
	testdatagenerators "quest-auth/tests/integration/core/test_data_generators"
)

// E2E 1: Login via HTTP and verify event stored
func (s *E2ESuite) TestLoginThroughAPI() {
	ctx := context.Background()

	// Arrange: create user via use case (lower layer)
	data := testdatagenerators.DefaultUserData()
	_, err := casesteps.RegisterUserStepData(ctx, s.TestDIContainer.RegisterUserHandler, data)
	s.Require().NoError(err)

	// Act: login via HTTP API using case_steps
	loginReq := casesteps.LoginHTTPRequest(data.ToLoginHTTPRequest())
	loginResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, loginReq)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, loginResp.StatusCode)

	// Verify login event exists using storage helper
	events, err := s.TestDIContainer.EventStorage.GetEventsByType(context.Background(), "user.login")
	s.Require().NoError(err)
	s.Assert().GreaterOrEqual(len(events), 1)
}

// E2E: Domain-level validation surfaces via HTTP (wrong password) using layered helpers
func (s *E2ESuite) TestLoginThroughAPI_WrongPassword_DomainValidation() {
	ctx := context.Background()

	// Arrange: create user via use case (lower layer)
	data := testdatagenerators.DefaultUserData()
	_, err := casesteps.RegisterUserStepData(ctx, s.TestDIContainer.RegisterUserHandler, data)
	s.Require().NoError(err)

	// Act: login via HTTP with wrong password
	loginReq := casesteps.LoginHTTPRequest(map[string]any{
		"email":    data.Email,
		"password": "wrong-password",
	})
	resp, _ := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, loginReq)

	// Assert: lowest domain validation error surfaced as 401 Unauthorized
	s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)
	s.Assert().Contains(resp.Body, "Unauthorized")
}
