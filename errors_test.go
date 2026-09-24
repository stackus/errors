package errors_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stackus/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func ExampleWrap() {
	err := errors.Wrap(errors.ErrNotFound, "message")
	fmt.Println(err)
	// Output: message
}

func ExampleWrap_multiple() {
	err := errors.Wrap(errors.ErrNotFound, "original message")
	err = errors.Wrap(err, "prefixed message")
	fmt.Println(err)
	// Output: prefixed message: original message
}

func TestWrapping(t *testing.T) {
	cause := errors.New("no rows found") // simulate sql.ErrNoRows
	err := errors.ErrNotFound.Wrap(cause, "find user")

	err = fmt.Errorf("repository: %w", err)
	err = errors.Wrap(err, "login")

	require.True(t, errors.Is(err, errors.ErrNotFound))
	require.Equal(t, "NOT_FOUND", errors.TypeCode(err))
	require.Equal(t, http.StatusNotFound, errors.HTTPCode(err))
}

func TestClassification(t *testing.T) {
	cause := errors.ErrNotFound.Msg("user not found")
	err := errors.ErrUnauthorized.Wrap(cause, "unauthorized access")

	require.True(t, errors.Is(err, errors.ErrNotFound))
	require.True(t, errors.Is(err, errors.ErrUnauthorized))
	require.Equal(t, "UNAUTHORIZED", errors.TypeCode(err))
	require.Equal(t, http.StatusUnauthorized, errors.HTTPCode(err))
	require.Equal(t, codes.Unauthenticated, errors.GRPCCode(err))
}

func TestJoinedErrors(t *testing.T) {
	same := errors.Join(
		errors.ErrNotFound,
		errors.ErrNotFound.Msg("missing record"),
	)

	require.Equal(t, http.StatusNotFound, errors.HTTPCode(same))

	mixed := errors.Join(
		errors.ErrNotFound,
		errors.ErrForbidden,
	)

	require.Equal(t, "UNKNOWN", errors.TypeCode(mixed))
	require.Equal(t, http.StatusInternalServerError, errors.HTTPCode(mixed))

	classified := errors.ErrInternal.Wrap(mixed, "multiple operations failed")

	require.Equal(t, "INTERNAL", errors.TypeCode(classified))
	require.True(t, errors.Is(classified, errors.ErrNotFound))
	require.True(t, errors.Is(classified, errors.ErrForbidden))
}
