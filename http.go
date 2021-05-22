package errors

import "net/http"

type HTTPCoder interface {
	HTTPCode() int
}

func (e Error) HTTPCode() int {
	switch e {
	case ErrOK:
		return http.StatusOK
	case ErrCanceled:
		return http.StatusRequestTimeout
	case ErrUnknown:
		return http.StatusInternalServerError
	case ErrInvalidArgument:
		return http.StatusBadRequest
	case ErrDeadlineExceeded:
		return http.StatusGatewayTimeout
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists:
		return http.StatusConflict
	case ErrPermissionDenied:
		return http.StatusForbidden
	case ErrResourceExhausted:
		return http.StatusTooManyRequests
	case ErrFailedPrecondition:
		return http.StatusBadRequest
	case ErrAborted:
		return http.StatusConflict
	case ErrOutOfRange:
		return http.StatusUnprocessableEntity
	case ErrUnimplemented:
		return http.StatusNotImplemented
	case ErrInternal:
		return http.StatusInternalServerError
	case ErrUnavailable:
		return http.StatusServiceUnavailable
	case ErrDataLoss:
		return http.StatusInternalServerError
	case ErrUnauthenticated:
		return http.StatusUnauthorized
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrConflict:
		return http.StatusConflict
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrUnprocessableEntity:
		return http.StatusUnprocessableEntity
	case ErrServer:
		return http.StatusInternalServerError
	case ErrClient:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
