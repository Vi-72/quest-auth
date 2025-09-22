package auth_grpc_tests

import (
	"context"
	grpcin "quest-auth/internal/adapters/in/grpc"
	"quest-auth/internal/core/application/usecases/queries"
	casesteps "quest-auth/tests/integration/core/case_steps"
	testdatagenerators "quest-auth/tests/integration/core/test_data_generators"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPC: Authenticate via gRPC using access token (random user)
func (s *Suite) TestAuthenticateThroughGRPC_RandomUser() {
	ctx := context.Background()

	// 1) Register random user via use case to get token
	userData := testdatagenerators.RandomUserData()
	reg, err := casesteps.RegisterUserStepData(ctx, s.TestDIContainer.RegisterUserHandler, userData)
	s.Require().NoError(err)

	// 2) Build gRPC auth handler and call Authenticate (real gRPC server method)
	authByToken := queries.NewAuthenticateByTokenHandler(s.TestDIContainer.JWTService)
	handler := grpcin.NewAuthHandler(authByToken)
	resp, err := handler.Authenticate(ctx, &v1.AuthenticateRequest{JwtToken: reg.AccessToken})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotNil(resp.User)

	// 3) Assert fields match claims
	s.Assert().Equal(reg.User.ID.String(), resp.User.Id)
	s.Assert().Equal(reg.User.Email, resp.User.Email)
	s.Assert().Equal(reg.User.Name, resp.User.Name)
	s.Assert().Equal(reg.User.Phone, resp.User.Phone)
}

// Validation: nil request should return InvalidArgument
func (s *Suite) TestAuthenticateThroughGRPC_NilRequest() {
	ctx := context.Background()
	authByToken := queries.NewAuthenticateByTokenHandler(s.TestDIContainer.JWTService)
	handler := grpcin.NewAuthHandler(authByToken)
	resp, err := handler.Authenticate(ctx, nil)
	s.Require().Error(err)
	s.Nil(resp)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Equal(codes.InvalidArgument, st.Code())
}

// Validation: empty jwt_token should return InvalidArgument
func (s *Suite) TestAuthenticateThroughGRPC_EmptyToken() {
	ctx := context.Background()
	authByToken := queries.NewAuthenticateByTokenHandler(s.TestDIContainer.JWTService)
	handler := grpcin.NewAuthHandler(authByToken)
	resp, err := handler.Authenticate(ctx, &v1.AuthenticateRequest{JwtToken: "   "})
	s.Require().Error(err)
	s.Nil(resp)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Equal(codes.InvalidArgument, st.Code())
}

// Domain-level: invalid token should surface as Unauthenticated at gRPC
func (s *Suite) TestAuthenticateThroughGRPC_InvalidToken_DomainError() {
	ctx := context.Background()
	authByToken := queries.NewAuthenticateByTokenHandler(s.TestDIContainer.JWTService)
	handler := grpcin.NewAuthHandler(authByToken)
	// malformed/invalid JWT (non-empty) to bypass handler empty-check and trigger lower-layer validation
	resp, err := handler.Authenticate(ctx, &v1.AuthenticateRequest{JwtToken: "invalid.jwt.token"})
	s.Require().Error(err)
	s.Nil(resp)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Equal(codes.Unauthenticated, st.Code())
}
