Feature: Audit API exchanges

  Background:
    Given a Referential "test" is created

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
      | remote_objectid_kind       | internal                   |
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
      | remote_objectid_kind               | internal                       |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name      | Test                                      |
      | ObjectIDs | "internal": "enRoute:StopPoint:SP:24:LOC" |
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
      | SubscriptionIdentifiers | ["6ba7b814-9dad-11d1-5-00c04fd430c8"] |
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
      | remote_objectid_kind               | internal                       |
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
