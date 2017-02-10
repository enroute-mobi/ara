Feature: Support SIRI request
  Background:
    Given a Referential "test" is created

  @wip
  # See Issue 2564
  Scenario: Handle a empty SIRI request
    When I send this SIRI request
      """
      """
    Then I should receive this SIRI reponse
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <S:Fault xmlns:ns4='http://www.w3.org/2003/05/soap-envelope'>
      <faultcode>S:Client</faultcode>
      <faultstring>Invalid Request</faultstring>
    </S:Fault>
  </S:Body>
</S:Envelope>
      """

  @wip
  # See Issue 2564
  Scenario: Handle an invalid SIRI request
    When I send this SIRI request
      """
Invalid Request
      """
    Then I should receive this SIRI reponse
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <S:Fault xmlns:ns4='http://www.w3.org/2003/05/soap-envelope'>
      <faultcode>S:Client</faultcode>
      <faultstring>Invalid Request</faultstring>
    </S:Fault>
  </S:Body>
</S:Envelope>
      """
