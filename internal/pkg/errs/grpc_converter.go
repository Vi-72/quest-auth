package errs

import "errors"

import "google.golang.org/grpc/codes"

// ToGRPC converts any error to a gRPC status code.
func ToGRPC(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	var gerr GRPCError
	if errors.As(err, &gerr) {
		return gerr.GRPCCode()
	}
	return codes.Internal
}
