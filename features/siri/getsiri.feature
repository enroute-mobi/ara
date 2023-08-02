Feature: Support SIRI GetSIRI

  Background:
    Given a Referential "test" is created

  Scenario: 2462 - Handle a SIRI GetSIRIService request with several StopMonitorings
    Given a SIRI Partner "test" exists with connectors [siri-service-request-broadcaster,siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | RATPDEV:Concerto |
      | remote_objectid_kind | internal         |
    And a StopArea exists with the following attributes:
      | Name      | Test 1                |
      | ObjectIDs | "internal": "boaarle" |
      | Monitored | true                  |
    And a Line exists with the following attributes:
      | Name      | Ligne 415                       |
      | ObjectIDs | "internal": "CdF:Line::415:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
      | Monitored | true                                    |
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
      | Monitored | true                  |
    And a Line exists with the following attributes:
      | Name      | Ligne 475                       |
      | ObjectIDs | "internal": "CdF:Line::475:LOC" |
    And a VehicleJourney exists with the following attributes:
      | ObjectIDs                             | "internal": "5CAR621689196575"                |
      | LineId                                | 6ba7b814-9dad-11d1-7-00c04fd430c8             |
      | Monitored                             | true                                          |
      | DestinationName                       | PARIS - Porte d'Orléans                       |
      | Attribute[DirectionName]              | Aller                                         |
      | DirectionType                         | cladebr                                       |
      | Reference[JourneyPatternRef]#ObjectId | "internal": "CdF:JourneyPattern::L475P53:LOC" |
      | Reference[DestinationRef]#ObjectId    | "internal": "parorle"                         |
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
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetSiriServiceResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <Answer>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>GetSIRIStopMonitoring:Test:0</siri:RequestMessageRef>
              <siri:Status>true</siri:Status>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>GetSIRIStopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:MonitoringRef>boaarle</siri:MonitoringRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>2017-01-12T10:52:46.042+01:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>SIRI:34852540</siri:ItemIdentifier>
                  <siri:MonitoringRef>boaarle</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>CdF:Line::415:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 415</siri:PublishedLineName>
                    <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>boaarle</siri:StopPointRef>
                      <siri:Order>44</siri:Order>
                      <siri:StopPointName>Test 1</siri:StopPointName>
                      <siri:VehicleAtStop>false</siri:VehicleAtStop>
                      <siri:DestinationDisplay>Méliès - Croix Bonnet</siri:DestinationDisplay>
                      <siri:AimedArrivalTime>2017-01-12T11:42:54.000+01:00</siri:AimedArrivalTime>
                      <siri:ExpectedArrivalTime>2017-01-12T11:42:54.000+01:00</siri:ExpectedArrivalTime>
                      <siri:ArrivalStatus>onTime</siri:ArrivalStatus>
                      <siri:AimedDepartureTime>2017-01-12T11:42:54.000+01:00</siri:AimedDepartureTime>
                      <siri:ExpectedDepartureTime>2017-01-12T11:42:54.000+01:00</siri:ExpectedDepartureTime>
                      <siri:DepartureStatus>onTime</siri:DepartureStatus>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
              </siri:StopMonitoringDelivery>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>GetSIRIStopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:MonitoringRef>cladebr</siri:MonitoringRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>2017-01-12T10:52:46.050+01:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>SIRI:34863800</siri:ItemIdentifier>
                  <siri:MonitoringRef>cladebr</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>CdF:Line::475:LOC</siri:LineRef>
                    <siri:DirectionRef>cladebr</siri:DirectionRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>5CAR621689196575</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:JourneyPatternRef>CdF:JourneyPattern::L475P53:LOC</siri:JourneyPatternRef>
                    <siri:PublishedLineName>Ligne 475</siri:PublishedLineName>
                    <siri:DirectionName>Aller</siri:DirectionName>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:9dd925e2cd515383ad6d975e761cea71ea1a79e7:LOC</siri:DestinationRef>
                    <siri:DestinationName>PARIS - Porte d'Orléans</siri:DestinationName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>cladebr</siri:StopPointRef>
                      <siri:Order>11</siri:Order>
                      <siri:StopPointName>Test 2</siri:StopPointName>
                      <siri:VehicleAtStop>false</siri:VehicleAtStop>
                      <siri:DestinationDisplay>PARIS - Porte d'Orléans</siri:DestinationDisplay>
                      <siri:AimedArrivalTime>2017-01-12T11:41:00.000+01:00</siri:AimedArrivalTime>
                      <siri:ExpectedArrivalTime>2017-01-12T11:41:00.000+01:00</siri:ExpectedArrivalTime>
                      <siri:ArrivalStatus>onTime</siri:ArrivalStatus>
                      <siri:AimedDepartureTime>2017-01-12T11:41:00.000+01:00</siri:AimedDepartureTime>
                      <siri:ExpectedDepartureTime>2017-01-12T11:41:00.000+01:00</siri:ExpectedDepartureTime>
                      <siri:DepartureStatus>onTime</siri:DepartureStatus>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
              </siri:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSiriServiceResponse>
        </S:Body>
      </S:Envelope>
      """
    Then an audit event should exist with these attributes:
      | Type            | SiriServiceRequest                                |
      | Protocol        | siri                                              |
      | Direction       | received                                          |
      | Status          | OK                                                |
      | Partner         | test                                              |
      | StopAreas       | ["boaarle", "cladebr"]                            |
      | Lines           | ["CdF:Line::415:LOC", "CdF:Line::475:LOC"]        |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201", "5CAR621689196575"] |

  Scenario: 2481 - Handle a GetSIRIService request with a StopMonitoring request on a unknown StopArea
    Given a SIRI Partner "test" exists with connectors [siri-service-request-broadcaster,siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | RATPDEV:Concerto |
      | remote_objectid_kind | internal         |
    And a StopArea exists with the following attributes:
      | Name      | Test 1                |
      | ObjectIDs | "internal": "boaarle" |
      | Monitored | true                  |
    And a Line exists with the following attributes:
      | Name      | Ligne 415                       |
      | ObjectIDs | "internal": "CdF:Line::415:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
      | Monitored | true                                    |
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
      | DirectionType                         | cladebr                                       |
      | Attribute[Monitored]                  | true                                          |
      | Reference[JourneyPatternRef]#ObjectId | "internal": "CdF:JourneyPattern::L475P53:LOC" |
      | Reference[DestinationRef]#ObjectId    | "internal": "parorle"                         |
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
                          <m0:MessageIdentifier>GetSIRIStopMonitoring:Test:1</m0:MessageIdentifier>
                          <m0:MonitoringRef>unknown</m0:MonitoringRef>
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
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetSiriServiceResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <Answer>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>GetSIRIStopMonitoring:Test:0</siri:RequestMessageRef>
              <siri:Status>false</siri:Status>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>GetSIRIStopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:MonitoringRef>boaarle</siri:MonitoringRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>2017-01-12T10:52:46.042+01:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>SIRI:34852540</siri:ItemIdentifier>
                  <siri:MonitoringRef>boaarle</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>CdF:Line::415:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 415</siri:PublishedLineName>
                    <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>boaarle</siri:StopPointRef>
                      <siri:Order>44</siri:Order>
                      <siri:StopPointName>Test 1</siri:StopPointName>
                      <siri:VehicleAtStop>false</siri:VehicleAtStop>
                      <siri:DestinationDisplay>Méliès - Croix Bonnet</siri:DestinationDisplay>
                      <siri:AimedArrivalTime>2017-01-12T11:42:54.000+01:00</siri:AimedArrivalTime>
                      <siri:ExpectedArrivalTime>2017-01-12T11:42:54.000+01:00</siri:ExpectedArrivalTime>
                      <siri:ArrivalStatus>onTime</siri:ArrivalStatus>
                      <siri:AimedDepartureTime>2017-01-12T11:42:54.000+01:00</siri:AimedDepartureTime>
                      <siri:ExpectedDepartureTime>2017-01-12T11:42:54.000+01:00</siri:ExpectedDepartureTime>
                      <siri:DepartureStatus>onTime</siri:DepartureStatus>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
              </siri:StopMonitoringDelivery>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>GetSIRIStopMonitoring:Test:1</siri:RequestMessageRef>
                <siri:MonitoringRef>unknown</siri:MonitoringRef>
                <siri:Status>false</siri:Status>
                <siri:ErrorCondition>
                  <siri:InvalidDataReferencesError>
                    <siri:ErrorText>StopArea not found: 'unknown'</siri:ErrorText>
                  </siri:InvalidDataReferencesError>
                </siri:ErrorCondition>
              </siri:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSiriServiceResponse>
        </S:Body>
      </S:Envelope>
      """
