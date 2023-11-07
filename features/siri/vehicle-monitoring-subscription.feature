Feature: Support SIRI VehicleMonitoring by subscription

  Background:
      Given a Referential "test" is created

  @ARA-1306
  Scenario: VehicleMonitoring subscription collect should send VehicleMonitoringSubscriptionRequest to partner
   Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-vehicle-monitoring-subscription-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_objectid_kind  | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
      | local_credential      | ara                   |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Test 1                             |
      | ObjectIDs | "internal": "RLA_Bus:Line::05:LOC" |
   And a minute has passed
   And 20 seconds have passed
   Then the SIRI server should have received a VehicleMonitoringSubscriptionRequest request with:
      | //siri:LineRef | RLA_Bus:Line::05:LOC |

  @ARA-1306
  Scenario: VehicleMonitoring subscription collect and partner CheckStatus is unavailable should not send VehicleMonitoringSubscriptionRequest to partner
   Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-vehicle-monitoring-subscription-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_objectid_kind  | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
      | local_credential      | ara                   |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Test 1                             |
      | ObjectIDs | "internal": "RLA_Bus:Line::05:LOC" |
   And a minute has passed
   And 10 seconds have passed
   Then the SIRI server should not have received a VehicleMonitoringSubscription request

  @ARA-1306
  Scenario: VehicleMonitoring subscription collect and partner CheckStatus is unavailable should send VehicleMonitoringSubscriptionRequest to partner whith setting collect.subscriptions.persistent
   Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-vehicle-monitoring-subscription-collector] and the following settings:
      | remote_url                       | http://localhost:8090 |
      | remote_credential                | test                  |
      | remote_objectid_kind             | internal              |
      | collect.include_lines            | RLA_Bus:Line::05:LOC  |
      | local_credential                 | ara                   |
      | collect.subscriptions.persistent | true                  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Test 1                             |
      | ObjectIDs | "internal": "RLA_Bus:Line::05:LOC" |
   And a minute has passed
   And 30 seconds have passed
   Then the SIRI server should have received a VehicleMonitoringSubscriptionRequest request with:
      | //siri:LineRef | RLA_Bus:Line::05:LOC |

  @ARA-1306
  Scenario: VehicleMonitoring subscription collect and partner CheckStatus is unavailable should send VehicleMonitoringSubscriptionRequest to partner whith setting collect.persistent
   Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-vehicle-monitoring-subscription-collector] and the following settings:
      | remote_url                       | http://localhost:8090 |
      | remote_credential                | test                  |
      | remote_objectid_kind             | internal              |
      | collect.include_lines            | RLA_Bus:Line::05:LOC  |
      | local_credential                 | ara                   |
      | collect.persistent               | true                  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Test 1                             |
      | ObjectIDs | "internal": "RLA_Bus:Line::05:LOC" |
   And a minute has passed
   And 10 seconds have passed
   Then the SIRI server should have received a VehicleMonitoringSubscriptionRequest request with:
      | //siri:LineRef | RLA_Bus:Line::05:LOC |

  @ARA-1236 @siri-valid
  Scenario: Send a VehicleMonitoring notification when a vehicle changes
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, siri-vehicle-monitoring-subscription-broadcaster] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | Ara                   |
      | local_credential      | Subscriber            |
      | remote_objectid_kind  | internal              |
      | sort_payload_for_test | true                  |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "Test:Line:3:LOC" |
      | Name      | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                |
      | ObjectIDs                | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored                | true                                      |
      | Attribute[DirectionName] | Direction Name                            |
    And a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "Test:Vehicle:201123:LOC" |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8     |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8     |
    And a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "Test:Vehicle:999:LOC" |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8     |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8     |
    And a Subscription exist with the following attributes:
      | Kind              | VehicleMonitoringBroadcast       |
      | SubscriberRef     | Subscriber                       |
      | ExternalId        | subscription-1                   |
      | ReferenceArray[0] | Line, "internal": "Test:Line:3:LOC" |
    When the Vehicle "internal:Test:Vehicle:201123:LOC" is edited with the following attributes:
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 234                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
    When the Vehicle "internal:Test:Vehicle:999:LOC" is edited with the following attributes:
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 12.234                            |
      | Latitude         | 8.678                             |
      | Bearing          | 126                               |
      | RecordedAtTime   | 2017-01-01T12:10:00.000Z          |
      | ValidUntilTime   | 2017-01-01T13:59:00.000Z          |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:NotifyVehicleMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:10.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            </ServiceDeliveryInfo>
            <Notification>
              <siri:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:10.000Z</siri:ResponseTimestamp>
                <siri:SubscriberRef>Subscriber</siri:SubscriberRef>
                <siri:SubscriptionRef>subscription-1</siri:SubscriptionRef>
                <siri:Status>true</siri:Status>
                <siri:VehicleActivity>
                  <siri:RecordedAtTime>2017-01-01T13:00:00.000Z</siri:RecordedAtTime>
                  <siri:ValidUntilTime>2017-01-01T14:00:00.000Z</siri:ValidUntilTime>
                  <siri:VehicleMonitoringRef>Test:Vehicle:201123:LOC</siri:VehicleMonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>Test:VehicleJourney:201:LOC</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DirectionName>Direction Name</siri:DirectionName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:VehicleLocation>
                      <siri:Longitude>1.234</siri:Longitude>
                      <siri:Latitude>5.678</siri:Latitude>
                    </siri:VehicleLocation>
                    <siri:Bearing>234</siri:Bearing>
                  </siri:MonitoredVehicleJourney>
                </siri:VehicleActivity>
                <siri:VehicleActivity>
                  <siri:RecordedAtTime>2017-01-01T12:10:00.000Z</siri:RecordedAtTime>
                  <siri:ValidUntilTime>2017-01-01T13:59:00.000Z</siri:ValidUntilTime>
                  <siri:VehicleMonitoringRef>Test:Vehicle:999:LOC</siri:VehicleMonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>Test:Line:3:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>Test:VehicleJourney:201:LOC</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DirectionName>Direction Name</siri:DirectionName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:VehicleLocation>
                      <siri:Longitude>12.234</siri:Longitude>
                      <siri:Latitude>8.678</siri:Latitude>
                    </siri:VehicleLocation>
                    <siri:Bearing>126</siri:Bearing>
                    </siri:MonitoredVehicleJourney>
                  </siri:VehicleActivity>
              </siri:VehicleMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </sw:NotifyVehicleMonitoring>
        </S:Body>
      </S:Envelope>
      """
    Then an audit event should exist with these attributes:
      | Type                    | NotifyVehicleMonitoring                             |
      | Direction               | sent                                                |
      | Protocol                | siri                                                |
      | Partner                 | test                                                |
      | Status                  | OK                                                  |
      | SubscriptionIdentifiers | ["subscription-1"]                                  |
      | Lines                   | ["Test:Line:3:LOC"]                                 |
      | Vehicles                | ["Test:Vehicle:201123:LOC", "Test:Vehicle:999:LOC"] |

  @ARA-1236 @siri-valid
  Scenario: Delete and recreate SIRI VehicleMonitoring request for subscription when receiving subscription with same existing number
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-vehicle-monitoring-subscription-broadcaster] and the following settings:
       | remote_url                        | http://localhost:8090 |
       | remote_credential                 | test                  |
       | local_credential                  | NINOXE:default        |
       | remote_objectid_kind              | internal              |
       | broadcast.subscription_persistent | true                  |
    And a minute has passed
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:VehicleMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>test1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:VehicleMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
              </siri:VehicleMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
              <siri:ChangeBeforeUpdates>PT1M</siri:ChangeBeforeUpdates>
            </siri:VehicleMonitoringSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>
    </S:Envelope>
      """
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Kind            | VehicleMonitoringBroadcast        |
      | ExternalId      | test1                             |
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-a-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:VehicleMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>test1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:VehicleMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-a-00c04fd430c8</siri:MessageIdentifier>
              </siri:VehicleMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
              <siri:ChangeBeforeUpdates>PT1M</siri:ChangeBeforeUpdates>
            </siri:VehicleMonitoringSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>
    </S:Envelope>
     """
    Then No Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Kind            | VehicleMonitoringBroadcast        |
      | ExternalId      | test1                             |
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Kind            | VehicleMonitoringBroadcast        |
      | ExternalId      | test1                             |
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-b-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:VehicleMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>test2</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:VehicleMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-b-00c04fd430c8</siri:MessageIdentifier>
              </siri:VehicleMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
              <siri:ChangeBeforeUpdates>PT1M</siri:ChangeBeforeUpdates>
            </siri:VehicleMonitoringSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>
    </S:Envelope>
     """
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Kind            | VehicleMonitoringBroadcast        |
      | ExternalId      | test1                             |
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Kind            | VehicleMonitoringBroadcast        |
      | ExternalId      | test2                             |

  @ARA-1236 @siri-valid
  Scenario: Handle a SIRI VehicleMonitoring request for subscription to all lines
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-vehicle-monitoring-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
    And a minute has passed
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:VehicleMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>test1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:VehicleMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
              </siri:VehicleMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
              <siri:ChangeBeforeUpdates>PT1M</siri:ChangeBeforeUpdates>
            </siri:VehicleMonitoringSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>
    </S:Envelope>
      """
    Then one Subscription exists with the following attributes:
      | Kind | VehicleMonitoringBroadcast |
    Then an audit event should exist with these attributes:
      | Type                    | VehicleMonitoringSubscriptionRequest |
      | Direction               | received                             |
      | Protocol                | siri                                 |
      | Partner                 | test                                 |
      | Status                  | OK                                   |
      | SubscriptionIdentifiers | ["test1"]                            |

  @ARA-1236 @siri-valid
  Scenario: Handle a SIRI VehicleMonitoring request for subscription for all existing lines in a referential only with same remote_objectid_kind
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-vehicle-monitoring-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
    And a Line exists with the following attributes:
      | ObjectIDs | "another": "NINOXE:Line:3:LOC"  |
      | Name      | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:A:BUS" |
      | Name      | Ligne A Bus                     |
    And a minute has passed
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:VehicleMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>6ba7b814-9dad-11d1--00c04fd430c8</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:VehicleMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
              </siri:VehicleMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
              <siri:ChangeBeforeUpdates>PT1M</siri:ChangeBeforeUpdates>
            </siri:VehicleMonitoringSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>
    </S:Envelope>
      """
    Then Subscriptions exist with the following resources:
      | internal | NINOXE:Line:A:BUS |
    Then No Subscriptions exist with the following resources:
      | internal | NINOXE:Line:3:LOC |

  @ARA-1236 @siri-valid
  Scenario: Handle a SIRI VehicleMonitoring request for subscription to a line not existing in a referential
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-vehicle-monitoring-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
    And a Line exists with the following attributes:
      | ObjectIDs | "another": "NINOXE:Line:3:LOC"  |
      | Name      | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:A:BUS" |
      | Name      | Ligne A Bus                     |
    And a minute has passed
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:VehicleMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>subscription-1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:VehicleMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
                <siri:LineRef>testLine</siri:LineRef>
              </siri:VehicleMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
              <siri:ChangeBeforeUpdates>PT1M</siri:ChangeBeforeUpdates>
            </siri:VehicleMonitoringSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>
    </S:Envelope>
      """
    Then I should receive this SIRI response
      """
     <?xml version='1.0' encoding='UTF-8'?> 
     <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
       <S:Body>
         <sw:SubscribeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
           <SubscriptionAnswerInfo>
             <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
             <siri:ResponderRef>test</siri:ResponderRef>
             <siri:RequestMessageRef xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xsi:type='siri:MessageRefStructure'>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:RequestMessageRef>
           </SubscriptionAnswerInfo>
           <Answer>
             <siri:ResponseStatus>
               <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
               <siri:RequestMessageRef>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:RequestMessageRef>
               <siri:SubscriberRef>test</siri:SubscriberRef>
               <siri:SubscriptionRef>subscription-1</siri:SubscriptionRef>
               <siri:Status>false</siri:Status>
               <siri:ErrorCondition>
                 <siri:InvalidDataReferencesError>
                   <siri:ErrorText>Unknown Line(s) testLine</siri:ErrorText>
                 </siri:InvalidDataReferencesError>
               </siri:ErrorCondition>
             </siri:ResponseStatus>
             <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
           </Answer>
           <AnswerExtension/>
         </sw:SubscribeResponse>
       </S:Body>
     </S:Envelope>
      """
    Then an audit event should exist with these attributes:
      | Type                    | VehicleMonitoringSubscriptionRequest |
      | Direction               | received                             |
      | Protocol                | siri                                 |
      | Partner                 | test                                 |
      | Status                  | Error                                |
      | Lines                   | ["testLine"]                         |

   @ARA-1236 @siri-valid
  Scenario: Handle a SIRI VehicleMonitoring request for subscription to a single line
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-vehicle-monitoring-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a minute has passed
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:VehicleMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>subscription-1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:VehicleMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
                <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
              </siri:VehicleMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
              <siri:ChangeBeforeUpdates>PT1M</siri:ChangeBeforeUpdates>
            </siri:VehicleMonitoringSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>
    </S:Envelope>
      """
    Then I should receive this SIRI response
      """
     <?xml version='1.0' encoding='UTF-8'?> 
     <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
       <S:Body>
         <sw:SubscribeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
           <SubscriptionAnswerInfo>
             <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
             <siri:ResponderRef>test</siri:ResponderRef>
             <siri:RequestMessageRef xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xsi:type='siri:MessageRefStructure'>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:RequestMessageRef>
           </SubscriptionAnswerInfo>
           <Answer>
             <siri:ResponseStatus>
               <siri:ResponseTimestamp>2017-01-01T12:01:00.000Z</siri:ResponseTimestamp>
               <siri:RequestMessageRef>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:RequestMessageRef>
               <siri:SubscriberRef>test</siri:SubscriberRef>
               <siri:SubscriptionRef>subscription-1</siri:SubscriptionRef>
               <siri:Status>true</siri:Status>
               <siri:ValidUntil>2017-01-03T12:03:00.000Z</siri:ValidUntil>
             </siri:ResponseStatus>
             <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
           </Answer>
           <AnswerExtension/>
         </sw:SubscribeResponse>
       </S:Body>
     </S:Envelope>
      """
    Then an audit event should exist with these attributes:
      | Type                      | VehicleMonitoringSubscriptionRequest |
      | Direction                 | received                             |
      | Protocol                  | siri                                 |
      | Partner                   | test                                 |
      | Status                    | OK                                   |
      | Lines                     | ["NINOXE:Line:3:LOC"]                |
      | SubscriptionIdentifiers   | ["subscription-1"]                   |

  Scenario: Create Vehicle Monitoring subscription by Line
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-vehicle-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_objectid_kind               | internal                       |
    And a Line exists with the following attributes:
      | Name      | Test                   |
      | ObjectIDs | "internal": "testLine" |
    When a minute has passed
    And a minute has passed
    And a minute has passed
    Then one Subscription exists with the following attributes:
      | Kind                      | VehicleMonitoringCollect |
      | Resources[0]/SubscribedAt | > 2017-01-01T12:01:00Z     |

  Scenario: Update a Vehicle after a VehicleMonitoringDelivery in a subscription
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
      | remote_objectid_kind               | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name      | Test                   |
      | ObjectIDs | "internal": "testLine" |
    And a Subscription exist with the following attributes:
      | Kind              | VehicleMonitoringCollect     |
      | ReferenceArray[0] | Line, "internal": "testLine" |
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
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</ns2:SubscriptionRef>
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
            </Notification>
            <SiriExtension />
          </ns6:NotifyVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
      And I see ara vehicles
      Then one Vehicle has the following attributes:
        | ObjectIDs | "internal": "TRANSDEV:Vehicle::1501:LOC" |
        | LineId    |        6ba7b814-9dad-11d1-3-00c04fd430c8 |
        | Bearing   |                                    171.0 |
        | Latitude  |                        48.99927561424598 |
        | Longitude |                       1.6770970859674874 |
      Then an audit event should exist with these attributes:
        | Type            | NotifyVehicleMonitoring                   |
        | Protocol        | siri                                      |
        | Direction       | received                                  |
        | Status          | OK                                        |
        | Partner         | test                                      |
        | Vehicles        | ["TRANSDEV:Vehicle::1501:LOC"]            |
        | VehicleJourneys | ["RDMANTOIS:VehicleJourney::6628652:LOC"] |
        | StopAreas       | ["50000016"]                              |
        | Lines           | ["testLine"]                              |

  @ARA-1101
  Scenario: Update a Vehicle after a VehicleMonitoringDelivery in a subscription using the partner setting siri.direction_type should update the DirectionRef
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
      | remote_objectid_kind               | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
      | siri.direction_type                | Aller, Retour                  |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name      | Test                   |
      | ObjectIDs | "internal": "testLine" |
    And a Subscription exist with the following attributes:
      | Kind              | VehicleMonitoringCollect     |
      | ReferenceArray[0] | Line, "internal": "testLine" |
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
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</ns2:SubscriptionRef>
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
                    <ns5:LineRef>C</ns5:LineRef>
                    <ns5:DirectionRef>Aller</ns5:DirectionRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>NAVINEO:DataFrame::1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>RDMANTOIS:VehicleJourney::6628652:LOC</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>RDMANTOIS:JourneyPattern::LCP37:LOC</ns5:JourneyPatternRef>
                    <ns5:JourneyPatternName>LCP37</ns5:JourneyPatternName>
                    <ns5:PublishedLineName>C</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
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
            </Notification>
            <SiriExtension />
          </ns6:NotifyVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
      Then the VehicleJourney "6ba7b814-9dad-11d1-a-00c04fd430c8" has the following attributes:
      | DirectionType | inbound |
