package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestError_Values(t *testing.T) {
	errors := []struct {
		e  Error
		t  string
		hc int
		gc codes.Code
	}{
		{ErrOK, "OK", 200, codes.OK},
		{ErrCanceled, "CANCELED", 408, codes.Canceled},
		{ErrUnknown, "UNKNOWN", 510, codes.Unknown},
		{ErrInvalidArgument, "INVALID_ARGUMENT", 400, codes.InvalidArgument},
		{ErrDeadlineExceeded, "DEADLINE_EXCEEDED", 504, codes.DeadlineExceeded},
		{ErrNotFound, "NOT_FOUND", 404, codes.NotFound},
		{ErrAlreadyExists, "ALREADY_EXISTS", 409, codes.AlreadyExists},
		{ErrPermissionDenied, "PERMISSION_DENIED", 403, codes.PermissionDenied},
		{ErrResourceExhausted, "RESOURCE_EXHAUSTED", 429, codes.ResourceExhausted},
		{ErrFailedPrecondition, "FAILED_PRECONDITION", 400, codes.FailedPrecondition},
		{ErrAborted, "ABORTED", 409, codes.Aborted},
		{ErrOutOfRange, "OUT_OF_RANGE", 422, codes.OutOfRange},
		{ErrUnimplemented, "UNIMPLEMENTED", 501, codes.Unimplemented},
		{ErrInternal, "INTERNAL", 500, codes.Internal},
		{ErrUnavailable, "UNAVAILABLE", 503, codes.Unavailable},
		{ErrDataLoss, "DATA_LOSS", 500, codes.DataLoss},
		{ErrUnauthenticated, "UNAUTHENTICATED", 401, codes.Unauthenticated},
		{ErrBadRequest, "BAD_REQUEST", 400, codes.InvalidArgument},
		{ErrUnauthorized, "UNAUTHORIZED", 401, codes.Unauthenticated},
		{ErrForbidden, "FORBIDDEN", 403, codes.PermissionDenied},
		{ErrMethodNotAllowed, "METHOD_NOT_ALLOWED", 405, codes.Unimplemented},
		{ErrRequestTimeout, "REQUEST_TIMEOUT", 408, codes.DeadlineExceeded},
		{ErrConflict, "CONFLICT", 409, codes.AlreadyExists},
		{ErrImATeapot, "IM_A_TEAPOT", 418, codes.Unknown},
		{ErrUnprocessableEntity, "UNPROCESSABLE_ENTITY", 422, codes.InvalidArgument},
		{ErrTooManyRequests, "TOO_MANY_REQUESTS", 429, codes.ResourceExhausted},
		{ErrUnavailableForLegalReasons, "UNAVAILABLE_FOR_LEGAL_REASONS", 451, codes.Unavailable},
		{ErrInternalServerError, "INTERNAL_SERVER_ERROR", 500, codes.Internal},
		{ErrNotImplemented, "NOT_IMPLEMENTED", 501, codes.Unimplemented},
		{ErrBadGateway, "BAD_GATEWAY", 502, codes.Aborted},
		{ErrServiceUnavailable, "SERVICE_UNAVAILABLE", 503, codes.Unavailable},
		{ErrGatewayTimeout, "GATEWAY_TIMEOUT", 504, codes.DeadlineExceeded},
	}
	for _, e := range errors {
		t.Run(e.t, func(t *testing.T) {
			assert.Equal(t, e.e.TypeCode(), e.t)
			assert.Equal(t, e.e.Error(), e.t)
			assert.Equal(t, e.e.HTTPCode(), e.hc)
			assert.Equal(t, e.e.GRPCCode(), e.gc)
		})
	}
}
