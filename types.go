package errors

// Errors
const (
	ErrOK                  Error = "OK"                   // HTTP: 200 GRPC: codes.OK
	ErrCanceled            Error = "CANCELED"             // HTTP: 408 GRPC: codes.Canceled
	ErrUnknown             Error = "UNKNOWN"              // HTTP: 500 GRPC: codes.Unknown
	ErrInvalidArgument     Error = "INVALID_ARGUMENT"     // HTTP: 400 GRPC: codes.InvalidArgument
	ErrDeadlineExceeded    Error = "DEADLINE_EXCEEDED"    // HTTP: 504 GRPC: codes.DeadlineExceeded
	ErrNotFound            Error = "NOT_FOUND"            // HTTP: 404 GRPC: codes.NotFound
	ErrAlreadyExists       Error = "ALREADY_EXISTS"       // HTTP: 409 GRPC: codes.AlreadyExists
	ErrPermissionDenied    Error = "PERMISSION_DENIED"    // HTTP: 403 GRPC: codes.PermissionDenied
	ErrResourceExhausted   Error = "RESOURCE_EXHAUSTED"   // HTTP: 429 GRPC: codes.ResourceExhausted
	ErrFailedPrecondition  Error = "FAILED_PRECONDITION"  // HTTP: 400 GRPC: codes.FailedPrecondition
	ErrAborted             Error = "ABORTED"              // HTTP: 409 GRPC: codes.Aborted
	ErrOutOfRange          Error = "OUT_OF_RANGE"         // HTTP: 422 GRPC: codes.OutOfRange
	ErrUnimplemented       Error = "UNIMPLEMENTED"        // HTTP: 501 GRPC: codes.Unimplemented
	ErrInternal            Error = "INTERNAL_ERROR"       // HTTP: 500 GRPC: codes.Internal
	ErrUnavailable         Error = "UNAVAILABLE"          // HTTP: 503 GRPC: codes.Unavailable
	ErrDataLoss            Error = "DATA_LOSS"            // HTTP: 500 GRPC: codes.DataLoss
	ErrUnauthenticated     Error = "UNAUTHENTICATED"      // HTTP: 401 GRPC: codes.Unauthenticated
	ErrBadRequest          Error = "BAD_REQUEST"          // HTTP: 400 GRPC: codes.InvalidArgument
	ErrConflict            Error = "CONFLICT"             // HTTP: 409 GRPC: codes.AlreadyExists
	ErrUnauthorized        Error = "UNAUTHORIZED"         // HTTP: 401 GRPC: codes.Unauthenticated
	ErrForbidden           Error = "FORBIDDEN"            // HTTP: 403 GRPC: codes.PermissionDenied
	ErrUnprocessableEntity Error = "UNPROCESSABLE_ENTITY" // HTTP: 422 GRPC: codes.InvalidArgument
	ErrServer              Error = "SERVER_ERROR"         // HTTP: 500 GRPC: codes.Internal
	ErrClient              Error = "CLIENT_ERROR"         // HTTP: 400 GRPC: codes.InvalidArgument
)
