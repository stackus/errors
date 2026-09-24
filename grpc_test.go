package errors_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stackus/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	require.Equal(t, "User not found", errors.PublicMessage(received))
	require.Contains(t, received.Error(), "User not found")

	s, ok := status.FromError(received)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
}

func TestGRPCTransmission(t *testing.T) {
	permission := errors.NewKind(
		"CUSTOM_PERMISSION", errors.ErrUnknown, "permission error",
		errors.WithGRPCCode(codes.PermissionDenied),
	)
	mixed := errors.NewKind(
		"CUSTOM_BAD_REQUEST", errors.ErrBadRequest, "error message",
		errors.WithHTTPCode(http.StatusForbidden),
		errors.WithGRPCCode(codes.Unimplemented),
		errors.WithPublicMessage("Request was rejected"),
	)
	registry, err := errors.NewRegistry(permission, mixed)
	require.NoError(t, err)

	type testCase struct {
		input      error
		registry   *errors.Registry
		typeCode   string
		httpCode   int
		grpcCode   codes.Code
		message    string
		identities []error
	}
	tests := map[string]testCase{
		"nil": {
			typeCode: "OK", httpCode: http.StatusOK, grpcCode: codes.OK,
		},
		"ordinary error": {
			input:    errors.Wrap(errors.New("test error"), "standard error"),
			typeCode: "UNKNOWN", httpCode: http.StatusInternalServerError, grpcCode: codes.Unknown,
			message: "Internal Server Error", identities: []error{errors.ErrUnknown},
		},
		"built-in error": {
			input:    errors.ErrNotImplemented,
			typeCode: "NOT_IMPLEMENTED", httpCode: http.StatusNotImplemented, grpcCode: codes.Unimplemented,
			message: "Internal Server Error", identities: []error{errors.ErrNotImplemented},
		},
		"configured gRPC code": {
			input: permission, registry: registry,
			typeCode: "CUSTOM_PERMISSION", httpCode: http.StatusInternalServerError, grpcCode: codes.PermissionDenied,
			message: "Internal Server Error", identities: []error{permission, errors.ErrUnknown},
		},
		"configured codes and public message": {
			input: mixed, registry: registry,
			typeCode: "CUSTOM_BAD_REQUEST", httpCode: http.StatusForbidden, grpcCode: codes.Unimplemented,
			message: "Request was rejected", identities: []error{mixed, errors.ErrBadRequest},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			received := errors.ReceiveGRPCError(errors.SendGRPCError(tt.input), tt.registry)
			require.Equal(t, tt.typeCode, errors.TypeCode(received))
			require.Equal(t, tt.httpCode, errors.HTTPCode(received))
			require.Equal(t, tt.grpcCode, errors.GRPCCode(received))

			if tt.input == nil {
				require.NoError(t, received)
				return
			}

			require.Error(t, received)
			for _, identity := range tt.identities {
				require.ErrorIs(t, received, identity)
			}
			s, ok := status.FromError(received)
			require.True(t, ok)
			require.Equal(t, tt.grpcCode, s.Code())
			require.Equal(t, tt.message, s.Message())
			require.Equal(t, fmt.Sprintf("rpc error: code = %s desc = %s", tt.grpcCode, tt.message), received.Error())
		})
	}
}

// SendGRPCError turns an error into a gRPC status error. The status message is
// the public message, so internal details stay on the server.
func ExampleSendGRPCError() {
	err := errors.ErrInternal.Wrap(errors.New("redis: connection pool exhausted"), "load session")

	s, _ := status.FromError(errors.SendGRPCError(err))
	fmt.Println(s.Code())
	fmt.Println(s.Message())
	// Output:
	// Internal
	// Internal Server Error
}

// With a registry, the received error matches the same Kind the server sent.
// Without one, it keeps the type code, category, HTTP code, and public
// message, but it doesn't match the Kind.
func ExampleReceiveGRPCError() {
	ErrCartExpired := errors.NewKind("CART_EXPIRED", errors.ErrFailedPrecondition, "cart expired")
	registry, _ := errors.NewRegistry(ErrCartExpired)

	sent := errors.SendGRPCError(ErrCartExpired.Msg("cart 9 expired at 12:00"))

	withRegistry := errors.ReceiveGRPCError(sent, registry)
	fmt.Println(errors.Is(withRegistry, ErrCartExpired), errors.TypeCode(withRegistry), errors.HTTPCode(withRegistry))

	withoutRegistry := errors.ReceiveGRPCError(sent)
	fmt.Println(errors.Is(withoutRegistry, ErrCartExpired), errors.Is(withoutRegistry, errors.ErrFailedPrecondition), errors.TypeCode(withoutRegistry))

	fmt.Println(errors.PublicMessage(withRegistry))
	// Output:
	// true CART_EXPIRED 400
	// false true CART_EXPIRED
	// Bad Request
}

// Register the kinds a client expects to receive, typically once at startup.
func ExampleNewRegistry() {
	ErrOrderNotFound := errors.NewKind("ORDER_NOT_FOUND", errors.ErrNotFound, "order not found")
	ErrOrderShipped := errors.NewKind("ORDER_SHIPPED", errors.ErrFailedPrecondition, "order already shipped")

	registry, err := errors.NewRegistry(ErrOrderNotFound, ErrOrderShipped)
	if err != nil {
		panic(err)
	}

	kind, ok := registry.Lookup("ORDER_SHIPPED")
	fmt.Println(ok, kind == ErrOrderShipped)
	// Output: true true
}

// A status error from google.golang.org/grpc/status is classified by its
// code, so it can be returned from a handler or inspected directly.
func ExampleSendGRPCError_statusError() {
	err := fmt.Errorf("get order: %w", status.Error(codes.NotFound, "order 7 not found"))
	fmt.Println(errors.TypeCode(err), errors.HTTPCode(err))

	s, _ := status.FromError(errors.SendGRPCError(err))
	fmt.Println(s.Code(), s.Message())
	// Output:
	// NOT_FOUND 404
	// NotFound Not Found
}

func TestGRPCStatusErrorsAreClassified(t *testing.T) {
	plain := status.Error(codes.NotFound, "order 7 not found")

	tests := map[string]error{
		"plain":   plain,
		"wrapped": fmt.Errorf("get order: %w", plain),
	}

	for name, err := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, "NOT_FOUND", errors.TypeCode(err))
			require.Equal(t, http.StatusNotFound, errors.HTTPCode(err))
			require.Equal(t, codes.NotFound, errors.GRPCCode(err))
			// a status without a classification detail keeps the safe default
			require.Equal(t, "Not Found", errors.PublicMessage(err))

			s, ok := status.FromError(errors.SendGRPCError(err))
			require.True(t, ok)
			require.Equal(t, codes.NotFound, s.Code())
			require.Equal(t, "Not Found", s.Message())
		})
	}
}

func TestGRPCStatusDetailsArePreserved(t *testing.T) {
	extra := wrapperspb.String("sku") // stands in for any detail, such as errdetails.BadRequest
	withDetail, err := status.New(codes.InvalidArgument, "bad sku").WithDetails(extra)
	require.NoError(t, err)

	hasExtra := func(t *testing.T, err error) bool {
		t.Helper()
		s, ok := status.FromError(err)
		require.True(t, ok)
		for _, d := range s.Details() {
			if v, ok := d.(*wrapperspb.StringValue); ok && v.GetValue() == "sku" {
				return true
			}
		}
		return false
	}

	sent := errors.SendGRPCError(withDetail.Err())
	require.True(t, hasExtra(t, sent))

	// forwarded through a second hop
	forwarded := errors.SendGRPCError(errors.Wrap(errors.ReceiveGRPCError(sent), "call inventory"))
	require.True(t, hasExtra(t, forwarded))
	require.Equal(t, codes.InvalidArgument, status.Code(forwarded))
	require.Equal(t, "INVALID_ARGUMENT", errors.TypeCode(errors.ReceiveGRPCError(forwarded)))

	// reclassified errors start a fresh status
	reclassified := errors.SendGRPCError(errors.ErrInternal.Wrap(errors.ReceiveGRPCError(sent), "call inventory"))
	require.False(t, hasExtra(t, reclassified))
	require.Equal(t, codes.Internal, status.Code(reclassified))
}

func TestGRPCUnregisteredKindKeepsClassification(t *testing.T) {
	kind := errors.NewKind(
		"PAYMENT_REQUIRED", errors.ErrBadRequest, "payment required",
		errors.WithHTTPCode(http.StatusPaymentRequired),
	)

	received := errors.ReceiveGRPCError(errors.SendGRPCError(kind))

	require.Equal(t, "PAYMENT_REQUIRED", errors.TypeCode(received))
	require.Equal(t, http.StatusPaymentRequired, errors.HTTPCode(received))
	require.Equal(t, codes.InvalidArgument, errors.GRPCCode(received))
	require.ErrorIs(t, received, errors.ErrBadRequest)
	require.NotErrorIs(t, received, errors.ErrInvalidArgument)
	require.NotErrorIs(t, received, kind)
}

func TestGRPCPublicMessageIsRestored(t *testing.T) {
	kind := errors.NewKind(
		"ORDER_SHIPPED", errors.ErrFailedPrecondition, "order shipped",
		errors.WithPublicMessage("The order has shipped"),
	)
	registry, err := errors.NewRegistry(kind)
	require.NoError(t, err)

	type testCase struct {
		err      error
		registry *errors.Registry
		message  string
	}
	tests := map[string]testCase{
		"override on a registered kind": {
			err:      errors.WrapPublicMessage(kind.Msg("order 7 shipped"), "Order 7 already shipped"),
			registry: registry,
			message:  "Order 7 already shipped",
		},
		"unregistered kind": {
			err:     kind.Msg("order 7 shipped"),
			message: "The order has shipped",
		},
		"override on a category": {
			err:     errors.WrapPublicMessage(errors.ErrInvalidArgument.Msg("sku failed checksum"), "Bad SKU"),
			message: "Bad SKU",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			received := errors.ReceiveGRPCError(errors.SendGRPCError(tt.err), tt.registry)
			require.Equal(t, tt.message, errors.PublicMessage(received))

			// second hop
			received = errors.ReceiveGRPCError(errors.SendGRPCError(errors.Wrap(received, "forward")), tt.registry)
			require.Equal(t, tt.message, errors.PublicMessage(received))
		})
	}
}

func TestGRPCRegistries(t *testing.T) {
	first := errors.NewKind("FIRST", errors.ErrNotFound, "first")
	second := errors.NewKind("SECOND", errors.ErrNotFound, "second")

	a, err := errors.NewRegistry(first)
	require.NoError(t, err)
	b, err := errors.NewRegistry(second)
	require.NoError(t, err)

	received := errors.ReceiveGRPCError(errors.SendGRPCError(second), a, b)
	require.ErrorIs(t, received, second)
}

func TestGRPCRegisteredKindWithDifferentHTTPCode(t *testing.T) {
	// the sender's definition uses a different HTTP code than the receiver's
	sender := errors.NewKind("RATE_LIMITED", errors.ErrResourceExhausted, "rate limited", errors.WithHTTPCode(http.StatusServiceUnavailable))
	receiver := errors.NewKind("RATE_LIMITED", errors.ErrResourceExhausted, "rate limited")

	registry, err := errors.NewRegistry(receiver)
	require.NoError(t, err)

	received := errors.ReceiveGRPCError(errors.SendGRPCError(sender), registry)
	require.ErrorIs(t, received, receiver)
	require.Equal(t, http.StatusTooManyRequests, errors.HTTPCode(received))
}

func TestGRPCNeverSendsOrReceivesSuccess(t *testing.T) {
	sent := errors.SendGRPCError(errors.ErrOK)
	require.Error(t, sent)
	require.Equal(t, codes.Unknown, status.Code(sent))

	st, err := status.New(codes.NotFound, "Not Found").WithDetails(&errors.ErrorType{
		TypeCode: "OK", HTTPCode: http.StatusOK, GRPCCode: int64(codes.NotFound),
	})
	require.NoError(t, err)

	received := errors.ReceiveGRPCError(st.Err())
	require.Equal(t, "NOT_FOUND", errors.TypeCode(received))
	require.Equal(t, http.StatusNotFound, errors.HTTPCode(received))
}
