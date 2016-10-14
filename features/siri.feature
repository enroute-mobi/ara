Feature: Siri server

  @server
  Scenario: Handle a SIRI Checkstatus request
    Given we send a checkstatus request with body "checkstatus-soap-request.xml"
    Then we should recieve a checkstatus response with body "checkstatus-soap-response.xml"
