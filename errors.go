package errors

import (
	stderrors "errors"
	"fmt"
)

// TypeCoder interface to extract an errors embeddable type as a string
type TypeCoder interface {
	TypeCode() string
}

// Error base error
type Error string

// Error implements error
func (e Error) Error() string {
	return string(e)
}

func (e Error) TypeCode() string {
	return string(e)
}

type embeddedError struct {
	e   error
	msg string
}

func (e embeddedError) Error() string {
	return e.msg
}

func (e embeddedError) Is(target error) bool {
	return stderrors.Is(e.e, target)
}

func (e embeddedError) As(target interface{}) bool {
	return stderrors.As(e.e, target)
}

// Wrap returns an error with msg wrapped with the supplied error
// If err is nil then Wrap returns nil
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(TypeCoder); ok {
		return embeddedError{err, msg}
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// Wrapf returns an error with a formatted msg wrapped with the supplied error
// If err is nil then Wrapf returns nil
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(TypeCoder); ok {
		return embeddedError{err, fmt.Sprintf(format, args...)}
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// TypeCode returns the embedded type for the given error or blank when nil or UNKNOWN otherwise
func TypeCode(err error) string {
	if err == nil {
		return ""
	}

	var e TypeCoder
	if stderrors.As(err, &e) {
		return e.TypeCode()
	}
	return "UNKNOWN"
}

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
