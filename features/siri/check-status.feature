Feature: Support SIRI CheckStatus
  Background:
    Given a Referential "test" is created

  @ARA-1023
  Scenario: Use OAuth 2.0 token to perform a SIRI CheckStatus request
    Given a SIRI server waits CheckStatus request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:CheckStatusResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <CheckStatusAnswerInfo>
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:ProducerRef>Ara</siri:ProducerRef>
            <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            <siri:RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:RequestMessageRef>
          </CheckStatusAnswerInfo>
          <Answer>
            <siri:Status>true</siri:Status>
            <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
          </Answer>
          <AnswerExtension/>
          </sw:CheckStatusResponse>
        </S:Body>
      </S:Envelope>
      """
    And an OAuth server waits request on "http://localhost:8091" and accepts client "ara" with
      | client_secret | oauth-secret            |
      | access_token  | oauth-access-token-test |
    And a Partner "test_partner" exists with connectors [siri-check-status-client] and the following settings:
      | remote_credential                         | Ara                   |
      | remote_authentication.oauth.client_id     | ara                   |
      | remote_authentication.oauth.client_secret | oauth-secret          |
      | remote_authentication.oauth.token_url     | http://localhost:8091 |
      | remote_url                                | http://localhost:8090 |
    When 30 seconds have passed
    Then the SIRI server should have received a CheckStatus request with the payload:
          """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <sw:CheckStatus xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
        <Request>
          <siri:RequestTimestamp>2017-01-01T12:00:30.000Z</siri:RequestTimestamp>
          <siri:RequestorRef>Ara</siri:RequestorRef>
          <siri:MessageIdentifier>6ba7b814-9dad-11d1-2-00c04fd430c8</siri:MessageIdentifier>
        </Request>
        <RequestExtension/>
        </sw:CheckStatus>
      </S:Body>
    </S:Envelope>
          """

  @ARA-1424
  Scenario: Use OAuth 2.0 token with Scopes to perform a SIRI CheckStatus request
    Given a SIRI server waits CheckStatus request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:CheckStatusResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <CheckStatusAnswerInfo>
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:ProducerRef>Ara</siri:ProducerRef>
            <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            <siri:RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:RequestMessageRef>
          </CheckStatusAnswerInfo>
          <Answer>
            <siri:Status>true</siri:Status>
            <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
          </Answer>
          <AnswerExtension/>
          </sw:CheckStatusResponse>
        </S:Body>
      </S:Envelope>
      """
    And an OAuth server waits request on "http://localhost:8091" and accepts client "ara" with
      | client_secret | oauth-secret            |
      | access_token  | oauth-access-token-test |
      | scopes        | scope-test              |
    And a Partner "test_partner" exists with connectors [siri-check-status-client] and the following settings:
      | remote_credential                         | Ara                   |
      | remote_authentication.oauth.client_id     | ara                   |
      | remote_authentication.oauth.client_secret | oauth-secret          |
      | remote_authentication.oauth.token_url     | http://localhost:8091 |
      | remote_url                                | http://localhost:8090 |
      | remote_authentication.oauth.scopes        | scope-test            |
    And show me ara partners
    When 30 seconds have passed
    Then the SIRI server should have received a CheckStatus request with the payload:
          """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <sw:CheckStatus xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
        <Request>
          <siri:RequestTimestamp>2017-01-01T12:00:30.000Z</siri:RequestTimestamp>
          <siri:RequestorRef>Ara</siri:RequestorRef>
          <siri:MessageIdentifier>6ba7b814-9dad-11d1-2-00c04fd430c8</siri:MessageIdentifier>
        </Request>
        <RequestExtension/>
        </sw:CheckStatus>
      </S:Body>
    </S:Envelope>
          """

  Scenario: 2460 - Handle a SIRI Checkstatus SOAP request
    Given a SIRI Partner "test" exists with connectors [siri-check-status-server] and the following settings:
      | local_credential | test |
    When I send this SIRI request to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test</siri:RequestorRef>
        <siri:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatusResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <CheckStatusAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:RequestMessageRef>
      </CheckStatusAnswerInfo>
      <Answer>
        <siri:Status>true</siri:Status>
        <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
      </Answer>
      <AnswerExtension/>
    </sw:CheckStatusResponse>
  </S:Body>
</S:Envelope>
      """


  @ARA-1025
  Scenario: Handle a SIRI Checkstatus request with raw envelope
    Given a SIRI Partner "test" exists with connectors [siri-check-status-server] and the following settings:
      | local_credential | test |
      | siri.envelope    | raw  |
    When I send this SIRI request to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns="http://www.siri.org.uk/siri" version="2.0">
  <CheckStatusRequest>
      <RequestTimestamp>2017-01-01T12:00:00.000Z</RequestTimestamp>
      <RequestorRef>test</RequestorRef>
      <MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</MessageIdentifier>
  </CheckStatusRequest>
</Siri>
"""
   Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='UTF-8'?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
<CheckStatusResponse>
    <ResponseTimestamp>2017-01-01T12:00:00.000Z</ResponseTimestamp>
    <ProducerRef>Ara</ProducerRef>
    <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ResponseMessageIdentifier>
    <RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</RequestMessageRef>
    <Status>true</Status>
    <ServiceStartedTime>2017-01-01T12:00:00.000Z</ServiceStartedTime>
</CheckStatusResponse>
</Siri>
"""

  @ARA-1055
  Scenario: Handle a SIRI Checkstatus request with raw envelope without special characters between tags
    Given a SIRI Partner "test" exists with connectors [siri-check-status-server] and the following settings:
      | local_credential | test |
      | siri.envelope    | raw  |
    When I send this SIRI request to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?><Siri xmlns="http://www.siri.org.uk/siri" version="2.0"><CheckStatusRequest><RequestTimestamp>2017-01-01T12:00:00.000Z</RequestTimestamp><RequestorRef>test</RequestorRef><MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</MessageIdentifier></CheckStatusRequest></Siri>
"""
   Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='UTF-8'?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
<CheckStatusResponse>
    <ResponseTimestamp>2017-01-01T12:00:00.000Z</ResponseTimestamp>
    <ProducerRef>Ara</ProducerRef>
    <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ResponseMessageIdentifier>
    <RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</RequestMessageRef>
    <Status>true</Status>
    <ServiceStartedTime>2017-01-01T12:00:00.000Z</ServiceStartedTime>
</CheckStatusResponse>
</Siri>
"""

  @ARA-1025
  Scenario: Send SIRI CheckStatus request with SOAP
    Given a SIRI server waits CheckStatus request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatusResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <CheckStatusAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>d3f94aa2-7b76-449b-aa18-50caf78f9dc7</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:RequestMessageRef>
      </CheckStatusAnswerInfo>
      <Answer>
        <siri:Status>true</siri:Status>
        <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
      </Answer>
      <AnswerExtension/>
    </sw:CheckStatusResponse>

  </S:Body>
</S:Envelope>
      """

    And a Partner "test_partner" exists with connectors [siri-check-status-client] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | Ara                   |

    When 30 seconds have passed
    Then the SIRI server should have received a CheckStatus request with the payload:
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
    <Request>
      <siri:RequestTimestamp>2017-01-01T12:00:30.000Z</siri:RequestTimestamp>
      <siri:RequestorRef>Ara</siri:RequestorRef>
      <siri:MessageIdentifier>6ba7b814-9dad-11d1-2-00c04fd430c8</siri:MessageIdentifier>
    </Request>
    <RequestExtension/>
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then the Partner "test_partner" in the Referential "test" has the operational status up

  @ARA-1025
  Scenario: Send SIRI CheckStatus request with raw envelope
    Given a raw SIRI server waits CheckStatus request on "http://localhost:8090" to respond with
      """
<?xml version="1.0" encoding="UTF-8"?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
<CheckStatusResponse>
    <ResponseTimestamp>2017-01-01T12:00:00.000Z</ResponseTimestamp>
    <ProducerRef>test</ProducerRef>
    <ResponseMessageIdentifier>6ba7b814-9dad-11d1-2-00c04fd430c8:</ResponseMessageIdentifier>
    <RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC</RequestMessageRef>
    <Status>true</Status>
    <ServiceStartedTime>2017-01-01T12:00:00.000Z</ServiceStartedTime>
</CheckStatusResponse>
</Siri>
      """
    And a Partner "test_partner" exists with connectors [siri-check-status-client] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | Ara                   |
      | siri.envelope     | raw                   |
    When 30 seconds have passed
    Then the SIRI server should have received a CheckStatus request with the payload:
      """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
<CheckStatusRequest>
      <RequestTimestamp>2017-01-01T12:00:30.000Z</RequestTimestamp>
      <RequestorRef>Ara</RequestorRef>
      <MessageIdentifier>6ba7b814-9dad-11d1-2-00c04fd430c8</MessageIdentifier>
  </CheckStatusRequest>
</Siri>
      """
    Then the Partner "test_partner" in the Referential "test" has the operational status up

  Scenario: Handle a SIRI Checkstatus request with invalid RequestorRef
    When I send this SIRI request
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>invalid</siri:RequestorRef>
        <siri:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <S:Fault>
      <faultcode>S:UnknownCredential</faultcode>
      <faultstring>RequestorRef Unknown 'invalid'</faultstring>
    </S:Fault>
  </S:Body>
</S:Envelope>
      """

  Scenario: Manage several local credentials on siri-check-status-server connecter
    Given a SIRI Partner "test" exists with connectors [siri-check-status-server] and the following settings:
      | local_credentials | test1,test2 |
    When I send this SIRI request to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test1</siri:RequestorRef>
        <siri:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatusResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <CheckStatusAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:RequestMessageRef>
      </CheckStatusAnswerInfo>
      <Answer>
        <siri:Status>true</siri:Status>
        <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
      </Answer>
      <AnswerExtension/>
    </sw:CheckStatusResponse>
  </S:Body>
</S:Envelope>
      """
    When I send this SIRI request to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test2</siri:RequestorRef>
        <siri:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatusResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <CheckStatusAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:RequestMessageRef>
      </CheckStatusAnswerInfo>
      <Answer>
        <siri:Status>true</siri:Status>
        <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
      </Answer>
      <AnswerExtension/>
    </sw:CheckStatusResponse>
  </S:Body>
</S:Envelope>
      """

  @ARA-1181
  Scenario: Handle a SIRI Checkstatus SOAP request with Rate Limit
    Given a SIRI Partner "test" exists with connectors [siri-check-status-server] and the following settings:
      | local_credential  | test |
      | rate_limit_per_ip | 2    |
    When I send this SIRI request 2 times to the Referential "test"
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test</siri:RequestorRef>
        <siri:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatusResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <CheckStatusAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:RequestMessageRef>
      </CheckStatusAnswerInfo>
      <Answer>
        <siri:Status>true</siri:Status>
        <siri:ServiceStartedTime>2017-01-01T12:00:00.000Z</siri:ServiceStartedTime>
      </Answer>
      <AnswerExtension/>
    </sw:CheckStatusResponse>
  </S:Body>
</S:Envelope>
      """
    When I send this SIRI request to the Referential "test" expecting an error
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri"
    xmlns:ns3="http://www.ifopt.org.uk/acsb"
    xmlns:ns4="http://www.ifopt.org.uk/ifopt"
    xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns6="http://scma/siri"
    xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test</siri:RequestorRef>
        <siri:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:CheckStatus>
  </S:Body>
</S:Envelope>
      """
    Then I should receive a HTTP error 429
