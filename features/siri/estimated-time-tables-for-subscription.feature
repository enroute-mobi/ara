Feature: Support SIRI EstimatedTimeTable by subscription

  Background:
    Given a Referential "test" is created

  @wip
    Scenario: 4233 - Manage a ETT Subscription

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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-collector] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
    And a Subscription exist with the following attributes:
      | Kind                      | EstimatedTimetable          |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                          | Passage 32                              |
      | ObjectIDs                     | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                        | 6ba7b814-9dad-11d1-6-00c04fd430c8       |
      | Attribute[DirectionRef]       | Aller                                   |
      | Attribute[OriginName]         | Le début                                |
      | Attribute[DestinationName]    | La fin.                                 |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z                                        |
      | Schedule[expected]#Arrival      | 2017-01-01T15:00:00.000Z                                        |
      | ArrivalStatus                   | onTime                                                              |
    And a minute has passed
    When I send this SIRI request
      """
      <ns1:NotifyEstimatedTimetable xmlns:ns1="http://wsdl.siri.org.uk">
       <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns5:ResponseTimestamp>
         <ns5:ProducerRef>Edwig</ns5:ProducerRef>
         <ns5:ResponseMessageIdentifier>NAVINEO:SM:NOT:427843</ns5:ResponseMessageIdentifier>
         <ns5:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns5:RequestMessageRef>
       </ServiceDeliveryInfo>
       <Notification xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns3:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
            <ns3:ResponseTimestamp>2017-03-29T16:47:53.039+02:00</ns3:ResponseTimestamp>
            <ns5:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns5:RequestMessageRef>
            <ns5:SubscriberRef>RATPDEV:Concerto</ns5:SubscriberRef>
            <ns5:SubscriptionRef>Edwig:Subscription::6ba7b814-9dad-11d1-38-00c04fd430c8:LOC</ns5:SubscriptionRef>
            <ns3:Status>true</ns3:Status>
            <ns3:EstimatedTimetable>
              <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
              <ns3:RequestMessageRef>EstimatedTimetable:Test:0</ns3:RequestMessageRef>
              <ns3:Status>true</ns3:Status>
              <ns3:EstimatedJourneyVersionFrame>
                <ns3:RecordedAtTime>2017-01-01T12:00:00.000Z</ns3:RecordedAtTime>
                <ns3:EstimatedVehicleJourney>
                  <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                  <ns3:DirectionRef>Aller</ns3:DirectionRef>
                  <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                  <ns3:EstimatedCalls>
                    <ns3:EstimatedCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:AimedArrivalTime>2017-01-01T15:00:00.000Z</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-01T15:01:00.000Z</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>Delayed</ns3:ArrivalStatus>
                    </ns3:EstimatedCall>
                  </ns3:EstimatedCalls>
                </ns3:EstimatedVehicleJourney>
              </ns3:EstimatedJourneyVersionFrame>
            </ns3:EstimatedTimetable>
         </ns3:EstimatedTimetableDelivery>
       </Notification>
      <NotifyExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
      </ns1:NotifyEstimatedTimetable>
      """
    Then a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z                                        |
      | Schedule[expected]#Arrival      | 2017-01-01T15:01:00.000Z                                        |
      | ArrivalStatus                   | Delayed                                                              |
