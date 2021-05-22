// Package errors builds on Go 1.13 errors adding HTTP and GRPC code to your errors.
//
// Wrap() and Wrapf()
//
// When the wrap functions are used with one of the defined Err* constants you get back
// an error that you're able to pass the error through a GRPC server and client or
// use to build HTTP error messages and set the HTTP status.
//
// Wrapping any error other than an Error will return an error with the message formatted
// as "<message>: <error>".
//
// Wrapping an Error will return an error with an unaltered error message.
//
// Transmitting errors over GRPC
//
// The errors produced with wrap, that have also been wrapped first with an Err* can be
// send with SendGRPCError() and received with ReceiveGRPCError().
//
// You may want to create and use GRPC server and client interceptors to avoid having to
// call the Send/Receive methods in every handler.
//
// The Err* constants are errors and can be used directly is desired.
package errors
