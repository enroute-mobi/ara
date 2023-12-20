Feature: Audit API exchanges

  Background:
    Given a Referential "test" is created

  @ARA-1096
  Scenario: Audit a Partner status changed to up
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      """
    And a Partner "test" exists with connectors [siri-check-status-client] and the following settings:
      | remote_url            | http://localhost:8090      |
      | remote_credential     | test                       |
      | remote_code_space  | internal                   |
    When a minute has passed
    Then an audit event should exist with these attributes:
      | NewStatus   | up            |
      | PartnerUUID | /{test-uuid}/ |
      | Slug        | test          |

  Scenario: Audit a received SIRI CheckStatus Request
    Given a Partner "test" exists with connectors [siri-check-status-server] and the following settings:
      | local_credential | test |
    When I send this SIRI request to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test</siri:RequestorRef>
        <siri:MessageIdentifier>enRoute:Message::test</siri:MessageIdentifier>
      </Request>
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then an audit event should exist with these attributes:
      | Type               | CheckStatusRequest    |
      | Direction          | received              |
      | Protocol           | siri                  |
      | Partner            | test                  |
      | Status             | OK                    |
      | RequestIdentifier  | enRoute:Message::test |
      | ResponseIdentifier | /{test-uuid}/ |
      | ProcessingTime     | 0                     |

  Scenario: Not audit SIRI CheckStatus Request for unknown partner
    When I send this SIRI request to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>dummy</siri:RequestorRef>
        <siri:MessageIdentifier>enRoute:Message::test</siri:MessageIdentifier>
      </Request>
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then an audit event should not exist with these attributes:
      | Type               | CheckStatusRequest    |
    And an audit event should exist with these attributes:
      | Protocol           | siri                                            |
      | Direction          | received                                        |
      | Status             | Error                                           |
      | ErrorDetails       | UnknownCredential: RequestorRef Unknown 'dummy' |

  Scenario: Audit a sent SIRI CheckStatus Request
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      """
    And a Partner "test" exists with connectors [siri-check-status-client] and the following settings:
      | remote_url                 | http://localhost:8090      |
      | remote_credential          | test                       |
      | remote_code_space       | internal                   |
    When a minute has passed
    Then an audit event should exist with these attributes:
      | Type               | CheckStatusRequest    |
      | Protocol           | siri                  |
      | Direction          | sent                  |
      | Status             | OK                    |
      | Partner            | test                  |
      | RequestIdentifier  | /{test-uuid}/         |
      | ResponseIdentifier | c464f588-5128-46c8-ac3f-8b8a465692ab |
      | ProcessingTime     | 0                     |

  Scenario: Audit a StopMonitoring Subscription request
    Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
      """
  <?xml version='1.0' encoding='utf-8'?>
  <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://wsdl.siri.org.uk">
  <S:Body>
    <ns1:SubscribeResponse xmlns:ns5="http://www.siri.org.uk/siri">
      <SubscriptionAnswerInfo>
        <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
        <ns5:ResponderRef>remote_credential</ns5:ResponderRef>
        <ns5:RequestMessageRef>Subscription:Test:0</ns5:RequestMessageRef>
        <ns5:ResponseMessageIdentifier>c464f588-5128-46c8-ac3f-8b8a465692ab</ns5:ResponseMessageIdentifier>
      </SubscriptionAnswerInfo>
      <Answer>
        <ns5:ResponseStatus>
            <ns5:ResponseTimestamp>2017-01-01T12:00:00+01:00</ns5:ResponseTimestamp>
            <ns5:RequestMessageRef>{LastRequestMessageRef}</ns5:RequestMessageRef>
            <ns5:SubscriberRef>SubscriberRef</ns5:SubscriberRef>
            <ns5:SubscriptionRef>SubscriptionIdentifier</ns5:SubscriptionRef>
            <ns5:Status>true</ns5:Status>
            <ns5:ValidUntil>2017-01-02T12:00:00+01:00</ns5:ValidUntil>
        </ns5:ResponseStatus>
      </Answer>
    </ns1:SubscribeResponse>
  </S:Body>
  </S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | remote_credential              |
      | local_credential                   | local_credential               |
      | remote_code_space               | internal                       |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name      | Test                                      |
      | Codes | "internal": "enRoute:StopPoint:SP:24:LOC" |
    When I wait that a Subscription has been created with the following attributes:
      | Kind      | StopMonitoringCollect |
    Then an audit event should exist with these attributes:
      | Type                    | StopMonitoringSubscriptionRequest     |
      | Protocol                | siri                                  |
      | Direction               | sent                                  |
      | Status                  | OK                                    |
      | Partner                 | test                                  |
      | RequestIdentifier       | /{test-uuid}/                         |
      | ResponseIdentifier      | c464f588-5128-46c8-ac3f-8b8a465692ab  |
      | ProcessingTime          | 0                                     |
      | SubscriptionIdentifiers | ["6ba7b814-9dad-11d1-4-00c04fd430c8"] |
      | StopAreas               | ["enRoute:StopPoint:SP:24:LOC"]       |

  @ARA-888
  Scenario: Audit a StopMonitoring Subscription request
    Given a SIRI server waits DeleteSubscription request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:DeleteSubscriptionResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
      <DeleteSubscriptionAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:40.000Z</siri:ResponseTimestamp>
        <siri:ResponderRef>remote_credential</siri:ResponderRef>
        <siri:RequestMessageRef>MessageIdentifier</siri:RequestMessageRef>
        <siri:ResponseMessageIdentifier>c464f588-5128-46c8-ac3f-8b8a465692ab</siri:ResponseMessageIdentifier>
      </DeleteSubscriptionAnswerInfo>
      <Answer>
        <siri:ResponseTimestamp>2017-01-01T12:00:40.000Z</siri:ResponseTimestamp>
        <siri:ResponderRef>remote_credential</siri:ResponderRef>
        <siri:RequestMessageRef>MessageIdentifier</siri:RequestMessageRef>
        <siri:TerminationResponseStatus>
          <siri:ResponseTimestamp>2017-01-01T12:00:40.000Z</siri:ResponseTimestamp>
          <siri:SubscriberRef>remote_credential</siri:SubscriberRef>
          <siri:SubscriptionRef/>
          <siri:Status>true</siri:Status>
        </siri:TerminationResponseStatus>
      </Answer>
      <AnswerExtension/>
    </sw:DeleteSubscriptionResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | remote_credential              |
      | local_credential                   | local_credential               |
      | remote_code_space               | internal                       |
    When I send this SIRI request to the Referential "test"
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <ns6:NotifyStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
          xmlns:ns3="http://www.ifopt.org.uk/acsb"
          xmlns:ns4="http://www.ifopt.org.uk/ifopt"
          xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns6="http://wsdl.siri.org.uk"
          xmlns:ns7="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>local_credential </ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>
              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>local_credential</ns2:SubscriberRef>
                <ns2:SubscriptionRef>dummy</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>enRoute:VehicleJourney:201-enRoute:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>enRoute:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>enRoute:Line:3:LOC</ns3:LineRef>
                    <ns3:DirectionRef>Left</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>enRoute:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>enRoute:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>enRoute:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OperatorRef>enRoute:Company:15563880:LOC</ns3:OperatorRef>
                    <ns3:ProductCategoryRef>0</ns3:ProductCategoryRef>
                    <ns3:VehicleFeatureRef>TRFC_M4_1</ns3:VehicleFeatureRef>
                    <ns3:OriginRef>enRoute:StopPoint:SP:42:LOC</ns3:OriginRef>
                    <ns3:OriginName>Magicien Noir</ns3:OriginName>
                    <ns3:DestinationRef>enRoute:StopPoint:SP:62:LOC</ns3:DestinationRef>
                    <ns3:DestinationName>Cimetière des Sauvages</ns3:DestinationName>
                    <ns3:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</ns3:OriginAimedDepartureTime>
                    <ns3:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</ns3:DestinationAimedArrivalTime>
                    <ns3:Monitored>true</ns3:Monitored>
                    <ns3:ProgressRate>normalProgress</ns3:ProgressRate>
                    <ns3:Delay>P0Y0M0DT0H0M0.000S</ns3:Delay>
                    <ns3:CourseOfJourneyRef>201</ns3:CourseOfJourneyRef>
                    <ns3:VehicleRef>enRoute:Vehicle:23:LOC</ns3:VehicleRef>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>enRoute:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R)</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T13:00:00.000+02:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>delayed</ns3:ArrivalStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns2:StopMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </ns6:NotifyStopMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then an audit event should exist with these attributes:
      | Type                    | DeleteSubscriptionRequest            |
      | Protocol                | siri                                 |
      | Direction               | sent                                 |
      | Status                  | OK                                   |
      | Partner                 | test                                 |
      | SubscriptionIdentifiers | ["dummy"]                            |
      | ProcessingTime          | 0                                    |

  @ARA-880
  Scenario: Audit a referential with hyphen in the slug
    Given a Referential "test-with-hyphen" is created
    And a Partner "test" exists in Referential "test-with-hyphen" with connectors [siri-check-status-server] and the following settings:
      | local_credential | test |
    When I send this SIRI request to the Referential "test-with-hyphen"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test</siri:RequestorRef>
        <siri:MessageIdentifier>enRoute:Message::test</siri:MessageIdentifier>
      </Request>
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then an audit event should exist with these attributes:
      | Type               | CheckStatusRequest |
      | Dataset            | cucumber_test_with_hyphen   |

  @ARA-880
  Scenario: Audit a referential without hyphen in the slug
    Given a Referential "test_without_hyphen" is created
    And a Partner "test" exists in Referential "test_without_hyphen" with connectors [siri-check-status-server] and the following settings:
      | local_credential | test |
    When I send this SIRI request to the Referential "test_without_hyphen"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test</siri:RequestorRef>
        <siri:MessageIdentifier>enRoute:Message::test</siri:MessageIdentifier>
      </Request>
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then an audit event should exist with these attributes:
      | Type               | CheckStatusRequest           |
      | Dataset            | cucumber_test_without_hyphen |

  @ARA-1060
  Scenario: Audit a recevied SIRI EstimatedTimetableSubscriptionRequest
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_code_space | internal              |
       | siri.envelope        | raw                   |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:A:BUS" |
      | Name      | Ligne A Bus                     |
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
    Then an audit event should exist with these attributes:
      | Type                                 | EstimatedTimetableSubscriptionRequest   |
      | Direction                            | received                                |
      | Protocol                             | siri                                    |
      | Partner                              | test                                    |
      | Status                               | OK                                      |
      | SubscriptionIdentifiers              | ["test1"]                               |

  @ARA-1086
  Scenario: Audit a ProductionTimetableRequest with duplicate SubscriptionIdentifier with an existing EstimatedTimetable subscription
    Given a raw SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-production-timetable-subscription-broadcaster] and the following settings:
       | remote_url                         | http://localhost:8090 |
       | remote_credential                  | ara                   |
       | local_credential                   | test                  |
       | remote_code_space               | internal              |
       | siri.envelope                      | raw                   |
       | broadcast.subscriptions.persistent | true                  |
    And a StopArea exists with the following attributes:
      | Name      | Test 24                                  |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Lines     | ["6ba7b814-9dad-11d1-4-00c04fd430c8"]    |
    And a StopArea exists with the following attributes:
      | Name      | Test 25                                  |
      | Codes | "internal": "NINOXE:StopPoint:SP:25:LOC" |
      | Lines     | ["6ba7b814-9dad-11d1-4-00c04fd430c8"]    |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                               | Passage 32                              |
      | Codes                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                             | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
      | DirectionType                      | outbound                                |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"              |
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
    Then an audit event should exist with these attributes:
      | Type                                 | ProductionTimetableSubscriptionRequest  |
      | Direction                            | received                                |
      | Protocol                             | siri                                    |
      | Partner                              | test                                    |
      | Status                               | Error                                   |

  @ARA-1152
  Scenario: Audit a send EstimatedTimetable subscription request
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
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_code_space | internal              |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Test                   |
      | Codes | "internal": "testLine" |
    And a minute has passed
    And a minute has passed
    And  an audit event should exist with these attributes:
      | Type      | EstimatedTimetableSubscriptionRequest |
      | Direction | sent                                  |
      | Protocol  | siri                                  |
      | Partner   | test                                  |
      | Status    | OK                                    |

  @ARA-1152
  Scenario: Audit a received SIRI EstimatedTimetable Notification
    Given a Partner "test" exists with connectors [siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | test                  |
      | remote_code_space | internal              |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect             |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |

    When I send this SIRI request to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
<sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
  <ServiceDeliveryInfo>
    <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
    <siri:ProducerRef>test</siri:ProducerRef>
    <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
    <siri:RequestMessageRef>enRoute:Message::test</siri:RequestMessageRef>
  </ServiceDeliveryInfo>
  <Notification>
    <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>6ba7b814-9dad-11d1-2-00c04fd430c8</siri:SubscriptionRef>
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
  <SiriExtension/>
</sw:NotifyEstimatedTimetable>
</S:Body>
</S:Envelope>
      """
    Then an audit event should exist with these attributes:
      | Type              | NotifyEstimatedTimetable |
      | Direction         | received                 |
      | Protocol          | siri                     |
      | Partner           | test                     |
      | RequestIdentifier | enRoute:Message::test    |
      | Lines             | ["NINOXE:Line:3:LOC"]    |

  @ARA-1385
  Scenario: Audit a StopMonitoringDelivery in a subscription with multiple StopMonitoringDelivery
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_code_space               | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
      | siri.direction_type                | Aller,Retour                   |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:25:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[1] | StopArea, "internal": "NINOXE:StopPoint:SP:25:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                              |
      | Codes                | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                   | 6ba7b814-9dad-11d1-c-00c04fd430c8       |
      | Monitored                | true                                    |
      | Attribute[DirectionName] | A Direction Name                        |
    And a minute has passed
    When I send this SIRI request
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <ns6:NotifyStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
          xmlns:ns3="http://www.ifopt.org.uk/acsb"
          xmlns:ns4="http://www.ifopt.org.uk/ifopt"
          xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns6="http://wsdl.siri.org.uk"
          xmlns:ns7="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns2:ResponseTimestamp>
              2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>
              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:DirectionRef>Aller</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OperatorRef>NINOXE:Company:15563880:LOC</ns3:OperatorRef>
                    <ns3:ProductCategoryRef>0</ns3:ProductCategoryRef>
                    <ns3:VehicleFeatureRef>TRFC_M4_1</ns3:VehicleFeatureRef>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:42:LOC</ns3:OriginRef>
                    <ns3:OriginName>Magicien Noir</ns3:OriginName>
                    <ns3:DestinationRef>NINOXE:StopPoint:SP:62:LOC</ns3:DestinationRef>
                    <ns3:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</ns3:OriginAimedDepartureTime>
                    <ns3:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</ns3:DestinationAimedArrivalTime>
                    <ns3:Monitored>true</ns3:Monitored>
                    <ns3:ProgressRate>normalProgress</ns3:ProgressRate>
                    <ns3:Delay>P0Y0M0DT0H0M0.000S</ns3:Delay>
                    <ns3:CourseOfJourneyRef>201</ns3:CourseOfJourneyRef>
                    <ns3:VehicleRef>NINOXE:Vehicle:23:LOC</ns3:VehicleRef>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R)</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T13:00:00.000+02:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>delayed</ns3:ArrivalStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns2:StopMonitoringDelivery>
              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-4</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:25:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:8:LOC</ns3:LineRef>
                    <ns3:DirectionRef>Aller</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OperatorRef>NINOXE:Company:15563880:LOC</ns3:OperatorRef>
                    <ns3:ProductCategoryRef>0</ns3:ProductCategoryRef>
                    <ns3:VehicleFeatureRef>TRFC_M4_1</ns3:VehicleFeatureRef>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:42:LOC</ns3:OriginRef>
                    <ns3:OriginName>Magicien Noir</ns3:OriginName>
                    <ns3:DestinationRef>NINOXE:StopPoint:SP:62:LOC</ns3:DestinationRef>
                    <ns3:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</ns3:OriginAimedDepartureTime>
                    <ns3:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</ns3:DestinationAimedArrivalTime>
                    <ns3:Monitored>true</ns3:Monitored>
                    <ns3:ProgressRate>normalProgress</ns3:ProgressRate>
                    <ns3:Delay>P0Y0M0DT0H0M0.000S</ns3:Delay>
                    <ns3:CourseOfJourneyRef>201</ns3:CourseOfJourneyRef>
                    <ns3:VehicleRef>NINOXE:Vehicle:24:LOC</ns3:VehicleRef>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:25:LOC</ns3:StopPointRef>
                      <ns3:Order>6</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R) new</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T13:00:00.000+02:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>delayed</ns3:ArrivalStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns2:StopMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </ns6:NotifyStopMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol           | siri                                    |
      | Direction          | received                                |
      | ResponseIdentifier | /{uuid}/                                |
      | Status             | OK                                      |
      | Type               | NotifyStopMonitoring                    |
      | StopAreas          | ["NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:25:LOC"] |
      | VehicleJourneys    | ["NINOXE:VehicleJourney:201", "NINOXE:VehicleJourney:202"] |
      | Lines              | ["NINOXE:Line:3:LOC", "NINOXE:Line:8:LOC"] |

  @ARA-1385
  Scenario: Audit a VehicleMonitoringDelivery in a subscription with multiple VehicleMonitoringDeliveries
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
      <AnswerExtension/>
    </ns1:SubscribeResponse>
  </S:Body>
  </S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-vehicle-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_code_space               | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name      | Test                   |
      | Codes | "internal": "testLine" |
    And a Line exists with the following attributes:
      | Name      | Test1                  |
      | Codes | "internal": "testLine1" |
    And a Subscription exist with the following attributes:
      | Kind              | VehicleMonitoringCollect     |
      | ReferenceArray[0] | Line, "internal": "testLine" |
    And a Subscription exist with the following attributes:
      | Kind              | VehicleMonitoringCollect     |
      | ReferenceArray[0] | Line, "internal": "testLine1" |
    And a minute has passed
    When I send this SIRI request
      """
      <?xml version='1.0' encoding='utf-8'?>
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <ns6:NotifyVehicleMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
          xmlns:ns3="http://www.ifopt.org.uk/acsb"
          xmlns:ns4="http://www.ifopt.org.uk/ifopt"
          xmlns:ns5="http://www.siri.org.uk/siri"
          xmlns:ns6="http://wsdl.siri.org.uk"
          xmlns:ns7="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>VehicleMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>
              <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns5:ResponseTimestamp>2022-06-25T15:08:14.940+02:00</ns5:ResponseTimestamp>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns5:Status>true</ns5:Status>
                <ns5:VehicleActivity>
                  <ns5:RecordedAtTime>2022-06-25T15:08:14.928+02:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>108</ns5:ItemIdentifier>
                  <ns5:ValidUntilTime>2022-06-25T16:08:14.928+02:00</ns5:ValidUntilTime>
                  <ns5:VehicleMonitoringRef>108</ns5:VehicleMonitoringRef>
                  <ns5:ProgressBetweenStops>
                    <ns5:LinkDistance>340.0</ns5:LinkDistance>
                    <ns5:Percentage>73.0</ns5:Percentage>
                  </ns5:ProgressBetweenStops>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>testLine</ns5:LineRef>
                    <ns5:DirectionRef>Aller</ns5:DirectionRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>NAVINEO:DataFrame::1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>RDMANTOIS:VehicleJourney::6628652:LOC</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>RDMANTOIS:JourneyPattern::LCP37:LOC</ns5:JourneyPatternRef>
                    <ns5:JourneyPatternName>LCP37</ns5:JourneyPatternName>
                    <ns5:PublishedLineName>testLine</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>OPERYORDM:Operator::OPERYORDM:LOC</ns5:OperatorRef>
                    <ns5:OriginRef>50000037</ns5:OriginRef>
                    <ns5:OriginName>Port Fouquet</ns5:OriginName>
                    <ns5:DestinationRef>50000031</ns5:DestinationRef>
                    <ns5:DestinationName>Mantes la Jolie Gare routière - Quai 2</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:VehicleLocation srsName="2154">
                      <ns5:Coordinates>603204 6878517</ns5:Coordinates>
                    </ns5:VehicleLocation>
                    <ns5:Bearing>171.0</ns5:Bearing>
                    <ns5:VehicleRef>TRANSDEV:Vehicle::1501:LOC</ns5:VehicleRef>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>50000016</ns5:StopPointRef>
                      <ns5:Order>9</ns5:Order>
                      <ns5:StopPointName>Hôpital F. Quesnay</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>MantesLJ Gare</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2022-06-25T15:05:00.000+02:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2022-06-25T15:08:27.000+02:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2022-06-25T15:05:00.000+02:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2022-06-25T15:08:27.000+02:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                  <ns5:Extensions/>
                </ns5:VehicleActivity>
              </ns5:VehicleMonitoringDelivery>
              <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns5:ResponseTimestamp>2022-06-25T15:08:14.940+02:00</ns5:ResponseTimestamp>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns5:Status>true</ns5:Status>
                <ns5:VehicleActivity>
                  <ns5:RecordedAtTime>2022-06-25T15:08:14.928+02:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>108</ns5:ItemIdentifier>
                  <ns5:ValidUntilTime>2022-06-25T16:08:14.928+02:00</ns5:ValidUntilTime>
                  <ns5:VehicleMonitoringRef>108</ns5:VehicleMonitoringRef>
                  <ns5:ProgressBetweenStops>
                    <ns5:LinkDistance>350.0</ns5:LinkDistance>
                    <ns5:Percentage>80.0</ns5:Percentage>
                  </ns5:ProgressBetweenStops>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>testLine1</ns5:LineRef>
                    <ns5:DirectionRef>Aller</ns5:DirectionRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>NAVINEO:DataFrame::1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>RDMANTOIS:VehicleJourney::6628658:LOC</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>RDMANTOIS:JourneyPattern::LCP38:LOC</ns5:JourneyPatternRef>
                    <ns5:JourneyPatternName>LCP38</ns5:JourneyPatternName>
                    <ns5:PublishedLineName>Test1</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>OPERYORDM:Operator::OPERYORDM:LOC</ns5:OperatorRef>
                    <ns5:OriginRef>50000037</ns5:OriginRef>
                    <ns5:OriginName>Port Fouquet</ns5:OriginName>
                    <ns5:DestinationRef>50000031</ns5:DestinationRef>
                    <ns5:DestinationName>Mantes la Jolie Gare routière - Quai 3</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:VehicleLocation srsName="2154">
                      <ns5:Coordinates>603204 6878517</ns5:Coordinates>
                    </ns5:VehicleLocation>
                    <ns5:Bearing>171.0</ns5:Bearing>
                    <ns5:VehicleRef>TRANSDEV:Vehicle::1502:LOC</ns5:VehicleRef>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>50000020</ns5:StopPointRef>
                      <ns5:Order>10</ns5:Order>
                      <ns5:StopPointName>Hôpital F. Quesnay</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>MantesLJ Gare</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2022-06-25T15:05:00.000+02:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2022-06-25T15:08:27.000+02:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2022-06-25T15:05:00.000+02:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2022-06-25T15:08:27.000+02:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                  <ns5:Extensions/>
                </ns5:VehicleActivity>
              </ns5:VehicleMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </ns6:NotifyVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
      Then one Vehicle has the following attributes:
        | Codes | "internal": "TRANSDEV:Vehicle::1501:LOC" |
        | LineId    |        6ba7b814-9dad-11d1-3-00c04fd430c8 |
        | Bearing   |                                    171.0 |
        | Latitude  |                        48.99927561424598 |
        | Longitude |                       1.6770970859674874 |
      Then an audit event should exist with these attributes:
        | Type            | NotifyVehicleMonitoring                                                            |
        | Protocol        | siri                                                                               |
        | Direction       | received                                                                           |
        | Status          | OK                                                                                 |
        | Partner         | test                                                                               |
        | Vehicles        | ["TRANSDEV:Vehicle::1501:LOC", "TRANSDEV:Vehicle::1502:LOC"]                       |
        | VehicleJourneys | ["RDMANTOIS:VehicleJourney::6628652:LOC", "RDMANTOIS:VehicleJourney::6628658:LOC"] |
        | StopAreas       | ["50000016", "50000020"]                                                           |
        | Lines           | ["testLine", "testLine1"]                                                          |

  @ARA-1385 
  Scenario: Audit a EstimatedTimetableDelivery in a subscription with multiple ExistimatedTimetableDeliveries
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
      <AnswerExtension/>
    </ns1:SubscribeResponse>
  </S:Body>
  </S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name      | Test                            |
      | Codes | "internal": "NINOXE:Line:3:LOC" |
    And a Line exists with the following attributes:
      | Name      | Test                            |
      | Codes | "internal": "NINOXE:Line:4:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect             |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect             |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:4:LOC" |
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
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:SubscriptionRef>
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
    <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
        <siri:EstimatedVehicleJourney>
          <siri:LineRef>NINOXE:Line:4:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
          <siri:DestinationRef>ThisIsTheEnd</siri:DestinationRef>
          <siri:EstimatedCalls>
            <siri:EstimatedCall>
              <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
              <siri:Order>8</siri:Order>
              <siri:StopPointName>Test2</siri:StopPointName>
              <siri:VehicleAtStop>false</siri:VehicleAtStop>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>onTime</siri:ArrivalStatus>
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
    And an audit event should exist with these attributes:
      | Protocol        | siri                                                         |
      | Direction       | received                                                     |
      | Status          | OK                                                           |
      | Type            | NotifyEstimatedTimetable                                     |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:25:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201", "NINOXE:VehicleJourney:202"]   |
      | Lines           | ["NINOXE:Line:3:LOC", "NINOXE:Line:4:LOC"]                   |
