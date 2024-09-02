![](https://github.com/stackus/errors/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/stackus/errors)](https://goreportcard.com/report/github.com/stackus/errors)
[![](https://godoc.org/github.com/stackus/errors?status.svg)](https://pkg.go.dev/github.com/stackus/errors)

# errors

Builds on Go 1.13 errors by adding HTTP statuses and GRPC codes to them.

## Installation

    go get -u github.com/stackus/errors

## Prerequisites

Go 1.13

## Embeddable codes

This library allows the use and helps facilitate the embedding of a type code, HTTP status, and GRPC code into errors
that can then be shared between services.

### Type Codes

Type codes are strings that are returned by any error that implements `errors.TypeCoder`.

    type TypeCoder interface {
        error
        TypeCode() string
    }

### HTTP Statuses

HTTP statuses are integer values that have defined in the `net/http` package and are returned by any error that
implements `errors.HTTPCoder`.

    type HTTPCoder interface {
        error
        HTTPCode() int
    }

### GRPC Codes

GRPC codes are `codes.Code` are int64 values defined in the `google.golang.org/grpc/codes` package and are returned by
any error that implements `errors.GRPCCoder`.

    type GRPCCoder interface {
        error
        GRPCCode() codes.Code
    }

### Packaged Error Types

The package also comes with many defined errors that are named in a way to reflect the GRPC code or HTTP status they
represent. The list of embeddable `errors.Error` types can be
found [here](https://github.com/stackus/errors/blob/master/types.go).

## Wrapping errors

The `errors.Wrap(error, string) error` function is used to wrap errors combining messages in most cases. However, when
the function is used with an error that has implemented `errors.TypeCoder` the message is not altered, and the error is
embedded instead.

```go
// Wrapping normal errors appends the error message
err := errors.Wrap(fmt.Errorf("sql error"), "error message")
fmt.Println(err) // Outputs: "error message: sql error"

// Wrapping errors.TypeCoder errors embeds the type
err := errors.Wrap(errors.ErrNotFound, "error message")
fmt.Println(err) // Outputs: "error message"

```

Wrapping multiple times will add additional prefixes to the error message.

```go
// Wrapping multiple times
err := errors.Wrap(errors.ErrNotFound, "error message")
err = errors.Wrap(err, "prefix")
err = errors.Wrap(err, "another")
fmt.Println(err) // Outputs: "another: prefix: error message"
```

### Wrapping using the errors.Err* errors

It is possible to use the package errors to wrap existing errors to add or override Type, HTTP code, or GRPC status codes.

```go
// Err will use the wrapped error .Error() output as the message
err := errors.ErrBadRequest.Err(fmt.Errorf("some error"))
// Msg and Msgf returns the Error with just the custom message applied
err = errors.ErrBadRequest.Msgf("%d total reasons", 7)
// Wrap and Wrapf will accept messages and simple wrap the error
err = errors.ErrUnauthorized.Wrap(err, "some message")
```

Both errors can be checked for using the `Is()` and `As()` methods when you wrap errors with the package errors this way.

## Getting type, HTTP status, or GRPC code

The Go 1.13 `errors.As(error, interface{}) bool` function from the standard `errors` package can be used to turn an
error into any of the three "Coder" interfaces documented above.

    err := errors.Wrap(errors.NotFound, "error message")
    var coder errors.TypeCoder
    if errors.As(err, &coder) {
        fmt.Println(coder.TypeCode()) // Outputs: "NOT_FOUND"
    }

> The functions `Is()`, `As()`, and `Unwrap()` from the standard `errors` package have all been made available in this package as proxies for convenience.

The functions `errors.TypeCode(error) string`, `errors.HTTPCode(error) int`, and `errors.GRPCCode(error) codes.Code` can
be used to fetch specific code. They're more convenient to use than the interfaces directly. The catch is they have
defined rules for the values they return.

#### errors.TypeCode(error) string

If the error implements or has wrapped an error that implements `errors.TypeCoder` it will return the code from that
error. If no error is found to support the interface then the string `"UNKNOWN"` is returned. Nil errors result in a
blank string being returned.

    fmt.Println(errors.TypeCode(errors.ErrNotFound)) // Outputs: "NOT_FOUND"
    fmt.Println(errors.TypeCode(fmt.Errorf("an error"))) // Outputs: "UNKNOWN"
    fmt.Println(errors.TypeCode(nil)) // Outputs: ""

#### errors.HTTPCode(error) int

If the error implements or has wrapped an error that implements `errors.HTTPCoder` it will return the status from that
error. If no error is found to support the interface then `http.StatusNotExtended` is returned. Nil errors result
in `http.StatusOK` being returned.

    fmt.Println(errors.HTTPCode(errors.ErrNotFound)) // Outputs: 404
    fmt.Println(errors.HTTPCode(fmt.Errorf("an error"))) // Outputs: 510
    fmt.Println(errors.HTTPCode(nil)) // Outputs: 200

#### errors.GRPCCode(error) codes.Code

If the error implements or has wrapped an error that implements `errors.GRPCCoder` it will return the code from that
error. If no error is found to support the interface then `codes.Unknown` is returned. Nil errors result in `codes.OK`
being returned.

    fmt.Println(errors.GRPCCode(errors.ErrNotFound)) // Outputs: 5
    fmt.Println(errors.GRPCCode(fmt.Errorf("an error"))) // Outputs: 2
    fmt.Println(errors.GRPCCode(nil)) // Outputs: 0

#### Why Unknown? Why not default to internal errors?

Part of the reason you'd want to use a library that adds code to your errors is because you want to better identify the
problems in your application. By marking un-coded errors as "Unknown" errors they'll stand out from any errors you've
marked as `codes.Internal` for example.

## Transmitting errors with GRPC

The functions `SendGRPCError(error) error` and `ReceiveGRPCError(error) error` provide a way to convert
a `status.Status` and its error into an error that provides codes and vice versa. You can use these in your server and
client handlers directly, or they can be used with GRPC interceptors.

Server Interceptor Example:

    // Unary only example
    func serverErrorUnaryInterceptor() grpc.UnaryServerInterceptor {
	    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		    return resp, errors.SendGRPCError(err)
    	}
    }

    server := grpc.NewServer(grpc.ChainUnaryInterceptor(serverErrorUnaryInterceptor()), ...others)

Client Interceptor Example:

    // Unary only example
    func clientErrorUnaryInterceptor() grpc.UnaryClientInterceptor {
	    return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		    return errors.ReceiveGRPCError(invoker(ctx, method, req, reply, cc, opts...))
    	}
    }

    cc, err := grpc.Dial(uri, grpc.WithChainUnaryInterceptor(clientErrorUnaryInterceptor()), ...others)

### Comparing received errors

Servers and clients may not always use a shared library when exchanging errors. In fact there isn't any requirement that
the server and client both use this library to exchange errors.

When comparing received errors with `errors.Is(error, error) bool` the checks are a little more loose. A received error
is considered to be the same if **ANY** of the codes are a match. This differs from a strict equality check for the
server before the error was sent.

The "Code" functions and the "Coder" interfaces continue to work the same on a client as they did on the server that
sent the error.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

MIT
