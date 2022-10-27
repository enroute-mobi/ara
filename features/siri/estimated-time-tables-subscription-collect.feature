Feature: Support SIRI EstimatedTimetable

  Background:
      Given a Referential "test" is created

  @ARA-1152
  Scenario: Create EstimatedTimetable subscription collect
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
        <ns5:RequestMessageRef>response</ns5:RequestMessageRef>
      </SubscriptionAnswerInfo>
      <Answer
        xmlns:ns2="http://www.ifopt.org.uk/acsb"
        xmlns:ns3="http://www.ifopt.org.uk/ifopt"
        xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
        xmlns:ns5="http://www.siri.org.uk/siri"
        xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseStatus>
            <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
            <ns5:RequestMessageRef>{LastRequestMessageRef}</ns5:RequestMessageRef>
            <ns5:SubscriberRef>test</ns5:SubscriberRef>
            <ns5:SubscriptionRef>6ba7b814-9dad-11d1-5-00c04fd430c8</ns5:SubscriptionRef>
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_objectid_kind | internal              |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Test                   |
      | ObjectIDs | "internal": "testLine" |
    And a minute has passed
    Then one Subscription exists with the following attributes:
      | Kind | EstimatedTimetableCollect |

#   @ARA-1152
#   Scenario: Update ara models after a EstimatedTimetableNotify in a subscription
#     Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
#       """
#   <?xml version='1.0' encoding='utf-8'?>
#   <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
#   <S:Body>
#     <ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
#       <SubscriptionAnswerInfo
#         xmlns:ns2="http://www.ifopt.org.uk/acsb"
#         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
#         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
#         xmlns:ns5="http://www.siri.org.uk/siri"
#         xmlns:ns6="http://wsdl.siri.org.uk/siri">
#         <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
#         <ns5:Address>http://appli.chouette.mobi/siri_france/siri</ns5:Address>
#         <ns5:ResponderRef>NINOXE:default</ns5:ResponderRef>
#         <ns5:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns5:MessageRefStructure">Subscription:Test:0</ns5:RequestMessageRef>
#       </SubscriptionAnswerInfo>
#       <Answer
#         xmlns:ns2="http://www.ifopt.org.uk/acsb"
#         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
#         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
#         xmlns:ns5="http://www.siri.org.uk/siri"
#         xmlns:ns6="http://wsdl.siri.org.uk/siri">
#         <ns5:ResponseStatus>
#             <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
#             <ns5:RequestMessageRef>Subscription:Test:0</ns5:RequestMessageRef>
#             <ns5:SubscriberRef>SubscriberRef</ns5:SubscriberRef>
#             <ns5:SubscriptionRef>SubscriptionIdentifier</ns5:SubscriptionRef>
#             <ns5:Status>true</ns5:Status>
#             <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
#         </ns5:ResponseStatus>
#         <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
#       </Answer>
#       <AnswerExtension/>
#     </ns1:SubscribeResponse>
#   </S:Body>
#   </S:Envelope>
#       """
#     And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
#       | remote_url           | http://localhost:8090 |
#       | remote_credential    | test                  |
#       | local_credential     | NINOXE:default        |
#       | remote_objectid_kind | internal              |
#     And 30 seconds have passed
#     And a Line exists with the following attributes:
#       | Name      | Test                            |
#       | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
#     And a Subscription exist with the following attributes:
#       | Kind              | EstimatedTimetableCollect             |
#       | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
#     And a minute has passed
#     When I send this SIRI request
#       """
# <?xml version='1.0' encoding='utf-8'?>
# <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
# <S:Body>
# <sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
#   <ServiceDeliveryInfo>
#     <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
#     <siri:ProducerRef>test</siri:ProducerRef>
#     <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
#     <siri:RequestMessageRef></siri:RequestMessageRef>
#   </ServiceDeliveryInfo>
#   <Notification>
#     <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
#       <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
#       <siri:RequestMessageRef></siri:RequestMessageRef>
#       <siri:SubscriberRef>subscriber</siri:SubscriberRef>
#       <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
#       <siri:Status>true</siri:Status>
#       <siri:EstimatedJourneyVersionFrame>
#         <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
#         <siri:EstimatedVehicleJourney>
#           <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
#           <siri:DirectionRef>Aller</siri:DirectionRef>
#           <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
#           <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
#           <siri:DestinationRef>ThisIsTheEnd</siri:DestinationRef>
#           <siri:EstimatedCalls>
#             <siri:EstimatedCall>
#               <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
#               <siri:Order>4</siri:Order>
#               <siri:StopPointName>Test</siri:StopPointName>
#               <siri:VehicleAtStop>false</siri:VehicleAtStop>
#               <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
#               <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
#             </siri:EstimatedCall>
#           </siri:EstimatedCalls>
#         </siri:EstimatedVehicleJourney>
#       </siri:EstimatedJourneyVersionFrame>
#     </siri:EstimatedTimetableDelivery>
#   </Notification>
#   <NotifyExtension/>
# </sw:NotifyEstimatedTimetable>
# </S:Body>
# </S:Envelope>
#       """
#     And 1 minute has passed
#     And I see ara vehicle_journeys
#     Then one VehicleJourney has the following attributes:
#       | ObjectIDs     | "internal": "NINOXE:VehicleJourney:201" |
#       | LineId        | 6ba7b814-9dad-11d1-1-00c04fd430c8       |
#       | DirectionType | Aller                                   |
#     And I see ara stop_areas
#     And one StopArea has the following attributes:
#       | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
#       | Name      | Test                                     |
#     And I see ara stop_visits
#     And one StopVisit has the following attributes:
#       | PassageOrder               | 4                        |
#       | VehicleAtStop              | false                    |
#       | ArrivalStatus              | Delayed                  |
#       | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
#       | VehicleJourneyId           | ????                     |
#       | StopAreaId                 | ????                     |
