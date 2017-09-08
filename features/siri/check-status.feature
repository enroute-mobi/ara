Feature: Support SIRI CheckStatus
  Background:
    Given a Referential "test" is created

  Scenario: 2460 - Handle a SIRI Checkstatus request
    Given a Partner "test" exists with connectors [siri-check-status-server] and the following settings:
      | local_credential | test |
    When I send this SIRI request to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test</siri:RequestorRef>
        <siri:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatusResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <CheckStatusAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Edwig</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:RequestMessageRef>
      </CheckStatusAnswerInfo>
      <Answer>
        <siri:Status>true</siri:Status>
        <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
      </Answer>
      <AnswerExtension/>
    </sw:CheckStatusResponse>
  </S:Body>
</S:Envelope>
      """

  Scenario: Handle a SIRI Checkstatus request with invalid RequestorRef
    When I send this SIRI request
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>invalid</siri:RequestorRef>
        <siri:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <S:Fault>
      <faultcode>S:UnknownCredential</faultcode>
      <faultstring>RequestorRef Unknown</faultstring>
    </S:Fault>
  </S:Body>
</S:Envelope>
      """
