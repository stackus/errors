package errors

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

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
	grpcErr := SendGRPCError(err)
	grpcErr = ReceiveGRPCError(grpcErr)
	expectedError = grpcErr
	return nil
}

func anErrorWithTheMessage(message string) error {
	expectedError = Wrap(expectedError, message)
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
	expectedError = Wrap(expectedError, message)
	return nil
}

func theErrorIsSentOverGRPC() error {
	grpcErr := SendGRPCError(expectedError)
	grpcErr = ReceiveGRPCError(grpcErr)
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
	got := http.StatusText(HTTPCode(expectedError))
	if got != httpStatus {
		return fmt.Errorf("expected HTTP status to be `%s` but got `%s`", httpStatus, got)
	}
	return nil
}

func theTypeCodeIs(typeCode string) error {
	got := TypeCode(expectedError)
	if got != typeCode {
		return fmt.Errorf("expected type code to be `%s` but got `%s`", typeCode, got)
	}
	return nil
}

func theGRPCCodeIs(grpcCode string) error {
	got := GRPCCode(expectedError).String()
	if got != grpcCode {
		return fmt.Errorf("expected GRPC code to be `%s` but got `%s`", grpcCode, got)
	}
	return nil
}

func theErrorIsA(errName string) error {
	if !Is(expectedError, convertErrNameToError(errName)) {
		return fmt.Errorf("expected error to be a `%s`", errName)
	}
	return nil
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		expectedError = ErrUnknown
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
		expectedError = ErrUnknown
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
	ctx.Step(`^wrapped with the message "([^"]*)"$`, wrappedWithTheMessage)
	ctx.Step(`^the error is sent over GRPC$`, theErrorIsSentOverGRPC)

	// Then
	ctx.Step(`^the Type code is "([^"]*)"$`, theTypeCodeIs)
	ctx.Step(`^the HTTP status is "([^"]*)"$`, theHTTPStatusIs)
	ctx.Step(`^the GRPC code is "([^"]*)"$`, theGRPCCodeIs)
	ctx.Step(`^the error message is "([^"]*)"$`, theErrorMessageIs)
	ctx.Step(`^the error is a "([^"]*)"$`, theErrorIsA)
}
