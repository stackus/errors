package errors

import (
	stderrors "errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
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

// SendGRPCError converts err into a gRPC status error for a server to return.
// The status code is [GRPCCode](err), and the message is [PublicMessage](err),
// so internal text never leaves the server. The type code, HTTP code, and
// category are attached as an [ErrorType] status detail, which
// [ReceiveGRPCError] reads on the client.
//
// Call it on every error a handler returns, usually from a server
// interceptor such as the ones in the grpcerrs subpackage. It returns nil
// for a nil err.
//
// A gRPC status error in err's chain, whether made with
// google.golang.org/grpc/status or received from another service, is
// classified by its code and detail. Unless an outer category or kind
// reclassifies it, its other status details are sent along too.
func SendGRPCError(err error) error {
	if err == nil {
		return nil
	}

	class, ok := resolveClassification(err)
	if !ok || class.grpcCode == codes.OK {
		class = unknownClassification()
	}

	// Deliberately do not use err.Error() here.
	message := PublicMessage(err)

	detail, detailErr := anypb.New(&ErrorType{
		TypeCode: class.typeCode,
		HTTPCode: int64(class.httpCode),
		GRPCCode: int64(class.grpcCode),
		Category: string(class.category),
	})

	s := status.New(class.grpcCode, message)
	if class.status != nil {
		// Keep the source status's other details.
		p := class.status.Proto()
		p.Code = int32(class.grpcCode)
		p.Message = message
		s = status.FromProto(p)
	}

	if detailErr == nil {
		s = withErrorTypeDetail(s, detail)
	}

	return s.Err()
}

// ReceiveGRPCError rebuilds a classified error from a gRPC status error, such
// as one returned by a client call. Call it on every error a client call
// returns, usually from a client interceptor such as the ones in the
// grpcerrs subpackage.
//
// The result is classified from the status code and the [ErrorType] detail
// that [SendGRPCError] attaches:
//
//   - A built-in category, such as ErrBadRequest, is restored as sent, with
//     the same codes and a match with [Is].
//   - A [Kind] found in one of the registries is restored as sent. The result
//     matches the Kind and its category, and has the local Kind's codes.
//     Registries are searched in order, and the first match is used.
//   - A Kind that isn't registered keeps its type code, category, and HTTP
//     code. The result matches the category, but not any Kind.
//
// [PublicMessage] of the result is the message the sender sent. The status
// code decides the result. When a detail disagrees with it, the detail is
// ignored. A status without a detail, for example from a server that does
// not use this package, is classified by its status code alone, and its
// message is not treated as public.
//
// The result's Error text is the status error's text, and its GRPCStatus
// method returns the received status, details included.
//
// It returns nil for nil and returns an error that isn't a gRPC status
// unchanged.
func ReceiveGRPCError(err error, registries ...*Registry) error {
	if err == nil {
		return nil
	}

	// status.FromError would rewrite the message of a wrapped status, so
	// find the status error itself.
	var gs interface{ GRPCStatus() *status.Status }
	if !stderrors.As(err, &gs) || gs.GRPCStatus() == nil {
		// An ordinary Go error is not a received gRPC status.
		return err
	}

	s := gs.GRPCStatus()

	return &receivedGRPCError{
		original: err,
		status:   s,
		class:    statusClassification(s, registries),
	}
}

// statusClassification classifies a gRPC status from its code and its first
// valid ErrorType detail.
func statusClassification(s *status.Status, registries []*Registry) errorClassification {
	actualCode := s.Code()
	if actualCode == codes.OK {
		actualCode = codes.Unknown
	}

	class := categoryClassification(categoryFromGRPCCode(actualCode))

	// The actual gRPC status is authoritative.
	class.grpcCode = actualCode

	for _, detail := range s.Details() {
		info, ok := detail.(*ErrorType)
		if !ok {
			continue
		}

		if decoded, ok := decodeGRPCClassification(info, actualCode, registries, class); ok {
			// The sender used PublicMessage as the status message.
			decoded.publicMessage = s.Message()
			class = decoded
			break
		}
	}

	class.status = s

	return class
}

// withErrorTypeDetail returns s with detail replacing any ErrorType details.
func withErrorTypeDetail(s *status.Status, detail *anypb.Any) *status.Status {
	p := s.Proto()

	details := make([]*anypb.Any, 0, len(p.Details)+1)
	for _, d := range p.Details {
		if !d.MessageIs((*ErrorType)(nil)) {
			details = append(details, d)
		}
	}

	p.Details = append(details, detail)

	return status.FromProto(p)
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

func decodeGRPCClassification(info *ErrorType, actualCode codes.Code, registries []*Registry, fallback errorClassification) (errorClassification, bool) {
	code := info.GetTypeCode()
	if len(code) > 128 || !typeCodeRe.MatchString(code) || code == string(ErrOK) {
		return fallback, false
	}

	// The metadata must not contradict the actual status.
	if info.GetGRPCCode() != int64(actualCode) {
		return fallback, false
	}

	sentCategory, sentCategoryOK := categoryByCode(info.GetCategory())

	// A registered application kind establishes exact identity. The local
	// definition decides the HTTP code.
	if kind, ok := lookupKind(registries, code); ok {
		class := kind.classification()

		if class.grpcCode == actualCode &&
			(info.GetCategory() == "" || info.GetCategory() == string(class.category)) {
			return class, true
		}
	}

	// Recognize built-in category codes independently
	// of application-specific registrations.
	if builtin, ok := categoryByCode(code); ok {
		class := categoryClassification(builtin)

		if class.grpcCode != actualCode ||
			int64(class.httpCode) != info.GetHTTPCode() {
			return fallback, false
		}

		return class, true
	}

	// Other application-specific codes keep the sender's type code,
	// category, and HTTP code, but they do not establish a known local
	// error identity.
	class := fallback
	class.typeCode = code

	if sentCategoryOK {
		class.category = sentCategory
		class.httpCode = sentCategory.HTTPCode()
	}

	if httpCode := info.GetHTTPCode(); httpCode >= 400 && httpCode <= 599 {
		class.httpCode = int(httpCode)
	}

	return class, true
}

func lookupKind(registries []*Registry, code string) (*Kind, bool) {
	for _, registry := range registries {
		if kind, ok := registry.Lookup(code); ok {
			return kind, true
		}
	}

	return nil, false
}
