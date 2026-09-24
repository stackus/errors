package errors_test

import (
	"testing"

	"github.com/stackus/errors"
	"github.com/stretchr/testify/require"
)

func TestPublicMessage(t *testing.T) {
	secret := errors.ErrInternal.Wrap(errors.New("database path: /private/psst.db"), "load authentication data")

	require.NotContains(t, errors.PublicMessage(secret), "/private/psst.db")
	require.Equal(t, "Internal Server Error", errors.PublicMessage(secret))

	err := errors.WrapPublicMessage(errors.ErrForbidden.Msg("user lacks permission to edit project 123"), "Permission denied")

	require.Equal(t, "Permission denied", errors.PublicMessage(err))
	require.Contains(t, err.Error(), "project 123")
}
