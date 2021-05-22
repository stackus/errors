package errors

import (
	stderrors "errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcError struct {
	e error
	s *status.Status
}

type GRPCCoder interface {
	GRPCCode() codes.Code
}

func (e Error) GRPCCode() codes.Code {
	switch e {
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
	case ErrBadRequest:
		return codes.InvalidArgument
	case ErrConflict:
		return codes.AlreadyExists
	case ErrUnauthorized:
		return codes.Unauthenticated
	case ErrForbidden:
		return codes.PermissionDenied
	case ErrUnprocessableEntity:
		return codes.InvalidArgument
	case ErrServer:
		return codes.Internal
	case ErrClient:
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}

func (e grpcError) Error() string {
	return e.e.Error()
}

func (e grpcError) GRPCStatus() *status.Status {
	return e.s
}

func (e grpcError) As(target interface{}) bool {
	return stderrors.As(e.e, target)
}

func (e grpcError) Is(target error) bool {
	return stderrors.Is(e.e, target)
}

// SendGRPCError ensures that the error being used is sent with the correct code applied
//
// Use in the server when sending errors.
// If err is nil then SendGRPCError returns nil.
func SendGRPCError(err error) error {
	if err == nil {
		return nil
	}

	// Already setup with a code
	if _, ok := status.FromError(err); ok {
		return err
	}

	var coder GRPCCoder
	if stderrors.As(err, &coder) {
		return status.Error(coder.GRPCCode(), err.Error())
	}

	// Return an unknown error code to help identify improperly coded errors
	return status.Error(codes.Unknown, err.Error())
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
			e: embed(ErrUnknown, err.Error()),
			s: s,
		}
	}

	return &grpcError{
		e: embed(codeToError(s.Code()), s.Message()),
		s: s,
	}
}

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
