package errors_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stackus/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestCodeDefaults(t *testing.T) {
	type testCase struct {
		err      error
		typeCode string
		httpCode int
		grpcCode codes.Code
	}
	tests := map[string]testCase{
		"nil":                {nil, "OK", http.StatusOK, codes.OK},
		"unclassified error": {errors.New("test error"), "UNKNOWN", http.StatusInternalServerError, codes.Unknown},
		// a non-nil error never reports success
		"ErrOK":            {errors.ErrOK, "UNKNOWN", http.StatusInternalServerError, codes.Unknown},
		"ErrOK message":    {errors.ErrOK.Msg("x"), "UNKNOWN", http.StatusInternalServerError, codes.Unknown},
		"ErrOK wrapper":    {errors.ErrOK.Wrap(errors.New("db"), "x"), "UNKNOWN", http.StatusInternalServerError, codes.Unknown},
		"success coder":    {successCoder{}, "UNKNOWN", http.StatusInternalServerError, codes.Unknown},
		"wrapped as ErrOK": {fmt.Errorf("x: %w", errors.ErrOK), "UNKNOWN", http.StatusInternalServerError, codes.Unknown},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.typeCode, errors.TypeCode(tt.err))
			require.Equal(t, tt.httpCode, errors.HTTPCode(tt.err))
			require.Equal(t, tt.grpcCode, errors.GRPCCode(tt.err))
		})
	}
}

// successCoder is an error that claims success codes.
type successCoder struct{}

func (successCoder) Error() string        { return "success?" }
func (successCoder) TypeCode() string     { return "OK" }
func (successCoder) HTTPCode() int        { return http.StatusOK }
func (successCoder) GRPCCode() codes.Code { return codes.OK }

func TestErrOKMethods(t *testing.T) {
	require.Equal(t, "OK", errors.ErrOK.TypeCode())
	require.Equal(t, http.StatusOK, errors.ErrOK.HTTPCode())
	require.Equal(t, codes.OK, errors.ErrOK.GRPCCode())
}

// TypeCode finds the classification anywhere in an error's chain. Errors
// without a classification are UNKNOWN, and context errors map to their
// matching categories.
func ExampleTypeCode() {
	fmt.Println(errors.TypeCode(errors.ErrNotFound))
	fmt.Println(errors.TypeCode(fmt.Errorf("load: %w", errors.ErrForbidden.Msg("no access"))))
	fmt.Println(errors.TypeCode(errors.New("something failed")))
	fmt.Println(errors.TypeCode(context.DeadlineExceeded))
	fmt.Println(errors.TypeCode(nil))
	// Output:
	// NOT_FOUND
	// FORBIDDEN
	// UNKNOWN
	// DEADLINE_EXCEEDED
	// OK
}

func ExampleHTTPCode() {
	fmt.Println(errors.HTTPCode(errors.ErrNotFound))
	fmt.Println(errors.HTTPCode(errors.ErrBadRequest.Msg("missing name")))
	fmt.Println(errors.HTTPCode(errors.New("something failed")))
	fmt.Println(errors.HTTPCode(fmt.Errorf("query: %w", context.Canceled)))
	fmt.Println(errors.HTTPCode(nil))
	// Output:
	// 404
	// 400
	// 500
	// 408
	// 200
}

func ExampleGRPCCode() {
	fmt.Println(errors.GRPCCode(errors.ErrNotFound))
	fmt.Println(errors.GRPCCode(errors.ErrUnauthorized.Msg("token expired")))
	fmt.Println(errors.GRPCCode(errors.New("something failed")))
	fmt.Println(errors.GRPCCode(nil))
	// Output:
	// NotFound
	// Unauthenticated
	// Unknown
	// OK
}

// ValidationError is an application error type. It implements TypeCoder and
// HTTPCoder, so the package can classify it without wrapping.
type ValidationError struct {
	Field string
}

func (e ValidationError) Error() string    { return "invalid field: " + e.Field }
func (e ValidationError) TypeCode() string { return "VALIDATION_FAILED" }
func (e ValidationError) HTTPCode() int    { return http.StatusUnprocessableEntity }

// Your own error types are classified by implementing the coder interfaces.
// Codes a type does not provide use the defaults; here the gRPC code is
// Unknown because ValidationError does not implement GRPCCoder.
func ExampleHTTPCoder() {
	err := fmt.Errorf("create account: %w", ValidationError{Field: "email"})

	fmt.Println(errors.TypeCode(err))
	fmt.Println(errors.HTTPCode(err))
	fmt.Println(errors.GRPCCode(err))
	fmt.Println(errors.PublicMessage(err))
	// Output:
	// VALIDATION_FAILED
	// 422
	// Unknown
	// Unprocessable Entity
}
