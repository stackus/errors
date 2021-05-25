package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

type customTestError struct {
	gc codes.Code
	hc int
	t  string
}

func (e customTestError) Error() string        { return e.t }
func (e customTestError) GRPCCode() codes.Code { return e.gc }
func (e customTestError) HTTPCode() int        { return e.hc }
func (e customTestError) TypeCode() string     { return e.t }

type embedTestError struct {
	t string
}

func (e embedTestError) Error() string    { return e.t }
func (e embedTestError) TypeCode() string { return e.t }

type httpTestError struct {
	hc int
}

func (e httpTestError) Error() string { return fmt.Sprintf("HTTP(%d)", e.hc) }
func (e httpTestError) HTTPCode() int { return e.hc }

type grpcTestError struct {
	gc codes.Code
}

func (e grpcTestError) Error() string        { return fmt.Sprintf("GRPC(%d)", e.gc) }
func (e grpcTestError) GRPCCode() codes.Code { return e.gc }
