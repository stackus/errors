package errors

import (
	"fmt"
	"regexp"

	"google.golang.org/grpc/codes"
)

type (
	Kind struct {
		id *struct{ _ byte } // non-zero identity ensures two independent kinds are never equal

		category Error
		typeCode string
		httpCode int
		grpcCode codes.Code

		message       string
		publicMessage string
	}

	KindOption func(*Kind)
)

var typeCodeRe = regexp.MustCompile(`^[A-Z][A-Z0-9_.]*$`)

func NewKind(typeCode string, category Error, message string, options ...KindOption) *Kind {
	if !typeCodeRe.MatchString(typeCode) {
		panic("errors: invalid type code: " + typeCode)
	}

	if _, ok := errorCategories[category]; !ok {
		panic("errors: invalid category: " + string(category))
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

func WithKind(err error, kind *Kind) error {
	if err == nil {
		return nil
	}

	if kind == nil {
		panic("errors: kind cannot be nil")
	}

	return kind.WithCause(err)
}

func (k *Kind) Error() string {
	return k.message
}

func (k *Kind) Is(target error) bool {
	if t, ok := target.(*Kind); ok && k == t {
		return k == t
	}

	if category, ok := target.(Error); ok {
		return k.category == category
	}

	return false
}

func (k *Kind) Msg(msg string) error {
	return &wrappedError{
		kind: k,
		msg:  msg,
	}
}

func (k *Kind) Msgf(format string, args ...any) error {
	return k.Msg(fmt.Sprintf(format, args...))
}

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

func (k *Kind) Wrapf(cause error, format string, args ...any) error {
	if cause == nil {
		return nil
	}

	return k.Wrap(cause, fmt.Sprintf(format, args...))
}

func (k *Kind) TypeCode() string {
	return k.typeCode
}

func (k *Kind) HTTPCode() int {
	return k.httpCode
}

func (k *Kind) GRPCCode() codes.Code {
	return k.grpcCode
}

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

// --- KindOptions

func WithHTTPCode(httpCode int) KindOption {
	return func(k *Kind) {
		k.httpCode = httpCode
	}
}

func WithGRPCCode(grpcCode codes.Code) KindOption {
	return func(k *Kind) {
		k.grpcCode = grpcCode
	}
}

func WithPublicMessage(publicMessage string) KindOption {
	return func(k *Kind) {
		if publicMessage == "" {
			panic("errors: public message cannot be empty")
		}
		k.publicMessage = publicMessage
	}
}
