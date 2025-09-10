Feature: Manages Partner templates

  Background:
    Given a Referential "test" is created

  Scenario: Create a Partner template
    Given a Partner template "test" exists with the following attributes:
      | CredentialType   | FormatMatching                                                                   |
      | LocalCredential  | %{value}                                                                         |
      | RemoteCredential | %{value}                                                                         |
      | ConnectorTypes   | ["siri-check-status-client","siri-estimated-timetable-subscription-broadcaster"] |
      | Settings         | {"remote_code_space":"internal"}                                                 |
    Then one Partner template has the following attributes:
      | Slug | test |

  Scenario: Update a Partner template
    Given a Partner template "test" exists with the following attributes:
      | CredentialType   | FormatMatching                                                                   |
      | LocalCredential  | %{value}                                                                         |
      | RemoteCredential | %{value}                                                                         |
      | ConnectorTypes   | ["siri-check-status-client","siri-estimated-timetable-subscription-broadcaster"] |
      | Settings         | {"remote_code_space":"internal"}                                                 |
    When the Partner template "test" is updated with the following attributes:
      | LocalCredential | test:%{value} |
    Then one Partner template has the following attributes:
      | Slug            | test          |
      | LocalCredential | test:%{value} |

  Scenario: Handle the creation of a template via PartnerTemplate
    Given a Partner template "test" exists with the following attributes:
      | CredentialType   | FormatMatching               |
      | LocalCredential  | test:local:%{value}          |
      | RemoteCredential | test:remote:%{value}         |
      | ConnectorTypes   | ["siri-check-status-server"] |
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
        <siri:RequestorRef>test:local:1234</siri:RequestorRef>
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
        <siri:ProducerRef>test:remote:1234</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-3-00c04fd430c8</siri:ResponseMessageIdentifier>
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
