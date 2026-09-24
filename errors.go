package errors

import (
	stderrors "errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

// Error is a built-in error category. Its string value is its type code.
//
// Use the Err* constants in three ways:
//
//   - Return one directly: return errors.ErrNotFound
//   - Match against one: errors.Is(err, errors.ErrNotFound)
//   - Classify an error: errors.ErrNotFound.Wrap(err, "find user")
//
// Errors built from a category match it with [Is], including a [Kind]
// created in that category. Values that are not registered categories report
// HTTP 500 and gRPC Unknown.
type Error string

// New returns an ordinary, unclassified error with the given message.
func New(message string) error {
	return stderrors.New(message)
}

// Error returns the category's type code as its error text.
func (e Error) Error() string {
	return string(e)
}

// TypeCode returns the category's type code.
func (e Error) TypeCode() string {
	return string(e)
}

// HTTPCode returns the category's HTTP status, 200 for [ErrOK], or 500 if it
// is unregistered.
func (e Error) HTTPCode() int {
	if e == ErrOK {
		return http.StatusOK
	}

	def, ok := errorCategories[e]
	if !ok {
		return http.StatusInternalServerError
	}

	return def.httpCode
}

// GRPCCode returns the category's gRPC code, codes.OK for [ErrOK], or Unknown
// if it is unregistered.
func (e Error) GRPCCode() codes.Code {
	if e == ErrOK {
		return codes.OK
	}

	cat, ok := errorCategories[e]
	if !ok {
		return codes.Unknown
	}

	return cat.grpcCode
}

// Msg returns a classified error with msg as its full error text and no cause.
func (e Error) Msg(msg string) error {
	return &wrappedError{
		category: e,
		msg:      msg,
	}
}

// Msgf formats a message and returns the same classified result as [Error.Msg].
func (e Error) Msgf(format string, args ...any) error {
	return &wrappedError{
		category: e,
		msg:      fmt.Sprintf(format, args...),
	}
}

// WithCause classifies cause and uses cause.Error() as the error text.
// It returns nil when cause is nil.
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

// Wrap classifies cause and prefixes its error text with msg. An empty msg
// leaves the cause's text unchanged. It returns nil when cause is nil.
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

// Wrapf formats msg and calls [Error.Wrap]. It returns nil when cause is nil.
func (e Error) Wrapf(cause error, format string, args ...any) error {
	if cause == nil {
		return nil
	}

	return e.Wrap(cause, fmt.Sprintf(format, args...))
}

func (e Error) classification() errorClassification {
	return categoryClassification(e)
}

// Wrap adds msg to err. The result keeps err's classification, and [Is] and
// [As] can still reach err.
//
// If err is a category such as ErrNotFound, msg becomes the complete error
// text, so errors.Wrap(errors.ErrNotFound, "user 42 not found") reads
// "user 42 not found". Any other error, including a [Kind], reads
// "msg: err". Wrap returns err unchanged when msg is empty, and nil when err
// is nil.
//
// To classify or reclassify err, use a category's or a kind's Wrap method
// instead.
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

// Wrapf formats msg and calls [Wrap]. It returns nil when err is nil.
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

// AsType returns the first error in err's chain assignable to E, following
// the standard library's errors.AsType behavior.
func AsType[E error](err error) (E, bool) {
	return stderrors.AsType[E](err)
}

// ErrUnsupported is the standard library's errors.ErrUnsupported sentinel.
var ErrUnsupported = stderrors.ErrUnsupported

// As calls the standard library's errors.As.
func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}

// Is calls the standard library's errors.Is.
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// Unwrap calls the standard library's errors.Unwrap.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// Join calls the standard library's errors.Join.
func Join(errs ...error) error {
	return stderrors.Join(errs...)
}
