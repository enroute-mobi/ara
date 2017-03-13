Feature: Support SIRI CheckStatus
  Background:
    Given a Referential "test" is created

  Scenario: Handle a SIRI Checkstatus request
    Given a Partner "test" exists with connectors [siri-check-status-server] and the following settings:
      | local_credential | test |
    When I send this SIRI request to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ns7:CheckStatus xmlns:ns2="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:ns7="http://wsdl.siri.org.uk">
      <Request>
        <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
        <ns2:RequestorRef>test</ns2:RequestorRef>
        <ns2:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</ns2:MessageIdentifier>
      </Request>
      <RequestExtension />
    </ns7:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ns7:CheckStatusResponse xmlns:ns2="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:ns7="http://wsdl.siri.org.uk"
    xmlns:ns8="http://wsdl.siri.org.uk/siri">
      <CheckStatusAnswerInfo>
        <ns2:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns2:ResponseTimestamp>
        <ns2:ProducerRef>Edwig</ns2:ProducerRef>
        <ns2:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ns2:ResponseMessageIdentifier>
        <ns2:RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</ns2:RequestMessageRef>
      </CheckStatusAnswerInfo>
      <Answer>
        <ns2:Status>true</ns2:Status>
        <ns2:ServiceStartedTime>2017-01-01T12:00:00.000Z</ns2:ServiceStartedTime>
      </Answer>
      <AnswerExtension />
    </ns7:CheckStatusResponse>
  </S:Body>
</S:Envelope>
      """

  Scenario: Handle a SIRI Checkstatus request with invalid RequestorRef
    When I send this SIRI request
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ns7:CheckStatus xmlns:ns2="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:ns7="http://wsdl.siri.org.uk">
      <Request>
        <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
        <ns2:RequestorRef>invalid</ns2:RequestorRef>
        <ns2:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</ns2:MessageIdentifier>
      </Request>
      <RequestExtension />
    </ns7:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <S:Fault xmlns:ns4='http://www.w3.org/2003/05/soap-envelope'>
      <faultcode>S:UnknownCredential</faultcode>
      <faultstring>RequestorRef Unknown</faultstring>
    </S:Fault>
  </S:Body>
</S:Envelope>
      """
