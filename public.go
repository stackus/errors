package errors

import (
	"net/http"
)

type (
	publicError struct {
		cause   error
		message string
	}

	// PublicError is the type code and client message for an error, as
	// returned by [Public]. It serializes without the error's internal text,
	// for example as JSON:
	//
	//	{"code":"USER_NOT_FOUND","message":"The user could not be found"}
	PublicError struct {
		Code    string `json:"code" xml:"code" yaml:"code" msgpack:"code"`
		Message string `json:"message" xml:"message" yaml:"message" msgpack:"message"`
	}
)

func (e publicError) Error() string {
	return e.cause.Error()
}

func (e publicError) Unwrap() error {
	return e.cause
}

// Public returns err's type code and client-facing message. For nil, the code
// is "OK" and the message is empty.
func Public(err error) PublicError {
	return PublicError{
		Code:    TypeCode(err),
		Message: PublicMessage(err),
	}
}

// WrapPublicMessage sets the message that [PublicMessage] returns for err,
// for cases where one error needs a message its kind or category doesn't
// provide. err's text, codes, and identity don't change. If err is later
// wrapped by a category or kind, that outer classification's message is used
// instead. WrapPublicMessage returns nil for a nil err and panics when
// message is empty for a non-nil err.
func WrapPublicMessage(err error, message string) error {
	if err == nil {
		return nil
	}

	if message == "" {
		panic("errors: public message cannot be empty")
	}

	return &publicError{
		cause:   err,
		message: message,
	}
}

// PublicMessage returns a message for err that is safe to show clients. It
// never returns err.Error(). It uses the first of these that applies:
//
//  1. A message set with [WrapPublicMessage].
//  2. The message set with [WithPublicMessage] on err's [Kind].
//  3. "Internal Server Error" for a 5xx HTTP code or an unclassified error.
//  4. The HTTP status text for err's code, such as "Not Found".
//
// It returns an empty string for nil.
func PublicMessage(err error) string {
	if err == nil {
		return ""
	}

	if msg, ok := publicMessageOverride(err); ok {
		return msg
	}

	c, ok := resolveClassification(err)
	if !ok {
		return "Internal Server Error"
	}

	if c.publicMessage != "" {
		return c.publicMessage
	}

	if c.httpCode >= 500 {
		return "Internal Server Error"
	}

	msg := http.StatusText(c.httpCode)
	if msg == "" {
		return "Request failed"
	}

	return msg
}

func publicMessageOverride(err error) (string, bool) {
	for err != nil {
		if e, ok := err.(*publicError); ok {
			return e.message, true
		}

		if _, ok := err.(classified); ok {
			return "", false
		}

		if _, ok := err.(multiUnwrapper); ok {
			return "", false
		}

		wrapped, ok := err.(singleUnwrapper)

		if !ok {
			return "", false
		}

		err = wrapped.Unwrap()
	}

	return "", false
}
