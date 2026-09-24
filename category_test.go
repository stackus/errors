package errors_test

import (
	"net/http"
	"testing"

	"github.com/stackus/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestBuiltinCodes(t *testing.T) {
	type testCase struct {
		err      errors.Error
		typeCode string
		httpCode int
		grpcCode codes.Code
	}
	tests := map[string]testCase{
		"ErrCanceled":                   {errors.ErrCanceled, "CANCELED", http.StatusRequestTimeout, codes.Canceled},
		"ErrUnknown":                    {errors.ErrUnknown, "UNKNOWN", http.StatusInternalServerError, codes.Unknown},
		"ErrInvalidArgument":            {errors.ErrInvalidArgument, "INVALID_ARGUMENT", http.StatusBadRequest, codes.InvalidArgument},
		"ErrDeadlineExceeded":           {errors.ErrDeadlineExceeded, "DEADLINE_EXCEEDED", http.StatusGatewayTimeout, codes.DeadlineExceeded},
		"ErrNotFound":                   {errors.ErrNotFound, "NOT_FOUND", http.StatusNotFound, codes.NotFound},
		"ErrAlreadyExists":              {errors.ErrAlreadyExists, "ALREADY_EXISTS", http.StatusConflict, codes.AlreadyExists},
		"ErrPermissionDenied":           {errors.ErrPermissionDenied, "PERMISSION_DENIED", http.StatusForbidden, codes.PermissionDenied},
		"ErrResourceExhausted":          {errors.ErrResourceExhausted, "RESOURCE_EXHAUSTED", http.StatusTooManyRequests, codes.ResourceExhausted},
		"ErrFailedPrecondition":         {errors.ErrFailedPrecondition, "FAILED_PRECONDITION", http.StatusBadRequest, codes.FailedPrecondition},
		"ErrAborted":                    {errors.ErrAborted, "ABORTED", http.StatusConflict, codes.Aborted},
		"ErrOutOfRange":                 {errors.ErrOutOfRange, "OUT_OF_RANGE", http.StatusUnprocessableEntity, codes.OutOfRange},
		"ErrUnimplemented":              {errors.ErrUnimplemented, "UNIMPLEMENTED", http.StatusNotImplemented, codes.Unimplemented},
		"ErrInternal":                   {errors.ErrInternal, "INTERNAL", http.StatusInternalServerError, codes.Internal},
		"ErrUnavailable":                {errors.ErrUnavailable, "UNAVAILABLE", http.StatusServiceUnavailable, codes.Unavailable},
		"ErrDataLoss":                   {errors.ErrDataLoss, "DATA_LOSS", http.StatusInternalServerError, codes.DataLoss},
		"ErrUnauthenticated":            {errors.ErrUnauthenticated, "UNAUTHENTICATED", http.StatusUnauthorized, codes.Unauthenticated},
		"ErrBadRequest":                 {errors.ErrBadRequest, "BAD_REQUEST", http.StatusBadRequest, codes.InvalidArgument},
		"ErrUnauthorized":               {errors.ErrUnauthorized, "UNAUTHORIZED", http.StatusUnauthorized, codes.Unauthenticated},
		"ErrForbidden":                  {errors.ErrForbidden, "FORBIDDEN", http.StatusForbidden, codes.PermissionDenied},
		"ErrMethodNotAllowed":           {errors.ErrMethodNotAllowed, "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, codes.Unimplemented},
		"ErrRequestTimeout":             {errors.ErrRequestTimeout, "REQUEST_TIMEOUT", http.StatusRequestTimeout, codes.DeadlineExceeded},
		"ErrConflict":                   {errors.ErrConflict, "CONFLICT", http.StatusConflict, codes.Aborted},
		"ErrImATeapot":                  {errors.ErrImATeapot, "IM_A_TEAPOT", http.StatusTeapot, codes.Unknown},
		"ErrUnprocessableEntity":        {errors.ErrUnprocessableEntity, "UNPROCESSABLE_ENTITY", http.StatusUnprocessableEntity, codes.InvalidArgument},
		"ErrTooManyRequests":            {errors.ErrTooManyRequests, "TOO_MANY_REQUESTS", http.StatusTooManyRequests, codes.ResourceExhausted},
		"ErrUnavailableForLegalReasons": {errors.ErrUnavailableForLegalReasons, "UNAVAILABLE_FOR_LEGAL_REASONS", http.StatusUnavailableForLegalReasons, codes.PermissionDenied},
		"ErrInternalServerError":        {errors.ErrInternalServerError, "INTERNAL_SERVER_ERROR", http.StatusInternalServerError, codes.Internal},
		"ErrNotImplemented":             {errors.ErrNotImplemented, "NOT_IMPLEMENTED", http.StatusNotImplemented, codes.Unimplemented},
		"ErrBadGateway":                 {errors.ErrBadGateway, "BAD_GATEWAY", http.StatusBadGateway, codes.Unavailable},
		"ErrServiceUnavailable":         {errors.ErrServiceUnavailable, "SERVICE_UNAVAILABLE", http.StatusServiceUnavailable, codes.Unavailable},
		"ErrGatewayTimeout":             {errors.ErrGatewayTimeout, "GATEWAY_TIMEOUT", http.StatusGatewayTimeout, codes.DeadlineExceeded},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.typeCode, errors.TypeCode(tt.err))
			require.Equal(t, tt.httpCode, errors.HTTPCode(tt.err))
			require.Equal(t, tt.grpcCode, errors.GRPCCode(tt.err))
		})
	}
}
