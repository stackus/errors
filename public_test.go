package errors_test

import (
	"encoding/json"
	"fmt"
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

// PublicMessage never returns an error's internal text. Server errors become
// "Internal Server Error", client errors use their HTTP status text, and a
// kind can supply its own message.
func ExamplePublicMessage() {
	ErrAccountLocked := errors.NewKind(
		"ACCOUNT_LOCKED",
		errors.ErrPermissionDenied,
		"account locked after failed logins",
		errors.WithPublicMessage("Your account is locked"),
	)

	fmt.Println(errors.PublicMessage(errors.ErrInternal.Wrap(errors.New("password authentication failed for user app"), "connect")))
	fmt.Println(errors.PublicMessage(errors.ErrNotFound.Msg("row 42 missing from users")))
	fmt.Println(errors.PublicMessage(ErrAccountLocked.Msg("5 failed logins from 10.0.0.9")))
	fmt.Println(errors.PublicMessage(errors.New("unclassified")))
	// Output:
	// Internal Server Error
	// Not Found
	// Your account is locked
	// Internal Server Error
}

// WrapPublicMessage sets a client-facing message for a single error. The
// error text and identity don't change.
func ExampleWrapPublicMessage() {
	err := errors.ErrInvalidArgument.Msg("sku ABC-123 failed checksum")
	err = errors.WrapPublicMessage(err, "The product code is not valid")

	fmt.Println(err)
	fmt.Println(errors.PublicMessage(err))
	fmt.Println(errors.Is(err, errors.ErrInvalidArgument))
	// Output:
	// sku ABC-123 failed checksum
	// The product code is not valid
	// true
}

// Public returns a value that can be serialized directly in a response.
func ExamplePublic() {
	err := errors.ErrConflict.Msg("version 3 does not match 4")

	b, _ := json.Marshal(errors.Public(err))
	fmt.Println(string(b))
	// Output: {"code":"CONFLICT","message":"Conflict"}
}
