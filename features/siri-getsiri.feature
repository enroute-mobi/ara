Feature: Support SIRI GetSIRI

  Background:
    Given a Referential "test" is created

  @wip
  Scenario: 2462 - Handle a SIRI GetSIRIService request with several StopMonitorings
    Given a Partner "test" exists with connectors [siri-request-broadcaster siri-stop-monitoring-request-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | Sqybus                |
      | remote_objectid_kind | internal              |
    And a StopArea exists with the following attributes:
      | Name      | Test 1                                                                    |
      | ObjectIDs | "internal": "boaarle", "external": "RATPDev:StopPoint:Q:eeft52df543d:LOC" |
    And a Line exists with the following attributes:
      | Name      | Ligne 415                       |
      | ObjectIDs | "internal": "CdF:Line::415:LOC" |
    And a VehicleJourney exists with the following attributes:
      | ObjectIDs                                 | "internal": "1STD721689197098"                 |
      | Attribute[DestinationName]                | Méliès - Croix Bonnet                          |
      | Attribute[DirectionName]                  | Aller                                          |
      | Attribute[DirectionRef]                   | boabonn                                        |
      | Attribute[Monitored]                      | false                                          |
      | Reference[JourneyPatternRef]#ObjectID     | "internal": "CdF:JourneyPattern::L415P289:LOC" |
      | Reference[DestinationRef]#ObjectID        | "internal": "boabonn"                          |
    And a StopVisit exists with the following attributes:
      | ArrivalStatus                        | onTime                        |
      | DepartureStatus                      | onTime                        |
      | ObjectIDs                            | "internal": "SIRI:34852540"   |
      | PassageOrder                         | 44                            |
      | RecordedAt                           | 2017-01-12T10:52:46.042+01:00 |
      | Schedule[aimed]#Arrival              | 2017-01-12T11:42:54.000+01:00 |
      | Schedule[aimed]#Departure            | 2017-01-12T11:42:54.000+01:00 |
      | Schedule[expected]#Arrival           | 2017-01-12T11:42:54.000+01:00 |
      | Schedule[expected]#Departure         | 2017-01-12T11:42:54.000+01:00 |
      | VehicleAtStop                        | false                         |
      | Attribute[DestinationDisplay]        | Méliès - Croix Bonnet         |
      | Reference["OperatorRef"]             | CdF:Company::410:LOC.         |
    And a StopArea exists with the following attributes:
      | Name      | Test 2                                                                     |
      | ObjectIDs | "internal": "cladebr", "external": "RATPDev:StopPoint:Q:875fdetgyh765:LOC" |
    And a Line exists with the following attributes:
      | Name      | Ligne 475                      |
      | ObjectIDs | "internal": "CdF:Line::475:LOC"|
    And a VehicleJourney exists with the following attributes:
      | ObjectIDs                                 | "internal": "5CAR621689196575"                |
      | Attribute[DestinationName]                | PARIS - Porte d'Orléans                       |
      | Attribute[DirectionName]                  | Aller                                         |
      | Attribute[DirectionRef]                   | cladebr                                       |
      | Attribute[Monitored]                      | true                                          |
      | Reference[JourneyPatternRef]#ObjectID     | "internal": "CdF:JourneyPattern::L475P53:LOC" |
      | Reference[DestinationRef]#ObjectID        | "internal": "parorle"                         |
    And a StopVisit exists with the following attributes:
      | ArrivalStatus                        | onTime                        |
      | DepartureStatus                      | onTime                        |
      | ObjectIDs                            | "internal": "SIRI:34863800"   |
      | PassageOrder                         | 11                            |
      | RecordedAt                           | 2017-01-12T10:52:46.050+01:00 |
      | Schedule[aimed]#Arrival              | 2017-01-12T11:41:00.000+01:00 |
      | Schedule[aimed]#Departure            | 2017-01-12T11:41:00.000+01:00 |
      | Schedule[expected]#Arrival           | 2017-01-12T11:41:00.000+01:00 |
      | Schedule[expected]#Departure         | 2017-01-12T11:41:00.000+01:00 |
      | VehicleAtStop                        | false                         |
      | Attribute[DestinationDisplay]        | PARIS - Porte d'Orléans       |
      | Reference["OperatorRef"]             | CdF:Company::410:LOC.         |
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
      <?xml version="1.0" encoding="UTF-8"?>
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" />
        <soap:Body>
          <ns1:GetSiriServiceResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-12T10:52:46.055+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>Sqybus</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>Sqybus:ResponseMessage::1:LOC</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>GetSIRIStopMonitoring:Test:0</ns5:RequestMessageRef>
              <ns5:Status>true</ns5:Status>
              <ns5:ErrorCondition />
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-12T10:52:46.054+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>GetSIRIStopMonitoring:Test:0</ns5:RequestMessageRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-12T10:52:46.042+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:34852540</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721689197098</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>Ligne 415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-12T11:42:54.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-12T11:42:54.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-12T11:42:54.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-12T11:42:54.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-12T10:52:46.054+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>GetSIRIStopMonitoring:Test:0</ns5:RequestMessageRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-12T10:52:46.050+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:34863800</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>cladebr</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::475:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>5CAR621689196575</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L475P53:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>Ligne 475</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>parorle</ns5:DestinationRef>
                    <ns5:DestinationName>PARIS - Porte d'Orléans</ns5:DestinationName>
                    <ns5:Monitored>false</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>cladebr</ns5:StopPointRef>
                      <ns5:Order>11</ns5:Order>
                      <ns5:StopPointName>Charles Debry</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>PARIS - Porte d'Orléans</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-12T11:41:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-12T11:41:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-12T11:41:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-12T11:41:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
          </ns1:GetSiriServiceResponse>
        </soap:Body>
      </soap:Envelope>
      """
