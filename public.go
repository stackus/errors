package errors

import (
	"net/http"
)

type (
	publicError struct {
		cause   error
		message string
	}

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

func Public(err error) PublicError {
	return PublicError{
		Code:    TypeCode(err),
		Message: PublicMessage(err),
	}
}

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
