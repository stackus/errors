Feature: Type coded errors
  Errors can be modified to contain type codes

  Scenario: treat errors without any Type code as unknown
    Given the error does not implement TypeCoder{}
    Then the Type code is "UNKNOWN"

  Scenario: package Errors have Type codes
    Given the error is "ErrNotFound"
    Then the Type code is "NOT_FOUND"

  Scenario: new errors can be configured with Type codes
    Given an error with Type code "CUSTOM"
    Then the Type code is "CUSTOM"

  Scenario: type codes can be sent over GRPC connections
    Given the error is "ErrNotImplemented"
    When the error is sent over GRPC
    Then the Type code is "UNIMPLEMENTED"
