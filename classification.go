package errors

import (
	"context"
	stderrors "errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	classified interface {
		classification() errorClassification
	}

	singleUnwrapper interface {
		Unwrap() error
	}

	multiUnwrapper interface {
		Unwrap() []error
	}
)

func resolveClassification(err error) (errorClassification, bool) {
	if err == nil {
		return errorClassification{}, false
	}

	// explicitly classified outer error returned as-is
	if e, ok := err.(classified); ok {
		return e.classification(), true
	}

	if c, ok := externalClassification(err); ok {
		return c, true
	}

	// a gRPC status error, such as one from status.Error or a client call
	if e, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
		if s := e.GRPCStatus(); s != nil {
			return statusClassification(s, nil), true
		}
	}

	// check for single unwrap
	if e, ok := err.(singleUnwrapper); ok {
		return resolveClassification(e.Unwrap())
	}

	// check for multi unwrap
	if multi, ok := err.(multiUnwrapper); ok {
		var selected errorClassification

		found := false

		for _, child := range multi.Unwrap() {
			curr, ok := resolveClassification(child)
			if !ok && child != nil {
				return unknownClassification(), true
			}

			if !ok {
				continue
			}

			if !found {
				selected = curr
				found = true
				continue
			}

			if !sameClassification(selected, curr) {
				return unknownClassification(), true
			}
		}

		return selected, found
	}

	// special case for context errors, which are not classified by default
	switch {
	case stderrors.Is(err, context.Canceled):
		return categoryClassification(ErrCanceled), true
	case stderrors.Is(err, context.DeadlineExceeded):
		return categoryClassification(ErrDeadlineExceeded), true
	}

	return errorClassification{}, false
}

func externalClassification(err error) (errorClassification, bool) {
	c := unknownClassification()

	found := false

	if e, ok := err.(TypeCoder); ok {
		c.typeCode = e.TypeCode()
		found = true
	}

	if e, ok := err.(HTTPCoder); ok {
		c.httpCode = e.HTTPCode()
		found = true
	}

	if e, ok := err.(GRPCCoder); ok {
		c.grpcCode = e.GRPCCode()
		found = true
	}

	if !found {
		return errorClassification{}, false
	}

	// A non-nil error never reports success.
	if c.typeCode == "" || c.typeCode == string(ErrOK) {
		c.typeCode = string(ErrUnknown)
	}

	if c.httpCode < 400 || c.httpCode > 599 {
		c.httpCode = http.StatusInternalServerError
	}

	if c.grpcCode == codes.OK {
		c.grpcCode = codes.Unknown
	}

	return c, true
}

func sameClassification(a, b errorClassification) bool {
	return a.kind == b.kind &&
		a.category == b.category &&
		a.typeCode == b.typeCode &&
		a.httpCode == b.httpCode &&
		a.grpcCode == b.grpcCode &&
		a.publicMessage == b.publicMessage
}
