package errors

import (
	"google.golang.org/grpc/codes"
)

type wrappedError struct {
	kind     *Kind
	category Error
	cause    error
	msg      string
}

func (e *wrappedError) Error() string {
	return e.msg
}

func (e *wrappedError) Unwrap() error {
	return e.cause
}

func (e *wrappedError) Is(target error) bool {
	if e.kind != nil {
		return e.kind.Is(target)
	}

	category, ok := target.(Error)
	return ok && e.category == category
}

func (e *wrappedError) TypeCode() string {
	return e.classification().typeCode
}

func (e *wrappedError) HTTPCode() int {
	return e.classification().httpCode
}

func (e *wrappedError) GRPCCode() codes.Code {
	return e.classification().grpcCode
}

func (e *wrappedError) classification() errorClassification {
	if e.kind != nil {
		return e.kind.classification()
	}

	return categoryClassification(e.category)
}
