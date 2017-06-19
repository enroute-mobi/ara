Feature: Support SIRI GeneralMessage

  Background:
      Given a Referential "test" is created

  # Should be fixed by 3800 and 3801
  @wip
  Scenario: 3008 - Handle a SIRI GetGeneralMessage request
    Given a Situation exists with the following attributes:
      | ObjectIDs               | "external" : "Edwig:InfoMessage::test:LOC"                                 |
      | RecordedAt              | 2017-01-01T03:30:06+02:00                                                  |
      | Version                 | 1                                                                          |
      | Channel                 | Perturbation                                                               |
      | ValidUntil              | 2017-01-01T20:30:06+02:00                                                  |
      | Messages[0]#MessageType | longMessage                                                                |
      | Messages[0]#MessageText | La nouvelle carte d'abonnement est disponible au points de vente du réseau |
    And a Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential     | NINOXE:default |
      | remote_objectid_kind | external       |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <SOAP-ENV:Header/>
      <S:Body>
        <ns7:GetGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns3="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://wsdl.siri.org.uk/siri" xmlns:ns7="http://wsdl.siri.org.uk">
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
          <ns8:GetGeneralMessageResponse xmlns:ns3="http://www.siri.org.uk/siri"
          xmlns:ns4="http://www.ifopt.org.uk/acsb"
          xmlns:ns5="http://www.ifopt.org.uk/ifopt"
          xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns7="http://scma/siri"
          xmlns:ns8="http://wsdl.siri.org.uk"
          xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
              <ns3:ProducerRef>Edwig</ns3:ProducerRef>
              <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>GeneralMessage:Test:0</ns3:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
                <ns3:Status>true</ns3:Status>
                <ns3:GeneralMessage>
                  <ns3:formatRef>STIF-IDF</ns3:formatRef>
                  <ns3:RecordedAtTime>2017-01-01T03:30:06.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>RATPDev:Item::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</ns3:ItemIdentifier> <!-- Shoud be fixed by 3800 -->
                  <ns3:InfoMessageIdentifier>Edwig:InfoMessage::test:LOC</ns3:InfoMessageIdentifier> <!-- Should be fixed by 3801 -->
                  <ns3:InfoMessageVersion>1</ns3:InfoMessageVersion>
                  <ns3:InfoChannelRef>Perturbation</ns3:InfoChannelRef>
                  <ns3:ValidUntilTime>2017-01-01T20:30:06.000+02:00</ns3:ValidUntilTime>
                  <ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns9:IDFLineSectionStructure">
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText>La nouvelle carte d'abonnement est disponible au points de vente du réseau</MessageText>
                    </Message>
                  </ns3:Content>
                </ns3:GeneralMessage>
              </ns3:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension />
          </ns8:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """

  @wip
  Scenario: 3009 - Handle a SIRI ServiceDelivery after GM Request to a SIRI server
    Given a SIRI server waits GeneralMessageRequest request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns8:GetGeneralMessageResponse xmlns:ns3="http://www.siri.org.uk/siri"
          xmlns:ns4="http://www.ifopt.org.uk/acsb"
          xmlns:ns5="http://www.ifopt.org.uk/ifopt"
          xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns7="http://scma/siri"
          xmlns:ns8="http://wsdl.siri.org.uk"
          xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>
              2017-03-29T16:48:00.993+02:00</ns3:ResponseTimestamp>
              <ns3:ProducerRef>NINOXE:default</ns3:ProducerRef>
              <ns3:Address>
              http://appli.chouette.mobi/siri_france/siri</ns3:Address>
              <ns3:ResponseMessageIdentifier>
              b28e8207-f030-4932-966c-3e6099fad4ef</ns3:ResponseMessageIdentifier>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>
                2017-03-29T16:48:00.039+02:00</ns3:ResponseTimestamp>
                <ns3:Status>true</ns3:Status>
                <ns3:GeneralMessage formatRef="FRANCE">
                  <ns3:RecordedAtTime>
                  2017-03-29T03:30:06.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>3477</ns3:ItemIdentifier>
                  <ns3:InfoMessageIdentifier>
                  NINOXE:GeneralMessage:27_1</ns3:InfoMessageIdentifier>
                  <ns3:InfoMessageVersion>1</ns3:InfoMessageVersion>
                  <ns3:InfoChannelRef>Commercial</ns3:InfoChannelRef>
                  <ns3:ValidUntilTime>
                  2017-03-29T20:30:06.000+02:00</ns3:ValidUntilTime>
                  <ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                  xsi:type="ns9:IDFGeneralMessageStructure">
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText xml:lang="NL">La nouvelle carte
                      d'abonnement est disponible au points de vente du
                      réseau</MessageText>
                    </Message>
                  </ns3:Content>
                </ns3:GeneralMessage>
                <ns3:GeneralMessage formatRef="FRANCE">
                  <ns3:RecordedAtTime>
                  2017-03-29T03:30:06.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>3471</ns3:ItemIdentifier>
                  <ns3:InfoMessageIdentifier>
                  NINOXE:GeneralMessage:21_1</ns3:InfoMessageIdentifier>
                  <ns3:InfoMessageVersion>1</ns3:InfoMessageVersion>
                  <ns3:InfoChannelRef>Commercial</ns3:InfoChannelRef>
                  <ns3:ValidUntilTime>
                  2017-03-29T22:30:06.000+02:00</ns3:ValidUntilTime>
                  <ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                  xsi:type="ns9:IDFGeneralMessageStructure">
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText xml:lang="NL">Les nouveaux tarifs sont
                      consultable sur le site internet</MessageText>
                    </Message>
                  </ns3:Content>
                </ns3:GeneralMessage>
              </ns3:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension />
          </ns8:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
    | remote_url           | http://localhost:8090 |
    | remote_credential    | ratpdev               |
    | remote_objectid_kind | internal              |
    And a minute has passed
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <SOAP-ENV:Header/>
      <S:Body>
        <ns7:GetGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns3="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://wsdl.siri.org.uk/siri" xmlns:ns7="http://wsdl.siri.org.uk">
          <ServiceRequestInfo>
            <ns2:RequestTimestamp>2017-03-29T16:49:00.311Z</ns2:RequestTimestamp>
            <ns2:RequestorRef>SQYBUS</ns2:RequestorRef>
            <ns2:MessageIdentifier>GeneralMessage:Test:0</ns2:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <ns2:RequestTimestamp>2017-03-29T16:49:00.311Z</ns2:RequestTimestamp>
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
          <ns8:GetGeneralMessageResponse xmlns:ns3="http://www.siri.org.uk/siri"
          xmlns:ns4="http://www.ifopt.org.uk/acsb"
          xmlns:ns5="http://www.ifopt.org.uk/ifopt"
          xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns7="http://scma/siri"
          xmlns:ns8="http://wsdl.siri.org.uk"
          xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>
              2017-03-29T16:49:30.993+02:00</ns3:ResponseTimestamp>
              <ns3:ProducerRef>RATPDev</ns3:ProducerRef>
              <ns3:Address>
              http://appli.chouette.mobi/siri_france/siri</ns3:Address>
              <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>RATPDev:Message::9dad:LOC</ns3:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>
                2017-03-29T16:48:30.039+02:00</ns3:ResponseTimestamp>
                <ns3:Status>true</ns3:Status>
                <ns3:GeneralMessage formatRef="FRANCE">
                  <ns3:RecordedAtTime>
                  2017-03-29T03:30:06.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>3477</ns3:ItemIdentifier>
                  <ns3:InfoMessageIdentifier>
                  NINOXE:GeneralMessage:27_1</ns3:InfoMessageIdentifier>
                  <ns3:InfoMessageVersion>1</ns3:InfoMessageVersion>
                  <ns3:InfoChannelRef>Commercial</ns3:InfoChannelRef>
                  <ns3:ValidUntilTime>
                  2017-03-29T20:30:06.000+02:00</ns3:ValidUntilTime>
                  <ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                  xsi:type="ns9:IDFGeneralMessageStructure">
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText xml:lang="NL">La nouvelle carte
                      d'abonnement est disponible au points de vente du
                      réseau</MessageText>
                    </Message>
                  </ns3:Content>
                </ns3:GeneralMessage>
                <ns3:GeneralMessage formatRef="FRANCE">
                  <ns3:RecordedAtTime>
                  2017-03-29T03:30:06.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>3471</ns3:ItemIdentifier>
                  <ns3:InfoMessageIdentifier>
                  NINOXE:GeneralMessage:21_1</ns3:InfoMessageIdentifier>
                  <ns3:InfoMessageVersion>1</ns3:InfoMessageVersion>
                  <ns3:InfoChannelRef>Commercial</ns3:InfoChannelRef>
                  <ns3:ValidUntilTime>
                  2017-03-29T22:30:06.000+02:00</ns3:ValidUntilTime>
                  <ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                  xsi:type="ns9:IDFGeneralMessageStructure">
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText xml:lang="NL">Les nouveaux tarifs sont
                      consultable sur le site internet</MessageText>
                    </Message>
                  </ns3:Content>
                </ns3:GeneralMessage>
              </ns3:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension />
          </ns8:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 3032 - Handle a GeneralMessage response (ServiceDelivery)
    Given a SIRI server waits GeneralMessageRequest request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns8:GetGeneralMessageResponse xmlns:ns3="http://www.siri.org.uk/siri"
          xmlns:ns4="http://www.ifopt.org.uk/acsb"
          xmlns:ns5="http://www.ifopt.org.uk/ifopt"
          xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns7="http://scma/siri"
          xmlns:ns8="http://wsdl.siri.org.uk"
          xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>
              2017-03-29T16:48:00.993+02:00</ns3:ResponseTimestamp>
              <ns3:ProducerRef>NINOXE:default</ns3:ProducerRef>
              <ns3:Address>
              http://appli.chouette.mobi/siri_france/siri</ns3:Address>
              <ns3:ResponseMessageIdentifier>
              b28e8207-f030-4932-966c-3e6099fad4ef</ns3:ResponseMessageIdentifier>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>
                2017-03-29T16:48:00.039+02:00</ns3:ResponseTimestamp>
                <ns3:Status>true</ns3:Status>
                <ns3:GeneralMessage formatRef="FRANCE">
                  <ns3:RecordedAtTime>
                  2017-03-29T03:30:06.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>3477</ns3:ItemIdentifier>
                  <ns3:InfoMessageIdentifier>
                  NINOXE:GeneralMessage:27_1</ns3:InfoMessageIdentifier>
                  <ns3:InfoMessageVersion>1</ns3:InfoMessageVersion>
                  <ns3:InfoChannelRef>Commercial</ns3:InfoChannelRef>
                  <ns3:ValidUntilTime>
                  2017-03-29T20:30:06.000+02:00</ns3:ValidUntilTime>
                  <ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                  xsi:type="ns9:IDFGeneralMessageStructure">
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText xml:lang="NL">La nouvelle carte d'abonnement est disponible au points de vente du réseau</MessageText>
                    </Message>
                  </ns3:Content>
                </ns3:GeneralMessage>
              </ns3:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension />
          </ns8:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ineo                  |
      | remote_objectid_kind | internal              |
    And a minute has passed
    When a minute has passed
    And the SIRI server has received a GeneralMessage request
    Then one Situation has the following attributes:
      | ObjectIDs               | "internal" : "NINOXE:GeneralMessage:27_1"                                                        |
      | RecordedAt              | 2017-03-29T03:30:06+02:00                                                  |
      | Version                 | 1                                                                          |
      | Channel                 | Commercial                                                                 |
      | ProducerRef             | NINOXE:default                                                             |
      | ValidUntil              | 2017-03-29T20:30:06+02:00                                                  |
      | Messages[0]#MessageType | longMessage                                                                |
      | Messages[0]#MessageText | La nouvelle carte d'abonnement est disponible au points de vente du réseau |
