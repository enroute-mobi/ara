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
      | ResponseIdentifier | RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC |
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
