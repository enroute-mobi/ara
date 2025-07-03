Feature: Support SIRI EstimatedTimetable

  Background:
    Given a Referential "test" is created

  @ARA-1139
  Scenario: Handle a SIRI EstimatedTimetable request with partner setting siri.passage_order set to visit_number should display the VisitNumber tag instead of Order tag
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential   | test         |
      | remote_code_space  | internal     |
      | siri.passage_order | visit_number |
    And a StopArea exists with the following attributes:
      | Name            | Tutute                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
    #retard d'une minute
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
                <ns2:LineDirection>
                  <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                </ns2:LineDirection>
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
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                    <siri:DestinationName>La fin.</siri:DestinationName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        <siri:VisitNumber>4</siri:VisitNumber>
                        <siri:StopPointName>Tutute</siri:StopPointName>
                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                        <siri:DepartureStatus>delayed</siri:DepartureStatus>
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

  Scenario: 3950 - Handle a SIRI EstimatedTimetable request
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name            | Tutute                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:26:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:27:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-2 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:05:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:06:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:26:LOC-3 |
      | PassageOrder                | 6                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:10:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:11:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-4 |
      | PassageOrder                | 7                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:16:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:16:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And a VehicleJourney exists with the following attributes:
      | Name                           | TEST                              |
      | Codes[internal]                | WITHOUT:STOP:VISITS               |
      | LineId                         | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin 2                          |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
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
                <ns2:LineDirection>
                  <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                </ns2:LineDirection>
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
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-d-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                    <siri:DestinationName>La fin.</siri:DestinationName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        <siri:Order>4</siri:Order>
                        <siri:StopPointName>Tutute</siri:StopPointName>
                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                        <siri:DepartureStatus>delayed</siri:DepartureStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                        <siri:Order>5</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:05:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:06:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:26:LOC</siri:StopPointRef>
                        <siri:Order>6</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:10:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:11:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:27:LOC</siri:StopPointRef>
                        <siri:Order>7</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
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
    And an audit event should exist with these attributes:
      | Protocol        | siri                                                                                                                  |
      | Direction       | received                                                                                                              |
      | Status          | OK                                                                                                                    |
      | Type            | EstimatedTimetableRequest                                                                                             |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC","NINOXE:StopPoint:SP:25:LOC","NINOXE:StopPoint:SP:26:LOC","NINOXE:StopPoint:SP:27:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201"]                                                                                         |
      | Lines           | ["NINOXE:Line:3:LOC"]                                                                                                 |

  Scenario: 3950 - Handle a SIRI EstimatedTimetable request with a recorded call
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential                  | test     |
      | remote_code_space                 | internal |
      | broadcast.recorded_calls.duration | 1h       |
    And a StopArea exists with the following attributes:
      | Name            | Tutute                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:26:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:27:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-2 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:05:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:06:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:26:LOC-3 |
      | PassageOrder                | 6                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:10:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:11:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-4 |
      | PassageOrder                | 7                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T11:16:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T11:16:00.000Z                               |
      | ArrivalStatus               | arrived                                                |
      | DepartureStatus             | departed                                               |
    #à l'heure
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
                <ns2:LineDirection>
                  <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                </ns2:LineDirection>
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
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                    <siri:DestinationName>La fin.</siri:DestinationName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:RecordedCalls>
                      <siri:RecordedCall>
                          <siri:StopPointRef>NINOXE:StopPoint:SP:27:LOC</siri:StopPointRef>
                          <siri:Order>7</siri:Order>
                          <siri:StopPointName>Test</siri:StopPointName>
                          <siri:AimedArrivalTime>2017-01-01T11:16:00.000Z</siri:AimedArrivalTime>
                          <siri:ExpectedArrivalTime>2017-01-01T11:16:00.000Z</siri:ExpectedArrivalTime>
                          <siri:ArrivalStatus>arrived</siri:ArrivalStatus>
                          <siri:DepartureStatus>departed</siri:DepartureStatus>
                        </siri:RecordedCall>
                    </siri:RecordedCalls>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        <siri:Order>4</siri:Order>
                        <siri:StopPointName>Tutute</siri:StopPointName>
                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                        <siri:DepartureStatus>delayed</siri:DepartureStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                        <siri:Order>5</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:05:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:06:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:26:LOC</siri:StopPointRef>
                        <siri:Order>6</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:10:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:11:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
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

  Scenario: 3950 - Handle a SIRI EstimatedTimetable request with the no rewrite setting
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential                           | test            |
      | remote_code_space                          | internal        |
      | broadcast.no_destinationref_rewriting_from | NoRewriteOrigin |
    And a StopArea exists with the following attributes:
      | Name            | Tutute                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:26:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:27:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Origin                         | NoRewriteOrigin                   |
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-2 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:05:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:06:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:26:LOC-3 |
      | PassageOrder                | 6                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:10:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:11:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-4 |
      | PassageOrder                | 7                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-7-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:16:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:16:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    #à l'heure
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
                <ns2:LineDirection>
                  <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                </ns2:LineDirection>
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
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DestinationRef>ThisIsTheEnd</siri:DestinationRef>
                    <siri:DestinationName>La fin.</siri:DestinationName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        <siri:Order>4</siri:Order>
                        <siri:StopPointName>Tutute</siri:StopPointName>
                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                        <siri:DepartureStatus>delayed</siri:DepartureStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                        <siri:Order>5</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:05:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:06:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:26:LOC</siri:StopPointRef>
                        <siri:Order>6</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:10:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:11:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:27:LOC</siri:StopPointRef>
                        <siri:Order>7</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
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

  @ARA-1298
  Scenario: Handle a SIRI EstimatedTimetable request with Partner remote_code_space changed
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name            | Test 1                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test 2                     |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test 3                     |
      | Codes[internal] | NINOXE:StopPoint:SP:26:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test 4                     |
      | Codes[external] | NINOXE:StopPoint:SP:27:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a Line exists with the following attributes:
      | Codes[external] | NINOXE:Line:A:BUS:LOC |
      | Name            | Ligne A Bus           |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[external]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-2 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:05:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:06:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:26:LOC-3 |
      | PassageOrder                | 6                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:10:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:11:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[external]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-4 |
      | PassageOrder                | 7                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-9-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:16:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:16:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    #à l'heure
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
                <ns2:LineDirection>
                  <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                </ns2:LineDirection>
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
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-e-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                    <siri:DestinationName>La fin.</siri:DestinationName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        <siri:Order>4</siri:Order>
                        <siri:StopPointName>Test 1</siri:StopPointName>
                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                        <siri:DepartureStatus>delayed</siri:DepartureStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                        <siri:Order>5</siri:Order>
                        <siri:StopPointName>Test 2</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:05:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:06:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:26:LOC</siri:StopPointRef>
                        <siri:Order>6</siri:Order>
                        <siri:StopPointName>Test 3</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:10:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:11:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
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
    And the Partner "test" is updated with the following settings:
      | local_credential  | test     |
      | remote_code_space | external |
    And a minute has passed
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
                <ns2:LineDirection>
                  <ns2:LineRef>NINOXE:Line:A:BUS:LOC</ns2:LineRef>
                </ns2:LineDirection>
              </ns2:Lines>
            </Request>
            <RequestExtension />
          </ns7:GetEstimatedTimetable>
        </S:Body>
      </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
      <S:Body>
      <sw:GetEstimatedTimetableResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-f-00c04fd430c8</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>EstimatedTimetable:Test:0</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:EstimatedTimetableDelivery version='2.0:FR-IDF-2.4'>
          <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>EstimatedTimetable:Test:0</siri:RequestMessageRef>
          <siri:Status>true</siri:Status>
          <siri:EstimatedJourneyVersionFrame>
          <siri:RecordedAtTime>2017-01-01T12:01:00.000Z</siri:RecordedAtTime>
            <siri:EstimatedVehicleJourney>
              <siri:LineRef>NINOXE:Line:A:BUS:LOC</siri:LineRef>
              <siri:DirectionRef>Aller</siri:DirectionRef>
              <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne A Bus</siri:PublishedLineName>
              <siri:DestinationRef>a8989abce31bae21da02c1c2cf42dd855cd86a1d</siri:DestinationRef>
              <siri:DestinationName>La fin.</siri:DestinationName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:EstimatedCalls>
                <siri:EstimatedCall>
                  <siri:StopPointRef>NINOXE:StopPoint:SP:27:LOC</siri:StopPointRef>
                  <siri:Order>7</siri:Order>
                  <siri:StopPointName>Test 4</siri:StopPointName>
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

  @ARA-1363
  Scenario: Handle a SIRI EstimatedTimetable request using the generator setting reference_vehicle_journey_identifier
    # Setting a Partner without default generators
    Given a Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential                                | test                             |
      | remote_code_space                               | internal                         |
      | generators.reference_vehicle_journey_identifier | ch:1:ServiceJourney:87_TAC:%{id} |
    And a StopArea exists with the following attributes:
      | Name            | Tutute                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[_default]                | 6ba7b814                          |
      | Codes[external]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
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
                        <ns2:LineDirection>
                          <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                        </ns2:LineDirection>
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
                      <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:ResponseMessageIdentifier>
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
                                  <siri:DatedVehicleJourneyRef>ch:1:ServiceJourney:87_TAC:6ba7b814</siri:DatedVehicleJourneyRef>
                                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                                  <siri:DestinationRef>a8989abce31bae21da02c1c2cf42dd855cd86a1d</siri:DestinationRef>
                                  <siri:DestinationName>La fin.</siri:DestinationName>
                                  <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                                  <siri:EstimatedCalls>
                                    <siri:EstimatedCall>
                                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                                          <siri:Order>4</siri:Order>
                                          <siri:StopPointName>Tutute</siri:StopPointName>
                                          <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                                          <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                                          <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                                          <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                                          <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                                          <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                                          <siri:DepartureStatus>delayed</siri:DepartureStatus>
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
    And an audit event should exist with these attributes:
      | Protocol        | siri                                    |
      | Direction       | received                                |
      | Status          | OK                                      |
      | Type            | EstimatedTimetableRequest               |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"]          |
      | VehicleJourneys | ["ch:1:ServiceJourney:87_TAC:6ba7b814"] |
      | Lines           | ["NINOXE:Line:3:LOC"]                   |

  @ARA-1363
  Scenario: Handle a SIRI EstimatedTimetable request using the default generator should send DatedVehicleJourneyRef according to default setting
    # Setting a "SIRI Partner" with default generators
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name            | Tutute                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[_default]                | 6ba7b814                          |
      | Codes[external]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
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
                        <ns2:LineDirection>
                          <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                        </ns2:LineDirection>
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
                    <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
                                <siri:DatedVehicleJourneyRef>RATPDev:VehicleJourney::6ba7b814:LOC</siri:DatedVehicleJourneyRef>
                                <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                                <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                                <siri:DestinationName>La fin.</siri:DestinationName>
                                <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                                <siri:EstimatedCalls>
                                  <siri:EstimatedCall>
                                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                                        <siri:Order>4</siri:Order>
                                        <siri:StopPointName>Tutute</siri:StopPointName>
                                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                                        <siri:DepartureStatus>delayed</siri:DepartureStatus>
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
    And an audit event should exist with these attributes:
      | Protocol        | siri                                     |
      | Direction       | received                                 |
      | Status          | OK                                       |
      | Type            | EstimatedTimetableRequest                |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"]           |
      | VehicleJourneys | ["RATPDev:VehicleJourney::6ba7b814:LOC"] |
      | Lines           | ["NINOXE:Line:3:LOC"]                    |

  @ARA-1475
  Scenario: Handle a SIRI EstimatedTimetable request using the setting broadcast.prefer_referent_stop_areas should broadcast Referent StopArea
    Given a Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential                     | test     |
      | remote_code_space                    | internal |
      | broadcast.prefer_referent_stop_areas | true     |
    And a StopArea exists with the following attributes:
      | Name            | Referent                         |
      | Codes[internal] | NINOXE:StopPoint:SP:Referent:LOC |
      # 6ba7b814-9dad-11d1-2-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Name            | Test                              |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC        |
      | ReferentId      | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      # 6ba7b814-9dad-11d1-3-00c04fd430c8
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
      # 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[_default]                | 6ba7b814                          |
      | Codes[external]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
      # 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
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
                  <ns2:LineDirection>
                    <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                  </ns2:LineDirection>
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
                <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:ResponseMessageIdentifier>
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
                            <siri:DatedVehicleJourneyRef>VehicleJourney:6ba7b814</siri:DatedVehicleJourneyRef>
                            <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                            <siri:DestinationRef>a8989abce31bae21da02c1c2cf42dd855cd86a1d</siri:DestinationRef>
                            <siri:DestinationName>La fin.</siri:DestinationName>
                            <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                            <siri:EstimatedCalls>
                              <siri:EstimatedCall>
                                  <siri:StopPointRef>NINOXE:StopPoint:SP:Referent:LOC</siri:StopPointRef>
                                    <siri:Order>4</siri:Order>
                                    <siri:StopPointName>Referent</siri:StopPointName>
                                    <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                                    <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                                    <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                                    <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                                    <siri:DepartureStatus>delayed</siri:DepartureStatus>
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

  @ARA-1466
  Scenario: Handle a SIRI EstimatedTimetable request with a StopVisit with status Cancelled in the future should be broadcasted in EstimatedCalls
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential                  | test     |
      | remote_code_space                 | internal |
      | broadcast.recorded_calls.duration |       1h |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:27:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-4 |
      | PassageOrder                | 7                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T14:16:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T14:16:00.000Z                               |
      | ArrivalStatus               | cancelled                                              |
      | DepartureStatus             | cancelled                                              |
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
                <ns2:LineDirection>
                  <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                </ns2:LineDirection>
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
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                    <siri:DestinationName>La fin.</siri:DestinationName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                          <siri:StopPointRef>NINOXE:StopPoint:SP:27:LOC</siri:StopPointRef>
                          <siri:Order>7</siri:Order>
                          <siri:StopPointName>Test</siri:StopPointName>
                          <siri:AimedArrivalTime>2017-01-01T14:16:00.000Z</siri:AimedArrivalTime>
                          <siri:ExpectedArrivalTime>2017-01-01T14:16:00.000Z</siri:ExpectedArrivalTime>
                          <siri:ArrivalStatus>cancelled</siri:ArrivalStatus>
                          <siri:DepartureStatus>cancelled</siri:DepartureStatus>
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

  @ARA-1493
  Scenario: Handle referent lines in a SIRI EstimatedTimetable request
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential      | test     |
      | remote_code_space     | internal |
      | sort_payload_for_test | true     |
    And a StopArea exists with the following attributes:
      | Name            | Tutute                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:26:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:27:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a Line exists with the following attributes:
      | Codes[external] | NINOXE:Line:4:LOC                 |
      | ReferentId      | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | Name            | Ligne 3 Metro                     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-2 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:05:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:06:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:26:LOC-3 |
      | PassageOrder                | 6                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:10:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:11:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-4 |
      | PassageOrder                | 7                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:16:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:16:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    #à l'heure
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
                <ns2:LineDirection>
                  <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                </ns2:LineDirection>
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
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-d-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                    <siri:DestinationName>La fin.</siri:DestinationName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        <siri:Order>4</siri:Order>
                        <siri:StopPointName>Tutute</siri:StopPointName>
                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                        <siri:DepartureStatus>delayed</siri:DepartureStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                        <siri:Order>5</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:05:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:06:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:26:LOC</siri:StopPointRef>
                        <siri:Order>6</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:10:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:11:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:27:LOC</siri:StopPointRef>
                        <siri:Order>7</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
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
    And an audit event should exist with these attributes:
      | Protocol        | siri                                                                                                                  |
      | Direction       | received                                                                                                              |
      | Status          | OK                                                                                                                    |
      | Type            | EstimatedTimetableRequest                                                                                             |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC","NINOXE:StopPoint:SP:25:LOC","NINOXE:StopPoint:SP:26:LOC","NINOXE:StopPoint:SP:27:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201"]                                                                                         |
      | Lines           | ["NINOXE:Line:3:LOC"]                                                                                                 |

  @ARA-1493
  Scenario: Handle a referent line family in a SIRI EstimatedTimetable request
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential      | test     |
      | remote_code_space     | internal |
      | sort_payload_for_test | true     |
    And a StopArea exists with the following attributes:
      | Name            | Tutute                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:26:LOC |
      | Monitored       | true                       |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:27:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a Line exists with the following attributes:
      | Codes[external] | NINOXE:Line:4:LOC                 |
      | ReferentId      | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | Name            | Ligne 3 Metro                     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:202         |
      | LineId                         | 6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-2 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:05:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:06:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:26:LOC-3 |
      | PassageOrder                | 6                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-9-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:10:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:11:00.000Z                               |
      | ArrivalStatus               | delayed                                                |
    #retard d'une minute
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-4 |
      | PassageOrder                | 7                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-9-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:16:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:16:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    #à l'heure
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
                <ns2:LineDirection>
                  <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                </ns2:LineDirection>
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
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-e-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                    <siri:DestinationName>La fin.</siri:DestinationName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        <siri:Order>4</siri:Order>
                        <siri:StopPointName>Tutute</siri:StopPointName>
                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                        <siri:DepartureStatus>delayed</siri:DepartureStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                        <siri:Order>5</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:05:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:06:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                    </siri:EstimatedCalls>
                  </siri:EstimatedVehicleJourney>
                  <siri:EstimatedVehicleJourney>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:DirectionRef>Aller</siri:DirectionRef>
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                    <siri:DestinationName>La fin.</siri:DestinationName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:26:LOC</siri:StopPointRef>
                        <siri:Order>6</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:10:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:11:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                      </siri:EstimatedCall>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:27:LOC</siri:StopPointRef>
                        <siri:Order>7</siri:Order>
                        <siri:StopPointName>Test</siri:StopPointName>
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
    And an audit event should exist with these attributes:
      | Protocol        | siri                                                                                                                  |
      | Direction       | received                                                                                                              |
      | Status          | OK                                                                                                                    |
      | Type            | EstimatedTimetableRequest                                                                                             |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC","NINOXE:StopPoint:SP:25:LOC","NINOXE:StopPoint:SP:26:LOC","NINOXE:StopPoint:SP:27:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201","NINOXE:VehicleJourney:202"]                                                                                         |
      | Lines           | ["NINOXE:Line:3:LOC"]                                                                                                 |

  @ARA-1534
  Scenario: Handle a SIRI EstimatedTimetable request with vehicle occupancy
    Given a SIRI Partner "test" exists with connectors [siri-estimated-timetable-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name            | Tutute                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
      # id 6ba7b814-9dad-11d1-2-00c04fd430c8
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
      # id 6ba7b814-9dad-11d1-3-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
      # id 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes[internal]               | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                  | 4                                                      |
      | VehicleAtStop                 | false                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop                 | false                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                               |
      | ArrivalStatus                 | delayed                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                               |
      | DepartureStatus               | delayed                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                            |
      # id 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a Vehicle exists with the following attributes:
      | Codes[internal] | Test:Vehicle:1:LOC                |
      | Longitude       | 1.234                             |
      | Latitude        | 5.678                             |
      | Bearing         | 123                               |
      | Occupancy       | fewSeatsAvailable                 |
      | Percentage      | 15.6                              |
      | RecordedAtTime  | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime  | 2017-01-01T14:00:00.000Z          |
      | NextStopVisitId | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
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
                <ns2:LineDirection>
                  <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                </ns2:LineDirection>
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
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                    <siri:DestinationName>La fin.</siri:DestinationName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        <siri:Order>4</siri:Order>
                        <siri:StopPointName>Tutute</siri:StopPointName>
                        <siri:Occupancy>fewSeatsAvailable</siri:Occupancy>
                        <siri:DestinationDisplay>Pouet-pouet</siri:DestinationDisplay>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</siri:ExpectedArrivalTime>
                        <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
                        <siri:AimedDepartureTime>2017-01-01T15:01:00.000Z</siri:AimedDepartureTime>
                        <siri:ExpectedDepartureTime>2017-01-01T15:02:00.000Z</siri:ExpectedDepartureTime>
                        <siri:DepartureStatus>delayed</siri:DepartureStatus>
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
