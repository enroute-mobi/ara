Feature: Support SIRI GeneralMessage by subscription

  Background:
      Given a Referential "test" is created

@wip
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
              <ns5:ResponseTimestamp>2017-01-01T12:01:00.000+02:00</ns5:ResponseTimestamp>
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
                  <ns5:ResponseTimestamp>2017-01-01T12:01:00.000+02:00</ns5:ResponseTimestamp>
                  <ns5:RequestMessageRef>GMSubscription:Test:0</ns5:RequestMessageRef>
                  <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
                  <ns5:SubscriptionRef>NINOXE:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ns5:SubscriptionRef>
                  <ns5:Status>true</ns5:Status>
                  <ns5:ValidUntil>2017-02-01T12:00:00.000+02:00</ns5:ValidUntil>
              </ns5:ResponseStatus>
              <ns5:ServiceStartedTime>2017-01-01T12:01:00.000+02:00</ns5:ServiceStartedTime>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:SubscribeResponse>
        </S:Body>
        </S:Envelope>
        """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-subscription-collector] and the following settings:
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
          <ns7: GeneralMessageSubscriptionRequest xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:ns3="http://www.ifopt.org.uk/acsb"
                                 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                                 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                                 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <SubscriptionRequestInfo>
              <ns2:RequestTimestamp>2017-01-01T12:00:00.000+02:00</ns2:RequestTimestamp>
              <ns2:RequestorRef>test</ns2:RequestorRef>
              <ns2:MessageIdentifier>GMSubscription:Test:0</ns2:MessageIdentifier>
            </SubscriptionRequestInfo>

            <Request version="2.0:FR-IDF-2.4">
              <ns2:RequestTimestamp>2017-01-01T12:00:00.000+02:00</ns2:RequestTimestamp>
              <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
              <ns5:SubscriptionRef>NINOXE:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ns5:SubscriptionRef>
              <ns5:InitialTerminationTime>2017-02-01T12:00:00.000+02:00</ns5:InitialTerminationTime>
            </Request>
            <RequestExtension />
          </ns7:GetStopMonitoring>
        </S:Body>
      </S:Envelope>
        """
    Then Then I should receive this SIRI response
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
          <ns5:ResponseTimestamp>2017-01-01T12:01:00.000+02:00</ns5:ResponseTimestamp>
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
              <ns5:ResponseTimestamp>2017-01-01T12:01:00.000+02:00</ns5:ResponseTimestamp>
              <ns5:RequestMessageRef>GMSubscription:Test:0</ns5:RequestMessageRef>
              <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
              <ns5:SubscriptionRef>NINOXE:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ns5:SubscriptionRef>
              <ns5:Status>true</ns5:Status>
              <ns5:ValidUntil>2017-02-01T12:00:00.000+02:00</ns5:ValidUntil>
          </ns5:ResponseStatus>
          <ns5:ServiceStartedTime>2017-01-01T12:01:00.000+02:00</ns5:ServiceStartedTime>
        </Answer>
        <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
      </ns1:SubscribeResponse>
    </S:Body>
    </S:Envelope>
    """

@wip
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
        And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-subscription-collector] and the following settings:
          | remote_url           | http://localhost:8090 |
          | remote_credential    | test                  |
          | local_credential     | NINOXE:default        |
          | remote_objectid_kind | internal              |
          And 30 seconds have passed
        And a StopArea exists with the following attributes:
          | Name      | Test                                     |
          | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
        And a Line exists with the following attributes:
          | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
          | Name      | Ligne 3 Metro                   |
        And a VehicleJourney exists with the following attributes:
          | Name      | Passage 32                              |
          | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
          | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
        And 15 seconds have passed
        And a StopVisit exists with the following attributes:
          | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
          | PassageOrder                    | 4                                                                    |
          | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
          | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
          | VehicleAtStop                   | true                                                                 |
          | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
          | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                             |
          | DepartureStatus                 | onTime                                                               |
        And a minute has passed
        When I send this SIRI request
          """
          <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
            <soap:Body>
              <ns6:NotifyStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
              xmlns:ns3="http://www.ifopt.org.uk/acsb"
              xmlns:ns4="http://www.ifopt.org.uk/ifopt"
              xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
              xmlns:ns6="http://wsdl.siri.org.uk"
              xmlns:ns7="http://wsdl.siri.org.uk/siri">
                <ServiceDeliveryInfo>
                  <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                  <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
                  <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
                  <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                </ServiceDeliveryInfo>

                <Notification>
                  <ns2:StopMonitoringDelivery version="1.3">
                    <ns2:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns2:ResponseTimestamp>
                    <ns2:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns2:RequestMessageRef>
                    <ns2:SubscriberRef>RATPDEV:Concerto</ns2:SubscriberRef>
                    <ns2:SubscriptionRef>Edwig:Subscription::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</ns2:SubscriptionRef>
                    <ns2:Status>true</ns2:Status>
                    <ns2:MonitoredStopVisitCancellation>
                      <ns2:RecordedAtTime>2017-05-15T13:26:10.116+02:00</ns2:RecordedAtTime>
                      <ns2:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns2:ItemRef>
                      <ns2:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns2:MonitoringRef>
                      <ns2:LineRef>NINOXE:Line:3:LOC</ns2:LineRef>
                    </ns2:MonitoredStopVisitCancellation>
                  </ns2:StopMonitoringDelivery>
                </Notification>
                <SiriExtension />
              </ns6:NotifyStopMonitoring>
            </soap:Body>
          </soap:Envelope>
          """
        Then the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" has the following attributes:
          | DepartureStatus | departed  |
          | ArrivalStatus   | cancelled |

