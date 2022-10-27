Feature: Support SIRI VehicleMonitoring by subscription

  Background:
      Given a Referential "test" is created

  Scenario: Create Vehicle Monitoring subscription by Line
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-vehicle-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_objectid_kind               | internal                       |
    And a Line exists with the following attributes:
      | Name      | Test                   |
      | ObjectIDs | "internal": "testLine" |
    When a minute has passed
    And a minute has passed
    And a minute has passed
    Then one Subscription exists with the following attributes:
      | Kind                      | VehicleMonitoringCollect |
      | Resources[0]/SubscribedAt | > 2017-01-01T12:01:00Z     |

  Scenario: Update a Vehicle after a VehicleMonitoringDelivery in a subscription
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-vehicle-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_objectid_kind               | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name      | Test                   |
      | ObjectIDs | "internal": "testLine" |
    And a Subscription exist with the following attributes:
      | Kind              | VehicleMonitoringCollect     |
      | ReferenceArray[0] | Line, "internal": "testLine" |
    And a minute has passed
    When I send this SIRI request
      """
      <?xml version='1.0' encoding='utf-8'?>
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <soap:Body>
          <ns6:NotifyVehicleMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
          xmlns:ns3="http://www.ifopt.org.uk/acsb"
          xmlns:ns4="http://www.ifopt.org.uk/ifopt"
          xmlns:ns5="http://www.siri.org.uk/siri"
          xmlns:ns6="http://wsdl.siri.org.uk"
          xmlns:ns7="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>VehicleMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>
              <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns5:ResponseTimestamp>2022-06-25T15:08:14.940+02:00</ns5:ResponseTimestamp>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns5:Status>true</ns5:Status>
                <ns5:VehicleActivity>
                  <ns5:RecordedAtTime>2022-06-25T15:08:14.928+02:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>108</ns5:ItemIdentifier>
                  <ns5:ValidUntilTime>2022-06-25T16:08:14.928+02:00</ns5:ValidUntilTime>
                  <ns5:VehicleMonitoringRef>108</ns5:VehicleMonitoringRef>
                  <ns5:ProgressBetweenStops>
                    <ns5:LinkDistance>340.0</ns5:LinkDistance>
                    <ns5:Percentage>73.0</ns5:Percentage>
                  </ns5:ProgressBetweenStops>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>C</ns5:LineRef>
                    <ns5:DirectionRef>Aller</ns5:DirectionRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>NAVINEO:DataFrame::1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>RDMANTOIS:VehicleJourney::6628652:LOC</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>RDMANTOIS:JourneyPattern::LCP37:LOC</ns5:JourneyPatternRef>
                    <ns5:JourneyPatternName>LCP37</ns5:JourneyPatternName>
                    <ns5:PublishedLineName>C</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>OPERYORDM:Operator::OPERYORDM:LOC</ns5:OperatorRef>
                    <ns5:OriginRef>50000037</ns5:OriginRef>
                    <ns5:OriginName>Port Fouquet</ns5:OriginName>
                    <ns5:DestinationRef>50000031</ns5:DestinationRef>
                    <ns5:DestinationName>Mantes la Jolie Gare routière - Quai 2</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:VehicleLocation srsName="2154">
                      <ns5:Coordinates>603204 6878517</ns5:Coordinates>
                    </ns5:VehicleLocation>
                    <ns5:Bearing>171.0</ns5:Bearing>
                    <ns5:VehicleRef>108</ns5:VehicleRef>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>50000016</ns5:StopPointRef>
                      <ns5:Order>9</ns5:Order>
                      <ns5:StopPointName>Hôpital F. Quesnay</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>MantesLJ Gare</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2022-06-25T15:05:00.000+02:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2022-06-25T15:08:27.000+02:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2022-06-25T15:05:00.000+02:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2022-06-25T15:08:27.000+02:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                  <ns5:Extensions/>
                </ns5:VehicleActivity>
              </ns5:VehicleMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </ns6:NotifyVehicleMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
      And I see ara vehicles
      Then one Vehicle has the following attributes:
      | ObjectIDs | "internal": "108"                 |
      | LineId    | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
      | Bearing   | 171.0                             |
      | Latitude  | 48.99927561424598                 |
      | Longitude | 1.6770970859674874                |
