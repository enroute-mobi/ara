Feature: Support SIRI GeneralMessage

  Background:
      Given a Referential "test" is created

  @wip
  Scenario: 3008 - Performs a SIRI GeneralMessage Request to a Partner
    Given a Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    When I send this GeneralMessageRequest
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
      <SOAP-ENV:Header/>
      <S:Body>
        <ns7:GetGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns3="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://wsdl.siri.org.uk/siri" xmlns:ns7="http://wsdl.siri.org.uk">
          <ServiceRequestInfo>
            <ns2:RequestTimestamp>2017-03-29T16:47:58.311Z</ns2:RequestTimestamp>
            <ns2:RequestorRef>NINOXE:default</ns2:RequestorRef>
            <ns2:MessageIdentifier>GeneralMessage:Test:0</ns2:MessageIdentifier>
          </ServiceRequestInfo>
          <Request version="2.0:FR-IDF-2.4">
            <ns2:RequestTimestamp>2017-03-29T16:47:58.311Z</ns2:RequestTimestamp>
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

  @wip
  Scenario: Handle a SIRI ServiceDelivery after GM Request to a SIRI server
    Given a SIRI server waits GeneralMessageRequest on "http://localhost:8090" to respond with
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
    | remote_url | http://localhost:8090 |
    | remote_credential | ratpdev |
    | remote_objectid_kind | internal |
    And a minute has passed
    When I receive this GeneralMessageRequest
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
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
    Then I send this SIRI ServiceDelivery
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
              <ns3:ResponseMessageIdentifier>
              b28e8207-f030-4932-966c-3e6099fad4ef</ns3:ResponseMessageIdentifier>
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
  @wip
  Scenario: Handle a GeneralMessage response (ServiceDelivery)
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
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
              </ns3:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension />
          </ns8:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ineo                  |
      | remote_objectid_kind | internal              |
    And a minute has passed
    When a minute has passed
    And the SIRI server has received a GeneralMessage Request
    Then a Situation exists with the following attributes
      | ObjectIDs                       | "internal" : "3477"                                                        |
      | RecordedAt                      | 2017-03-29T03:30:06.000+02:00                                              |
      | Attribute[ProducerRef]          | NINOXE:default                                                             |
      | Attribute[InfoMessageVersion]   | 1                                                                          |
      | Attribute[InfoChannelRef]       | Commercial                                                                 |
      | Attribute[ValidUntilTime]       | 2017-03-29T20:30:06.000+02:00                                              |
      | Attribute[MessageType]          | longMessage                                                                |
      | Attribute[MessageText]          | La nouvelle carte d'abonnement est disponible au points de vente du réseau |
