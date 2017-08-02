Feature: Support SIRI EstimatedTimeTable

  Background:
      Given a Referential "test" is created

@wip
  Scenario: 3950 - Handle a SIRI EstimatedTimeTable request
    Given a Partner "test" exists with connectors [siri-estimated-timetable-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:25:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:26:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:27:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                          | Passage 32                              |
      | ObjectIDs                     | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                        | 6ba7b814-9dad-11d1-6-00c04fd430c8       |
      | Attribute[DirectionRef]       | Aller                                   |
      | Attribute[OriginName]         | Le début                                |
      | Attribute[DestinationName]    | La fin.                                 |
    And a StopVisit exists with the following attributes:
    #retard d'une minute
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T13:00:00.000+02:00                                        |
      | Schedule[expected]#Arrival      | 2017-01-01T13:01:00.000+02:00                                        |
      | ArrivalStatus                   | Delayed                                                              |
    And a StopVisit exists with the following attributes:
    #retard d'une minute
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-3" |
      | PassageOrder                    | 5                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T13:05:00.000+02:00                                        |
      | Schedule[expected]#Arrival      | 2017-01-01T13:06:00.000+02:00                                        |
      | ArrivalStatus                   | Delayed                                                              |
    And a StopVisit exists with the following attributes:
    #retard d'une minute
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:26:LOC-3" |
      | PassageOrder                    | 6                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T13:10:00.000+02:00                                        |
      | Schedule[expected]#Arrival      | 2017-01-01T13:11:00.000+02:00                                        |
      | ArrivalStatus                   | Delayed                                                              |
    And a StopVisit exists with the following attributes:
    #à l'heure
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-3" |
      | PassageOrder                    | 7                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T13:16:00.000+02:00                                        |
      | Schedule[expected]#Arrival      | 2017-01-01T13:16:00.000+02:00                                        |
      | ArrivalStatus                   | onTime                                                               |
    And I see edwig vehicle_journeys
    And I see edwig stop_visits
    And I see edwig lines
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetEstimatedTimetable xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:ns3="http://www.ifopt.org.uk/acsb"
                                 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                                 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                                 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <ServiceRequestInfo>
              <ns2:RequestTimestamp>2017-01-01T12:00:00.000+02:00</ns2:RequestTimestamp>
              <ns2:RequestorRef>test</ns2:RequestorRef>
              <ns2:MessageIdentifier>EstimatedTimetable:Test:0</ns2:MessageIdentifier>
            </ServiceRequestInfo>

            <Request version="2.0:FR-IDF-2.4">
              <ns2:RequestTimestamp>2017-01-01T12:00:00.000+02:00</ns2:RequestTimestamp>
              <ns2:MessageIdentifier>EstimatedTimetable:Test:0</ns2:MessageIdentifier>
              <ns2:Lines>
                <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
              </ns2:Lines>
            </Request>
            <RequestExtension />
          </ns7:GetEstimatedTimetable>
        </S:Body>
      </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns8:GetEstimatedTimetableResponse xmlns:ns3="http://www.siri.org.uk/siri"
                                         xmlns:ns4="http://www.ifopt.org.uk/acsb"
                                         xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                                         xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                                         xmlns:ns7="http://scma/siri"
                                         xmlns:ns8="http://wsdl.siri.org.uk"
                                         xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>2017-01-01T12:00:00.000+02:00</ns3:ResponseTimestamp>
              <ns3:ProducerRef>Edwig</ns3:ProducerRef>
              <ns3:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>EstimatedTimetable:Test:0</ns3:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>2017-01-01T12:00:00.000+02:00</ns3:ResponseTimestamp>
                <ns3:RequestMessageRef>EstimatedTimetable:Test:0</ns3:RequestMessageRef>
                <ns3:Status>true</ns3:Status>
                <ns3:EstimatedJourneyVersionFrame>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:EstimatedVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:DirectionRef>Aller</ns3:DirectionRef>
                    <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:20:LOC-3</ns3:OriginRef>
                    <ns3:OriginName>Le début</ns3:OriginName>
                    <ns3:DestinationRef>NINOXE:StopPoint:SP:27:LOC</ns3:DestinationRef>
                    <ns3:DestinationName>La fin.</ns3:DestinationName>
                    <ns3:EstimatedCalls>
                      <ns3:EstimatedCall>
                        <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                        <ns3:Order>4</ns3:Order>
                        <ns3:AimedArrivalTime>2017-01-01T13:00:00.000+02:00</ns3:AimedArrivalTime>
                        <ns3:ActualArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ActualArrivalTime>
                        <ns3:ArrivalStatus>Delayed</ns3:ArrivalStatus>
                      </ns3:EstimatedCall>
                      <ns3:EstimatedCall>
                        <ns3:StopPointRef>NINOXE:StopPoint:SP:25:LOC</ns3:StopPointRef>
                        <ns3:Order>5</ns3:Order>
                        <ns3:AimedArrivalTime>2017-01-01T13:00:05.000+02:00</ns3:AimedArrivalTime>
                        <ns3:ActualArrivalTime>2017-01-01T13:01:06.000+02:00</ns3:ActualArrivalTime>
                        <ns3:ArrivalStatus>Delayed</ns3:ArrivalStatus>
                      </ns3:EstimatedCall>
                      <ns3:EstimatedCall>
                        <ns3:StopPointRef>NINOXE:StopPoint:SP:26:LOC</ns3:StopPointRef>
                        <ns3:Order>6</ns3:Order>
                        <ns3:AimedArrivalTime>2017-01-01T13:10:00.000+02:00</ns3:AimedArrivalTime>
                        <ns3:ActualArrivalTime>2017-01-01T13:11:00.000+02:00</ns3:ActualArrivalTime>
                        <ns3:ArrivalStatus>Delayed</ns3:ArrivalStatus>
                      </ns3:EstimatedCall>
                      <ns3:EstimatedCall>
                        <ns3:StopPointRef>NINOXE:StopPoint:SP:27:LOC</ns3:StopPointRef>
                        <ns3:Order>7</ns3:Order>
                        <ns3:AimedArrivalTime>2017-01-01T13:16:00.000+02:00</ns3:AimedArrivalTime>
                        <ns3:ActualArrivalTime>2017-01-01T13:16:00.000+02:00</ns3:ActualArrivalTime>
                        <ns3:ArrivalStatus>onTime</ns3:ArrivalStatus>
                      </ns3:EstimatedCall>
                    </ns3:EstimatedCalls>
                  </ns3:EstimatedVehicleJourney>
                </ns3:EstimatedJourneyVersionFrame>
              </ns3:EstimatedTimetableDelivery>
            </Answer>
          </ns8:GetEstimatedTimetableResponse>
        </S:Body>
      </S:Envelope>
      """
