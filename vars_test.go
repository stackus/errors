package errors

import (
	stderrors "errors"
	"fmt"
	"net/http"

	_ "github.com/cucumber/godog"
	"google.golang.org/grpc/codes"
)

type typeTestError struct {
	t string
	e error
}

func (e typeTestError) Error() string              { return e.t }
func (e typeTestError) TypeCode() string           { return e.t }
func (e typeTestError) Is(err error) bool          { return stderrors.Is(e.e, err) }
func (e typeTestError) As(target interface{}) bool { return stderrors.As(e.e, target) }

type httpTestError struct {
	hc int
	e  error
}

func (e httpTestError) Error() string              { return fmt.Sprintf("HTTP(%d)", e.hc) }
func (e httpTestError) HTTPCode() int              { return e.hc }
func (e httpTestError) Is(err error) bool          { return stderrors.Is(e.e, err) }
func (e httpTestError) As(target interface{}) bool { return stderrors.As(e.e, target) }

type grpcTestError struct {
	gc codes.Code
	e  error
}

func (e grpcTestError) Error() string              { return fmt.Sprintf("GRPC(%d)", e.gc) }
func (e grpcTestError) GRPCCode() codes.Code       { return e.gc }
func (e grpcTestError) Is(err error) bool          { return stderrors.Is(e.e, err) }
func (e grpcTestError) As(target interface{}) bool { return stderrors.As(e.e, target) }

func convertGRPCStringToCode(grpcCode string) codes.Code {
	switch grpcCode {
	case "codes.OK":
		return codes.OK
	case "codes.Canceled":
		return codes.Canceled
	case "codes.Unknown":
		return codes.Unknown
	case "codes.InvalidArgument":
		return codes.InvalidArgument
	case "codes.DeadlineExceeded":
		return codes.DeadlineExceeded
	case "codes.NotFound":
		return codes.NotFound
	case "codes.AlreadyExists":
		return codes.AlreadyExists
	case "codes.PermissionDenied":
		return codes.PermissionDenied
	case "codes.ResourceExhausted":
		return codes.ResourceExhausted
	case "codes.FailedPrecondition":
		return codes.FailedPrecondition
	case "codes.Aborted":
		return codes.Aborted
	case "codes.OutOfRange":
		return codes.OutOfRange
	case "codes.Unimplemented":
		return codes.Unimplemented
	case "codes.Internal":
		return codes.Internal
	case "codes.Unavailable":
		return codes.Unavailable
	case "codes.DataLoss":
		return codes.DataLoss
	case "codes.Unauthenticated":
		return codes.Unauthenticated
	default:
		return codes.Unknown
	}
}

func convertHTTPStringToInt(httpStatus string) int {
	switch httpStatus {
	case "http.StatusOK":
		return http.StatusOK
	case "http.StatusBadRequest":
		return http.StatusBadRequest
	case "http.StatusUnauthorized":
		return http.StatusUnauthorized
	case "http.StatusPaymentRequired":
		return http.StatusPaymentRequired
	case "http.StatusForbidden":
		return http.StatusForbidden
	case "http.StatusNotFound":
		return http.StatusNotFound
	case "http.StatusMethodNotAllowed":
		return http.StatusMethodNotAllowed
	case "http.StatusNotAcceptable":
		return http.StatusNotAcceptable
	case "http.StatusProxyAuthRequired":
		return http.StatusProxyAuthRequired
	case "http.StatusRequestTimeout":
		return http.StatusRequestTimeout
	case "http.StatusConflict":
		return http.StatusConflict
	case "http.StatusGone":
		return http.StatusGone
	case "http.StatusLengthRequired":
		return http.StatusLengthRequired
	case "http.StatusPreconditionFailed":
		return http.StatusPreconditionFailed
	case "http.StatusRequestEntityTooLarge":
		return http.StatusRequestEntityTooLarge
	case "http.StatusRequestURITooLong":
		return http.StatusRequestURITooLong
	case "http.StatusUnsupportedMediaType":
		return http.StatusUnsupportedMediaType
	case "http.StatusRequestedRangeNotSatisfiable":
		return http.StatusRequestedRangeNotSatisfiable
	case "http.StatusExpectationFailed":
		return http.StatusExpectationFailed
	case "http.StatusTeapot":
		return http.StatusTeapot
	case "http.StatusMisdirectedRequest":
		return http.StatusMisdirectedRequest
	case "http.StatusUnprocessableEntity":
		return http.StatusUnprocessableEntity
	case "http.StatusLocked":
		return http.StatusLocked
	case "http.StatusFailedDependency":
		return http.StatusFailedDependency
	case "http.StatusTooEarly":
		return http.StatusTooEarly
	case "http.StatusUpgradeRequired":
		return http.StatusUpgradeRequired
	case "http.StatusPreconditionRequired":
		return http.StatusPreconditionRequired
	case "http.StatusTooManyRequests":
		return http.StatusTooManyRequests
	case "http.StatusRequestHeaderFieldsTooLarge":
		return http.StatusRequestHeaderFieldsTooLarge
	case "http.StatusUnavailableForLegalReasons":
		return http.StatusUnavailableForLegalReasons
	case "http.StatusInternalServerError":
		return http.StatusInternalServerError
	case "http.StatusNotImplemented":
		return http.StatusNotImplemented
	case "http.StatusBadGateway":
		return http.StatusBadGateway
	case "http.StatusServiceUnavailable":
		return http.StatusServiceUnavailable
	case "http.StatusGatewayTimeout":
		return http.StatusGatewayTimeout
	case "http.StatusHTTPVersionNotSupported":
		return http.StatusHTTPVersionNotSupported
	case "http.StatusVariantAlsoNegotiates":
		return http.StatusVariantAlsoNegotiates
	case "http.StatusInsufficientStorage":
		return http.StatusInsufficientStorage
	case "http.StatusLoopDetected":
		return http.StatusLoopDetected
	case "http.StatusNotExtended":
		return http.StatusNotExtended
	case "http.StatusNetworkAuthenticationRequired":
		return http.StatusNetworkAuthenticationRequired
	default:
		return http.StatusNotExtended
	}
}

func convertErrNameToError(errName string) Error {
	switch errName {
	case "ErrOK":
		return ErrOK
	case "ErrCanceled":
		return ErrCanceled
	case "ErrUnknown":
		return ErrUnknown
	case "ErrInvalidArgument":
		return ErrInvalidArgument
	case "ErrDeadlineExceeded":
		return ErrDeadlineExceeded
	case "ErrNotFound":
		return ErrNotFound
	case "ErrAlreadyExists":
		return ErrAlreadyExists
	case "ErrPermissionDenied":
		return ErrPermissionDenied
	case "ErrResourceExhausted":
		return ErrResourceExhausted
	case "ErrFailedPrecondition":
		return ErrFailedPrecondition
	case "ErrAborted":
		return ErrAborted
	case "ErrOutOfRange":
		return ErrOutOfRange
	case "ErrUnimplemented":
		return ErrUnimplemented
	case "ErrInternal":
		return ErrInternal
	case "ErrUnavailable":
		return ErrUnavailable
	case "ErrDataLoss":
		return ErrDataLoss
	case "ErrUnauthenticated":
		return ErrUnauthenticated
	case "ErrBadRequest":
		return ErrBadRequest
	case "ErrUnauthorized":
		return ErrUnauthorized
	case "ErrForbidden":
		return ErrForbidden
	case "ErrMethodNotAllowed":
		return ErrMethodNotAllowed
	case "ErrRequestTimeout":
		return ErrRequestTimeout
	case "ErrConflict":
		return ErrConflict
	case "ErrImATeapot":
		return ErrImATeapot
	case "ErrUnprocessableEntity":
		return ErrUnprocessableEntity
	case "ErrTooManyRequests":
		return ErrTooManyRequests
	case "ErrUnavailableForLegalReasons":
		return ErrUnavailableForLegalReasons
	case "ErrInternalServerError":
		return ErrInternalServerError
	case "ErrNotImplemented":
		return ErrNotImplemented
	case "ErrBadGateway":
		return ErrBadGateway
	case "ErrServiceUnavailable":
		return ErrServiceUnavailable
	case "ErrGatewayTimeout":
		return ErrGatewayTimeout
	default:
		return ErrUnknown
	}
}
