Feature: Support SIRI EstimatedTimetable

  Background:
      Given a Referential "test" is created

  @ARA-1306
  Scenario: EstimatedTimetable subscription collect should send EstimatedTimetableSubscriptionRequest to partner
   Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url            | http://localhost:8090                |
      | remote_credential     | test                                 |
      | remote_code_space     | internal                             |
      | collect.include_lines | RLA_Bus:Line::05:LOC,RLA_TRAM::A:LOC |
      | local_credential      | ara                                  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name            | Test 1               |
      | Codes[internal] | RLA_Bus:Line::05:LOC |
    And a Line exists with the following attributes:
      | Name            | Test 2          |
      | Codes[internal] | RLA_TRAM::A:LOC |
   And a minute has passed
   And 20 seconds have passed
   Then the SIRI server should have received a EstimatedTimetableSubscriptionRequest request with:
     | //siri:LineRef | ["RLA_Bus:Line::05:LOC", "RLA_TRAM::A:LOC"] |

  @ARA-1306
  Scenario: EstimatedTimetable subscription collect and partner CheckStatus is unavailable should not send EstimatedTimetableSubscriptionRequest to partner
   Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
      | local_credential      | ara                   |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name            | Test 1               |
      | Codes[internal] | RLA_Bus:Line::05:LOC |
   And a minute has passed
   And 10 seconds have passed
   Then the SIRI server should not have received a EstimatedTimetableSubscriptionRequest request

  @ARA-1306
  Scenario: EstimatedTimetable subscription collect and partner CheckStatus is unavailable should send EstimatedTimetableSubscriptionRequest to partner whith setting collect.subscriptions.persistent
   Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url                       | http://localhost:8090 |
      | remote_credential                | test                  |
      | remote_code_space                | internal              |
      | collect.include_lines            | RLA_Bus:Line::05:LOC  |
      | local_credential                 | ara                   |
      | collect.subscriptions.persistent | true                  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name            | Test 1               |
      | Codes[internal] | RLA_Bus:Line::05:LOC |
   And a minute has passed
   And 30 seconds have passed
   Then the SIRI server should have received a EstimatedTimetableSubscriptionRequest request with:
      | //siri:LineRef | RLA_Bus:Line::05:LOC |

  @ARA-1306
  Scenario: EstimatedTimetable subscription collect and partner CheckStatus is unavailable should send EstimatedTimetableSubscriptionRequest to partner whith setting collect.persistent
   Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
      | local_credential      | ara                   |
      | collect.persistent    | true                  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name            | Test 1               |
      | Codes[internal] | RLA_Bus:Line::05:LOC |
   And a minute has passed
   And 10 seconds have passed
   Then the SIRI server should have received a EstimatedTimetableSubscriptionRequest request with:
      | //siri:LineRef | RLA_Bus:Line::05:LOC |

  @ARA-1152
  Scenario: Create EstimatedTimetable subscription collect
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
        <ns5:RequestMessageRef>response</ns5:RequestMessageRef>
      </SubscriptionAnswerInfo>
      <Answer
        xmlns:ns2="http://www.ifopt.org.uk/acsb"
        xmlns:ns3="http://www.ifopt.org.uk/ifopt"
        xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
        xmlns:ns5="http://www.siri.org.uk/siri"
        xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseStatus>
            <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
            <ns5:RequestMessageRef>{LastRequestMessageRef}</ns5:RequestMessageRef>
            <ns5:SubscriberRef>test</ns5:SubscriberRef>
            <ns5:SubscriptionRef>6ba7b814-9dad-11d1-5-00c04fd430c8</ns5:SubscriptionRef>
            <ns5:Status>true</ns5:Status>
            <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
        </ns5:ResponseStatus>
        <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
      </Answer>
      <AnswerExtension/>
    </ns1:SubscribeResponse>
  </S:Body>
  </S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name            | Test     |
      | Codes[internal] | testLine |
    And 10 seconds have passed
    Then one Subscription exists with the following attributes:
      | Kind | EstimatedTimetableCollect |

  @ARA-1179
  Scenario: Check EstimatedTimetable subscription collect payload
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | local_url             | http://test           |
      | remote_credential     | test                  |
      | local_credential      | NINOXE:default        |
      | remote_code_space     | internal              |
      | sort_payload_for_test | true                  |
    And a Line exists with the following attributes:
      | Name            | Test1     |
      | Codes[internal] | testLine1 |
    And a Line exists with the following attributes:
      | Name            | Test2     |
      | Codes[internal] | testLine2 |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect     |
      | ExternalId        | ExternalId                    |
      | ReferenceArray[0] | Line, "internal": "testLine1" |
      | ReferenceArray[1] | Line, "internal": "testLine2" |
    And 5 seconds have passed
    Then the SIRI server should receive this response
    """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <SubscriptionRequestInfo>
  <siri:RequestTimestamp>2017-01-01T12:00:05.000Z</siri:RequestTimestamp>
  <siri:RequestorRef>test</siri:RequestorRef>
  <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
        <siri:ConsumerAddress>http://test</siri:ConsumerAddress>
      </SubscriptionRequestInfo>
      <Request>
  <siri:EstimatedTimetableSubscriptionRequest>
  <siri:SubscriberRef>test</siri:SubscriberRef>
  <siri:SubscriptionIdentifier>6ba7b814-9dad-11d1-4-00c04fd430c8</siri:SubscriptionIdentifier>
  <siri:InitialTerminationTime>2017-01-03T12:00:05.000Z</siri:InitialTerminationTime>
  <siri:EstimatedTimetableRequest version="2.0:FR-IDF-2.4">
  <siri:RequestTimestamp>2017-01-01T12:00:05.000Z</siri:RequestTimestamp>
  <siri:MessageIdentifier>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:MessageIdentifier>
  <siri:Lines>
  <siri:LineDirection>
    <siri:LineRef>testLine1</siri:LineRef>
  </siri:LineDirection>
  <siri:LineDirection>
    <siri:LineRef>testLine2</siri:LineRef>
  </siri:LineDirection>
  </siri:Lines>
  </siri:EstimatedTimetableRequest>
  <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
  <siri:ChangeBeforeUpdates>PT1M</siri:ChangeBeforeUpdates>
  </siri:EstimatedTimetableSubscriptionRequest>
      </Request>
      <RequestExtension />
    </ws:Subscribe>
  </S:Body>
</S:Envelope>
    """

  @ARA-1152 @ARA-1310 @ARA-1825
  Scenario: Create ara models after a EstimatedTimetableNotify in a subscription
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
        <ns5:RequestMessageRef>Subscription:Test:0</ns5:RequestMessageRef>
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
      <AnswerExtension/>
    </ns1:SubscribeResponse>
  </S:Body>
  </S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name            | Test              |
      | Codes[internal] | NINOXE:Line:3:LOC |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect             |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a minute has passed
    When I send this SIRI request
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
<sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
  <ServiceDeliveryInfo>
    <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
    <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
    <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
    <siri:RequestMessageRef></siri:RequestMessageRef>
  </ServiceDeliveryInfo>
  <Notification>
    <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
        <siri:EstimatedVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
          <siri:Cancellation>true</siri:Cancellation>
          <siri:DestinationRef>ThisIsTheEnd</siri:DestinationRef>
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
  <SiriExtension/>
</sw:NotifyEstimatedTimetable>
</S:Body>
</S:Envelope>
      """
    And 30 seconds have passed
    Then one VehicleJourney has the following attributes:
      | Codes[internal] | NINOXE:VehicleJourney:201         |
      | LineId          | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | DirectionType   | Aller                             |
      | Id              | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
      | Cancellation    | true                              |
    And one StopArea has the following attributes:
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC        |
      | Name            | Test                              |
      | Id              | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
    And one StopVisit has the following attributes:
      | Codes[internal]            | NINOXE:VehicleJourney:201-4       |
      | PassageOrder               | 4                                 |
      | VehicleAtStop              | false                             |
      | ArrivalStatus              | delayed                           |
      | DepartureStatus            | nil                               |
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01Z              |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
      | StopAreaId                 | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
    And an audit event should exist with these attributes:
      | Protocol        | siri                           |
      | Direction       | received                       |
      | Status          | OK                             |
      | Type            | NotifyEstimatedTimetable       |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201"]  |
      | Lines           | ["NINOXE:Line:3:LOC"]          |

  @ARA-1152 @ARA-1310
  Scenario: Update a StopVisit after an EstimatedTimetableNotify in a subscription
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
        <ns5:RequestMessageRef>Subscription:Test:0</ns5:RequestMessageRef>
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
      <AnswerExtension/>
    </ns1:SubscribeResponse>
  </S:Body>
  </S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name            | Test              |
      | Codes[internal] | NINOXE:Line:3:LOC |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                        |
      | Codes[other]             | Test:VehicleJourney:201:LOC       |
      | LineId                   | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored                | true                              |
      | Attributes[DirectionName] | Direction Name                    |
    And a StopArea exists with the following attributes:
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Name            | Test                       |
    And a StopVisit exists with the following attributes:
      | Codes[internal]            | NINOXE:VehicleJourney:201-4       |
      | PassageOrder               | 4                                 |
      | VehicleAtStop              | false                             |
      | ArrivalStatus              | onTime                            |
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01Z              |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | StopAreaId                 | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect             |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a minute has passed
    When I send this SIRI request
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
<sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
  <ServiceDeliveryInfo>
    <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
    <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
    <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
    <siri:RequestMessageRef></siri:RequestMessageRef>
  </ServiceDeliveryInfo>
  <Notification>
    <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:SubscriptionRef>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
        <siri:EstimatedVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
          <siri:DestinationRef>ThisIsTheEnd</siri:DestinationRef>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:EstimatedCalls>
            <siri:EstimatedCall>
              <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
              <siri:Order>4</siri:Order>
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:ExpectedArrivalTime>2017-01-01T15:10:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
            </siri:EstimatedCall>
          </siri:EstimatedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
    </siri:EstimatedTimetableDelivery>
  </Notification>
  <SiriExtension/>
</sw:NotifyEstimatedTimetable>
</S:Body>
</S:Envelope>
      """
    Then the StopVisit "internal:NINOXE:VehicleJourney:201-4" has the following attributes:
      | Codes[internal]            | NINOXE:VehicleJourney:201-4 |
      | PassageOrder               | 4                           |
      | VehicleAtStop              | false                       |
      | ArrivalStatus              | delayed                     |
      | Schedule[expected]#Arrival | 2017-01-01T15:10:01Z        |

  @ARA-1411
  Scenario: RAW EstimatedTimetable subscription collect should send EstimatedTimetableSubscriptionRequest to partner
   Given a raw SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
      | local_credential      | ara                   |
      | siri.envelope         | raw                   |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name            | Test 1               |
      | Codes[internal] | RLA_Bus:Line::05:LOC |
   And a minute has passed
   And 20 seconds have passed
   Then the SIRI server should have received a raw EstimatedTimetableSubscriptionRequest request with:
     | //siri:LineRef | RLA_Bus:Line::05:LOC |

  @ARA-1411
  Scenario: Create ara models after a RAW EstimatedTimetableDelivery in a subscription
    Given a raw SIRI server waits Subscribe request on "http://localhost:8090" to respond with
      """
<?xml version="1.0" encoding="utf-8"?>
<Siri xmlns="http://www.siri.org.uk/siri" version="2.0">
<SubscriptionResponse>
        <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
        <ResponderRef>NINOXE:default</ResponderRef>
        <ResponseStatus>
            <ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ResponseTimestamp>
            <RequestMessageRef>{LastRequestMessageRef}</RequestMessageRef>
            <SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</SubscriptionRef>
            <Status>true</Status>
            <ValidUntil>2016-09-22T08:01:20.227+02:00</ValidUntil>
        </ResponseStatus>
        <ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ServiceStartedTime>
</SubscriptionResponse>
</Siri>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name            | Test              |
      | Codes[internal] | NINOXE:Line:3:LOC |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect             |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a minute has passed
    Then one Subscription exists with the following attributes:
      | Kind                      | EstimatedTimetableCollect |
      | Resources[0]/SubscribedAt | > 2017-01-01T12:01:00Z    |
    When I send this SIRI request
      """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns='http://www.siri.org.uk/siri'>
 <ServiceDelivery>
    <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
    <ProducerRef>NINOXE:default</ProducerRef>
    <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</ResponseMessageIdentifier>
    <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</SubscriptionRef>
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
          <DestinationRef>ThisIsTheEnd</DestinationRef>
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
    And 30 seconds have passed
    Then one VehicleJourney has the following attributes:
      | Codes[internal] | NINOXE:VehicleJourney:201         |
      | LineId          | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | DirectionType   | Aller                             |
      | Id              | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
    And one StopArea has the following attributes:
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC        |
      | Name            | Test                              |
      | Id              | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
    And one StopVisit has the following attributes:
      | Codes[internal]            | NINOXE:VehicleJourney:201-4       |
      | PassageOrder               | 4                                 |
      | ArrivalStatus              | delayed                           |
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01Z              |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
      | StopAreaId                 | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
    And an audit event should exist with these attributes:
      | Protocol        | siri                           |
      | Direction       | received                       |
      | Status          | OK                             |
      | Type            | NotifyEstimatedTimetable       |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201"]  |
      | Lines           | ["NINOXE:Line:3:LOC"]          |

  @wip @ARA-1465
  Scenario: Create ara models after an EstimatedTimetableDelivery with multiple Estimated VehicleJourneys in a subscription
    Given a raw SIRI server waits Subscribe request on "http://localhost:8090" to respond with
      """
<?xml version="1.0" encoding="utf-8"?>
<Siri xmlns="http://www.siri.org.uk/siri" version="2.0">
<SubscriptionResponse>
        <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
        <ResponderRef>NINOXE:default</ResponderRef>
        <ResponseStatus>
            <ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ResponseTimestamp>
            <RequestMessageRef>{LastRequestMessageRef}</RequestMessageRef>
            <SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</SubscriptionRef>
            <Status>true</Status>
            <ValidUntil>2016-09-22T08:01:20.227+02:00</ValidUntil>
        </ResponseStatus>
        <ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ServiceStartedTime>
</SubscriptionResponse>
</Siri>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name            | Test              |
      | Codes[internal] | NINOXE:Line:3:LOC |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect             |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a minute has passed
    Then one Subscription exists with the following attributes:
      | Kind                      | EstimatedTimetableCollect |
      | Resources[0]/SubscribedAt | > 2017-01-01T12:01:00Z    |
    When I send this SIRI request
      """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns='http://www.siri.org.uk/siri'>
 <ServiceDelivery>
    <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
    <ProducerRef>NINOXE:default</ProducerRef>
    <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</ResponseMessageIdentifier>
    <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</SubscriptionRef>
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
          <DestinationRef>ThisIsTheEnd</DestinationRef>
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
        <EstimatedVehicleJourney>
          <LineRef>NINOXE:Line:4:LOC</LineRef>
          <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
          <DestinationRef>ThisIsTheEnd</DestinationRef>
          <OperatorRef>CdF:Company::410:LOC</OperatorRef>
          <EstimatedCalls>
            <EstimatedCall>
              <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
              <Order>3</Order>
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
    And 30 seconds have passed
    Then one Line has the following attributes:
      | Codes[internal] | NINOXE:Line:4:LOC                 |
      | Id              | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
    Then one VehicleJourney has the following attributes:
      | Codes[internal] | NINOXE:VehicleJourney:201         |
      | LineId          | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | DirectionType   | Aller                             |
      | Id              | 6ba7b814-9dad-11d1-a-00c04fd430c8 |
    And one VehicleJourney has the following attributes:
      | Codes[internal] | NINOXE:VehicleJourney:202         |
      | LineId          | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
      | DirectionType   | Aller                             |
      | Id              | 6ba7b814-9dad-11d1-b-00c04fd430c8 |
    And one StopArea has the following attributes:
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC        |
      | Name            | Test                              |
      | Id              | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
    And one StopVisit has the following attributes:
      | Codes[internal]            | NINOXE:VehicleJourney:201-4       |
      | PassageOrder               | 4                                 |
      | ArrivalStatus              | delayed                           |
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01Z              |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-a-00c04fd430c8 |
      | StopAreaId                 | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
    And one StopVisit has the following attributes:
      | Codes[internal]            | NINOXE:VehicleJourney:202-3       |
      | PassageOrder               | 3                                 |
      | ArrivalStatus              | delayed                           |
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01Z              |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-b-00c04fd430c8 |
      | StopAreaId                 | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
    And an audit event should exist with these attributes:
      | Protocol        | siri                                                      |
      | Direction       | received                                                  |
      | Status          | OK                                                        |
      | Type            | NotifyEstimatedTimetable                                  |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"]                            |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201","NINOXE:VehicleJourney:202"] |
      | Lines           | ["NINOXE:Line:3:LOC","NINOXE:Line:4:LOC"]                 |
