package errs

import "google.golang.org/grpc/codes"

// GRPCError represents an error that can be converted to gRPC code.
type GRPCError interface {
	error
	GRPCCode() codes.Code
}

func (e *DomainValidationError) GRPCCode() codes.Code { return codes.InvalidArgument }

func (e *NotFoundError) GRPCCode() codes.Code { return codes.NotFound }

func (e *InfrastructureError) GRPCCode() codes.Code { return codes.Internal }

func (e *JWTValidationError) GRPCCode() codes.Code { return codes.Unauthenticated }
