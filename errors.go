package errors

import (
	stderrors "errors"
	"fmt"
	"io"
)

// Error base error
type Error string

// Error implements error
func (e Error) Error() string {
	return string(e)
}

// Format implements fmt.Formatter
//
// Error is embedded without modification to the error message
// with Wrap() and Wrapf() or manually with fmt.Errorf() using "%-w".
func (e Error) Format(s fmt.State, v rune) {
	switch v {
	case 's':
		_, _ = io.WriteString(s, string(e))
	case 'v':
		if !s.Flag('-') {
			_, _ = io.WriteString(s, string(e))
		}
	}
}

// Wrap returns an error with msg wrapped with the supplied error
// If err is nil then Wrap returns nil
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(Error); ok {
		return embed(e, msg)
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// Wrapf returns an error with a formatted msg wrapped with the supplied error
// If err is nil then Wrapf returns nil
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(Error); ok {
		return embedf(e, format, args...)
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// Message displays the string value for Error prefixed to the existing error message
//
// Prefixed as "Error: message"
//
// If err is not an Error or it hasn't wrapped one then there will no modifications made
// to the message.
func Message(err error) string {
	var e Error
	if stderrors.As(err, &e) {
		return fmt.Sprintf("%s: %s", e.Error(), err.Error())
	}
	return err.Error()
}

// embed returns an error with msg wrapped with the supplied Error without suffixing it
// If err is nil then embed returns nil
func embed(err Error, msg string) error {
	return fmt.Errorf("%-w%s", err, msg)
}

// embedf returns an error with a formatted msg wrapped with the supplied Error without suffixing it
// If err is nil then embedf returns nil
func embedf(err Error, format string, args ...interface{}) error {
	return fmt.Errorf("%-w%s", err, fmt.Sprintf(format, args...))
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
