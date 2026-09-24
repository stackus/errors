package errors_test

import (
	"fmt"
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

func TestKindCodeOptions(t *testing.T) {
	type testCase struct {
		kind     *errors.Kind
		typeCode string
		httpCode int
		grpcCode codes.Code
	}
	tests := map[string]testCase{
		"custom type code": {
			kind:     errors.NewKind("CUSTOM", errors.ErrUnknown, "custom error"),
			typeCode: "CUSTOM", httpCode: http.StatusInternalServerError, grpcCode: codes.Unknown,
		},
		"custom HTTP code": {
			kind:     errors.NewKind("CUSTOM_HTTP", errors.ErrUnknown, "custom HTTP error", errors.WithHTTPCode(http.StatusTeapot)),
			typeCode: "CUSTOM_HTTP", httpCode: http.StatusTeapot, grpcCode: codes.Unknown,
		},
		"custom gRPC code": {
			kind:     errors.NewKind("CUSTOM_GRPC", errors.ErrUnknown, "custom gRPC error", errors.WithGRPCCode(codes.AlreadyExists)),
			typeCode: "CUSTOM_GRPC", httpCode: http.StatusInternalServerError, grpcCode: codes.AlreadyExists,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.typeCode, errors.TypeCode(tt.kind))
			require.Equal(t, tt.httpCode, errors.HTTPCode(tt.kind))
			require.Equal(t, tt.grpcCode, errors.GRPCCode(tt.kind))
		})
	}
}

func TestNewKindValidation(t *testing.T) {
	type testCase struct {
		newKind   func()
		panicText string
	}
	tests := map[string]testCase{
		"invalid type code": {
			func() { _ = errors.NewKind("invalid", errors.ErrNotFound, "missing") },
			"errors: invalid type code: invalid",
		},
		"invalid category": {
			func() { _ = errors.NewKind("CUSTOM", errors.Error("MISSING"), "missing") },
			"errors: invalid category: MISSING",
		},
		"reserved OK type code": {
			func() { _ = errors.NewKind("OK", errors.ErrNotFound, "missing") },
			"errors: reserved kind type code: OK",
		},
		"reserved type code": {
			func() { _ = errors.NewKind("NOT_FOUND", errors.ErrNotFound, "missing") },
			"errors: reserved kind type code: NOT_FOUND",
		},
		"empty message": {
			func() { _ = errors.NewKind("CUSTOM", errors.ErrNotFound, "") },
			"errors: message cannot be empty",
		},
		"HTTP code below error range": {
			func() {
				_ = errors.NewKind("CUSTOM", errors.ErrNotFound, "missing", errors.WithHTTPCode(http.StatusFound))
			},
			"errors: invalid http code: 302",
		},
		"HTTP code above error range": {
			func() { _ = errors.NewKind("CUSTOM", errors.ErrNotFound, "missing", errors.WithHTTPCode(600)) },
			"errors: invalid http code: 600",
		},
		"successful gRPC code": {
			func() { _ = errors.NewKind("CUSTOM", errors.ErrNotFound, "missing", errors.WithGRPCCode(codes.OK)) },
			"errors: invalid grpc code: OK",
		},
		"empty public message": {
			func() { _ = errors.NewKind("CUSTOM", errors.ErrNotFound, "missing", errors.WithPublicMessage("")) },
			"errors: public message cannot be empty",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.PanicsWithValue(t, tt.panicText, tt.newKind)
		})
	}
}

func TestKindWrapping(t *testing.T) {
	kind := errors.NewKind("STORAGE_FAILURE", errors.ErrInternal, "storage failed", errors.WithPublicMessage("Storage unavailable"))
	cause := errors.New("disk failed")
	require.Equal(t, "storage failed", kind.Error())
	require.Equal(t, "STORAGE_FAILURE", kind.TypeCode())
	require.Equal(t, http.StatusInternalServerError, kind.HTTPCode())
	require.Equal(t, codes.Internal, kind.GRPCCode())
	require.Equal(t, "Storage unavailable", kind.PublicError())

	type testCase struct {
		wrap      func() error
		message   string
		wantCause error
	}
	tests := map[string]testCase{
		"message":              {func() error { return kind.Msg("custom message") }, "custom message", nil},
		"formatted message":    {func() error { return kind.Msgf("record %d", 42) }, "record 42", nil},
		"cause":                {func() error { return kind.WithCause(cause) }, "disk failed", cause},
		"wrapped cause":        {func() error { return kind.Wrap(cause, "load") }, "load: disk failed", cause},
		"formatted wrap":       {func() error { return kind.Wrapf(cause, "load %d", 42) }, "load 42: disk failed", cause},
		"WithKind convenience": {func() error { return errors.WithKind(cause, kind) }, "disk failed", cause},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.wrap()
			require.Equal(t, tt.message, err.Error())
			require.Equal(t, tt.wantCause, errors.Unwrap(err))
			require.ErrorIs(t, err, kind)
			require.ErrorIs(t, err, errors.ErrInternal)
			require.Equal(t, "STORAGE_FAILURE", errors.TypeCode(err))
		})
	}
}

func TestKindNilCause(t *testing.T) {
	kind := errors.NewKind("STORAGE_FAILURE", errors.ErrInternal, "storage failed")
	tests := map[string]func() error{
		"WithCause": func() error { return kind.WithCause(nil) },
		"Wrap":      func() error { return kind.Wrap(nil, "load") },
		"Wrapf":     func() error { return kind.Wrapf(nil, "load %d", 42) },
		"WithKind":  func() error { return errors.WithKind(nil, kind) },
	}

	for name, call := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, call())
		})
	}

	require.PanicsWithValue(t, "errors: kind cannot be nil", func() {
		_ = errors.WithKind(errors.New("disk failed"), nil)
	})
}

// Declare kinds as package-level variables and use them like any other
// sentinel error. Each kind is its own identity, and it also matches the
// category it belongs to.
func ExampleNewKind() {
	var (
		ErrOrderNotFound   = errors.NewKind("ORDER_NOT_FOUND", errors.ErrNotFound, "order not found")
		ErrInvoiceNotFound = errors.NewKind("INVOICE_NOT_FOUND", errors.ErrNotFound, "invoice not found")
	)

	err := ErrOrderNotFound.Msgf("order %d not found", 7)

	fmt.Println(errors.Is(err, ErrOrderNotFound))
	fmt.Println(errors.Is(err, ErrInvoiceNotFound))
	fmt.Println(errors.Is(err, errors.ErrNotFound))
	fmt.Println(errors.TypeCode(err), errors.HTTPCode(err), errors.GRPCCode(err))
	// Output:
	// true
	// false
	// true
	// ORDER_NOT_FOUND 404 NotFound
}

// Options override the codes a kind inherits from its category and set the
// message that clients see.
func ExampleNewKind_options() {
	ErrPaymentRequired := errors.NewKind(
		"PAYMENT_REQUIRED",
		errors.ErrFailedPrecondition,
		"account has no active subscription",
		errors.WithHTTPCode(http.StatusPaymentRequired),
		errors.WithPublicMessage("A subscription is required"),
	)
	ErrQuotaExceeded := errors.NewKind(
		"QUOTA_EXCEEDED",
		errors.ErrForbidden,
		"storage quota exceeded",
		errors.WithGRPCCode(codes.ResourceExhausted),
	)

	fmt.Println(errors.HTTPCode(ErrPaymentRequired), errors.GRPCCode(ErrPaymentRequired))
	fmt.Println(errors.PublicMessage(ErrPaymentRequired))
	fmt.Println(errors.HTTPCode(ErrQuotaExceeded), errors.GRPCCode(ErrQuotaExceeded))
	fmt.Println(errors.PublicMessage(ErrQuotaExceeded))
	// Output:
	// 402 FailedPrecondition
	// A subscription is required
	// 403 ResourceExhausted
	// Forbidden
}

// Wrap classifies a cause with the kind and adds context to its message.
func ExampleKind_Wrap() {
	ErrEmailTaken := errors.NewKind("EMAIL_TAKEN", errors.ErrAlreadyExists, "email already registered")

	cause := errors.New(`pq: duplicate key value violates unique constraint "users_email_key"`)
	err := ErrEmailTaken.Wrap(cause, "register user")

	fmt.Println(err)
	fmt.Println(errors.Is(err, ErrEmailTaken))
	fmt.Println(errors.HTTPCode(err))
	// Output:
	// register user: pq: duplicate key value violates unique constraint "users_email_key"
	// true
	// 409
}

// Msgf creates an error of the kind with a detailed message and no cause.
func ExampleKind_Msgf() {
	ErrInsufficientFunds := errors.NewKind("INSUFFICIENT_FUNDS", errors.ErrFailedPrecondition, "insufficient funds")

	err := ErrInsufficientFunds.Msgf("balance %d is less than %d", 30, 45)

	fmt.Println(err)
	fmt.Println(errors.TypeCode(err))
	// Output:
	// balance 30 is less than 45
	// INSUFFICIENT_FUNDS
}

// WithKind classifies an error with a kind and keeps its message. It returns
// nil for a nil error, so it can wrap a call's result directly.
func ExampleWithKind() {
	ErrStorage := errors.NewKind("STORAGE_UNAVAILABLE", errors.ErrUnavailable, "storage unavailable")

	save := func(fail bool) error {
		if fail {
			return errors.New("s3: request timed out")
		}
		return nil
	}

	fmt.Println(errors.WithKind(save(false), ErrStorage))

	err := errors.WithKind(save(true), ErrStorage)
	fmt.Println(err)
	fmt.Println(errors.TypeCode(err))
	// Output:
	// <nil>
	// s3: request timed out
	// STORAGE_UNAVAILABLE
}
