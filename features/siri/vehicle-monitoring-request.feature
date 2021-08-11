Feature: Support SIRI VehicleMonitoring by request

  Background:
      Given a Referential "test" is created

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
      | remote_objectid_kind  | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Test 1                             |
      | ObjectIDs | "internal": "RLA_Bus:Line::05:LOC" |
    When a minute has passed
    And the SIRI server has received a GetVehicleMonitoring request
    Then one StopArea has the following attributes:
      | ObjectIDs | "internal": "RLA_Bus:StopPoint:BP:PASTO8:LOC" |
      | Name      | Carabacel                                     |
    And one Line has the following attributes:
      | ObjectIDs | "internal": "RLA_Bus:Line::05:LOC" |
      | Name      | Test 1                             |
    And one VehicleJourney has the following attributes:
      | ObjectIDs                             | "internal": "RLA_Bus:VehicleJourney::2978464:LOC" |
      | Reference[OriginRef]#ObjectId         | "internal": "RLA_Bus:StopPoint:BP:DELOY0:LOC"     |
      | OriginName                            | Deloye / Dubouchage                               |
      | Reference[DestinationRef]#ObjectId    | "internal": "RLA_Bus:StopPoint:BP:RIMIE9:LOC"     |
      | DestinationName                       | Rimiez Saint-George                               |
      | Reference[JourneyPatternRef]#ObjectId | "internal": "RLA_Bus:JourneyPattern::L05P99:LOC"  |
      | Monitored                             | false                                             |
    And one Vehicle has the following attributes:
      | ObjectIDs      | "internal": "RLA290"          |
      | SRSName        | EPSG:2154                     |
      | Coordinates    | 1044593 6298716               |
      | DriverRef      | "5753"                        |
      | Bearing        | 287.0                         |
      | LinkDistance   | "349.0"                       |
      | Percentage     | "70.0"                        |
      | ValidUntilTime | 2021-08-02T08:50:27.733+02:00 |
