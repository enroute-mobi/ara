Feature: Support SIRI StopDiscovery

  Background:
    Given a Referential "test" is created

  Scenario: 2464 - Handle a SIRI StopDiscovery request
    Given a Partner "test" exists with connectors [siri-stop-points-discovery-request-broadcaster] and the following settings:
      | local_credential     | test                  |
      | remote_objectid_kind | internal              |
    And a StopArea exists with the following attributes:
      | Name                  | Test                                     |
      | ObjectIDs             | "internal": "NINOXE:StopPoint:BP:6:LOC"   |
    And a StopArea exists with the following attributes:
      | Name                  | Test 2                                    |
      | ObjectIDs             | "internal": "NINOXE:StopPoint:SP:16:LOC"   |
    And a StopArea exists with the following attributes:
      | Name                  | Test 3                                    |
      | ObjectIDs             | "internal": "NINOXE:StopPoint:BP:7:LOC"   |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns7:StopPointsDiscovery xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:ns3="http://www.ifopt.org.uk/acsb"
                           xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                           xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                           xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <Request>
              <ns2:RequestTimestamp>2017-03-03T11:28:00.359Z</ns2:RequestTimestamp>
              <ns2:RequestorRef>RATPDev</ns2:RequestorRef>
            </Request>
            <RequestExtension />
          </ns7:StopPointsDiscovery>
        </S:Body>
      </S:Envelope>
        """
    Then I should receive this SIRI response
      """
      <?xml version="1.0" encoding="UTF-8"?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENV='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <ns8:StopPointsDiscoveryResponse xmlns:ns8="http://wsdl.siri.org.uk" xmlns:ns3="http://www.siri.org.uk/siri" xmlns:ns4="http://www.ifopt.org.uk/acsb" xmlns:ns5="http://www.ifopt.org.uk/ifopt" xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns7="http://scma/siri" xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>2017-03-03T11:28:30.359Z</ns3:ResponseTimestamp>
              <ns3:ProducerRef>test</ns3:ProducerRef>
              <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7:LOC</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>RATPDev:Message::3dfg56:LOC</ns3:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer version="2.0">
              <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
              <ns3:Status>true</ns3:Status>
              <ns3:AnnotatedStopPointRef>
                <ns3:StopPointRef>NINOXE:StopPoint:BP:6:LOC</ns3:StopPointRef>
                <ns3:StopName>Test</ns3:StopName>
              </ns3:AnnotatedStopPointRef>
              <ns3:AnnotatedStopPointRef>
                <ns3:StopPointRef>NINOXE:StopPoint:BP:7:LOC</ns3:StopPointRef>
                <ns3:StopName>Test 3</ns3:StopName>
              </ns3:AnnotatedStopPointRef>
              <ns3:AnnotatedStopPointRef>
                <ns3:StopPointRef>NINOXE:StopPoint:SP:16:LOC</ns3:StopPointRef>
                <ns3:StopName>Test 2</ns3:StopName>
              </ns3:AnnotatedStopPointRef>
            </Answer>
            <AnswerExtension />
          </ns8:StopPointsDiscoveryResponse>
        </S:Body>
      </S:Envelope>
      """
