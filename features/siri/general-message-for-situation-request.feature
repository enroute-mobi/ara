Feature: Support SIRI GeneralMessage for Situation

  Background:
      Given a Referential "test" is created

  # Should be fixed by #3797
  @wip
  Scenario: Ignore situations associatd to the Commercial channel
    Given a Situation exists with the following attributes:
      | ObjectIDs               | "internal" : "1"                                 |
      | RecordedAt              | 2017-01-01T03:30:06+02:00                        |
      | Version                 | 1                                                |
      | Channel                 | Commercial                                       |
      | ValidUntil              | 2017-01-01T20:30:06+02:00                        |
      | Messages[0]#MessageType | longMessage                                      |
      | Messages[0]#MessageText | Les situations commercials doivent être ignorées |
    And a Situation exists with the following attributes:
      | ObjectIDs               | "internal" : "2"          |
      | RecordedAt              | 2017-01-01T03:30:06+02:00 |
      | Version                 | 1                         |
      | Channel                 | Perturbations             |
      | ValidUntil              | 2017-01-01T20:30:06+02:00 |
      | Messages[0]#MessageType | longMessage               |
      | Messages[0]#MessageText | Les autres non            |
    And a Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential     | TestPartner |
      | remote_objectid_kind | internal    |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <SOAP-ENV:Header/>
      <S:Body>
        <ns7:GetGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns3="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://wsdl.siri.org.uk/siri" xmlns:ns7="http://wsdl.siri.org.uk">
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
              <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>GeneralMessage:Test:0</ns3:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
                <ns3:Status>true</ns3:Status>

                <ns3:GeneralMessage>
                  <ns3:formatRef>STIF-IDF</ns3:formatRef>
                  <ns3:RecordedAtTime>2017-01-01T03:30:06.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>Edwig:Item::2:LOC</ns3:ItemIdentifier>
                  <ns3:InfoMessageIdentifier>Edwig:InfoMessage::2:LOC</ns3:InfoMessageIdentifier>
                  <ns3:InfoMessageVersion>1</ns3:InfoMessageVersion>
                  <ns3:InfoChannelRef>Perturbations</ns3:InfoChannelRef>
                  <ns3:ValidUntilTime>2017-01-01T20:30:06.000+02:00</ns3:ValidUntilTime>
                  <ns3:Content xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xsi:type='ns9:IDFLineSectionStructure'>
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText>Les autres non</MessageText>
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
