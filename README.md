![](https://github.com/stackus/errors/workflows/CI/badge.svg)
[![](https://godoc.org/github.com/stackus/errors?status.svg)](https://pkg.go.dev/github.com/stackus/errors)

# errors

Give each Go error a **type code**, an **HTTP status**, and a **gRPC code**, and
keep them as the error is wrapped, returned, and sent between services.

```go
var ErrUserNotFound = errors.NewKind("USER_NOT_FOUND", errors.ErrNotFound, "user not found")

err := ErrUserNotFound.Wrap(sql.ErrNoRows, "find user 42")

errors.TypeCode(err)                // "USER_NOT_FOUND"
errors.HTTPCode(err)                // 404
errors.GRPCCode(err)                // codes.NotFound
errors.Is(err, ErrUserNotFound)     // true
errors.Is(err, errors.ErrNotFound)  // true
errors.Is(err, sql.ErrNoRows)       // true
```

Without this, services usually map errors to status codes in each handler.
Those mappings drift apart and lose detail at every gRPC hop. With this
package, you classify an error once, where it happens:

- **One classification, every transport.** HTTP handlers call `HTTPCode(err)`,
  gRPC servers return `SendGRPCError(err)`, and logs record `TypeCode(err)`.
- **Your own sentinel errors.** `NewKind` creates errors like
  `USER_NOT_FOUND` that you check with `errors.Is`. They also match a broader
  category, such as `ErrNotFound`, so generic code can handle them.
- **Classifications cross gRPC.** A client gets an error that still matches
  the server's `Kind`, with the same type code, HTTP status, and gRPC code. A
  gateway can forward it to HTTP without its own mapping table.
- **Safe messages for clients.** `PublicMessage` and `Public` give clients a
  message without internal details such as SQL errors or file paths.
- **Works with the standard library.** Errors work with `errors.Is`,
  `errors.As`, and `fmt.Errorf("%w")`. The package re-exports `Is`, `As`,
  `AsType`, `Unwrap`, `Join`, and `New`, so it can replace the standard import.

## Installation

```sh
go get github.com/stackus/errors
```

Requires Go 1.26.8 or later.

## Contents

- [Quick start](#quick-start)
- [Defining your errors](#defining-your-errors)
- [Classifying errors and adding context](#classifying-errors-and-adding-context)
- [Checking errors](#checking-errors)
- [Reading codes](#reading-codes)
- [Public messages](#public-messages)
- [HTTP responses](#http-responses)
- [gRPC](#grpc)
- [Using your own error types](#using-your-own-error-types)
- [Built-in categories](#built-in-categories)

## Quick start

An error is classified where it enters the application. Code above it adds
context, and the caller reads the codes and checks what the error is:

```go
var ErrUserNotFound = errors.NewKind(
    "USER_NOT_FOUND",
    errors.ErrNotFound,
    "user not found",
    errors.WithPublicMessage("The user could not be found"),
)

// Repository: classify the driver error.
func (r *Repo) FindUser(ctx context.Context, id string) (*User, error) {
    row := r.db.QueryRowContext(ctx, "SELECT ... WHERE id = $1", id)
    if err := row.Scan(...); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound.Wrap(err, "find user "+id)
        }
        return nil, errors.ErrInternal.Wrap(err, "find user "+id)
    }
    ...
}

// Service: add context. The classification is preserved.
func (s *Service) GetProfile(ctx context.Context, id string) (*Profile, error) {
    user, err := s.repo.FindUser(ctx, id)
    if err != nil {
        return nil, errors.Wrap(err, "get profile")
    }
    ...
}
```

```go
err := svc.GetProfile(ctx, "42")

fmt.Println(err)                                // get profile: find user 42: sql: no rows in result set
fmt.Println(errors.TypeCode(err))               // USER_NOT_FOUND
fmt.Println(errors.HTTPCode(err))               // 404
fmt.Println(errors.GRPCCode(err))               // NotFound
fmt.Println(errors.Is(err, ErrUserNotFound))    // true
fmt.Println(errors.Is(err, errors.ErrNotFound)) // true
fmt.Println(errors.PublicMessage(err))          // The user could not be found
```

## Defining your errors

### Built-in categories

The package includes a category for every gRPC code and for common HTTP
statuses. Examples are `ErrNotFound`, `ErrInvalidArgument`, `ErrConflict`, and
`ErrUnauthorized`. [See the full table](#built-in-categories). A category is
already a sentinel error:

```go
if user == nil {
    return errors.ErrNotFound
}
```

Categories are enough for most generic failures, such as "not found",
"invalid input", or "internal".

### Kinds: your own sentinel errors

Create a `Kind` for an error that callers need to tell apart from others, or
that should have its own code in logs, metrics, and API responses. Declare
kinds as package-level variables, just like `var ErrX = errors.New(...)`:

```go
var (
    ErrOrderNotFound = errors.NewKind("ORDER_NOT_FOUND", errors.ErrNotFound, "order not found")
    ErrOrderShipped  = errors.NewKind("ORDER_SHIPPED", errors.ErrFailedPrecondition, "order already shipped")
)
```

A kind has three parts:

- **Type code.** An uppercase string such as `ORDER_NOT_FOUND`. It can contain
  `A-Z`, `0-9`, `_`, and `.`. It can't be the same as a built-in category's
  code.
- **Category.** A built-in category. The kind inherits its HTTP and gRPC codes,
  and errors made from the kind match the category with `errors.Is`.
- **Message.** The error text of the kind itself.

Options adjust a kind when it's created:

```go
var ErrPaymentRequired = errors.NewKind(
    "PAYMENT_REQUIRED",
    errors.ErrFailedPrecondition,
    "account has no active subscription",
    errors.WithHTTPCode(http.StatusPaymentRequired),        // 402 instead of 400
    errors.WithPublicMessage("A subscription is required"), // shown to clients
)

var ErrQuotaExceeded = errors.NewKind(
    "QUOTA_EXCEEDED",
    errors.ErrForbidden,
    "storage quota exceeded",
    errors.WithGRPCCode(codes.ResourceExhausted),           // instead of PermissionDenied
)
```

`NewKind` panics on invalid input, such as a bad type code or an HTTP code
outside 400–599. Kinds are declared at package level, so these mistakes
surface when the program starts.

## Classifying errors and adding context

Categories and kinds have the same methods for creating errors:

| Method                     | Result text      | Keeps a cause | Use it to                                    |
|----------------------------|------------------|---------------|----------------------------------------------|
| `Msg(msg)` / `Msgf(...)`   | `msg`            | no            | create a new error with a detailed message   |
| `WithCause(err)`           | `err.Error()`    | yes           | classify an error without changing its text  |
| `Wrap(err, msg)` / `Wrapf` | `msg: err`       | yes           | classify an error and add context            |

```go
errors.ErrInvalidArgument.Msg("page size must be positive")
ErrOrderShipped.Msgf("order %d shipped at %s", id, shippedAt)
errors.ErrUnavailable.WithCause(err)
ErrEmailTaken.Wrap(pgErr, "register user")
```

`Wrap`, `Wrapf`, and `WithCause` return `nil` when the error is `nil`.
`errors.WithKind(err, kind)` does the same as `kind.WithCause(err)`. Use it to
classify a call's result directly:

```go
return errors.WithKind(s.store.Put(ctx, obj), ErrStorageUnavailable)
```

### Add context without changing the classification

After an error is classified, add context as it moves up through your code.
The package-level `errors.Wrap` / `errors.Wrapf` and `fmt.Errorf` with `%w`
all keep the classification and the cause:

```go
err := errors.ErrNotFound.Wrap(sql.ErrNoRows, "find user")
err = fmt.Errorf("repository: %w", err)
err = errors.Wrap(err, "login")

fmt.Println(err)                  // login: repository: find user: sql: no rows in result set
fmt.Println(errors.TypeCode(err)) // NOT_FOUND
```

When `errors.Wrap` is given a category directly, it uses the message as the
whole error text:

```go
err := errors.Wrap(errors.ErrNotFound, "user 42 not found")
fmt.Println(err) // user 42 not found
```

### Reclassifying

The outermost classification wins. To change what an error means as it
crosses a boundary, wrap it with another category or kind:

```go
// A missing user during login is an authentication failure, not a 404.
err := errors.ErrUnauthenticated.Wrap(err, "login")
```

The original error is still in the chain, so `errors.Is(err, ErrUserNotFound)`
is still `true`.

## Checking errors

Use `errors.Is` with a kind to match exactly, or with a category to match
anything in that category:

```go
switch {
case errors.Is(err, ErrOrderShipped):
    // this specific failure
case errors.Is(err, errors.ErrNotFound):
    // any "not found", including ErrOrderNotFound
}
```

Two kinds are always different errors, even when they're in the same category.

Use `errors.As` or the generic `errors.AsType` to get an error value or an
interface from the chain:

```go
if coder, ok := errors.AsType[errors.TypeCoder](err); ok {
    log.Printf("type=%s", coder.TypeCode())
}
```

## Reading codes

`errors.TypeCode(err)`, `errors.HTTPCode(err)`, and `errors.GRPCCode(err)`
look through the whole error chain:

| Error                                            | `TypeCode`          | `HTTPCode` | `GRPCCode`         |
|--------------------------------------------------|---------------------|------------|--------------------|
| category or kind, anywhere in the chain          | its type code       | its status | its code           |
| gRPC status error (`status.Error`, client call)  | from its code and detail | from its code and detail | its code |
| `context.Canceled`                               | `CANCELED`          | 408        | `Canceled`         |
| `context.DeadlineExceeded`                       | `DEADLINE_EXCEEDED` | 504        | `DeadlineExceeded` |
| `errors.Join(...)` where all children agree      | the shared code     | shared     | shared             |
| `errors.Join(...)` where children differ         | `UNKNOWN`           | 500        | `Unknown`          |
| unclassified (`errors.New`, third-party errors)  | `UNKNOWN`           | 500        | `Unknown`          |
| `nil`                                            | `OK`                | 200        | `OK`               |
| `ErrOK` returned or wrapped as an error          | `UNKNOWN`           | 500        | `Unknown`          |

If an error has several classifications in its chain, the outermost one is
used. To choose the classification of a joined error, wrap the join:
`errors.ErrInvalidArgument.WithCause(errors.Join(errs...))`.

### Why UNKNOWN and not INTERNAL?

Unclassified errors are reported as `UNKNOWN` so that they stay separate from
errors you have deliberately marked `ErrInternal`. When `UNKNOWN` shows up in
your logs, dashboards, or client responses, it points to an error that hasn't
been classified yet.

A non-nil error never reports success: `HTTPCode(err)` is always 400–599 and
`GRPCCode(err)` is never `OK` when `err != nil`. Only `nil` gives `OK`, 200,
and `codes.OK`.

## Public messages

Error text often includes SQL errors, hostnames, or file paths that clients
shouldn't see. `errors.PublicMessage(err)` returns a message that is safe to
send. It never returns `err.Error()`. It uses the first of these that applies:

1. A message set on this error with `errors.WrapPublicMessage(err, msg)`.
2. The kind's `WithPublicMessage` text.
3. `"Internal Server Error"` for a 5xx status or an unclassified error.
4. The HTTP status text, such as `"Not Found"` or `"Conflict"`.

```go
err := errors.ErrInternal.Wrap(dbErr, "load user")
errors.PublicMessage(err) // "Internal Server Error"

err = errors.WrapPublicMessage(errors.ErrInvalidArgument.Msg("sku ABC-123 failed checksum"),
    "The product code is not valid")
errors.PublicMessage(err) // "The product code is not valid"
err.Error()               // "sku ABC-123 failed checksum"
```

`errors.Public(err)` returns a `PublicError{Code, Message}` that can be
serialized directly:

```json
{"code":"USER_NOT_FOUND","message":"The user could not be found"}
```

## HTTP responses

A single helper can turn any error into a response:

```go
func writeError(w http.ResponseWriter, err error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(errors.HTTPCode(err))
    _ = json.NewEncoder(w).Encode(errors.Public(err))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.svc.GetUser(r.Context(), r.PathValue("id"))
    if err != nil {
        slog.ErrorContext(r.Context(), "get user", "err", err, "type", errors.TypeCode(err))
        writeError(w, err)
        return
    }
    ...
}
```

| Error                                                     | Response                                                                   |
|-----------------------------------------------------------|----------------------------------------------------------------------------|
| `ErrUserNotFound.Wrap(sql.ErrNoRows, "find user 42")`     | `404 {"code":"USER_NOT_FOUND","message":"The user could not be found"}`    |
| `errors.ErrInternal.Wrap(dialErr, "load user")`           | `500 {"code":"INTERNAL","message":"Internal Server Error"}`                |
| `errors.ErrUnprocessableEntity.Msg("email: missing @")`   | `422 {"code":"UNPROCESSABLE_ENTITY","message":"Unprocessable Entity"}`     |

## gRPC

`SendGRPCError` converts an error into a gRPC status on the server.
`ReceiveGRPCError` converts it back on the client. The status contains:

- the gRPC code from `GRPCCode(err)`;
- the message from `PublicMessage(err)`, so internal text stays on the server;
- a status detail with the type code, HTTP status, and category;
- any other details from a status error in the chain, such as
  `errdetails.BadRequest`.

On the client, the received error has the sender's type code, category, HTTP
status, gRPC code, and public message. That holds across any number of hops.

### Interceptors

The `grpcerrs` subpackage applies these functions to every call:

```go
import "github.com/stackus/errors/grpcerrs"

server := grpc.NewServer(
    grpc.ChainUnaryInterceptor(grpcerrs.UnaryServerInterceptor()),
    grpc.ChainStreamInterceptor(grpcerrs.StreamServerInterceptor()),
)

conn, err := grpc.NewClient(target,
    grpc.WithChainUnaryInterceptor(grpcerrs.UnaryClientInterceptor(registry)),
    grpc.WithChainStreamInterceptor(grpcerrs.StreamClientInterceptor(registry)),
)
```

Handlers return this package's errors, and callers check them as if the
server were local:

```go
// server
func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
    order, err := s.svc.GetOrder(ctx, req.GetId())
    if err != nil {
        return nil, err // e.g. orders.ErrOrderNotFound.Wrap(sql.ErrNoRows, "get order 7")
    }
    ...
}

// client
_, err := client.GetOrder(ctx, &pb.GetOrderRequest{Id: 7})

errors.Is(err, orders.ErrOrderNotFound)  // true, with the kind registered
errors.Is(err, errors.ErrNotFound)       // true
errors.TypeCode(err)                     // "ORDER_NOT_FOUND"
errors.HTTPCode(err)                     // 404, so a gateway can pass it on
errors.PublicMessage(err)                // "Order not found"
err.Error()                              // "rpc error: code = NotFound desc = Order not found"
```

The stream client interceptor converts errors from opening the stream and
from `Header`, `SendMsg`, `RecvMsg`, and `CloseSend`. `io.EOF` is returned
unchanged.

### Restoring kinds with a Registry

Every received error matches its built-in category and keeps its type code,
HTTP status, and public message. To also match your own kinds with
`errors.Is`, the client needs the same `Kind` values, usually from a package
shared with the server, registered in a `Registry`:

```go
registry, err := errors.NewRegistry(
    orders.ErrOrderNotFound,
    orders.ErrOrderShipped,
)
```

A received kind is matched by type code, gRPC code, and category. The local
kind's definition supplies the HTTP status. You can pass several registries;
they are searched in order.

### Plain gRPC status errors

Status errors made with `google.golang.org/grpc/status` are classified by
their code, whether they are returned by a handler or received without the
client interceptor, and even when wrapped:

```go
err := fmt.Errorf("get order: %w", status.Error(codes.NotFound, "order 7 not found"))
errors.TypeCode(err)      // "NOT_FOUND"
errors.HTTPCode(err)      // 404
errors.PublicMessage(err) // "Not Found"
```

A status without this package's detail, such as one from `status.Error` or
from a service that doesn't use this package, keeps its details. Its message
is not treated as public, because it may contain internal text. Use
`WrapPublicMessage` to choose the message clients see.

### Details

- The gRPC status code decides the result. If the detail disagrees with it,
  the detail is ignored.
- Wrapping a received error with a category or kind reclassifies it. The
  upstream status details are then not forwarded.
- `SendGRPCError(nil)` and `ReceiveGRPCError(nil)` return `nil`.
  `ReceiveGRPCError` returns errors that aren't gRPC statuses unchanged.

## Using your own error types

An existing error type can be classified by implementing one or more of these
interfaces. It doesn't need to be wrapped:

```go
type TypeCoder interface { error; TypeCode() string }
type HTTPCoder interface { error; HTTPCode() int }
type GRPCCoder interface { error; GRPCCode() codes.Code }
```

```go
type ValidationError struct{ Field string }

func (e ValidationError) Error() string    { return "invalid field: " + e.Field }
func (e ValidationError) TypeCode() string { return "VALIDATION_FAILED" }
func (e ValidationError) HTTPCode() int    { return http.StatusUnprocessableEntity }

err := fmt.Errorf("create account: %w", ValidationError{Field: "email"})
errors.TypeCode(err) // "VALIDATION_FAILED"
errors.HTTPCode(err) // 422
errors.GRPCCode(err) // codes.Unknown (not implemented)
```

Codes a type doesn't provide use the `UNKNOWN` defaults. An HTTP code outside
400–599 is treated as 500, and a gRPC code of `OK` as `Unknown`. Every error in
this package implements all three interfaces.

## Built-in categories

Categories named for gRPC codes:

| Category                | Type code             | HTTP | gRPC                 |
|-------------------------|-----------------------|------|----------------------|
| `ErrCanceled`           | `CANCELED`            | 408  | `Canceled`           |
| `ErrUnknown`            | `UNKNOWN`             | 500  | `Unknown`            |
| `ErrInvalidArgument`    | `INVALID_ARGUMENT`    | 400  | `InvalidArgument`    |
| `ErrDeadlineExceeded`   | `DEADLINE_EXCEEDED`   | 504  | `DeadlineExceeded`   |
| `ErrNotFound`           | `NOT_FOUND`           | 404  | `NotFound`           |
| `ErrAlreadyExists`      | `ALREADY_EXISTS`      | 409  | `AlreadyExists`      |
| `ErrPermissionDenied`   | `PERMISSION_DENIED`   | 403  | `PermissionDenied`   |
| `ErrResourceExhausted`  | `RESOURCE_EXHAUSTED`  | 429  | `ResourceExhausted`  |
| `ErrFailedPrecondition` | `FAILED_PRECONDITION` | 400  | `FailedPrecondition` |
| `ErrAborted`            | `ABORTED`             | 409  | `Aborted`            |
| `ErrOutOfRange`         | `OUT_OF_RANGE`        | 422  | `OutOfRange`         |
| `ErrUnimplemented`      | `UNIMPLEMENTED`       | 501  | `Unimplemented`      |
| `ErrInternal`           | `INTERNAL`            | 500  | `Internal`           |
| `ErrUnavailable`        | `UNAVAILABLE`         | 503  | `Unavailable`        |
| `ErrDataLoss`           | `DATA_LOSS`           | 500  | `DataLoss`           |
| `ErrUnauthenticated`    | `UNAUTHENTICATED`     | 401  | `Unauthenticated`    |

Categories named for HTTP statuses:

| Category                        | Type code                       | HTTP | gRPC                |
|---------------------------------|---------------------------------|------|---------------------|
| `ErrBadRequest`                 | `BAD_REQUEST`                   | 400  | `InvalidArgument`   |
| `ErrUnauthorized`               | `UNAUTHORIZED`                  | 401  | `Unauthenticated`   |
| `ErrForbidden`                  | `FORBIDDEN`                     | 403  | `PermissionDenied`  |
| `ErrMethodNotAllowed`           | `METHOD_NOT_ALLOWED`            | 405  | `Unimplemented`     |
| `ErrRequestTimeout`             | `REQUEST_TIMEOUT`               | 408  | `DeadlineExceeded`  |
| `ErrConflict`                   | `CONFLICT`                      | 409  | `Aborted`           |
| `ErrGone`                       | `GONE`                          | 410  | `NotFound`          |
| `ErrUnsupportedMediaType`       | `UNSUPPORTED_MEDIA_TYPE`        | 415  | `InvalidArgument`   |
| `ErrImATeapot`                  | `IM_A_TEAPOT`                   | 418  | `Unknown`           |
| `ErrUnprocessableEntity`        | `UNPROCESSABLE_ENTITY`          | 422  | `InvalidArgument`   |
| `ErrTooManyRequests`            | `TOO_MANY_REQUESTS`             | 429  | `ResourceExhausted` |
| `ErrUnavailableForLegalReasons` | `UNAVAILABLE_FOR_LEGAL_REASONS` | 451  | `PermissionDenied`  |
| `ErrInternalServerError`        | `INTERNAL_SERVER_ERROR`         | 500  | `Internal`          |
| `ErrNotImplemented`             | `NOT_IMPLEMENTED`               | 501  | `Unimplemented`     |
| `ErrBadGateway`                 | `BAD_GATEWAY`                   | 502  | `Unavailable`       |
| `ErrServiceUnavailable`         | `SERVICE_UNAVAILABLE`           | 503  | `Unavailable`       |
| `ErrGatewayTimeout`             | `GATEWAY_TIMEOUT`               | 504  | `DeadlineExceeded`  |

Each category is a separate error, so `ErrBadRequest` and `ErrInvalidArgument`
don't match each other with `errors.Is`, even though both use HTTP 400 and gRPC
`InvalidArgument`. A category received over gRPC keeps its identity when the
sender used this package.

`ErrOK` names the codes of a `nil` error: `OK`, 200, and `codes.OK`. It is
not an error to return. Returned or wrapped as an error, it is classified as
`UNKNOWN`, so a failure can't be reported as a success. It can't be used as a
kind's category or type code.

More runnable examples are in the
[package documentation](https://pkg.go.dev/github.com/stackus/errors#pkg-examples).

## Contributing

Pull requests are welcome. Please include tests for behavior changes.

## License

[MIT](LICENSE)
