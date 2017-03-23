Feature: Support SIRI GetSIRI

  Background:
    Given a Referential "test" is created

  Scenario: 2462 - Handle a SIRI GetSIRIService request with several StopMonitorings
    Given a Partner "test" exists with connectors [siri-service-request-broadcaster,siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | RATPDEV:Concerto |
      | remote_objectid_kind | internal         |
    And a StopArea exists with the following attributes:
      | Name      | Test 1                |
      | ObjectIDs | "internal": "boaarle" |
    And a Line exists with the following attributes:
      | Name      | Ligne 415                       |
      | ObjectIDs | "internal": "CdF:Line::415:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
    And a StopVisit exists with the following attributes:
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | ArrivalStatus                 | onTime                            |
      | DepartureStatus               | onTime                            |
      | ObjectIDs                     | "internal": "SIRI:34852540"       |
      | PassageOrder                  | 44                                |
      | RecordedAt                    | 2017-01-12T10:52:46.042+01:00     |
      | Schedule[aimed]#Arrival       | 2017-01-12T11:42:54.000+01:00     |
      | Schedule[aimed]#Departure     | 2017-01-12T11:42:54.000+01:00     |
      | Schedule[expected]#Arrival    | 2017-01-12T11:42:54.000+01:00     |
      | Schedule[expected]#Departure  | 2017-01-12T11:42:54.000+01:00     |
      | VehicleAtStop                 | false                             |
      | Attribute[DestinationDisplay] | Méliès - Croix Bonnet             |
      | Reference["OperatorRef"]      | CdF:Company::410:LOC.             |
    And a StopArea exists with the following attributes:
      | Name      | Test 2                |
      | ObjectIDs | "internal": "cladebr" |
    And a Line exists with the following attributes:
      | Name      | Ligne 475                       |
      | ObjectIDs | "internal": "CdF:Line::475:LOC" |
    And a VehicleJourney exists with the following attributes:
      | ObjectIDs                             | "internal": "5CAR621689196575"                |
      | LineId                                | 6ba7b814-9dad-11d1-7-00c04fd430c8             |
      | Attribute[DestinationName]            | PARIS - Porte d'Orléans                       |
      | Attribute[DirectionName]              | Aller                                         |
      | Attribute[DirectionRef]               | cladebr                                       |
      | Attribute[Monitored]                  | true                                          |
      | Reference[JourneyPatternRef]#ObjectID | "internal": "CdF:JourneyPattern::L475P53:LOC" |
      | Reference[DestinationRef]#ObjectID    | "internal": "parorle"                         |
    And a StopVisit exists with the following attributes:
      | StopAreaId                    | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
      | ArrivalStatus                 | onTime                            |
      | DepartureStatus               | onTime                            |
      | ObjectIDs                     | "internal": "SIRI:34863800"       |
      | PassageOrder                  | 11                                |
      | RecordedAt                    | 2017-01-12T10:52:46.050+01:00     |
      | Schedule[aimed]#Arrival       | 2017-01-12T11:41:00.000+01:00     |
      | Schedule[aimed]#Departure     | 2017-01-12T11:41:00.000+01:00     |
      | Schedule[expected]#Arrival    | 2017-01-12T11:41:00.000+01:00     |
      | Schedule[expected]#Departure  | 2017-01-12T11:41:00.000+01:00     |
      | VehicleAtStop                 | false                             |
      | Attribute[DestinationDisplay] | PARIS - Porte d'Orléans           |
      | Reference["OperatorRef"]      | CdF:Company::410:LOC.             |
    When I send this SIRI request
      """
      <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:m0="http://www.siri.org.uk/siri" xmlns:m1="http://www.ifopt.org.uk/acsb">
          <SOAP-ENV:Body>
              <m:GetSiriService xmlns:m="http://wsdl.siri.org.uk">
                  <Request>
                      <m0:ServiceRequestContext>
                      </m0:ServiceRequestContext>
                      <m0:RequestTimestamp>2001-12-17T09:30:47Z</m0:RequestTimestamp>
                      <m0:RequestorRef>RATPDEV:Concerto</m0:RequestorRef>
                      <m0:MessageIdentifier>GetSIRIStopMonitoring:Test:0</m0:MessageIdentifier>
                      <m0:DelegatorRef/>
                      <m0:StopMonitoringRequest version="2.0:FR-IDF-2.4">
                          <m0:RequestTimestamp>2017-01-10T16:30:47Z</m0:RequestTimestamp>
                          <m0:MessageIdentifier>GetSIRIStopMonitoring:Test:0</m0:MessageIdentifier>
                          <m0:MonitoringRef>boaarle</m0:MonitoringRef>
                          <m0:StopVisitTypes>all</m0:StopVisitTypes>
                      </m0:StopMonitoringRequest>
                      <m0:StopMonitoringRequest version="2.0:FR-IDF-2.4">
                          <m0:RequestTimestamp>2017-01-10T16:30:47Z</m0:RequestTimestamp>
                          <m0:MessageIdentifier>GetSIRIStopMonitoring:Test:0</m0:MessageIdentifier>
                          <m0:MonitoringRef>cladebr</m0:MonitoringRef>
                          <m0:StopVisitTypes>all</m0:StopVisitTypes>
                      </m0:StopMonitoringRequest>
                  </Request>
              </m:GetSiriService>
          </SOAP-ENV:Body>
      </SOAP-ENV:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns1:GetSiriServiceResponse xmlns:ns1="http://wsdl.siri.org.uk">
          <Answer xmlns:ns3="http://www.siri.org.uk/siri"
                  xmlns:ns4="http://www.ifopt.org.uk/acsb"
                  xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                  xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                  xmlns:ns7="http://scma/siri"
                  xmlns:ns8="http://wsdl.siri.org.uk"
                  xmlns:ns9="http://wsdl.siri.org.uk/siri">
              <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
              <ns3:ProducerRef>Edwig</ns3:ProducerRef>
              <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>GetSIRIStopMonitoring:Test:0</ns3:RequestMessageRef>
              <ns3:Status>true</ns3:Status>
              <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
                <ns3:RequestMessageRef>GetSIRIStopMonitoring:Test:0</ns3:RequestMessageRef>
                <ns3:Status>true</ns3:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2017-01-12T10:52:46.042+01:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>SIRI:34852540</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>boaarle</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>CdF:Line::415:LOC</ns3:LineRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:PublishedLineName>Ligne 415</ns3:PublishedLineName>
                    <ns3:VehicleJourneyName>Passage 32</ns3:VehicleJourneyName>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>boaarle</ns3:StopPointRef>
                      <ns3:Order>44</ns3:Order>
                      <ns3:StopPointName>Test 1</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:DestinationDisplay>Méliès - Croix Bonnet</ns3:DestinationDisplay>
                      <ns3:AimedArrivalTime>2017-01-12T11:42:54.000+01:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-12T11:42:54.000+01:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>onTime</ns3:ArrivalStatus>
                      <ns3:AimedDepartureTime>2017-01-12T11:42:54.000+01:00</ns3:AimedDepartureTime>
                      <ns3:ExpectedDepartureTime>2017-01-12T11:42:54.000+01:00</ns3:ExpectedDepartureTime>
                      <ns3:DepartureStatus>onTime</ns3:DepartureStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns3:StopMonitoringDelivery>
              <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
                <ns3:RequestMessageRef>GetSIRIStopMonitoring:Test:0</ns3:RequestMessageRef>
                <ns3:Status>true</ns3:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2017-01-12T10:52:46.050+01:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>SIRI:34863800</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>cladebr</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>CdF:Line::475:LOC</ns3:LineRef>
                    <ns3:DirectionRef>cladebr</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>5CAR621689196575</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>RATPDev:JourneyPattern::3eb50093508f11b474950daa6b0b8632660a32a6:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 475</ns3:PublishedLineName>
                    <ns3:DirectionName>Aller</ns3:DirectionName>
                    <ns3:DestinationRef>RATPDev:StopPoint:Q:9dd925e2cd515383ad6d975e761cea71ea1a79e7:LOC</ns3:DestinationRef>
                    <ns3:DestinationName>PARIS - Porte d'Orléans</ns3:DestinationName>
                    <ns3:Monitored>true</ns3:Monitored>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>cladebr</ns3:StopPointRef>
                      <ns3:Order>11</ns3:Order>
                      <ns3:StopPointName>Test 2</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:DestinationDisplay>PARIS - Porte d'Orléans</ns3:DestinationDisplay>
                      <ns3:AimedArrivalTime>2017-01-12T11:41:00.000+01:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-12T11:41:00.000+01:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>onTime</ns3:ArrivalStatus>
                      <ns3:AimedDepartureTime>2017-01-12T11:41:00.000+01:00</ns3:AimedDepartureTime>
                      <ns3:ExpectedDepartureTime>2017-01-12T11:41:00.000+01:00</ns3:ExpectedDepartureTime>
                      <ns3:DepartureStatus>onTime</ns3:DepartureStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns3:StopMonitoringDelivery>
            </Answer>
          </ns1:GetSiriServiceResponse>
        </S:Body>
      </S:Envelope>
      """
