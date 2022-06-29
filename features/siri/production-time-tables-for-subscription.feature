Feature: Support SIRI ProductionTimeTable by subscription

  Background:
    Given a Referential "test" is created

  Scenario: Handle a SOAP SIRI ProductionTimeTable subscription to all lines
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-production-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | ara                   |
       | local_credential     | test                  |
       | remote_objectid_kind | internal              |
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
      | Attribute[DirectionType]           | Aller                                   |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a ScheduledStopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z                                             |
    And a ScheduledStopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T17:00:00.000Z                                             |
    And a minute has passed
    And I send this SIRI request
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
  <sw:Subscribe xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri' xmlns:sws='http://wsdl.siri.org.uk/siri'>
    <SubscriptionRequestInfo>
      <siri:RequestTimestamp>2017-01-01T12:01:05.000Z</siri:RequestTimestamp>      
      <siri:RequestorRef>test</siri:RequestorRef>
      <siri:MessageIdentifier>1</siri:MessageIdentifier>
    </SubscriptionRequestInfo>
    <Request>
      <siri:ProductionTimetableSubscriptionRequest>
        <siri:SubscriberRef>test</siri:SubscriberRef>
        <siri:SubscriptionIdentifier>1</siri:SubscriptionIdentifier>
        <siri:InitialTerminationTime>2017-01-03T12:01:05.000Z</siri:InitialTerminationTime>
        <siri:ProductionTimetableRequest>
          <siri:RequestTimestamp>2017-01-01T12:01:05.000Z</siri:RequestTimestamp>
        </siri:ProductionTimetableRequest>
      </siri:ProductionTimetableSubscriptionRequest>
    </Request>
    <RequestExtension/>
  </sw:Subscribe>
</S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
<sw:SubscribeResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
    <SubscriptionAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
        <siri:Address></siri:Address>
        <siri:ResponderRef>ara</siri:ResponderRef>
        <siri:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="siri:MessageRefStructure">1</siri:RequestMessageRef>
    </SubscriptionAnswerInfo>
    <Answer>
        <siri:ResponseStatus>
            <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
            <siri:RequestMessageRef></siri:RequestMessageRef>
            <siri:SubscriberRef>test</siri:SubscriberRef>
            <siri:SubscriptionRef>1</siri:SubscriptionRef>
            <siri:Status>true</siri:Status>
            <siri:ValidUntil>2017-01-03T12:01:05.000Z</siri:ValidUntil>
        </siri:ResponseStatus>
        <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
    </Answer>
<AnswerExtension />
</sw:SubscribeResponse>
</S:Body>
</S:Envelope>
      """
    And 2 minutes have passed
    Then the SIRI server should receive this response
      """
<?xml version='1.0' encoding='UTF-8'?>
<S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
  <S:Body>
    <sw:NotifyProductionTimetable xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:03:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>ara</siri:ProducerRef>
      </ServiceDeliveryInfo>
      <Notification>
        <siri:ProductionTimetableDelivery version='2.0:FR-IDF-2.4'>
          <siri:ResponseTimestamp>2017-01-01T12:03:00.000Z</siri:ResponseTimestamp>
          <siri:SubscriptionRef>1</siri:SubscriptionRef>
          <siri:Status>true</siri:Status>
          <siri:DatedTimetableVersionFrame>
            <siri:RecordedAtTime>2017-01-01T12:03:00.000Z</siri:RecordedAtTime>
            <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
            <siri:DirectionRef>Aller</siri:DirectionRef>
            <siri:FirstOrLastJourney>unspecified</siri:FirstOrLastJourney>
            <siri:DatedVehicleJourney>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>2017-01-01</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:FirstOrLastJourney>unspecified</siri:FirstOrLastJourney>
              <siri:DatedCalls>
                <siri:DatedCall>
                  <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                  <siri:StopPointName>Test 24</siri:StopPointName>
                  <siri:AimedArrivalTime>2017-01-01T15:00:00.000Z</siri:AimedArrivalTime>
                </siri:DatedCall>
                <siri:DatedCall>
                  <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                  <siri:StopPointName>Test 25</siri:StopPointName>
                  <siri:AimedArrivalTime>2017-01-01T17:00:00.000Z</siri:AimedArrivalTime>
                </siri:DatedCall>
              </siri:DatedCalls>
            </siri:DatedVehicleJourney>
          </siri:DatedTimetableVersionFrame>
        </siri:ProductionTimetableDelivery>
      </Notification>
      <NotifyExtension/>
    </sw:NotifyProductionTimetable>
  </S:Body>
</S:Envelope>
      """

  Scenario: Handle a raw SIRI ProductionTimeTable subscription to all lines
    Given a raw SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-production-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | ara                   |
       | local_credential     | test                  |
       | remote_objectid_kind | internal              |
       | siri.envelope        | raw                   |
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
      | Attribute[DirectionType]           | Retour                                  |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a ScheduledStopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z                                             |
    And a ScheduledStopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T17:00:00.000Z                                             |
    And a minute has passed
    And I send this SIRI request
      """
<?xml version="1.0" encoding="utf-8"?>
<Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
   <SubscriptionRequest>
      <RequestTimestamp>2022-02-09T02:15:23.690717Z</RequestTimestamp>
      <RequestorRef>test</RequestorRef>
      <ProductionTimetableSubscriptionRequest>
         <SubscriptionIdentifier>1</SubscriptionIdentifier>
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
      <SubscriptionRef>1</SubscriptionRef>
      <Status>true</Status>
      <ValidUntil>2022-02-10T02:50:00.000Z</ValidUntil>
    </ResponseStatus>
    <ServiceStartedTime>2017-01-01T12:00:00.000Z</ServiceStartedTime>
  </SubscriptionResponse>
</Siri>
      """
    And 2 minutes have passed
    Then the SIRI server should receive this response
      """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
  <ServiceDelivery>
    <ResponseTimestamp>2017-01-01T12:03:00.000Z</ResponseTimestamp>
    <ProducerRef>ara</ProducerRef>
    <ProductionTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:03:00.000Z</ResponseTimestamp>
      <SubscriptionRef>1</SubscriptionRef>
      <Status>true</Status>
      <DatedTimetableVersionFrame>
        <RecordedAtTime>2017-01-01T12:03:00.000Z</RecordedAtTime>
        <LineRef>NINOXE:Line:3:LOC</LineRef>
        <DirectionRef>ch:1:Direction:R</DirectionRef>
        <FirstOrLastJourney>unspecified</FirstOrLastJourney>
        <DatedVehicleJourney>
          <FramedVehicleJourneyRef>
            <DataFrameRef>2017-01-01</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
          <PublishedLineName>Ligne 3 Metro</PublishedLineName>
          <OperatorRef>CdF:Company::410:LOC</OperatorRef>
          <FirstOrLastJourney>unspecified</FirstOrLastJourney>
          <DatedCalls>
            <DatedCall>
              <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
              <StopPointName>Test 24</StopPointName>
              <AimedArrivalTime>2017-01-01T15:00:00.000Z</AimedArrivalTime>
            </DatedCall>
            <DatedCall>
              <StopPointRef>NINOXE:StopPoint:SP:25:LOC</StopPointRef>
              <StopPointName>Test 25</StopPointName>
              <AimedArrivalTime>2017-01-01T17:00:00.000Z</AimedArrivalTime>
            </DatedCall>
          </DatedCalls>
        </DatedVehicleJourney>
      </DatedTimetableVersionFrame>
    </ProductionTimetableDelivery>
  </ServiceDelivery>
</Siri>
      """

  Scenario: Handle a raw SIRI ProductionTimeTable subscription on an unknown line
    Given a raw SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-production-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | ara                   |
       | local_credential     | test                  |
       | remote_objectid_kind | internal              |
       | siri.envelope        | raw                   |
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
      | Attribute[DirectionType]           | Aller                                   |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a minute has passed
    And I send this SIRI request
      """
<?xml version="1.0" encoding="utf-8"?>
<Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
   <SubscriptionRequest>
      <RequestTimestamp>2022-02-09T02:15:23.690717Z</RequestTimestamp>
      <RequestorRef>test</RequestorRef>
      <ProductionTimetableSubscriptionRequest>
         <SubscriptionIdentifier>1</SubscriptionIdentifier>
         <InitialTerminationTime>2022-02-10T02:50:00Z</InitialTerminationTime>
         <ProductionTimetableRequest>
            <RequestTimestamp>2022-02-09T02:15:23.690717Z</RequestTimestamp>
            <ValidityPeriod>
               <StartTime>2022-02-09T03:30:00Z</StartTime>
               <EndTime>2022-02-10T04:30:00Z</EndTime>
            </ValidityPeriod>
            <Lines>
              <LineDirection>
                <LineRef>NINOXE:Line:UNKNOWN:LOC</LineRef>
              </LineDirection>
            </Lines>
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
      <SubscriptionRef>1</SubscriptionRef>
      <Status>false</Status>
      <ErrorCondition>
        <InvalidDataReferencesError>
          <ErrorText>Unknown Line(s) NINOXE:Line:UNKNOWN:LOC</ErrorText>
        </InvalidDataReferencesError>
      </ErrorCondition>
    </ResponseStatus>
    <ServiceStartedTime>2017-01-01T12:00:00.000Z</ServiceStartedTime>
  </SubscriptionResponse>
</Siri>
      """

  @ARA-1107
  Scenario: Handle a raw SIRI ProductionTimeTable subscription to all lines with StopArea having a Parent with Partner ObjectiId
    Given a raw SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-production-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | ara                   |
       | local_credential     | test                  |
       | remote_objectid_kind | didok                 |
       | siri.envelope        | raw                   |
    And a StopArea exists with the following attributes:
    # "6ba7b814-9dad-11d1-2-00c04fd430c8"
      | Name      | Parent                                             |
      | ObjectIDs | "didok": "fr:1:StopPlace:OURA2:StopArea:log351672" |
    And a StopArea exists with the following attributes:
    # "6ba7b814-9dad-11d1-3-00c04fd430c8"
      | Name      | Child                                 |
      | ObjectIDs | "internal": "vlgabon1"                |
      | Lines     | ["6ba7b814-9dad-11d1-4-00c04fd430c8"] |
      | ParentId  | 6ba7b814-9dad-11d1-2-00c04fd430c8     |
    And a Line exists with the following attributes:
    # "6ba7b814-9dad-11d1-4-00c04fd430c8"
      | ObjectIDs | "didok": "NINOXE:Line:3:LOC"    |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
    # "6ba7b814-9dad-11d1-5-00c04fd430c8"
      | Name                               | Passage 32                              |
      | ObjectIDs                          | "didok": "NINOXE:VehicleJourney:201"    |
      | LineId                             | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
      | Attribute[DirectionType]           | Aller                                   |
      | Reference[DestinationRef]#ObjectId | "internal": "ThisIsTheEnd"              |
    And a ScheduledStopVisit exists with the following attributes:
    # "6ba7b814-9dad-11d1-6-00c04fd430c8"
      | PassageOrder                    | 4                                  |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8  |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8  |
      | VehicleAtStop                   | false                              |
      | Reference[OperatorRef]#ObjectId | "didok": "CdF:Company::410:LOC"    |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z           |
    And a minute has passed
    And I send this SIRI request
      """
<?xml version="1.0" encoding="utf-8"?>
<Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
   <SubscriptionRequest>
      <RequestTimestamp>2022-02-09T02:15:23.690717Z</RequestTimestamp>
      <RequestorRef>test</RequestorRef>
      <ProductionTimetableSubscriptionRequest>
         <SubscriptionIdentifier>1</SubscriptionIdentifier>
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
      <SubscriptionRef>1</SubscriptionRef>
      <Status>true</Status>
      <ValidUntil>2022-02-10T02:50:00.000Z</ValidUntil>
    </ResponseStatus>
    <ServiceStartedTime>2017-01-01T12:00:00.000Z</ServiceStartedTime>
  </SubscriptionResponse>
</Siri>
      """
    And 2 minutes have passed
    Then the SIRI server should receive this response
      """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
  <ServiceDelivery>
    <ResponseTimestamp>2017-01-01T12:03:00.000Z</ResponseTimestamp>
    <ProducerRef>ara</ProducerRef>
    <ProductionTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:03:00.000Z</ResponseTimestamp>
      <SubscriptionRef>1</SubscriptionRef>
      <Status>true</Status>
      <DatedTimetableVersionFrame>
        <RecordedAtTime>2017-01-01T12:03:00.000Z</RecordedAtTime>
        <LineRef>NINOXE:Line:3:LOC</LineRef>
        <DirectionRef>ch:1:Direction:H</DirectionRef>
        <FirstOrLastJourney>unspecified</FirstOrLastJourney>
        <DatedVehicleJourney>
          <FramedVehicleJourneyRef>
            <DataFrameRef>2017-01-01</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
          <PublishedLineName>Ligne 3 Metro</PublishedLineName>
          <OperatorRef>CdF:Company::410:LOC</OperatorRef>
          <FirstOrLastJourney>unspecified</FirstOrLastJourney>
          <DatedCalls>
            <DatedCall>
              <StopPointRef>fr:1:StopPlace:OURA2:StopArea:log351672</StopPointRef>
              <StopPointName>Parent</StopPointName>
              <AimedArrivalTime>2017-01-01T15:00:00.000Z</AimedArrivalTime>
            </DatedCall>
          </DatedCalls>
        </DatedVehicleJourney>
      </DatedTimetableVersionFrame>
    </ProductionTimetableDelivery>
  </ServiceDelivery>
</Siri>
      """

  Scenario: Handle a raw SIRI ProductionTimeTable subscription on a specific line
    Given a raw SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-production-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | ara                   |
       | local_credential     | test                  |
       | remote_objectid_kind | internal              |
       | siri.envelope        | raw                   |
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
      | Attribute[DirectionType]           | Aller                                   |
      | Reference[DestinationRef]#ObjectId | "external": "ThisIsTheEnd"              |
    And a minute has passed
    And I send this SIRI request
      """
<?xml version="1.0" encoding="utf-8"?>
<Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
   <SubscriptionRequest>
      <RequestTimestamp>2022-02-09T02:15:23.690717Z</RequestTimestamp>
      <RequestorRef>test</RequestorRef>
      <ProductionTimetableSubscriptionRequest>
         <SubscriptionIdentifier>1</SubscriptionIdentifier>
         <InitialTerminationTime>2022-02-10T02:50:00Z</InitialTerminationTime>
         <ProductionTimetableRequest>
            <RequestTimestamp>2022-02-09T02:15:23.690717Z</RequestTimestamp>
            <ValidityPeriod>
               <StartTime>2022-02-09T03:30:00Z</StartTime>
               <EndTime>2022-02-10T04:30:00Z</EndTime>
            </ValidityPeriod>
            <Lines>
              <LineDirection>
                <LineRef>NINOXE:Line:3:LOC</LineRef>
              </LineDirection>
            </Lines>
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
      <SubscriptionRef>1</SubscriptionRef>
      <Status>true</Status>
      <ValidUntil>2022-02-10T02:50:00.000Z</ValidUntil>
    </ResponseStatus>
    <ServiceStartedTime>2017-01-01T12:00:00.000Z</ServiceStartedTime>
  </SubscriptionResponse>
</Siri>
      """
