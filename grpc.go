package errors

import (
	stderrors "errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCCoder interface {
	error
	GRPCCode() codes.Code
}

func (e Error) GRPCCode() codes.Code {
	switch e {
	// GRPC Errors
	case ErrOK:
		return codes.OK
	case ErrCanceled:
		return codes.Canceled
	case ErrUnknown:
		return codes.Unknown
	case ErrInvalidArgument:
		return codes.InvalidArgument
	case ErrDeadlineExceeded:
		return codes.DeadlineExceeded
	case ErrNotFound:
		return codes.NotFound
	case ErrAlreadyExists:
		return codes.AlreadyExists
	case ErrPermissionDenied:
		return codes.PermissionDenied
	case ErrResourceExhausted:
		return codes.ResourceExhausted
	case ErrFailedPrecondition:
		return codes.FailedPrecondition
	case ErrAborted:
		return codes.Aborted
	case ErrOutOfRange:
		return codes.OutOfRange
	case ErrUnimplemented:
		return codes.Unimplemented
	case ErrInternal:
		return codes.Internal
	case ErrUnavailable:
		return codes.Unavailable
	case ErrDataLoss:
		return codes.DataLoss
	case ErrUnauthenticated:
		return codes.Unauthenticated

	// HTTP Errors
	case ErrBadRequest:
		return codes.InvalidArgument
	case ErrUnauthorized:
		return codes.Unauthenticated
	case ErrForbidden:
		return codes.PermissionDenied
	case ErrMethodNotAllowed:
		return codes.Unimplemented
	case ErrRequestTimeout:
		return codes.DeadlineExceeded
	case ErrConflict:
		return codes.AlreadyExists
	case ErrGone:
		return codes.NotFound
	case ErrUnsupportedMediaType:
		return codes.InvalidArgument
	case ErrImATeapot:
		return codes.Unknown
	case ErrUnprocessableEntity:
		return codes.InvalidArgument
	case ErrTooManyRequests:
		return codes.ResourceExhausted
	case ErrUnavailableForLegalReasons:
		return codes.Unavailable
	case ErrInternalServerError:
		return codes.Internal
	case ErrNotImplemented:
		return codes.Unimplemented
	case ErrBadGateway:
		return codes.Aborted
	case ErrServiceUnavailable:
		return codes.Unavailable
	case ErrGatewayTimeout:
		return codes.DeadlineExceeded
	default:
		return codes.Internal
	}
}

func (e Error) GRPCStatus() *status.Status {
	return errToStatus(e)
}

func (e embeddedError) GRPCStatus() *status.Status {
	return errToStatus(e)
}

type grpcError struct {
	gc codes.Code
	hc int
	m  string
	t  string
	s  *status.Status
}

func (e grpcError) Error() string {
	return e.m
}

func (e grpcError) GRPCStatus() *status.Status {
	return e.s
}

func (e grpcError) HTTPCode() int {
	return e.hc
}

func (e grpcError) GRPCCode() codes.Code {
	return e.gc
}

func (e grpcError) TypeCode() string {
	return e.t
}

// Is returns true if any of TypeCoder, HTTPCoder, GRPCCoder are a match between the error and target
func (e grpcError) Is(target error) bool {
	if t, ok := target.(GRPCCoder); ok && e.gc == t.GRPCCode() {
		return true
	}
	if t, ok := target.(HTTPCoder); ok && e.hc == t.HTTPCode() {
		return true
	}
	if t, ok := target.(TypeCoder); ok && e.t == t.TypeCode() {
		return true
	}
	return false
}

// GRPCCode returns the GRPC code for the given error or codes.OK when nil or codes.Unknown otherwise
func GRPCCode(err error) codes.Code {
	if err == nil {
		return ErrOK.GRPCCode()
	}
	var e GRPCCoder
	if stderrors.As(err, &e) {
		return e.GRPCCode()
	}
	return ErrUnknown.GRPCCode()
}

// SendGRPCError ensures that the error being used is sent with the correct code applied
//
// Use in the server when sending errors.
// If err is nil then SendGRPCError returns nil.
func SendGRPCError(err error) error {
	if err == nil {
		return nil
	}

	// Already setup with a grpcCode
	if _, ok := status.FromError(err); ok {
		return err
	}

	s := errToStatus(err)

	return s.Err()
}

// ReceiveGRPCError recreates the error with the coded Error reapplied
//
// Non-nil results can be used as both Error and *status.Status. Methods
// errors.Is()/errors.As(), and status.Convert()/status.FromError() will
// continue to work.
//
// Use in the clients when receiving errors.
// If err is nil then ReceiveGRPCError returns nil.
func ReceiveGRPCError(err error) error {
	if err == nil {
		return nil
	}

	s, ok := status.FromError(err)
	if !ok {
		return &grpcError{
			gc: ErrUnknown.GRPCCode(),
			hc: ErrUnknown.HTTPCode(),
			m:  err.Error(),
			t:  ErrUnknown.TypeCode(),
			s:  s,
		}
	}

	grpcCode := s.Code()
	httpCode := ErrUnknown.HTTPCode()
	embedType := codeToError(grpcCode).TypeCode()

	for _, detail := range s.Details() {
		switch d := detail.(type) {
		case *ErrorType:
			embedType = d.TypeCode
			grpcCode = codes.Code(d.GRPCCode)
			httpCode = int(d.HTTPCode)
		}
	}

	return &grpcError{
		gc: grpcCode,
		hc: httpCode,
		m:  s.Message(),
		s:  s,
		t:  embedType,
	}
}

// convert a code to a known Error type;
func codeToError(code codes.Code) Error {
	switch code {
	case codes.OK:
		return ErrOK
	case codes.Canceled:
		return ErrCanceled
	case codes.Unknown:
		return ErrUnknown
	case codes.InvalidArgument:
		return ErrInvalidArgument
	case codes.DeadlineExceeded:
		return ErrDeadlineExceeded
	case codes.NotFound:
		return ErrNotFound
	case codes.AlreadyExists:
		return ErrAlreadyExists
	case codes.PermissionDenied:
		return ErrPermissionDenied
	case codes.ResourceExhausted:
		return ErrResourceExhausted
	case codes.FailedPrecondition:
		return ErrFailedPrecondition
	case codes.Aborted:
		return ErrAborted
	case codes.OutOfRange:
		return ErrOutOfRange
	case codes.Unimplemented:
		return ErrUnimplemented
	case codes.Internal:
		return ErrInternal
	case codes.Unavailable:
		return ErrUnavailable
	case codes.DataLoss:
		return ErrDataLoss
	case codes.Unauthenticated:
		return ErrUnauthenticated
	default:
		return ErrInternal
	}
}

// convert an error into a gRPC *status.Status
func errToStatus(err error) *status.Status {
	grpcCode := ErrUnknown.GRPCCode()
	httpCode := ErrUnknown.HTTPCode()
	typeCode := ErrUnknown.TypeCode()

	// Set the grpcCode based on GRPCCoder output; otherwise leave as Unknown
	var grpcCoder GRPCCoder
	if stderrors.As(err, &grpcCoder) {
		grpcCode = grpcCoder.GRPCCode()
	}

	// short circuit building detailed errors if the code is OK
	if grpcCode == codes.OK {
		return status.New(codes.OK, "")
	}

	// Set the httpCode based on HTTPCoder output; otherwise leave as Unknown
	var httpCoder HTTPCoder
	if stderrors.As(err, &httpCoder) {
		httpCode = httpCoder.HTTPCode()
	}

	// Embed the specific error "type"; otherwise leave as "UNKNOWN"
	var typeCoder TypeCoder
	if stderrors.As(err, &typeCoder) {
		typeCode = typeCoder.TypeCode()
	}

	errInfo := &ErrorType{
		TypeCode: typeCode,
		GRPCCode: int64(grpcCode),
		HTTPCode: int64(httpCode),
	}

	s, _ := status.New(grpcCode, err.Error()).WithDetails(errInfo)

	return s
}
