Feature: Support SIRI StopMonitoring by subscription

  Background:
      Given a Referential "test" is created

  Scenario: 3258 - Update a StopVisit after a StopMonitoringDelivery in a subscription
    Given a Partner "test" exists with connectors [siri-stop-monitoring-subscription-collector] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
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
                <ns2:SubscriptionRef>RELAIS:Subscription::64479:LOC</ns2:SubscriptionRef>
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
                      <ns3:StopPointRef>NINOXE:StopPoint:Q:50:LOC</ns3:StopPointRef>
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
    Then a StopVisit exists with the following attributes:
      | ObjectIDs                    | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | Schedule[expected]#Arrival   | 2017-01-01T13:01:00.000Z                                             |
      | Schedule[expected]#Departure | 2017-01-01T13:02:00.000Z                                             |

  Scenario: 3737 - Manage a MonitoredStopVisitCancellation
    Given a Partner "test" exists with connectors [siri-stop-monitoring-subscription-collector] and the following settings:
      | local_credential     | NINOXE:default |
      | remote_objectid_kind | internal       |
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
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                             |
      | DepartureStatus                 | onTime                                                               |
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
                <ns2:SubscriptionRef>Edwig:Subscription::ea7004a9-4eb1-4dbc-9e13-93afe28b41a7:LOC</ns2:SubscriptionRef>
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
      | DepartureStatus | cancelled |
      | ArrivalStatus   | cancelled |

  Scenario: Manage several StopMonitoringDeliveries in NotifyStopMonitoring
    Given a Partner "test" exists with connectors [siri-stop-monitoring-subscription-collector] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | ObjectIDs | "internal": "first" |
    And a StopArea exists with the following attributes:
      | ObjectIDs | "internal": "second" |
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
              <ns2:ProducerRef>test</ns2:ProducerRef>
              <ns2:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd59</ns2:ResponseMessageIdentifier>
              <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Notification>

              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::64479:LOC</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>first</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>first</ns3:MonitoringRef>
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
                      <ns3:StopPointRef>first</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R)</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:ExpectedArrivalTime>2017-01-01T13:01:00.000+02:00</ns3:ExpectedArrivalTime>
                      <ns3:ArrivalStatus>delayed</ns3:ArrivalStatus>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns2:StopMonitoringDelivery>

              <ns2:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns2:ResponseTimestamp>2017-05-15T13:26:12.798+02:00</ns2:ResponseTimestamp>
                <ns2:RequestMessageRef>StopMonitoring:TestDelivery:0</ns2:RequestMessageRef>
                <ns2:SubscriberRef>RELAIS</ns2:SubscriberRef>
                <ns2:SubscriptionRef>RELAIS:Subscription::64479:LOC</ns2:SubscriptionRef>
                <ns2:Status>true</ns2:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2016-09-22T07:56:53.000+02:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>second</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>second</ns3:MonitoringRef>
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
                      <ns3:StopPointRef>second</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Elf Sylvain - Métro (R)</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:ExpectedArrivalTime>2017-01-01T14:01:00.000+02:00</ns3:ExpectedArrivalTime>
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
    Then the StopVisit "internal:first" has the following attributes:
      | Schedule[expected]#Arrival   | 2017-01-01T13:01:00+02:00   |
    And the StopVisit "internal:second" has the following attributes:
      | Schedule[expected]#Arrival   | 2017-01-01T14:01:00+02:00  |
