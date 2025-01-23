Feature: Support SIRI GeneralMessage for Situation

  Background:
      Given a Referential "test" is created

  @ARA-1362
  Scenario: 3797 - Do not ignore Situations associated to other keywords than Commercial/Perturbation/Information
    Given a Situation exists with the following attributes:
      | Codes                        | "internal" : "1"                              |
      | RecordedAt                   | 2017-01-01T03:30:06+02:00                     |
      | Version                      | 1                                             |
      | Keywords                     | ["Others"]                                    |
      | ValidityPeriods[0]#StartTime | 2016-01-01T20:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00                     |
      | ReportType                   | general                                       |
      | Description[DefaultValue]    | La nouvelle carte d'abonnement est disponible |
      | Description[Translations]#EN | The new pass is available                     |
      | Summary[Translations]#FR     | Nouveau pass Navigo                           |
      | Summary[Translations]#EN     | New pass Navigo                               |
      | Affects[Line]                | 6ba7b814-9dad-11d1-2-00c04fd430c8             |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a SIRI Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential  | TestPartner |
      | remote_code_space | internal    |
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
             <siri:ProducerRef>Ara</siri:ProducerRef>
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
                 <siri:InfoMessageIdentifier>1</siri:InfoMessageIdentifier>
                 <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                 <siri:InfoChannelRef>Information</siri:InfoChannelRef>
                 <siri:ValidUntilTime>2017-01-01T20:30:06.000+02:00</siri:ValidUntilTime>
                 <siri:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                               xsi:type="stif:IDFGeneralMessageStructure">
                   <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                   <siri:Message>
                     <siri:MessageType>shortMessage</siri:MessageType>
                     <siri:MessageText xml:lang='EN'>New pass Navigo</siri:MessageText>
                     <siri:MessageText xml:lang='FR'>Nouveau pass Navigo</siri:MessageText>
                   </siri:Message>
                   <siri:Message>
                     <siri:MessageType>textOnly</siri:MessageType>
                     <siri:MessageText>La nouvelle carte d'abonnement est disponible</siri:MessageText>
                     <siri:MessageText xml:lang='EN'>The new pass is available</siri:MessageText>
                   </siri:Message>
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
      | Codes                                                                              | "external" : "test"                                                        |
      | RecordedAt                                                                         | 2017-01-01T03:30:06+02:00                                                  |
      | Version                                                                            | 1                                                                          |
      | Keywords                                                                           | ["Commercial"]                                                             |
      | ValidityPeriods[0]#EndTime                                                         | 2017-01-01T20:30:06+02:00                                                  |
      | Summary[DefaultValue]                                                              | Carte abonnement                                                           |
      | Description[DefaultValue]                                                          | La nouvelle carte d'abonnement est disponible au points de vente du réseau |
      | Affects[StopArea]                                                                  | 6ba7b814-9dad-11d1-3-00c04fd430c8                                          |
      | Affects[Line]                                                                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                          |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedDestinations[0]/StopAreaId | 6ba7b814-9dad-11d1-3-00c04fd430c8                                          |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/LastStopId     | 6ba7b814-9dad-11d1-4-00c04fd430c8                                          |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/FirstStopId    | 6ba7b814-9dad-11d1-3-00c04fd430c8                                          |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[0]/RouteRef         | Route:66:LOC                                                               |
    And a Line exists with the following attributes:
      | Codes | "external": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test last stop                           |
      | Codes | "external": "NINOXE:StopPoint:SP:25:LOC" |
    And a SIRI Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential  | NINOXE:default |
      | remote_code_space | external       |
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
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4" xmlns:stif="http://wsdl.siri.org.uk/siri">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:GeneralMessage formatRef="STIF-IDF">
                  <siri:RecordedAtTime>2017-01-01T03:30:06.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>RATPDev:Item::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:ItemIdentifier>
                  <siri:InfoMessageIdentifier>test</siri:InfoMessageIdentifier>
                  <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                  <siri:InfoChannelRef>Commercial</siri:InfoChannelRef>
                  <siri:ValidUntilTime>2017-01-01T20:30:06.000+02:00</siri:ValidUntilTime>
                  <siri:Content xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xsi:type='stif:IDFGeneralMessageStructure'>
                    <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:DestinationRef>NINOXE:StopPoint:SP:24:LOC</siri:DestinationRef>
                    <siri:RouteRef>Route:66:LOC</siri:RouteRef>
                    <siri:LineSection>
                      <siri:FirstStop>NINOXE:StopPoint:SP:24:LOC</siri:FirstStop>
                      <siri:LastStop>NINOXE:StopPoint:SP:25:LOC</siri:LastStop>
                      <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    </siri:LineSection>
                    <siri:Message>
                      <siri:MessageType>shortMessage</siri:MessageType>
                      <siri:MessageText>Carte abonnement</siri:MessageText>
                    </siri:Message>
                    <siri:Message>
                      <siri:MessageType>textOnly</siri:MessageType>
                      <siri:MessageText>La nouvelle carte d'abonnement est disponible au points de vente du réseau</siri:MessageText>
                    </siri:Message>
                  </siri:Content>
                </siri:GeneralMessage>
              </siri:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol  | siri                                                         |
      | Direction | received                                                     |
      | Status    | OK                                                           |
      | Type      | GeneralMessageRequest                                        |
      | StopAreas | ["NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:25:LOC"] |
      | Lines     | ["NINOXE:Line:3:LOC"]                                        |

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
                  <siri:RecordedAtTime>2017-03-29T03:30:06.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>3477</siri:ItemIdentifier>
                  <siri:InfoMessageIdentifier>NINOXE:GeneralMessage:27_1</siri:InfoMessageIdentifier>
                  <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                  <siri:InfoChannelRef>Commercial</siri:InfoChannelRef>
                  <siri:ValidUntilTime>2017-03-29T20:50:06.000+02:00</siri:ValidUntilTime>
                  <siri:Content>
                   <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                   <siri:LineRef>1234</siri:LineRef>
                   <siri:DestinationRef>destinationRef1</siri:DestinationRef>
                   <siri:DestinationRef>destinationRef2</siri:DestinationRef>
                   <siri:RouteRef>Route:66:LOC</siri:RouteRef>
                    <siri:LineSection>
                      <siri:FirstStop>NINOXE:StopPoint:SP:25:LOC</siri:FirstStop>
                      <siri:LastStop>NINOXE:StopPoint:SP:26:LOC</siri:LastStop>
                      <siri:LineRef>1234</siri:LineRef>
                    </siri:LineSection>
                    <siri:Message>
                      <siri:MessageType>shortMessage</siri:MessageType>
                      <siri:MessageText xml:lang='EN'>New pass Navigo</siri:MessageText>
                      <siri:MessageText xml:lang='FR'>Nouveau pass Navigo</siri:MessageText>
                   </siri:Message>
                   <siri:Message>
                      <siri:MessageType>textOnly</siri:MessageType>
                      <siri:MessageText>La nouvelle carte d'abonnement est disponible</siri:MessageText>
                      <siri:MessageText xml:lang='EN'>The new pass is available</siri:MessageText>
                   </siri:Message>
                   </siri:Content>
                </siri:GeneralMessage>
                <siri:GeneralMessage formatRef="FRANCE">
                  <siri:RecordedAtTime>2017-03-29T03:30:06.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>3478</siri:ItemIdentifier>
                  <siri:InfoMessageIdentifier>NINOXE:GeneralMessage:27_2</siri:InfoMessageIdentifier>
                  <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                  <siri:InfoChannelRef>Commercial</siri:InfoChannelRef>
                  <siri:ValidUntilTime>2017-03-29T20:50:06.000+02:00</siri:ValidUntilTime>
                  <siri:Content>
                    <siri:LineRef>5678</siri:LineRef>
                    <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText>carte d'abonnement</MessageText>
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
      | remote_url        | http://localhost:8090 |
      | remote_credential | ineo                  |
      | remote_code_space | internal              |
    And a Line exists with the following attributes:
      | Name              | Test              |
      | Codes             | "internal":"1234" |
      | CollectSituations | true              |
    And a Line exists with the following attributes:
      | Name              | Test              |
      | Codes             | "internal":"5678" |
      | CollectSituations | true              |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test1                         |
      | Codes | "internal": "destinationRef1" |
    And a StopArea exists with the following attributes:
      | Name  | Test2                         |
      | Codes | "internal": "destinationRef2" |
    And a StopArea exists with the following attributes:
      | Name  | firstStop                                |
      | Codes | "internal": "NINOXE:StopPoint:SP:25:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | lastStop                                 |
      | Codes | "internal": "NINOXE:StopPoint:SP:26:LOC" |
    And a minute has passed
    When a minute has passed
    And the SIRI server has received a GeneralMessage request
    Then one Situation has the following attributes:
      | Codes                                                                              | "internal" : "NINOXE:GeneralMessage:27_1"     |
      | RecordedAt                                                                         | 2017-03-29T03:30:06+02:00                     |
      | Version                                                                            | 1                                             |
      | Keywords                                                                           | ["Commercial"]                                |
      | Progress                                                                           | published                                     |
      | ProducerRef                                                                        | NINOXE:default                                |
      | ValidityPeriods[0]#StartTime                                                       | 2017-03-29T03:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime                                                         | 2017-03-29T20:50:06+02:00                     |
      | Description[DefaultValue]                                                          | La nouvelle carte d'abonnement est disponible |
      | Description[Translations]#EN                                                       | The new pass is available                     |
      | Summary[Translations]#FR                                                           | Nouveau pass Navigo                           |
      | Summary[Translations]#EN                                                           | New pass Navigo                               |
      | Affects[StopArea]                                                                  | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
      | Affects[Line]                                                                      | 6ba7b814-9dad-11d1-2-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedDestinations[0]/StopAreaId | 6ba7b814-9dad-11d1-5-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedDestinations[1]/StopAreaId | 6ba7b814-9dad-11d1-6-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/FirstStop      | 6ba7b814-9dad-11d1-7-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/LastStop       | 6ba7b814-9dad-11d1-8-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[0]/RouteRef         | Route:66:LOC                                  |
    Then one Situation has the following attributes:
      | Codes                        | "internal" : "NINOXE:GeneralMessage:27_2" |
      | RecordedAt                   | 2017-03-29T03:30:06+02:00                 |
      | Version                      | 1                                         |
      | Keywords                     | ["Commercial"]                            |
      | ProducerRef                  | NINOXE:default                            |
      | Progress                     | published                                 |
      | ValidityPeriods[0]#StartTime | 2017-03-29T03:30:06+02:00                 |
      | ValidityPeriods[0]#EndTime   | 2017-03-29T20:50:06+02:00                 |
      | Description[DefaultValue]    | carte d'abonnement                        |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-4-00c04fd430c8         |
      | Affects[Line]                | 6ba7b814-9dad-11d1-3-00c04fd430c8         |
    And an audit event should exist with these attributes:
      | Protocol  | siri                                                                                                                             |
      | Direction | sent                                                                                                                             |
      | Status    | OK                                                                                                                               |
      | Type      | GeneralMessageRequest                                                                                                            |
      | StopAreas | ["destinationRef1", "destinationRef2", "NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:25:LOC", "NINOXE:StopPoint:SP:26:LOC"] |
      | Lines     | ["1234", "5678"]                                                                                                                 |

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
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:Status>true</siri:Status>
                <siri:GeneralMessage>
                  <siri:formatRef>STIF-IDF</siri:formatRef>
                  <siri:RecordedAtTime>2017-01-01T03:35:00.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>RATPDev:Item::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:ItemIdentifier>
                  <siri:InfoMessageIdentifier>NINOXE:GeneralMessage:27_1</siri:InfoMessageIdentifier>
                  <siri:InfoMessageVersion>2</siri:InfoMessageVersion>
                  <siri:InfoChannelRef>Commercial</siri:InfoChannelRef>
                  <siri:ValidUntilTime>2017-01-07T23:30:06.000+02:00</siri:ValidUntilTime>
                  <siri:Content>
                    <Message>
                      <MessageType>textOnly</MessageType>
                      <MessageText>Points de vente du réseau ouvert demain matin pour la nouvelle carte d'abonnement</MessageText>
                    </Message>
                    <Message>
                      <MessageType>shortMessage</MessageType>
                      <MessageText>Points de vente du réseau ouverts</MessageText>
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
      | Codes                      | "external" : "NINOXE:GeneralMessage:27_1"                                  |
      | RecordedAt                 | 2017-01-01T03:30:06+02:00                                                  |
      | Version                    | 1                                                                          |
      | Channel                    | Commercial                                                                 |
      | ReportType                 | general                                                                    |
      | ValidityPeriods[0]#EndTime | 2017-01-01T20:30:06+02:00                                                  |
      | Summary[DefaultValue]      | La nouvelle carte d'abonnement est disponible au points de vente du réseau |
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | ineo                  |
      | remote_code_space | external              |
    And a Line exists with the following attributes:
      | Name              | Test              |
      | Codes             | "external":"1234" |
      | CollectSituations | true              |
    And a minute has passed
    When a minute has passed
    And the SIRI server has received a GeneralMessage request
    Then one Situation has the following attributes:
      | Codes                        | "external" : "NINOXE:GeneralMessage:27_1"                                         |
      | RecordedAt                   | 2017-01-01T03:35:00+02:00                                                         |
      | Version                      | 2                                                                                 |
      | Keywords                     | ["Commercial"]                                                                    |
      | ReportType                   | general                                                                           |
      | ValidityPeriods[0]#StartTime | 2017-01-01T03:35:00+02:00                                                         |
      | ValidityPeriods[0]#EndTime   | 2017-01-07T23:30:06+02:00                                                         |
      | Summary[DefaultValue]        | Points de vente du réseau ouverts                                                 |
      | Description[DefaultValue]    | Points de vente du réseau ouvert demain matin pour la nouvelle carte d'abonnement |

  Scenario: 3882 - GeneralMessageResponse empty with an expired Situation
    Given a Situation exists with the following attributes:
      | Codes                      | "external" : "test"                                                        |
      | RecordedAt                 | 2017-01-01T03:30:06+02:00                                                  |
      | Version                    | 1                                                                          |
      | Channel                    | Perturbation                                                               |
      | Messages[0]#MessageType    | longMessage                                                                |
      | Messages[0]#MessageText    | La nouvelle carte d'abonnement est disponible au points de vente du réseau |
      | ValidityPeriods[0]#EndTime | 2017-01-01T01:01:00+02:00                                                  |
      | Description[DefaultValue]  | La nouvelle carte d'abonnement est disponible au points de vente du réseau |
    And a SIRI Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential  | NINOXE:default |
      | remote_code_space | external       |
    And a Line exists with the following attributes:
      | Name              | Test              |
      | Codes             | "external":"1234" |
      | CollectSituations | true              |
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
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
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
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</siri:MessageIdentifier>
    </ServiceRequestInfo>
    <Request version='2.0:FR-IDF-2.4'>
      <siri:RequestTimestamp>2017-01-01T12:01:00.000Z</siri:RequestTimestamp>
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</siri:MessageIdentifier>
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
      And a SIRI Partner "test" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
        | remote_url                      | http://localhost:8090 |
        | remote_credential               | test                  |
        | local_credential                | NINOXE:default        |
        | remote_code_space               | internal              |
        | collect.filter_general_messages | true                  |
        | collect.include_lines           | 1234                  |
      And 30 seconds have passed
      And a Line exists with the following attributes:
        | Name              | Test              |
        | Codes             | "internal":"1234" |
        | CollectSituations | true              |
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
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:MessageIdentifier>
    </ServiceRequestInfo>
    <Request version='2.0:FR-IDF-2.4'>
      <siri:RequestTimestamp>2017-01-01T12:00:40.000Z</siri:RequestTimestamp>
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:MessageIdentifier>
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
      And a SIRI Partner "test" exists with connectors [siri-check-status-client, siri-general-message-request-collector] and the following settings:
        | remote_url                      | http://localhost:8090 |
        | remote_credential               | test                  |
        | local_credential                | NINOXE:default        |
        | remote_code_space               | internal              |
        | collect.filter_general_messages | true                  |
      And 30 seconds have passed
      And a StopArea exists with the following attributes:
        | Name              | Test              |
        | Codes             | "internal":"1234" |
        | CollectSituations | true              |
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
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:MessageIdentifier>
    </ServiceRequestInfo>
    <Request version='2.0:FR-IDF-2.4'>
      <siri:RequestTimestamp>2017-01-01T12:00:40.000Z</siri:RequestTimestamp>
      <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:MessageIdentifier>
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

  @siri-valid @ARA-1443
  Scenario: Collect GeneralMessage with internal tags
    Given a SIRI server waits GeneralMessageRequest request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetGeneralMessageResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-03-29T16:48:00.993+02:00</siri:ResponseTimestamp>
              <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
              <siri:Address>http://appli.chouette.mobi/siri_france/siri</siri:Address>
              <siri:ResponseMessageIdentifier>b28e8207-f030-4932-966c-3e6099fad4ef</siri:ResponseMessageIdentifier>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-03-29T16:48:00.039+02:00</siri:ResponseTimestamp>
                <siri:Status>true</siri:Status>
                <siri:GeneralMessage formatRef="FRANCE">
                  <siri:RecordedAtTime>2017-03-29T03:30:06.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>3477</siri:ItemIdentifier>
                  <siri:InfoMessageIdentifier>NINOXE:GeneralMessage:27_1</siri:InfoMessageIdentifier>
                  <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                  <siri:InfoChannelRef>Commercial</siri:InfoChannelRef>
                  <siri:ValidUntilTime>2017-03-29T20:50:06.000+02:00</siri:ValidUntilTime>
                  <siri:Content>
                   <siri:LineRef>1234</siri:LineRef>
                    <Message>
                      <MessageType>longMessage</MessageType>
                      <MessageText xml:lang="NL">La nouvelle carte d'abonnement est disponible</MessageText>
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
      | remote_url                       | http://localhost:8090 |
      | remote_credential                | ineo                  |
      | remote_code_space                | internal              |
      | collect.situations.internal_tags | first,second          |
    And a Line exists with the following attributes:
      | Name              | Test              |
      | Codes             | "internal":"1234" |
      | CollectSituations | true              |
    And a minute has passed
    When a minute has passed
    And the SIRI server has received a GeneralMessage request
    Then one Situation has the following attributes:
      | Codes                        | "internal" : "NINOXE:GeneralMessage:27_1"     |
      | InternalTags                 | ["first","second"]                            |

  @ARA-1444
  Scenario: Broadcast Situation GeneralMessage with internal tags
    Given a Situation exists with the following attributes:
      | Codes                      | "external" : "test"                           |
      | RecordedAt                 | 2017-01-01T03:30:06+02:00                     |
      | Version                    | 1                                             |
      | InternalTags               | ["first","second"]                            |
      | Keywords                   | ["Commercial"]                                |
      | ValidityPeriods[0]#EndTime | 2017-01-01T20:30:06+02:00                     |
      | Description[DefaultValue]  | La nouvelle carte d'abonnement est disponible |
      | Affects[StopArea]          | 6ba7b814-9dad-11d1-3-00c04fd430c8             |
      | Affects[Line]              | 6ba7b814-9dad-11d1-2-00c04fd430c8             |
    And a Line exists with the following attributes:
      | Codes | "external": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a SIRI Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential                   | NINOXE:default |
      | remote_code_space                  | external       |
      | broadcast.situations.internal_tags | first, another |
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
              <siri:ProducerRef>Ara</siri:ProducerRef>
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
                  <siri:InfoMessageIdentifier>test</siri:InfoMessageIdentifier>
                  <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                  <siri:InfoChannelRef>Commercial</siri:InfoChannelRef>
                  <siri:ValidUntilTime>2017-01-01T20:30:06.000+02:00</siri:ValidUntilTime>
                  <siri:Content xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xsi:type='stif:IDFGeneralMessageStructure'>
                    <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:Message>
                      <siri:MessageType>textOnly</siri:MessageType>
                      <siri:MessageText>La nouvelle carte d'abonnement est disponible</siri:MessageText>
                    </siri:Message>
                  </siri:Content>
                </siri:GeneralMessage>
              </siri:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
      """

  @ARA-1542
  Scenario: Handle a SIRI GeneralMessage with partner setting broadcast.situations.time_to_live outside broadcast period wihout RequestTimeStamp should not broadcast situation
    Given a Situation exists with the following attributes:
      | Codes                        | "external" : "test"               |
      | RecordedAt                   | 2017-01-01T03:30:06+02:00         |
      | Version                      | 1                                 |
      | Keywords                     | ["Commercial", "Test"]            |
      | ReportType                   | general                           |
      | ParticipantRef               | "535"                             |
      | VersionedAt                  | 2017-01-01T01:02:03+02:00         |
      | Progress                     | published                         |
      | Reality                      | test                              |
      | ValidityPeriods[0]#StartTime | 2017-01-01T03:10:06+02:00         |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T03:14:06+02:00         |
      | Affects[AllLines]            |                                   |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a SIRI Partner "test" exists with connectors [siri-general-message-request-broadcaster] and the following settings:
      | local_credential                  | NINOXE:default |
      | remote_code_space                 | external       |
      | broadcast.situations.time_to_live | 5m             |
    # Situation BroadcastPeriod() ends at 2017-01-01T03:14:06+02:00, and requestPeriod will start at 2017-01-01T05:15:06+02:00
    # so the Situation should not be broadcasted
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
    # Situation BroadcastPeriod() ends at 2017-01-01T03:14:06+02:00, and requestPeriod will start at 2017-01-01T05:15:06+02:00
    # so the Situation should not be broadcasted
    Then I should receive this SIRI response
    """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetGeneralMessageResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:GeneralMessageDelivery version='2.0:FR-IDF-2.4' xmlns:stif='http://wsdl.siri.org.uk/siri'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>GeneralMessage:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
              </siri:GeneralMessageDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetGeneralMessageResponse>
        </S:Body>
      </S:Envelope>
    """
