Feature: Support SIRI GeneralMessage for Situation

  Background:
      Given a Referential "test" is created

  Scenario: 3797 - Ignore situations associatd to the Commercial channel
    Given a Situation exists with the following attributes:
      | ObjectIDs               | "internal" : "1"                                 |
      | RecordedAt              | 2017-01-01T03:30:06+02:00                        |
      | Version                 | 1                                                |
      | Channel                 | Commercial                                       |
      | ValidUntil              | 2017-01-01T20:30:06+02:00                        |
      | Messages[0]#MessageType | longMessage                                      |
      | Messages[0]#MessageText | Les situations commercials doivent être ignorées |
      | References[0]           | LineRef:{"internal":"NINOXE:Line:3:LOC"}         |
    And a Situation exists with the following attributes:
      | ObjectIDs               | "internal" : "2"                         |
      | RecordedAt              | 2017-01-01T03:30:06+02:00                |
      | Version                 | 1                                        |
      | Channel                 | Perturbations                            |
      | ValidUntil              | 2017-01-01T20:30:06+02:00                |
      | Messages[0]#MessageType | longMessage                              |
      | Messages[0]#MessageText | Les autres non                           |
      | References[0]           | LineRef:{"internal":"NINOXE:Line:3:LOC"} |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential     | TestPartner |
      | remote_objectid_kind | internal    |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ns7:GetGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri" xmlns:siri="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://wsdl.siri.org.uk/siri" xmlns:ns7="http://wsdl.siri.org.uk">
          <ServiceRequestInfo>
            <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
            <ns2:RequestorRef>TestPartner</ns2:RequestorRef>
            <ns2:MessageIdentifier>GeneralMessage:Test:0</ns2:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
            <ns2:MessageIdentifier>GeneralMessage:Test:0</ns2:MessageIdentifier>
            <ns2:Extensions>
              <ns6:IDFGeneralMessageRequestFilter>
              </ns6:IDFGeneralMessageRequestFilter>
            </ns2:Extensions>
          </Request>
          <RequestExtension/>
        </ns7:GetGeneralMessage>
      </S:Body>
    </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetGeneralMessageResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Edwig</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4" xmlns:stif="http://wsdl.siri.org.uk/siri">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:GeneralMessage formatRef="STIF-IDF">
                  <siri:RecordedAtTime>2017-01-01T03:30:06.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>RATPDev:Item::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ItemIdentifier>
                  <siri:InfoMessageIdentifier>2</siri:InfoMessageIdentifier>
                  <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                  <siri:InfoChannelRef>Perturbations</siri:InfoChannelRef>
                  <siri:ValidUntilTime>2017-01-01T20:30:06.000+02:00</siri:ValidUntilTime>
                  <siri:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
          				xsi:type="stif:IDFGeneralMessageStructure">
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText>Les autres non</MessageText>
                    </Message>
                  </siri:Content>
                </siri:GeneralMessage>
              </siri:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 3008 - Handle a SIRI GetGeneralMessage request
    Given a Situation exists with the following attributes:
      | ObjectIDs               | "external" : "test"                                                        |
      | RecordedAt              | 2017-01-01T03:30:06+02:00                                                  |
      | Version                 | 1                                                                          |
      | Channel                 | Perturbation                                                               |
      | ValidUntil              | 2017-01-01T20:30:06+02:00                                                  |
      | Messages[0]#MessageType | longMessage                                                                |
      | Messages[0]#MessageText | La nouvelle carte d'abonnement est disponible au points de vente du réseau |
      | References[0]           | LineRef:{"external":"NINOXE:Line:3:LOC"}                                   |
    And a Line exists with the following attributes:
      | ObjectIDs | "external": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential     | NINOXE:default |
      | remote_objectid_kind | external       |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ns7:GetGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri" xmlns:siri="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://wsdl.siri.org.uk/siri" xmlns:ns7="http://wsdl.siri.org.uk">
          <ServiceRequestInfo>
            <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
            <ns2:RequestorRef>NINOXE:default</ns2:RequestorRef>
            <ns2:MessageIdentifier>GeneralMessage:Test:0</ns2:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
            <ns2:MessageIdentifier>GeneralMessage:Test:0</ns2:MessageIdentifier>
            <ns2:Extensions>
              <ns6:IDFGeneralMessageRequestFilter>
              </ns6:IDFGeneralMessageRequestFilter>
            </ns2:Extensions>
          </Request>
          <RequestExtension/>
        </ns7:GetGeneralMessage>
      </S:Body>
    </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetGeneralMessageResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Edwig</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4" xmlns:stif="http://wsdl.siri.org.uk/siri">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:GeneralMessage formatRef="STIF-IDF">
                  <siri:RecordedAtTime>2017-01-01T03:30:06.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>RATPDev:Item::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:ItemIdentifier>
                  <siri:InfoMessageIdentifier>test</siri:InfoMessageIdentifier>
                  <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                  <siri:InfoChannelRef>Perturbation</siri:InfoChannelRef>
                  <siri:ValidUntilTime>2017-01-01T20:30:06.000+02:00</siri:ValidUntilTime>
                  <siri:Content xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xsi:type='stif:IDFGeneralMessageStructure'>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText>La nouvelle carte d'abonnement est disponible au points de vente du réseau</MessageText>
                    </Message>
                  </siri:Content>
                </siri:GeneralMessage>
              </siri:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 3032 - Handle a GeneralMessage response (ServiceDelivery)
    Given a SIRI server waits GeneralMessageRequest request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetGeneralMessageResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>
              2017-03-29T16:48:00.993+02:00</siri:ResponseTimestamp>
              <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
              <siri:Address>http://appli.chouette.mobi/siri_france/siri</siri:Address>
              <siri:ResponseMessageIdentifier>b28e8207-f030-4932-966c-3e6099fad4ef</siri:ResponseMessageIdentifier>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-03-29T16:48:00.039+02:00</siri:ResponseTimestamp>
                <siri:Status>true</siri:Status>
                <siri:GeneralMessage formatRef="FRANCE">
                  <siri:RecordedAtTime>
                  2017-03-29T03:30:06.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>3477</siri:ItemIdentifier>
                  <siri:InfoMessageIdentifier>
                  NINOXE:GeneralMessage:27_1</siri:InfoMessageIdentifier>
                  <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                  <siri:InfoChannelRef>Commercial</siri:InfoChannelRef>
                  <siri:ValidUntilTime>
                  2017-03-29T20:30:06.000+02:00</siri:ValidUntilTime>
                  <siri:Content>
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText xml:lang="NL">La nouvelle carte d'abonnement est disponible au points de vente du réseau</MessageText>
                    </Message>
                  </siri:Content>
                </siri:GeneralMessage>
              </siri:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension />
          </sw:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ineo                  |
      | remote_objectid_kind | internal              |
    And a Line exists with the following attributes:
      | Name                   | Test              |
      | ObjectIDs              | "internal":"1234" |
      | CollectGeneralMessages | true              |
    And a minute has passed
    When a minute has passed
    And the SIRI server has received a GeneralMessage request
    Then one Situation has the following attributes:
      | ObjectIDs               | "internal" : "NINOXE:GeneralMessage:27_1"                                  |
      | RecordedAt              | 2017-03-29T03:30:06+02:00                                                  |
      | Version                 | 1                                                                          |
      | Channel                 | Commercial                                                                 |
      | ProducerRef             | NINOXE:default                                                             |
      | ValidUntil              | 2017-03-29T20:30:06+02:00                                                  |
      | Messages[0]#MessageType | longMessage                                                                |
      | Messages[0]#MessageText | La nouvelle carte d'abonnement est disponible au points de vente du réseau |

  Scenario: 3864 - Modification of a Situation after a GetGeneralMessageResponse
    Given a SIRI server waits GeneralMessageRequest request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetGeneralMessageResponse xmlns:siri="http://www.siri.org.uk/siri"
          xmlns:ns4="http://www.ifopt.org.uk/acsb"
          xmlns:ns5="http://www.ifopt.org.uk/ifopt"
          xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns7="http://scma/siri"
          xmlns:sw="http://wsdl.siri.org.uk"
          xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Edwig</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:Status>true</siri:Status>
                <siri:GeneralMessage>
                  <siri:formatRef>STIF-IDF</siri:formatRef>
                  <siri:RecordedAtTime>2017-01-01T03:30:06.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>RATPDev:Item::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:ItemIdentifier>
                  <siri:InfoMessageIdentifier>NINOXE:GeneralMessage:27_1</siri:InfoMessageIdentifier>
                  <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                  <siri:InfoChannelRef>Perturbation</siri:InfoChannelRef>
                  <siri:ValidUntilTime>2017-01-07T20:30:06.000+02:00</siri:ValidUntilTime>
                  <siri:Content>
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText>La nouvelle carte d'abonnement est disponible au points de vente du réseau</MessageText>
                    </Message>
                  </siri:Content>
                </siri:GeneralMessage>
              </siri:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension />
          </sw:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """
    And a Situation exists with the following attributes:
      | ObjectIDs               | "external" : "NINOXE:GeneralMessage:27_1"                                  |
      | RecordedAt              | 2017-01-01T03:30:06+02:00                                                  |
      | Version                 | 1                                                                          |
      | Channel                 | Perturbation                                                               |
      | ValidUntil              | 2017-01-01T20:30:06+02:00                                                  |
      | Messages[0]#MessageType | longMessage                                                                |
      | Messages[0]#MessageText | La nouvelle carte d'abonnement est disponible au points de vente du réseau |
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ineo                  |
      | remote_objectid_kind | internal              |
    And a Line exists with the following attributes:
      | Name                   | Test              |
      | ObjectIDs              | "internal":"1234" |
      | CollectGeneralMessages | true              |
    And a minute has passed
    When a minute has passed
    And the SIRI server has received a GeneralMessage request
    Then a Situation exists with the following attributes:
      | ObjectIDs               | "external" : "NINOXE:GeneralMessage:27_1"                                  |
      | RecordedAt              | 2017-01-01T03:30:06+02:00                                                  |
      | Version                 | 1                                                                          |
      | Channel                 | Perturbation                                                               |
      | ValidUntil              | 2017-01-07T20:30:06+02:00                                                  |
      | Messages[0]#MessageType | longMessage                                                                |
      | Messages[0]#MessageText | La nouvelle carte d'abonnement est disponible au points de vente du réseau |

  Scenario: 3882 - GeneralMessageResponse empty with an expired Situation
    Given a Situation exists with the following attributes:
      | ObjectIDs               | "external" : "test"                                                        |
      | RecordedAt              | 2017-01-01T03:30:06+02:00                                                  |
      | Version                 | 1                                                                          |
      | Channel                 | Perturbation                                                               |
      | ValidUntil              | 2017-01-01T12:01:00+02:00                                                  |
      | Messages[0]#MessageType | longMessage                                                                |
      | Messages[0]#MessageText | La nouvelle carte d'abonnement est disponible au points de vente du réseau |
    And a Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential     | NINOXE:default |
      | remote_objectid_kind | external       |
    And a Line exists with the following attributes:
      | Name                   | Test              |
      | ObjectIDs              | "internal":"1234" |
      | CollectGeneralMessages | true              |
    And a minute has passed
    And a minute has passed
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <SOAP-ENV:Header/>
      <S:Body>
        <ns7:GetGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri" xmlns:siri="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://wsdl.siri.org.uk/siri" xmlns:ns7="http://wsdl.siri.org.uk">
          <ServiceRequestInfo>
            <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
            <ns2:RequestorRef>NINOXE:default</ns2:RequestorRef>
            <ns2:MessageIdentifier>GeneralMessage:Test:0</ns2:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
            <ns2:MessageIdentifier>GeneralMessage:Test:0</ns2:MessageIdentifier>
            <ns2:Extensions>
              <ns6:IDFGeneralMessageRequestFilter>
              </ns6:IDFGeneralMessageRequestFilter>
            </ns2:Extensions>
          </Request>
          <RequestExtension/>
        </ns7:GetGeneralMessage>
      </S:Body>
    </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetGeneralMessageResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Edwig</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4" xmlns:stif="http://wsdl.siri.org.uk/siri">
                <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
              </siri:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """

  @wip
  Scenario: Manage a Request without filter
    Given a SIRI server waits GetGeneralMessage request on "http://localhost:8090" to respond with
    """
    """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_objectid_kind | internal              |
    And 30 seconds have passed
    And 30 seconds have passed
    And the SIRI server has received a GetGeneralMessage request
    Then the SIRI server should receive this response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
  <sw:GetGeneralMessage xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri' xmlns:sws='http://wsdl.siri.org.uk/siri'>
    <ServiceRequestInfo>
      <siri:RequestTimestamp>2017-01-01T12:01:00.000Z</siri:RequestTimestamp>
      <siri:RequestorRef>test</siri:RequestorRef>
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:MessageIdentifier>
    </ServiceRequestInfo>
    <Request version='2.0:FR-IDF-2.4'>
      <siri:RequestTimestamp>2017-01-01T12:01:00.000Z</siri:RequestTimestamp>
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:MessageIdentifier>
      <siri:Extensions>
        <sws:IDFGeneralMessageRequestFilter>
        </sws:IDFGeneralMessageRequestFilter>
      </siri:Extensions>
    </Request>
    <RequestExtension/>
  </sw:GetGeneralMessage>
</S:Body>
</S:Envelope>
      """

  Scenario: Manage a Request with a Line filter
    Given a SIRI server waits GetGeneralMessage request on "http://localhost:8090" to respond with
    """
    """
      And a Partner "test" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
        | remote_url                      | http://localhost:8090 |
        | remote_credential               | test                  |
        | local_credential                | NINOXE:default        |
        | remote_objectid_kind            | internal              |
        | collect.filter_general_messages | true                  |
        | collect.include_lines           | 1234                  |
      And 30 seconds have passed
      And a Line exists with the following attributes:
        | Name                   | Test              |
        | ObjectIDs              | "internal":"1234" |
        | CollectGeneralMessages | true              |
      And 10 seconds have passed
      And the SIRI server has received a GetGeneralMessage request
    Then the SIRI server should receive this response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
  <sw:GetGeneralMessage xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri' xmlns:sws='http://wsdl.siri.org.uk/siri'>
    <ServiceRequestInfo>
      <siri:RequestTimestamp>2017-01-01T12:00:40.000Z</siri:RequestTimestamp>
      <siri:RequestorRef>test</siri:RequestorRef>
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:MessageIdentifier>
    </ServiceRequestInfo>
    <Request version='2.0:FR-IDF-2.4'>
      <siri:RequestTimestamp>2017-01-01T12:00:40.000Z</siri:RequestTimestamp>
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:MessageIdentifier>
      <siri:Extensions>
        <sws:IDFGeneralMessageRequestFilter>
          <siri:LineRef>1234</siri:LineRef>
        </sws:IDFGeneralMessageRequestFilter>
      </siri:Extensions>
    </Request>
    <RequestExtension/>
  </sw:GetGeneralMessage>
</S:Body>
</S:Envelope>
      """

  Scenario: Manage a Request with a StopArea filter
    Given a SIRI server waits GetGeneralMessage request on "http://localhost:8090" to respond with
    """
    """
      And a Partner "test" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
        | remote_url                      | http://localhost:8090 |
        | remote_credential               | test                  |
        | local_credential                | NINOXE:default        |
        | remote_objectid_kind            | internal              |
        | collect.filter_general_messages | true                  |
      And 30 seconds have passed
      And a StopArea exists with the following attributes:
        | Name                   | Test              |
        | ObjectIDs              | "internal":"1234" |
        | CollectGeneralMessages | true              |
      And 10 seconds have passed
      And the SIRI server has received a GetGeneralMessage request
    Then the SIRI server should receive this response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
  <sw:GetGeneralMessage xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri' xmlns:sws='http://wsdl.siri.org.uk/siri'>
    <ServiceRequestInfo>
      <siri:RequestTimestamp>2017-01-01T12:00:40.000Z</siri:RequestTimestamp>
      <siri:RequestorRef>test</siri:RequestorRef>
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:MessageIdentifier>
    </ServiceRequestInfo>
    <Request version='2.0:FR-IDF-2.4'>
      <siri:RequestTimestamp>2017-01-01T12:00:40.000Z</siri:RequestTimestamp>
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:MessageIdentifier>
      <siri:Extensions>
        <sws:IDFGeneralMessageRequestFilter>
          <siri:StopPointRef>1234</siri:StopPointRef>
        </sws:IDFGeneralMessageRequestFilter>
      </siri:Extensions>
    </Request>
    <RequestExtension/>
  </sw:GetGeneralMessage>
</S:Body>
</S:Envelope>
      """
