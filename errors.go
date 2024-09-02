package errors

import (
	stderrors "errors"
	"fmt"
)

// TypeCoder interface to extract an errors embeddable type as a string
type TypeCoder interface {
	error
	TypeCode() string
}

// Error base error
type Error string

// Error implements error
func (e Error) Error() string {
	if e == ErrOK {
		return ""
	}
	return string(e)
}

func (e Error) TypeCode() string {
	if e == ErrOK {
		return ""
	}
	return string(e)
}

// Err overrides or adds Type,HTTP,GRPC information for the passed in error
// while leaving Is() and As() functionality unchanged
func (e Error) Err(err error) error {
	if err == nil {
		return nil
	}
	return embeddedError{te: e, e: err, msg: err.Error()}
}

// Msg sets a custom message for the Error
func (e Error) Msg(msg string) error {
	return embeddedError{e: e, msg: msg}
}

// Msgf sets a custom message for formatting for the Error
func (e Error) Msgf(format string, args ...interface{}) error {
	return embeddedError{e: e, msg: fmt.Sprintf(format, args...)}
}

// Wrap an error with message while overriding or adding Type,HTTP,GRPC information
// while leaving Is() and As() functionality unchanged
func (e Error) Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return embeddedError{te: e, e: err, msg: msg}
}

// Wrapf an error with message while overriding or adding Type,HTTP,GRPC information
// while leaving Is() and As() functionality unchanged
func (e Error) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return embeddedError{te: e, e: err, msg: fmt.Sprintf(format, args...)}
}

type embeddedError struct {
	e   error  // original error to be embedded
	te  error  // overriding error type
	msg string // for the humans
}

func (e embeddedError) Error() string {
	return e.msg
}

func (e embeddedError) TypeCode() string {
	var typeCoder TypeCoder
	if e.te != nil && stderrors.As(e.te, &typeCoder) {
		return typeCoder.TypeCode()
	}
	if e.e != nil && stderrors.As(e.e, &typeCoder) {
		return typeCoder.TypeCode()
	}
	return ErrUnknown.TypeCode()
}

func (e embeddedError) Is(target error) bool {
	if e.te != nil && stderrors.Is(e.te, target) {
		return true
	}
	if e.e != nil && stderrors.Is(e.e, target) {
		return true
	}
	return false
}

func (e embeddedError) As(target interface{}) bool {
	if e.te != nil && stderrors.As(e.te, target) {
		return true
	}
	if e.e != nil && stderrors.As(e.e, target) {
		return true
	}
	return false
}

// Wrap returns an error with msg wrapped with the supplied error
// If err is nil then Wrap returns nil
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	switch err.(type) {
	case embeddedError:
		return embeddedError{e: err, msg: fmt.Sprintf("%s: %s", msg, err.Error())}
	case TypeCoder:
		return embeddedError{te: err, msg: msg}
	default:
		return embeddedError{e: err, te: ErrInternalServerError, msg: fmt.Sprintf("%s: %s", msg, err.Error())}
	}
}

// Wrapf returns an error with a formatted msg wrapped with the supplied error
// If err is nil then Wrapf returns nil
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	switch err.(type) {
	case embeddedError:
		return embeddedError{e: err, msg: fmt.Sprintf("%s: %s", fmt.Sprintf(format, args...), err.Error())}
	case TypeCoder:
		return embeddedError{te: err, msg: fmt.Sprintf(format, args...)}
	default:
		return embeddedError{
			e:   err,
			te:  ErrInternalServerError,
			msg: fmt.Sprintf("%s: %s", fmt.Sprintf(format, args...), err.Error()),
		}
	}
}

// TypeCode returns the embedded type for the given error or blank when nil or UNKNOWN otherwise
func TypeCode(err error) string {
	if err == nil {
		return ErrOK.TypeCode()
	}

	var e TypeCoder
	if stderrors.As(err, &e) {
		return e.TypeCode()
	}
	return ErrUnknown.TypeCode()
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
