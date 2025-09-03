package unit

import (
	"context"
	"testing"

	authpb "quest-auth/api/proto"
	grpcHandler "quest-auth/internal/adapters/in/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthHandler_Authenticate_ValidationErrors(t *testing.T) {
	// Создаем handler с nil зависимостями для тестирования только валидации
	handler := &grpcHandler.AuthHandler{}

	t.Run("nil request returns InvalidArgument", func(t *testing.T) {
		ctx := context.Background()

		resp, err := handler.Authenticate(ctx, nil)

		assert.Nil(t, resp)
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "request cannot be nil")
	})

	t.Run("empty JWT token returns InvalidArgument", func(t *testing.T) {
		ctx := context.Background()
		req := &authpb.AuthenticateRequest{
			JwtToken: "",
		}

		resp, err := handler.Authenticate(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "jwt_token is required")
	})

	t.Run("whitespace JWT token returns InvalidArgument", func(t *testing.T) {
		ctx := context.Background()
		req := &authpb.AuthenticateRequest{
			JwtToken: "   ",
		}

		resp, err := handler.Authenticate(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "jwt_token is required")
	})
}
