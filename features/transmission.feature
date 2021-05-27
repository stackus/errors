Feature: Transmission of errors over GRPC
  Errors with codes can be sent to GRPC clients

  Scenario: nil errors are not turned into problems
    Given the error is nil
    When the error is sent over GRPC
    Then the GRPC code is "OK"
    Then the HTTP status is "OK"
    Then the Type code is ""

  Scenario: standard errors are treated as unknowns
    Given an error with the message "standard error"
    When the error is sent over GRPC
    Then the GRPC code is "Unknown"
    Then the HTTP status is "Not Extended"
    Then the Type code is "UNKNOWN"
    Then the error message is "standard error"

  Scenario: GRPC errors do not pick up extra info
    Given an error with GRPC code "codes.PermissionDenied"
    When the error is sent over GRPC
    Then the GRPC code is "PermissionDenied"
    Then the HTTP status is "Not Extended"
    Then the Type code is "UNKNOWN"

  Scenario: all codes can be transmitted
    Given an error with GRPC code "codes.Unimplemented"
    Given an error with HTTP status "http.StatusForbidden"
    Given an error with Type code "BAD_REQUEST"
    Given an error with the message "error message"
    When the error is sent over GRPC
    Then the HTTP status is "Forbidden"
    Then the Type code is "BAD_REQUEST"
    Then the GRPC code is "Unimplemented"
    Then the error message is "error message"
    Then the error is a "ErrNotImplemented"
    Then the error is a "ErrPermissionDenied"
    Then the error is a "ErrBadRequest"