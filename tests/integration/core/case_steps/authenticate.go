package casesteps

import (
	"context"
	authpb "quest-auth/api/grpc/proto"

	grpcin "quest-auth/internal/adapters/in/grpc"
	"quest-auth/internal/core/application/usecases/queries"
	"quest-auth/internal/core/ports"
)

// AuthenticateByTokenStep invokes the gRPC Authenticate handler using provided JWT service and token
func AuthenticateByTokenStep(ctx context.Context, jwtService ports.JWTService, token string) (*authpb.AuthenticateResponse, error) {
	authenticateByToken := queries.NewAuthenticateByTokenHandler(jwtService)
	handler := grpcin.NewAuthHandler(authenticateByToken)
	return handler.Authenticate(ctx, &authpb.AuthenticateRequest{JwtToken: token})
}
