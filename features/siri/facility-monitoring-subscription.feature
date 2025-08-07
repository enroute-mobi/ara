Feature: Support SIRI FacilityMonitoring by subscription

  Background:
      Given a Referential "test" is created

  @ARA-1755
  Scenario: Handle a SIRI FacilityMonitoring request for subscription to a single Facility
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-facility-monitoring-subscription-broadcaster] and the following settings:
       | remote_url        | http://localhost:8090 |
       | remote_credential | test                  |
       | local_credential  | NINOXE:default        |
       | remote_code_space | internal              |
    And a Facility exists with the following attributes:
      | Codes[internal] | NINOXE:Facility:1:LOC |
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
            <siri:FacilityMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>subscription-1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:FacilityMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
                <siri:FacilityRef>NINOXE:Facility:1:LOC</siri:FacilityRef>
              </siri:FacilityMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
            </siri:FacilityMonitoringSubscriptionRequest>
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
      | Type                    | FacilityMonitoringSubscriptionRequest |
      | Direction               | received                              |
      | Protocol                | siri                                  |
      | Partner                 | test                                  |
      | Status                  | OK                                    |
      | SubscriptionIdentifiers | ["subscription-1"]                    |

  @ARA-1755
  Scenario: Handle a SIRI FacilityMonitoring request for subscription to an unknown facility
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-facility-monitoring-subscription-broadcaster] and the following settings:
       | remote_url        | http://localhost:8090 |
       | remote_credential | test                  |
       | local_credential  | NINOXE:default        |
       | remote_code_space | internal              |
    And a Facility exists with the following attributes:
      | Codes[internal] | NINOXE:Facility:1:LOC |
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
            <siri:FacilityMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>subscription-1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:FacilityMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
                <siri:FacilityRef>WRONG</siri:FacilityRef>
              </siri:FacilityMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
            </siri:FacilityMonitoringSubscriptionRequest>
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
                    <siri:ErrorText>Unknown Facility(ies) WRONG</siri:ErrorText>
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
      | Type                    | FacilityMonitoringSubscriptionRequest |
      | Direction               | received                              |
      | Protocol                | siri                                  |
      | Partner                 | test                                  |
      | Status                  | Error                                 |

  @ARA-1755
  Scenario: Send a FacilityMonitoring notification when a facility status changes
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-facility-monitoring-subscription-broadcaster] and the following settings:
       | remote_url        | http://localhost:8090 |
       | remote_credential | test                  |
       | local_credential  | NINOXE:default        |
       | remote_code_space | internal              |
    And a Facility exists with the following attributes:
      | Codes[internal] | NINOXE:Facility:1:LOC |
    And a Subscription exist with the following attributes:
      | Kind              | FacilityMonitoringBroadcast                   |
      | SubscriberRef     | Subscriber                                    |
      | ExternalId        | subscription-1                                |
      | ReferenceArray[0] | Facility, "internal": "NINOXE:Facility:1:LOC" |
    When the Facility "internal:NINOXE:Facility:1:LOC" is edited with the following attributes:
      | Status | available  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:NotifyFacilityMonitoring xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:10.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>test</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-4-00c04fd430c8</siri:ResponseMessageIdentifier>
            </ServiceDeliveryInfo>
            <Notification>
              <siri:FacilityMonitoringDelivery version='2.0:FR-IDF-2.4'>
                <siri:ResponseTimestamp>2017-01-01T12:00:10.000Z</siri:ResponseTimestamp>
                <siri:SubscriberRef>Subscriber</siri:SubscriberRef>
                <siri:SubscriptionRef>subscription-1</siri:SubscriptionRef>
                <siri:Status>true</siri:Status>
                <siri:FacilityCondition>
                  <siri:FacilityRef>NINOXE:Facility:1:LOC</siri:FacilityRef>
                  <siri:FacilityStatus>
                    <siri:Status>available</siri:Status>
                  </siri:FacilityStatus>
                </siri:FacilityCondition>
              </siri:FacilityMonitoringDelivery>
            </Notification>
            <SiriExtension/>
          </sw:NotifyFacilityMonitoring>
        </S:Body>
      </S:Envelope>
      """
    Then an audit event should exist with these attributes:
      | Type                    | NotifyFacilityMonitoring |
      | Direction               | sent                     |
      | Protocol                | siri                     |
      | Partner                 | test                     |
      | Status                  | OK                       |
      | SubscriptionIdentifiers | ["subscription-1"]       |

  @ARA-1755
  Scenario: Handle a SIRI FacilityMonitoring subscription for all existing facilities in a referential having same remote_code_space
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-facility-monitoring-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
   And a Facility exists with the following attributes:
      | Codes[internal] | NINOXE:Facility:1:LOC |
   And a Facility exists with the following attributes:
      | Codes[internal] | NINOXE:Facility:2:LOC |
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
            <siri:FacilityMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>subscription-1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:FacilityMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
              </siri:FacilityMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
            </siri:FacilityMonitoringSubscriptionRequest>
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
    Then Subscriptions exist with the following resources:
      | internal | NINOXE:Facility:1:LOC |
      | internal | NINOXE:Facility:2:LOC |
