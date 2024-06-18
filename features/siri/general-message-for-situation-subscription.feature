Feature: Support SIRI GeneralMessage by subscription

  Background:
      Given a Referential "test" is created

  @ARA-1362
  Scenario: 3863 - Manage a GM Subscription
    Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
       """
   <?xml version='1.0' encoding='utf-8'?>
   <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
   <S:Body>
     <ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
       <SubscriptionAnswerInfo
         xmlns:ns2="http://www.ifopt.org.uk/acsb"
         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
         xmlns:ns5="http://www.siri.org.uk/siri"
         xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
         <ns5:Address>http://appli.chouette.mobi/siri_france/siri</ns5:Address>
         <ns5:ResponderRef>NINOXE:default</ns5:ResponderRef>
         <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</ns5:RequestMessageRef>
       </SubscriptionAnswerInfo>
       <Answer
         xmlns:ns2="http://www.ifopt.org.uk/acsb"
         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
         xmlns:ns5="http://www.siri.org.uk/siri"
         xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseStatus>
             <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
             <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</ns5:RequestMessageRef>
             <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
             <ns5:SubscriptionRef>6ba7b814-9dad-11d1-6-00c04fd430c8</ns5:SubscriptionRef>
             <ns5:Status>true</ns5:Status>
             <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
         </ns5:ResponseStatus>
         <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
       </Answer>
       <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
     </ns1:SubscribeResponse>
   </S:Body>
   </S:Envelope>
   """
      And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-general-message-subscription-collector] and the following settings:
        | remote_url                      | http://localhost:8090 |
        | remote_credential               | test                  |
        | local_credential                | NINOXE:default        |
        | remote_code_space            | internal              |
        | collect.filter_general_messages | true                  |
        | collect.include_lines           | NINOXE:Line::3:LOC    |
      And 30 seconds have passed
      And a Line exists with the following attributes:
        | Name              | Test                            |
        | Codes         | "internal":"NINOXE:Line::3:LOC" |
        | CollectSituations | true                            |
      And a Line exists with the following attributes:
        | Name              | Test                            |
        | Codes             | "internal":"NINOXE:Line::4:LOC" |
        | CollectSituations | true                            |
      And a StopArea exists with the following attributes:
        | Name              | Test                                    |
        | Codes         | "internal":"NINOXE:StopPoint:SP:24:LOC" |
        | CollectSituations | true                                    |
      And a StopArea exists with the following attributes:
        | Name              | Test                                    |
        | Codes         | "internal":"NINOXE:StopPoint:SP:12:LOC" |
        | CollectSituations | true                                    |
      And 10 seconds have passed
      And 5 seconds have passed
      And a Situation exists with the following attributes:
        | Codes                  | "internal" : "NINOXE:GeneralMessage:27_1" |
        | RecordedAt                 | 2017-01-01T03:30:06+02:00                 |
        | Version                    | 1                                         |
        | Keywords                   | ["Perturbation"]                          |
        | ValidityPeriods[0]#EndTime | 2017-01-01T20:30:06+02:00                 |
        | Description[DefaultValue]  | Les autres non                            |
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
    <S:Body>
      <ns1:NotifyGeneralMessage xmlns:ns1="http://wsdl.siri.org.uk">
       <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns5:ResponseTimestamp>
         <ns5:ProducerRef>NINOXE:default</ns5:ProducerRef>
         <ns5:ResponseMessageIdentifier>NAVINEO:SM:NOT:427843</ns5:ResponseMessageIdentifier>
         <ns5:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns5:RequestMessageRef>
       </ServiceDeliveryInfo>
       <Notification xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
            <ns3:ResponseTimestamp>2017-03-29T16:47:53.039+02:00</ns3:ResponseTimestamp>
            <ns5:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns5:RequestMessageRef>
            <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
            <ns5:SubscriptionRef>6ba7b814-9dad-11d1-7-00c04fd430c8</ns5:SubscriptionRef>
            <ns3:Status>true</ns3:Status>
            <ns3:GeneralMessage>
               <ns3:RecordedAtTime>2017-03-01T03:30:06.000+01:00</ns3:RecordedAtTime>
               <ns3:ItemIdentifier>3477</ns3:ItemIdentifier>
               <ns3:InfoMessageIdentifier>NINOXE:GeneralMessage:27_1</ns3:InfoMessageIdentifier>
               <ns3:InfoMessageVersion>2</ns3:InfoMessageVersion>
               <ns3:formatRef>STIF-IDF</ns3:formatRef>
               <ns3:InfoChannelRef>Commercial</ns3:InfoChannelRef>
               <ns3:ValidUntilTime>2017-03-29T03:30:06.000+01:00</ns3:ValidUntilTime>
               <ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns9:IDFGeneralMessageStructure">
               <LineRef>NINOXE:Line::3:LOC</LineRef>
                  <Message>
                    <MessageType>textOnly</MessageType>
                    <MessageText xml:lang="NL">La nouvelle carte d'abonnement est disponible au points de vente du réseau</MessageText>
                  </Message>
                  <LineSection>
                    <FirstStop>NINOXE:StopPoint:SP:24:LOC</FirstStop>
                    <LastStop>NINOXE:StopPoint:SP:12:LOC</LastStop>
                    <LineRef>NINOXE:Line::3:LOC</LineRef>
                  </LineSection>
               </ns3:Content>
            </ns3:GeneralMessage>
            <ns3:GeneralMessage>
               <ns3:RecordedAtTime>2017-03-01T03:30:06.000+01:00</ns3:RecordedAtTime>
               <ns3:ItemIdentifier>3477</ns3:ItemIdentifier>
               <ns3:InfoMessageIdentifier>NINOXE:GeneralMessage:27_2</ns3:InfoMessageIdentifier>
               <ns3:InfoMessageVersion>2</ns3:InfoMessageVersion>
               <ns3:formatRef>STIF-IDF</ns3:formatRef>
               <ns3:InfoChannelRef>Commercial</ns3:InfoChannelRef>
               <ns3:ValidUntilTime>2017-03-29T03:30:06.000+01:00</ns3:ValidUntilTime>
               <ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns9:IDFGeneralMessageStructure">
               <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
                  <Message>
                    <MessageType>textOnly</MessageType>
                    <MessageText xml:lang="NL">carte d'abonnement</MessageText>
                  </Message>
                  <LineRef>NINOXE:Line::4:LOC</LineRef>
               </ns3:Content>
            </ns3:GeneralMessage>
         </ns3:GeneralMessageDelivery>
       </Notification>
       <SiriExtension/>
      </ns1:NotifyGeneralMessage>
    </S:Body>
    </S:Envelope>
      """
    Then one Situation has the following attributes:
      | Codes                                                                         | "internal" : "NINOXE:GeneralMessage:27_1" |
      | Keywords                                                                      | ["Commercial"]                            |
      | ReportType                                                                    | general                                   |
      | Progress                                                                      | published                                 |
      | ValidityPeriods[0]#StartTime                                                  | 2017-03-01T03:30:06+01:00                 |
      | ValidityPeriods[0]#EndTime                                                    | 2017-03-29T03:30:06+01:00                 |
      | Version                                                                       | 2                                         |
      | Affects[Line]                                                                 | 6ba7b814-9dad-11d1-3-00c04fd430c8         |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedSections[0]/LastStop  | 6ba7b814-9dad-11d1-6-00c04fd430c8         |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedSections[0]/FirstStop | 6ba7b814-9dad-11d1-5-00c04fd430c8         |
    Then one Situation has the following attributes:
      | Codes                        | "internal" : "NINOXE:GeneralMessage:27_2" |
      | Keywords                     | ["Commercial"]                            |
      | ReportType                   | general                                   |
      | Progress                     | published                                 |
      | ValidityPeriods[0]#StartTime | 2017-03-01T03:30:06+01:00                 |
      | ValidityPeriods[0]#EndTime   | 2017-03-29T03:30:06+01:00                 |
      | Version                      | 2                                         |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-5-00c04fd430c8         |
      | Affects[Line]                | 6ba7b814-9dad-11d1-4-00c04fd430c8         |
    And an audit event should exist with these attributes:
      | Protocol  | siri                                                         |
      | Direction | received                                                     |
      | Status    | OK                                                           |
      | Type      | NotifyGeneralMessage                                         |
      | StopAreas | ["NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:12:LOC"] |
      | Lines     | ["NINOXE:Line::3:LOC", "NINOXE:Line::4:LOC"]                 |

  Scenario: 3865 - Manage a InfoMessageCancellation
    Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
    """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
          <SubscriptionAnswerInfo
            xmlns:ns2="http://www.ifopt.org.uk/acsb"
            xmlns:ns3="http://www.ifopt.org.uk/ifopt"
            xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
            xmlns:ns5="http://www.siri.org.uk/siri"
            xmlns:ns6="http://wsdl.siri.org.uk/siri">
            <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
            <ns5:Address>http://appli.chouette.mobi/siri_france/siri</ns5:Address>
            <ns5:ResponderRef>NINOXE:default</ns5:ResponderRef>
            <ns5:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns5:MessageRefStructure">Subscription:Test:0</ns5:RequestMessageRef>
          </SubscriptionAnswerInfo>
          <Answer
            xmlns:ns2="http://www.ifopt.org.uk/acsb"
            xmlns:ns3="http://www.ifopt.org.uk/ifopt"
            xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
            xmlns:ns5="http://www.siri.org.uk/siri"
            xmlns:ns6="http://wsdl.siri.org.uk/siri">
            <ns5:ResponseStatus>
                <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>Subscription:Test:0</ns5:RequestMessageRef>
                <ns5:SubscriberRef>SubscriberRef</ns5:SubscriberRef>
                <ns5:SubscriptionRef>SubscriptionIdentifier</ns5:SubscriptionRef>
                <ns5:Status>true</ns5:Status>
                <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
            </ns5:ResponseStatus>
            <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
          </Answer>
          <AnswerExtension/>
        </ns1:SubscribeResponse>
      </S:Body>
      </S:Envelope>
    """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-general-message-subscription-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind | GeneralMessageCollect |
    And a Situation exists with the following attributes:
      | Codes                  | "internal" : "2"          |
      | RecordedAt                 | 2017-01-01T03:30:06+02:00 |
      | Version                    | 1                         |
      | Channel                    | Perturbations             |
      | ValidityPeriods[0]#EndTIme | 2017-01-01T20:30:06+02:00 |
      | Descriptions[DefaultValue] | Les autres non            |
    When I send this SIRI request
    """
    <?xml version='1.0' encoding='utf-8'?>
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <ns6:NotifyGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri"
        xmlns:ns3="http://www.ifopt.org.uk/acsb"
        xmlns:ns4="http://www.ifopt.org.uk/ifopt"
        xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
        xmlns:ns6="http://wsdl.siri.org.uk"
        xmlns:ns7="http://wsdl.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
            <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
            <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
            <ns2:RequestMessageRef>GeneralMessage:TestDelivery:0</ns2:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Notification>
            <ns2:GeneralMessageDelivery version="1.3">
              <ns2:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns2:ResponseTimestamp>
              <ns2:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns2:RequestMessageRef>
              <ns2:SubscriberRef>RATPDEV:Concerto</ns2:SubscriberRef>
              <ns2:SubscriptionRef>6ba7b814-9dad-11d1-3-00c04fd430c8</ns2:SubscriptionRef>
              <ns2:Status>true</ns2:Status>
              <ns2:GeneralMessageCancellation>
                <ns2:RecordedAtTime>2017-05-15T13:26:10.116+02:00</ns2:RecordedAtTime>
                <ns2:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns2:ItemRef>
                <ns2:InfoMessageIdentifier>2</ns2:InfoMessageIdentifier>
              </ns2:GeneralMessageCancellation>
            </ns2:GeneralMessageDelivery>
          </Notification>
          <SiriExtension />
        </ns6:NotifyGeneralMessage>
      </soap:Body>
    </soap:Envelope>
    """
    Then a Situation "internal:2" should not exist in Referential "test"

  Scenario: Brodcast a GeneralMessage Notification after modification of a Situation
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, siri-general-message-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_code_space | internal              |
    And a Subscription exist with the following attributes:
      | Kind              | GeneralMessageBroadcast                     |
      | ExternalId        | externalId                                  |
      | SubscriberRef     | subscriber                                  |
      | ReferenceArray[0] | Situation, "SituationResource": "Situation" |
    And a Line exists with the following attributes:
        | Name              | Test              |
        | Codes         | "internal":"1234" |
        | CollectSituations | true              |
    And a Situation exists with the following attributes:
      | Codes                                                                           | "internal" : "NINOXE:GeneralMessage:27_1" |
      | RecordedAt                                                                          | 2017-01-01T03:30:06+02:00                 |
      | Version                                                                             | 1                                         |
      | Keywords                                                                            | ["Perturbation"]                          |
      | ValidityPeriods[0]#EndTime                                                          | 2017-01-01T20:30:06+02:00                 |
      | Description[DefaultValue]                                                           | a very very very long message             |
      | Affects[Line]                                                                       | 6ba7b814-9dad-11d1-3-00c04fd430c8         |
      | Affects[StopArea]                                                                   | 6ba7b814-9dad-11d1-5-00c04fd430c8         |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedDestinations[0]/StopAreaId  | 6ba7b814-9dad-11d1-6-00c04fd430c8         |
    And a StopArea exists with the following attributes:
        | Name              | Test                                    |
        | Codes         | "internal":"NINOXE:StopPoint:SP:24:LOC" |
        | CollectSituations | true                                    |
    And a StopArea exists with the following attributes:
        | Name              | Test                                    |
        | Codes         | "internal":"NINOXE:StopPoint:SP:12:LOC" |
        | CollectSituations | true                                    |
    And 10 seconds have passed
    When the Situation "6ba7b814-9dad-11d1-4-00c04fd430c8" is edited with the following attributes:
      | RecordedAt                 | 2017-01-01T03:50:06+02:00              |
      | ValidityPeriods[0]#EndTime | 2017-10-24T20:30:06+02:00              |
      | Description[DefaultValue]  | an ANOTHER very very very long message |
      | Version                    | 2                                      |
    And 15 seconds have passed
    Then the SIRI server should receive this response
    """
     <?xml version='1.0' encoding='utf-8'?>
     <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
       <S:Body>
         <sw:NotifyGeneralMessage xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
           <ServiceDeliveryInfo>
             <siri:ResponseTimestamp>2017-01-01T12:00:25.000Z</siri:ResponseTimestamp>
             <siri:ProducerRef>test</siri:ProducerRef>
             <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
             <siri:RequestMessageRef></siri:RequestMessageRef>
           </ServiceDeliveryInfo>
           <Notification>
             <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4" xmlns:stif="http://wsdl.siri.org.uk/siri">
               <siri:ResponseTimestamp>2017-01-01T12:00:25.000Z</siri:ResponseTimestamp>
               <siri:RequestMessageRef></siri:RequestMessageRef>
               <siri:SubscriberRef>subscriber</siri:SubscriberRef>
               <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
               <siri:Status>true</siri:Status>
               <siri:GeneralMessage formatRef="STIF-IDF">
                 <siri:RecordedAtTime>2017-01-01T03:50:06.000+02:00</siri:RecordedAtTime>
                 <siri:ItemIdentifier>RATPDev:Item::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ItemIdentifier>
                 <siri:InfoMessageIdentifier>NINOXE:GeneralMessage:27_1</siri:InfoMessageIdentifier>
                 <siri:InfoMessageVersion>2</siri:InfoMessageVersion>
                 <siri:InfoChannelRef>Perturbation</siri:InfoChannelRef>
                 <siri:ValidUntilTime>2017-10-24T20:30:06.000+02:00</siri:ValidUntilTime>
                 <siri:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                               xsi:type="stif:IDFGeneralMessageStructure">
                   <siri:LineRef>1234</siri:LineRef>
                   <siri:DestinationRef>NINOXE:StopPoint:SP:12:LOC</siri:DestinationRef>
                   <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                   <Message>
                     <MessageType>textOnly</MessageType>
                     <MessageText>an ANOTHER very very very long message</MessageText>
                   </Message>
                 </siri:Content>
               </siri:GeneralMessage>
             </siri:GeneralMessageDelivery>
           </Notification>
           <SiriExtension />
         </sw:NotifyGeneralMessage>
       </S:Body>
     </S:Envelope>
    """
    And an audit event should exist with these attributes:
      | Protocol  | siri                                                         |
      | Direction | sent                                                         |
      | Status    | OK                                                           |
      | Type      | NotifyGeneralMessage                                         |
      | StopAreas | ["NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:12:LOC"] |
      | Lines     | ["1234"]                                                     |

  Scenario: Brodcast a GeneralMessage Notification when keywords does not contains Perturbation/Information/Commercial but ReportType is type incident should broadcast as Pertubation
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, siri-general-message-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_code_space | internal              |
    And a Subscription exist with the following attributes:
      | Kind              | GeneralMessageBroadcast                     |
      | ExternalId        | externalId                                  |
      | SubscriberRef     | subscriber                                  |
      | ReferenceArray[0] | Situation, "SituationResource": "Situation" |
    And a Line exists with the following attributes:
      | Name              | Test              |
      | Codes         | "internal":"1234" |
      | CollectSituations | true              |
    And a Situation exists with the following attributes:
      | Codes                  | "internal" : "NINOXE:GeneralMessage:27_1" |
      | RecordedAt                 | 2017-01-01T03:30:06+02:00                 |
      | Version                    | 1                                         |
      | Keywords                   | ["Other"]                                 |
      | ValidityPeriods[0]#EndTime | 2017-01-01T20:30:06+02:00                 |
      | ReportType                 | incident                                  |
      | Description[DefaultValue]  | a very very very long message             |
      | Affects[Line]              | 6ba7b814-9dad-11d1-3-00c04fd430c8         |
    And 20 seconds have passed
    Then the SIRI server should receive this response
    """
     <?xml version='1.0' encoding='utf-8'?>
     <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
       <S:Body>
         <sw:NotifyGeneralMessage xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
           <ServiceDeliveryInfo>
             <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
             <siri:ProducerRef>test</siri:ProducerRef>
             <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
             <siri:RequestMessageRef></siri:RequestMessageRef>
           </ServiceDeliveryInfo>
           <Notification>
             <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4" xmlns:stif="http://wsdl.siri.org.uk/siri">
               <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
               <siri:RequestMessageRef></siri:RequestMessageRef>
               <siri:SubscriberRef>subscriber</siri:SubscriberRef>
               <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
               <siri:Status>true</siri:Status>
               <siri:GeneralMessage formatRef="STIF-IDF">
                 <siri:RecordedAtTime>2017-01-01T03:30:06.000+02:00</siri:RecordedAtTime>
                 <siri:ItemIdentifier>RATPDev:Item::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ItemIdentifier>
                 <siri:InfoMessageIdentifier>NINOXE:GeneralMessage:27_1</siri:InfoMessageIdentifier>
                 <siri:InfoMessageVersion>1</siri:InfoMessageVersion>
                 <siri:InfoChannelRef>Perturbation</siri:InfoChannelRef>
                 <siri:ValidUntilTime>2017-01-01T20:30:06.000+02:00</siri:ValidUntilTime>
                 <siri:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                               xsi:type="stif:IDFGeneralMessageStructure">
                   <siri:LineRef>1234</siri:LineRef>
                   <Message>
                     <MessageType>textOnly</MessageType>
                     <MessageText>a very very very long message</MessageText>
                   </Message>
                 </siri:Content>
               </siri:GeneralMessage>
             </siri:GeneralMessageDelivery>
           </Notification>
           <SiriExtension />
         </sw:NotifyGeneralMessage>
       </S:Body>
     </S:Envelope>
    """

#   Scenario: Manage a Subscription without filter
#     Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
#     """
#     """
#     And a Partner "test" exists with connectors [siri-check-status-client, siri-general-message-subscription-collector] and the following settings:
#       | remote_url           | http://localhost:8090 |
#       | remote_credential    | test                  |
#       | local_credential     | NINOXE:default        |
#       | remote_code_space | internal              |
#     And 30 seconds have passed
#     And 30 seconds have passed
#     And 5 seconds have passed
#     And the SIRI server has received a Subscribe request
#     Then the SIRI server should receive this response
#       """
# <?xml version='1.0' encoding='utf-8'?>
# <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
# <S:Body>
#   <sw:Subscribe xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri' xmlns:sws='http://wsdl.siri.org.uk/siri'>
#     <SubscriptionRequestInfo>
#       <siri:RequestTimestamp>2017-01-01T12:01:05.000Z</siri:RequestTimestamp>
#       <siri:RequestorRef>test</siri:RequestorRef>
#       <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:MessageIdentifier>
#     </SubscriptionRequestInfo>
#     <Request>
#       <siri:GeneralMessageSubscriptionRequest>
#         <siri:SubscriberRef>test</siri:SubscriberRef>
#         <siri:SubscriptionIdentifier>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:SubscriptionIdentifier>
#         <siri:InitialTerminationTime>2017-01-03T12:01:05.000Z</siri:InitialTerminationTime>
#         <siri:GeneralMessageRequest version='2.0:FR-IDF-2.4'>
#           <siri:RequestTimestamp>2017-01-01T12:01:05.000Z</siri:RequestTimestamp>
#           <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:MessageIdentifier>
#             <siri:Extensions>
#               <sws:IDFGeneralMessageRequestFilter>
#               </sws:IDFGeneralMessageRequestFilter>
#             </siri:Extensions>
#         </siri:GeneralMessageRequest>
#       </siri:GeneralMessageSubscriptionRequest>
#     </Request>
#     <RequestExtension/>
#   </sw:Subscribe>
# </S:Body>
# </S:Envelope>
#       """

#   Scenario: Manage a Subscription with a Line filter
#     Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
#     """
#     """
#       And a Partner "test" exists with connectors [siri-check-status-client, siri-general-message-subscription-collector] and the following settings:
#         | remote_url                      | http://localhost:8090 |
#         | remote_credential               | test                  |
#         | local_credential                | NINOXE:default        |
#         | remote_code_space            | internal              |
#         | collect.filter_general_messages | true                  |
#       And 30 seconds have passed
#       And a Line exists with the following attributes:
#         | Name                   | Test              |
#         | Codes              | "internal":"1234" |
#         | CollectSituations | true              |
#       And 10 seconds have passed
#       And 5 seconds have passed
#       And the SIRI server has received a Subscribe request
#     Then the SIRI server should receive this response
#       """
# <?xml version='1.0' encoding='utf-8'?>
# <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
# <S:Body>
#   <sw:Subscribe xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri' xmlns:sws='http://wsdl.siri.org.uk/siri'>
#     <SubscriptionRequestInfo>
#       <siri:RequestTimestamp>2017-01-01T12:00:45.000Z</siri:RequestTimestamp>
#       <siri:RequestorRef>test</siri:RequestorRef>
#       <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:MessageIdentifier>
#     </SubscriptionRequestInfo>
#     <Request>
#       <siri:GeneralMessageSubscriptionRequest>
#         <siri:SubscriberRef>test</siri:SubscriberRef>
#         <siri:SubscriptionIdentifier>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:SubscriptionIdentifier>
#         <siri:InitialTerminationTime>2017-01-03T12:00:45.000Z</siri:InitialTerminationTime>
#         <siri:GeneralMessageRequest version='2.0:FR-IDF-2.4'>
#           <siri:RequestTimestamp>2017-01-01T12:00:45.000Z</siri:RequestTimestamp>
#           <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:MessageIdentifier>
#             <siri:Extensions>
#               <sws:IDFGeneralMessageRequestFilter>
#                 <siri:LineRef>1234</siri:LineRef>
#               </sws:IDFGeneralMessageRequestFilter>
#             </siri:Extensions>
#         </siri:GeneralMessageRequest>
#       </siri:GeneralMessageSubscriptionRequest>
#     </Request>
#     <RequestExtension/>
#   </sw:Subscribe>
# </S:Body>
# </S:Envelope>
#       """

#   Scenario: Manage a Subscription with a StopArea filter
#     Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
#     """
#     """
#       And a Partner "test" exists with connectors [siri-check-status-client, siri-general-message-subscription-collector] and the following settings:
#         | remote_url                      | http://localhost:8090 |
#         | remote_credential               | test                  |
#         | local_credential                | NINOXE:default        |
#         | remote_code_space            | internal              |
#         | collect.filter_general_messages | true                  |
#       And 30 seconds have passed
#       And a StopArea exists with the following attributes:
#         | Name                   | Test              |
#         | Codes              | "internal":"1234" |
#         | CollectGeneralMessages | true              |
#       And 10 seconds have passed
#       And 5 seconds have passed
#       And the SIRI server has received a Subscribe request
#     Then the SIRI server should receive this response
#       """
# <?xml version='1.0' encoding='utf-8'?>
# <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
# <S:Body>
#   <sw:Subscribe xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri' xmlns:sws='http://wsdl.siri.org.uk/siri'>
#     <SubscriptionRequestInfo>
#       <siri:RequestTimestamp>2017-01-01T12:00:45.000Z</siri:RequestTimestamp>
#       <siri:RequestorRef>test</siri:RequestorRef>
#       <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:MessageIdentifier>
#     </SubscriptionRequestInfo>
#     <Request>
#       <siri:GeneralMessageSubscriptionRequest>
#         <siri:SubscriberRef>test</siri:SubscriberRef>
#         <siri:SubscriptionIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:SubscriptionIdentifier>
#         <siri:InitialTerminationTime>2017-01-03T12:00:45.000Z</siri:InitialTerminationTime>
#         <siri:GeneralMessageRequest version='2.0:FR-IDF-2.4'>
#           <siri:RequestTimestamp>2017-01-01T12:00:45.000Z</siri:RequestTimestamp>
#           <siri:MessageIdentifier>RATPDev:Message::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:MessageIdentifier>
#             <siri:Extensions>
#               <sws:IDFGeneralMessageRequestFilter>
#                 <siri:StopPointRef>1234</siri:StopPointRef>
#               </sws:IDFGeneralMessageRequestFilter>
#             </siri:Extensions>
#         </siri:GeneralMessageRequest>
#       </siri:GeneralMessageSubscriptionRequest>
#     </Request>
#     <RequestExtension/>
#   </sw:Subscribe>
# </S:Body>
# </S:Envelope>
#       """

  @ARA-957
  Scenario: Send DeleteSubscriptionRequests
    Given a SIRI server on "http://localhost:8090"
      And a Partner "test" exists with connectors [siri-general-message-subscription-collector] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_code_space | internal              |
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <ns6:NotifyGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri"
        xmlns:ns3="http://www.ifopt.org.uk/acsb"
        xmlns:ns4="http://www.ifopt.org.uk/ifopt"
        xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
        xmlns:ns6="http://wsdl.siri.org.uk"
        xmlns:ns7="http://wsdl.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
            <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
            <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
            <ns2:RequestMessageRef>GeneralMessage:TestDelivery:0</ns2:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Notification>
            <ns2:GeneralMessageDelivery version="1.3">
              <ns2:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns2:ResponseTimestamp>
              <ns2:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns2:RequestMessageRef>
              <ns2:SubscriberRef>RATPDEV:Concerto</ns2:SubscriberRef>
              <ns2:SubscriptionRef>6ba7b814-9dad-11d1-3-00c04fd430c8</ns2:SubscriptionRef>
              <ns2:Status>true</ns2:Status>
              <ns2:GeneralMessageCancellation>
                <ns2:RecordedAtTime>2017-05-15T13:26:10.116+02:00</ns2:RecordedAtTime>
                <ns2:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns2:ItemRef>
                <ns2:InfoMessageIdentifier>2</ns2:InfoMessageIdentifier>
              </ns2:GeneralMessageCancellation>
            </ns2:GeneralMessageDelivery>
          </Notification>
          <SiriExtension />
        </ns6:NotifyGeneralMessage>
      </soap:Body>
    </soap:Envelope>
    """
    Then the SIRI server should have received 1 DeleteSubscription request
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <ns6:NotifyGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri"
        xmlns:ns3="http://www.ifopt.org.uk/acsb"
        xmlns:ns4="http://www.ifopt.org.uk/ifopt"
        xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
        xmlns:ns6="http://wsdl.siri.org.uk"
        xmlns:ns7="http://wsdl.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
            <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
            <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
            <ns2:RequestMessageRef>GeneralMessage:TestDelivery:0</ns2:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Notification>
            <ns2:GeneralMessageDelivery version="1.3">
              <ns2:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns2:ResponseTimestamp>
              <ns2:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns2:RequestMessageRef>
              <ns2:SubscriberRef>RATPDEV:Concerto</ns2:SubscriberRef>
              <ns2:SubscriptionRef>6ba7b814-9dad-11d1-3-00c04fd430c8</ns2:SubscriptionRef>
              <ns2:Status>true</ns2:Status>
              <ns2:GeneralMessageCancellation>
                <ns2:RecordedAtTime>2017-05-15T13:26:10.116+02:00</ns2:RecordedAtTime>
                <ns2:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns2:ItemRef>
                <ns2:InfoMessageIdentifier>2</ns2:InfoMessageIdentifier>
              </ns2:GeneralMessageCancellation>
            </ns2:GeneralMessageDelivery>
          </Notification>
          <SiriExtension />
        </ns6:NotifyGeneralMessage>
      </soap:Body>
    </soap:Envelope>
    """
    Then the SIRI server should not have received 2 DeleteSubscription requests
    When 6 minutes have passed
      And I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <ns6:NotifyGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri"
        xmlns:ns3="http://www.ifopt.org.uk/acsb"
        xmlns:ns4="http://www.ifopt.org.uk/ifopt"
        xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
        xmlns:ns6="http://wsdl.siri.org.uk"
        xmlns:ns7="http://wsdl.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
            <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
            <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
            <ns2:RequestMessageRef>GeneralMessage:TestDelivery:0</ns2:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Notification>
            <ns2:GeneralMessageDelivery version="1.3">
              <ns2:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns2:ResponseTimestamp>
              <ns2:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns2:RequestMessageRef>
              <ns2:SubscriberRef>RATPDEV:Concerto</ns2:SubscriberRef>
              <ns2:SubscriptionRef>6ba7b814-9dad-11d1-3-00c04fd430c8</ns2:SubscriptionRef>
              <ns2:Status>true</ns2:Status>
              <ns2:GeneralMessageCancellation>
                <ns2:RecordedAtTime>2017-05-15T13:26:10.116+02:00</ns2:RecordedAtTime>
                <ns2:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns2:ItemRef>
                <ns2:InfoMessageIdentifier>2</ns2:InfoMessageIdentifier>
              </ns2:GeneralMessageCancellation>
            </ns2:GeneralMessageDelivery>
          </Notification>
          <SiriExtension />
        </ns6:NotifyGeneralMessage>
      </soap:Body>
    </soap:Envelope>
    """
    Then the SIRI server should have received 2 DeleteSubscription requests

  @ARA-1256
  Scenario: Delete and recreate subscription when receiving subscription with same existing number
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-general-message-subscription-broadcaster] and the following settings:
       | remote_url                         | http://localhost:8090 |
       | remote_credential                  | test                  |
       | local_credential                   | NINOXE:default        |
       | remote_code_space               | internal              |
       | broadcast.subscriptions.persistent | true                  |
      And a Line exists with the following attributes:
        | Name              | Test              |
        | Codes             | "internal":"1234" |
        | CollectSituations | true              |
    And a minute has passed
    When I send this SIRI request
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:Subscribe xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri" xmlns:sws="http://wsdl.siri.org.uk/siri">
      <SubscriptionRequestInfo>
	<siri:RequestTimestamp>2017-01-01T12:00:45.000Z</siri:RequestTimestamp>
	<siri:RequestorRef>NINOXE:default</siri:RequestorRef>
	<siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
      </SubscriptionRequestInfo>
      <Request>
	<siri:GeneralMessageSubscriptionRequest>
	  <siri:SubscriberRef>test</siri:SubscriberRef>
	  <siri:SubscriptionIdentifier>1</siri:SubscriptionIdentifier>
	  <siri:InitialTerminationTime>2017-01-03T12:00:45.000Z</siri:InitialTerminationTime>
	  <siri:GeneralMessageRequest version="2.0:FR-IDF-2.4">
	    <siri:RequestTimestamp>2017-01-01T12:00:45.000Z</siri:RequestTimestamp>
	    <siri:MessageIdentifier>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:MessageIdentifier>
	    <siri:Extensions>
	      <sws:IDFGeneralMessageRequestFilter>
		<siri:LineRef>1234</siri:LineRef>
	      </sws:IDFGeneralMessageRequestFilter>
	    </siri:Extensions>
	  </siri:GeneralMessageRequest>
	</siri:GeneralMessageSubscriptionRequest>
      </Request>
      <RequestExtension/>
    </sw:Subscribe>
  </S:Body>
</S:Envelope>
    """
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Kind            | GeneralMessageBroadcast           |
      | ExternalId      | 1                                 |
    When I send this SIRI request
    """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:Subscribe xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri" xmlns:sws="http://wsdl.siri.org.uk/siri">
      <SubscriptionRequestInfo>
	<siri:RequestTimestamp>2017-01-01T12:01:45.000Z</siri:RequestTimestamp>
	<siri:RequestorRef>NINOXE:default</siri:RequestorRef>
	<siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
      </SubscriptionRequestInfo>
      <Request>
	<siri:GeneralMessageSubscriptionRequest>
	  <siri:SubscriberRef>test</siri:SubscriberRef>
	  <siri:SubscriptionIdentifier>1</siri:SubscriptionIdentifier>
	  <siri:InitialTerminationTime>2017-01-03T12:00:45.000Z</siri:InitialTerminationTime>
	  <siri:GeneralMessageRequest version="2.0:FR-IDF-2.4">
	    <siri:RequestTimestamp>2017-01-01T12:00:45.000Z</siri:RequestTimestamp>
	    <siri:MessageIdentifier>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:MessageIdentifier>
	    <siri:Extensions>
	      <sws:IDFGeneralMessageRequestFilter>
		<siri:LineRef>1234</siri:LineRef>
	      </sws:IDFGeneralMessageRequestFilter>
	    </siri:Extensions>
	  </siri:GeneralMessageRequest>
	</siri:GeneralMessageSubscriptionRequest>
      </Request>
      <RequestExtension/>
    </sw:Subscribe>
  </S:Body>
</S:Envelope>
      """
    Then No Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Kind       | GeneralMessageBroadcast           |
      | ExternalId      | 1                                 |
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Kind       | GeneralMessageBroadcast           |
      | ExternalId      | 1                                 |
    When I send this SIRI request
    """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:Subscribe xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri" xmlns:sws="http://wsdl.siri.org.uk/siri">
      <SubscriptionRequestInfo>
	<siri:RequestTimestamp>2017-01-01T12:01:45.000Z</siri:RequestTimestamp>
	<siri:RequestorRef>NINOXE:default</siri:RequestorRef>
	<siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
      </SubscriptionRequestInfo>
      <Request>
	<siri:GeneralMessageSubscriptionRequest>
	  <siri:SubscriberRef>test</siri:SubscriberRef>
	  <siri:SubscriptionIdentifier>2</siri:SubscriptionIdentifier>
	  <siri:InitialTerminationTime>2017-01-03T12:00:45.000Z</siri:InitialTerminationTime>
	  <siri:GeneralMessageRequest version="2.0:FR-IDF-2.4">
	    <siri:RequestTimestamp>2017-01-01T12:00:45.000Z</siri:RequestTimestamp>
	    <siri:MessageIdentifier>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:MessageIdentifier>
	    <siri:Extensions>
	      <sws:IDFGeneralMessageRequestFilter>
		<siri:LineRef>1234</siri:LineRef>
	      </sws:IDFGeneralMessageRequestFilter>
	    </siri:Extensions>
	  </siri:GeneralMessageRequest>
	</siri:GeneralMessageSubscriptionRequest>
      </Request>
      <RequestExtension/>
    </sw:Subscribe>
  </S:Body>
</S:Envelope>
      """
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Kind            | GeneralMessageBroadcast           |
      | ExternalId      | 1                                 |
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | Kind            | GeneralMessageBroadcast           |
      | ExternalId      | 2                                 |

  @ARA-1378
  Scenario: Returns empty SOAP response on General Message notification
    Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
       """
   <?xml version='1.0' encoding='utf-8'?>
   <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
   <S:Body>
     <ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
       <SubscriptionAnswerInfo
         xmlns:ns2="http://www.ifopt.org.uk/acsb"
         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
         xmlns:ns5="http://www.siri.org.uk/siri"
         xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
         <ns5:Address>http://appli.chouette.mobi/siri_france/siri</ns5:Address>
         <ns5:ResponderRef>NINOXE:default</ns5:ResponderRef>
         <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</ns5:RequestMessageRef>
       </SubscriptionAnswerInfo>
       <Answer
         xmlns:ns2="http://www.ifopt.org.uk/acsb"
         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
         xmlns:ns5="http://www.siri.org.uk/siri"
         xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseStatus>
             <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
             <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</ns5:RequestMessageRef>
             <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
             <ns5:SubscriptionRef>6ba7b814-9dad-11d1-5-00c04fd430c8</ns5:SubscriptionRef>
             <ns5:Status>true</ns5:Status>
             <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
         </ns5:ResponseStatus>
         <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
       </Answer>
       <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
     </ns1:SubscribeResponse>
   </S:Body>
   </S:Envelope>
   """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-general-message-subscription-collector] and the following settings:
      | remote_url                               | http://localhost:8090 |
      | remote_credential                        | test                  |
      | local_credential                         | NINOXE:default        |
      | remote_code_space                     | internal              |
      | collect.filter_general_messages          | true                  |
      | collect.include_lines                    | 1234                  |
      | siri.soap.empty_response_on_notification | true                  |
    And 30 seconds have passed
      And a Line exists with the following attributes:
        | Name              | Test              |
        | Codes             | "internal":"1234" |
        | CollectSituations | true              |
      And 10 seconds have passed
      And 5 seconds have passed
      And a Situation exists with the following attributes:
        | Codes               | "internal" : "NINOXE:GeneralMessage:27_1" |
        | RecordedAt              | 2017-01-01T03:30:06+02:00                 |
        | Version                 | 1                                         |
        | Channel                 | Perturbations                             |
        | ValidUntil              | 2017-01-01T20:30:06+02:00                 |
        | Messages[0]#MessageType | longMessage                               |
        | Messages[0]#MessageText | Les autres non                            |
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
    <S:Body>
      <ns1:NotifyGeneralMessage xmlns:ns1="http://wsdl.siri.org.uk">
       <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns5:ResponseTimestamp>
         <ns5:ProducerRef>NINOXE:default</ns5:ProducerRef>
         <ns5:ResponseMessageIdentifier>NAVINEO:SM:NOT:427843</ns5:ResponseMessageIdentifier>
         <ns5:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns5:RequestMessageRef>
       </ServiceDeliveryInfo>
       <Notification xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
            <ns3:ResponseTimestamp>2017-03-29T16:47:53.039+02:00</ns3:ResponseTimestamp>
            <ns5:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns5:RequestMessageRef>
            <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
            <ns5:SubscriptionRef>6ba7b814-9dad-11d1-5-00c04fd430c8</ns5:SubscriptionRef>
            <ns3:Status>true</ns3:Status>
            <ns3:GeneralMessage>
               <ns3:RecordedAtTime>2017-03-29T03:30:06.000+01:00</ns3:RecordedAtTime>
               <ns3:ItemIdentifier>3477</ns3:ItemIdentifier>
               <ns3:InfoMessageIdentifier>NINOXE:GeneralMessage:27_1</ns3:InfoMessageIdentifier>
               <ns3:InfoMessageVersion>2</ns3:InfoMessageVersion>
               <ns3:formatRef>STIF-IDF</ns3:formatRef>
               <ns3:InfoChannelRef>Commercial</ns3:InfoChannelRef>
               <ns3:ValidUntilTime>2017-03-29T03:30:06.000+01:00</ns3:ValidUntilTime>
               <ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns9:IDFGeneralMessageStructure">
                  <Message>
                    <MessageType>longMessage</MessageType>
                    <MessageText xml:lang="NL">La nouvelle carte d'abonnement est disponible au points de vente du réseau</MessageText>
                  </Message>
                  <LineSection>
                    <FirstStop>NINOXE:StopPoint:SP:24:LOC</FirstStop>
                    <LastStop>NINOXE:StopPoint:SP:12:LOC</LastStop>
                    <LineRef>NINOXE:Line::3:LOC</LineRef>
                  </LineSection>
               </ns3:Content>
            </ns3:GeneralMessage>
         </ns3:GeneralMessageDelivery>
       </Notification>
      <SiriExtension/>
      </ns1:NotifyGeneralMessage>
   </S:Body>
   </S:Envelope>
      """
    Then a Situation exists with the following attributes:
        | Codes | "internal" : "NINOXE:GeneralMessage:27_1" |
        | Channel   | Commercial                                |
    And I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?> 
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
        </S:Body>
      </S:Envelope>
      """

  @ARA-1417 @siri-valid
  Scenario: Log SubscriptionIdentifiers for a GeneralMessageSubscriptionRequest
    Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
       """
   <?xml version='1.0' encoding='utf-8'?>
   <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
   <S:Body>
     <ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
       <SubscriptionAnswerInfo
         xmlns:ns2="http://www.ifopt.org.uk/acsb"
         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
         xmlns:ns5="http://www.siri.org.uk/siri"
         xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
         <ns5:Address>http://appli.chouette.mobi/siri_france/siri</ns5:Address>
         <ns5:ResponderRef>NINOXE:default</ns5:ResponderRef>
         <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</ns5:RequestMessageRef>
       </SubscriptionAnswerInfo>
       <Answer
         xmlns:ns2="http://www.ifopt.org.uk/acsb"
         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
         xmlns:ns5="http://www.siri.org.uk/siri"
         xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseStatus>
             <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
             <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</ns5:RequestMessageRef>
             <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
             <ns5:SubscriptionRef>6ba7b814-9dad-11d1-6-00c04fd430c8</ns5:SubscriptionRef>
             <ns5:Status>true</ns5:Status>
             <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
         </ns5:ResponseStatus>
         <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
       </Answer>
       <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
     </ns1:SubscribeResponse>
   </S:Body>
   </S:Envelope>
   """
      And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-general-message-subscription-collector] and the following settings:
        | remote_url                      | http://localhost:8090 |
        | remote_credential               | test                  |
        | local_credential                | NINOXE:default        |
        | remote_code_space            | internal              |
        | collect.filter_general_messages | true                  |
        | collect.include_lines           | NINOXE:Line::3:LOC    |
      And a Line exists with the following attributes:
        | Name              | Test                            |
        | Codes         | "internal":"NINOXE:Line::3:LOC" |
        | CollectSituations | true                            |
      And a Subscription exist with the following attributes:
        | Kind              | GeneralMessageCollect                  |
        | ReferenceArray[0] | Line, "internal": "NINOXE:Line::3:LOC" |
      And 30 seconds have passed
      And a minute has passed
      And an audit event should exist with these attributes:
      | Protocol                | siri                                  |
      | Direction               | sent                                  |
      | Status                  | OK                                    |
      | Type                    | GeneralMessageSubscriptionRequest     |
      | Lines                   | ["NINOXE:Line::3:LOC"]                |
      | SubscriptionIdentifiers | ["6ba7b814-9dad-11d1-3-00c04fd430c8"] |

  @ARA-1443
  Scenario: Collect GeneralMessage with internal tags
    Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
       """
   <?xml version='1.0' encoding='utf-8'?>
   <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
   <S:Body>
     <ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
       <SubscriptionAnswerInfo
         xmlns:ns2="http://www.ifopt.org.uk/acsb"
         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
         xmlns:ns5="http://www.siri.org.uk/siri"
         xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
         <ns5:Address>http://appli.chouette.mobi/siri_france/siri</ns5:Address>
         <ns5:ResponderRef>NINOXE:default</ns5:ResponderRef>
         <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</ns5:RequestMessageRef>
       </SubscriptionAnswerInfo>
       <Answer
         xmlns:ns2="http://www.ifopt.org.uk/acsb"
         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
         xmlns:ns5="http://www.siri.org.uk/siri"
         xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseStatus>
             <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
             <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</ns5:RequestMessageRef>
             <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
             <ns5:SubscriptionRef>6ba7b814-9dad-11d1-5-00c04fd430c8</ns5:SubscriptionRef>
             <ns5:Status>true</ns5:Status>
             <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
         </ns5:ResponseStatus>
         <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
       </Answer>
       <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
     </ns1:SubscribeResponse>
   </S:Body>
   </S:Envelope>
   """
      And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-general-message-subscription-collector] and the following settings:
        | remote_url                       | http://localhost:8090 |
        | remote_credential                | test                  |
        | local_credential                 | NINOXE:default        |
        | remote_code_space                | internal              |
        | collect.filter_general_messages  | true                  |
        | collect.include_lines            | NINOXE:Line::3:LOC    |
        | collect.situations.internal_tags | first,second          |
      And 30 seconds have passed
      And a Line exists with the following attributes:
        | Name              | Test                            |
        | Codes             | "internal":"NINOXE:Line::4:LOC" |
        | CollectSituations | true                            |
      And a StopArea exists with the following attributes:
        | Name              | Test                                    |
        | Codes             | "internal":"NINOXE:StopPoint:SP:24:LOC" |
        | CollectSituations | true                                    |
      And 10 seconds have passed
      And 5 seconds have passed
      And show me ara subscriptions for partner "test"
      When I send this SIRI request
        """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
    <S:Body>
      <ns1:NotifyGeneralMessage xmlns:ns1="http://wsdl.siri.org.uk">
       <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns5:ResponseTimestamp>
         <ns5:ProducerRef>NINOXE:default</ns5:ProducerRef>
         <ns5:ResponseMessageIdentifier>NAVINEO:SM:NOT:427843</ns5:ResponseMessageIdentifier>
         <ns5:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns5:RequestMessageRef>
       </ServiceDeliveryInfo>
       <Notification xmlns="http://www.siri.org.ukt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
            <ns3:ResponseTimestamp>2017-03-29T16:47:53.039+02:00</ns3:ResponseTimestamp>
            <ns5:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns5:RequestMessageRef>
            <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
            <ns5:SubscriptionRef>6ba7b814-9dad-11d1-5-00c04fd430c8</ns5:SubscriptionRef>
            <ns3:Status>true</ns3:Status>
            <ns3:GeneralMessage>
               <ns3:RecordedAtTime>2017-03-01T03:30:06.000+01:00</ns3:RecordedAtTime>
               <ns3:ItemIdentifier>3477</ns3:ItemIdentifier>
               <ns3:InfoMessageIdentifier>NINOXE:GeneralMessage:27_2</ns3:InfoMessageIdentifier>
               <ns3:InfoMessageVersion>2</ns3:InfoMessageVersion>
               <ns3:formatRef>STIF-IDF</ns3:formatRef>
               <ns3:InfoChannelRef>Commercial</ns3:InfoChannelRef>
               <ns3:ValidUntilTime>2017-03-29T03:30:06.000+01:00</ns3:ValidUntilTime>
               <ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns9:IDFGeneralMessageStructure">
               <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
                  <Message>
                    <MessageType>textOnly</MessageType>
                    <MessageText xml:lang="NL">carte d'abonnement</MessageText>
                  </Message>
                  <LineRef>NINOXE:Line::4:LOC</LineRef>
               </ns3:Content>
            </ns3:GeneralMessage>
         </ns3:GeneralMessageDelivery>
       </Notification>
       <SiriExtension/>
      </ns1:NotifyGeneralMessage>
    </S:Body>
    </S:Envelope>
      """
    Then one Situation has the following attributes:
      | Codes                        | "internal" : "NINOXE:GeneralMessage:27_2" |
      | InternalTags                 | ["first","second"]                        |

  @ARA-1444    
  Scenario: Broadcast a GeneralMessage Notification after modification of a Situation with matching InternalTags
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, siri-general-message-subscription-broadcaster] and the following settings:
       | remote_url                         | http://localhost:8090 |
       | remote_credential                  | test                  |
       | local_credential                   | NINOXE:default        |
       | remote_code_space                  | internal              |
       | broadcast.situations.internal_tags | first,another         |
    And a Subscription exist with the following attributes:
      | Kind              | GeneralMessageBroadcast                     |
      | ExternalId        | externalId                                  |
      | SubscriberRef     | subscriber                                  |
      | ReferenceArray[0] | Situation, "SituationResource": "Situation" |
    And a Line exists with the following attributes:
      | Name              | Test              |
      | Codes             | "internal":"1234" |
      | CollectSituations | true              |
    And a Situation exists with the following attributes:
      | Codes                      | "internal" : "NINOXE:GeneralMessage:27_1" |
      | RecordedAt                 | 2017-01-01T03:30:06+02:00                 |
      | Version                    | 1                                         |
      | InternalTags               | ["first","second"]                        |
      | Keywords                   | ["Perturbation"]                          |
      | ValidityPeriods[0]#EndTime | 2017-01-01T20:30:06+02:00                 |
      | Description[DefaultValue]  | a very very very long message             |
      | Affects[Line]              | 6ba7b814-9dad-11d1-3-00c04fd430c8         |
      | Affects[StopArea]          | 6ba7b814-9dad-11d1-5-00c04fd430c8         |
    And a StopArea exists with the following attributes:
        | Name              | Test                                    |
        | Codes         | "internal":"NINOXE:StopPoint:SP:24:LOC" |
        | CollectSituations | true                                    |
    And 10 seconds have passed
    When the Situation "6ba7b814-9dad-11d1-4-00c04fd430c8" is edited with the following attributes:
      | RecordedAt                 | 2017-01-01T03:50:06+02:00              |
      | ValidityPeriods[0]#EndTime | 2017-10-24T20:30:06+02:00              |
      | Description[DefaultValue]  | an ANOTHER very very very long message |
      | Version                    | 2                                      |
    And 15 seconds have passed
    Then the SIRI server should receive this response
    """
     <?xml version='1.0' encoding='utf-8'?>
     <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
       <S:Body>
         <sw:NotifyGeneralMessage xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
           <ServiceDeliveryInfo>
             <siri:ResponseTimestamp>2017-01-01T12:00:25.000Z</siri:ResponseTimestamp>
             <siri:ProducerRef>test</siri:ProducerRef>
             <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
             <siri:RequestMessageRef></siri:RequestMessageRef>
           </ServiceDeliveryInfo>
           <Notification>
             <siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4" xmlns:stif="http://wsdl.siri.org.uk/siri">
               <siri:ResponseTimestamp>2017-01-01T12:00:25.000Z</siri:ResponseTimestamp>
               <siri:RequestMessageRef></siri:RequestMessageRef>
               <siri:SubscriberRef>subscriber</siri:SubscriberRef>
               <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
               <siri:Status>true</siri:Status>
               <siri:GeneralMessage formatRef="STIF-IDF">
                 <siri:RecordedAtTime>2017-01-01T03:50:06.000+02:00</siri:RecordedAtTime>
                 <siri:ItemIdentifier>RATPDev:Item::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ItemIdentifier>
                 <siri:InfoMessageIdentifier>NINOXE:GeneralMessage:27_1</siri:InfoMessageIdentifier>
                 <siri:InfoMessageVersion>2</siri:InfoMessageVersion>
                 <siri:InfoChannelRef>Perturbation</siri:InfoChannelRef>
                 <siri:ValidUntilTime>2017-10-24T20:30:06.000+02:00</siri:ValidUntilTime>
                 <siri:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                               xsi:type="stif:IDFGeneralMessageStructure">
                   <siri:LineRef>1234</siri:LineRef>
                   <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                   <Message>
                     <MessageType>textOnly</MessageType>
                     <MessageText>an ANOTHER very very very long message</MessageText>
                   </Message>
                 </siri:Content>
               </siri:GeneralMessage>
             </siri:GeneralMessageDelivery>
           </Notification>
           <SiriExtension />
         </sw:NotifyGeneralMessage>
       </S:Body>
     </S:Envelope>
    """

  @ARA-1444    
  Scenario: Do not broadcast a GeneralMessage Notification after modification of a Situation with no matching InternalTags
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, siri-general-message-subscription-broadcaster] and the following settings:
      | remote_url                         | http://localhost:8090 |
      | remote_credential                  | test                  |
      | local_credential                   | NINOXE:default        |
      | remote_code_space                  | internal              |
      | broadcast.situations.internal_tags | first,another         |
    And a Subscription exist with the following attributes:
      | Kind              | GeneralMessageBroadcast                     |
      | ExternalId        | externalId                                  |
      | SubscriberRef     | subscriber                                  |
      | ReferenceArray[0] | Situation, "SituationResource": "Situation" |
    And a Line exists with the following attributes:
      | Name              | Test              |
      | Codes             | "internal":"1234" |
      | CollectSituations | true              |
    And a Situation exists with the following attributes:
      | Codes                      | "internal" : "NINOXE:GeneralMessage:27_1" |
      | RecordedAt                 | 2017-01-01T03:30:06+02:00                 |
      | Version                    | 1                                         |
      | InternalTags               | ["wrong"]                                 |
      | Keywords                   | ["Perturbation"]                          |
      | ValidityPeriods[0]#EndTime | 2017-01-01T20:30:06+02:00                 |
      | Description[DefaultValue]  | a very very very long message             |
      | Affects[Line]              | 6ba7b814-9dad-11d1-3-00c04fd430c8         |
      | Affects[StopArea]          | 6ba7b814-9dad-11d1-5-00c04fd430c8         |
    And a StopArea exists with the following attributes:
      | Name              | Test                                    |
      | Codes             | "internal":"NINOXE:StopPoint:SP:24:LOC" |
      | CollectSituations | true                                    |
    And 10 seconds have passed
    When the Situation "6ba7b814-9dad-11d1-4-00c04fd430c8" is edited with the following attributes:
      | RecordedAt                 | 2017-01-01T03:50:06+02:00              |
      | ValidityPeriods[0]#EndTime | 2017-10-24T20:30:06+02:00              |
      | Description[DefaultValue]  | an ANOTHER very very very long message |
      | Version                    | 2                                      |
    And 15 seconds have passed
    Then the SIRI server should not have received a NotifyGeneralMessage request
