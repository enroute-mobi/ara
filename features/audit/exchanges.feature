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
      | Timestamp          | 2017-01-01T12:00:00Z  |
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
      | Timestamp          | 2017-01-01T12:00:00Z                            |
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
    And an audit event should exist with these attributes:
      | Protocol           | siri                  |
      | Direction          | sent                  |
      | Status             | OK                    |
      | Partner            | test                  |
      | Timestamp          | 2017-01-01T12:00:00Z  |
      | RequestIdentifier  | RATPDev:Message::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC |
      | ResponseIdentifier | c464f588-5128-46c8-ac3f-8b8a465692ab |
      | Timestamp          | 2017-01-01T12:01:00Z  |
      | ProcessingTime     | 0                     |
