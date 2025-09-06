// HANDLER LAYER INTEGRATION TESTS
// Tests for Authenticate handler orchestration logic

package auth_handler_tests

import (
	"context"

	authpb "quest-auth/api/proto"
	grpcin "quest-auth/internal/adapters/in/grpc"
	"quest-auth/internal/core/application/usecases/commands"
	"quest-auth/internal/core/application/usecases/queries"
)

// Tests gRPC Authenticate handler via in-memory invocation
func (s *Suite) TestAuthenticateGRPC_Success() {
	ctx := context.Background()

	// 1) Register to obtain a valid access token
	regRes, err := s.TestDIContainer.RegisterUserHandler.Handle(ctx, commands.RegisterUserCommand{
		Email:    "grpcuser@example.com",
		Phone:    "+1234567870",
		Name:     "GRPC User",
		Password: "securepassword123",
	})
	s.Require().NoError(err)

	// 2) Build gRPC handler with same JWT service
	authByToken := queries.NewAuthenticateByTokenHandler(s.TestDIContainer.JWTService)
	handler := grpcin.NewAuthHandler(authByToken)

	// 3) Invoke Authenticate
	resp, err := handler.Authenticate(ctx, &authpb.AuthenticateRequest{JwtToken: regRes.AccessToken})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotNil(resp.User)

	s.Assert().Equal(regRes.User.ID.String(), resp.User.Id)
	s.Assert().Equal(regRes.User.Email, resp.User.Email)
	s.Assert().Equal(regRes.User.Name, resp.User.Name)
	s.Assert().Equal(regRes.User.Phone, resp.User.Phone)
	s.Require().NotNil(resp.User.CreatedAt)
}

//authenticate_grpc_test.go перепиши тесты в этой файле
//1) сделай тесты которые проверяют именно нендлер вызов его (вызывай через шаги)
//2) добавь тесты поторые проверяют валидацию именно котрые прописаны в хендлере

//добавь во все тесты такие стрки подсказки в начале
// REPOSITORY LAYER INTEGRATION TESTS
// Tests for repository implementations and database interactions
