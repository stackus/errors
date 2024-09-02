package errors

import (
	stderrors "errors"
	"net/http"
)

type HTTPCoder interface {
	error
	HTTPCode() int
}

func (e Error) HTTPCode() int {
	switch e {
	// GRPC Errors
	case ErrOK:
		return http.StatusOK
	case ErrCanceled:
		return http.StatusRequestTimeout
	case ErrUnknown:
		return http.StatusNotExtended
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

	// HTTP Errors
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case ErrRequestTimeout:
		return http.StatusRequestTimeout
	case ErrConflict:
		return http.StatusConflict
	case ErrGone:
		return http.StatusGone
	case ErrUnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case ErrImATeapot:
		return 418 // teapot support
	case ErrUnprocessableEntity:
		return http.StatusUnprocessableEntity
	case ErrTooManyRequests:
		return http.StatusTooManyRequests
	case ErrUnavailableForLegalReasons:
		return http.StatusUnavailableForLegalReasons
	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrNotImplemented:
		return http.StatusNotImplemented
	case ErrBadGateway:
		return http.StatusBadGateway
	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrGatewayTimeout:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

// HTTPCode returns the HTTP status for the given error or http.StatusOK when nil or http.StatusNotExtended otherwise
func HTTPCode(err error) int {
	if err == nil {
		return ErrOK.HTTPCode()
	}

	var e HTTPCoder
	if stderrors.As(err, &e) {
		return e.HTTPCode()
	}
	return ErrUnknown.HTTPCode()
}
