Feature: Support SIRI StopDiscovery

  Background:
    Given a Referential "test" is created

  Scenario: 2464 3292 - Handle a SIRI StopDiscovery request
    Given a Partner "test" exists with connectors [siri-stop-points-discovery-request-broadcaster] and the following settings:
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
    And a StopArea exists with the following attributes:
      | Name      | Test                                                                      |
      | ObjectIDs | "internal": "NINOXE:StopPoint:BP:6:LOC"                                   |
      | Lines     | ["6ba7b814-9dad-11d1-2-00c04fd430c8","6ba7b814-9dad-11d1-3-00c04fd430c8"] |
    And a StopArea exists with the following attributes:
      | Name      | Test 2                                   |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:16:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test 3                                  |
      | ObjectIDs | "internal": "NINOXE:StopPoint:BP:7:LOC" |
      | Lines     | ["6ba7b814-9dad-11d1-4-00c04fd430c8"]   |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns7:StopPointsDiscovery xmlns:ns2="http://www.siri.org.uk/siri"
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
          </ns7:StopPointsDiscovery>
        </S:Body>
        </S:Envelope>
        """
    Then I should receive this SIRI response
      """
      <?xml version="1.0" encoding="UTF-8"?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <Answer version="2.0:FR-IDF-2.4">
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:Address>address</siri:Address>
            <siri:RequestMessageRef>STIF:Message::2345Fsdfrg35df:LOC</siri:RequestMessageRef>
            <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            <siri:Status>true</siri:Status>
              <siri:AnnotatedStopPointRef>
                <siri:StopPointRef>NINOXE:StopPoint:BP:6:LOC</siri:StopPointRef>
                <siri:Monitored>true</siri:Monitored>
                <siri:StopName>Test</siri:StopName>
                <siri:Lines>
                  <siri:LineRef>STIF:Line::C00272:</siri:LineRef>
                  <siri:LineRef>STIF:Line::C00273:</siri:LineRef>
                </siri:Lines>
              </siri:AnnotatedStopPointRef>
              <siri:AnnotatedStopPointRef>
                <siri:StopPointRef>NINOXE:StopPoint:BP:7:LOC</siri:StopPointRef>
                <siri:Monitored>true</siri:Monitored>
                <siri:StopName>Test 3</siri:StopName>
                <siri:Lines>
                  <siri:LineRef>STIF:Line::C00274:</siri:LineRef>
                </siri:Lines>
              </siri:AnnotatedStopPointRef>
            </Answer>
            <AnswerExtension/>
          </sw:StopPointsDiscoveryResponse>
        </S:Body>
      </S:Envelope>
      """
