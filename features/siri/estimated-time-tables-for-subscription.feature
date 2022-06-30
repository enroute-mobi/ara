Feature: Support SIRI EstimatedTimeTable by subscription

  Background:
    Given a Referential "test" is created

  @ARA-1060
  Scenario: Handle a raw SIRI EstimatedTimeTable request for subscription for all existing lines in a referential having same remote_objectid_kind
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
       | siri.envelope        |  raw                  |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:A:BUS" |
      | Name      | Ligne A Bus                     |
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
    Then Subscriptions exist with the following attributes:
      | internal | NINOXE:Line:3:LOC |
      | internal | NINOXE:Line:A:BUS |

  @ARA-1060
  Scenario: Handle a raw SIRI EstimatedTimeTable request for subscription for all existing lines in a referential only with same remote_objectid_kind
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
       | siri.envelope        | raw                   |
    And a Line exists with the following attributes:
      | ObjectIDs | "another": "NINOXE:Line:3:LOC"  |
      | Name      | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:A:BUS" |
      | Name      | Ligne A Bus                     |
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
    Then Subscriptions exist with the following attributes:
      | internal | NINOXE:Line:A:BUS |
    Then No Subscriptions exist with the following attributes:
      | internal | NINOXE:Line:3:LOC |

  Scenario: 4234 - Handle a SOAP SIRI EstimatedTimeTable request for subscription
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
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
  Scenario: Handle a raw SIRI EstimatedTimeTable request for subscription
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
       | siri.envelope        | raw                   |
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

  Scenario: 4235 - Manage a ETT Notify after modification of a StopVisit
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                               | Passage 32                              |
      | ObjectIDs                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                             | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
      | Attribute[DirectionRef]            | Aller                                   |
      | Attribute[OriginName]              | Le début                                |
      | Attribute[DestinationName]         | La fin.                                 |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:00:00.000Z                                             |
      | ArrivalStatus                   | onTime                                                               |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
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
    <siri:RequestMessageRef></siri:RequestMessageRef>
  </ServiceDeliveryInfo>
  <Notification>
    <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
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
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:VehicleAtStop>false</siri:VehicleAtStop>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
            </siri:EstimatedCall>
          </siri:EstimatedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
    </siri:EstimatedTimetableDelivery>
  </Notification>
  <NotifyExtension />
</sw:NotifyEstimatedTimetable>
</S:Body>
</S:Envelope>
      """

  @ARA-1062 
  Scenario: Manage a ETT Notify after modification of a StopVisit with StopVisit departure time within the broadcast.recorded_calls.duration
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
       | remote_url                        | http://localhost:8090 |
       | remote_credential                 | test                  |
       | local_credential                  | NINOXE:default        |
       | remote_objectid_kind              | internal              |
       | broadcast.recorded_calls.duration | 1h                    |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                               | Passage 32                              |
      | ObjectIDs                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                             | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
      | Attribute[DirectionRef]            | Aller                                   |
      | Attribute[OriginName]              | Le début                                |
      | Attribute[DestinationName]         | La fin.                                 |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T12:00:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:00:00.000Z                                             |
      | ArrivalStatus                   | onTime                                                               |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
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
    <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
    <siri:RequestMessageRef></siri:RequestMessageRef>
  </ServiceDeliveryInfo>
  <Notification>
    <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:15.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:15.000Z</siri:RecordedAtTime>
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
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:VehicleAtStop>false</siri:VehicleAtStop>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
            </siri:EstimatedCall>
          </siri:EstimatedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
    </siri:EstimatedTimetableDelivery>
  </Notification>
  <NotifyExtension />
</sw:NotifyEstimatedTimetable>
</S:Body>
</S:Envelope>
      """
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival   | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus                | arrived                  |
      | DepartureStatus              | departed                 |
      | Schedule[expected]#Departure | 2017-01-01T15:01:11.000Z |
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
    <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
    <siri:RequestMessageRef></siri:RequestMessageRef>
  </ServiceDeliveryInfo>
  <Notification>
    <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
        <siri:EstimatedVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
          <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
          <siri:RecordedCalls>
            <siri:RecordedCall>
              <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
              <siri:Order>4</siri:Order>
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>arrived</siri:ArrivalStatus>
              <siri:ExpectedDepartureTime>2017-01-01T15:01:11.000Z</siri:ExpectedDepartureTime>
              <siri:DepartureStatus>departed</siri:DepartureStatus>
            </siri:RecordedCall>
          </siri:RecordedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
    </siri:EstimatedTimetableDelivery>
  </Notification>
  <NotifyExtension />
</sw:NotifyEstimatedTimetable>
</S:Body>
</S:Envelope>
      """

  @ARA-1062 
  Scenario: Manage a ETT Notify after modification of a StopVisit with StopVisit departure time oustide the broadcast.recorded_calls.duration
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
       | remote_url                        | http://localhost:8090 |
       | remote_credential                 | test                  |
       | local_credential                  | NINOXE:default        |
       | remote_objectid_kind              | internal              |
       | broadcast.recorded_calls.duration | 1h                    |
    Given a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                               | Passage 32                              |
      | ObjectIDs                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                             | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
      | Attribute[DirectionRef]            | Aller                                   |
      | Attribute[OriginName]              | Le début                                |
      | Attribute[DestinationName]         | La fin.                                 |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T01:00:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T01:10:00.000Z                                             |
      | ArrivalStatus                   | onTime                                                               |
      | Schedule[expected]#Departure    | 2017-01-01T01:20:00.000Z                                             |
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | ArrivalStatus   | arrived  |
      | DepartureStatus | departed |
    And 10 seconds have passed
    Then the SIRI server should not have received a NotifyEstimatedTimetable request

  @ARA-1025
  Scenario: Manage a raw ETT Notify after modification of a StopVisit
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
       | siri.envelope        | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                               | Passage 32                              |
      | ObjectIDs                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                             | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
      # | Attribute[DirectionRef]            | Aller                                   |
      | Attribute[OriginName]              | Le début                                |
      | Attribute[DestinationName]         | La fin.                                 |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:00:00.000Z                                             |
      | ArrivalStatus                   | onTime                                                               |
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
      <DirectionRef>unknown</DirectionRef>
      <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
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
       | remote_objectid_kind                       | internal              |
       | broadcast.no_destinationref_rewriting_from | NoRewriteOrigin       |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Origin                             | NoRewriteOrigin                         |
      | Name                               | Passage 32                              |
      | ObjectIDs                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                             | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
      | Attribute[DirectionRef]            | Aller                                   |
      | Attribute[OriginName]              | Le début                                |
      | Attribute[DestinationName]         | La fin.                                 |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:00:00.000Z                                             |
      | ArrivalStatus                   | onTime                                                               |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
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
    <siri:RequestMessageRef></siri:RequestMessageRef>
  </ServiceDeliveryInfo>
  <Notification>
    <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
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
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:VehicleAtStop>false</siri:VehicleAtStop>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
            </siri:EstimatedCall>
          </siri:EstimatedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
    </siri:EstimatedTimetableDelivery>
  </Notification>
  <NotifyExtension />
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
       | remote_objectid_kind               | internal              |
       | siri.envelope                      | raw                   |
       | broadcast.subscriptions.persistent | true                  |
    And a StopArea exists with the following attributes:
      | Name      | Test 24                                  |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Lines     | ["6ba7b814-9dad-11d1-4-00c04fd430c8"]    |
    And a StopArea exists with the following attributes:
      | Name      | Test 25                                  |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:25:LOC" |
      | Lines     | ["6ba7b814-9dad-11d1-4-00c04fd430c8"]    |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                               | Passage 32                              |
      | ObjectIDs                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                             | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
      | Attribute[DirectionRef]            | Aller                                   |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast            |
      | ExternalId        | SpecialExternalId                      |
      | SubscriberRef     | subscriber                             |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC"  |
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