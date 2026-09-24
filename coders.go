package errors

import (
	"google.golang.org/grpc/codes"
)

type (
	// TypeCoder is implemented by errors that report a type code.
	//
	// Every error in this package implements TypeCoder, [HTTPCoder], and
	// [GRPCCoder]. Your own error types can implement any of them to be
	// classified without wrapping. Codes a type does not provide use the
	// UNKNOWN defaults. Because a non-nil error never reports success, a type
	// code of "OK" is treated as "UNKNOWN", an HTTP code outside 400–599 as
	// 500, and a gRPC code of OK as Unknown.
	TypeCoder interface {
		error
		TypeCode() string
	}

	// HTTPCoder is implemented by errors that report an HTTP status code.
	// See [TypeCoder].
	HTTPCoder interface {
		error
		HTTPCode() int
	}

	// GRPCCoder is implemented by errors that report a gRPC status code.
	// See [TypeCoder].
	GRPCCoder interface {
		error
		GRPCCode() codes.Code
	}
)

// TypeCode returns err's type code. It uses the outermost classification in
// err's chain; a joined error with different classifications is "UNKNOWN". It
// returns "UNKNOWN" for an unclassified error and "OK" for nil. A non-nil
// error never returns "OK", even [ErrOK] itself. See the package
// documentation for the full rules.
func TypeCode(err error) string {
	if err == nil {
		return ErrOK.TypeCode()
	}

	c, ok := resolveClassification(err)
	if !ok {
		return ErrUnknown.TypeCode()
	}

	return c.typeCode
}

// HTTPCode returns err's HTTP status, using the same rules as [TypeCode]. It
// returns 200 for nil, 500 for an unclassified error, and a status in
// 400–599 for any non-nil error.
func HTTPCode(err error) int {
	if err == nil {
		return ErrOK.HTTPCode()
	}

	c, ok := resolveClassification(err)
	if !ok {
		return ErrUnknown.HTTPCode()
	}

	return c.httpCode
}

// GRPCCode returns err's gRPC code, using the same rules as [TypeCode]. It
// returns codes.OK for nil, codes.Unknown for an unclassified error, and never
// codes.OK for a non-nil error.
func GRPCCode(err error) codes.Code {
	if err == nil {
		return ErrOK.GRPCCode()
	}

	c, ok := resolveClassification(err)
	if !ok {
		return ErrUnknown.GRPCCode()
	}

	return c.grpcCode
}
