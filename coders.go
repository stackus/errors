package errors

import (
	"google.golang.org/grpc/codes"
)

type (
	// TypeCoder interface to extract an errors embeddable type as a string
	TypeCoder interface {
		error
		TypeCode() string
	}

	HTTPCoder interface {
		error
		HTTPCode() int
	}

	GRPCCoder interface {
		error
		GRPCCode() codes.Code
	}
)

// TypeCode returns the embedded type for the given error or blank when nil or UNKNOWN otherwise
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

// HTTPCode returns the HTTP status for the given error or http.StatusOK when nil or http.StatusNotExtended otherwise
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

// GRPCCode returns the GRPC code for the given error or codes.OK when nil or codes.Unknown otherwise
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
