package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestError_HTTPCode(t *testing.T) {
	assert.Equal(t, ErrOK.HTTPCode(), http.StatusOK)
	assert.Equal(t, ErrCanceled.HTTPCode(), http.StatusRequestTimeout)
	assert.Equal(t, ErrUnknown.HTTPCode(), http.StatusInternalServerError)
	assert.Equal(t, ErrInvalidArgument.HTTPCode(), http.StatusBadRequest)
	assert.Equal(t, ErrDeadlineExceeded.HTTPCode(), http.StatusGatewayTimeout)
	assert.Equal(t, ErrNotFound.HTTPCode(), http.StatusNotFound)
	assert.Equal(t, ErrAlreadyExists.HTTPCode(), http.StatusConflict)
	assert.Equal(t, ErrPermissionDenied.HTTPCode(), http.StatusForbidden)
	assert.Equal(t, ErrResourceExhausted.HTTPCode(), http.StatusTooManyRequests)
	assert.Equal(t, ErrFailedPrecondition.HTTPCode(), http.StatusBadRequest)
	assert.Equal(t, ErrAborted.HTTPCode(), http.StatusConflict)
	assert.Equal(t, ErrOutOfRange.HTTPCode(), http.StatusUnprocessableEntity)
	assert.Equal(t, ErrUnimplemented.HTTPCode(), http.StatusNotImplemented)
	assert.Equal(t, ErrInternal.HTTPCode(), http.StatusInternalServerError)
	assert.Equal(t, ErrUnavailable.HTTPCode(), http.StatusServiceUnavailable)
	assert.Equal(t, ErrDataLoss.HTTPCode(), http.StatusInternalServerError)
	assert.Equal(t, ErrUnauthenticated.HTTPCode(), http.StatusUnauthorized)
	assert.Equal(t, ErrBadRequest.HTTPCode(), http.StatusBadRequest)
	assert.Equal(t, ErrConflict.HTTPCode(), http.StatusConflict)
	assert.Equal(t, ErrUnauthorized.HTTPCode(), http.StatusUnauthorized)
	assert.Equal(t, ErrForbidden.HTTPCode(), http.StatusForbidden)
	assert.Equal(t, ErrUnprocessableEntity.HTTPCode(), http.StatusUnprocessableEntity)
	assert.Equal(t, ErrServer.HTTPCode(), http.StatusInternalServerError)
	assert.Equal(t, ErrClient.HTTPCode(), http.StatusBadRequest)
}

func TestError_GRPCCode(t *testing.T) {
	assert.Equal(t, ErrOK.GRPCCode(), codes.OK)
	assert.Equal(t, ErrCanceled.GRPCCode(), codes.Canceled)
	assert.Equal(t, ErrUnknown.GRPCCode(), codes.Unknown)
	assert.Equal(t, ErrInvalidArgument.GRPCCode(), codes.InvalidArgument)
	assert.Equal(t, ErrDeadlineExceeded.GRPCCode(), codes.DeadlineExceeded)
	assert.Equal(t, ErrNotFound.GRPCCode(), codes.NotFound)
	assert.Equal(t, ErrAlreadyExists.GRPCCode(), codes.AlreadyExists)
	assert.Equal(t, ErrPermissionDenied.GRPCCode(), codes.PermissionDenied)
	assert.Equal(t, ErrResourceExhausted.GRPCCode(), codes.ResourceExhausted)
	assert.Equal(t, ErrFailedPrecondition.GRPCCode(), codes.FailedPrecondition)
	assert.Equal(t, ErrAborted.GRPCCode(), codes.Aborted)
	assert.Equal(t, ErrOutOfRange.GRPCCode(), codes.OutOfRange)
	assert.Equal(t, ErrUnimplemented.GRPCCode(), codes.Unimplemented)
	assert.Equal(t, ErrInternal.GRPCCode(), codes.Internal)
	assert.Equal(t, ErrUnavailable.GRPCCode(), codes.Unavailable)
	assert.Equal(t, ErrDataLoss.GRPCCode(), codes.DataLoss)
	assert.Equal(t, ErrUnauthenticated.GRPCCode(), codes.Unauthenticated)
	assert.Equal(t, ErrBadRequest.GRPCCode(), codes.InvalidArgument)
	assert.Equal(t, ErrConflict.GRPCCode(), codes.AlreadyExists)
	assert.Equal(t, ErrUnauthorized.GRPCCode(), codes.Unauthenticated)
	assert.Equal(t, ErrForbidden.GRPCCode(), codes.PermissionDenied)
	assert.Equal(t, ErrUnprocessableEntity.GRPCCode(), codes.InvalidArgument)
	assert.Equal(t, ErrServer.GRPCCode(), codes.Internal)
	assert.Equal(t, ErrClient.GRPCCode(), codes.InvalidArgument)
}
