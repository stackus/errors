package errors_test

import (
	"net/http"
	"testing"

	"github.com/stackus/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestKindIdentity(t *testing.T) {
	a := errors.NewKind("FOO_NOT_FOUND", errors.ErrNotFound, "foo was not found")
	b := errors.NewKind("BAR_NOT_FOUND", errors.ErrNotFound, "bar was not found")

	require.True(t, errors.Is(a, a))
	require.False(t, errors.Is(a, b))

	require.True(t, errors.Is(a, errors.ErrNotFound))

	require.Equal(t, "FOO_NOT_FOUND", errors.TypeCode(a))
	require.Equal(t, http.StatusNotFound, errors.HTTPCode(a))
	require.Equal(t, codes.NotFound, errors.GRPCCode(a))
}
