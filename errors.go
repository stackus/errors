package errors

import (
	stderrors "errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

// Error base error
type Error string

func New(message string) error {
	return stderrors.New(message)
}

func (e Error) Error() string {
	return string(e)
}

func (e Error) TypeCode() string {
	return string(e)
}

func (e Error) HTTPCode() int {
	def, ok := errorCategories[e]
	if !ok {
		return http.StatusInternalServerError
	}

	return def.httpCode
}

func (e Error) GRPCCode() codes.Code {
	cat, ok := errorCategories[e]
	if !ok {
		return codes.Unknown
	}

	return cat.grpcCode
}

// Msg sets a custom message for the Error
func (e Error) Msg(msg string) error {
	return &wrappedError{
		category: e,
		msg:      msg,
	}
}

// Msgf sets a custom message for formatting for the Error
func (e Error) Msgf(format string, args ...any) error {
	return &wrappedError{
		category: e,
		msg:      fmt.Sprintf(format, args...),
	}
}

func (e Error) WithCause(cause error) error {
	if cause == nil {
		return nil
	}
	return &wrappedError{
		category: e,
		cause:    cause,
		msg:      cause.Error(),
	}
}

func (e Error) Wrap(cause error, msg string) error {
	if cause == nil {
		return nil
	}

	return &wrappedError{
		category: e,
		cause:    cause,
		msg:      contextualMessage(msg, cause),
	}
}

func (e Error) Wrapf(cause error, format string, args ...any) error {
	if cause == nil {
		return nil
	}

	return e.Wrap(cause, fmt.Sprintf(format, args...))
}

func (e Error) classification() errorClassification {
	return categoryClassification(e)
}

// Wrap returns an error with msg wrapped with the supplied error
//
// If err is nil then Wrap returns nil
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	if msg == "" {
		return err
	}

	if e, ok := err.(Error); ok {
		return &wrappedError{
			category: e,
			cause:    err,
			msg:      msg,
		}
	}

	return fmt.Errorf("%s: %w", msg, err)
}

// Wrapf returns an error with a formatted msg wrapped with the supplied error
//
// If err is nil then Wrapf returns nil
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return Wrap(err, fmt.Sprintf(format, args...))
}

func contextualMessage(msg string, cause error) string {
	if cause == nil {
		return msg
	}

	if msg == "" {
		return cause.Error()
	}

	return msg + ": " + cause.Error()
}

// Go 1.26 convenience

func AsType[E error](err error) (E, bool) {
	return stderrors.AsType[E](err)
}

// Go 1.21 convenience

var ErrUnsupported = stderrors.ErrUnsupported

// Go 1.13 convenience

// As implements the standard errors.As for convenience
func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}

// Is implements the standard errors.Is for convenience
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// Unwrap implements the standard errors.Wrap for convenience
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// Join implements the standard errors.Join for convenience
func Join(errs ...error) error {
	return stderrors.Join(errs...)
}
