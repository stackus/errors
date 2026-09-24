package errors_test

import (
	"testing"

	"github.com/stackus/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCRoundTrip(t *testing.T) {
	kind := errors.NewKind(
		"USER_NOT_FOUND",
		errors.ErrNotFound,
		"user not found",
		errors.WithPublicMessage("User not found"),
	)

	registry, err := errors.NewRegistry(kind)
	require.NoError(t, err)

	sent := errors.SendGRPCError(kind)
	received := errors.ReceiveGRPCError(sent, registry)

	require.True(t, errors.Is(received, kind))
	require.True(t, errors.Is(received, errors.ErrNotFound))
	require.Equal(t, "USER_NOT_FOUND", errors.TypeCode(received))
	require.Equal(t, codes.NotFound, errors.GRPCCode(received))

	s, ok := status.FromError(received)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
}
