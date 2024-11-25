Feature: Support SIRI VehicleMonitoring by request

  Background:
      Given a Referential "test" is created

  @ARA-1306
  Scenario: VehicleMonitoring request collect should send GetVehicleMonitoring request to partner
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-vehicle-monitoring-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
      | local_credential      | ara                   |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name  | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
   And a minute has passed
   Then the SIRI server should have received 1 GetVehicleMonitoring request

  @ARA-1306
  Scenario: VehicleMonitoring request collect and partner CheckStatus is unavailable should not send GetVehicleMonitoring request to partner
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-vehicle-monitoring-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
      | local_credential      | ara                   |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name  | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
   And a minute has passed
   Then the SIRI server should not have received a GetVehicleMonitoring request

  @ARA-1306
  Scenario: VehicleMonitoring request collect and partner CheckStatus is unavailable should send GetVehicleMonitoring request to partner whith setting collect.persistent
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-vehicle-monitoring-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
      | local_credential      | ara                   |
      | collect.persistent    | true                  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name  | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
   And a minute has passed
   Then the SIRI server should have received 1 GetVehicleMonitoring request

  @siri-valid @ARA-1234
  Scenario: Handle a SIRI VehicleMonitoring request with fallback on generic connector remote_code_space
   Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential                                              | test     |
      | remote_code_space                                             | internal |
      | siri-vehicle-monitoring-request-broadcaster.remote_code_space | other    |
    Given a Line exists with the following attributes:
      | Codes | "other": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro              |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                             |
      | Codes                    | "other": "Test:VehicleJourney:201:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8      |
      | Monitored                | true                                   |
      | Attribute[DirectionName] | Direction Name                         |
    And a Vehicle exists with the following attributes:
      | Codes            | "other": "Test:Vehicle:201123:LOC" |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8  |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8  |
      | Longitude        | 1.234                              |
      | Latitude         | 5.678                              |
      | Bearing          | 123                                |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z           |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z           |
   When I send this SIRI request
     """
     <?xml version='1.0' encoding='UTF-8'?>
     <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceRequestInfo>
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>test</siri:RequestorRef>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
          </Request>
          <RequestExtension />
        </sw:GetVehicleMonitoring>
      </soap:Body>
    </soap:Envelope>
      """
    Then I should receive this SIRI response
    """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:VehicleActivity>
                  <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                  <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                  <siri:VehicleMonitoringRef>Test:Vehicle:201123:LOC</siri:VehicleMonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>Test:VehicleJourney:201:LOC</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DirectionName>Direction Name</siri:DirectionName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:VehicleLocation>
                      <siri:Longitude>1.234</siri:Longitude>
                      <siri:Latitude>5.678</siri:Latitude>
                    </siri:VehicleLocation>
                    <siri:Bearing>123</siri:Bearing>
                  </siri:MonitoredVehicleJourney>
                </siri:VehicleActivity>
              </siri:VehicleMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetVehicleMonitoringResponse>
        </S:Body>
      </S:Envelope>
    """

  @siri-valid @ARA-1234
  Scenario: Handle a SIRI VehicleMonitoring request with multiple connector setting siri-vehicle-monitoring-request-broadcaster.vehicle_journey_remote_code_space
   Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential                                                      | test          |
      | remote_code_space                                                     | internal      |
      | siri-vehicle-monitoring-request-broadcaster.vehicle_remote_code_space | other, other2 |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                |
      | Codes                    | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored                | true                                      |
      | Attribute[DirectionName] | Direction Name                            |
    And a Vehicle exists with the following attributes:
      | Codes            | "other": "Test:Vehicle:201123:LOC" |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8  |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8  |
      | Longitude        | 1.234                              |
      | Latitude         | 5.678                              |
      | Bearing          | 123                                |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z           |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z           |
      | Occupancy        | seatsAvailable                     |
   When I send this SIRI request
     """
     <?xml version='1.0' encoding='UTF-8'?>
     <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceRequestInfo>
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>test</siri:RequestorRef>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
          </Request>
          <RequestExtension />
        </sw:GetVehicleMonitoring>
      </soap:Body>
    </soap:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:VehicleActivity>
                  <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                  <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                  <siri:VehicleMonitoringRef>Test:Vehicle:201123:LOC</siri:VehicleMonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>Test:VehicleJourney:201:LOC</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DirectionName>Direction Name</siri:DirectionName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:VehicleLocation>
                      <siri:Longitude>1.234</siri:Longitude>
                      <siri:Latitude>5.678</siri:Latitude>
                    </siri:VehicleLocation>
                    <siri:Bearing>123</siri:Bearing>
                    <siri:Occupancy>seatsAvailable</siri:Occupancy>
                  </siri:MonitoredVehicleJourney>
                </siri:VehicleActivity>
              </siri:VehicleMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetVehicleMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

  @siri-valid @ARA-1234
  Scenario: Handle a SIRI VehicleMonitoring request with unmatching code kind
   Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
     | local_credential  | test  |
     | remote_code_space | wrong |
   Given a Line exists with the following attributes:
     | Codes | "internal": "Test:Line:3:LOC" |
     | Name  | Ligne 3 Metro                 |
   And a VehicleJourney exists with the following attributes:
     | Name                     | Passage 32                                |
     | Codes                    | "internal": "Test:VehicleJourney:201:LOC" |
     | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
     | Monitored                | true                                      |
     | Attribute[DirectionName] | Direction Name                            |
   And a Vehicle exists with the following attributes:
     | Codes            | "other": "Test:Vehicle:201123:LOC" |
     | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8  |
     | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8  |
     | Longitude        | 1.234                              |
     | Latitude         | 5.678                              |
     | Bearing          | 123                                |
     | RecordedAtTime   | 2017-01-01T13:00:00.000Z           |
     | ValidUntilTime   | 2017-01-01T14:00:00.000Z           |
   When I send this SIRI request
     """
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceRequestInfo>
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>test</siri:RequestorRef>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
          </Request>
          <RequestExtension />
        </sw:GetVehicleMonitoring>
      </soap:Body>
    </soap:Envelope>
      """
   Then I should receive this SIRI response
      """
      <?xml version="1.0" encoding="UTF-8"?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
                <siri:Status>false</siri:Status>
                <siri:ErrorCondition>
                  <siri:InvalidDataReferencesError>
                    <siri:ErrorText>Line Test:Line:3:LOC not found</siri:ErrorText>
                  </siri:InvalidDataReferencesError>
                </siri:ErrorCondition>
              </siri:VehicleMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetVehicleMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """
    Then an audit event should exist with these attributes:
        | Type              | VehicleMonitoringRequest       |
        | Protocol          | siri                           |
        | Direction         | received                       |
        | Status            | Error                          |
        | Partner           | test                           |
        | Vehicles          | []                             |
        | RequestIdentifier | Test:1234::LOC                 |
        | Lines             | ["Test:Line:3:LOC"]            |
        | ErrorDetails      | Line Test:Line:3:LOC not found |

  @siri-valid @ARA-1590
  Scenario: Handle a SIRI VehicleMonitoring request with Referent Line with Fallback on vehicle remote codeSpace (RatpCap case)
    Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential                                                           | test                |
      | remote_code_space                                                          | internal            |
      | siri-vehicle-monitoring-request-broadcaster.vehicle_remote_code_space      | rdmantois, rdbievre |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Referent-1" |
      | Name  | Line Referent 1          |
    # 6ba7b814-9dad-11d1-2-00c04fd430c8
    And a Line exists with the following attributes:
      | Codes      | "rdbievre": "Line-1"              |
      | Name       | Ligne 1                           |
      | ReferentId | 6ba7b814-9dad-11d1-2-00c04fd430c8 | # Line Referent 1
      # 6ba7b814-9dad-11d1-3-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                                                    |
      | Codes                    | "rdbievre": "bievre-VehicleJourney", "internal": "STIF:bievre-VehicleJourney" |
      | LineId                   | 6ba7b814-9dad-11d1-3-00c04fd430c8                                             |
      | Monitored                | true                                                                          |
      | Attribute[DirectionName] | Direction Name                                                                |
    # 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a Vehicle exists with the following attributes:
      | Codes            | "rdbievre": "bievre-Vehicle"      |
      | LineId           | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 123                               |
      | Occupancy        | seatsAvailable                    |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | DriverRef        | "1233"                            |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-7-00c04fd430c8 |
    # 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Codes      | "internal": "Stop-1"              |
      | Name       | Stop 1                            |
      | ReferentId | 6ba7b814-9dad-11d1-b-00c04fd430c8 | # Stop Referent
    # 6ba7b814-9dad-11d1-6-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "bievre-VehicleJourney-bievre-Vehicle" |
      | PassageOrder                  | 4                                                  |
      | VehicleAtStop                 | false                                              |
      | StopAreaId                    | 6ba7b814-9dad-11d1-6-00c04fd430c8                  |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-4-00c04fd430c8                  |
      | VehicleAtStop                 | false                                              |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                 |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                           |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                           |
      | ArrivalStatus                 | delayed                                            |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                           |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                           |
      | DepartureStatus               | delayed                                            |
      | Attribute[DestinationDisplay] | Pouet-pouet                                        |
    # 6ba7b814-9dad-11d1-7-00c04fd430c8
    And a Line exists with the following attributes:
      | Codes      | "rdmantois": "Line-2"             |
      | Name       | Line 2                            |
      | ReferentId | 6ba7b814-9dad-11d1-2-00c04fd430c8 | # Line Referent 1
    # 6ba7b814-9dad-11d1-8-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                                                        |
      | Codes                    | "rdmantois": "mantois-VehicleJourney", "internal": "STIF:mantois-VehicleJourney"  |
      | LineId                   | 6ba7b814-9dad-11d1-8-00c04fd430c8                                                 |
      | Monitored                | true                                                                              |
      | Attribute[DirectionName] | Another Direction Name                                                            |
    # 6ba7b814-9dad-11d1-9-00c04fd430c8
    And a Vehicle exists with the following attributes:
      | Codes            | "rdmantois": "mantois-Vehicle"    |
      | LineId           | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
      | Longitude        | 3.232                             |
      | Latitude         | 8.329                             |
      | Bearing          | 355                               |
      | Occupancy        | fewSeatsAvailable                 |
      | RecordedAtTime   | 2017-01-01T14:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T15:00:00.000Z          |
      | DriverRef        | "567"                             |
    # 6ba7b814-9dad-11d1-a-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Codes | "internal": "Stop-Referent-1" |
      | Name  | Stop Referent                 |
    When I send this SIRI request
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceRequestInfo>
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:RequestorRef>test</siri:RequestorRef>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            </ServiceRequestInfo>
            <Request version="2.0:FR-IDF-2.4">
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
              <siri:LineRef>Referent-1</siri:LineRef>
            </Request>
            <RequestExtension />
          </sw:GetVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then I should receive this SIRI response
    """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
          <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:ProducerRef>Ara</siri:ProducerRef>
          <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-c-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
          <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
        <Answer>
          <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
          <siri:Status>true</siri:Status>
          <siri:VehicleActivity>
            <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
            <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
            <siri:VehicleMonitoringRef>bievre-Vehicle</siri:VehicleMonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>Referent-1</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>STIF:bievre-VehicleJourney</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Line Referent 1</siri:PublishedLineName>
              <siri:DirectionName>Direction Name</siri:DirectionName>
              <siri:Monitored>true</siri:Monitored>
              <siri:VehicleLocation>
                <siri:Longitude>1.234</siri:Longitude>
                <siri:Latitude>5.678</siri:Latitude>
              </siri:VehicleLocation>
              <siri:Bearing>123</siri:Bearing>
              <siri:Occupancy>seatsAvailable</siri:Occupancy>
              <siri:DriverRef>1233</siri:DriverRef>
              <siri:MonitoredCall>
                <siri:StopPointRef>Stop-1</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Stop 1</siri:StopPointName>
                <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                <siri:ExpectedArrivalTime>2017-01-01T15:02:00.000Z</siri:ExpectedArrivalTime>
                <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                <siri:ExpectedDepartureTime>2017-01-01T15:01:00.000Z</siri:ExpectedDepartureTime>
                <siri:DepartureStatus>delayed</siri:DepartureStatus>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:VehicleActivity>
          <siri:VehicleActivity>
            <siri:RecordedAtTime>2017-01-01T14:00:00.000Z</siri:RecordedAtTime>
            <siri:ValidUntilTime>2017-01-01T15:00:00.000Z</siri:ValidUntilTime>
            <siri:VehicleMonitoringRef>mantois-Vehicle</siri:VehicleMonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>Referent-1</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>STIF:mantois-VehicleJourney</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Line Referent 1</siri:PublishedLineName>
              <siri:DirectionName>Another Direction Name</siri:DirectionName>
              <siri:Monitored>true</siri:Monitored>
              <siri:VehicleLocation>
                <siri:Longitude>3.232</siri:Longitude>
                <siri:Latitude>8.329</siri:Latitude>
              </siri:VehicleLocation>
              <siri:Bearing>355</siri:Bearing>
              <siri:Occupancy>fewSeatsAvailable</siri:Occupancy>
              <siri:DriverRef>567</siri:DriverRef>
            </siri:MonitoredVehicleJourney>
          </siri:VehicleActivity>
          </siri:VehicleMonitoringDelivery>
        </Answer>
        <AnswerExtension/>
          </sw:GetVehicleMonitoringResponse>
        </S:Body>
      </S:Envelope>
    """
    Then an audit event should exist with these attributes:
      | Type              | VehicleMonitoringRequest                                      |
      | Protocol          | siri                                                          |
      | Direction         | received                                                      |
      | Status            | OK                                                            |
      | Partner           | test                                                          |
      | Vehicles          | ["bievre-Vehicle", "mantois-Vehicle"]                         |
      | RequestIdentifier | Test:1234::LOC                                                |
      | Lines             | ["Referent-1"]                                                |
      | VehicleJourneys   | ["STIF:bievre-VehicleJourney", "STIF:mantois-VehicleJourney"] |

  @siri-valid @ARA-1234
  Scenario: Send all the vehicles in respond to a SIRI VehicleMonitoring request
    Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential      | test     |
      | remote_code_space     | internal |
      | sort_payload_for_test | true     |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                                |
      | Codes     | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored | true                                      |
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 33                                |
      | Codes                             | "internal": "Test:VehicleJourney:202:LOC" |
      | LineId                            | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored                         | true                                      |
      | Reference[DestinationRef]#Code    | "internal": "Test:StopPoint:Destination"  |
      | Reference[JourneyPatternRef]#Code | "internal": "Test:JourneyPattern:1"       |
      | Reference[OriginRef]#Code         | "internal": "Test:StopPoint:Origin"       |
      | OriginName                        | Origin Name                               |
      | DestinationName                   | Destination Name                          |
      | DirectionName                     | Direction Name                            |
      | DirectionType                     | outbound                                  |
      | Attribute[JourneyPatternName]     | Journey Pattern Name                      |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:1:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | DriverRef        | Driver1                           |
      | Bearing          | 120                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | LinkDistance     | 12                                |
      | Percentage       | 42                                |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:2:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | DriverRef        | Driver2                           |
      | Bearing          | 153                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | LinkDistance     | 34                                |
      | Percentage       | 55                                |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:3:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | DriverRef        | Driver3                           |
      | Bearing          | 163                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | LinkDistance     | 56                                |
      | Percentage       | 21                                |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Name  | Carabacel                                |
      # 6ba7b814-9dad-11d1-8-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "Test:VehicleJourney:202:LOC-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                  | 4                                                                      |
      | VehicleAtStop                 | false                                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-8-00c04fd430c8                                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-4-00c04fd430c8                                      |
      | VehicleAtStop                 | false                                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                                               |
      | ArrivalStatus                 | delayed                                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                                               |
      | DepartureStatus               | delayed                                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                                            |
      # 6ba7b814-9dad-11d1-9-00c04fd430c8
    When I send this SIRI request
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceRequestInfo>
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:RequestorRef>test</siri:RequestorRef>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            </ServiceRequestInfo>
            <Request version="2.0:FR-IDF-2.4">
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
              <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
            </Request>
            <RequestExtension />
          </sw:GetVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then I should receive this SIRI response
    """
    <?xml version='1.0' encoding='UTF-8'?>
    <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
      <S:Body>
        <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:ProducerRef>Ara</siri:ProducerRef>
            <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Answer>
            <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
              <siri:Status>true</siri:Status>
              <siri:VehicleActivity>
                <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                <siri:VehicleMonitoringRef>Test:Vehicle:1:LOC</siri:VehicleMonitoringRef>
                <siri:ProgressBetweenStops>
                  <siri:LinkDistance>12</siri:LinkDistance>
                  <siri:Percentage>42</siri:Percentage>
                </siri:ProgressBetweenStops>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>Test:VehicleJourney:201:LOC</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:VehicleLocation>
                    <siri:Longitude>1.234</siri:Longitude>
                    <siri:Latitude>5.678</siri:Latitude>
                  </siri:VehicleLocation>
                  <siri:Bearing>120</siri:Bearing>
                  <siri:DriverRef>Driver1</siri:DriverRef>
                </siri:MonitoredVehicleJourney>
              </siri:VehicleActivity>
              <siri:VehicleActivity>
                <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                <siri:VehicleMonitoringRef>Test:Vehicle:2:LOC</siri:VehicleMonitoringRef>
                <siri:ProgressBetweenStops>
                  <siri:LinkDistance>34</siri:LinkDistance>
                  <siri:Percentage>55</siri:Percentage>
                </siri:ProgressBetweenStops>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>Test:VehicleJourney:201:LOC</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:VehicleLocation>
                    <siri:Longitude>1.234</siri:Longitude>
                    <siri:Latitude>5.678</siri:Latitude>
                  </siri:VehicleLocation>
                  <siri:Bearing>153</siri:Bearing>
                  <siri:DriverRef>Driver2</siri:DriverRef>
                </siri:MonitoredVehicleJourney>
              </siri:VehicleActivity>
              <siri:VehicleActivity>
                <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                <siri:VehicleMonitoringRef>Test:Vehicle:3:LOC</siri:VehicleMonitoringRef>
                <siri:ProgressBetweenStops>
                  <siri:LinkDistance>56</siri:LinkDistance>
                  <siri:Percentage>21</siri:Percentage>
                </siri:ProgressBetweenStops>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>Test:VehicleJourney:202:LOC</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:JourneyPatternRef>Test:JourneyPattern:1</siri:JourneyPatternRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:OriginRef>RATPDev:StopPoint:Q:488317b5b41cb7ba0a4812c18b312f0e2b986852:LOC</siri:OriginRef>
                  <siri:OriginName>Origin Name</siri:OriginName>
                  <siri:DestinationRef>RATPDev:StopPoint:Q:7bef317e38443efe7d8e8e7f3b7b59881b2e3be0:LOC</siri:DestinationRef>
                  <siri:DestinationName>Destination Name</siri:DestinationName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:VehicleLocation>
                    <siri:Longitude>1.234</siri:Longitude>
                    <siri:Latitude>5.678</siri:Latitude>
                  </siri:VehicleLocation>
                  <siri:Bearing>163</siri:Bearing>
                  <siri:DriverRef>Driver3</siri:DriverRef>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Carabacel</siri:StopPointName>
                      <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                      <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                      <siri:ExpectedArrivalTime>2017-01-01T15:02:00.000Z</siri:ExpectedArrivalTime>
                      <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                      <siri:ExpectedDepartureTime>2017-01-01T15:01:00.000Z</siri:ExpectedDepartureTime>
                      <siri:DepartureStatus>delayed</siri:DepartureStatus>
                    </siri:MonitoredCall>
                </siri:MonitoredVehicleJourney>
              </siri:VehicleActivity>
            </siri:VehicleMonitoringDelivery>
          </Answer>
          <AnswerExtension/>
        </sw:GetVehicleMonitoringResponse>
      </S:Body>
    </S:Envelope>
    """
    Then an audit event should exist with these attributes:
        | Type              | VehicleMonitoringRequest                                           |
        | Protocol          | siri                                                               |
        | Direction         | received                                                           |
        | Status            | OK                                                                 |
        | Partner           | test                                                               |
        | Vehicles          | ["Test:Vehicle:1:LOC", "Test:Vehicle:2:LOC", "Test:Vehicle:3:LOC"] |
        | RequestIdentifier | Test:1234::LOC                                                     |
        | Lines             | ["Test:Line:3:LOC"]                                                |
        | VehicleJourneys   | ["Test:VehicleJourney:202:LOC", "Test:VehicleJourney:201:LOC"]     |

  @siri-valid @ARA-1384
  Scenario: Handle a SIRI VehicleMonitoring request with Vehicle filter
    Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential      | test     |
      | remote_code_space     | internal |
      | sort_payload_for_test | true     |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                                |
      | Codes     | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored | true                                      |
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 33                                |
      | Codes                             | "internal": "Test:VehicleJourney:202:LOC" |
      | LineId                            | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored                         | true                                      |
      | Reference[DestinationRef]#Code    | "internal": "Test:StopPoint:Destination"  |
      | Reference[JourneyPatternRef]#Code | "internal": "Test:JourneyPattern:1"       |
      | Reference[OriginRef]#Code         | "internal": "Test:StopPoint:Origin"       |
      | OriginName                        | Origin Name                               |
      | DestinationName                   | Destination Name                          |
      | DirectionName                     | Direction Name                            |
      | DirectionType                     | outbound                                  |
      | Attribute[JourneyPatternName]     | Journey Pattern Name                      |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:1:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | DriverRef        | Driver1                           |
      | Bearing          | 120                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | LinkDistance     | 12                                |
      | Percentage       | 42                                |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:2:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | DriverRef        | Driver2                           |
      | Bearing          | 153                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | LinkDistance     | 34                                |
      | Percentage       | 55                                |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:3:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | DriverRef        | Driver3                           |
      | Bearing          | 163                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | LinkDistance     | 56                                |
      | Percentage       | 21                                |
    When I send this SIRI request
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceRequestInfo>
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:RequestorRef>test</siri:RequestorRef>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            </ServiceRequestInfo>
            <Request version="2.0:FR-IDF-2.4">
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
              <siri:VehicleRef>Test:Vehicle:1:LOC</siri:VehicleRef>
            </Request>
            <RequestExtension />
          </sw:GetVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then I should receive this SIRI response
    """
    <?xml version='1.0' encoding='UTF-8'?>
    <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
      <S:Body>
        <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:ProducerRef>Ara</siri:ProducerRef>
            <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Answer>
            <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
              <siri:Status>true</siri:Status>
              <siri:VehicleActivity>
                <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                <siri:VehicleMonitoringRef>Test:Vehicle:1:LOC</siri:VehicleMonitoringRef>
                <siri:ProgressBetweenStops>
                  <siri:LinkDistance>12</siri:LinkDistance>
                  <siri:Percentage>42</siri:Percentage>
                </siri:ProgressBetweenStops>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>Test:VehicleJourney:201:LOC</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:VehicleLocation>
                    <siri:Longitude>1.234</siri:Longitude>
                    <siri:Latitude>5.678</siri:Latitude>
                  </siri:VehicleLocation>
                  <siri:Bearing>120</siri:Bearing>
                  <siri:DriverRef>Driver1</siri:DriverRef>
                </siri:MonitoredVehicleJourney>
              </siri:VehicleActivity>
             </siri:VehicleMonitoringDelivery>
          </Answer>
          <AnswerExtension/>
        </sw:GetVehicleMonitoringResponse>
      </S:Body>
    </S:Envelope>
    """
    Then an audit event should exist with these attributes:
        | Type              | VehicleMonitoringRequest        |
        | Protocol          | siri                            |
        | Direction         | received                        |
        | Status            | OK                              |
        | Partner           | test                            |
        | Vehicles          | ["Test:Vehicle:1:LOC"]          |
        | RequestIdentifier | Test:1234::LOC                  |
        | Lines             | ["Test:Line:3:LOC"]             |
        | VehicleJourneys   | ["Test:VehicleJourney:201:LOC"] |

  @siri-valid @ARA-1384
  Scenario: Handle a SIRI VehicleMonitoring request with Vehicle filter with unmatching code kind
   Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
     | local_credential  | test     |
     | remote_code_space | internal |
   Given a Line exists with the following attributes:
     | Codes | "internal": "Test:Line:3:LOC" |
     | Name  | Ligne 3 Metro                 |
   And a VehicleJourney exists with the following attributes:
     | Name                     | Passage 32                                |
     | Codes                    | "internal": "Test:VehicleJourney:201:LOC" |
     | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
     | Monitored                | true                                      |
     | Attribute[DirectionName] | Direction Name                            |
   And a Vehicle exists with the following attributes:
     | Codes            | "other": "Test:Vehicle:201123:LOC" |
     | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8  |
     | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8  |
     | Longitude        | 1.234                              |
     | Latitude         | 5.678                              |
     | Bearing          | 123                                |
     | RecordedAtTime   | 2017-01-01T13:00:00.000Z           |
     | ValidUntilTime   | 2017-01-01T14:00:00.000Z           |
   When I send this SIRI request
     """
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceRequestInfo>
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>test</siri:RequestorRef>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            <siri:VehicleRef>Test:Vehicle:201123:LOC</siri:VehicleRef>
          </Request>
          <RequestExtension />
        </sw:GetVehicleMonitoring>
      </soap:Body>
    </soap:Envelope>
      """
   Then I should receive this SIRI response
      """
      <?xml version="1.0" encoding="UTF-8"?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
                <siri:Status>false</siri:Status>
                <siri:ErrorCondition>
                  <siri:InvalidDataReferencesError>
                    <siri:ErrorText>Vehicle Test:Vehicle:201123:LOC not found</siri:ErrorText>
                  </siri:InvalidDataReferencesError>
                </siri:ErrorCondition>
              </siri:VehicleMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetVehicleMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """
    Then an audit event should exist with these attributes:
        | Type              | VehicleMonitoringRequest                  |
        | Protocol          | siri                                      |
        | Direction         | received                                  |
        | Status            | Error                                     |
        | Partner           | test                                      |
        | Vehicles          | ["Test:Vehicle:201123:LOC"]               |
        | RequestIdentifier | Test:1234::LOC                            |
        | Lines             | []                                        |
        | ErrorDetails      | Vehicle Test:Vehicle:201123:LOC not found |

  @siri-valid @ARA-1384
  Scenario: Handle a SIRI VehicleMonitoring request without Vehicle or Line filter should return an Error
   Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
     | local_credential  | test     |
     | remote_code_space | internal |
   Given a Line exists with the following attributes:
     | Codes | "internal": "Test:Line:3:LOC" |
     | Name  | Ligne 3 Metro                 |
   And a VehicleJourney exists with the following attributes:
     | Name                     | Passage 32                                |
     | Codes                    | "internal": "Test:VehicleJourney:201:LOC" |
     | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
     | Monitored                | true                                      |
     | Attribute[DirectionName] | Direction Name                            |
   And a Vehicle exists with the following attributes:
     | Codes            | "internal": "Test:Vehicle:201123:LOC" |
     | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8     |
     | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8     |
     | Longitude        | 1.234                                 |
     | Latitude         | 5.678                                 |
     | Bearing          | 123                                   |
     | RecordedAtTime   | 2017-01-01T13:00:00.000Z              |
     | ValidUntilTime   | 2017-01-01T14:00:00.000Z              |
   When I send this SIRI request
     """
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceRequestInfo>
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>test</siri:RequestorRef>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
          </Request>
          <RequestExtension />
        </sw:GetVehicleMonitoring>
      </soap:Body>
    </soap:Envelope>
      """
   Then I should receive this SIRI response
      """
      <?xml version="1.0" encoding="UTF-8"?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
                <siri:Status>false</siri:Status>
                <siri:ErrorCondition>
                  <siri:InvalidDataReferencesError>
                    <siri:ErrorText>VehicleMonitoringRequest must have one LineRef OR one VehicleRef</siri:ErrorText>
                  </siri:InvalidDataReferencesError>
                </siri:ErrorCondition>
              </siri:VehicleMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetVehicleMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """
    Then an audit event should exist with these attributes:
        | Type              | VehicleMonitoringRequest                                         |
        | Protocol          | siri                                                             |
        | Direction         | received                                                         |
        | Status            | Error                                                            |
        | Partner           | test                                                             |
        | Vehicles          | []                                                               |
        | RequestIdentifier | Test:1234::LOC                                                   |
        | Lines             | []                                                               |
        | ErrorDetails      | VehicleMonitoringRequest must have one LineRef OR one VehicleRef |

  @siri-valid @ARA-1234
  Scenario: Handle a SIRI VehicleMonitoring request with Vehicle filter with fallback on generic connector remote_code_space
   Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential                                              | test     |
      | remote_code_space                                             | internal |
      | siri-vehicle-monitoring-request-broadcaster.remote_code_space | other    |
    Given a Line exists with the following attributes:
      | Codes | "other": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro              |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                             |
      | Codes                    | "other": "Test:VehicleJourney:201:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8      |
      | Monitored                | true                                   |
      | Attribute[DirectionName] | Direction Name                         |
    And a Vehicle exists with the following attributes:
      | Codes            | "other": "Test:Vehicle:201123:LOC" |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8  |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8  |
      | Longitude        | 1.234                              |
      | Latitude         | 5.678                              |
      | Bearing          | 123                                |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z           |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z           |
   When I send this SIRI request
     """
     <?xml version='1.0' encoding='UTF-8'?>
     <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceRequestInfo>
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>test</siri:RequestorRef>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
            <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            <siri:VehicleRef>Test:Vehicle:201123:LOC</siri:VehicleRef>
          </Request>
          <RequestExtension />
        </sw:GetVehicleMonitoring>
      </soap:Body>
    </soap:Envelope>
      """
    Then I should receive this SIRI response
    """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:VehicleActivity>
                  <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                  <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                  <siri:VehicleMonitoringRef>Test:Vehicle:201123:LOC</siri:VehicleMonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>Test:VehicleJourney:201:LOC</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DirectionName>Direction Name</siri:DirectionName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:VehicleLocation>
                      <siri:Longitude>1.234</siri:Longitude>
                      <siri:Latitude>5.678</siri:Latitude>
                    </siri:VehicleLocation>
                    <siri:Bearing>123</siri:Bearing>
                  </siri:MonitoredVehicleJourney>
                </siri:VehicleActivity>
              </siri:VehicleMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetVehicleMonitoringResponse>
        </S:Body>
      </S:Envelope>
    """

  Scenario: Performs a SIRI VehicleMonitoring request to a Partner
    Given a SIRI server waits GetVehicleMonitoring request on "http://localhost:8090" to respond with
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns1:GetVehicleMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>RLA_Bus</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>RLA_Bus:ResponseMessage::23833:LOC</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
          <ns5:Status>true</ns5:Status>
          <ns5:VehicleActivity>
            <ns5:RecordedAtTime>2021-08-02T08:50:27.733+02:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>290</ns5:ItemIdentifier>
            <ns5:ValidUntilTime>2021-08-02T09:50:27.733+02:00</ns5:ValidUntilTime>
            <ns5:VehicleMonitoringRef>290</ns5:VehicleMonitoringRef>
            <ns5:ProgressBetweenStops>
              <ns5:LinkDistance>349.0</ns5:LinkDistance>
              <ns5:Percentage>70.0</ns5:Percentage>
            </ns5:ProgressBetweenStops>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>RLA_Bus:Line::05:LOC</ns5:LineRef>
              <ns5:DirectionRef>Aller</ns5:DirectionRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>RLA_Bus:DataFrame::1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>RLA_Bus:VehicleJourney::2978464:LOC</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>RLA_Bus:JourneyPattern::L05P99:LOC</ns5:JourneyPatternRef>
              <ns5:JourneyPatternName>L05P99</ns5:JourneyPatternName>
              <ns5:PublishedLineName>05</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>RLA_Bus:Operator::RLA:LOC</ns5:OperatorRef>
              <ns5:OriginRef>RLA_Bus:StopPoint:BP:DELOY0:LOC</ns5:OriginRef>
              <ns5:OriginName>Deloye / Dubouchage</ns5:OriginName>
              <ns5:DestinationRef>RLA_Bus:StopPoint:BP:RIMIE9:LOC</ns5:DestinationRef>
              <ns5:DestinationName>Rimiez Saint-George</ns5:DestinationName>
              <ns5:Monitored>false</ns5:Monitored>
              <ns5:VehicleLocation srsName="EPSG:2154">
                <ns5:Coordinates>1044593 6298716</ns5:Coordinates>
              </ns5:VehicleLocation>
              <ns5:Bearing>287.0</ns5:Bearing>
              <ns5:VehicleRef>RLA290</ns5:VehicleRef>
              <ns5:DriverRef>5753</ns5:DriverRef>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>RLA_Bus:StopPoint:BP:PASTO8:LOC</ns5:StopPointRef>
                <ns5:Order>6</ns5:Order>
                <ns5:StopPointName>Carabacel</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Rimiez Saint-George</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2021-08-02T07:38:42.000+02:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>delayed</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2021-08-02T07:38:42.000+02:00</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>delayed</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney><ns5:Extensions/></ns5:VehicleActivity>
        </ns5:VehicleMonitoringDelivery>
      </Answer><AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/></ns1:GetVehicleMonitoringResponse>
  </soap:Body>
</soap:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-vehicle-monitoring-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name  | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
    When a minute has passed
    And the SIRI server has received a GetVehicleMonitoring request
    Then one StopArea has the following attributes:
      | Codes | "internal": "RLA_Bus:StopPoint:BP:PASTO8:LOC" |
      | Name  | Carabacel                                     |
    And one Line has the following attributes:
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
      | Name  | Test 1                             |
    And one VehicleJourney has the following attributes:
      | Codes                             | "internal": "RLA_Bus:VehicleJourney::2978464:LOC" |
      | Reference[OriginRef]#Code         | "internal": "RLA_Bus:StopPoint:BP:DELOY0:LOC"     |
      | OriginName                        | Deloye / Dubouchage                               |
      | Reference[DestinationRef]#Code    | "internal": "RLA_Bus:StopPoint:BP:RIMIE9:LOC"     |
      | DestinationName                   | Rimiez Saint-George                               |
      | Reference[JourneyPatternRef]#Code | "internal": "RLA_Bus:JourneyPattern::L05P99:LOC"  |
      | Monitored                         | false                                             |
    And one Vehicle has the following attributes:
      | Codes          | "internal": "RLA290"          |
      | Longitude      | 7.276192074052043             |
      | Latitude       | 43.70347861870634             |
      | DriverRef      | "5753"                        |
      | Bearing        | 287.0                         |
      | LinkDistance   | 349.0                         |
      | Percentage     | 70.0                          |
      | ValidUntilTime | 2021-08-02T08:50:27.733+02:00 |
    And an audit event should exist with these attributes:
      | Protocol        | siri                                    |
      | Direction       | sent                                    |
      | Status          | OK                                      |
      | Type            | VehicleMonitoringRequest                |
      | StopAreas       | ["RLA_Bus:StopPoint:BP:PASTO8:LOC"]     |
      | VehicleJourneys | ["RLA_Bus:VehicleJourney::2978464:LOC"] |
      | Lines           | ["RLA_Bus:Line::05:LOC"]                |
      | Vehicles        | ["RLA290"]                              |

  Scenario: Collect Vehicle Position with numeric srsName
    Given a SIRI server waits GetVehicleMonitoring request on "http://localhost:8090" to respond with
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns1:GetVehicleMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>RLA_Bus</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>RLA_Bus:ResponseMessage::23833:LOC</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
          <ns5:Status>true</ns5:Status>
          <ns5:VehicleActivity>
            <ns5:RecordedAtTime>2021-08-02T08:50:27.733+02:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>290</ns5:ItemIdentifier>
            <ns5:ValidUntilTime>2021-08-02T09:50:27.733+02:00</ns5:ValidUntilTime>
            <ns5:VehicleMonitoringRef>290</ns5:VehicleMonitoringRef>
            <ns5:ProgressBetweenStops>
              <ns5:LinkDistance>349.0</ns5:LinkDistance>
              <ns5:Percentage>70.0</ns5:Percentage>
            </ns5:ProgressBetweenStops>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>RLA_Bus:Line::05:LOC</ns5:LineRef>
              <ns5:DirectionRef>Aller</ns5:DirectionRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>RLA_Bus:DataFrame::1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>RLA_Bus:VehicleJourney::2978464:LOC</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>RLA_Bus:JourneyPattern::L05P99:LOC</ns5:JourneyPatternRef>
              <ns5:JourneyPatternName>L05P99</ns5:JourneyPatternName>
              <ns5:PublishedLineName>05</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>RLA_Bus:Operator::RLA:LOC</ns5:OperatorRef>
              <ns5:OriginRef>RLA_Bus:StopPoint:BP:DELOY0:LOC</ns5:OriginRef>
              <ns5:OriginName>Deloye / Dubouchage</ns5:OriginName>
              <ns5:DestinationRef>RLA_Bus:StopPoint:BP:RIMIE9:LOC</ns5:DestinationRef>
              <ns5:DestinationName>Rimiez Saint-George</ns5:DestinationName>
              <ns5:Monitored>false</ns5:Monitored>
              <ns5:VehicleLocation srsName="2154">
                <ns5:Coordinates>1044593 6298716</ns5:Coordinates>
              </ns5:VehicleLocation>
              <ns5:Bearing>287.0</ns5:Bearing>
              <ns5:VehicleRef>RLA290</ns5:VehicleRef>
              <ns5:DriverRef>5753</ns5:DriverRef>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>RLA_Bus:StopPoint:BP:PASTO8:LOC</ns5:StopPointRef>
                <ns5:Order>6</ns5:Order>
                <ns5:StopPointName>Carabacel</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Rimiez Saint-George</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2021-08-02T07:38:42.000+02:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>delayed</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2021-08-02T07:38:42.000+02:00</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>delayed</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney><ns5:Extensions/></ns5:VehicleActivity>
        </ns5:VehicleMonitoringDelivery>
      </Answer><AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/></ns1:GetVehicleMonitoringResponse>
  </soap:Body>
</soap:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-vehicle-monitoring-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name  | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
    When a minute has passed
    And the SIRI server has received a GetVehicleMonitoring request
    Then one StopArea has the following attributes:
      | Codes | "internal": "RLA_Bus:StopPoint:BP:PASTO8:LOC" |
      | Name  | Carabacel                                     |
    And one Line has the following attributes:
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
      | Name  | Test 1                             |
    And one VehicleJourney has the following attributes:
      | Codes                             | "internal": "RLA_Bus:VehicleJourney::2978464:LOC" |
      | Reference[OriginRef]#Code         | "internal": "RLA_Bus:StopPoint:BP:DELOY0:LOC"     |
      | OriginName                        | Deloye / Dubouchage                               |
      | Reference[DestinationRef]#Code    | "internal": "RLA_Bus:StopPoint:BP:RIMIE9:LOC"     |
      | DestinationName                   | Rimiez Saint-George                               |
      | Reference[JourneyPatternRef]#Code | "internal": "RLA_Bus:JourneyPattern::L05P99:LOC"  |
      | Monitored                         | false                                             |
    And one Vehicle has the following attributes:
      | Codes          | "internal": "RLA290"          |
      | Longitude      | 7.276192074052043             |
      | Latitude       | 43.70347861870634             |
      | DriverRef      | "5753"                        |
      | Bearing        | 287.0                         |
      | LinkDistance   | 349.0                         |
      | Percentage     | 70.0                          |
      | ValidUntilTime | 2021-08-02T08:50:27.733+02:00 |

  Scenario: Collect Vehicle Position with collect.default_srs_name setting
    Given a SIRI server waits GetVehicleMonitoring request on "http://localhost:8090" to respond with
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns1:GetVehicleMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>RLA_Bus</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>RLA_Bus:ResponseMessage::23833:LOC</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
          <ns5:Status>true</ns5:Status>
          <ns5:VehicleActivity>
            <ns5:RecordedAtTime>2021-08-02T08:50:27.733+02:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>290</ns5:ItemIdentifier>
            <ns5:ValidUntilTime>2021-08-02T09:50:27.733+02:00</ns5:ValidUntilTime>
            <ns5:VehicleMonitoringRef>290</ns5:VehicleMonitoringRef>
            <ns5:ProgressBetweenStops>
              <ns5:LinkDistance>349.0</ns5:LinkDistance>
              <ns5:Percentage>70.0</ns5:Percentage>
            </ns5:ProgressBetweenStops>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>RLA_Bus:Line::05:LOC</ns5:LineRef>
              <ns5:DirectionRef>Aller</ns5:DirectionRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>RLA_Bus:DataFrame::1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>RLA_Bus:VehicleJourney::2978464:LOC</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>RLA_Bus:JourneyPattern::L05P99:LOC</ns5:JourneyPatternRef>
              <ns5:JourneyPatternName>L05P99</ns5:JourneyPatternName>
              <ns5:PublishedLineName>05</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>RLA_Bus:Operator::RLA:LOC</ns5:OperatorRef>
              <ns5:OriginRef>RLA_Bus:StopPoint:BP:DELOY0:LOC</ns5:OriginRef>
              <ns5:OriginName>Deloye / Dubouchage</ns5:OriginName>
              <ns5:DestinationRef>RLA_Bus:StopPoint:BP:RIMIE9:LOC</ns5:DestinationRef>
              <ns5:DestinationName>Rimiez Saint-George</ns5:DestinationName>
              <ns5:Monitored>false</ns5:Monitored>
              <ns5:VehicleLocation>
                <ns5:Coordinates>1044593 6298716</ns5:Coordinates>
              </ns5:VehicleLocation>
              <ns5:Bearing>287.0</ns5:Bearing>
              <ns5:VehicleRef>RLA290</ns5:VehicleRef>
              <ns5:DriverRef>5753</ns5:DriverRef>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>RLA_Bus:StopPoint:BP:PASTO8:LOC</ns5:StopPointRef>
                <ns5:Order>6</ns5:Order>
                <ns5:StopPointName>Carabacel</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Rimiez Saint-George</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2021-08-02T07:38:42.000+02:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>delayed</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2021-08-02T07:38:42.000+02:00</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>delayed</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney><ns5:Extensions/></ns5:VehicleActivity>
        </ns5:VehicleMonitoringDelivery>
      </Answer><AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/></ns1:GetVehicleMonitoringResponse>
  </soap:Body>
</soap:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-vehicle-monitoring-request-collector] and the following settings:
      | remote_url               | http://localhost:8090 |
      | remote_credential        | test                  |
      | remote_code_space        | internal              |
      | collect.include_lines    | RLA_Bus:Line::05:LOC  |
      | collect.default_srs_name | EPSG:2154             |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name  | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
    When a minute has passed
    And the SIRI server has received a GetVehicleMonitoring request
    Then one StopArea has the following attributes:
      | Codes | "internal": "RLA_Bus:StopPoint:BP:PASTO8:LOC" |
      | Name  | Carabacel                                     |
    And one Line has the following attributes:
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
      | Name  | Test 1                             |
    And one VehicleJourney has the following attributes:
      | Codes                             | "internal": "RLA_Bus:VehicleJourney::2978464:LOC" |
      | Reference[OriginRef]#Code         | "internal": "RLA_Bus:StopPoint:BP:DELOY0:LOC"     |
      | OriginName                        | Deloye / Dubouchage                               |
      | Reference[DestinationRef]#Code    | "internal": "RLA_Bus:StopPoint:BP:RIMIE9:LOC"     |
      | DestinationName                   | Rimiez Saint-George                               |
      | Reference[JourneyPatternRef]#Code | "internal": "RLA_Bus:JourneyPattern::L05P99:LOC"  |
      | Monitored                         | false                                             |
    And one Vehicle has the following attributes:
      | Codes          | "internal": "RLA290"          |
      | Longitude      | 7.276192074052043             |
      | Latitude       | 43.70347861870634             |
      | DriverRef      | "5753"                        |
      | Bearing        | 287.0                         |
      | LinkDistance   | 349.0                         |
      | Percentage     | 70.0                          |
      | ValidUntilTime | 2021-08-02T08:50:27.733+02:00 |

  Scenario: Collect vehicle next Stop in SIRI VehicleMonitoring delivery
    Given a SIRI server waits GetVehicleMonitoring request on "http://localhost:8090" to respond with
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns1:GetVehicleMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>RLA_Bus</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>RLA_Bus:ResponseMessage::23833:LOC</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
          <ns5:Status>true</ns5:Status>
          <ns5:VehicleActivity>
            <ns5:RecordedAtTime>2021-08-02T08:50:27.733+02:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>290</ns5:ItemIdentifier>
            <ns5:ValidUntilTime>2021-08-02T09:50:27.733+02:00</ns5:ValidUntilTime>
            <ns5:VehicleMonitoringRef>290</ns5:VehicleMonitoringRef>
            <ns5:ProgressBetweenStops>
              <ns5:LinkDistance>349.0</ns5:LinkDistance>
              <ns5:Percentage>70.0</ns5:Percentage>
            </ns5:ProgressBetweenStops>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>RLA_Bus:Line::05:LOC</ns5:LineRef>
              <ns5:DirectionRef>Aller</ns5:DirectionRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>RLA_Bus:DataFrame::1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>RLA_Bus:VehicleJourney::2978464:LOC</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>RLA_Bus:JourneyPattern::L05P99:LOC</ns5:JourneyPatternRef>
              <ns5:JourneyPatternName>L05P99</ns5:JourneyPatternName>
              <ns5:PublishedLineName>05</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>RLA_Bus:Operator::RLA:LOC</ns5:OperatorRef>
              <ns5:OriginRef>RLA_Bus:StopPoint:BP:DELOY0:LOC</ns5:OriginRef>
              <ns5:OriginName>Deloye / Dubouchage</ns5:OriginName>
              <ns5:DestinationRef>RLA_Bus:StopPoint:BP:RIMIE9:LOC</ns5:DestinationRef>
              <ns5:DestinationName>Rimiez Saint-George</ns5:DestinationName>
              <ns5:Monitored>false</ns5:Monitored>
              <ns5:VehicleLocation srsName="EPSG:2154">
                <ns5:Coordinates>1044593 6298716</ns5:Coordinates>
              </ns5:VehicleLocation>
              <ns5:Bearing>287.0</ns5:Bearing>
              <ns5:VehicleRef>RLA290</ns5:VehicleRef>
              <ns5:DriverRef>5753</ns5:DriverRef>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>RLA_Bus:StopPoint:BP:PASTO8:LOC</ns5:StopPointRef>
                <ns5:Order>6</ns5:Order>
                <ns5:StopPointName>Carabacel</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Rimiez Saint-George</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2021-08-02T07:38:42.000+02:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>delayed</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2021-08-02T07:38:42.000+02:00</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>delayed</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney><ns5:Extensions/></ns5:VehicleActivity>
        </ns5:VehicleMonitoringDelivery>
      </Answer><AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/></ns1:GetVehicleMonitoringResponse>
  </soap:Body>
</soap:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-vehicle-monitoring-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name  | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
      # 6ba7b814-9dad-11d1-3-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Codes  | "internal": "RLA_Bus:VehicleJourney::2978464:LOC" |
      | LineId | 6ba7b814-9dad-11d1-3-00c04fd430c8                 |
      # 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Codes | "internal": "RLA_Bus:StopPoint:BP:PASTO8:LOC" |
      | Name  | Carabacel                                     |
      # 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes            | "internal": "RLA920-RLA_Bus:StopPoint:BP:PASTO8:LOC" |
      | PassageOrder     | 6                                                    |
      | StopAreaId       | 6ba7b814-9dad-11d1-5-00c04fd430c8                    |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8                    |
      # 6ba7b814-9dad-11d1-6-00c04fd430c8
    When a minute has passed
    And the SIRI server has received a GetVehicleMonitoring request
    Then one Vehicle has the following attributes:
      | Codes            | "internal": "RLA290"              |
      | LineId           | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | StopAreaId       | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-6-00c04fd430c8 |

  Scenario: Update vehicle next Stop in SIRI VehicleMonitoring delivery
    Given a SIRI server waits GetVehicleMonitoring request on "http://localhost:8090" to respond with
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns1:GetVehicleMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>RLA_Bus</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>RLA_Bus:ResponseMessage::23833:LOC</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
          <ns5:Status>true</ns5:Status>
          <ns5:VehicleActivity>
            <ns5:RecordedAtTime>2021-08-02T08:50:27.733+02:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>290</ns5:ItemIdentifier>
            <ns5:ValidUntilTime>2021-08-02T09:50:27.733+02:00</ns5:ValidUntilTime>
            <ns5:VehicleMonitoringRef>290</ns5:VehicleMonitoringRef>
            <ns5:ProgressBetweenStops>
              <ns5:LinkDistance>349.0</ns5:LinkDistance>
              <ns5:Percentage>70.0</ns5:Percentage>
            </ns5:ProgressBetweenStops>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>RLA_Bus:Line::05:LOC</ns5:LineRef>
              <ns5:DirectionRef>Aller</ns5:DirectionRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>RLA_Bus:DataFrame::1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>RLA_Bus:VehicleJourney::2978464:LOC</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>RLA_Bus:JourneyPattern::L05P99:LOC</ns5:JourneyPatternRef>
              <ns5:JourneyPatternName>L05P99</ns5:JourneyPatternName>
              <ns5:PublishedLineName>05</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>RLA_Bus:Operator::RLA:LOC</ns5:OperatorRef>
              <ns5:OriginRef>RLA_Bus:StopPoint:BP:DELOY0:LOC</ns5:OriginRef>
              <ns5:OriginName>Deloye / Dubouchage</ns5:OriginName>
              <ns5:DestinationRef>RLA_Bus:StopPoint:BP:RIMIE9:LOC</ns5:DestinationRef>
              <ns5:DestinationName>Rimiez Saint-George</ns5:DestinationName>
              <ns5:Monitored>false</ns5:Monitored>
              <ns5:VehicleLocation srsName="EPSG:2154">
                <ns5:Coordinates>1044593 6298716</ns5:Coordinates>
              </ns5:VehicleLocation>
              <ns5:Bearing>287.0</ns5:Bearing>
              <ns5:VehicleRef>RLA290</ns5:VehicleRef>
              <ns5:DriverRef>5753</ns5:DriverRef>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>RLA_Bus:StopPoint:BP:CAL05:LOC</ns5:StopPointRef>
                <ns5:Order>7</ns5:Order>
                <ns5:StopPointName>Vieux Port</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Rimiez Saint-George</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2021-08-02T07:38:42.000+02:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>delayed</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2021-08-02T07:38:42.000+02:00</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>delayed</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney><ns5:Extensions/></ns5:VehicleActivity>
        </ns5:VehicleMonitoringDelivery>
      </Answer><AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/></ns1:GetVehicleMonitoringResponse>
  </soap:Body>
</soap:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-vehicle-monitoring-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name  | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
      # 6ba7b814-9dad-11d1-3-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Codes  | "internal": "RLA_Bus:VehicleJourney::2978464:LOC" |
      | LineId | 6ba7b814-9dad-11d1-3-00c04fd430c8                 |
      # 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Codes | "internal": "RLA_Bus:StopPoint:BP:PASTO8:LOC" |
      | Name  | Carabacel                                     |
      # 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Codes | "internal": "RLA_Bus:StopPoint:BP:CAL05:LOC" |
      | Name  | Vieux Port                                   |
      # 6ba7b814-9dad-11d1-6-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes            | "internal": "RLA920-RLA_Bus:StopPoint:BP:PASTO8:LOC" |
      | PassageOrder     | 6                                                    |
      | StopAreaId       | 6ba7b814-9dad-11d1-5-00c04fd430c8                    |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8                    |
      # 6ba7b814-9dad-11d1-7-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes            | "internal": "RLA920-RLA_Bus:StopPoint:BP:CAL05:LOC" |
      | PassageOrder     | 7                                                   |
      | StopAreaId       | 6ba7b814-9dad-11d1-6-00c04fd430c8                   |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8                   |
      # 6ba7b814-9dad-11d1-8-00c04fd430c8
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "RLA290"              |
      | LineId           | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | StopAreaId       | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      # 6ba7b814-9dad-11d1-9-00c04fd430c8
    When a minute has passed
    And the SIRI server has received a GetVehicleMonitoring request
    Then the Vehicle "6ba7b814-9dad-11d1-9-00c04fd430c8" has the following attributes:
      | Codes            | "internal": "RLA290"              |
      | LineId           | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | StopAreaId       | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-8-00c04fd430c8 |

  @ARA-1298 @siri-valid
  Scenario: Performs a SIRI VehicleMonitoring request to a Partner
    Given a SIRI server waits GetVehicleMonitoring request on "http://localhost:8090" to respond with
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns1:GetVehicleMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>RLA_Bus</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>RLA_Bus:ResponseMessage::23833:LOC</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
          <ns5:Status>true</ns5:Status>
          <ns5:VehicleActivity>
            <ns5:RecordedAtTime>2021-08-02T08:50:27.733+02:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>290</ns5:ItemIdentifier>
            <ns5:ValidUntilTime>2021-08-02T09:50:27.733+02:00</ns5:ValidUntilTime>
            <ns5:VehicleMonitoringRef>290</ns5:VehicleMonitoringRef>
            <ns5:ProgressBetweenStops>
              <ns5:LinkDistance>349.0</ns5:LinkDistance>
              <ns5:Percentage>70.0</ns5:Percentage>
            </ns5:ProgressBetweenStops>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>RLA_Bus:Line::05:LOC</ns5:LineRef>
              <ns5:DirectionRef>Aller</ns5:DirectionRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>RLA_Bus:DataFrame::1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>RLA_Bus:VehicleJourney::2978464:LOC</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>RLA_Bus:JourneyPattern::L05P99:LOC</ns5:JourneyPatternRef>
              <ns5:JourneyPatternName>L05P99</ns5:JourneyPatternName>
              <ns5:PublishedLineName>05</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>RLA_Bus:Operator::RLA:LOC</ns5:OperatorRef>
              <ns5:OriginRef>RLA_Bus:StopPoint:BP:DELOY0:LOC</ns5:OriginRef>
              <ns5:OriginName>Deloye / Dubouchage</ns5:OriginName>
              <ns5:DestinationRef>RLA_Bus:StopPoint:BP:RIMIE9:LOC</ns5:DestinationRef>
              <ns5:DestinationName>Rimiez Saint-George</ns5:DestinationName>
              <ns5:Monitored>false</ns5:Monitored>
              <ns5:VehicleLocation srsName="EPSG:2154">
                <ns5:Coordinates>1044593 6298716</ns5:Coordinates>
              </ns5:VehicleLocation>
              <ns5:Bearing>287.0</ns5:Bearing>
              <ns5:VehicleRef>RLA290</ns5:VehicleRef>
              <ns5:DriverRef>5753</ns5:DriverRef>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>RLA_Bus:StopPoint:BP:PASTO8:LOC</ns5:StopPointRef>
                <ns5:Order>6</ns5:Order>
                <ns5:StopPointName>Carabacel</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Rimiez Saint-George</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2021-08-02T07:38:42.000+02:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>delayed</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2021-08-02T07:38:42.000+02:00</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>delayed</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney><ns5:Extensions/></ns5:VehicleActivity>
        </ns5:VehicleMonitoringDelivery>
      </Answer><AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/></ns1:GetVehicleMonitoringResponse>
  </soap:Body>
</soap:Envelope>
      """
    And a Line exists with the following attributes:
      | Codes | "external": "RLA_Bus:Line::05:LOC" |
    And a Partner "test" exists with connectors [siri-check-status-client, siri-vehicle-monitoring-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | remote_code_space | internal              |
    And a minute has passed
    And a minute has passed
    Then the SIRI server should not have received a GetVehicleMonitoring request
    And the Partner "test" is updated with the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_code_space | external              |
      | remote_credential | test                  |
    When a minute has passed
    And a minute has passed
    Then the SIRI server should have received 1 GetVehicleMonitoring request

  @siri-valid @ARA-1298
  Scenario: Handle a SIRI VehicleMonitoring request with Partner remote_code_space changed
    Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential      | test     |
      | remote_code_space     | internal |
      | sort_payload_for_test | true     |
    And a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro                 |
      # Id 6ba7b814-9dad-11d1-2-00c04fd430c8
    And a Line exists with the following attributes:
      | Codes | "external": "Test:Line:A:BUS:LOC" |
      | Name  | Ligne A Bus                       |
      # Id 6ba7b814-9dad-11d1-3-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                                |
      | Codes     | "external": "Test:VehicleJourney:201:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8         |
      | Monitored | true                                      |
      # Id 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 33                                |
      | Codes                             | "internal": "Test:VehicleJourney:202:LOC" |
      | LineId                            | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored                         | true                                      |
      | Reference[DestinationRef]#Code    | "internal": "Test:StopPoint:Destination"  |
      | Reference[JourneyPatternRef]#Code | "internal": "Test:JourneyPattern:1"       |
      | Reference[OriginRef]#Code         | "internal": "Test:StopPoint:Origin"       |
      | OriginName                        | Origin Name                               |
      | DestinationName                   | Destination Name                          |
      | DirectionName                     | Direction Name                            |
      | DirectionType                     | outbound                                  |
      | Attribute[JourneyPatternName]     | Journey Pattern Name                      |
      # Id 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:1:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | DriverRef        | Driver1                           |
      | Bearing          | 120                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | LinkDistance     | 12                                |
      | Percentage       | 42                                |
    And a Vehicle exists with the following attributes:
      | Codes            | "external": "Test:Vehicle:2:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | DriverRef        | Driver2                           |
      | Bearing          | 153                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | LinkDistance     | 34                                |
      | Percentage       | 55                                |
    When I send this SIRI request
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceRequestInfo>
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:RequestorRef>test</siri:RequestorRef>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            </ServiceRequestInfo>
            <Request version="2.0:FR-IDF-2.4">
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
              <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
            </Request>
            <RequestExtension />
          </sw:GetVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then I should receive this SIRI response
    """
    <?xml version='1.0' encoding='UTF-8'?>
    <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
      <S:Body>
        <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:ProducerRef>Ara</siri:ProducerRef>
            <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Answer>
            <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
              <siri:Status>true</siri:Status>
              <siri:VehicleActivity>
                <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                <siri:VehicleMonitoringRef>Test:Vehicle:1:LOC</siri:VehicleMonitoringRef>
                <siri:ProgressBetweenStops>
                  <siri:LinkDistance>12</siri:LinkDistance>
                  <siri:Percentage>42</siri:Percentage>
                </siri:ProgressBetweenStops>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>Test:VehicleJourney:202:LOC</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:JourneyPatternRef>Test:JourneyPattern:1</siri:JourneyPatternRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:OriginRef>RATPDev:StopPoint:Q:488317b5b41cb7ba0a4812c18b312f0e2b986852:LOC</siri:OriginRef>
                  <siri:OriginName>Origin Name</siri:OriginName>
                  <siri:DestinationRef>RATPDev:StopPoint:Q:7bef317e38443efe7d8e8e7f3b7b59881b2e3be0:LOC</siri:DestinationRef>
                  <siri:DestinationName>Destination Name</siri:DestinationName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:VehicleLocation>
                    <siri:Longitude>1.234</siri:Longitude>
                    <siri:Latitude>5.678</siri:Latitude>
                  </siri:VehicleLocation>
                  <siri:Bearing>120</siri:Bearing>
                  <siri:DriverRef>Driver1</siri:DriverRef>
                </siri:MonitoredVehicleJourney>
              </siri:VehicleActivity>
            </siri:VehicleMonitoringDelivery>
          </Answer>
          <AnswerExtension/>
        </sw:GetVehicleMonitoringResponse>
      </S:Body>
    </S:Envelope>
    """
    And the Partner "test" is updated with the following settings:
      | local_credential  | test     |
      | remote_code_space | external |
    And a minute has passed
    When I send this SIRI request
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceRequestInfo>
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:RequestorRef>test</siri:RequestorRef>
              <siri:MessageIdentifier>Test:123456::LOC</siri:MessageIdentifier>
            </ServiceRequestInfo>
            <Request version="2.0:FR-IDF-2.4">
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:MessageIdentifier>Test:123456::LOC</siri:MessageIdentifier>
              <siri:LineRef>Test:Line:A:BUS:LOC</siri:LineRef>
            </Request>
            <RequestExtension />
          </sw:GetVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then I should receive this SIRI response
    """
     <?xml version='1.0' encoding='UTF-8'?>
     <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
       <S:Body>
         <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
           <ServiceDeliveryInfo>
             <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
             <siri:ProducerRef>Ara</siri:ProducerRef>
             <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-9-00c04fd430c8</siri:ResponseMessageIdentifier>
             <siri:RequestMessageRef>Test:123456::LOC</siri:RequestMessageRef>
           </ServiceDeliveryInfo>
           <Answer>
             <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
               <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
               <siri:RequestMessageRef>Test:123456::LOC</siri:RequestMessageRef>
               <siri:Status>true</siri:Status>
               <siri:VehicleActivity>
                 <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                 <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                 <siri:VehicleMonitoringRef>Test:Vehicle:2:LOC</siri:VehicleMonitoringRef>
                 <siri:ProgressBetweenStops>
                   <siri:LinkDistance>34</siri:LinkDistance>
                   <siri:Percentage>55</siri:Percentage>
                 </siri:ProgressBetweenStops>
                 <siri:MonitoredVehicleJourney>
                   <siri:LineRef>Test:Line:A:BUS:LOC</siri:LineRef>
                   <siri:FramedVehicleJourneyRef>
                     <siri:DataFrameRef>2017-01-01</siri:DataFrameRef>
                     <siri:DatedVehicleJourneyRef>Test:VehicleJourney:201:LOC</siri:DatedVehicleJourneyRef>
                   </siri:FramedVehicleJourneyRef>
                   <siri:PublishedLineName>Ligne A Bus</siri:PublishedLineName>
                   <siri:Monitored>true</siri:Monitored>
                   <siri:VehicleLocation>
                     <siri:Longitude>1.234</siri:Longitude>
                     <siri:Latitude>5.678</siri:Latitude>
                   </siri:VehicleLocation>
                   <siri:Bearing>153</siri:Bearing>
                   <siri:DriverRef>Driver2</siri:DriverRef>
                 </siri:MonitoredVehicleJourney>
               </siri:VehicleActivity>
             </siri:VehicleMonitoringDelivery>
           </Answer>
           <AnswerExtension/>
         </sw:GetVehicleMonitoringResponse>
       </S:Body>
     </S:Envelope>
    """

  @ARA-1363 @siri-valid
  Scenario: Handle a SIRI VehicleMonitoring request using the generator setting reference_vehicle_journey_identifier
    # Setting a Partner without default generators
    Given a Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential                                | test                             |
      | remote_code_space                            | internal                         |
      | sort_payload_for_test                           | true                             |
      | generators.reference_vehicle_journey_identifier | ch:1:ServiceJourney:87_TAC:%{id} |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name      | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                                                        |
      | Codes | "_default": "6ba7b814", "external": "Test:VehicleJourney:201:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-2-00c04fd430c8                                 |
      | Monitored | true                                                              |
    And a Vehicle exists with the following attributes:
      | Codes        | "internal": "Test:Vehicle:1:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | DriverRef        | Driver1                           |
      | Bearing          | 120                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | LinkDistance     | 12                                |
      | Percentage       | 42                                |
    When I send this SIRI request
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceRequestInfo>
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:RequestorRef>test</siri:RequestorRef>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            </ServiceRequestInfo>
            <Request version="2.0:FR-IDF-2.4">
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
              <siri:VehicleRef>Test:Vehicle:1:LOC</siri:VehicleRef>
            </Request>
            <RequestExtension />
          </sw:GetVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then I should receive this SIRI response
    """
    <?xml version='1.0' encoding='UTF-8'?>
    <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
      <S:Body>
        <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:ProducerRef>Ara</siri:ProducerRef>
            <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:ResponseMessageIdentifier>
            <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Answer>
            <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
              <siri:Status>true</siri:Status>
              <siri:VehicleActivity>
                <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                <siri:VehicleMonitoringRef>Test:Vehicle:1:LOC</siri:VehicleMonitoringRef>
                <siri:ProgressBetweenStops>
                  <siri:LinkDistance>12</siri:LinkDistance>
                  <siri:Percentage>42</siri:Percentage>
                </siri:ProgressBetweenStops>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>2017-01-01</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>ch:1:ServiceJourney:87_TAC:6ba7b814</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:VehicleLocation>
                    <siri:Longitude>1.234</siri:Longitude>
                    <siri:Latitude>5.678</siri:Latitude>
                  </siri:VehicleLocation>
                  <siri:Bearing>120</siri:Bearing>
                  <siri:DriverRef>Driver1</siri:DriverRef>
                </siri:MonitoredVehicleJourney>
              </siri:VehicleActivity>
             </siri:VehicleMonitoringDelivery>
          </Answer>
          <AnswerExtension/>
        </sw:GetVehicleMonitoringResponse>
      </S:Body>
    </S:Envelope>
    """
    Then an audit event should exist with these attributes:
        | Type              | VehicleMonitoringRequest                |
        | Protocol          | siri                                    |
        | Direction         | received                                |
        | Status            | OK                                      |
        | Partner           | test                                    |
        | Vehicles          | ["Test:Vehicle:1:LOC"]                  |
        | RequestIdentifier | Test:1234::LOC                          |
        | Lines             | ["Test:Line:3:LOC"]                     |
        | VehicleJourneys   | ["ch:1:ServiceJourney:87_TAC:6ba7b814"] |

  @ARA-1363 @siri-valid
  Scenario: Handle a SIRI VehicleMonitoring request using the default generator should send DatedVehicleJourneyRef according to default setting
    # Setting a "SIRI Partner" with default generators
    Given a SIRI Partner "test" exists with connectors [siri-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential      | test     |
      | remote_code_space  | internal |
      | sort_payload_for_test | true     |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name      | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                                                        |
      | Codes | "_default": "6ba7b814", "external": "Test:VehicleJourney:201:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-2-00c04fd430c8                                 |
      | Monitored | true                                                              |
    And a Vehicle exists with the following attributes:
      | Codes        | "internal": "Test:Vehicle:1:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | DriverRef        | Driver1                           |
      | Bearing          | 120                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | LinkDistance     | 12                                |
      | Percentage       | 42                                |
    When I send this SIRI request
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <sw:GetVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceRequestInfo>
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:RequestorRef>test</siri:RequestorRef>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
            </ServiceRequestInfo>
            <Request version="2.0:FR-IDF-2.4">
              <siri:RequestTimestamp>2006-01-02T15:04:05.000Z</siri:RequestTimestamp>
              <siri:MessageIdentifier>Test:1234::LOC</siri:MessageIdentifier>
              <siri:VehicleRef>Test:Vehicle:1:LOC</siri:VehicleRef>
            </Request>
            <RequestExtension />
          </sw:GetVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then I should receive this SIRI response
    """
    <?xml version='1.0' encoding='UTF-8'?>
    <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
      <S:Body>
        <sw:GetVehicleMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:ProducerRef>Ara</siri:ProducerRef>
            <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Answer>
            <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:RequestMessageRef>Test:1234::LOC</siri:RequestMessageRef>
              <siri:Status>true</siri:Status>
              <siri:VehicleActivity>
                <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                <siri:VehicleMonitoringRef>Test:Vehicle:1:LOC</siri:VehicleMonitoringRef>
                <siri:ProgressBetweenStops>
                  <siri:LinkDistance>12</siri:LinkDistance>
                  <siri:Percentage>42</siri:Percentage>
                </siri:ProgressBetweenStops>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>RATPDev:VehicleJourney::6ba7b814:LOC</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:VehicleLocation>
                    <siri:Longitude>1.234</siri:Longitude>
                    <siri:Latitude>5.678</siri:Latitude>
                  </siri:VehicleLocation>
                  <siri:Bearing>120</siri:Bearing>
                  <siri:DriverRef>Driver1</siri:DriverRef>
                </siri:MonitoredVehicleJourney>
              </siri:VehicleActivity>
             </siri:VehicleMonitoringDelivery>
          </Answer>
          <AnswerExtension/>
        </sw:GetVehicleMonitoringResponse>
      </S:Body>
    </S:Envelope>
    """
    Then an audit event should exist with these attributes:
        | Type              | VehicleMonitoringRequest                 |
        | Protocol          | siri                                     |
        | Direction         | received                                 |
        | Status            | OK                                       |
        | Partner           | test                                     |
        | Vehicles          | ["Test:Vehicle:1:LOC"]                   |
        | RequestIdentifier | Test:1234::LOC                           |
        | Lines             | ["Test:Line:3:LOC"]                      |
        | VehicleJourneys   | ["RATPDev:VehicleJourney::6ba7b814:LOC"] |

  @siri-valid @ARA-1591
  Scenario: VehicleMonitoringDelivery with Status false for one Delivery is logged as an Error status in BigQuery
    Given a SIRI server waits GetVehicleMonitoring request on "http://localhost:8090" to respond with
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns1:GetVehicleMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>RLA_Bus</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>RLA_Bus:ResponseMessage::23833:LOC</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
          <ns5:Status>false</ns5:Status>
          <ns5:ErrorCondition>
            <ns5:AllowedResourceUsageExceededError>
              <ns5:ErrorText>too many requets</ns5:ErrorText>
            </ns5:AllowedResourceUsageExceededError>
          </ns5:ErrorCondition>
        </ns5:VehicleMonitoringDelivery>
      </Answer><AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/></ns1:GetVehicleMonitoringResponse>
  </soap:Body>
</soap:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-vehicle-monitoring-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name  | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
    When a minute has passed
    And the SIRI server has received a GetVehicleMonitoring request
    And an audit event should exist with these attributes:
      | Protocol  | siri                     |
      | Direction | sent                     |
      | Status    | Error                    |
      | Type      | VehicleMonitoringRequest |
