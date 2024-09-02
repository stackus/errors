package errors

// Errors named in line with GRPC codes and some that overlap with HTTP statuses
const (
	ErrOK                 Error = "OK"                  // HTTP: 200 GRPC: codes.OK
	ErrCanceled           Error = "CANCELED"            // HTTP: 408 GRPC: codes.Canceled
	ErrUnknown            Error = "UNKNOWN"             // HTTP: 510 GRPC: codes.Unknown
	ErrInvalidArgument    Error = "INVALID_ARGUMENT"    // HTTP: 400 GRPC: codes.InvalidArgument
	ErrDeadlineExceeded   Error = "DEADLINE_EXCEEDED"   // HTTP: 504 GRPC: codes.DeadlineExceeded
	ErrNotFound           Error = "NOT_FOUND"           // HTTP: 404 GRPC: codes.NotFound
	ErrAlreadyExists      Error = "ALREADY_EXISTS"      // HTTP: 409 GRPC: codes.AlreadyExists
	ErrPermissionDenied   Error = "PERMISSION_DENIED"   // HTTP: 403 GRPC: codes.PermissionDenied
	ErrResourceExhausted  Error = "RESOURCE_EXHAUSTED"  // HTTP: 429 GRPC: codes.ResourceExhausted
	ErrFailedPrecondition Error = "FAILED_PRECONDITION" // HTTP: 400 GRPC: codes.FailedPrecondition
	ErrAborted            Error = "ABORTED"             // HTTP: 409 GRPC: codes.Aborted
	ErrOutOfRange         Error = "OUT_OF_RANGE"        // HTTP: 422 GRPC: codes.OutOfRange
	ErrUnimplemented      Error = "UNIMPLEMENTED"       // HTTP: 501 GRPC: codes.Unimplemented
	ErrInternal           Error = "INTERNAL"            // HTTP: 500 GRPC: codes.Internal
	ErrUnavailable        Error = "UNAVAILABLE"         // HTTP: 503 GRPC: codes.Unavailable
	ErrDataLoss           Error = "DATA_LOSS"           // HTTP: 500 GRPC: codes.DataLoss
	ErrUnauthenticated    Error = "UNAUTHENTICATED"     // HTTP: 401 GRPC: codes.Unauthenticated
)

// Errors named in line with HTTP statuses
const (
	ErrBadRequest                 Error = "BAD_REQUEST"                   // HTTP: 400 GRPC: codes.InvalidArgument
	ErrUnauthorized               Error = "UNAUTHORIZED"                  // HTTP: 401 GRPC: codes.Unauthenticated
	ErrForbidden                  Error = "FORBIDDEN"                     // HTTP: 403 GRPC: codes.PermissionDenied
	ErrMethodNotAllowed           Error = "METHOD_NOT_ALLOWED"            // HTTP: 405 GRPC: codes.Unimplemented
	ErrRequestTimeout             Error = "REQUEST_TIMEOUT"               // HTTP: 408 GRPC: codes.DeadlineExceeded
	ErrConflict                   Error = "CONFLICT"                      // HTTP: 409 GRPC: codes.AlreadyExists
	ErrGone                       Error = "GONE"                          // HTTP: 410 GRPC: codes.NotFound
	ErrUnsupportedMediaType       Error = "UNSUPPORTED_MEDIA_TYPE"        // HTTP: 415 GRPC: codes.InvalidArgument
	ErrImATeapot                  Error = "IM_A_TEAPOT"                   // HTTP: 418 GRPC: codes.Unknown
	ErrUnprocessableEntity        Error = "UNPROCESSABLE_ENTITY"          // HTTP: 422 GRPC: codes.InvalidArgument
	ErrTooManyRequests            Error = "TOO_MANY_REQUESTS"             // HTTP: 429 GRPC: codes.ResourceExhausted
	ErrUnavailableForLegalReasons Error = "UNAVAILABLE_FOR_LEGAL_REASONS" // HTTP: 451 GRPC: codes.Unavailable
	ErrInternalServerError        Error = "INTERNAL_SERVER_ERROR"         // HTTP: 500 GRPC: codes.Internal
	ErrNotImplemented             Error = "NOT_IMPLEMENTED"               // HTTP: 501 GRPC: codes.Unimplemented
	ErrBadGateway                 Error = "BAD_GATEWAY"                   // HTTP: 502 GRPC: codes.Aborted
	ErrServiceUnavailable         Error = "SERVICE_UNAVAILABLE"           // HTTP: 503 GRPC: codes.Unavailable
	ErrGatewayTimeout             Error = "GATEWAY_TIMEOUT"               // HTTP: 504 GRPC: codes.DeadlineExceeded
)
