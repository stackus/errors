package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type receivedGRPCError struct {
	original error
	status   *status.Status

	class errorClassification
}

func (e *receivedGRPCError) Error() string {
	return e.original.Error()
}

func (e *receivedGRPCError) Unwrap() error {
	return e.original
}

func (e *receivedGRPCError) Is(target error) bool {
	if e.class.kind != nil {
		return e.class.kind.Is(target)
	}

	category, ok := target.(Error)
	return ok && e.class.category == category
}

func (e *receivedGRPCError) GRPCStatus() *status.Status {
	return e.status
}

func (e *receivedGRPCError) TypeCode() string {
	return e.class.typeCode
}

func (e *receivedGRPCError) HTTPCode() int {
	return e.class.httpCode
}

func (e *receivedGRPCError) GRPCCode() codes.Code {
	return e.class.grpcCode
}

func (e *receivedGRPCError) classification() errorClassification {
	return e.class
}

func SendGRPCError(err error) error {
	if err == nil {
		return nil
	}

	class, ok := resolveClassification(err)
	if !ok {
		class = unknownClassification()
	}

	if class.grpcCode == codes.OK {
		class = unknownClassification()
	}

	// Deliberately do not use err.Error() here.
	message := PublicMessage(err)

	s := status.New(class.grpcCode, message)

	detail := &ErrorType{
		TypeCode: class.typeCode,
		HTTPCode: int64(class.httpCode),
		GRPCCode: int64(class.grpcCode),
		Category: string(class.category),
	}

	withDetails, detailErr := s.WithDetails(detail)
	if detailErr == nil {
		s = withDetails
	}

	return s.Err()
}

func ReceiveGRPCError(err error, registries ...*Registry) error {
	if err == nil {
		return nil
	}

	s, ok := status.FromError(err)
	if !ok {
		// An ordinary Go error is not a received gRPC status.
		return err
	}

	actualCode := s.Code()

	class := categoryClassification(categoryFromGRPCCode(actualCode))

	// The actual gRPC status is authoritative.
	class.grpcCode = actualCode

	var registry *Registry
	if len(registries) != 0 {
		registry = registries[0]
	}

	for _, detail := range s.Details() {
		info, ok := detail.(*ErrorType)
		if !ok {
			continue
		}

		class = decodeGRPCClassification(
			info,
			actualCode,
			registry,
			class,
		)

		// Use the first valid classification detail.
		break
	}

	return &receivedGRPCError{
		original: err,
		status:   s,
		class:    class,
	}
}

func categoryFromGRPCCode(code codes.Code) Error {
	switch code {
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
		return ErrUnknown
	}
}

func decodeGRPCClassification(info *ErrorType, actualCode codes.Code, registry *Registry, fallback errorClassification) errorClassification {
	if info == nil {
		return fallback
	}

	code := info.GetTypeCode()
	if !typeCodeRe.MatchString(code) {
		return fallback
	}

	if len(code) > 128 {
		return fallback
	}

	// The metadata must not contradict the actual status.
	if info.GetGRPCCode() != int64(actualCode) {
		return fallback
	}

	// A registered application kind establishes exact identity.
	if kind, ok := registry.Lookup(code); ok {
		class := kind.classification()

		if class.grpcCode != actualCode {
			return fallback
		}

		if int64(class.httpCode) != info.GetHTTPCode() {
			return fallback
		}

		if info.GetCategory() != "" &&
			info.GetCategory() != string(class.category) {
			return fallback
		}

		return class
	}

	// Recognize built-in category codes independently
	// of application-specific registrations.
	if builtin, ok := categoryByCode(code); ok {
		class := categoryClassification(builtin)

		if class.grpcCode != actualCode {
			return fallback
		}

		if int64(class.httpCode) != info.GetHTTPCode() {
			return fallback
		}

		return class
	}

	// Unknown application-specific codes are preserved as
	// diagnostic metadata, but they do not establish a known
	// local error identity.
	fallback.typeCode = code

	return fallback
}

//func (e Error) GRPCStatus() *status.Status {
//	return errToStatus(e)
//}
//
//func (e *wrappedError) GRPCStatus() *status.Status {
//	return errToStatus(e)
//}

//type grpcError struct {
//	gc codes.Code
//	hc int
//	m  string
//	t  string
//	s  *status.Status
//}
//
//func (e grpcError) Error() string {
//	return e.m
//}
//
//func (e grpcError) GRPCStatus() *status.Status {
//	return e.s
//}
//
//func (e grpcError) HTTPCode() int {
//	return e.hc
//}
//
//func (e grpcError) GRPCCode() codes.Code {
//	return e.gc
//}
//
//func (e grpcError) TypeCode() string {
//	return e.t
//}
//
//// Is returns true if any of TypeCoder, HTTPCoder, GRPCCoder are a match between the error and target
//func (e grpcError) Is(target error) bool {
//	if t, ok := target.(GRPCCoder); ok && e.gc == t.GRPCCode() {
//		return true
//	}
//	if t, ok := target.(HTTPCoder); ok && e.hc == t.HTTPCode() {
//		return true
//	}
//	if t, ok := target.(TypeCoder); ok && e.t == t.TypeCode() {
//		return true
//	}
//	return false
//}

// SendGRPCError ensures that the error being used is sent with the correct code applied
//
// Use in the server when sending errors.
// If err is nil then SendGRPCError returns nil.
//func SendGRPCError(err error) error {
//	if err == nil {
//		return nil
//	}
//
//	// Already setup with a grpcCode
//	if _, ok := status.FromError(err); ok {
//		return err
//	}
//
//	s := errToStatus(err)
//
//	return s.Err()
//}

// ReceiveGRPCError recreates the error with the coded Error reapplied
//
// Non-nil results can be used as both Error and *status.Status. Methods
// errors.Is()/errors.As(), and status.Convert()/status.FromError() will
// continue to work.
//
// Use in the clients when receiving errors.
// If err is nil then ReceiveGRPCError returns nil.
//func ReceiveGRPCError(err error) error {
//	if err == nil {
//		return nil
//	}
//
//	s, ok := status.FromError(err)
//	if !ok {
//		return &grpcError{
//			gc: ErrUnknown.GRPCCode(),
//			hc: ErrUnknown.HTTPCode(),
//			m:  err.Error(),
//			t:  ErrUnknown.TypeCode(),
//			s:  s,
//		}
//	}
//
//	grpcCode := s.Code()
//	httpCode := ErrUnknown.HTTPCode()
//	embedType := codeToError(grpcCode).TypeCode()
//
//	for _, detail := range s.Details() {
//		switch d := detail.(type) {
//		case *ErrorType:
//			embedType = d.TypeCode
//			grpcCode = codes.Code(d.GRPCCode)
//			httpCode = int(d.HTTPCode)
//		}
//	}
//
//	return &grpcError{
//		gc: grpcCode,
//		hc: httpCode,
//		m:  s.Message(),
//		s:  s,
//		t:  embedType,
//	}
//}

//// convert a code to a known Error type;
//func codeToError(code codes.Code) Error {
//	switch code {
//	case codes.OK:
//		return ErrOK
//	case codes.Canceled:
//		return ErrCanceled
//	case codes.Unknown:
//		return ErrUnknown
//	case codes.InvalidArgument:
//		return ErrInvalidArgument
//	case codes.DeadlineExceeded:
//		return ErrDeadlineExceeded
//	case codes.NotFound:
//		return ErrNotFound
//	case codes.AlreadyExists:
//		return ErrAlreadyExists
//	case codes.PermissionDenied:
//		return ErrPermissionDenied
//	case codes.ResourceExhausted:
//		return ErrResourceExhausted
//	case codes.FailedPrecondition:
//		return ErrFailedPrecondition
//	case codes.Aborted:
//		return ErrAborted
//	case codes.OutOfRange:
//		return ErrOutOfRange
//	case codes.Unimplemented:
//		return ErrUnimplemented
//	case codes.Internal:
//		return ErrInternal
//	case codes.Unavailable:
//		return ErrUnavailable
//	case codes.DataLoss:
//		return ErrDataLoss
//	case codes.Unauthenticated:
//		return ErrUnauthenticated
//	default:
//		return ErrInternal
//	}
//}
//
//// convert an error into a gRPC *status.Status
//func errToStatus(err error) *status.Status {
//	grpcCode := ErrUnknown.GRPCCode()
//	httpCode := ErrUnknown.HTTPCode()
//	typeCode := ErrUnknown.TypeCode()
//
//	// Set the grpcCode based on GRPCCoder output; otherwise leave as Unknown
//	var grpcCoder GRPCCoder
//	if stderrors.As(err, &grpcCoder) {
//		grpcCode = grpcCoder.GRPCCode()
//	}
//
//	// short circuit building detailed errors if the code is OK
//	if grpcCode == codes.OK {
//		return status.New(codes.OK, "")
//	}
//
//	// Set the httpCode based on HTTPCoder output; otherwise leave as Unknown
//	var httpCoder HTTPCoder
//	if stderrors.As(err, &httpCoder) {
//		httpCode = httpCoder.HTTPCode()
//	}
//
//	// Embed the specific error "type"; otherwise leave as "UNKNOWN"
//	var typeCoder TypeCoder
//	if stderrors.As(err, &typeCoder) {
//		typeCode = typeCoder.TypeCode()
//	}
//
//	errInfo := &ErrorType{
//		TypeCode: typeCode,
//		GRPCCode: int64(grpcCode),
//		HTTPCode: int64(httpCode),
//	}
//
//	s, _ := status.New(grpcCode, err.Error()).WithDetails(errInfo)
//
//	return s
//}
