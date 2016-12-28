Feature: Siri server

  @server
  Scenario: Handle a SIRI Checkstatus request
    Given a Referential "test" exists
      And we send a checkstatus request for referential "test"
    Then we should recieve a positive checkstatus response
