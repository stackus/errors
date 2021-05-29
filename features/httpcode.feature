Feature: HTTP status errors
  Errors can be modified to contain HTTP statuses

  Scenario: nil errors are not problems
    Given the error is nil
    Then the HTTP status is "OK"

  Scenario: treat errors without any HTTP status as unknown
    Given the error does not implement HTTPCoder{}
    Then the HTTP status is "Not Extended"

  Scenario: package Errors have HTTP statuses
    Given the error is "ErrNotFound"
    Then the HTTP status is "Not Found"

  Scenario: new errors can be configured with HTTP statuses
    Given an error with HTTP status "http.StatusTeapot"
    Then the HTTP status is "I'm a teapot"

  Scenario: HTTP statuses can be sent over GRPC connections
    Given the error is "ErrNotImplemented"
    When the error is sent over GRPC
    Then the HTTP status is "Not Implemented"