Feature: Support SIRI StopPointsDiscovery

  Background:
    Given a Referential "test" is created

  @ARA-1095
  Scenario: Handle a SIRI StopPointsDiscovery request with ReferentID having wrong remoteCode kind
    Given a Partner "test" exists with connectors [siri-stop-points-discovery-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_code_space | internal |
      | local_url            | address  |
    And a Line exists with the following attributes:
      | Name      | Line 1                          |
      | Codes | "internal":"STIF:Line::C00272:" |
    And a Line exists with the following attributes:
      | Name      | Line 2                          |
      | Codes | "internal":"STIF:Line::C00273:" |
    And a StopArea exists with the following attributes:
      | Name      | Stop1Referent                                                             |
      | Codes | "internal": "NINOXE:StopPoint:BP:11:LOC"                                  |
      | Lines     | ["6ba7b814-9dad-11d1-2-00c04fd430c8","6ba7b814-9dad-11d1-3-00c04fd430c8"] |
    And a StopArea exists with the following attributes:
      | Name       | Stop2                                    |
      | Codes  | "internal": "NINOXE:StopPoint:SP:22:LOC" |
      | Lines      | ["6ba7b814-9dad-11d1-3-00c04fd430c8"]    |
      | ReferentID | 6ba7b814-9dad-11d1-4-00c04fd430c8        |
    And a StopArea exists with the following attributes:
      | Name      | Stop3Referent                                                              |
      | Codes | "wrong": "NINOXE:StopPoint:BP:33:LOC"                                      |
      | Lines     | ["6ba7b814-9dad-11d1-2-00c04fd430c8", "6ba7b814-9dad-11d1-3-00c04fd430c8"] |
    And a StopArea exists with the following attributes:
      | Name       | Stop4                                    |
      | Codes  | "internal": "NINOXE:StopPoint:BP:44:LOC" |
      | Lines      | ["6ba7b814-9dad-11d1-2-00c04fd430c8"]    |
      | ReferentID | 6ba7b814-9dad-11d1-6-00c04fd430c8        |
    And a StopArea exists with the following attributes:
      | Name       | Stop5                                    |
      | Codes  | "internal": "NINOXE:StopPoint:BP:55:LOC" |
      | Lines      | ["6ba7b814-9dad-11d1-3-00c04fd430c8"]    |
      | ReferentID | 6ba7b814-9dad-11d1-6-00c04fd430c8        |
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
            <Answer version="2.0">
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:Status>true</siri:Status>
              <siri:AnnotatedStopPointRef>
                <siri:StopPointRef>NINOXE:StopPoint:BP:11:LOC</siri:StopPointRef>
                <siri:Monitored>true</siri:Monitored>
                <siri:StopName>Stop1Referent</siri:StopName>
                <siri:Lines>
                  <siri:LineRef>STIF:Line::C00272:</siri:LineRef>
                  <siri:LineRef>STIF:Line::C00273:</siri:LineRef>
                </siri:Lines>
              </siri:AnnotatedStopPointRef>
              <siri:AnnotatedStopPointRef>
                <siri:StopPointRef>NINOXE:StopPoint:BP:44:LOC</siri:StopPointRef>
                <siri:Monitored>true</siri:Monitored>
                <siri:StopName>Stop4</siri:StopName>
                <siri:Lines>
                  <siri:LineRef>STIF:Line::C00272:</siri:LineRef>
                </siri:Lines>
              </siri:AnnotatedStopPointRef>
              <siri:AnnotatedStopPointRef>
                <siri:StopPointRef>NINOXE:StopPoint:BP:55:LOC</siri:StopPointRef>
                <siri:Monitored>true</siri:Monitored>
                <siri:StopName>Stop5</siri:StopName>
                <siri:Lines>
                  <siri:LineRef>STIF:Line::C00273:</siri:LineRef>
                </siri:Lines>
              </siri:AnnotatedStopPointRef>
            </Answer>
            <AnswerExtension/>
          </sw:StopPointsDiscoveryResponse>
        </S:Body>
      </S:Envelope>
      """
  
  @ARA-1095
  Scenario: Handle a SIRI StopPointsDiscovery request with ReferentID
    Given a Partner "test" exists with connectors [siri-stop-points-discovery-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_code_space | internal |
      | local_url            | address  |
    And a Line exists with the following attributes:
      | Name      | Line 1                          |
      | Codes | "internal":"STIF:Line::C00272:" |
    And a Line exists with the following attributes:
      | Name      | Line 2                          |
      | Codes | "internal":"STIF:Line::C00273:" |
    And a Line exists with the following attributes:
      | Name      | Line 3                          |
      | Codes | "internal":"STIF:Line::C00274:" |
    And a StopArea exists with the following attributes:
      | Name      | Stop1Referent                                                                                                  |
      | Codes | "internal": "NINOXE:StopPoint:BP:11:LOC"                                                                       |
      | Lines     | ["6ba7b814-9dad-11d1-2-00c04fd430c8","6ba7b814-9dad-11d1-3-00c04fd430c8", "6ba7b814-9dad-11d1-4-00c04fd430c8"] |
    And a StopArea exists with the following attributes:
      | Name       | Stop2                                 |
      | Codes  | "internal": "NINOXE:StopPoint:SP:22:LOC"  |
      | Lines      | ["6ba7b814-9dad-11d1-3-00c04fd430c8"] |
      | ReferentID | 6ba7b814-9dad-11d1-5-00c04fd430c8     |
    And a StopArea exists with the following attributes:
      | Name       | Stop3Referent                                                              |
      | Codes  | "internal": "NINOXE:StopPoint:BP:33:LOC"                                   |
      | Lines      | ["6ba7b814-9dad-11d1-2-00c04fd430c8", "6ba7b814-9dad-11d1-4-00c04fd430c8"] |
    And a StopArea exists with the following attributes:
      | Name       | Stop4                                    |
      | Codes  | "internal": "NINOXE:StopPoint:BP:44:LOC" |
      | Lines      | ["6ba7b814-9dad-11d1-2-00c04fd430c8"]    |
      | ReferentID | 6ba7b814-9dad-11d1-7-00c04fd430c8        |
    And a StopArea exists with the following attributes:
      | Name       | Stop5                                    |
      | Codes  | "internal": "NINOXE:StopPoint:BP:55:LOC" |
      | Lines      | ["6ba7b814-9dad-11d1-2-00c04fd430c8"]    |
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
            <Answer version="2.0">
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:Status>true</siri:Status>
              <siri:AnnotatedStopPointRef>
                <siri:StopPointRef>NINOXE:StopPoint:BP:11:LOC</siri:StopPointRef>
                <siri:Monitored>true</siri:Monitored>
                <siri:StopName>Stop1Referent</siri:StopName>
                <siri:Lines>
                  <siri:LineRef>STIF:Line::C00272:</siri:LineRef>
                  <siri:LineRef>STIF:Line::C00273:</siri:LineRef>
                  <siri:LineRef>STIF:Line::C00274:</siri:LineRef>
                </siri:Lines>
              </siri:AnnotatedStopPointRef>
              <siri:AnnotatedStopPointRef>
                <siri:StopPointRef>NINOXE:StopPoint:BP:33:LOC</siri:StopPointRef>
                <siri:Monitored>true</siri:Monitored>
                <siri:StopName>Stop3Referent</siri:StopName>
                <siri:Lines>
                  <siri:LineRef>STIF:Line::C00272:</siri:LineRef>
                  <siri:LineRef>STIF:Line::C00274:</siri:LineRef>
                </siri:Lines>
              </siri:AnnotatedStopPointRef>
              <siri:AnnotatedStopPointRef>
                <siri:StopPointRef>NINOXE:StopPoint:BP:55:LOC</siri:StopPointRef>
                <siri:Monitored>true</siri:Monitored>
                <siri:StopName>Stop5</siri:StopName>
                <siri:Lines>
                  <siri:LineRef>STIF:Line::C00272:</siri:LineRef>
                </siri:Lines>
              </siri:AnnotatedStopPointRef>
            </Answer>
            <AnswerExtension/>
          </sw:StopPointsDiscoveryResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 2464 3292 - Handle a SIRI StopPointsDiscovery request
    Given a Partner "test" exists with connectors [siri-stop-points-discovery-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_code_space | internal |
      | local_url            | address  |
    And a Line exists with the following attributes:
      | Name      | Line 1                          |
      | Codes | "internal":"STIF:Line::C00272:" |
    And a Line exists with the following attributes:
      | Name      | Line 2                          |
      | Codes | "internal":"STIF:Line::C00273:" |
    And a Line exists with the following attributes:
      | Name      | Line 3                          |
      | Codes | "internal":"STIF:Line::C00274:" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                                                      |
      | Codes | "internal": "NINOXE:StopPoint:BP:6:LOC"                                   |
      | Lines     | ["6ba7b814-9dad-11d1-2-00c04fd430c8","6ba7b814-9dad-11d1-3-00c04fd430c8"] |
    And a StopArea exists with the following attributes:
      | Name      | Test 2                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:16:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test 3                                  |
      | Codes | "internal": "NINOXE:StopPoint:BP:7:LOC" |
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
            <Answer version="2.0">
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
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

  Scenario: StopPointsDiscovery collect
    Given a SIRI server waits StopPointsDiscovery request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
      <siri:Status>true</siri:Status>
        <siri:AnnotatedStopPointRef>
          <siri:StopPointRef>NINOXE:StopPoint:BP:6:LOC</siri:StopPointRef>
          <siri:Monitored>true</siri:Monitored>
          <siri:StopName>Test</siri:StopName>
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
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-points-discovery-request-collector] and the following settings:
      | remote_url                 | http://localhost:8090      |
      | remote_credential          | test                       |
      | remote_code_space       | internal                   |
      | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
    And a minute has passed
    Then a StopArea "internal":"NINOXE:StopPoint:BP:6:LOC" should exist
    And a StopArea "internal":"NINOXE:StopPoint:BP:7:LOC" should exist

  Scenario: Collect Stop Areas discovered by StopPointsDiscovery (ARA-862)
    Given a SIRI server "A" waits StopPointsDiscovery request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
        <siri:Status>true</siri:Status>
        <siri:AnnotatedStopPointRef>
          <siri:StopPointRef>StopArea:A</siri:StopPointRef>
          <siri:Monitored>true</siri:Monitored>
          <siri:StopName>Stop A</siri:StopName>
        </siri:AnnotatedStopPointRef>
      </Answer>
      <AnswerExtension/>
    </sw:StopPointsDiscoveryResponse>
  </S:Body>
</S:Envelope>
      """
    And a SIRI server "B" waits StopPointsDiscovery request on "http://localhost:8091" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
        <siri:Status>true</siri:Status>
        <siri:AnnotatedStopPointRef>
          <siri:StopPointRef>StopArea:B</siri:StopPointRef>
          <siri:Monitored>true</siri:Monitored>
          <siri:StopName>Stop B</siri:StopName>
        </siri:AnnotatedStopPointRef>
      </Answer>
      <AnswerExtension/>
    </sw:StopPointsDiscoveryResponse>
  </S:Body>
</S:Envelope>
      """
    And the SIRI server "A" waits a GetStopMonitoring request to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>enRoute</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>enRoute:ResponseMessage::6ba7b814-9dad-11d1-e-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>enRoute:Message::2345Fsdfrg35df:LOC</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery>
          <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
          <siri:Status>true</siri:Status>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>2017-01-01T11:47:15.600+01:00</siri:RecordedAtTime>
            <siri:ItemIdentifier>enRoute:Item::63a1ebcfe85a0e9a548c91a611cfb572f4a545af:LOC</siri:ItemIdentifier>
            <siri:MonitoringRef>StopArea:A</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:MonitoredCall>
                <siri:StopPointRef>StopArea:A</siri:StopPointRef>
                <siri:Order>44</siri:Order>
                <siri:StopPointName>Stop A</siri:StopPointName>
                <siri:AimedArrivalTime>2017-01-01T13:43:05.000+01:00</siri:AimedArrivalTime>
                <siri:ExpectedArrivalTime>2017-01-01T13:43:05.000+01:00</siri:ExpectedArrivalTime>
                <siri:ArrivalStatus>onTime</siri:ArrivalStatus>
                <siri:AimedDepartureTime>2017-01-01T13:43:05.000+01:00</siri:AimedDepartureTime>
                <siri:ExpectedDepartureTime>2017-01-01T13:43:05.000+01:00</siri:ExpectedDepartureTime>
                <siri:DepartureStatus>onTime</siri:DepartureStatus>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "partner_a" exists with connectors [siri-check-status-client, siri-stop-points-discovery-request-collector, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                        | http://localhost:8090      |
      | remote_credential                 | test                       |
      | remote_code_space              | internal                   |
      | collect.use_discovered_stop_areas | true                       |
    And a Partner "partner_b" exists with connectors [siri-check-status-client, siri-stop-points-discovery-request-collector, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                        | http://localhost:8091      |
      | remote_credential                 | test                       |
      | remote_code_space              | internal                   |
      | collect.use_discovered_stop_areas | true                       |
    When 2 minutes have passed
    And 2 minutes have passed
    Then the "A" SIRI server should have received a GetStopMonitoring request with:
      | //siri:MonitoringRef | StopArea:A |
    And the "B" SIRI server should have received a GetStopMonitoring request with:
      | //siri:MonitoringRef | StopArea:B |

  @ARA-1030
  Scenario: Collect Stop Areas discovered by StopPointsDiscovery with one partner setting collect.exclude_stop_areas
    Given a SIRI server "A" waits StopPointsDiscovery request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:Status>true</siri:Status>
        <siri:AnnotatedStopPointRef>
          <siri:StopPointRef>StopArea:A</siri:StopPointRef>
          <siri:Monitored>true</siri:Monitored>
          <siri:StopName>Stop A</siri:StopName>
        </siri:AnnotatedStopPointRef>
      </Answer>
      <AnswerExtension/>
    </sw:StopPointsDiscoveryResponse>
  </S:Body>
</S:Envelope>
      """
    And a SIRI server "B" waits StopPointsDiscovery request on "http://localhost:8091" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:Status>true</siri:Status>
        <siri:AnnotatedStopPointRef>
          <siri:StopPointRef>StopArea:A</siri:StopPointRef>
          <siri:Monitored>true</siri:Monitored>
          <siri:StopName>Stop A</siri:StopName>
        </siri:AnnotatedStopPointRef>
      </Answer>
      <AnswerExtension/>
    </sw:StopPointsDiscoveryResponse>
  </S:Body>
</S:Envelope>
      """
    And the SIRI server "A" waits a GetStopMonitoring request to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>enRoute</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>enRoute:ResponseMessage::6ba7b814-9dad-11d1-e-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>enRoute:Message::2345Fsdfrg35df:LOC</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery>
          <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>enRoute:Message::2345Fsdfrg35df:LOC</siri:RequestMessageRef>
          <siri:MonitoringRef>StopArea:A</siri:MonitoringRef>
          <siri:Status>true</siri:Status>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>2017-01-01T11:47:15.600+01:00</siri:RecordedAtTime>
            <siri:ItemIdentifier>enRoute:Item::63a1ebcfe85a0e9a548c91a611cfb572f4a545af:LOC</siri:ItemIdentifier>
            <siri:MonitoringRef>StopArea:A</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:MonitoredCall>
                <siri:StopPointRef>StopArea:A</siri:StopPointRef>
                <siri:Order>44</siri:Order>
                <siri:StopPointName>Stop A</siri:StopPointName>
                <siri:AimedArrivalTime>2017-01-01T13:43:05.000+01:00</siri:AimedArrivalTime>
                <siri:ExpectedArrivalTime>2017-01-01T13:43:05.000+01:00</siri:ExpectedArrivalTime>
                <siri:ArrivalStatus>onTime</siri:ArrivalStatus>
                <siri:AimedDepartureTime>2017-01-01T13:43:05.000+01:00</siri:AimedDepartureTime>
                <siri:ExpectedDepartureTime>2017-01-01T13:43:05.000+01:00</siri:ExpectedDepartureTime>
                <siri:DepartureStatus>onTime</siri:DepartureStatus>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "partner_a" exists with connectors [siri-check-status-client, siri-stop-points-discovery-request-collector, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | remote_code_space              | internal              |
      | collect.use_discovered_stop_areas | true                  |
      | collect.priority                  | 2                     |
    And a Partner "partner_b" exists with connectors [siri-check-status-client, siri-stop-points-discovery-request-collector, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                        | http://localhost:8091 |
      | remote_credential                 | test                  |
      | remote_code_space              | internal              |
      | collect.exclude_stop_areas        | StopArea:A            |
      | collect.priority                  | 1                     |
      | collect.use_discovered_stop_areas | true                  |
    When a minute has passed
    Then one StopArea has the following attributes:
      | Codes | "internal": "StopArea:A" |
      | Monitored | true                     |


  @ARA-1030
  Scenario: Collect Stop Areas discovered by StopPointsDiscovery with one partner setting collect.exclude_lines
    Given a SIRI server "A" waits StopPointsDiscovery request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:Status>true</siri:Status>
        <siri:AnnotatedStopPointRef>
          <siri:StopPointRef>StopArea:A</siri:StopPointRef>
          <siri:Monitored>true</siri:Monitored>
          <siri:StopName>Stop A</siri:StopName>
          <siri:Lines>
          <siri:LineRef>Test:Line:2:LOC</siri:LineRef>
          </siri:Lines>
        </siri:AnnotatedStopPointRef>
      </Answer>
      <AnswerExtension/>
    </sw:StopPointsDiscoveryResponse>
  </S:Body>
</S:Envelope>
      """
    And a SIRI server "B" waits StopPointsDiscovery request on "http://localhost:8091" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:Status>true</siri:Status>
        <siri:AnnotatedStopPointRef>
          <siri:StopPointRef>StopArea:A</siri:StopPointRef>
          <siri:Monitored>true</siri:Monitored>
          <siri:StopName>Stop A</siri:StopName>
          <siri:Lines>
          <siri:LineRef>Test:Line:1:LOC</siri:LineRef>
          </siri:Lines>
        </siri:AnnotatedStopPointRef>
      </Answer>
      <AnswerExtension/>
      </sw:StopPointsDiscoveryResponse>
  </S:Body>
</S:Envelope>
      """
    And the SIRI server "A" waits a GetStopMonitoring request to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>enRoute</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>enRoute:ResponseMessage::6ba7b814-9dad-11d1-e-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>enRoute:Message::2345Fsdfrg35df:LOC</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery>
          <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>enRoute:Message::2345Fsdfrg35df:LOC</siri:RequestMessageRef>
          <siri:MonitoringRef>StopArea:A</siri:MonitoringRef>
          <siri:Status>true</siri:Status>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>2017-01-01T11:47:15.600+01:00</siri:RecordedAtTime>
            <siri:ItemIdentifier>enRoute:Item::63a1ebcfe85a0e9a548c91a611cfb572f4a545af:LOC</siri:ItemIdentifier>
            <siri:MonitoringRef>StopArea:A</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:MonitoredCall>
                <siri:StopPointRef>StopArea:A</siri:StopPointRef>
                <siri:Order>44</siri:Order>
                <siri:StopPointName>Stop A</siri:StopPointName>
                <siri:AimedArrivalTime>2017-01-01T13:43:05.000+01:00</siri:AimedArrivalTime>
                <siri:ExpectedArrivalTime>2017-01-01T13:43:05.000+01:00</siri:ExpectedArrivalTime>
                <siri:ArrivalStatus>onTime</siri:ArrivalStatus>
                <siri:AimedDepartureTime>2017-01-01T13:43:05.000+01:00</siri:AimedDepartureTime>
                <siri:ExpectedDepartureTime>2017-01-01T13:43:05.000+01:00</siri:ExpectedDepartureTime>
                <siri:DepartureStatus>onTime</siri:DepartureStatus>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "partner_a" exists with connectors [siri-check-status-client, siri-stop-points-discovery-request-collector, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | remote_code_space              | internal              |
      | collect.exclude_lines             | Test:line:1:LOC       |
      | collect.priority                  | 1                     |
      | collect.use_discovered_stop_areas | true                  |
    And a Partner "partner_b" exists with connectors [siri-check-status-client, siri-stop-points-discovery-request-collector, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                        | http://localhost:8091 |
      | remote_credential                 | test                  |
      | remote_code_space              | internal              |
      | collect.priority                  | 2                     |
      | collect.use_discovered_stop_areas | true                  |
    When a minute has passed
    Then one StopArea has the following attributes:
      | Codes | "internal": "StopArea:A"   |
      | Monitored | true                       |

  @ARA-1298 @siri-valid
  Scenario: StopPointsDiscovery collect with Partner remote_code_space changed
    Given a SIRI server waits StopPointsDiscovery request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
      <siri:Status>true</siri:Status>
        <siri:AnnotatedStopPointRef>
          <siri:StopPointRef>NINOXE:StopPoint:BP:6:LOC</siri:StopPointRef>
          <siri:Monitored>true</siri:Monitored>
          <siri:StopName>Test</siri:StopName>
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
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-points-discovery-request-collector] and the following settings:
      | remote_url                 | http://localhost:8090      |
      | remote_credential          | test                       |
      | remote_code_space       | internal                   |
    And a minute has passed
    Then a StopArea "internal":"NINOXE:StopPoint:BP:6:LOC" should exist
    And a StopArea "internal":"NINOXE:StopPoint:BP:7:LOC" should exist
    And the Partner "test" is updated with the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | remote_code_space | external              |
    And a SIRI server waits StopPointsDiscovery request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <Answer version="2.0">
      <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
      <siri:Status>true</siri:Status>
        <siri:AnnotatedStopPointRef>
          <siri:StopPointRef>NINOXE:StopPoint:BP:6:LOC</siri:StopPointRef>
          <siri:Monitored>true</siri:Monitored>
          <siri:StopName>Test</siri:StopName>
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
    And a minute has passed
    Then a StopArea "internal":"NINOXE:StopPoint:BP:6:LOC" should exist
    And a StopArea "internal":"NINOXE:StopPoint:BP:7:LOC" should exist
    Then a StopArea "external":"NINOXE:StopPoint:BP:6:LOC" should exist
    And a StopArea "external":"NINOXE:StopPoint:BP:7:LOC" should exist

  @ARA-1298 @siri-valid
  Scenario: Handle a SIRI StopPointsDiscovery request with Partner remote_code_space changed
    Given a Partner "test" exists with connectors [siri-stop-points-discovery-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_code_space | internal |
      | local_url            | address  |
    And a Line exists with the following attributes:
      | Name      | Line 1                          |
      | Codes | "internal":"STIF:Line::C00272:" |
    And a Line exists with the following attributes:
      | Name      | Line 2                          |
      | Codes | "internal":"STIF:Line::C00273:" |
    And a Line exists with the following attributes:
      | Name      | Line 3                          |
      | Codes | "external":"STIF:Line::C00274:"         |
    And a StopArea exists with the following attributes:
      | Name      | Test                                                                      |
      | Codes | "internal": "NINOXE:StopPoint:BP:6:LOC"                                   |
      | Lines     | ["6ba7b814-9dad-11d1-2-00c04fd430c8","6ba7b814-9dad-11d1-3-00c04fd430c8"] |
    And a StopArea exists with the following attributes:
      | Name      | Test 2                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:16:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test 3                                  |
      | Codes | "external": "NINOXE:StopPoint:BP:7:LOC" |
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
            <Answer version="2.0">
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
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
            </Answer>
            <AnswerExtension/>
          </sw:StopPointsDiscoveryResponse>
        </S:Body>
      </S:Envelope>
      """
    And the Partner "test" is updated with the following settings:
      | local_credential     | test     |
      | remote_code_space | external |
      | local_url            | address  |
    And a minute has passed
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
<?xml version='1.0' encoding='UTF-8'?>
<S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
  <S:Body>
    <sw:StopPointsDiscoveryResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
      <Answer version='2.0'>
        <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
        <siri:Status>true</siri:Status>
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

  @ARA-1545 @siri-valid
  Scenario: Handle a SIRI StopPointsDiscovery request with StopArea ReferentId and Lines ReferentId
    Given a Partner "test" exists with connectors [siri-stop-points-discovery-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | external |
      | local_url         | address  |
    And a Line exists with the following attributes:
      | Name  | Line Referent 1             |
      | Codes | "external":"Referent-1"  |
    And a Line exists with the following attributes:
      | Name       | Line 1                            |
      | Codes      | "internal":"Line-1"               |
      | ReferentId | 6ba7b814-9dad-11d1-2-00c04fd430c8 | # Line Referent 1
    And a Line exists with the following attributes:
      | Name  | Line Referent 2         |
      | Codes | "external":"Referent-2" |
    And a Line exists with the following attributes:
      | Name       | Line 2                            |
      | Codes      | "internal":"Line-2"               |
      | ReferentId | 6ba7b814-9dad-11d1-4-00c04fd430c8 | # Line Referent 2
    And a Line exists with the following attributes:
      | Name       | Line 3                            |
      | Codes      | "internal":"Line-3"               |
      | ReferentId | 6ba7b814-9dad-11d1-4-00c04fd430c8 | # Line Referent 2
    And a Line exists with the following attributes:
      | Name  | Line 4                                   |
      | Codes | "internal":"Line-4", "external":"Line-4" |
    And a StopArea exists with the following attributes:
      | Name  | Stop Referent                      |
      | Codes | "external": "Stop-Referent-1"  |
    And a StopArea exists with the following attributes:
      | Name       | Stop 1                                                                     |
      | Codes      | "internal": "Stop-1"                                                       |
      | Lines      | ["6ba7b814-9dad-11d1-3-00c04fd430c8", "6ba7b814-9dad-11d1-7-00c04fd430c8"] | # Line 1, Line 4
      | ReferentID | 6ba7b814-9dad-11d1-8-00c04fd430c8                                          | # Stop Referent
    And a StopArea exists with the following attributes:
      | Name       | Stop 2                                                                     |
      | Codes      | "internal": "Stop-2"                                                       |
      | Lines      | ["6ba7b814-9dad-11d1-5-00c04fd430c8", "6ba7b814-9dad-11d1-6-00c04fd430c8"] | # Line 2, Line 3
      | ReferentID | 6ba7b814-9dad-11d1-8-00c04fd430c8                                          | # Stop Referent
    And a StopArea exists with the following attributes:
      | Name  | Stop 3                                                                                                          |
      | Codes | "internal": "Stop-3", "external": "Stop-3"                                                                      |
      | Lines | ["6ba7b814-9dad-11d1-3-00c04fd430c8", "6ba7b814-9dad-11d1-5-00c04fd430c8", "6ba7b814-9dad-11d1-7-00c04fd430c8"] | # Line 1, Line 2, Line 4
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns7:StopPointsDiscovery xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <Request>
              <ns2:RequestTimestamp>2017-03-03T11:28:00.359Z</ns2:RequestTimestamp>
              <ns2:RequestorRef>test</ns2:RequestorRef>
              <ns2:MessageIdentifier>STIF:Message::2345Fsdfrg35df:LOC</ns2:MessageIdentifier>
            </Request>
            <RequestExtension/>
          </ns7:StopPointsDiscovery>
        </S:Body>
        </S:Envelope>
        """
    Then I should receive this SIRI response
        """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:StopPointsDiscoveryResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <Answer version='2.0'>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:Status>true</siri:Status>
              <siri:AnnotatedStopPointRef>
                <siri:StopPointRef>Stop-3</siri:StopPointRef>
                <siri:Monitored>true</siri:Monitored>
                <siri:StopName>Stop 3</siri:StopName>
                <siri:Lines>
                  <siri:LineRef>Line-4</siri:LineRef>
                  <siri:LineRef>Referent-1</siri:LineRef>
                  <siri:LineRef>Referent-2</siri:LineRef>
                </siri:Lines>
              </siri:AnnotatedStopPointRef>
              <siri:AnnotatedStopPointRef>
                <siri:StopPointRef>Stop-Referent-1</siri:StopPointRef>
                <siri:Monitored>true</siri:Monitored>
                <siri:StopName>Stop Referent</siri:StopName>
                <siri:Lines>
                  <siri:LineRef>Line-4</siri:LineRef>
                  <siri:LineRef>Referent-1</siri:LineRef>
                  <siri:LineRef>Referent-2</siri:LineRef>
                </siri:Lines>
              </siri:AnnotatedStopPointRef>
            </Answer>
            <AnswerExtension/>
          </sw:StopPointsDiscoveryResponse>
        </S:Body>
      </S:Envelope>
        """
