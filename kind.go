package errors

import (
	"fmt"
	"regexp"

	"google.golang.org/grpc/codes"
)

type (
	// Kind is an application-specific sentinel error inside a built-in
	// [Error] category. Create kinds with [NewKind] and declare them as
	// package-level variables:
	//
	//	var ErrUserNotFound = errors.NewKind("USER_NOT_FOUND", errors.ErrNotFound, "user not found")
	//
	// Kinds are compared by pointer, so two kinds with the same type code are
	// still different errors. Errors created from a kind (with Msg, Wrap,
	// WithCause, and so on) match both the kind and its category with [Is].
	// Its type code, HTTP code, gRPC code, and public message are set when the
	// kind is created and cannot change.
	Kind struct {
		id *struct{ _ byte } // non-zero identity ensures two independent kinds are never equal

		category Error
		typeCode string
		httpCode int
		grpcCode codes.Code

		message       string
		publicMessage string
	}

	// KindOption configures a Kind during [NewKind].
	KindOption func(*Kind)
)

var typeCodeRe = regexp.MustCompile(`^[A-Z][A-Z0-9_.]*$`)

// NewKind creates an application-specific error with a distinct type code.
// The type code must start with an uppercase letter and contain only uppercase
// letters, digits, underscores, or periods. Its category supplies default
// HTTP and gRPC codes; options can override them. message is the error text
// of the kind itself and of errors made with its WithCause method.
//
// NewKind panics for an invalid type code or one that matches a built-in
// category or ErrOK, an unregistered category (including ErrOK), an empty message, an
// HTTP code outside 400–599, or a gRPC code of OK. Kinds are normally
// declared as package-level variables, so these mistakes surface when the
// program starts.
func NewKind(typeCode string, category Error, message string, options ...KindOption) *Kind {
	if !typeCodeRe.MatchString(typeCode) {
		panic("errors: invalid type code: " + typeCode)
	}

	if _, ok := errorCategories[category]; !ok {
		panic("errors: invalid category: " + string(category))
	}

	if typeCode == string(ErrOK) {
		panic("errors: reserved kind type code: " + typeCode)
	}

	for builtin := range errorCategories {
		if typeCode == string(builtin) {
			panic("errors: reserved kind type code: " + typeCode)
		}
	}

	if message == "" {
		panic("errors: message cannot be empty")
	}

	k := &Kind{
		category: category,
		typeCode: typeCode,
		httpCode: category.HTTPCode(),
		grpcCode: category.GRPCCode(),
		message:  message,
	}

	for _, option := range options {
		option(k)
	}

	if k.httpCode < 400 || k.httpCode > 599 {
		panic("errors: invalid http code: " + fmt.Sprint(k.httpCode))
	}

	if k.grpcCode == codes.OK {
		panic("errors: invalid grpc code: " + k.grpcCode.String())
	}

	return k
}

// WithKind classifies err with kind and keeps err's text and its place in the
// chain. It is the same as kind.WithCause(err), except that it returns nil
// for a nil err before checking kind, so it can wrap a call's result
// directly. It panics if kind is nil and err is not.
func WithKind(err error, kind *Kind) error {
	if err == nil {
		return nil
	}

	if kind == nil {
		panic("errors: kind cannot be nil")
	}

	return kind.WithCause(err)
}

// Error returns the Kind's default error message.
func (k *Kind) Error() string {
	return k.message
}

// Is reports whether target is this Kind or this Kind's category.
func (k *Kind) Is(target error) bool {
	if t, ok := target.(*Kind); ok && k == t {
		return k == t
	}

	if category, ok := target.(Error); ok {
		return k.category == category
	}

	return false
}

// Msg returns this Kind with msg as its full error text and no cause.
func (k *Kind) Msg(msg string) error {
	return &wrappedError{
		kind: k,
		msg:  msg,
	}
}

// Msgf formats a message and calls [Kind.Msg].
func (k *Kind) Msgf(format string, args ...any) error {
	return k.Msg(fmt.Sprintf(format, args...))
}

// WithCause wraps cause with this Kind and uses cause.Error() as the error
// text. It returns nil when cause is nil.
func (k *Kind) WithCause(cause error) error {
	if cause == nil {
		return nil
	}

	return &wrappedError{
		kind:  k,
		cause: cause,
		msg:   cause.Error(),
	}
}

// Wrap wraps cause with this Kind and prefixes its error text with msg. It
// returns nil when cause is nil.
func (k *Kind) Wrap(cause error, msg string) error {
	if cause == nil {
		return nil
	}

	return &wrappedError{
		kind:  k,
		cause: cause,
		msg:   contextualMessage(msg, cause),
	}
}

// Wrapf formats msg and calls [Kind.Wrap]. It returns nil when cause is nil.
func (k *Kind) Wrapf(cause error, format string, args ...any) error {
	if cause == nil {
		return nil
	}

	return k.Wrap(cause, fmt.Sprintf(format, args...))
}

// TypeCode returns the Kind's application type code.
func (k *Kind) TypeCode() string {
	return k.typeCode
}

// HTTPCode returns the Kind's HTTP status.
func (k *Kind) HTTPCode() int {
	return k.httpCode
}

// GRPCCode returns the Kind's gRPC status code.
func (k *Kind) GRPCCode() codes.Code {
	return k.grpcCode
}

// PublicError returns the message set with [WithPublicMessage], or an empty
// string if none was set. Use [PublicMessage] to get a client message for any
// error.
func (k *Kind) PublicError() string {
	return k.publicMessage
}

func (k *Kind) classification() errorClassification {
	return errorClassification{
		kind:          k,
		category:      k.category,
		typeCode:      k.typeCode,
		httpCode:      k.httpCode,
		grpcCode:      k.grpcCode,
		publicMessage: k.publicMessage,
	}
}

// WithHTTPCode overrides the category's HTTP status for a new Kind. NewKind
// panics if the resulting status is outside 400–599.
func WithHTTPCode(httpCode int) KindOption {
	return func(k *Kind) {
		k.httpCode = httpCode
	}
}

// WithGRPCCode overrides the category's gRPC code for a new Kind. NewKind
// panics if the resulting code is OK.
func WithGRPCCode(grpcCode codes.Code) KindOption {
	return func(k *Kind) {
		k.grpcCode = grpcCode
	}
}

// WithPublicMessage sets the message that [PublicMessage] and [SendGRPCError]
// show clients for this Kind. Without it, clients see the HTTP status text
// for 4xx codes and "Internal Server Error" for 5xx codes. It panics if
// publicMessage is empty.
func WithPublicMessage(publicMessage string) KindOption {
	return func(k *Kind) {
		if publicMessage == "" {
			panic("errors: public message cannot be empty")
		}
		k.publicMessage = publicMessage
	}
}
