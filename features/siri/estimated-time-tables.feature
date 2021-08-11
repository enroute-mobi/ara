Feature: Support SIRI EstimatedTimeTable

  Background:
      Given a Referential "test" is created

  Scenario: 3950 - Handle a SIRI EstimatedTimeTable request
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Tutute                                   |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:25:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:26:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:27:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                               | Passage 32                              |
      | ObjectIDs                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                             | 6ba7b814-9dad-11d1-6-00c04fd430c8       |
      | Attribute[DirectionRef]            | Aller                                   |
      | Attribute[OriginName]              | Le début                                |
      | Attribute[DestinationName]         | La fin.                                 |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a StopVisit exists with the following attributes:
    #retard d'une minute
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | VehicleAtStop                   | false                                                                |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:01:00.000Z                                             |
      | ArrivalStatus                   | Delayed                                                              |
      | Schedule[aimed]#Departure       | 2017-01-01T15:01:00.000Z                                             |
      | Schedule[expected]#Departure    | 2017-01-01T15:02:00.000Z                                             |
      | DepartureStatus                 | Delayed                                                              |
      | Attribute[DestinationDisplay]   | Pouet-pouet                                                          |
    And a StopVisit exists with the following attributes:
    #retard d'une minute
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-2" |
      | PassageOrder                    | 5                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:05:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:06:00.000Z                                             |
      | ArrivalStatus                   | Delayed                                                              |
    And a StopVisit exists with the following attributes:
    #retard d'une minute
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:26:LOC-3" |
      | PassageOrder                    | 6                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:10:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:11:00.000Z                                             |
      | ArrivalStatus                   | Delayed                                                              |
    And a StopVisit exists with the following attributes:
    #à l'heure
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-4" |
      | PassageOrder                    | 7                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:16:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:16:00.000Z                                             |
      | ArrivalStatus                   | onTime                                                               |
    And I see ara vehicle_journeys
    And I see ara stop_visits
    And I see ara lines
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetEstimatedTimetable xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:siri="http://www.ifopt.org.uk/acsb"
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
          <sw:GetEstimatedTimetableResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-c-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>EstimatedTimetable:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>EstimatedTimetable:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:EstimatedJourneyVersionFrame>
                  <siri:RecordedAtTime>2017-01-01T12:00:00.000Z</siri:RecordedAtTime>
                  <siri:EstimatedVehicleJourney>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:DirectionRef>Aller</siri:DirectionRef>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        <siri:Order>4</siri:Order>
                        <siri:StopPointName>Tutute</siri:StopPointName>
                        <siri:VehicleAtStop>false</siri:VehicleAtStop>
                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                        <siri:DepartureStatus>Delayed</siri:DepartureStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                        <siri:Order>5</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:VehicleAtStop>false</siri:VehicleAtStop>
                        <siri:AimedArrivalTime>2017-01-01T15:05:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:06:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:26:LOC</siri:StopPointRef>
                        <siri:Order>6</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:VehicleAtStop>false</siri:VehicleAtStop>
                        <siri:AimedArrivalTime>2017-01-01T15:10:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:11:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:27:LOC</siri:StopPointRef>
                        <siri:Order>7</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:VehicleAtStop>false</siri:VehicleAtStop>
                        <siri:AimedArrivalTime>2017-01-01T15:16:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:16:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>onTime</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                    </siri:EstimatedCalls>
                  </siri:EstimatedVehicleJourney>
                </siri:EstimatedJourneyVersionFrame>
              </siri:EstimatedTimetableDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetEstimatedTimetableResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 3950 - Handle a SIRI EstimatedTimeTable request with the no rewrite setting
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential                           | test            |
      | remote_objectid_kind                       | internal        |
      | broadcast.no_destinationref_rewriting_from | NoRewriteOrigin |
    And a StopArea exists with the following attributes:
      | Name      | Tutute                                   |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:25:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:26:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:27:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Origin                             | NoRewriteOrigin                         |
      | Name                               | Passage 32                              |
      | ObjectIDs                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                             | 6ba7b814-9dad-11d1-6-00c04fd430c8       |
      | Attribute[DirectionRef]            | Aller                                   |
      | Attribute[OriginName]              | Le début                                |
      | Attribute[DestinationName]         | La fin.                                 |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a StopVisit exists with the following attributes:
    #retard d'une minute
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | VehicleAtStop                   | false                                                                |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:01:00.000Z                                             |
      | ArrivalStatus                   | Delayed                                                              |
      | Schedule[aimed]#Departure       | 2017-01-01T15:01:00.000Z                                             |
      | Schedule[expected]#Departure    | 2017-01-01T15:02:00.000Z                                             |
      | DepartureStatus                 | Delayed                                                              |
      | Attribute[DestinationDisplay]   | Pouet-pouet                                                          |
    And a StopVisit exists with the following attributes:
    #retard d'une minute
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-2" |
      | PassageOrder                    | 5                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:05:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:06:00.000Z                                             |
      | ArrivalStatus                   | Delayed                                                              |
    And a StopVisit exists with the following attributes:
    #retard d'une minute
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:26:LOC-3" |
      | PassageOrder                    | 6                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:10:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:11:00.000Z                                             |
      | ArrivalStatus                   | Delayed                                                              |
    And a StopVisit exists with the following attributes:
    #à l'heure
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-4" |
      | PassageOrder                    | 7                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:16:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:16:00.000Z                                             |
      | ArrivalStatus                   | onTime                                                               |
    And I see ara vehicle_journeys
    And I see ara stop_visits
    And I see ara lines
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetEstimatedTimetable xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:siri="http://www.ifopt.org.uk/acsb"
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
          <sw:GetEstimatedTimetableResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-c-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>EstimatedTimetable:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>EstimatedTimetable:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:EstimatedJourneyVersionFrame>
                  <siri:RecordedAtTime>2017-01-01T12:00:00.000Z</siri:RecordedAtTime>
                  <siri:EstimatedVehicleJourney>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:DirectionRef>Aller</siri:DirectionRef>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:DestinationRef>ThisIsTheEnd</siri:DestinationRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        <siri:Order>4</siri:Order>
                        <siri:StopPointName>Tutute</siri:StopPointName>
                        <siri:VehicleAtStop>false</siri:VehicleAtStop>
                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                        <siri:DepartureStatus>Delayed</siri:DepartureStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                        <siri:Order>5</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:VehicleAtStop>false</siri:VehicleAtStop>
                        <siri:AimedArrivalTime>2017-01-01T15:05:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:06:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:26:LOC</siri:StopPointRef>
                        <siri:Order>6</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:VehicleAtStop>false</siri:VehicleAtStop>
                        <siri:AimedArrivalTime>2017-01-01T15:10:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:11:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:27:LOC</siri:StopPointRef>
                        <siri:Order>7</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:VehicleAtStop>false</siri:VehicleAtStop>
                        <siri:AimedArrivalTime>2017-01-01T15:16:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:16:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>onTime</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                    </siri:EstimatedCalls>
                  </siri:EstimatedVehicleJourney>
                </siri:EstimatedJourneyVersionFrame>
              </siri:EstimatedTimetableDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetEstimatedTimetableResponse>
        </S:Body>
      </S:Envelope>
      """
