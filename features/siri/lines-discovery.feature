Feature: Support SIRI LinesDiscovery

  Background:
    Given a Referential "test" is created

  Scenario: 4397 - Handle a SIRI LinesDiscovery request
    Given a Partner "test" exists with connectors [siri-lines-discovery-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
      | local_url            | address  |
    And a Line exists with the following attributes:
      | Name      | Line 1                          |
      | ObjectIDs | "internal":"STIF:Line::C00272:" |
    And a Line exists with the following attributes:
      | Name      | Line 2                          |
      | ObjectIDs | "internal":"STIF:Line::C00273:" |
    And a Line exists with the following attributes:
      | Name      | Line 3                          |
      | ObjectIDs | "internal":"STIF:Line::C00274:" |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns7:LinesDiscovery xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
                           xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                           xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                           xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <Request>
              <ns2:RequestTimestamp>2017-03-03T11:28:00.359Z</ns2:RequestTimestamp>
              <ns2:RequestorRef>test</ns2:RequestorRef>
              <ns2:MessageIdentifier>STIF:Message::2345Fsdfrg35df:LOC</ns2:MessageIdentifier>
            </Request>
            <RequestExtension />
          </ns7:LinesDiscovery>
        </S:Body>
        </S:Envelope>
        """
    Then I should receive this SIRI response
      """
      <?xml version="1.0" encoding="UTF-8"?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:LinesDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <Answer version="2.0">
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:Status>true</siri:Status>
              <siri:AnnotatedLineRef>
                <siri:LineRef>STIF:Line::C00272:</siri:LineRef>
                <siri:LineName>Line 1</siri:LineName>
                <siri:Monitored>true</siri:Monitored>
              </siri:AnnotatedLineRef>
              <siri:AnnotatedLineRef>
                <siri:LineRef>STIF:Line::C00273:</siri:LineRef>
                <siri:LineName>Line 2</siri:LineName>
                <siri:Monitored>true</siri:Monitored>
              </siri:AnnotatedLineRef>
              <siri:AnnotatedLineRef>
                <siri:LineRef>STIF:Line::C00274:</siri:LineRef>
                <siri:LineName>Line 3</siri:LineName>
                <siri:Monitored>true</siri:Monitored>
              </siri:AnnotatedLineRef>
            </Answer>
            <AnswerExtension />
          </sw:LinesDiscoveryResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: LinesDiscovery collect
    Given a SIRI server waits LinesDiscovery request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:LinesDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
      <siri:Status>true</siri:Status>
        <siri:AnnotatedLineRef>
          <siri:LineRef>NINOXE:Line:BP:6:LOC</siri:LineRef>
          <siri:LineName>Test</siri:LineName>
          <siri:Monitored>true</siri:Monitored>
        </siri:AnnotatedLineRef>
        <siri:AnnotatedLineRef>
          <siri:LineRef>NINOXE:Line:BP:7:LOC</siri:LineRef>
          <siri:LineName>Test 3</siri:LineName>
          <siri:Monitored>true</siri:Monitored>
        </siri:AnnotatedLineRef>
      </Answer>
      <AnswerExtension/>
    </sw:LinesDiscoveryResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-lines-discovery-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_objectid_kind  | internal              |
    And a minute has passed
    Then a Line "internal":"NINOXE:Line:BP:6:LOC" should exist
    And a Line "internal":"NINOXE:Line:BP:7:LOC" should exist

  @ARA-1298 @siri-vali
  Scenario: LinesDiscovery collect
    Given a SIRI server waits LinesDiscovery request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:LinesDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
      <siri:Status>true</siri:Status>
        <siri:AnnotatedLineRef>
          <siri:LineRef>NINOXE:Line:BP:6:LOC</siri:LineRef>
          <siri:LineName>Test</siri:LineName>
          <siri:Monitored>true</siri:Monitored>
        </siri:AnnotatedLineRef>
        <siri:AnnotatedLineRef>
          <siri:LineRef>NINOXE:Line:BP:7:LOC</siri:LineRef>
          <siri:LineName>Test 3</siri:LineName>
          <siri:Monitored>true</siri:Monitored>
        </siri:AnnotatedLineRef>
      </Answer>
      <AnswerExtension/>
    </sw:LinesDiscoveryResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-lines-discovery-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_objectid_kind  | internal              |
    And a minute has passed
    Then a Line "internal":"NINOXE:Line:BP:6:LOC" should exist
    And a Line "internal":"NINOXE:Line:BP:7:LOC" should exist
    And the Partner "test" is updated with the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | remote_objectid_kind | external              |
    And a SIRI server waits LinesDiscovery request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:LinesDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
      <siri:Status>true</siri:Status>
        <siri:AnnotatedLineRef>
          <siri:LineRef>NINOXE:Line:BP:6:LOC</siri:LineRef>
          <siri:LineName>Test</siri:LineName>
          <siri:Monitored>true</siri:Monitored>
        </siri:AnnotatedLineRef>
        <siri:AnnotatedLineRef>
          <siri:LineRef>NINOXE:Line:BP:7:LOC</siri:LineRef>
          <siri:LineName>Test 3</siri:LineName>
          <siri:Monitored>true</siri:Monitored>
        </siri:AnnotatedLineRef>
      </Answer>
      <AnswerExtension/>
    </sw:LinesDiscoveryResponse>
  </S:Body>
</S:Envelope>
      """
    And a minute has passed
    Then a Line "internal":"NINOXE:Line:BP:6:LOC" should exist
    And a Line "internal":"NINOXE:Line:BP:7:LOC" should exist
    Then a Line "external":"NINOXE:Line:BP:6:LOC" should exist
    And a Line "external":"NINOXE:Line:BP:7:LOC" should exist
