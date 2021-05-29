Feature: Builtin errors have all the right codes

  Scenario Outline: verify error codes
    Given the error is "<error>"
    Then the Type code is "<type code>"
    And the HTTP status is "<http code>"
    And the GRPC code is "<grpc code>"

    Scenarios:
      | error                         | type code                     | http code                     | grpc code          |
      | ErrOK                         |                               | OK                            | OK                 |
      | ErrCanceled                   | CANCELED                      | Request Timeout               | Canceled           |
      | ErrUnknown                    | UNKNOWN                       | Not Extended                  | Unknown            |
      | ErrInvalidArgument            | INVALID_ARGUMENT              | Bad Request                   | InvalidArgument    |
      | ErrDeadlineExceeded           | DEADLINE_EXCEEDED             | Gateway Timeout               | DeadlineExceeded   |
      | ErrNotFound                   | NOT_FOUND                     | Not Found                     | NotFound           |
      | ErrAlreadyExists              | ALREADY_EXISTS                | Conflict                      | AlreadyExists      |
      | ErrPermissionDenied           | PERMISSION_DENIED             | Forbidden                     | PermissionDenied   |
      | ErrResourceExhausted          | RESOURCE_EXHAUSTED            | Too Many Requests             | ResourceExhausted  |
      | ErrFailedPrecondition         | FAILED_PRECONDITION           | Bad Request                   | FailedPrecondition |
      | ErrAborted                    | ABORTED                       | Conflict                      | Aborted            |
      | ErrOutOfRange                 | OUT_OF_RANGE                  | Unprocessable Entity          | OutOfRange         |
      | ErrUnimplemented              | UNIMPLEMENTED                 | Not Implemented               | Unimplemented      |
      | ErrInternal                   | INTERNAL                      | Internal Server Error         | Internal           |
      | ErrUnavailable                | UNAVAILABLE                   | Service Unavailable           | Unavailable        |
      | ErrDataLoss                   | DATA_LOSS                     | Internal Server Error         | DataLoss           |
      | ErrUnauthenticated            | UNAUTHENTICATED               | Unauthorized                  | Unauthenticated    |
      | ErrBadRequest                 | BAD_REQUEST                   | Bad Request                   | InvalidArgument    |
      | ErrUnauthorized               | UNAUTHORIZED                  | Unauthorized                  | Unauthenticated    |
      | ErrForbidden                  | FORBIDDEN                     | Forbidden                     | PermissionDenied   |
      | ErrMethodNotAllowed           | METHOD_NOT_ALLOWED            | Method Not Allowed            | Unimplemented      |
      | ErrRequestTimeout             | REQUEST_TIMEOUT               | Request Timeout               | DeadlineExceeded   |
      | ErrConflict                   | CONFLICT                      | Conflict                      | AlreadyExists      |
      | ErrImATeapot                  | IM_A_TEAPOT                   | I'm a teapot                  | Unknown            |
      | ErrUnprocessableEntity        | UNPROCESSABLE_ENTITY          | Unprocessable Entity          | InvalidArgument    |
      | ErrTooManyRequests            | TOO_MANY_REQUESTS             | Too Many Requests             | ResourceExhausted  |
      | ErrUnavailableForLegalReasons | UNAVAILABLE_FOR_LEGAL_REASONS | Unavailable For Legal Reasons | Unavailable        |
      | ErrInternalServerError        | INTERNAL_SERVER_ERROR         | Internal Server Error         | Internal           |
      | ErrNotImplemented             | NOT_IMPLEMENTED               | Not Implemented               | Unimplemented      |
      | ErrBadGateway                 | BAD_GATEWAY                   | Bad Gateway                   | Aborted            |
      | ErrServiceUnavailable         | SERVICE_UNAVAILABLE           | Service Unavailable           | Unavailable        |
      | ErrGatewayTimeout             | GATEWAY_TIMEOUT               | Gateway Timeout               | DeadlineExceeded   |
