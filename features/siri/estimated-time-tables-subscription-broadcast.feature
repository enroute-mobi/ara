Feature: Support SIRI EstimatedTimetable by subscription

  Background:
    Given a Referential "test" is created

  @ARA-1060
  Scenario: Handle a raw SIRI EstimatedTimetable request for subscription for all existing lines in a referential having same remote_code_space
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:A:BUS |
      | Name            | Ligne A Bus       |
    And a minute has passed
    When I send this SIRI request
      """
      <?xml version="1.0" encoding="utf-8"?>
      <Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
      <SubscriptionRequest>
      <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>
      <RequestorRef>NINOXE:default</RequestorRef>
      <EstimatedTimetableSubscriptionRequest>
         <SubscriptionIdentifier>test1</SubscriptionIdentifier>
         <InitialTerminationTime>2022-02-10T02:00:00Z</InitialTerminationTime>
         <EstimatedTimetableRequest>
            <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>
            <PreviewInterval>PT3H0S</PreviewInterval>
         </EstimatedTimetableRequest>
         <ChangeBeforeUpdates>PT30S</ChangeBeforeUpdates>
      </EstimatedTimetableSubscriptionRequest>
      </SubscriptionRequest>
      </Siri>
      """
    Then I should receive this SIRI response
      """
      <?xml version="1.0" encoding="utf-8"?>
      <Siri xmlns="http://www.siri.org.uk/siri" version="2.0">
      <SubscriptionResponse>
        <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
        <ResponderRef>test</ResponderRef>
        <ResponseStatus>
            <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
            <SubscriptionRef>test1</SubscriptionRef>
            <Status>true</Status>
            <ValidUntil>2022-02-10T02:00:00.000Z</ValidUntil>
        </ResponseStatus>
        <ServiceStartedTime>2017-01-01T12:00:00.000Z</ServiceStartedTime>
      </SubscriptionResponse>
      </Siri>
      """
    Then a Subscription exist with the following attributes:
      | internal | NINOXE:Line:3:LOC |
      | internal | NINOXE:Line:A:BUS |

  @ARA-1060
  Scenario: Handle a raw SIRI EstimatedTimetable request for subscription for all existing lines in a referential only with same remote_code_space
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And a Line exists with the following attributes:
      | Codes[another] | NINOXE:Line:3:LOC |
      | Name           | Ligne 3 Metro     |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:A:BUS |
      | Name            | Ligne A Bus       |
    And a minute has passed
    When I send this SIRI request
      """
      <?xml version="1.0" encoding="utf-8"?>
      <Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
      <SubscriptionRequest>
      <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>
      <RequestorRef>NINOXE:default</RequestorRef>
      <EstimatedTimetableSubscriptionRequest>
         <SubscriptionIdentifier>test1</SubscriptionIdentifier>
         <InitialTerminationTime>2022-02-10T02:00:00Z</InitialTerminationTime>
         <EstimatedTimetableRequest>
            <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>
            <PreviewInterval>PT3H0S</PreviewInterval>
         </EstimatedTimetableRequest>
         <ChangeBeforeUpdates>PT30S</ChangeBeforeUpdates>
      </EstimatedTimetableSubscriptionRequest>
      </SubscriptionRequest>
      </Siri>
      """
    Then I should receive this SIRI response
      """
      <?xml version="1.0" encoding="utf-8"?>
      <Siri xmlns="http://www.siri.org.uk/siri" version="2.0">
      <SubscriptionResponse>
        <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
        <ResponderRef>test</ResponderRef>
        <ResponseStatus>
            <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
            <SubscriptionRef>test1</SubscriptionRef>
            <Status>true</Status>
            <ValidUntil>2022-02-10T02:00:00.000Z</ValidUntil>
        </ResponseStatus>
        <ServiceStartedTime>2017-01-01T12:00:00.000Z</ServiceStartedTime>
      </SubscriptionResponse>
      </Siri>
      """
    Then Subscriptions exist with the following resources:
      | internal | NINOXE:Line:A:BUS |
    Then No Subscriptions exist with the following resources:
      | internal | NINOXE:Line:3:LOC |

  Scenario: 4234 - Handle a SOAP SIRI EstimatedTimetable request for subscription
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And a minute has passed
    When I send this SIRI request
      """
      <SOAP-ENV:Envelope
      xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
      <SOAP-ENV:Header/>
      <SOAP-ENV:Body>
      <ws:Subscribe xmlns:siri="http://www.siri.org.uk/siri" xmlns:ws="http://wsdl.siri.org.uk">
      <SubscriptionRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:01:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
      </SubscriptionRequestInfo>
      <Request>
        <siri:EstimatedTimetableSubscriptionRequest>
          <siri:SubscriptionIdentifier>NINOXE:default</siri:SubscriptionIdentifier>
          <siri:InitialTerminationTime>2017-01-01T13:00:00.000Z</siri:InitialTerminationTime>
          <siri:EstimatedTimetableRequest>
            <siri:RequestTimestamp>2017-01-01T12:01:00.000Z</siri:RequestTimestamp>
            <siri:PreviewInterval>PT23H</siri:PreviewInterval>
          </siri:EstimatedTimetableRequest>
          <siri:ChangeBeforeUpdates>PT3M</siri:ChangeBeforeUpdates>
        </siri:EstimatedTimetableSubscriptionRequest>
      </Request>
      <RequestExtension/>
      </ws:Subscribe>
      </SOAP-ENV:Body>
      </SOAP-ENV:Envelope>
      """
    Then one Subscription exists with the following attributes:
      | Kind | EstimatedTimetableBroadcast |

  @ARA-1025
  Scenario: Handle a raw SIRI EstimatedTimetable request for subscription
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And a minute has passed
    When I send this SIRI request
      """
      <?xml version="1.0" encoding="utf-8"?>
      <Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
      <SubscriptionRequest>
      <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>
      <RequestorRef>NINOXE:default</RequestorRef>
      <EstimatedTimetableSubscriptionRequest>
         <SubscriptionIdentifier>1</SubscriptionIdentifier>
         <InitialTerminationTime>2017-01-01T14:00:00Z</InitialTerminationTime>
         <EstimatedTimetableRequest>
            <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>
            <PreviewInterval>PT23H</PreviewInterval>
         </EstimatedTimetableRequest>
         <ChangeBeforeUpdates>PT30S</ChangeBeforeUpdates>
      </EstimatedTimetableSubscriptionRequest>
      </SubscriptionRequest>
      </Siri>
      """
    Then one Subscription exists with the following attributes:
      | Kind | EstimatedTimetableBroadcast |

  @ARA-1256
  Scenario: Delete and recreate subscription when receiving subscription with same existing number
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                         | http://localhost:8090 |
      | remote_credential                  | test                  |
      | local_credential                   | NINOXE:default        |
      | remote_code_space                  | internal              |
      | siri.envelope                      | raw                   |
      | broadcast.subscriptions.persistent | true                  |
    And a minute has passed
    When I send this SIRI request
      """
      <?xml version="1.0" encoding="utf-8"?>
      <Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
      <SubscriptionRequest>
      <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>
      <RequestorRef>NINOXE:default</RequestorRef>
      <EstimatedTimetableSubscriptionRequest>
         <SubscriptionIdentifier>1</SubscriptionIdentifier>
         <InitialTerminationTime>2017-01-01T14:00:00Z</InitialTerminationTime>
         <EstimatedTimetableRequest>
            <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>
            <PreviewInterval>PT23H</PreviewInterval>
         </EstimatedTimetableRequest>
         <ChangeBeforeUpdates>PT30S</ChangeBeforeUpdates>
      </EstimatedTimetableSubscriptionRequest>
      </SubscriptionRequest>
      </Siri>
      """
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Kind            | EstimatedTimetableBroadcast       |
      | ExternalId      |                                 1 |
    When I send this SIRI request
      """
      <?xml version="1.0" encoding="utf-8"?>
      <Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
      <SubscriptionRequest>
      <RequestTimestamp>2017-01-01T12:02:00.000Z</RequestTimestamp>
      <RequestorRef>NINOXE:default</RequestorRef>
      <EstimatedTimetableSubscriptionRequest>
         <SubscriptionIdentifier>1</SubscriptionIdentifier>
         <InitialTerminationTime>2017-01-01T14:00:00Z</InitialTerminationTime>
         <EstimatedTimetableRequest>
            <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>
            <PreviewInterval>PT23H</PreviewInterval>
         </EstimatedTimetableRequest>
         <ChangeBeforeUpdates>PT30S</ChangeBeforeUpdates>
      </EstimatedTimetableSubscriptionRequest>
      </SubscriptionRequest>
      </Siri>
      """
    Then No Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Kind            | EstimatedTimetableBroadcast       |
      | ExternalId      |                                 1 |
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Kind            | EstimatedTimetableBroadcast       |
      | ExternalId      |                                 1 |
    When I send this SIRI request
      """
      <?xml version="1.0" encoding="utf-8"?>
      <Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
      <SubscriptionRequest>
      <RequestTimestamp>2017-01-01T12:02:00.000Z</RequestTimestamp>
      <RequestorRef>NINOXE:default</RequestorRef>
      <EstimatedTimetableSubscriptionRequest>
         <SubscriptionIdentifier>2</SubscriptionIdentifier>
         <InitialTerminationTime>2017-01-01T14:00:00Z</InitialTerminationTime>
         <EstimatedTimetableRequest>
            <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>
            <PreviewInterval>PT23H</PreviewInterval>
         </EstimatedTimetableRequest>
         <ChangeBeforeUpdates>PT30S</ChangeBeforeUpdates>
      </EstimatedTimetableSubscriptionRequest>
      </SubscriptionRequest>
      </Siri>
      """
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Kind            | EstimatedTimetableBroadcast       |
      | ExternalId      |                                 1 |
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Kind            | EstimatedTimetableBroadcast       |
      | ExternalId      |                                 2 |

  Scenario: 4235 - Manage a ETT Notify after modification of a StopVisit
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[OriginRef]#Code      | "external": "ThisIsTheBeginning"  |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>

      </ServiceDeliveryInfo>
      <Notification>
      <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
        <siri:EstimatedVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
          <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
          <siri:OriginRef>RATPDev:StopPoint:Q:437bd677531592a8d5138228dd571274689b24e0:LOC</siri:OriginRef>
          <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
          <siri:DestinationName>La fin.</siri:DestinationName>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:EstimatedCalls>
            <siri:EstimatedCall>
              <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
              <siri:Order>4</siri:Order>
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
            </siri:EstimatedCall>
          </siri:EstimatedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
      </siri:EstimatedTimetableDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyEstimatedTimetable>
      </S:Body>
      </S:Envelope>
      """

  @ARA-1134
  Scenario: Manage a raw ETT notify after modification of a StopVisit broadcasting the PublishedLineName as line number
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | local_credential                  | NINOXE:default        |
      | remote_code_space                 | internal              |
      | broadcast.recorded_calls.duration | 1h                    |
      | siri.envelope                     | raw                   |
      | siri.line.published_name          | number                |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
      | Number          | L3M               |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
       <RecordedAtTime>2017-01-01T12:00:15.000Z</RecordedAtTime>
       <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
         <PublishedLineName>L3M</PublishedLineName>
         <DestinationName>La fin.</DestinationName>
         <OperatorRef>CdF:Company::410:LOC</OperatorRef>
         <EstimatedCalls>
           <EstimatedCall>
             <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
             <Order>4</Order>
             <StopPointName>Test</StopPointName>
             <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
             <ArrivalStatus>delayed</ArrivalStatus>
           </EstimatedCall>
         </EstimatedCalls>
         <IsCompleteStopSequence>false</IsCompleteStopSequence>
       </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  @ARA-1134
  Scenario: Manage a raw ETT notify after modification of a StopVisit with line not having a number and settings siri.line.published_name set to number
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | local_credential                  | NINOXE:default        |
      | remote_code_space                 | internal              |
      | broadcast.recorded_calls.duration |                    1h |
      | siri.envelope                     | raw                   |
      | siri.line.published_name          | number                |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
       <RecordedAtTime>2017-01-01T12:00:15.000Z</RecordedAtTime>
       <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
         <PublishedLineName>Ligne 3 Metro</PublishedLineName>
         <DestinationName>La fin.</DestinationName>
         <OperatorRef>CdF:Company::410:LOC</OperatorRef>
         <EstimatedCalls>
           <EstimatedCall>
             <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
             <Order>4</Order>
             <StopPointName>Test</StopPointName>
             <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
             <ArrivalStatus>delayed</ArrivalStatus>
           </EstimatedCall>
         </EstimatedCalls>
         <IsCompleteStopSequence>false</IsCompleteStopSequence>
       </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  @ARA-1139
  Scenario: Manage a raw ETT notify after modification of a StopVisit using the settings siri.passage_order set to visit_number should display the VisitNumber tag instead of Order tag
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url         | http://localhost:8090 |
      | remote_credential  | test                  |
      | local_credential   | NINOXE:default        |
      | remote_code_space  | internal              |
      | siri.envelope      | raw                   |
      | siri.passage_order | visit_number          |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
       <RecordedAtTime>2017-01-01T12:00:15.000Z</RecordedAtTime>
       <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
         <PublishedLineName>Ligne 3 Metro</PublishedLineName>
         <DestinationName>La fin.</DestinationName>
         <OperatorRef>CdF:Company::410:LOC</OperatorRef>
         <EstimatedCalls>
           <EstimatedCall>
             <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
             <VisitNumber>4</VisitNumber>
             <StopPointName>Test</StopPointName>
             <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
             <ArrivalStatus>delayed</ArrivalStatus>
           </EstimatedCall>
         </EstimatedCalls>
         <IsCompleteStopSequence>false</IsCompleteStopSequence>
       </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  @ARA-1243
  Scenario: Manage a raw ETT Notify after modification of a StopVisit with StopVisit departure time within the broadcast.recorded_calls.duration must order StopVisits by Order
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | local_credential                  | NINOXE:default        |
      | remote_code_space                 | internal              |
      | broadcast.recorded_calls.duration | 1h                    |
      | siri.envelope                     | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name            | Test1                      |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name            | Test2                      |
      | Codes[internal] | NINOXE:StopPoint:SP:26:LOC |
      # "Id":"6ba7b814-9dad-11d1-5-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name            | Test3                      |
      | Codes[internal] | NINOXE:StopPoint:SP:27:LOC |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
      # "Id":"6ba7b814-9dad-11d1-8-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 1                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T01:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T01:02:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
      # "Id":"6ba7b814-9dad-11d1-9-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1 |
      | PassageOrder                | 2                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::420:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:05:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
      # "Id":"6ba7b814-9dad-11d1-a-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes[internal] | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:26:LOC-1" |
      | PassageOrder                |                                                                    3 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-8-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::430:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T12:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:10:00.000Z |
      | ArrivalStatus               | onTime                                                               |
      # "Id":"6ba7b814-9dad-11d1-b-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes[internal] | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:27:LOC-1" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-8-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::440:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T12:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:15:00.000Z |
      | ArrivalStatus               | onTime                                                               |
      # "Id":"6ba7b814-9dad-11d1-c-00c04fd430c8"
    And 5 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-9-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival   | 2017-01-01T11:01:01.000Z |
      | ArrivalStatus                | arrived                  |
      | DepartureStatus              | departed                 |
      | Schedule[expected]#Departure | 2017-01-01T11:01:11.000Z |
    When the StopVisit "6ba7b814-9dad-11d1-a-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival   | 2017-01-01T15:11:01.000Z |
      | ArrivalStatus                | arrived                  |
      | DepartureStatus              | departed                 |
      | Schedule[expected]#Departure | 2017-01-01T15:15:11.000Z |
    When the StopVisit "6ba7b814-9dad-11d1-b-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:16:01.000Z |
      | ArrivalStatus              | delayed                  |
    When the StopVisit "6ba7b814-9dad-11d1-c-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:26:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:10.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-e-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:10.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
      <RecordedAtTime>2017-01-01T12:00:10.000Z</RecordedAtTime>
      <EstimatedVehicleJourney>
        <LineRef>NINOXE:Line:3:LOC</LineRef>
        <DirectionRef>Aller</DirectionRef>
        <FramedVehicleJourneyRef>
          <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
          <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
        </FramedVehicleJourneyRef>
        <PublishedLineName>Ligne 3 Metro</PublishedLineName>
        <DestinationName>La fin.</DestinationName>
        <OperatorRef>CdF:Company::410:LOC</OperatorRef>
        <RecordedCalls>
          <RecordedCall>
            <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
            <Order>1</Order>
            <StopPointName>Test</StopPointName>
            <ExpectedArrivalTime>2017-01-01T11:01:01.000Z</ExpectedArrivalTime>
            <ArrivalStatus>arrived</ArrivalStatus>
            <ExpectedDepartureTime>2017-01-01T11:01:11.000Z</ExpectedDepartureTime>
            <DepartureStatus>departed</DepartureStatus>
          </RecordedCall>
          <RecordedCall>
            <StopPointRef>NINOXE:StopPoint:SP:25:LOC</StopPointRef>
            <Order>2</Order>
            <StopPointName>Test1</StopPointName>
            <ExpectedArrivalTime>2017-01-01T15:11:01.000Z</ExpectedArrivalTime>
            <ArrivalStatus>arrived</ArrivalStatus>
            <ExpectedDepartureTime>2017-01-01T15:15:11.000Z</ExpectedDepartureTime>
            <DepartureStatus>departed</DepartureStatus>
          </RecordedCall>
        </RecordedCalls>
        <EstimatedCalls>
          <EstimatedCall>
            <StopPointRef>NINOXE:StopPoint:SP:26:LOC</StopPointRef>
            <Order>3</Order>
            <StopPointName>Test2</StopPointName>
            <ExpectedArrivalTime>2017-01-01T15:16:01.000Z</ExpectedArrivalTime>
            <ArrivalStatus>delayed</ArrivalStatus>
          </EstimatedCall>
          <EstimatedCall>
            <StopPointRef>NINOXE:StopPoint:SP:27:LOC</StopPointRef>
            <Order>4</Order>
            <StopPointName>Test3</StopPointName>
            <ExpectedArrivalTime>2017-01-01T15:26:01.000Z</ExpectedArrivalTime>
            <ArrivalStatus>delayed</ArrivalStatus>
          </EstimatedCall>
        </EstimatedCalls>
        <IsCompleteStopSequence>false</IsCompleteStopSequence>
      </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """
    And an audit event should exist with these attributes:
      | Protocol        | siri                                                                                                                  |
      | Direction       | sent                                                                                                                  |
      | Status          | OK                                                                                                                    |
      | Type            | NotifyEstimatedTimetable                                                                                              |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC","NINOXE:StopPoint:SP:25:LOC","NINOXE:StopPoint:SP:26:LOC","NINOXE:StopPoint:SP:27:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201"]                                                                                         |
      | Lines           | ["NINOXE:Line:3:LOC"]                                                                                                 |

  @ARA-1126
  Scenario: Manage a raw ETT Notify after modification of a StopVisit with StopVisit departure time within the broadcast.recorded_calls.duration
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | local_credential                  | NINOXE:default        |
      | remote_code_space                 | internal              |
      | broadcast.recorded_calls.duration | 1h                    |
      | siri.envelope                     | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a StopArea exists with the following attributes:
      | Name            | Test1                      |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T12:20:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 5 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-7-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:10.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:10.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
       <RecordedAtTime>2017-01-01T12:00:10.000Z</RecordedAtTime>
       <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
         <PublishedLineName>Ligne 3 Metro</PublishedLineName>
         <DestinationName>La fin.</DestinationName>
         <OperatorRef>CdF:Company::410:LOC</OperatorRef>
         <EstimatedCalls>
           <EstimatedCall>
             <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
             <Order>4</Order>
             <StopPointName>Test</StopPointName>
             <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
             <ArrivalStatus>delayed</ArrivalStatus>
           </EstimatedCall>
         </EstimatedCalls>
         <IsCompleteStopSequence>false</IsCompleteStopSequence>
       </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """
    When the StopVisit "6ba7b814-9dad-11d1-8-00c04fd430c8" is edited with the following attributes:
      | ArrivalStatus | delayed |
    When the StopVisit "6ba7b814-9dad-11d1-7-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival   | 2017-01-01T11:40:01.000Z |
      | ArrivalStatus                | arrived                  |
      | DepartureStatus              | departed                 |
      | Schedule[expected]#Departure | 2017-01-01T11:45:11.000Z |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-b-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
        <RecordedAtTime>2017-01-01T12:00:15.000Z</RecordedAtTime>
        <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
          <PublishedLineName>Ligne 3 Metro</PublishedLineName>
          <DestinationName>La fin.</DestinationName>
          <OperatorRef>CdF:Company::410:LOC</OperatorRef>
          <RecordedCalls>
            <RecordedCall>
              <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
              <Order>4</Order>
              <StopPointName>Test</StopPointName>
              <ExpectedArrivalTime>2017-01-01T11:40:01.000Z</ExpectedArrivalTime>
              <ArrivalStatus>arrived</ArrivalStatus>
              <ExpectedDepartureTime>2017-01-01T11:45:11.000Z</ExpectedDepartureTime>
              <DepartureStatus>departed</DepartureStatus>
            </RecordedCall>
          </RecordedCalls>
          <EstimatedCalls>
            <EstimatedCall>
              <StopPointRef>NINOXE:StopPoint:SP:25:LOC</StopPointRef>
              <Order>5</Order>
              <StopPointName>Test1</StopPointName>
              <AimedArrivalTime>2017-01-01T12:00:00.000Z</AimedArrivalTime>
              <ExpectedArrivalTime>2017-01-01T15:00:00.000Z</ExpectedArrivalTime>
              <ArrivalStatus>delayed</ArrivalStatus>
            </EstimatedCall>
          </EstimatedCalls>
          <IsCompleteStopSequence>false</IsCompleteStopSequence>
        </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """
    When the StopVisit "6ba7b814-9dad-11d1-8-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival   | 2017-01-01T11:51:01.000Z |
      | ArrivalStatus                | arrived                  |
      | DepartureStatus              | departed                 |
      | Schedule[expected]#Departure | 2017-01-01T11:52:11.000Z |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-c-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
        <RecordedAtTime>2017-01-01T12:00:20.000Z</RecordedAtTime>
        <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
          <PublishedLineName>Ligne 3 Metro</PublishedLineName>
          <DestinationName>La fin.</DestinationName>
          <OperatorRef>CdF:Company::410:LOC</OperatorRef>
          <RecordedCalls>
            <RecordedCall>
              <StopPointRef>NINOXE:StopPoint:SP:25:LOC</StopPointRef>
              <Order>5</Order>
              <StopPointName>Test1</StopPointName>
              <ExpectedArrivalTime>2017-01-01T11:51:01.000Z</ExpectedArrivalTime>
              <ArrivalStatus>arrived</ArrivalStatus>
              <ExpectedDepartureTime>2017-01-01T11:52:11.000Z</ExpectedDepartureTime>
              <DepartureStatus>departed</DepartureStatus>
            </RecordedCall>
          </RecordedCalls>
          <IsCompleteStopSequence>false</IsCompleteStopSequence>
        </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  @ARA-1062
  Scenario: Manage a ETT Notify after modification of a StopVisit with StopVisit departure time within the broadcast.recorded_calls.duration
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | local_credential                  | NINOXE:default        |
      | remote_code_space                 | internal              |
      | broadcast.recorded_calls.duration | 1h                    |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:15.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>

      </ServiceDeliveryInfo>
      <Notification>
      <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:15.000Z</siri:ResponseTimestamp>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:15.000Z</siri:RecordedAtTime>
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
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
            </siri:EstimatedCall>
          </siri:EstimatedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
      </siri:EstimatedTimetableDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyEstimatedTimetable>
      </S:Body>
      </S:Envelope>
      """
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival   | 2017-01-01T11:31:01.000Z |
      | ArrivalStatus                | arrived                  |
      | DepartureStatus              | departed                 |
      | Schedule[expected]#Departure | 2017-01-01T11:32:11.000Z |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>

      </ServiceDeliveryInfo>
      <Notification>
      <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
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
              <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
              <siri:Order>4</siri:Order>
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:ExpectedArrivalTime>2017-01-01T11:31:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>arrived</siri:ArrivalStatus>
              <siri:ExpectedDepartureTime>2017-01-01T11:32:11.000Z</siri:ExpectedDepartureTime>
              <siri:DepartureStatus>departed</siri:DepartureStatus>
            </siri:RecordedCall>
          </siri:RecordedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
      </siri:EstimatedTimetableDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyEstimatedTimetable>
      </S:Body>
      </S:Envelope>
      """

  Scenario: Manage a raw ETT Notify after modification of a StopVisit with StopVisit departure time oustide the broadcast.recorded_calls.duration
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | local_credential                  | NINOXE:default        |
      | remote_code_space                 | internal              |
      | broadcast.recorded_calls.duration | 1h                    |
      | siry.envelope                     | raw                   |
    Given a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]              | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                 | 4                                                      |
      | StopAreaId                   | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId             | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop                | false                                                  |
      | Reference[OperatorRef]#Code  | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival      | 2017-01-01T01:00:00.000Z                               |
      | Schedule[expected]#Arrival   | 2017-01-01T01:10:00.000Z                               |
      | ArrivalStatus                | onTime                                                 |
      | Schedule[expected]#Departure | 2017-01-01T01:20:00.000Z                               |
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | ArrivalStatus   | arrived  |
      | DepartureStatus | departed |
    And 10 seconds have passed
    Then the SIRI server should not have received a NotifyEstimatedTimetable request

  Scenario: Manage a ETT Notify after modification of a StopVisit with StopVisit departure time oustide the broadcast.recorded_calls.duration
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | local_credential                  | NINOXE:default        |
      | remote_code_space                 | internal              |
      | broadcast.recorded_calls.duration | 1h                    |
    Given a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]              | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                 | 4                                                      |
      | StopAreaId                   | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId             | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop                | false                                                  |
      | Reference[OperatorRef]#Code  | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival      | 2017-01-01T01:00:00.000Z                               |
      | Schedule[expected]#Arrival   | 2017-01-01T01:10:00.000Z                               |
      | ArrivalStatus                | onTime                                                 |
      | Schedule[expected]#Departure | 2017-01-01T01:20:00.000Z                               |
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | ArrivalStatus   | arrived  |
      | DepartureStatus | departed |
    And 10 seconds have passed
    Then the SIRI server should not have received a NotifyEstimatedTimetable request

  @ARA-1107
  Scenario: Manage a raw ETT Notify with StopArea having a Parent with with Partner Code after modification of a StopVisit
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | didok                 |
      | siri.envelope     | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast        |
      | ExternalId        | externalId                         |
      | SubscriberRef     | subscriber                         |
      | ReferenceArray[0] | Line, "didok": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name         | Parent                                  |
      | Codes[didok] | fr:1:StopPlace:OURA2:StopArea:log351672 |
      | Monitored    | true                                    |
    And a StopArea exists with the following attributes:
      | Name            | Child1                            |
      | Codes[internal] | vlgabon1                          |
      | ParentId        | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                              |
    And a Line exists with the following attributes:
      | Codes[didok] | NINOXE:Line:3:LOC |
      | Name         | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[didok]                   | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | DirectionType                  | outbound                          |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "internal": "ThisIsTheEnd"        |
      | Attribute[VehicleMode]         | bus                               |
    And a StopVisit exists with the following attributes:
      | PassageOrder                | 4                                 |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | false                             |
      | Reference[OperatorRef]#Code | "didok": "CdF:Company::410:LOC"   |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:00:00.000Z          |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z          |
      | ArrivalStatus               | onTime                            |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-7-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
      <RecordedAtTime>2017-01-01T12:00:20.000Z</RecordedAtTime>
      <EstimatedVehicleJourney>
      <LineRef>NINOXE:Line:3:LOC</LineRef>
      <DirectionRef>outbound</DirectionRef>
      <FramedVehicleJourneyRef>
        <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
        <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
      </FramedVehicleJourneyRef>
      <VehicleMode>bus</VehicleMode>
      <PublishedLineName>Ligne 3 Metro</PublishedLineName>
      <DestinationName>La fin.</DestinationName>
      <OperatorRef>CdF:Company::410:LOC</OperatorRef>
      <EstimatedCalls>
        <EstimatedCall>
          <StopPointRef>fr:1:StopPlace:OURA2:StopArea:log351672</StopPointRef>
      <Order>4</Order>
          <StopPointName>Parent</StopPointName>
          <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
          <ArrivalStatus>delayed</ArrivalStatus>
        </EstimatedCall>
      </EstimatedCalls>
      <IsCompleteStopSequence>false</IsCompleteStopSequence>
      </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  @ARA-1025
  Scenario: Manage a raw ETT Notify after modification of a StopVisit
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | outbound                          |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
      | Attribute[VehicleMode]         | bus                               |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
      <RecordedAtTime>2017-01-01T12:00:20.000Z</RecordedAtTime>
      <EstimatedVehicleJourney>
      <LineRef>NINOXE:Line:3:LOC</LineRef>
      <DirectionRef>outbound</DirectionRef>
      <FramedVehicleJourneyRef>
        <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
        <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
      </FramedVehicleJourneyRef>
      <VehicleMode>bus</VehicleMode>
      <PublishedLineName>Ligne 3 Metro</PublishedLineName>
      <DestinationName>La fin.</DestinationName>
      <OperatorRef>CdF:Company::410:LOC</OperatorRef>
      <EstimatedCalls>
        <EstimatedCall>
          <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
      <Order>4</Order>
          <StopPointName>Test</StopPointName>
          <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
          <ArrivalStatus>delayed</ArrivalStatus>
        </EstimatedCall>
      </EstimatedCalls>
      <IsCompleteStopSequence>false</IsCompleteStopSequence>
      </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  Scenario: 4235 - Manage a ETT Notify after modification of a StopVisit withe the no rewrite setting
    Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
       <ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
         <SubscriptionAnswerInfo
           xmlns:ns2="http://www.ifopt.org.uk/acsb"
           xmlns:ns3="http://www.ifopt.org.uk/ifopt"
           xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
           xmlns:ns5="http://www.siri.org.uk/siri"
           xmlns:ns6="http://wsdl.siri.org.uk/siri">
           <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
           <ns5:Address>http://appli.chouette.mobi/siri_france/siri</ns5:Address>
           <ns5:ResponderRef>NINOXE:default</ns5:ResponderRef>
           <ns5:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns5:MessageRefStructure">Subscription:Test:0</ns5:RequestMessageRef>
         </SubscriptionAnswerInfo>
         <Answer
           xmlns:ns2="http://www.ifopt.org.uk/acsb"
           xmlns:ns3="http://www.ifopt.org.uk/ifopt"
           xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
           xmlns:ns5="http://www.siri.org.uk/siri"
           xmlns:ns6="http://wsdl.siri.org.uk/siri">
           <ns5:ResponseStatus>
               <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
               <ns5:RequestMessageRef>Subscription:Test:0</ns5:RequestMessageRef>
               <ns5:SubscriberRef>SubscriberRef</ns5:SubscriberRef>
               <ns5:SubscriptionRef>SubscriptionIdentifier</ns5:SubscriptionRef>
               <ns5:Status>true</ns5:Status>
               <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
           </ns5:ResponseStatus>
           <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
         </Answer>
         <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
       </ns1:SubscribeResponse>
      </S:Body>
      </S:Envelope>
      """
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                                 | http://localhost:8090 |
      | remote_credential                          | test                  |
      | local_credential                           | NINOXE:default        |
      | remote_code_space                          | internal              |
      | broadcast.no_destinationref_rewriting_from | NoRewriteOrigin       |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Origin                         | NoRewriteOrigin                   |
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
    </ServiceDeliveryInfo>
      <Notification>
      <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
       <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
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
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
            </siri:EstimatedCall>
          </siri:EstimatedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
      </siri:EstimatedTimetableDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyEstimatedTimetable>
      </S:Body>
      </S:Envelope>
      """

  @ARA-1086
  Scenario: Handle a raw SIRI error if subscriptions are made using same ExternalId
    Given a raw SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-production-timetable-subscription-broadcaster] and the following settings:
      | remote_url                         | http://localhost:8090 |
      | remote_credential                  | ara                   |
      | local_credential                   | test                  |
      | remote_code_space                  | internal              |
      | siri.envelope                      | raw                   |
      | broadcast.subscriptions.persistent | true                  |
    And a StopArea exists with the following attributes:
      | Name            | Test 24                               |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC            |
      | Lines           | ["6ba7b814-9dad-11d1-4-00c04fd430c8"] |
    And a StopArea exists with the following attributes:
      | Name            | Test 25                               |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC            |
      | Lines           | ["6ba7b814-9dad-11d1-4-00c04fd430c8"] |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | SpecialExternalId                     |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    When a minute has passed
    And I send this SIRI request
      """
      <?xml version="1.0" encoding="utf-8"?>
      <Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
      <SubscriptionRequest>
      <RequestTimestamp>2022-02-09T02:15:23.690717Z</RequestTimestamp>
      <RequestorRef>test</RequestorRef>
      <ProductionTimetableSubscriptionRequest>
         <SubscriptionIdentifier>SpecialExternalId</SubscriptionIdentifier>
         <InitialTerminationTime>2022-02-10T02:50:00Z</InitialTerminationTime>
         <ProductionTimetableRequest>
            <RequestTimestamp>2022-02-09T02:15:23.690717Z</RequestTimestamp>
            <ValidityPeriod>
               <StartTime>2022-02-09T03:30:00Z</StartTime>
               <EndTime>2022-02-10T04:30:00Z</EndTime>
            </ValidityPeriod>
         </ProductionTimetableRequest>
      </ProductionTimetableSubscriptionRequest>
      </SubscriptionRequest>
      </Siri>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <SubscriptionResponse>
      <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
      <ResponderRef>ara</ResponderRef>
      <ResponseStatus>
      <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
      <SubscriptionRef>SpecialExternalId</SubscriptionRef>
      <Status>false</Status>
      <ErrorCondition>
        <OtherError number="2">
          <ErrorText>[BAD_REQUEST] Subscription Id SpecialExternalId already exists</ErrorText>
        </OtherError>
      </ErrorCondition>
      </ResponseStatus>
      <ServiceStartedTime>2017-01-01T12:00:00.000Z</ServiceStartedTime>
      </SubscriptionResponse>
      </Siri>
      """
    And one Subscription exists with the following attributes:
      | Kind | EstimatedTimetableBroadcast |

  @ARA-1101
  Scenario: Manage a raw ETT Notify after modification of a StopVisit with partner setting siri.direction_type should broadcast the DirectionRef with setting value
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url          | http://localhost:8090             |
      | remote_credential   | test                              |
      | local_credential    | NINOXE:default                    |
      | remote_code_space   | internal                          |
      | siri.envelope       | raw                               |
      | siri.direction_type | ch:1:Direction:R,ch:1:Direction:H |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | outbound                          |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
      | Attribute[VehicleMode]         | bus                               |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
      <RecordedAtTime>2017-01-01T12:00:20.000Z</RecordedAtTime>
      <EstimatedVehicleJourney>
      <LineRef>NINOXE:Line:3:LOC</LineRef>
      <DirectionRef>ch:1:Direction:H</DirectionRef>
      <FramedVehicleJourneyRef>
        <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
        <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
      </FramedVehicleJourneyRef>
      <VehicleMode>bus</VehicleMode>
      <PublishedLineName>Ligne 3 Metro</PublishedLineName>
      <DestinationName>La fin.</DestinationName>
      <OperatorRef>CdF:Company::410:LOC</OperatorRef>
      <EstimatedCalls>
        <EstimatedCall>
          <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
      <Order>4</Order>
          <StopPointName>Test</StopPointName>
          <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
          <ArrivalStatus>delayed</ArrivalStatus>
        </EstimatedCall>
      </EstimatedCalls>
      <IsCompleteStopSequence>false</IsCompleteStopSequence>
      </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  @ARA-1219
  Scenario: Check IsCompleteSequence if we son't broadcast an old StopVisit
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | local_credential                  | NINOXE:default        |
      | remote_code_space                 | internal              |
      | broadcast.recorded_calls.duration | 1h                    |
      | siri.envelope                     | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a StopArea exists with the following attributes:
      | Name            | Test1                      |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]              | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                 | 4                                                      |
      | StopAreaId                   | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId             | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop                | false                                                  |
      | Reference[OperatorRef]#Code  | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival      | 2017-01-01T01:00:00.000Z                               |
      | Schedule[expected]#Arrival   | 2017-01-01T01:10:00.000Z                               |
      | ArrivalStatus                | onTime                                                 |
      | Schedule[expected]#Departure | 2017-01-01T01:20:00.000Z                               |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::420:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:20:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:30:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
      # "Id":"6ba7b814-9dad-11d1-8-00c04fd430c8"
    And 5 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-8-00c04fd430c8" is edited with the following attributes:
      | ArrivalStatus | arrived |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:10.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:10.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
        <RecordedAtTime>2017-01-01T12:00:10.000Z</RecordedAtTime>
        <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
          <PublishedLineName>Ligne 3 Metro</PublishedLineName>
          <DestinationName>La fin.</DestinationName>
          <OperatorRef>CdF:Company::420:LOC</OperatorRef>
          <EstimatedCalls>
            <EstimatedCall>
              <StopPointRef>NINOXE:StopPoint:SP:25:LOC</StopPointRef>
              <Order>5</Order>
              <StopPointName>Test1</StopPointName>
              <AimedArrivalTime>2017-01-01T12:20:00.000Z</AimedArrivalTime>
              <ExpectedArrivalTime>2017-01-01T15:30:00.000Z</ExpectedArrivalTime>
              <ArrivalStatus>arrived</ArrivalStatus>
            </EstimatedCall>
          </EstimatedCalls>
          <IsCompleteStopSequence>false</IsCompleteStopSequence>
        </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  @ARA-1219
  Scenario: Check IsCompleteSequence if we son't broadcast an old StopVisit
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a StopArea exists with the following attributes:
      | Name            | Test1                      |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a ScheduledStopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T01:00:00.000Z                               |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a ScheduledStopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::420:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:20:00.000Z                               |
      # "Id":"6ba7b814-9dad-11d1-8-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::420:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:20:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:30:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
      # "Id":"6ba7b814-9dad-11d1-9-00c04fd430c8"
    And 5 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-9-00c04fd430c8" is edited with the following attributes:
      | ArrivalStatus | arrived |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:10.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-b-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:10.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
        <RecordedAtTime>2017-01-01T12:00:10.000Z</RecordedAtTime>
        <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
          <PublishedLineName>Ligne 3 Metro</PublishedLineName>
          <DestinationName>La fin.</DestinationName>
          <OperatorRef>CdF:Company::420:LOC</OperatorRef>
          <EstimatedCalls>
            <EstimatedCall>
              <StopPointRef>NINOXE:StopPoint:SP:25:LOC</StopPointRef>
              <Order>5</Order>
              <StopPointName>Test1</StopPointName>
              <AimedArrivalTime>2017-01-01T12:20:00.000Z</AimedArrivalTime>
              <ExpectedArrivalTime>2017-01-01T15:30:00.000Z</ExpectedArrivalTime>
              <ArrivalStatus>arrived</ArrivalStatus>
            </EstimatedCall>
          </EstimatedCalls>
          <IsCompleteStopSequence>false</IsCompleteStopSequence>
        </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  @ARA-1234
  Scenario: Broadcast full ETT StopVisits when receiving all stops and with ScheduledStopVisits
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | default                        |
      | remote_code_space                  | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
      | siri.envelope                      | raw                            |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name            | FIRST                      |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name            | SECOND                     |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
      # "SubscriptionRef":"RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC"
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
      # "SubscriptionRef":"RELAIS:Subscription::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC"
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                        |
      | Codes[internal]          | NINOXE:VehicleJourney:201         |
      | LineId                   | 6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | Monitored                | true                              |
      | Attribute[DirectionName] | A Direction Name                  |
      # "Id":"6ba7b814-9dad-11d1-8-00c04fd430c8"
    And a ScheduledStopVisit exists with the following attributes:
      | Codes[internal]         | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder            | 4                                                      |
      | StopAreaId              | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId        | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop           | false                                                  |
      | Schedule[aimed]#Arrival | 2017-01-01T15:00:00.000Z                               |
      # "Id":"6ba7b814-9dad-11d1-9-00c04fd430c8"
    And a ScheduledStopVisit exists with the following attributes:
      | Codes[internal]         | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-3 |
      | PassageOrder            | 5                                                      |
      | StopAreaId              | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId        | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop           | false                                                  |
      | Schedule[aimed]#Arrival | 2017-01-01T15:20:00.000Z                               |
      # "Id":"6ba7b814-9dad-11d1-d-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes[internal]            | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder               | 4                                                      |
      | StopAreaId                 | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop              | true                                                   |
      | Schedule[aimed]#Arrival    | 2017-01-01T13:00:00.000+02:00                          |
      | Schedule[expected]#Arrival | 2017-01-01T13:01:00.000+02:00                          |
    And 10 seconds have passed
    When the VehicleJourney "6ba7b814-9dad-11d1-8-00c04fd430c8" is edited with the following attributes:
      | HasCompleteStopSequence | true |
    And a StopVisit exists with the following attributes:
      | Codes[internal]            | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-3 |
      | PassageOrder               | 5                                                      |
      | StopAreaId                 | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop              | false                                                  |
      | Schedule[aimed]#Arrival    | 2017-01-01T15:00:00.000+02:00                          |
      | Schedule[expected]#Arrival | 2017-01-01T15:01:00.000+02:00                          |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
        <ServiceDelivery>
          <ResponseTimestamp>2017-01-01T12:00:50.000Z</ResponseTimestamp>
          <ProducerRef>test</ProducerRef>
          <ResponseMessageIdentifier>6ba7b814-9dad-11d1-e-00c04fd430c8</ResponseMessageIdentifier>
          <EstimatedTimetableDelivery>
            <ResponseTimestamp>2017-01-01T12:00:50.000Z</ResponseTimestamp>
            <SubscriberRef>subscriber</SubscriberRef>
            <SubscriptionRef>externalId</SubscriptionRef>
            <Status>true</Status>
            <EstimatedJourneyVersionFrame>
              <RecordedAtTime>2017-01-01T12:00:50.000Z</RecordedAtTime>
              <EstimatedVehicleJourney>
                <LineRef>NINOXE:Line:3:LOC</LineRef>
                <DirectionRef>unknown</DirectionRef>
                <FramedVehicleJourneyRef>
                  <DataFrameRef>2017-01-01</DataFrameRef>
                  <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
                </FramedVehicleJourneyRef>
                <PublishedLineName>Ligne 3 Metro</PublishedLineName>
                <EstimatedCalls>
                  <EstimatedCall>
                    <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
                    <Order>4</Order>
                    <StopPointName>FIRST</StopPointName>
                    <AimedArrivalTime>2017-01-01T13:00:00.000+02:00</AimedArrivalTime>
                    <ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ExpectedArrivalTime>
                  </EstimatedCall>
                  <EstimatedCall>
                    <StopPointRef>NINOXE:StopPoint:SP:25:LOC</StopPointRef>
                    <Order>5</Order>
                    <StopPointName>SECOND</StopPointName>
                    <AimedArrivalTime>2017-01-01T15:00:00.000+02:00</AimedArrivalTime>
                    <ExpectedArrivalTime>2017-01-01T15:01:00.000+02:00</ExpectedArrivalTime>
                  </EstimatedCall>
                </EstimatedCalls>
                <IsCompleteStopSequence>true</IsCompleteStopSequence>
              </EstimatedVehicleJourney>
            </EstimatedJourneyVersionFrame>
          </EstimatedTimetableDelivery>
        </ServiceDelivery>
      </Siri>
      """

  @ARA-1234
  Scenario: Do not broadcast full ETT StopVisits after a first full broadcast when editing a StopVisit with ScheduledStopVisits
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | default                        |
      | remote_code_space                  | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name            | FIRST                      |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name            | SECOND                     |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
      # "SubscriptionRef":"RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC"
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
      # "SubscriptionRef":"RELAIS:Subscription::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC"
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                        |
      | Codes[internal]          | NINOXE:VehicleJourney:201         |
      | LineId                   | 6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | Monitored                | true                              |
      | Attribute[DirectionName] | A Direction Name                  |
      | HasCompleteStopSequence  | true                              |
      # "Id":"6ba7b814-9dad-11d1-8-00c04fd430c8"
    And a ScheduledStopVisit exists with the following attributes:
      | Codes[internal]         | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder            | 4                                                      |
      | StopAreaId              | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId        | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop           | false                                                  |
      | Schedule[aimed]#Arrival | 2017-01-01T15:00:00.000Z                               |
      # "Id":"6ba7b814-9dad-11d1-9-00c04fd430c8"
    And a ScheduledStopVisit exists with the following attributes:
      | Codes[internal]         | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-3 |
      | PassageOrder            | 5                                                      |
      | StopAreaId              | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId        | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop           | false                                                  |
      | Schedule[aimed]#Arrival | 2017-01-01T15:20:00.000Z                               |
      # "Id":"6ba7b814-9dad-11d1-a-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes[internal]            | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder               | 4                                                      |
      | StopAreaId                 | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop              | true                                                   |
      | Schedule[aimed]#Arrival    | 2017-01-01T13:00:00.000+02:00                          |
      | Schedule[expected]#Arrival | 2017-01-01T13:01:00.000+02:00                          |
      # "Id":"6ba7b814-9dad-11d1-b-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes[internal]            | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-3 |
      | PassageOrder               | 5                                                      |
      | StopAreaId                 | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-8-00c04fd430c8                      |
      | VehicleAtStop              | false                                                  |
      | Schedule[aimed]#Arrival    | 2017-01-01T15:00:00.000+02:00                          |
      | Schedule[expected]#Arrival | 2017-01-01T15:01:00.000+02:00                          |
      # "Id":"6ba7b814-9dad-11d1-c-00c04fd430c8"
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-c-00c04fd430c8" is edited with the following attributes:
      | VehicleAtStop | true |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>test</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-e-00c04fd430c8</siri:ResponseMessageIdentifier>
            </ServiceDeliveryInfo>
            <Notification>
              <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
                <siri:SubscriberRef>subscriber</siri:SubscriberRef>
                <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
                <siri:Status>true</siri:Status>
                <siri:EstimatedJourneyVersionFrame>
                  <siri:RecordedAtTime>2017-01-01T12:00:50.000Z</siri:RecordedAtTime>
                  <siri:EstimatedVehicleJourney>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:DirectionRef>unknown</siri:DirectionRef>
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                     <siri:EstimatedCalls>
                      <siri:EstimatedCall>
                        <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                        <siri:Order>5</siri:Order>
                        <siri:StopPointName>SECOND</siri:StopPointName>
                        <siri:AimedArrivalTime>2017-01-01T15:00:00.000+02:00</siri:AimedArrivalTime>
                        <siri:ExpectedArrivalTime>2017-01-01T15:01:00.000+02:00</siri:ExpectedArrivalTime>
                      </siri:EstimatedCall>
                    </siri:EstimatedCalls>
                  </siri:EstimatedVehicleJourney>
                </siri:EstimatedJourneyVersionFrame>
              </siri:EstimatedTimetableDelivery>
            </Notification>
            <SiriExtension />
          </sw:NotifyEstimatedTimetable>
        </S:Body>
      </S:Envelope>
      """

  @ARA-1366
  Scenario: Handle a raw SIRI EstimatedTimetable subscription to all lines with a StopVisit having a VehicleJourneyId not existing should not broadcast the associated EstimatedCall
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:25:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | DirectionType                  | outbound                          |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
      | Attribute[VehicleMode]         | bus                               |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-2 |
      | PassageOrder                | 5                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-20-00c04fd430c8                     |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:20:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:30:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    # StopVisit with unknown VehicleJourneyId
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-7-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    When the StopVisit "6ba7b814-9dad-11d1-8-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:45:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
      <RecordedAtTime>2017-01-01T12:00:20.000Z</RecordedAtTime>
      <EstimatedVehicleJourney>
      <LineRef>NINOXE:Line:3:LOC</LineRef>
      <DirectionRef>outbound</DirectionRef>
      <FramedVehicleJourneyRef>
        <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
        <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
      </FramedVehicleJourneyRef>
      <VehicleMode>bus</VehicleMode>
      <PublishedLineName>Ligne 3 Metro</PublishedLineName>
      <DestinationName>La fin.</DestinationName>
      <OperatorRef>CdF:Company::410:LOC</OperatorRef>
      <EstimatedCalls>
        <EstimatedCall>
          <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
      <Order>4</Order>
          <StopPointName>Test</StopPointName>
          <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
          <ArrivalStatus>delayed</ArrivalStatus>
        </EstimatedCall>
      </EstimatedCalls>
      <IsCompleteStopSequence>false</IsCompleteStopSequence>
      </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  @ARA-1363
  Scenario: Manage a raw ETT notify after modification of a StopVisit using the generator setting reference_vehicle_journey_identifier
    Given a SIRI server on "http://localhost:8090"
    # Setting a Partner without default generators
    And a Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                                      | http://localhost:8090            |
      | remote_credential                               | test                             |
      | local_credential                                | NINOXE:default                   |
      | remote_code_space                               | internal                         |
      | broadcast.recorded_calls.duration               | 1h                               |
      | siri.envelope                                   | raw                              |
      | generators.reference_vehicle_journey_identifier | ch:1:ServiceJourney:87_TAC:%{id} |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
      | Number          | L3M               |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[_default]                | 6ba7b814                          |
      | Codes[external]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>6ba7b814-9dad-11d1-8-00c04fd430c8</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
       <RecordedAtTime>2017-01-01T12:00:15.000Z</RecordedAtTime>
       <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>2017-01-01</DataFrameRef>
            <DatedVehicleJourneyRef>ch:1:ServiceJourney:87_TAC:6ba7b814</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
         <PublishedLineName>Ligne 3 Metro</PublishedLineName>
         <DestinationName>La fin.</DestinationName>
         <OperatorRef>CdF:Company::410:LOC</OperatorRef>
         <EstimatedCalls>
           <EstimatedCall>
             <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
             <Order>4</Order>
             <StopPointName>Test</StopPointName>
             <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
             <ArrivalStatus>delayed</ArrivalStatus>
           </EstimatedCall>
         </EstimatedCalls>
         <IsCompleteStopSequence>false</IsCompleteStopSequence>
       </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """
    And an audit event should exist with these attributes:
      | Protocol        | siri                                    |
      | Direction       | sent                                    |
      | Status          | OK                                      |
      | Type            | NotifyEstimatedTimetable                |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"]          |
      | VehicleJourneys | ["ch:1:ServiceJourney:87_TAC:6ba7b814"] |
      | Lines           | ["NINOXE:Line:3:LOC"]                   |

  @ARA-1363
  Scenario: Manage a raw ETT notify after modification of a StopVisit using the default generator should send DatedVehicleJourneyRef according to default setting
    Given a SIRI server on "http://localhost:8090"
    # Setting a "SIRI Partner" with default generators
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | local_credential                  | NINOXE:default        |
      | remote_code_space                 | internal              |
      | broadcast.recorded_calls.duration | 1h                    |
      | siri.envelope                     | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
      | Number          | L3M               |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[_default]                | 6ba7b814                          |
      | Codes[external]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
       <RecordedAtTime>2017-01-01T12:00:15.000Z</RecordedAtTime>
       <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>RATPDev:VehicleJourney::6ba7b814:LOC</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
         <PublishedLineName>Ligne 3 Metro</PublishedLineName>
         <DestinationName>La fin.</DestinationName>
         <OperatorRef>CdF:Company::410:LOC</OperatorRef>
         <EstimatedCalls>
           <EstimatedCall>
             <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
             <Order>4</Order>
             <StopPointName>Test</StopPointName>
             <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
             <ArrivalStatus>delayed</ArrivalStatus>
           </EstimatedCall>
         </EstimatedCalls>
         <IsCompleteStopSequence>false</IsCompleteStopSequence>
       </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """
    And an audit event should exist with these attributes:
      | Protocol        | siri                                     |
      | Direction       | sent                                     |
      | Status          | OK                                       |
      | Type            | NotifyEstimatedTimetable                 |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"]           |
      | VehicleJourneys | ["RATPDev:VehicleJourney::6ba7b814:LOC"] |

  @ARA-1475
  Scenario: Manage a raw ETT notify after modification of a StopVisit using the setting broadcast.prefer_referent_stop_areas should broadcast Referent StopArea
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                           | http://localhost:8090 |
      | remote_credential                    | test                  |
      | local_credential                     | NINOXE:default        |
      | remote_code_space                    | internal              |
      | broadcast.recorded_calls.duration    | 1h                    |
      | siri.envelope                        | raw                   |
      | broadcast.prefer_referent_stop_areas | true                  |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Referent                         |
      | Codes[internal] | NINOXE:StopPoint:SP:Referent:LOC |
    And a StopArea exists with the following attributes:
      | Name            | Test                              |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC        |
      | ReferentId      | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
      | Number          | L3M               |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[_default]                | 6ba7b814                          |
      | Codes[external]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-7-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
      <ServiceDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <ProducerRef>test</ProducerRef>
      <ResponseMessageIdentifier>6ba7b814-9dad-11d1-9-00c04fd430c8</ResponseMessageIdentifier>
      <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:15.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
       <RecordedAtTime>2017-01-01T12:00:15.000Z</RecordedAtTime>
       <EstimatedVehicleJourney>
         <LineRef>NINOXE:Line:3:LOC</LineRef>
         <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>2017-01-01</DataFrameRef>
            <DatedVehicleJourneyRef>VehicleJourney:6ba7b814</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
         <PublishedLineName>Ligne 3 Metro</PublishedLineName>
         <DestinationName>La fin.</DestinationName>
         <OperatorRef>CdF:Company::410:LOC</OperatorRef>
         <EstimatedCalls>
           <EstimatedCall>
             <StopPointRef>NINOXE:StopPoint:SP:Referent:LOC</StopPointRef>
             <Order>4</Order>
             <StopPointName>Referent</StopPointName>
             <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
             <ArrivalStatus>delayed</ArrivalStatus>
           </EstimatedCall>
         </EstimatedCalls>
         <IsCompleteStopSequence>false</IsCompleteStopSequence>
       </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
      </EstimatedTimetableDelivery>
      </ServiceDelivery>
      </Siri>
      """

  @ARA-1466
  Scenario: Manage a raw ETT Notify after modification of StopVisit with status Cancelled in the future should be broadcasted in EstimatedCalls
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                        | http://localhost:8090 |
      | remote_credential                 | test                  |
      | local_credential                  | NINOXE:default        |
      | remote_code_space                 | internal              |
      | broadcast.recorded_calls.duration | 1h                    |
      | siri.envelope                     | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 1                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T12:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T12:02:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 5 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival   | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus                | cancelled                |
      | DepartureStatus              | cancelled                |
      | Schedule[expected]#Departure | 2017-01-01T15:01:11.000Z |
    And 5 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
       <ServiceDelivery>
         <ResponseTimestamp>2017-01-01T12:00:10.000Z</ResponseTimestamp>
         <ProducerRef>test</ProducerRef>
         <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</ResponseMessageIdentifier>
         <EstimatedTimetableDelivery>
           <ResponseTimestamp>2017-01-01T12:00:10.000Z</ResponseTimestamp>
           <SubscriberRef>subscriber</SubscriberRef>
           <SubscriptionRef>externalId</SubscriptionRef>
           <Status>true</Status>
           <EstimatedJourneyVersionFrame>
           <RecordedAtTime>2017-01-01T12:00:10.000Z</RecordedAtTime>
           <EstimatedVehicleJourney>
             <LineRef>NINOXE:Line:3:LOC</LineRef>
             <DirectionRef>Aller</DirectionRef>
             <FramedVehicleJourneyRef>
               <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
               <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
             </FramedVehicleJourneyRef>
             <PublishedLineName>Ligne 3 Metro</PublishedLineName>
             <DestinationName>La fin.</DestinationName>
             <OperatorRef>CdF:Company::410:LOC</OperatorRef>
             <EstimatedCalls>
               <EstimatedCall>
                 <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
                 <Order>1</Order>
                 <StopPointName>Test</StopPointName>
                 <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
                 <ArrivalStatus>cancelled</ArrivalStatus>
                 <ExpectedDepartureTime>2017-01-01T15:01:11.000Z</ExpectedDepartureTime>
                 <DepartureStatus>cancelled</DepartureStatus>
               </EstimatedCall>
             </EstimatedCalls>
             <IsCompleteStopSequence>false</IsCompleteStopSequence>
           </EstimatedVehicleJourney>
           </EstimatedJourneyVersionFrame>
         </EstimatedTimetableDelivery>
       </ServiceDelivery>
      </Siri>
      """

  @ARA-1493
  Scenario: Handle referent lines a ETT Notify after modification of a StopVisit
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | local_credential      | NINOXE:default        |
      | remote_code_space     | internal              |
      | sort_payload_for_test | true                  |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a Line exists with the following attributes:
      | Codes[external] | NINOXE:Line:4:LOC                 |
      | ReferentId      | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Name            | Ligne 3 Metro                     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-7-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
      </ServiceDeliveryInfo>
      <Notification>
      <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
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
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
            </siri:EstimatedCall>
          </siri:EstimatedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
      </siri:EstimatedTimetableDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyEstimatedTimetable>
      </S:Body>
      </S:Envelope>
      """

  @ARA-1534
  Scenario: Manage a ETT Notify after modification of a StopVisit with Vehicle occupancy
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | Aller                             |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And a Vehicle exists with the following attributes:
      | Codes[internal] | Test:Vehicle:1:LOC                |
      | Longitude       | 1.234                             |
      | Latitude        | 5.678                             |
      | Bearing         | 123                               |
      | Occupancy       | fewSeatsAvailable                 |
      | Percentage      | 15.6                              |
      | RecordedAtTime  | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime  | 2017-01-01T14:00:00.000Z          |
      | NextStopVisitId | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
      </ServiceDeliveryInfo>
      <Notification>
      <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
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
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:Occupancy>fewSeatsAvailable</siri:Occupancy>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
            </siri:EstimatedCall>
          </siri:EstimatedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
      </siri:EstimatedTimetableDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyEstimatedTimetable>
      </S:Body>
      </S:Envelope>
      """

  @ARA-1534
  Scenario: Manage a raw ETT Notify after modification of a StopVisit with Vehicle occupancy
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                           | Passage 32                        |
      | Codes[internal]                | NINOXE:VehicleJourney:201         |
      | LineId                         | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | DirectionType                  | outbound                          |
      | Attribute[OriginName]          | Le début                          |
      | Attribute[DestinationName]     | La fin.                           |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"        |
      | Attribute[VehicleMode]         | bus                               |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | false                                                  |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:00:00.000Z                               |
      | Schedule[expected]#Arrival  | 2017-01-01T15:00:00.000Z                               |
      | ArrivalStatus               | onTime                                                 |
    And a Vehicle exists with the following attributes:
      | Codes[internal] | Test:Vehicle:1:LOC                |
      | Longitude       | 1.234                             |
      | Latitude        | 5.678                             |
      | Bearing         | 123                               |
      | Occupancy       | fewSeatsAvailable                 |
      | Percentage      | 15.6                              |
      | RecordedAtTime  | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime  | 2017-01-01T14:00:00.000Z          |
      | NextStopVisitId | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
        <ServiceDelivery>
          <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
          <ProducerRef>test</ProducerRef>
          <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</ResponseMessageIdentifier>
          <EstimatedTimetableDelivery>
            <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
            <SubscriberRef>subscriber</SubscriberRef>
            <SubscriptionRef>externalId</SubscriptionRef>
            <Status>true</Status>
            <EstimatedJourneyVersionFrame>
              <RecordedAtTime>2017-01-01T12:00:20.000Z</RecordedAtTime>
              <EstimatedVehicleJourney>
                <LineRef>NINOXE:Line:3:LOC</LineRef>
                <DirectionRef>outbound</DirectionRef>
                <FramedVehicleJourneyRef>
                  <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
                  <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
                </FramedVehicleJourneyRef>
                <VehicleMode>bus</VehicleMode>
                <PublishedLineName>Ligne 3 Metro</PublishedLineName>
                <DestinationName>La fin.</DestinationName>
                <OperatorRef>CdF:Company::410:LOC</OperatorRef>
                <EstimatedCalls>
                  <EstimatedCall>
                    <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
                    <Order>4</Order>
                    <StopPointName>Test</StopPointName>
                    <Occupancy>fewSeatsAvailable</Occupancy>
                    <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
                    <ArrivalStatus>delayed</ArrivalStatus>
                  </EstimatedCall>
                </EstimatedCalls>
                <IsCompleteStopSequence>false</IsCompleteStopSequence>
              </EstimatedVehicleJourney>
            </EstimatedJourneyVersionFrame>
          </EstimatedTimetableDelivery>
        </ServiceDelivery>
      </Siri>
      """
