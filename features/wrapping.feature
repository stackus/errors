Feature: Wrapping errors
  Errors can be wrapped to embed codes or to prefix messages

  Scenario: error messages can be prefixed
    Given an error with the message "some error"
    When wrapped with the message "more context"
    When wrapped with the message "even more context"
    Then the error message is "even more context: more context: some error"

  Scenario: errors with Type codes are embedded
    Given an error with Type code "CUSTOM"
    When wrapped with the message "more context"
    Then the error message is "more context"
    Then the Type code is "CUSTOM"

  Scenario: errors can be very custom
    Given an error with HTTP status "http.StatusNotFound"
    Given an error with GRPC code "codes.Internal"
    Given an error with Type code "VERY CUSTOM"
    Then the GRPC code is "Internal"
    Then the HTTP status is "Not Found"
    Then the Type code is "VERY CUSTOM"

  Scenario: error information can be overridden
    Given the error is "ErrBadRequest"
    When wrapped with the error "ErrForbidden" and message "some error"
    Then the HTTP status is "Forbidden"
    And the error message is "some error"
    And the error is a "ErrBadRequest"
    And the error is a "ErrForbidden"
