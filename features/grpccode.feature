Feature: grpc coded errors
  Errors can be modified to contain grpc codes

  Scenario: nil errors are not problems
    Given the error is nil
    Then the GRPC code is "OK"

  Scenario: treat errors without any GRPC code as unknown
    Given the error does not implement GRPCCoder{}
    Then the GRPC code is "Unknown"

  Scenario: package Errors have GRPC codes
    Given the error is "ErrNotFound"
    Then the GRPC code is "NotFound"

  Scenario: new errors can be configured with GRPC codes
    Given an error with GRPC code "codes.AlreadyExists"
    Then the GRPC code is "AlreadyExists"

  Scenario: GRPC codes can be sent over GRPC connections
    Given the error is "ErrNotImplemented"
    When the error is sent over GRPC
    Then the GRPC code is "Unimplemented"