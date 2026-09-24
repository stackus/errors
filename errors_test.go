package errors_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stackus/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

// ErrUserNotFound is an application sentinel error. It is a distinct error
// identity that also belongs to the built-in ErrNotFound category, so it maps
// to HTTP 404 and gRPC NotFound.
var ErrUserNotFound = errors.NewKind(
	"USER_NOT_FOUND",
	errors.ErrNotFound,
	"user not found",
	errors.WithPublicMessage("The user could not be found"),
)

// errNoRows stands in for a driver error such as sql.ErrNoRows.
var errNoRows = errors.New("sql: no rows in result set")

// This example follows an error from a repository up through a service. The
// repository classifies the driver error at the boundary, the service adds
// context, and the caller reads the codes and checks the error's identity.
func Example() {
	findUser := func(id string) error {
		// Classify the driver error where it enters the application.
		return ErrUserNotFound.Wrap(errNoRows, "find user "+id)
	}

	getProfile := func(id string) error {
		// Add context on the way up; the classification is preserved.
		return errors.Wrap(findUser(id), "get profile")
	}

	err := getProfile("42")

	fmt.Println(err)
	fmt.Println(errors.TypeCode(err))
	fmt.Println(errors.HTTPCode(err))
	fmt.Println(errors.GRPCCode(err))
	fmt.Println(errors.Is(err, ErrUserNotFound))
	fmt.Println(errors.Is(err, errors.ErrNotFound))
	fmt.Println(errors.Is(err, errNoRows))
	fmt.Println(errors.PublicMessage(err))
	// Output:
	// get profile: find user 42: sql: no rows in result set
	// USER_NOT_FOUND
	// 404
	// NotFound
	// true
	// true
	// true
	// The user could not be found
}

// This example writes an error as an HTTP response. The status comes from
// HTTPCode, and the body is the client-safe code and message from Public.
func Example_httpHandler() {
	writeError := func(w http.ResponseWriter, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errors.HTTPCode(err))
		_ = json.NewEncoder(w).Encode(errors.Public(err))
	}

	handle := func(err error) {
		rec := httptest.NewRecorder()
		writeError(rec, err)
		fmt.Print(rec.Code, " ", rec.Body.String())
	}

	handle(ErrUserNotFound.Wrap(errNoRows, "find user 42"))
	handle(errors.ErrInternal.Wrap(errors.New("dial tcp 10.0.0.5:5432: connection refused"), "load user"))
	handle(errors.ErrUnprocessableEntity.Msg("email: missing @"))
	// Output:
	// 404 {"code":"USER_NOT_FOUND","message":"The user could not be found"}
	// 500 {"code":"INTERNAL","message":"Internal Server Error"}
	// 422 {"code":"UNPROCESSABLE_ENTITY","message":"Unprocessable Entity"}
}

// Wrapping a built-in category directly uses the message as the complete error
// text. The result still matches the category.
func ExampleWrap() {
	err := errors.Wrap(errors.ErrNotFound, "user 42 not found")
	fmt.Println(err)
	fmt.Println(errors.Is(err, errors.ErrNotFound))
	// Output:
	// user 42 not found
	// true
}

// Each additional Wrap adds another prefix.
func ExampleWrap_multiple() {
	err := errors.Wrap(errors.ErrNotFound, "original message")
	err = errors.Wrap(err, "prefixed message")
	fmt.Println(err)
	// Output: prefixed message: original message
}

// Wrapping any other error produces "msg: cause" and keeps the cause and its
// classification reachable.
func ExampleWrap_ordinaryError() {
	err := errors.Wrap(errors.ErrForbidden.Msg("project is archived"), "update project")
	fmt.Println(err)
	fmt.Println(errors.TypeCode(err))
	// Output:
	// update project: project is archived
	// FORBIDDEN
}

// Use a category's Wrap to classify an error that came from somewhere else,
// such as a database driver or another library.
func ExampleError_Wrap() {
	err := errors.ErrNotFound.Wrap(errNoRows, "find order 7")
	fmt.Println(err)
	fmt.Println(errors.HTTPCode(err))
	fmt.Println(errors.Is(err, errors.ErrNotFound))
	fmt.Println(errors.Is(err, errNoRows))
	// Output:
	// find order 7: sql: no rows in result set
	// 404
	// true
	// true
}

// Msg creates a classified error from scratch.
func ExampleError_Msg() {
	err := errors.ErrInvalidArgument.Msg("page size must be positive")
	fmt.Println(err)
	fmt.Println(errors.GRPCCode(err))
	// Output:
	// page size must be positive
	// InvalidArgument
}

func ExampleError_Msgf() {
	err := errors.ErrTooManyRequests.Msgf("limit of %d requests per minute exceeded", 60)
	fmt.Println(err)
	fmt.Println(errors.HTTPCode(err))
	// Output:
	// limit of 60 requests per minute exceeded
	// 429
}

// WithCause classifies an error without changing its text.
func ExampleError_WithCause() {
	err := errors.ErrUnavailable.WithCause(errors.New("dial tcp: connection refused"))
	fmt.Println(err)
	fmt.Println(errors.TypeCode(err))
	// Output:
	// dial tcp: connection refused
	// UNAVAILABLE
}

// When a joined error contains different classifications, the join is
// UNKNOWN. Classify the join itself to choose the result.
func ExampleJoin() {
	joined := errors.Join(
		errors.ErrNotFound.Msg("user not found"),
		errors.ErrForbidden.Msg("project is private"),
	)
	fmt.Println(errors.TypeCode(joined))

	err := errors.ErrInvalidArgument.WithCause(joined)
	fmt.Println(errors.TypeCode(err))
	fmt.Println(errors.Is(err, errors.ErrForbidden))
	// Output:
	// UNKNOWN
	// INVALID_ARGUMENT
	// true
}

// AsType finds the first error in the chain that implements an interface, such
// as TypeCoder.
func ExampleAsType() {
	err := fmt.Errorf("checkout: %w", errors.ErrConflict.Msg("cart changed"))

	if coder, ok := errors.AsType[errors.TypeCoder](err); ok {
		fmt.Println(coder.TypeCode())
	}
	// Output: CONFLICT
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

func TestWrappedCategoryPreservesCause(t *testing.T) {
	err := errors.ErrBadRequest.Wrap(errors.ErrForbidden, "some error")

	require.Equal(t, http.StatusBadRequest, errors.HTTPCode(err))
	require.Equal(t, "some error: FORBIDDEN", err.Error())
	require.ErrorIs(t, err, errors.ErrBadRequest)
	require.ErrorIs(t, err, errors.ErrForbidden)
}

func TestErrorMessagesAndWrapping(t *testing.T) {
	cause := errors.New("record missing")
	type testCase struct {
		build      func() error
		message    string
		wantCause  error
		classified bool
	}
	tests := map[string]testCase{
		"formatted category message":     {func() error { return errors.ErrNotFound.Msgf("record %d", 42) }, "record 42", nil, true},
		"category with cause":            {func() error { return errors.ErrNotFound.WithCause(cause) }, "record missing", cause, true},
		"category wrap without prefix":   {func() error { return errors.ErrNotFound.Wrap(cause, "") }, "record missing", cause, true},
		"formatted category wrap":        {func() error { return errors.ErrNotFound.Wrapf(cause, "lookup %d", 42) }, "lookup 42: record missing", cause, true},
		"generic wrap of category":       {func() error { return errors.Wrap(errors.ErrNotFound, "lookup") }, "lookup", errors.ErrNotFound, true},
		"generic wrap of ordinary error": {func() error { return errors.Wrap(cause, "lookup") }, "lookup: record missing", cause, false},
		"formatted generic wrap":         {func() error { return errors.Wrapf(cause, "lookup %d", 42) }, "lookup 42: record missing", cause, false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.build()
			require.Equal(t, tt.message, err.Error())
			require.Equal(t, tt.wantCause, errors.Unwrap(err))
			if tt.classified {
				require.ErrorIs(t, err, errors.ErrNotFound)
				require.Equal(t, "NOT_FOUND", errors.TypeCode(err))
			}
		})
	}
}

func TestErrorWrappingNilAndEmpty(t *testing.T) {
	tests := map[string]func() error{
		"category WithCause nil": func() error { return errors.ErrNotFound.WithCause(nil) },
		"category Wrap nil":      func() error { return errors.ErrNotFound.Wrap(nil, "lookup") },
		"category Wrapf nil":     func() error { return errors.ErrNotFound.Wrapf(nil, "lookup %d", 42) },
		"generic Wrap nil":       func() error { return errors.Wrap(nil, "lookup") },
		"generic Wrapf nil":      func() error { return errors.Wrapf(nil, "lookup %d", 42) },
	}
	for name, call := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, call())
		})
	}

	cause := errors.New("record missing")
	require.Equal(t, cause, errors.Wrap(cause, ""))
}

func TestErrorAsHelpers(t *testing.T) {
	kind := errors.NewKind("MISSING_RECORD", errors.ErrNotFound, "record missing")
	err := fmt.Errorf("lookup: %w", kind)

	var target *errors.Kind
	require.True(t, errors.As(err, &target))
	require.Same(t, kind, target)
	got, ok := errors.AsType[*errors.Kind](err)
	require.True(t, ok)
	require.Same(t, kind, got)
}
