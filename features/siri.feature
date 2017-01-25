Feature: Siri server
  Background:
    Given a Referential "test" is created

  Scenario: Handle a SIRI Checkstatus request
    When we send a checkstatus request for referential "test"
    Then we should receive a positive checkstatus response
