Feature: Support SIRI request
  Background:
    Given a Referential "test" is created

  @skip-siri-valid
  Scenario: Handle a empty SIRI request
    When I send this SIRI request
      """
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <S:Fault>
      <faultcode>S:Client</faultcode>
      <faultstring>Invalid Request: empty body</faultstring>
    </S:Fault>
  </S:Body>
</S:Envelope>
      """

  @skip-siri-valid
  Scenario: Handle an invalid SIRI request
    When I send this SIRI request
      """
Invalid Request
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <S:Fault>
      <faultcode>S:Client</faultcode>
      <faultstring>Invalid Request: failed to parse xml input</faultstring>
    </S:Fault>
  </S:Body>
</S:Envelope>
      """
