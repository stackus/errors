// Package errors gives Go errors a type code, an HTTP status, and a gRPC code,
// so one error value can drive logs, HTTP responses, and gRPC statuses. It
// works alongside the standard library's errors package. [Is], [As],
// [AsType], [Unwrap], and [Join] are provided so it can replace that import.
//
// # Categories
//
// The Err* constants, such as [ErrNotFound] and [ErrInvalidArgument], are
// built-in categories of type [Error]. Each has a type code and HTTP and gRPC
// mappings (see the constants' comments). A category can be returned
// directly, matched with [Is], or used to classify another error:
//
//	return errors.ErrNotFound.Wrap(err, "find user")
//
// # Defining application errors
//
// [NewKind] creates a sentinel error with its own type code inside a category.
// Declare kinds as package-level variables:
//
//	var ErrUserNotFound = errors.NewKind("USER_NOT_FOUND", errors.ErrNotFound, "user not found")
//
// An error created from a kind matches both the kind and its category with
// [Is]. A kind inherits its category's HTTP and gRPC codes; [WithHTTPCode] and
// [WithGRPCCode] override them, and [WithPublicMessage] sets the message shown
// to clients.
//
// # Classifying and adding context
//
// [Error] and [*Kind] have the same set of constructors:
//
//   - Msg and Msgf create a new error with the given message and no cause.
//   - WithCause classifies an existing error and keeps its message.
//   - Wrap and Wrapf classify an existing error and prefix its message.
//
// Classify errors where they enter your application, for example when a
// database driver returns sql.ErrNoRows. After that, add context with [Wrap],
// [Wrapf], or fmt.Errorf with %w. Each keeps the classification and the
// original cause.
//
// # Reading codes
//
// [TypeCode], [HTTPCode], and [GRPCCode] return an error's codes. They follow
// these rules:
//
//   - The outermost classification in the chain wins, so wrapping a
//     classified error with a different category reclassifies it.
//   - context.Canceled and context.DeadlineExceeded are classified as
//     [ErrCanceled] and [ErrDeadlineExceeded].
//   - A joined error takes the classification its children share. When
//     they differ, or one is unclassified, the join is UNKNOWN.
//   - An unclassified error is UNKNOWN, HTTP 500, gRPC Unknown. Unknown
//     errors are kept separate from Internal ones on purpose: they show which
//     failures have not been classified yet.
//   - A nil error is OK, HTTP 200, gRPC OK, the codes named by [ErrOK].
//   - A non-nil error never reports success. [ErrOK] returned or wrapped as
//     an error is UNKNOWN, HTTP 500, gRPC Unknown.
//
// # Custom error types
//
// An error type of your own is classified when it implements [TypeCoder],
// [HTTPCoder], or [GRPCCoder]. Codes it does not provide use the UNKNOWN
// defaults.
//
// # Public messages
//
// An error's text often contains details that clients should not see.
// [PublicMessage] returns a client-safe message, and [Public] returns a code
// and message pair that can be serialized. To set the message for one error,
// use [WrapPublicMessage].
//
// # HTTP
//
// Set the response status from [HTTPCode] and write [Public] as the body. See
// the httpHandler example.
//
// # gRPC
//
// On the server, [SendGRPCError] converts an error into a gRPC status. The
// status carries the public message, the classification, and any details
// from a status error in the chain. On the client, [ReceiveGRPCError]
// rebuilds a classified error with the sender's type code, category, HTTP
// code, and public message. To also restore an application [Kind]'s
// identity, pass a [Registry] that contains it.
//
// The github.com/stackus/errors/grpcerrs package provides server and client
// interceptors that apply these functions to every call.
//
// gRPC status errors, such as those made with status.Error, are classified
// by their code wherever they appear in a chain.
package errors
