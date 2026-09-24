package errors_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"
	"google.golang.org/grpc/codes"
)

type typeTestError struct {
	t string
	e error
}

func (e typeTestError) Error() string      { return e.t }
func (e typeTestError) TypeCode() string   { return e.t }
func (e typeTestError) Is(err error) bool  { return stderrors.Is(e.e, err) }
func (e typeTestError) As(target any) bool { return stderrors.As(e.e, target) }

type httpTestError struct {
	hc int
	e  error
}

func (e httpTestError) Error() string              { return fmt.Sprintf("HTTP(%d)", e.hc) }
func (e httpTestError) HTTPCode() int              { return e.hc }
func (e httpTestError) Is(err error) bool          { return stderrors.Is(e.e, err) }
func (e httpTestError) As(target interface{}) bool { return stderrors.As(e.e, target) }

type grpcTestError struct {
	gc codes.Code
	e  error
}

func (e grpcTestError) Error() string              { return fmt.Sprintf("GRPC(%d)", e.gc) }
func (e grpcTestError) GRPCCode() codes.Code       { return e.gc }
func (e grpcTestError) Is(err error) bool          { return stderrors.Is(e.e, err) }
func (e grpcTestError) As(target interface{}) bool { return stderrors.As(e.e, target) }

func convertGRPCStringToCode(grpcCode string) codes.Code {
	switch grpcCode {
	case "codes.OK":
		return codes.OK
	case "codes.Canceled":
		return codes.Canceled
	case "codes.Unknown":
		return codes.Unknown
	case "codes.InvalidArgument":
		return codes.InvalidArgument
	case "codes.DeadlineExceeded":
		return codes.DeadlineExceeded
	case "codes.NotFound":
		return codes.NotFound
	case "codes.AlreadyExists":
		return codes.AlreadyExists
	case "codes.PermissionDenied":
		return codes.PermissionDenied
	case "codes.ResourceExhausted":
		return codes.ResourceExhausted
	case "codes.FailedPrecondition":
		return codes.FailedPrecondition
	case "codes.Aborted":
		return codes.Aborted
	case "codes.OutOfRange":
		return codes.OutOfRange
	case "codes.Unimplemented":
		return codes.Unimplemented
	case "codes.Internal":
		return codes.Internal
	case "codes.Unavailable":
		return codes.Unavailable
	case "codes.DataLoss":
		return codes.DataLoss
	case "codes.Unauthenticated":
		return codes.Unauthenticated
	default:
		return codes.Unknown
	}
}

func convertHTTPStringToInt(httpStatus string) int {
	switch httpStatus {
	case "http.StatusOK":
		return http.StatusOK
	case "http.StatusBadRequest":
		return http.StatusBadRequest
	case "http.StatusUnauthorized":
		return http.StatusUnauthorized
	case "http.StatusPaymentRequired":
		return http.StatusPaymentRequired
	case "http.StatusForbidden":
		return http.StatusForbidden
	case "http.StatusNotFound":
		return http.StatusNotFound
	case "http.StatusMethodNotAllowed":
		return http.StatusMethodNotAllowed
	case "http.StatusNotAcceptable":
		return http.StatusNotAcceptable
	case "http.StatusProxyAuthRequired":
		return http.StatusProxyAuthRequired
	case "http.StatusRequestTimeout":
		return http.StatusRequestTimeout
	case "http.StatusConflict":
		return http.StatusConflict
	case "http.StatusGone":
		return http.StatusGone
	case "http.StatusLengthRequired":
		return http.StatusLengthRequired
	case "http.StatusPreconditionFailed":
		return http.StatusPreconditionFailed
	case "http.StatusRequestEntityTooLarge":
		return http.StatusRequestEntityTooLarge
	case "http.StatusRequestURITooLong":
		return http.StatusRequestURITooLong
	case "http.StatusUnsupportedMediaType":
		return http.StatusUnsupportedMediaType
	case "http.StatusRequestedRangeNotSatisfiable":
		return http.StatusRequestedRangeNotSatisfiable
	case "http.StatusExpectationFailed":
		return http.StatusExpectationFailed
	case "http.StatusTeapot":
		return http.StatusTeapot
	case "http.StatusMisdirectedRequest":
		return http.StatusMisdirectedRequest
	case "http.StatusUnprocessableEntity":
		return http.StatusUnprocessableEntity
	case "http.StatusLocked":
		return http.StatusLocked
	case "http.StatusFailedDependency":
		return http.StatusFailedDependency
	case "http.StatusTooEarly":
		return http.StatusTooEarly
	case "http.StatusUpgradeRequired":
		return http.StatusUpgradeRequired
	case "http.StatusPreconditionRequired":
		return http.StatusPreconditionRequired
	case "http.StatusTooManyRequests":
		return http.StatusTooManyRequests
	case "http.StatusRequestHeaderFieldsTooLarge":
		return http.StatusRequestHeaderFieldsTooLarge
	case "http.StatusUnavailableForLegalReasons":
		return http.StatusUnavailableForLegalReasons
	case "http.StatusInternalServerError":
		return http.StatusInternalServerError
	case "http.StatusNotImplemented":
		return http.StatusNotImplemented
	case "http.StatusBadGateway":
		return http.StatusBadGateway
	case "http.StatusServiceUnavailable":
		return http.StatusServiceUnavailable
	case "http.StatusGatewayTimeout":
		return http.StatusGatewayTimeout
	case "http.StatusHTTPVersionNotSupported":
		return http.StatusHTTPVersionNotSupported
	case "http.StatusVariantAlsoNegotiates":
		return http.StatusVariantAlsoNegotiates
	case "http.StatusInsufficientStorage":
		return http.StatusInsufficientStorage
	case "http.StatusLoopDetected":
		return http.StatusLoopDetected
	case "http.StatusNotExtended":
		return http.StatusNotExtended
	case "http.StatusNetworkAuthenticationRequired":
		return http.StatusNetworkAuthenticationRequired
	default:
		return http.StatusNotExtended
	}
}

func convertErrNameToError(errName string) errors.Error {
	switch errName {
	case "ErrOK":
		return errors.ErrOK
	case "ErrCanceled":
		return errors.ErrCanceled
	case "ErrUnknown":
		return errors.ErrUnknown
	case "ErrInvalidArgument":
		return errors.ErrInvalidArgument
	case "ErrDeadlineExceeded":
		return errors.ErrDeadlineExceeded
	case "ErrNotFound":
		return errors.ErrNotFound
	case "ErrAlreadyExists":
		return errors.ErrAlreadyExists
	case "ErrPermissionDenied":
		return errors.ErrPermissionDenied
	case "ErrResourceExhausted":
		return errors.ErrResourceExhausted
	case "ErrFailedPrecondition":
		return errors.ErrFailedPrecondition
	case "ErrAborted":
		return errors.ErrAborted
	case "ErrOutOfRange":
		return errors.ErrOutOfRange
	case "ErrUnimplemented":
		return errors.ErrUnimplemented
	case "ErrInternal":
		return errors.ErrInternal
	case "ErrUnavailable":
		return errors.ErrUnavailable
	case "ErrDataLoss":
		return errors.ErrDataLoss
	case "ErrUnauthenticated":
		return errors.ErrUnauthenticated
	case "ErrBadRequest":
		return errors.ErrBadRequest
	case "ErrUnauthorized":
		return errors.ErrUnauthorized
	case "ErrForbidden":
		return errors.ErrForbidden
	case "ErrMethodNotAllowed":
		return errors.ErrMethodNotAllowed
	case "ErrRequestTimeout":
		return errors.ErrRequestTimeout
	case "ErrConflict":
		return errors.ErrConflict
	case "ErrImATeapot":
		return errors.ErrImATeapot
	case "ErrUnprocessableEntity":
		return errors.ErrUnprocessableEntity
	case "ErrTooManyRequests":
		return errors.ErrTooManyRequests
	case "ErrUnavailableForLegalReasons":
		return errors.ErrUnavailableForLegalReasons
	case "ErrInternalServerError":
		return errors.ErrInternalServerError
	case "ErrNotImplemented":
		return errors.ErrNotImplemented
	case "ErrBadGateway":
		return errors.ErrBadGateway
	case "ErrServiceUnavailable":
		return errors.ErrServiceUnavailable
	case "ErrGatewayTimeout":
		return errors.ErrGatewayTimeout
	default:
		return errors.ErrUnknown
	}
}

func TestMain(m *testing.M) {
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	opts := godog.Options{
		Format:   format,
		Paths:    []string{"features"},
		NoColors: true,
	}

	status := godog.TestSuite{
		Name:                 "errors",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

// var expectedMessage string
var expectedError error

func theErrorDoesNotImplementCoderInterfaces() error {
	expectedError = fmt.Errorf("%s", expectedError)
	return nil
}

func theErrorIs(errName string) error {
	err := convertErrNameToError(errName)
	expectedError = err
	return nil
}

func theErrorIsNil() error {
	expectedError = nil
	return nil
}

func theErrorSentByAGRPCServerWas(errName string) error {
	err := convertErrNameToError(errName)
	grpcErr := errors.SendGRPCError(err)
	grpcErr = errors.ReceiveGRPCError(grpcErr)
	expectedError = grpcErr
	return nil
}

func anErrorWithTheMessage(message string) error {
	expectedError = errors.Wrap(expectedError, message)
	return nil
}

func anErrorWithTypeCode(typeCode string) error {
	expectedError = typeTestError{t: typeCode, e: expectedError}
	return nil
}

func anErrorWithHTTPStatus(httpStatus string) error {
	expectedError = httpTestError{hc: convertHTTPStringToInt(httpStatus), e: expectedError}
	return nil
}

func anErrorWithGRPCCode(grpcCode string) error {
	expectedError = grpcTestError{gc: convertGRPCStringToCode(grpcCode), e: expectedError}
	return nil
}

func wrappedWithTheMessage(message string) error {
	expectedError = errors.Wrap(expectedError, message)
	return nil
}

func wrappedWithTheError(errName, message string) error {
	err := convertErrNameToError(errName)
	if origerr, ok := expectedError.(errors.Error); ok {
		expectedError = origerr.Wrap(err, message)
	}
	return nil
}

func theErrorIsSentOverGRPC() error {
	grpcErr := errors.SendGRPCError(expectedError)
	grpcErr = errors.ReceiveGRPCError(grpcErr)
	expectedError = grpcErr
	return nil
}

func theErrorMessageIs(message string) error {
	if expectedError.Error() != message {
		return fmt.Errorf("expected message to be `%s` but got `%s`", message, expectedError.Error())
	}
	return nil
}

func theHTTPStatusIs(httpStatus string) error {
	got := http.StatusText(errors.HTTPCode(expectedError))
	if got != httpStatus {
		return fmt.Errorf("expected HTTP status to be `%s` but got `%s`", httpStatus, got)
	}
	return nil
}

func theTypeCodeIs(typeCode string) error {
	got := errors.TypeCode(expectedError)
	if got != typeCode {
		return fmt.Errorf("expected type code to be `%s` but got `%s`", typeCode, got)
	}
	return nil
}

func theGRPCCodeIs(grpcCode string) error {
	got := errors.GRPCCode(expectedError).String()
	if got != grpcCode {
		return fmt.Errorf("expected GRPC code to be `%s` but got `%s`", grpcCode, got)
	}
	return nil
}

func theErrorIsA(errName string) error {
	if !errors.Is(expectedError, convertErrNameToError(errName)) {
		return fmt.Errorf("expected error to be a `%s`", errName)
	}
	return nil
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		expectedError = errors.ErrUnknown
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		expectedError = stderrors.New("test error")
		return ctx, nil
	})

	// Given
	ctx.Step(`^the error does not implement (:?Type|HTTP|GRPC)Coder{}$`, theErrorDoesNotImplementCoderInterfaces)
	ctx.Step(`^the error is "([^"]*)"$`, theErrorIs)
	ctx.Step(`^the error is nil$`, theErrorIsNil)
	ctx.Step(`^the error sent by a GRPC server was "([^"]*)"$`, theErrorSentByAGRPCServerWas)
	ctx.Step(`^an error with the message "([^"]*)"$`, anErrorWithTheMessage)
	ctx.Step(`^an error with Type code "([^"]*)"$`, anErrorWithTypeCode)
	ctx.Step(`^an error with HTTP status "([^"]*)"$`, anErrorWithHTTPStatus)
	ctx.Step(`^an error with GRPC code "([^"]*)"$`, anErrorWithGRPCCode)

	// When
	ctx.Step(`^wrapped with the error "([^"]*)" and message "([^"]*)"$`, wrappedWithTheError)
	ctx.Step(`^wrapped with the message "([^"]*)"$`, wrappedWithTheMessage)
	ctx.Step(`^the error is sent over GRPC$`, theErrorIsSentOverGRPC)

	// Then
	ctx.Step(`^the Type code is "([^"]*)"$`, theTypeCodeIs)
	ctx.Step(`^the HTTP status is "([^"]*)"$`, theHTTPStatusIs)
	ctx.Step(`^the GRPC code is "([^"]*)"$`, theGRPCCodeIs)
	ctx.Step(`^the error message is "([^"]*)"$`, theErrorMessageIs)
	ctx.Step(`^the error is a "([^"]*)"$`, theErrorIsA)
}
