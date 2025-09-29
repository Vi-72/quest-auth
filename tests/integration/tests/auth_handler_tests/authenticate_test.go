// HANDLER LAYER INTEGRATION TESTS
// Tests for Authenticate handler orchestration logic

package auth_handler_tests

import (
	"context"

	authv1 "github.com/Vi-72/quest-auth/api/grpc/sdk/go/proto/auth/v1"

	grpcin "github.com/Vi-72/quest-auth/internal/adapters/in/grpc"
	"github.com/Vi-72/quest-auth/internal/core/application/usecases/queries"
	casesteps "github.com/Vi-72/quest-auth/tests/integration/core/case_steps"
	testdatagenerators "github.com/Vi-72/quest-auth/tests/integration/core/test_data_generators"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 1) Success: call handler via steps (AuthenticateByTokenStep)
func (s *Suite) TestAuthenticateHandler_Success_UsingSteps() {
	ctx := context.Background()

	// Pre-condition: register random user via use case to get token
	data := testdatagenerators.RandomUserData()
	reg, err := casesteps.RegisterUserStepData(ctx, s.TestDIContainer.RegisterUserHandler, data)
	s.Require().NoError(err)

	// Act: call Authenticate handler via case step
	resp, err := casesteps.AuthenticateByTokenStep(ctx, s.TestDIContainer.JWTService, reg.AccessToken)
	// Assert
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotNil(resp.User)

	s.Assert().Equal(reg.User.ID.String(), resp.User.Id)
	s.Assert().Equal(reg.User.Email, resp.User.Email)
	s.Assert().Equal(reg.User.Name, resp.User.Name)
	s.Assert().Equal(reg.User.Phone, resp.User.Phone)
}

// 2) Validation: req == nil -> InvalidArgument
func (s *Suite) TestAuthenticateHandler_Validation_NilRequest() {
	ctx := context.Background()
	// Pre-condition: build handler
	authByToken := queries.NewAuthenticateByTokenHandler(s.TestDIContainer.JWTService)
	handler := grpcin.NewAuthHandler(authByToken)
	// Act
	resp, err := handler.Authenticate(ctx, nil)
	// Assert
	s.Require().Error(err)
	s.Nil(resp)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Equal(codes.InvalidArgument, st.Code())
}

// 2) Validation: empty jwt_token -> InvalidArgument
func (s *Suite) TestAuthenticateHandler_Validation_EmptyToken() {
	ctx := context.Background()
	// Pre-condition: build handler
	authByToken := queries.NewAuthenticateByTokenHandler(s.TestDIContainer.JWTService)
	handler := grpcin.NewAuthHandler(authByToken)
	// Act
	resp, err := handler.Authenticate(ctx, &authv1.AuthenticateRequest{JwtToken: "   "})
	// Assert
	s.Require().Error(err)
	s.Nil(resp)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Equal(codes.InvalidArgument, st.Code())
}
