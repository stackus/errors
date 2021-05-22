![](https://github.com/stackus/errors/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/stackus/errors)](https://goreportcard.com/report/github.com/stackus/errors)
[![](https://godoc.org/github.com/stackus/errors?status.svg)](https://pkg.go.dev/github.com/stackus/errors)

# errors

Builds on Go 1.13 errors by adding HTTP statuses and GRPC codes to them.

## Installation

    go get -u github.com/stackus/errors

## Prerequisites

Go 1.13

## Adding HTTP status and GRPC codes to your errors

The `errors.Wrap(error, string) error` function is used to embed an `errors.Error` or to wrap other errors. When used
with an `errors.Error` the desired message is not altered. Wrapping other errors will prefix the message before the
wrapped error message.

    err := errors.Wrap(errors.ErrNotFound, "found nothing")
    fmt.Println(err) // Outputs: "found nothing"
    err = errors.Wrap(err, "a prefixed message")
    fmt.Println(err) // Outputs: "a prefixed message: found nothing"

## Displaying the embedded Error message

To display the embedded message you can use `errors.Message(error) string`. If the error is or has wrapped
an `errors.Error` then its type will be prefixed to the message.

    err := errors.Wrap(errors.ErrNotFound, "still nothing")
    fmt.Println(errors.Message(err)) // Outputs: "NOT_FOUND: still nothing"

No matter how many errors have been wrapped, the embedded `errors.Error` will be shown as the first prefix.

## Is(err, target) and As(err, target)

You can check for an embedded `error.Error` with `Is(err, target) bool`.

    err := errors.Wrap(errors.Wrap(errors.ErrNotFound, "nothing found"), "a prefix")
    if errors.Is(err, errors.ErrNotFound) {
        fmt.Println(errors.Message(err)) // Outputs: "NOT_FOUND: a prefix: nothing found"
    }

The methods `Is()`, `As()`, and `Unwrap()` from the standard `errors` package have all been made available in this
package as proxies for convenience.

## HTTP Statuses

The wrapped `errors.Error` can be checked with `errors.As()` and the `errors.HTTPCoder` interface to locate the HTTP
status.

    err := errors.Wrap(errors.ErrNotFound, "found nothing")
    var coder errors.HTTPCoder
    if errors.As(err, &coder) {
        fmt.Println(coder.HTTPCode()) // Outputs: 404
    }

## GRPC Codes

A similar method can be used to get the GRPC codes with the `errors.GRPCCoder` interface.

    err := errors.Wrap(errors.NotFound, "found nothing")
    var coder errors.GRPCCoder
    if errors.As(err, &coder) {
        fmt.Println(coder.GRPCCode()) // Outputs: 5
    }

## Transmitting errors with GRPC

The methods `SendGRPCError(error) error` and `ReceiveGRPCError(error) error` provide a way to convert a status.Status
and its error into an `errors.Error` and vice versa. You can use these in your server and client handlers directly, or
they can be used with GRPC interceptors.

Server Example:

    func serverErrorUnaryInterceptor() grpc.UnaryServerInterceptor {
	    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		    return resp, errors.SendGRPCError(err)
    	}
    }

    server := grpc.NewServer(grpc.ChainUnaryInterceptor(serverErrorUnaryInterceptor()), ...others)

Client Example:

    func clientErrorUnaryInterceptor() grpc.UnaryClientInterceptor {
	    return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		    return errors.ReceiveGRPCError(invoker(ctx, method, req, reply, cc, opts...))
    	}
    }

    cc, err := grpc.Dial(uri, grpc.WithChainUnaryInterceptor(clientErrorUnaryInterceptor()), ...others)

There is no requirement that both the server and client use this library to benefit from coded errors.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

MIT