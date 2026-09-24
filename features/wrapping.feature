Feature: Wrapping errors
  Errors can be wrapped to embed codes or to prefix messages

  Scenario: error information cannot be overridden
    Given the error is "ErrBadRequest"
    When wrapped with the error "ErrForbidden" and message "some error"
    Then the HTTP status is "Bad Request"
    And the error message is "some error: FORBIDDEN"
    And the error is a "ErrBadRequest"
    And the error is a "ErrForbidden"
