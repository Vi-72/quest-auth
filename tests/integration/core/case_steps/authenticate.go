package casesteps

import (
	"context"
	authpb "github.com/Vi-72/quest-auth/api/grpc/sdk/go/proto/auth/v1"

	grpcin "github.com/Vi-72/quest-auth/internal/adapters/in/grpc"
	"github.com/Vi-72/quest-auth/internal/core/application/usecases/queries"
	"github.com/Vi-72/quest-auth/internal/core/ports"
)

// AuthenticateByTokenStep invokes the gRPC Authenticate handler using provided JWT service and token
func AuthenticateByTokenStep(ctx context.Context, jwtService ports.JWTService, token string) (*authpb.AuthenticateResponse, error) {
	authenticateByToken := queries.NewAuthenticateByTokenHandler(jwtService)
	handler := grpcin.NewAuthHandler(authenticateByToken)
	return handler.Authenticate(ctx, &authpb.AuthenticateRequest{JwtToken: token})
}
