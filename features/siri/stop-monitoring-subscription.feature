Feature: Support SIRI StopMonitoring by subscription

  Background:
    Given a Referential "test" is created

  @ARA-1101
  Scenario: Update VehicleJourney after a StopMonitoringDelivery in a subscription using the partner setting siri.direction_type should update the DirectionRef
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_code_space                  | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
      | siri.direction_type                | Aller,Retour                   |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                              |
      | Codes                    | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                   |       6ba7b814-9dad-11d1-a-00c04fd430c8 |
      | Monitored                | true                                    |
      | Attribute[DirectionName] | A Direction Name                        |
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
              <ns2:ResponseTimestamp>
              2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>
              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:DirectionRef>Aller</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OperatorRef>NINOXE:Company:15563880:LOC</ns3:OperatorRef>
                    <ns3:ProductCategoryRef>0</ns3:ProductCategoryRef>
                    <ns3:VehicleFeatureRef>TRFC_M4_1</ns3:VehicleFeatureRef>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:42:LOC</ns3:OriginRef>
                    <ns3:OriginName>Magicien Noir</ns3:OriginName>
                    <ns3:DestinationRef>NINOXE:StopPoint:SP:62:LOC</ns3:DestinationRef>
                    <ns3:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</ns3:OriginAimedDepartureTime>
                    <ns3:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</ns3:DestinationAimedArrivalTime>
                    <ns3:Monitored>true</ns3:Monitored>
                    <ns3:ProgressRate>normalProgress</ns3:ProgressRate>
                    <ns3:Delay>P0Y0M0DT0H0M0.000S</ns3:Delay>
                    <ns3:CourseOfJourneyRef>201</ns3:CourseOfJourneyRef>
                    <ns3:VehicleRef>NINOXE:Vehicle:23:LOC</ns3:VehicleRef>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R)</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T13:00:00.000+02:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>delayed</ns3:ArrivalStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns2:StopMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </ns6:NotifyStopMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then the StopVisit "6ba7b814-9dad-11d1-a-00c04fd430c8" has the following attributes:
      | Collected     | true  |
      | PassageOrder  |     4 |
      | VehicleAtStop | false |
    Then the VehicleJourney "6ba7b814-9dad-11d1-5-00c04fd430c8" has the following attributes:
      | DirectionType                          | inbound                       |
      | Attribute[VehicleFeatureRef]           | TRFC_M4_1                     |
      | Attribute[Delay]                       | P0Y0M0DT0H0M0.000S            |
      | Attribute[DestinationAimedArrivalTime] | 2016-09-22T08:02:00.000+02:00 |
      | Attribute[DirectionName]               | Mago-Cime OMNI                |
      | Attribute[OriginAimedDepartureTime]    | 2016-09-22T07:50:00.000+02:00 |
      | Attribute[ProductCategoryRef]          | 0                             |

  @ARA-1150
  Scenario: Update a StopVisit and all VehicleJourney Attributes after a StopMonitoringDelivery in a subscription
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_code_space                  | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                              |
      | Codes                    | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                   |       6ba7b814-9dad-11d1-a-00c04fd430c8 |
      | Monitored                | true                                    |
      | Attribute[DirectionName] | A Direction Name                        |
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
              <ns2:ResponseTimestamp>
              2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>
              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:DirectionRef>Left</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OperatorRef>NINOXE:Company:15563880:LOC</ns3:OperatorRef>
                    <ns3:ProductCategoryRef>0</ns3:ProductCategoryRef>
                    <ns3:VehicleFeatureRef>TRFC_M4_1</ns3:VehicleFeatureRef>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:42:LOC</ns3:OriginRef>
                    <ns3:OriginName>Magicien Noir</ns3:OriginName>
                    <ns3:DestinationRef>NINOXE:StopPoint:SP:62:LOC</ns3:DestinationRef>
                    <ns3:DestinationName>Cimetière des Sauvages</ns3:DestinationName>
                    <ns3:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</ns3:OriginAimedDepartureTime>
                    <ns3:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</ns3:DestinationAimedArrivalTime>
                    <ns3:Monitored>true</ns3:Monitored>
                    <ns3:ProgressRate>normalProgress</ns3:ProgressRate>
                    <ns3:Delay>P0Y0M0DT0H0M0.000S</ns3:Delay>
                    <ns3:CourseOfJourneyRef>201</ns3:CourseOfJourneyRef>
                    <ns3:VehicleRef>NINOXE:Vehicle:23:LOC</ns3:VehicleRef>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R)</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T13:00:00.000+02:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>delayed</ns3:ArrivalStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns2:StopMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </ns6:NotifyStopMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then the StopVisit "6ba7b814-9dad-11d1-a-00c04fd430c8" has the following attributes:
      | Collected     | true  |
      | PassageOrder  |     4 |
      | VehicleAtStop | false |
    Then the VehicleJourney "6ba7b814-9dad-11d1-5-00c04fd430c8" has the following attributes:
      | OriginName                             | Magicien Noir                 |
      | DestinationName                        | Cimetière des Sauvages        |
      | Attribute[DirectionName]               | Mago-Cime OMNI                |
      | Attribute[Delay]                       | P0Y0M0DT0H0M0.000S            |
      | Attribute[DestinationAimedArrivalTime] | 2016-09-22T08:02:00.000+02:00 |
      | Attribute[OriginAimedDepartureTime]    | 2016-09-22T07:50:00.000+02:00 |
      | Attribute[ProductCategoryRef]          | 0                             |
      | Attribute[VehicleFeatureRef]           | TRFC_M4_1                     |

  @ARA-1200
  Scenario: Do not update a VehicleJourney DirectionType after a StopMonitoringDelivery having no DirectionRef tag
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_code_space                  | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Name          | Passage 32                              |
      | Codes         | "internal": "NINOXE:VehicleJourney:201" |
      | LineId        |       6ba7b814-9dad-11d1-a-00c04fd430c8 |
      | Monitored     | true                                    |
      | DirectionType | Aller                                   |
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
              <ns2:ResponseTimestamp>
              2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>
              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:VehicleRef>NINOXE:Vehicle:23:LOC</ns3:VehicleRef>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R)</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T13:00:00.000+02:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>delayed</ns3:ArrivalStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns2:StopMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </ns6:NotifyStopMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then the VehicleJourney "6ba7b814-9dad-11d1-5-00c04fd430c8" has the following attributes:
      | DirectionType | Aller |

  Scenario: 3258 - Update a StopVisit after a StopMonitoringDelivery in a subscription
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_code_space                  | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
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
              <ns2:ResponseTimestamp>
              2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>
              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:DirectionRef>Left</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OperatorRef>NINOXE:Company:15563880:LOC</ns3:OperatorRef>
                    <ns3:ProductCategoryRef>0</ns3:ProductCategoryRef>
                    <ns3:VehicleFeatureRef>TRFC_M4_1</ns3:VehicleFeatureRef>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:42:LOC</ns3:OriginRef>
                    <ns3:OriginName>Magicien Noir</ns3:OriginName>
                    <ns3:DestinationRef>NINOXE:StopPoint:SP:62:LOC</ns3:DestinationRef>
                    <ns3:DestinationName>Cimetière des Sauvages</ns3:DestinationName>
                    <ns3:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</ns3:OriginAimedDepartureTime>
                    <ns3:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</ns3:DestinationAimedArrivalTime>
                    <ns3:Monitored>true</ns3:Monitored>
                    <ns3:ProgressRate>normalProgress</ns3:ProgressRate>
                    <ns3:Delay>P0Y0M0DT0H0M0.000S</ns3:Delay>
                    <ns3:CourseOfJourneyRef>201</ns3:CourseOfJourneyRef>
                    <ns3:VehicleRef>NINOXE:Vehicle:23:LOC</ns3:VehicleRef>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R)</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T13:00:00.000+02:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>delayed</ns3:ArrivalStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns2:StopMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </ns6:NotifyStopMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then the StopVisit "6ba7b814-9dad-11d1-a-00c04fd430c8" has the following attributes:
      | Collected     | true  |
      | PassageOrder  |     4 |
      | VehicleAtStop | false |
    And an audit event should exist with these attributes:
      | Protocol           | siri                           |
      | Direction          | received                       |
      | ResponseIdentifier | /{uuid}/                       |
      | Status             | OK                             |
      | Type               | NotifyStopMonitoring           |
      | StopAreas          | ["NINOXE:StopPoint:SP:24:LOC"] |
      | VehicleJourneys    | ["NINOXE:VehicleJourney:201"]  |
      | Lines              | ["NINOXE:Line:3:LOC"]          |

  Scenario: 3737 - Manage a MonitoredStopVisitCancellation
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
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090       |
      | remote_credential                  | test                        |
      | local_credential                   | NINOXE:default              |
      | remote_code_space                  | internal                    |
      | generators.subscription_identifier | Ara:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopVisit exists with the following attributes:
      | Codes                    | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder             |                                                                    4 |
      | StopAreaId               |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Schedule[actual]#Arrival |                                             2017-01-01T13:00:00.000Z |
      | DepartureStatus          | onTime                                                               |
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
                <ns2:SubscriptionRef>Ara:Subscription::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</ns2:SubscriptionRef>
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
      | DepartureStatus | departed |
      | ArrivalStatus   | arrived  |

  Scenario: Handle multiple StopAreas in Subscription
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
            <ns5:RequestMessageRef>Ara:Subscription::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</ns5:RequestMessageRef>
            <ns5:SubscriberRef>SubscriberRef</ns5:SubscriberRef>
            <ns5:SubscriptionRef>6ba7b814-9dad-11d1-5-00c04fd430c8SubscriptionIdentifier</ns5:SubscriptionRef>
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
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | test                  |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name  | Test 1                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test 2                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:25:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | ReferenceArray[1] | StopArea, "internal": "NINOXE:StopPoint:SP:25:LOC" |
    When a minute has passed
    And a minute has passed
    Then the SIRI server should have received a SubscriptionRequest request with 2 "StopMonitoringRequest"

  Scenario: 3737 - Manage a MonitoredStopVisitCancellation
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
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name   | Passage 32                              |
      | Codes  | "internal": "NINOXE:VehicleJourney:201" |
      | LineId |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
    And a StopVisit exists with the following attributes:
      | Codes                        | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                 |                                                                    4 |
      | StopAreaId                   |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId             |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop                | true                                                                 |
      | Reference[OperatorRef]#Code  | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[expected]#Arrival   |                                             2017-01-01T13:00:00.000Z |
      | Schedule[actual]#Arrival     |                                             2017-01-01T13:00:00.000Z |
      | Schedule[expected]#Departure |                                             2017-01-01T13:00:00.000Z |
      | DepartureStatus              | cancelled                                                            |
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
              <ns2:ResponseTimestamp>2017-05-15T13:30:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>

            <Notification>
              <ns2:StopMonitoringDelivery version="1.3">
                <ns2:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RATPDEV:Concerto</ns2:SubscriberRef>
                <ns2:SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</ns2:SubscriptionRef>
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
      | Collected                    | false                         |
      | ArrivalStatus                | arrived                       |
      | DepartureStatus              | cancelled                     |
      | Schedule[expected]#Arrival   |          2017-01-01T13:00:00Z |
      | Schedule[expected]#Departure |          2017-01-01T13:00:00Z |
      | Schedule[actual]#Arrival     |          2017-01-01T13:00:00Z |
      | Schedule[actual]#Departure   | 2017-05-15T13:26:10.116+02:00 |

  Scenario: 3737 - Manage a MonitoredStopVisitCancellation without RecordedAtTime
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
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name   | Passage 32                              |
      | Codes  | "internal": "NINOXE:VehicleJourney:201" |
      | LineId |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
    And a StopVisit exists with the following attributes:
      | Codes                        | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                 |                                                                    4 |
      | StopAreaId                   |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId             |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop                | true                                                                 |
      | Reference[OperatorRef]#Code  | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[expected]#Arrival   |                                             2017-01-01T13:00:00.000Z |
      | Schedule[actual]#Arrival     |                                             2017-01-01T13:00:00.000Z |
      | Schedule[expected]#Departure |                                             2017-01-01T13:00:00.000Z |
      | DepartureStatus              | cancelled                                                            |
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
              <ns2:ResponseTimestamp>2017-05-15T13:30:10.116+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>

            <Notification>
              <ns2:StopMonitoringDelivery version="1.3">
                <ns2:ResponseTimestamp>2017-06-19T16:04:25.983+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>RATPDev:Message::f9c8aa9e-df4d-4a8e-9e25-61f717f13e12:LOC</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RATPDEV:Concerto</ns2:SubscriberRef>
                <ns2:SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns2:MonitoredStopVisitCancellation>
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
      | Collected                    | false                         |
      | ArrivalStatus                | arrived                       |
      | DepartureStatus              | cancelled                     |
      | Schedule[expected]#Arrival   |          2017-01-01T13:00:00Z |
      | Schedule[expected]#Departure |          2017-01-01T13:00:00Z |
      | Schedule[actual]#Arrival     |          2017-01-01T13:00:00Z |
      | Schedule[actual]#Departure   | 2017-06-19T16:04:25.983+02:00 |

  Scenario: 4448 - Manage a SM Notify after modification of a StopVisit
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
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, test-stop-monitoring-request-collector, siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | SubscriberRef     | subscriber                                         |
      | ExternalId        | externalId                                         |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 32                                      |
      | Codes                             | "internal": "NINOXE:VehicleJourney:201"         |
      | LineId                            |               6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored                         | true                                            |
      | DirectionType                     | Aller                                           |
      | OriginName                        | Le début                                        |
      | DestinationName                   | La fin. <TER>                                   |
      | Reference[DestinationRef]#Code    | "external": "ThisIsTheEnd"                      |
      | Reference[JourneyPatternRef]#Code | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | DataFrameRef                  | abcd                                                                 |
      | PassageOrder                  |                                                                    4 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop                 | false                                                                |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival       |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival    |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus                 | onTime                                                               |
      | Attribute[DestinationDisplay] | Cergy le haut & arret <RER>                                          |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And 10 seconds have passed
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | cancelled                |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope
      xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyStopMonitoring
      xmlns:sw="http://wsdl.siri.org.uk"
      xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Notification>
      <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
      <siri:Status>true</siri:Status>
      <siri:MonitoredStopVisit>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1</siri:ItemIdentifier>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
        <siri:MonitoredVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:FramedVehicleJourneyRef>
            <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
            <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
          </siri:FramedVehicleJourneyRef>
          <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
          <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:OriginName>Le début</siri:OriginName>
          <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
          <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
          <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
          <siri:Monitored>true</siri:Monitored>
          <siri:MonitoredCall>
            <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
            <siri:Order>4</siri:Order>
            <siri:StopPointName>Test</siri:StopPointName>
            <siri:VehicleAtStop>false</siri:VehicleAtStop>
            <siri:DestinationDisplay>Cergy le haut &amp; arret &lt;RER&gt;</siri:DestinationDisplay>
            <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
            <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
          </siri:MonitoredCall>
        </siri:MonitoredVehicleJourney>
      </siri:MonitoredStopVisit>
      <siri:MonitoredStopVisitCancellation>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1</siri:ItemRef>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
        <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
        <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
        <siri:VehicleJourneyRef>
          <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
          <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
        </siri:VehicleJourneyRef>
      </siri:MonitoredStopVisitCancellation>
      </siri:StopMonitoringDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyStopMonitoring>
      </S:Body>
      </S:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol        | siri                           |
      | Direction       | sent                           |
      | Status          | OK                             |
      | Type            | NotifyStopMonitoring           |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201"]  |
      | Lines           | ["NINOXE:Line:3:LOC"]          |

  Scenario: 4448 - Manage a SM Notify after modification of a StopVisit with the RewriteJourneyPatternRef setting
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
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, test-stop-monitoring-request-collector, siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url                            | http://localhost:8090 |
      | remote_credential                     | test                  |
      | local_credential                      | NINOXE:default        |
      | remote_code_space                     | internal              |
      | broadcast.rewrite_journey_pattern_ref | true                  |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | SubscriberRef     | subscriber                                         |
      | ExternalId        | externalId                                         |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 32                                      |
      | Codes                             | "internal": "NINOXE:VehicleJourney:201"         |
      | LineId                            |               6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored                         | true                                            |
      | DirectionType                     | Aller                                           |
      | OriginName                        | Le début                                        |
      | DestinationName                   | La fin.                                         |
      | Reference[DestinationRef]#Code    | "external": "ThisIsTheEnd"                      |
      | Reference[JourneyPatternRef]#Code | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And 10 seconds have passed
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | cancelled                |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope
      xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyStopMonitoring
      xmlns:sw="http://wsdl.siri.org.uk"
      xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Notification>
      <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
      <siri:Status>true</siri:Status>
      <siri:MonitoredStopVisit>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1</siri:ItemIdentifier>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
        <siri:MonitoredVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:FramedVehicleJourneyRef>
            <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
            <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
          </siri:FramedVehicleJourneyRef>
          <siri:JourneyPatternRef>RATPDev:JourneyPattern::775b650b33aa71eaa01222ccf88a68ce23b58eff:LOC</siri:JourneyPatternRef>
          <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:OriginName>Le début</siri:OriginName>
          <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
          <siri:DestinationName>La fin.</siri:DestinationName>
          <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
          <siri:Monitored>true</siri:Monitored>
          <siri:MonitoredCall>
            <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
            <siri:Order>4</siri:Order>
            <siri:StopPointName>Test</siri:StopPointName>
            <siri:VehicleAtStop>false</siri:VehicleAtStop>
            <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
            <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
          </siri:MonitoredCall>
        </siri:MonitoredVehicleJourney>
      </siri:MonitoredStopVisit>
      <siri:MonitoredStopVisitCancellation>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1</siri:ItemRef>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
        <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
        <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
        <siri:VehicleJourneyRef>
          <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
          <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
        </siri:VehicleJourneyRef>
      </siri:MonitoredStopVisitCancellation>
      </siri:StopMonitoringDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyStopMonitoring>
      </S:Body>
      </S:Envelope>
      """

  Scenario: 4448 - Manage a SM Notify after modification of a StopVisit with the no DestinationRef rewrite setting
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
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, test-stop-monitoring-request-collector, siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url                                 | http://localhost:8090 |
      | remote_credential                          | test                  |
      | local_credential                           | NINOXE:default        |
      | remote_code_space                          | internal              |
      | broadcast.no_destinationref_rewriting_from | NoRewriteOrigin       |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | SubscriberRef     | subscriber                                         |
      | ExternalId        | externalId                                         |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Origin                         | NoRewriteOrigin                         |
      | Name                           | Passage 32                              |
      | Codes                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                         |       6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored                      | true                                    |
      | DirectionType                  | Aller                                   |
      | OriginName                     | Le début                                |
      | DestinationName                | La fin.                                 |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"              |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And 10 seconds have passed
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | cancelled                |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope
      xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyStopMonitoring
      xmlns:sw="http://wsdl.siri.org.uk"
      xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Notification>
      <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
      <siri:Status>true</siri:Status>
      <siri:MonitoredStopVisit>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1</siri:ItemIdentifier>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
        <siri:MonitoredVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:FramedVehicleJourneyRef>
            <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
            <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
          </siri:FramedVehicleJourneyRef>
          <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:OriginName>Le début</siri:OriginName>
          <siri:DestinationRef>ThisIsTheEnd</siri:DestinationRef>
          <siri:DestinationName>La fin.</siri:DestinationName>
          <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
          <siri:Monitored>true</siri:Monitored>
          <siri:MonitoredCall>
            <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
            <siri:Order>4</siri:Order>
            <siri:StopPointName>Test</siri:StopPointName>
            <siri:VehicleAtStop>false</siri:VehicleAtStop>
            <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
            <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
          </siri:MonitoredCall>
        </siri:MonitoredVehicleJourney>
      </siri:MonitoredStopVisit>
      <siri:MonitoredStopVisitCancellation>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1</siri:ItemRef>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
        <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
        <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
        <siri:VehicleJourneyRef>
          <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
          <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
        </siri:VehicleJourneyRef>
      </siri:MonitoredStopVisitCancellation>
      </siri:StopMonitoringDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyStopMonitoring>
      </S:Body>
      </S:Envelope>
      """

  Scenario: 4448 - Manage a SM Notify after modification of a StopVisit with the no DataFrameRef rewrite setting
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
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, test-stop-monitoring-request-collector, siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url                               | http://localhost:8090 |
      | remote_credential                        | test                  |
      | local_credential                         | NINOXE:default        |
      | remote_code_space                        | internal              |
      | broadcast.no_dataframeref_rewriting_from | NoRewriteOrigin       |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | SubscriberRef     | subscriber                                         |
      | ExternalId        | externalId                                         |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Origin                         | NoRewriteOrigin                         |
      | Name                           | Passage 32                              |
      | Codes                          | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                         |       6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored                      | true                                    |
      | DirectionType                  | Aller                                   |
      | OriginName                     | Le début                                |
      | DestinationName                | La fin.                                 |
      | Reference[DestinationRef]#Code | "external": "ThisIsTheEnd"              |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And 10 seconds have passed
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | cancelled                |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope
      xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyStopMonitoring
      xmlns:sw="http://wsdl.siri.org.uk"
      xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Notification>
      <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
      <siri:Status>true</siri:Status>
      <siri:MonitoredStopVisit>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1</siri:ItemIdentifier>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
        <siri:MonitoredVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:FramedVehicleJourneyRef>
            <siri:DataFrameRef>abcd</siri:DataFrameRef>
            <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
          </siri:FramedVehicleJourneyRef>
          <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:OriginName>Le début</siri:OriginName>
          <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
          <siri:DestinationName>La fin.</siri:DestinationName>
          <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
          <siri:Monitored>true</siri:Monitored>
          <siri:MonitoredCall>
            <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
            <siri:Order>4</siri:Order>
            <siri:StopPointName>Test</siri:StopPointName>
            <siri:VehicleAtStop>false</siri:VehicleAtStop>
            <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
            <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
          </siri:MonitoredCall>
        </siri:MonitoredVehicleJourney>
      </siri:MonitoredStopVisit>
      <siri:MonitoredStopVisitCancellation>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1</siri:ItemRef>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
        <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
        <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
        <siri:VehicleJourneyRef>
          <siri:DataFrameRef>abcd</siri:DataFrameRef>
          <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
        </siri:VehicleJourneyRef>
      </siri:MonitoredStopVisitCancellation>
      </siri:StopMonitoringDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyStopMonitoring>
      </S:Body>
      </S:Envelope>
      """

  Scenario: Manage a DeleteSubscription Request
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
    And a Partner "test" exists with connectors [siri-check-status-client, test-stop-monitoring-request-collector, siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 32                                      |
      | Codes                             | "internal": "NINOXE:VehicleJourney:201"         |
      | LineId                            |               6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored                         | true                                            |
      | DirectionType                     | Aller                                           |
      | OriginName                        | Le début                                        |
      | DestinationName                   | La fin.                                         |
      | Reference[DestinationRef]#Code    | "external": "ThisIsTheEnd"                      |
      | Reference[JourneyPatternRef]#Code | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And 10 seconds have passed
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
      <SOAP-ENV:Header />
      <S:Body>
      <sw:DeleteSubscription xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <DeleteSubscriptionInfo>
        <siri:RequestTimestamp>2006-01-02T15:04:05.000Z07:00</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
        <siri:MessageIdentifier>MessageIdentifier</siri:MessageIdentifier>
      </DeleteSubscriptionInfo>
      <Request version="2.0:FR-IDF-2.4">
        <siri:All/>
      </Request>
      <RequestExtension/>
      </sw:DeleteSubscription>
      </S:Body>
      </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:DeleteSubscriptionResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
      <DeleteSubscriptionAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:40.000Z</siri:ResponseTimestamp>
        <siri:ResponderRef>test</siri:ResponderRef>
        <siri:RequestMessageRef>MessageIdentifier</siri:RequestMessageRef>
      </DeleteSubscriptionAnswerInfo>
      <Answer>
        <siri:ResponseTimestamp>2017-01-01T12:00:40.000Z</siri:ResponseTimestamp>
        <siri:ResponderRef>test</siri:ResponderRef>
        <siri:RequestMessageRef>MessageIdentifier</siri:RequestMessageRef>
        <siri:TerminationResponseStatus>
          <siri:ResponseTimestamp>2017-01-01T12:00:40.000Z</siri:ResponseTimestamp>
          <siri:SubscriberRef>NINOXE:default</siri:SubscriberRef>
          <siri:SubscriptionRef/>
          <siri:Status>true</siri:Status>
        </siri:TerminationResponseStatus>
      </Answer>
      <AnswerExtension/>
      </sw:DeleteSubscriptionResponse>
      </S:Body>
      </S:Envelope>
      """

  @ARA-957
  Scenario: Send DeleteSubscriptionRequests
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
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
                <ns2:SubscriptionRef>Ara:Subscription::6ba7b814-9dad-11d1-33-00c04fd430c8:LOC</ns2:SubscriptionRef>
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
    Then the SIRI server should have received 1 DeleteSubscription request
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
                <ns2:SubscriptionRef>Ara:Subscription::6ba7b814-9dad-11d1-33-00c04fd430c8:LOC</ns2:SubscriptionRef>
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
    Then the SIRI server should not have received 2 DeleteSubscription requests
    When 6 minutes have passed
    And I send this SIRI request
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
                <ns2:SubscriptionRef>Ara:Subscription::6ba7b814-9dad-11d1-33-00c04fd430c8:LOC</ns2:SubscriptionRef>
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
    Then the SIRI server should have received 2 DeleteSubscription requests

  @ARA-1252 @ARA-1234
  Scenario: Update HasCompleteStopSequence of a VehicleJourney when receiving all StopVisits with ScheduledStopVisits
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
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_code_space                  | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:25:LOC" |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
      # "SubscriptionRef":"RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC"
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:25:LOC" |
      # "SubscriptionRef":"RELAIS:Subscription::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC"
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                              |
      | Codes                    | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                   |       6ba7b814-9dad-11d1-a-00c04fd430c8 |
      | Monitored                | true                                    |
      | Attribute[DirectionName] | A Direction Name                        |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a ScheduledStopVisit exists with the following attributes:
      | Codes                   | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder            |                                                                    1 |
      | StopAreaId              |                                    6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId        |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop           | false                                                                |
      | Schedule[aimed]#Arrival |                                             2017-01-01T15:00:00.000Z |
      # "Id":"6ba7b814-9dad-11d1-8-00c04fd430c8"
    And a ScheduledStopVisit exists with the following attributes:
      | Codes                   | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-3" |
      | PassageOrder            |                                                                    2 |
      | StopAreaId              |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId        |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop           | false                                                                |
      | Schedule[aimed]#Arrival |                                             2017-01-01T15:20:00.000Z |
      # "Id":"6ba7b814-9dad-11d1-9-00c04fd430c8"
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
              <ns2:ResponseTimestamp>
              2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>
              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:DirectionRef>Left</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:42:LOC</ns3:OriginRef>
                    <ns3:DestinationRef>NINOXE:StopPoint:SP:62:LOC</ns3:DestinationRef>
                    <ns3:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</ns3:OriginAimedDepartureTime>
                    <ns3:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</ns3:DestinationAimedArrivalTime>
                    <ns3:Monitored>true</ns3:Monitored>
                    <ns3:CourseOfJourneyRef>201</ns3:CourseOfJourneyRef>
                    <ns3:VehicleRef>NINOXE:Vehicle:23:LOC</ns3:VehicleRef>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R)</ns3:StopPointName>
                      <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T13:00:00.000+02:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>delayed</ns3:ArrivalStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns2:StopMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </ns6:NotifyStopMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then the VehicleJourney "6ba7b814-9dad-11d1-7-00c04fd430c8" has the following attributes:
      | HasCompleteStopSequence | false                                 |
      | StopVisits              | ["6ba7b814-9dad-11d1-f-00c04fd430c8"] |
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
              <ns2:ResponseTimestamp>
              2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
              <ns2:ProducerRef>NINOXE:default</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>
              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:25:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:DirectionRef>Left</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:42:LOC</ns3:OriginRef>
                    <ns3:DestinationRef>NINOXE:StopPoint:SP:62:LOC</ns3:DestinationRef>
                    <ns3:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</ns3:OriginAimedDepartureTime>
                    <ns3:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</ns3:DestinationAimedArrivalTime>
                    <ns3:Monitored>true</ns3:Monitored>
                    <ns3:CourseOfJourneyRef>201</ns3:CourseOfJourneyRef>
                    <ns3:VehicleRef>NINOXE:Vehicle:23:LOC</ns3:VehicleRef>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R)</ns3:StopPointName>
                      <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T13:00:00.000+02:00</ns3:AimedArrivalTime>
                      <ns3:ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>onTime</ns3:ArrivalStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns2:StopMonitoringDelivery>
            </Notification>
            <SiriExtension />
          </ns6:NotifyStopMonitoring>
        </soap:Body>
      </soap:Envelope>
      """
    Then the VehicleJourney "6ba7b814-9dad-11d1-7-00c04fd430c8" has the following attributes:
      | HasCompleteStopSequence | true                                                                        |
      | StopVisits              | ["6ba7b814-9dad-11d1-f-00c04fd430c8", "6ba7b814-9dad-11d1-14-00c04fd430c8"] |

  @ARA-1256
  Scenario: Delete and recreate subscription when receiving subscription with same existing number
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url                         | http://localhost:8090 |
      | remote_credential                  | test                  |
      | local_credential                   | NINOXE:default        |
      | remote_code_space                  | internal              |
      | broadcast.subscriptions.persistent | true                  |
    And a StopArea exists with the following attributes:
      | Name  | Test                   |
      | Codes | "internal": "coicogn2" |
    And a minute has passed
    When I send this SIRI request
      """
      <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <SOAP-ENV:Header />
      <SOAP-ENV:Body>
      <ws:Subscribe>
      <SubscriptionRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
        <siri:MessageIdentifier>Ara:Message::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC</siri:MessageIdentifier>
        <siri:ConsumerAddress>https://ara-staging.af83.io/test/siri</siri:ConsumerAddress>
      </SubscriptionRequestInfo>
      <Request>
        <StopMonitoringSubscriptionRequest>
          <SubscriberRef>RATPDEV:Concerto</SubscriberRef>
          <SubscriptionIdentifier>1</SubscriptionIdentifier>
          <InitialTerminationTime>2009-11-10T23:00:00.000Z</InitialTerminationTime>
          <StopMonitoringRequest>
            <MessageIdentifier>28679112-9dad-11d1-2-00c04fd430c8</MessageIdentifier>
            <RequestTimestamp>2017-01-01T12:00:00.000Z</RequestTimestamp>
            <MonitoringRef>coicogn2</MonitoringRef>
            <StopVisitTypes>all</StopVisitTypes>
          </StopMonitoringRequest>
          <IncrementalUpdates>true</IncrementalUpdates>
          <ChangeBeforeUpdates>PT1M</ChangeBeforeUpdates>
        </StopMonitoringSubscriptionRequest>
      </Request>
      <RequestExtension />
      </ws:Subscribe>
      </SOAP-ENV:Body>
      </SOAP-ENV:Envelope>
      """
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Kind            | StopMonitoringBroadcast           |
      | ExternalId      |                                 1 |
    When I send this SIRI request
      """
      <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <SOAP-ENV:Header />
      <SOAP-ENV:Body>
      <ws:Subscribe>
        <SubscriptionRequestInfo>
          <siri:RequestTimestamp>2017-01-01T12:02:00.000Z</siri:RequestTimestamp>
          <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
          <siri:MessageIdentifier>Ara:Message::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC</siri:MessageIdentifier>
          <siri:ConsumerAddress>https://ara-staging.af83.io/test/siri</siri:ConsumerAddress>
        </SubscriptionRequestInfo>
        <Request>
          <StopMonitoringSubscriptionRequest>
            <SubscriberRef>RATPDEV:Concerto</SubscriberRef>
            <SubscriptionIdentifier>1</SubscriptionIdentifier>
            <InitialTerminationTime>2009-11-10T23:00:00.000Z</InitialTerminationTime>
            <StopMonitoringRequest>
              <MessageIdentifier>28679112-9dad-11d1-2-00c04fd430c8</MessageIdentifier>
              <RequestTimestamp>2017-01-01T12:00:00.000Z</RequestTimestamp>
              <MonitoringRef>coicogn2</MonitoringRef>
              <StopVisitTypes>all</StopVisitTypes>
            </StopMonitoringRequest>
            <IncrementalUpdates>true</IncrementalUpdates>
            <ChangeBeforeUpdates>PT1M</ChangeBeforeUpdates>
          </StopMonitoringSubscriptionRequest>
        </Request>
        <RequestExtension />
      </ws:Subscribe>
      </SOAP-ENV:Body>
      </SOAP-ENV:Envelope>
      """
    Then No Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Kind            | StopMonitoringBroadcast           |
      | ExternalId      |                                 1 |
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Kind            | StopMonitoringBroadcast           |
      | ExternalId      |                                 1 |
    When I send this SIRI request
      """
      <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <SOAP-ENV:Header />
      <SOAP-ENV:Body>
      <ws:Subscribe>
        <SubscriptionRequestInfo>
          <siri:RequestTimestamp>2017-01-01T12:02:00.000Z</siri:RequestTimestamp>
          <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
          <siri:MessageIdentifier>Ara:Message::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC</siri:MessageIdentifier>
          <siri:ConsumerAddress>https://ara-staging.af83.io/test/siri</siri:ConsumerAddress>
        </SubscriptionRequestInfo>
        <Request>
          <StopMonitoringSubscriptionRequest>
            <SubscriberRef>RATPDEV:Concerto</SubscriberRef>
            <SubscriptionIdentifier>2</SubscriptionIdentifier>
            <InitialTerminationTime>2009-11-10T23:00:00.000Z</InitialTerminationTime>
            <StopMonitoringRequest>
              <MessageIdentifier>28679112-9dad-11d1-2-00c04fd430c8</MessageIdentifier>
              <RequestTimestamp>2017-01-01T12:00:00.000Z</RequestTimestamp>
              <MonitoringRef>coicogn2</MonitoringRef>
              <StopVisitTypes>all</StopVisitTypes>
            </StopMonitoringRequest>
            <IncrementalUpdates>true</IncrementalUpdates>
            <ChangeBeforeUpdates>PT1M</ChangeBeforeUpdates>
          </StopMonitoringSubscriptionRequest>
        </Request>
        <RequestExtension />
      </ws:Subscribe>
      </SOAP-ENV:Body>
      </SOAP-ENV:Envelope>
      """
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Kind            | StopMonitoringBroadcast           |
      | ExternalId      |                                 1 |
    Then one Subscription exists with the following attributes:
      | SubscriptionRef | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | Kind            | StopMonitoringBroadcast           |
      | ExternalId      |                                 2 |

  @ARA-1306
  Scenario: StopMonitoring subscription collect should send StopMonitoringSubscription request to partner
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | remote_code_space | internal              |
      | local_credential  | ara                   |
    And a StopArea exists with the following attributes:
      | Name  | Arletty               |
      | Codes | "internal": "boaarle" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect           |
      | ReferenceArray[0] | StopArea, "internal": "boaarle" |
    And a minute has passed
    Then the SIRI server should have received a StopMonitoringSubscriptionRequest request with:
      | //siri:MonitoringRef | boaarle |

  @ARA-1306
  Scenario: StopMonitoring subscription collect and partner CheckStatus is unavailable should not send StopMonitoringSubscription request to partner
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | remote_code_space | internal              |
      | local_credential  | ara                   |
    And a StopArea exists with the following attributes:
      | Name  | Arletty               |
      | Codes | "internal": "boaarle" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect           |
      | ReferenceArray[0] | StopArea, "internal": "boaarle" |
    And a minute has passed
    Then the SIRI server should not have received a StopMonitoringSubscription request

  @ARA-1306
  Scenario: StopMonitoring subscription collect and partner CheckStatus is unavailable should send StopMonitoringSubscription request to partner whith setting collect.subscriptions.persistent
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                       | http://localhost:8090 |
      | remote_credential                | test                  |
      | remote_code_space                | internal              |
      | local_credential                 | ara                   |
      | collect.subscriptions.persistent | true                  |
    And a StopArea exists with the following attributes:
      | Name  | Arletty               |
      | Codes | "internal": "boaarle" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect           |
      | ReferenceArray[0] | StopArea, "internal": "boaarle" |
    And a minute has passed
    Then the SIRI server should have received a StopMonitoringSubscriptionRequest request with:
      | //siri:MonitoringRef | boaarle |

  @ARA-1306
  Scenario: StopMonitoring subscription collect and partner CheckStatus is unavailable should send StopMonitoringSubscription request to partner whith setting collect.persistent
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url         | http://localhost:8090 |
      | remote_credential  | test                  |
      | remote_code_space  | internal              |
      | local_credential   | ara                   |
      | collect.persistent | true                  |
    And a StopArea exists with the following attributes:
      | Name  | Arletty               |
      | Codes | "internal": "boaarle" |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect           |
      | ReferenceArray[0] | StopArea, "internal": "boaarle" |
    And a minute has passed
    And 10 seconds have passed
    Then the SIRI server should have received a StopMonitoringSubscriptionRequest request with:
      | //siri:MonitoringRef | boaarle |

  @ARA-1324
  Scenario: Manage a SM Notify after modification of a StopVisit if the StopVisit is the next Stop Visit of a Vehicle should broadcasr Vehicle information
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
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, test-stop-monitoring-request-collector, siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | SubscriberRef     | subscriber                                         |
      | ExternalId        | externalId                                         |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:25:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
      # 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:25:LOC" |
      | Monitored | true                                     |
      # 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
      # 6ba7b814-9dad-11d1-6-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 32                                      |
      | Codes                             | "internal": "NINOXE:VehicleJourney:201"         |
      | LineId                            |               6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | Monitored                         | true                                            |
      | DirectionType                     | Aller                                           |
      | OriginName                        | Le début                                        |
      | DestinationName                   | La fin.                                         |
      | Reference[DestinationRef]#Code    | "external": "ThisIsTheEnd"                      |
      | Reference[JourneyPatternRef]#Code | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
      # 6ba7b814-9dad-11d1-7-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
      # 6ba7b814-9dad-11d1-8-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    5 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T18:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T20:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
      # 6ba7b814-9dad-11d1-9-00c04fd430c8
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "RLA290"              |
      | LineId           | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | StopAreaId       | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
      | Bearing          |                            121.55 |
      | Latitude         |                             55.55 |
      | Longitude        |                         111.11111 |
      | Occupancy        | seatsAvailable                    |
    And 10 seconds have passed
    When the StopVisit "6ba7b814-9dad-11d1-9-00c04fd430c8" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T20:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope
      xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyStopMonitoring
      xmlns:sw="http://wsdl.siri.org.uk"
      xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-c-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Notification>
      <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:MonitoringRef>NINOXE:StopPoint:SP:25:LOC</siri:MonitoringRef>
      <siri:Status>true</siri:Status>
      <siri:MonitoredStopVisit>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1</siri:ItemIdentifier>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:25:LOC</siri:MonitoringRef>
        <siri:MonitoredVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:FramedVehicleJourneyRef>
            <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
            <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
          </siri:FramedVehicleJourneyRef>
          <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
          <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:OriginName>Le début</siri:OriginName>
          <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
          <siri:DestinationName>La fin.</siri:DestinationName>
          <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
          <siri:Monitored>true</siri:Monitored>
          <siri:Occupancy>seatsAvailable</siri:Occupancy>
          <siri:VehicleLocation>
             <siri:Longitude>111.11111</siri:Longitude>
             <siri:Latitude>55.55</siri:Latitude>
          </siri:VehicleLocation>
          <siri:Bearing>121.55</siri:Bearing>
          <siri:MonitoredCall>
            <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
            <siri:Order>5</siri:Order>
            <siri:StopPointName>Test</siri:StopPointName>
            <siri:VehicleAtStop>false</siri:VehicleAtStop>
            <siri:ExpectedArrivalTime>2017-01-01T20:01:01.000Z</siri:ExpectedArrivalTime>
            <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
          </siri:MonitoredCall>
        </siri:MonitoredVehicleJourney>
      </siri:MonitoredStopVisit>
      </siri:StopMonitoringDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyStopMonitoring>
      </S:Body>
      </S:Envelope>
      """

  @ARA-1363
  Scenario: Manage a SM Notify after modification of a StopVisit using the generator setting reference_vehicle_journey_identifier
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
    # Setting a Partner without default generators
    And a Partner "test" exists with connectors [siri-check-status-client, test-stop-monitoring-request-collector, siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url                                      | http://localhost:8090             |
      | remote_credential                               | test                              |
      | local_credential                                | NINOXE:default                    |
      | remote_code_space                               | internal                          |
      | generators.reference_vehicle_journey_identifier | ch:1:ServiceJourney:87_TAC:%{id}a |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | SubscriberRef     | subscriber                                         |
      | ExternalId        | externalId                                         |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 32                                                      |
      | Codes                             | "_default": "6ba7b814", "external": "NINOXE:VehicleJourney:201" |
      | LineId                            |                               6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored                         | true                                                            |
      | DirectionType                     | Aller                                                           |
      | OriginName                        | Le début                                                        |
      | DestinationName                   | La fin. <TER>                                                   |
      | Reference[DestinationRef]#Code    | "external": "ThisIsTheEnd"                                      |
      | Reference[JourneyPatternRef]#Code | "internal": "NINOXE:JourneyPattern:3_42_62:LOC"                 |
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | DataFrameRef                  | abcd                                                                 |
      | PassageOrder                  |                                                                    4 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop                 | false                                                                |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival       |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival    |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus                 | onTime                                                               |
      | Attribute[DestinationDisplay] | Cergy le haut & arret <RER>                                          |
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope
      xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyStopMonitoring
      xmlns:sw="http://wsdl.siri.org.uk"
      xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:40.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-8-00c04fd430c8</siri:ResponseMessageIdentifier>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Notification>
      <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:40.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
      <siri:Status>true</siri:Status>
      <siri:MonitoredStopVisit>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1</siri:ItemIdentifier>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
        <siri:MonitoredVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:FramedVehicleJourneyRef>
            <siri:DataFrameRef>2017-01-01</siri:DataFrameRef>
            <siri:DatedVehicleJourneyRef>ch:1:ServiceJourney:87_TAC:6ba7b814a</siri:DatedVehicleJourneyRef>
          </siri:FramedVehicleJourneyRef>
          <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
          <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:OriginName>Le début</siri:OriginName>
          <siri:DestinationRef>a8989abce31bae21da02c1c2cf42dd855cd86a1d</siri:DestinationRef>
          <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
          <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
          <siri:Monitored>true</siri:Monitored>
          <siri:MonitoredCall>
            <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
            <siri:Order>4</siri:Order>
            <siri:StopPointName>Test</siri:StopPointName>
            <siri:VehicleAtStop>false</siri:VehicleAtStop>
            <siri:DestinationDisplay>Cergy le haut &amp; arret &lt;RER&gt;</siri:DestinationDisplay>
            <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
            <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
          </siri:MonitoredCall>
        </siri:MonitoredVehicleJourney>
      </siri:MonitoredStopVisit>
      </siri:StopMonitoringDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyStopMonitoring>
      </S:Body>
      </S:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol        | siri                                     |
      | Direction       | sent                                     |
      | Status          | OK                                       |
      | Type            | NotifyStopMonitoring                     |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"]           |
      | VehicleJourneys | ["ch:1:ServiceJourney:87_TAC:6ba7b814a"] |
      | Lines           | ["NINOXE:Line:3:LOC"]                    |

  @ARA-1363
  Scenario: Manage a SM Notify after modification of a StopVisit using the default generator should send DatedVehicleJourneyRef according to default setting
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
    # Setting a Partner with default generators
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, test-stop-monitoring-request-collector, siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | SubscriberRef     | subscriber                                         |
      | ExternalId        | externalId                                         |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 32                                                      |
      | Codes                             | "_default": "6ba7b814", "external": "NINOXE:VehicleJourney:201" |
      | LineId                            |                               6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored                         | true                                                            |
      | DirectionType                     | Aller                                                           |
      | OriginName                        | Le début                                                        |
      | DestinationName                   | La fin. <TER>                                                   |
      | Reference[DestinationRef]#Code    | "external": "ThisIsTheEnd"                                      |
      | Reference[JourneyPatternRef]#Code | "internal": "NINOXE:JourneyPattern:3_42_62:LOC"                 |
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | DataFrameRef                  | abcd                                                                 |
      | PassageOrder                  |                                                                    4 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop                 | false                                                                |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival       |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival    |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus                 | onTime                                                               |
      | Attribute[DestinationDisplay] | Cergy le haut & arret <RER>                                          |
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope
      xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyStopMonitoring
      xmlns:sw="http://wsdl.siri.org.uk"
      xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:40.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>test</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Notification>
      <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:40.000Z</siri:ResponseTimestamp>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      <siri:SubscriberRef>subscriber</siri:SubscriberRef>
      <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
      <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
      <siri:Status>true</siri:Status>
      <siri:MonitoredStopVisit>
        <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
        <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1</siri:ItemIdentifier>
        <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
        <siri:MonitoredVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:FramedVehicleJourneyRef>
            <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
            <siri:DatedVehicleJourneyRef>RATPDev:VehicleJourney::6ba7b814:LOC</siri:DatedVehicleJourneyRef>
          </siri:FramedVehicleJourneyRef>
          <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
          <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:OriginName>Le début</siri:OriginName>
          <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
          <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
          <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
          <siri:Monitored>true</siri:Monitored>
          <siri:MonitoredCall>
            <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
            <siri:Order>4</siri:Order>
            <siri:StopPointName>Test</siri:StopPointName>
            <siri:VehicleAtStop>false</siri:VehicleAtStop>
            <siri:DestinationDisplay>Cergy le haut &amp; arret &lt;RER&gt;</siri:DestinationDisplay>
            <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
            <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
          </siri:MonitoredCall>
        </siri:MonitoredVehicleJourney>
      </siri:MonitoredStopVisit>
      </siri:StopMonitoringDelivery>
      </Notification>
      <SiriExtension />
      </sw:NotifyStopMonitoring>
      </S:Body>
      </S:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol        | siri                                     |
      | Direction       | sent                                     |
      | Status          | OK                                       |
      | Type            | NotifyStopMonitoring                     |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"]           |
      | VehicleJourneys | ["RATPDev:VehicleJourney::6ba7b814:LOC"] |
      | Lines           | ["NINOXE:Line:3:LOC"]                    |

  @ARA-1493
  Scenario: Manage referent lines in a SM Notify after modification of a StopVisit
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
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, test-stop-monitoring-request-collector, siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | SubscriberRef     | subscriber                                         |
      | ExternalId        | externalId                                         |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | Codes      | "external": "NINOXE:Line:4:LOC"   |
      | ReferentId | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Name       | Ligne 3 Metro                     |
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 32                                      |
      | Codes                             | "internal": "NINOXE:VehicleJourney:201"         |
      | LineId                            |               6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | Monitored                         | true                                            |
      | DirectionType                     | Aller                                           |
      | OriginName                        | Le début                                        |
      | DestinationName                   | La fin. <TER>                                   |
      | Reference[DestinationRef]#Code    | "external": "ThisIsTheEnd"                      |
      | Reference[JourneyPatternRef]#Code | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | DataFrameRef                  | abcd                                                                 |
      | PassageOrder                  |                                                                    4 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop                 | false                                                                |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival       |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival    |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus                 | onTime                                                               |
      | Attribute[DestinationDisplay] | Cergy le haut & arret <RER>                                          |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And 10 seconds have passed
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | cancelled                |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:NotifyStopMonitoring  xmlns:sw="http://wsdl.siri.org.uk"
                                    xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
            <siri:ProducerRef>test</siri:ProducerRef>
            <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-b-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            <siri:RequestMessageRef></siri:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Notification>
            <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
              <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
              <siri:RequestMessageRef></siri:RequestMessageRef>
              <siri:SubscriberRef>subscriber</siri:SubscriberRef>
              <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
              <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
              <siri:Status>true</siri:Status>
              <siri:MonitoredStopVisit>
                <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1</siri:ItemIdentifier>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                  <siri:DirectionRef>Aller</siri:DirectionRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                  <siri:OriginName>Le début</siri:OriginName>
                  <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                  <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
                  <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:MonitoredCall>
                    <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                    <siri:Order>4</siri:Order>
                    <siri:StopPointName>Test</siri:StopPointName>
                    <siri:VehicleAtStop>false</siri:VehicleAtStop>
                    <siri:DestinationDisplay>Cergy le haut &amp; arret &lt;RER&gt;</siri:DestinationDisplay>
                    <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
                    <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
                  </siri:MonitoredCall>
                </siri:MonitoredVehicleJourney>
              </siri:MonitoredStopVisit>
              <siri:MonitoredStopVisitCancellation>
                <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                <siri:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1</siri:ItemRef>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                <siri:VehicleJourneyRef>
                  <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                  <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                </siri:VehicleJourneyRef>
              </siri:MonitoredStopVisitCancellation>
            </siri:StopMonitoringDelivery>
          </Notification>
          <SiriExtension />
          </sw:NotifyStopMonitoring>
        </S:Body>
      </S:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol        | siri                           |
      | Direction       | sent                           |
      | Status          | OK                             |
      | Type            | NotifyStopMonitoring           |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201"]  |
      | Lines           | ["NINOXE:Line:3:LOC"]          |

  @ARA-1493
  Scenario: Manage a referent line family in a SM Notify after modification of a StopVisit
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
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, test-stop-monitoring-request-collector, siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | local_credential      | NINOXE:default        |
      | remote_code_space     | internal              |
      | sort_payload_for_test | true                  |
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | SubscriberRef     | subscriber                                         |
      | ExternalId        | externalId                                         |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | Codes      | "external": "NINOXE:Line:4:LOC"   |
      | ReferentId | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Name       | Ligne 3 Metro                     |
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 32                                      |
      | Codes                             | "internal": "NINOXE:VehicleJourney:201"         |
      | LineId                            |               6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored                         | true                                            |
      | DirectionType                     | Aller                                           |
      | OriginName                        | Le début                                        |
      | DestinationName                   | La fin. <TER>                                   |
      | Reference[DestinationRef]#Code    | "external": "ThisIsTheEnd"                      |
      | Reference[JourneyPatternRef]#Code | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Name                              | Passage 32                                      |
      | Codes                             | "internal": "NINOXE:VehicleJourney:202"         |
      | LineId                            |               6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | Monitored                         | true                                            |
      | DirectionType                     | Aller                                           |
      | OriginName                        | Le début                                        |
      | DestinationName                   | La fin. <TER>                                   |
      | Reference[DestinationRef]#Code    | "external": "ThisIsTheEnd"                      |
      | Reference[JourneyPatternRef]#Code | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | DataFrameRef                  | abcd                                                                 |
      | PassageOrder                  |                                                                    4 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop                 | false                                                                |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival       |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival    |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus                 | onTime                                                               |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" |
      | DataFrameRef                | abcd                                                                 |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-8-00c04fd430c8 |
      | VehicleAtStop               | false                                                                |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T15:00:00.000Z |
      | Schedule[expected]#Arrival  |                                             2017-01-01T15:00:00.000Z |
      | ArrivalStatus               | onTime                                                               |
    And 10 seconds have passed
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
    When the StopVisit "internal:NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1" is edited with the following attributes:
      | Schedule[expected]#Arrival | 2017-01-01T15:01:01.000Z |
      | ArrivalStatus              | Delayed                  |
    And 10 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:NotifyStopMonitoring  xmlns:sw="http://wsdl.siri.org.uk"
                                    xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
            <siri:ProducerRef>test</siri:ProducerRef>
            <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-c-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            <siri:RequestMessageRef></siri:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Notification>
            <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
              <siri:ResponseTimestamp>2017-01-01T12:00:50.000Z</siri:ResponseTimestamp>
              <siri:RequestMessageRef></siri:RequestMessageRef>
              <siri:SubscriberRef>subscriber</siri:SubscriberRef>
              <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
              <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
              <siri:Status>true</siri:Status>
              <siri:MonitoredStopVisit>
                <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1</siri:ItemIdentifier>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                  <siri:DirectionRef>Aller</siri:DirectionRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                  <siri:OriginName>Le début</siri:OriginName>
                  <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                  <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
                  <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:MonitoredCall>
                    <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                    <siri:Order>4</siri:Order>
                    <siri:StopPointName>Test</siri:StopPointName>
                    <siri:VehicleAtStop>false</siri:VehicleAtStop>
                    <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
                    <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
                  </siri:MonitoredCall>
                </siri:MonitoredVehicleJourney>
              </siri:MonitoredStopVisit>
              <siri:MonitoredStopVisit>
                <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-1</siri:ItemIdentifier>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                  <siri:DirectionRef>Aller</siri:DirectionRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                  <siri:OriginName>Le début</siri:OriginName>
                  <siri:DestinationRef>RATPDev:StopPoint:Q:a8989abce31bae21da02c1c2cf42dd855cd86a1d:LOC</siri:DestinationRef>
                  <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
                  <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:MonitoredCall>
                    <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                    <siri:Order>4</siri:Order>
                    <siri:StopPointName>Test</siri:StopPointName>
                    <siri:VehicleAtStop>false</siri:VehicleAtStop>
                    <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
                    <siri:ArrivalStatus>Delayed</siri:ArrivalStatus>
                  </siri:MonitoredCall>
                </siri:MonitoredVehicleJourney>
              </siri:MonitoredStopVisit>
            </siri:StopMonitoringDelivery>
          </Notification>
          <SiriExtension />
          </sw:NotifyStopMonitoring>
        </S:Body>
      </S:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol        | siri                           |
      | Direction       | sent                           |
      | Status          | OK                             |
      | Type            | NotifyStopMonitoring           |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201","NINOXE:VehicleJourney:202"]  |
      | Lines           | ["NINOXE:Line:3:LOC"]          |
