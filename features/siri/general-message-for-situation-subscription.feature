Feature: Support SIRI GeneralMessage by subscription

  Background:
      Given a Referential "test" is created

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
         <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</ns5:RequestMessageRef>
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
      And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-general-message-subscription-collector] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
      And 30 seconds have passed
      And a minute has passed
      And a Situation exists with the following attributes:
       | ObjectIDs               | "internal" : "NINOXE:GeneralMessage:27_1" |
       | RecordedAt              | 2017-01-01T03:30:06+02:00                 |
       | Version                 | 1                                         |
       | Channel                 | Perturbations                             |
       | ValidUntil              | 2017-01-01T20:30:06+02:00                 |
       | Messages[0]#MessageType | longMessage                               |
       | Messages[0]#MessageText | Les autres non                            |
      And a minute has passed
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
      <NotifyExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
      </ns1:NotifyGeneralMessage>
   </S:Body>
   </S:Envelope>
      """
    Then the Situation "6ba7b814-9dad-11d1-8-00c04fd430c8" has the following attributes:
       | Channel | Commercial                        |

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
          <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
        </ns1:SubscribeResponse>
      </S:Body>
      </S:Envelope>
    """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-general-message-subscription-collector-monitoring-subscription-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_objectid_kind | internal              |
    And 30 seconds have passed
    And a minute has passed
    And a Subscription exist with the following attributes:
      | Kind | GeneralMessage |
    And a Situation exists with the following attributes:
      | ObjectIDs               | "internal" : "2"          |
      | RecordedAt              | 2017-01-01T03:30:06+02:00 |
      | Version                 | 1                         |
      | Channel                 | Perturbations             |
      | ValidUntil              | 2017-01-01T20:30:06+02:00 |
      | Messages[0]#MessageType | longMessage               |
      | Messages[0]#MessageText | Les autres non            |
    When I send this SIRI request
    """
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
              <ns2:SubscriptionRef>Edwig:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</ns2:SubscriptionRef>
              <ns2:Status>true</ns2:Status>
              <ns2:GeneralMessageCancellation>
                <ns2:RecordedAtTime>2017-05-15T13:26:10.116+02:00</ns2:RecordedAtTime>
                <ns2:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns2:ItemRef>
              </ns2:GeneralMessageCancellation>
            </ns2:GeneralMessageDelivery>
          </Notification>
          <SiriExtension />
        </ns6:NotifyGeneralMessage>
      </soap:Body>
    </soap:Envelope>
    """
    Then a Situation "internal":"2" should not exist in Referential "test"
