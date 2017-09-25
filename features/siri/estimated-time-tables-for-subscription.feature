Feature: Support SIRI EstimatedTimeTable by subscription

  Background:
    Given a Referential "test" is created

  Scenario: 4234 - Handle a SIRI EstimatedTimeTable request for subscription
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-collector] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
    And a minute has passed
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
          xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:EstimatedTimetableSubscriptionRequest xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:ns3="http://www.ifopt.org.uk/acsb"
                                 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                                 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                                 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <SubscriptionRequestInfo>
              <ns2:RequestTimestamp>2017-01-01T12:01:00.000Z</ns2:RequestTimestamp>
              <ns2:RequestorRef>NINOXE:default</ns2:RequestorRef>
              <ns2:MessageIdentifier>ETTSubscription:Test:0</ns2:MessageIdentifier>
              <ns2:ConsumerAddress>https://edwig-staging.af83.io/test/siri</ns2:ConsumerAddress>
            </SubscriptionRequestInfo>
            <Request version="2.0:FR-IDF-2.4">
              <EstimatedTimetableSubscriptionRequest>
                <ns2:RequestTimestamp>2017-01-01T12:01:00.000Z</ns2:RequestTimestamp>
                <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
                <ns5:SubscriptionRef>NINOXE:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ns5:SubscriptionRef>
                <MessageIdentifier>28679112-9dad-11d1-2-00c04fd430c8</MessageIdentifier>
                <ns5:InitialTerminationTime>2017-01-01T13:00:00.000Z</ns5:InitialTerminationTime>
              </EstimatedTimetableSubscriptionRequest>
            </Request>
            <RequestExtension />
          </ns7:EstimatedTimetableSubscriptionRequest>
        </S:Body>
      </S:Envelope>
      """
    Then a Subscription exist with the following attributes:
      | Kind                      | EstimatedTimetable          |


  Scenario: 4235 - Manage a ETT Notify after modification of a StopVisit
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-estimated-timetable-subscription-broadcaster] and the following settings:
       | remote_url           | http://localhost:8090 |
       | remote_credential    | test                  |
       | local_credential     | NINOXE:default        |
       | remote_objectid_kind | internal              |
    And a Subscription exist with the following attributes:
      | Kind                      | EstimatedTimetable                     |
      | ReferenceArray[0]           | Line, "internal": "NINOXE:Line:3:LOC"  |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                          | Passage 32                              |
      | ObjectIDs                     | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                        | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
      | Attribute[DirectionRef]       | Aller                                   |
      | Attribute[OriginName]         | Le début                                |
      | Attribute[DestinationName]    | La fin.                                 |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | false                                                                |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000Z                                             |
      | Schedule[expected]#Arrival      | 2017-01-01T15:00:00.000Z                                             |
      | ArrivalStatus                   | onTime                                                               |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival      | 2017-01-01T15:01:01.000Z                                             |
      | ArrivalStatus                   | Delayed                                                              |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
<sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>test</siri:ProducerRef>
		<siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
	</ServiceDeliveryInfo>
	<Notification>
		<siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
			<siri:RequestMessageRef></siri:RequestMessageRef>
			<siri:SubscriberRef>test</siri:SubscriberRef>
			<siri:SubscriptionRef>6ba7b814-9dad-11d1-2-00c04fd430c8</siri:SubscriptionRef>
			<siri:Status>true</siri:Status>
			<siri:EstimatedJourneyVersionFrame>
				<siri:RecordedAtTime>2017-01-01T12:00:00.000Z</siri:RecordedAtTime>
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
					<siri:DirectionRef>Aller</siri:DirectionRef>
					<siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
					<siri:EstimatedCalls>
						<siri:EstimatedCall>
							<siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
							<siri:Order>4</siri:Order>
							<siri:StopPointName>Test</siri:StopPointName>
							<siri:VehicleAtStop>false</siri:VehicleAtStop>
							<siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
							<siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
						</siri:EstimatedCall>
					</siri:EstimatedCalls>
				</siri:EstimatedVehicleJourney>
			</siri:EstimatedJourneyVersionFrame>
		</siri:EstimatedTimetableDelivery>
	</Notification>
	<NotificationExtension/>
</sw:NotifyEstimatedTimetable>
</S:Body>
</S:Envelope>
      """
