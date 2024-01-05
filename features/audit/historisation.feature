Feature: Audit API exchanges

  Background:
    Given a Referential "test" is created

  @ARA-1241
  Scenario: Audit a event for a Stop Visit when departure status is departed
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
      | Name      | Origin                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Latitude  | 48.8566                                  |
      | Longitude | 2.3522                                   |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name  | Mago-Cime OMNI                           |
      | Codes | "internal": "NINOXE:StopPoint:SP:62:LOC" |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
      # "SubscriptionRef":"RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC"
    And a Line exists with the following attributes:
      | Codes  | "internal": "NINOXE:Line:3:LOC", "ddip": "L3:LOC" |
      | Name   | Ligne 3 Metro                                     |
      | Number | L3                                                |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"      
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                              |
      | Codes                    | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                   | 6ba7b814-9dad-11d1-6-00c04fd430c8       |
      | Monitored                | true                                    |
      | Attribute[DirectionName] | A Direction Name                        |
      | Attribute[VehicleMode]   | bus                                     |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"
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
                    <ns3:DirectionRef>aller</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:24:LOC</ns3:OriginRef>
                    <ns3:OriginName>Origin</ns3:OriginName>
                    <ns3:DestinationName>Mago-Cime OMNI</ns3:DestinationName>                    
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
                      <ns3:ArrivalStatus>arrived</ns3:ArrivalStatus>
                      <ns3:DepartureStatus>departed</ns3:DepartureStatus>
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
    Then an audit event should exist with these attributes:
      | StopVisitUUID                 | /{test-uuid}/                                                                            |
      | AimedArrivalTime              | 2017-01-01T13:00:00+02:00                                                                |
      | ExpectedArrivalTime           | 2017-01-01T13:01:00+02:00                                                                |
      | DepartureStatus               | departed                                                                                 |
      | ArrivalStatus                 | arrived                                                                                  |
      | StopAreaName                  | Origin                                                                                   |
      | StopAreaCodes                 | [{"Kind"=>"internal", "Value"=>"NINOXE:StopPoint:SP:24:LOC"}]                                 |
      | StopAreaCoordinates           | POINT(2.352200 48.856600)                                                                |
      | LineName                      | Ligne 3 Metro                                                                            |
      | LineNumber                    | L3                                                                                       |
      | TransportMode                 | bus                                                                                      |
      | LineCodes                     | [{"Kind"=>"internal", "Value"=>"NINOXE:Line:3:LOC"},{"Kind"=>"ddip", "Value"=>"L3:LOC"}]               |
      | VehicleJourneyDirectionType   | aller                                                                                    |
      | VehicleJourneyDestinationName | Mago-Cime OMNI                                                                           |
      | VehicleJourneyOriginName      | Origin                                                                                   |
      | VehicleJourneyCodes           | [{"Kind"=>"internal", "Value"=>"NINOXE:VehicleJourney:201"}]                                  |

  @ARA-1241
  Scenario: Audit a event for a Stop Visit when arrival and departure statuses are cancelled
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
      | remote_code_space               | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name      | Origin                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Latitude  | 48.8566                                  |
      | Longitude | 2.3522                                   |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name      | Mago-Cime OMNI                           |
      | Codes | "internal": "NINOXE:StopPoint:SP:62:LOC" |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
      # "SubscriptionRef":"RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC"
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
      | Number    | L3                              |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"      
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                              |
      | Codes                | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                   | 6ba7b814-9dad-11d1-6-00c04fd430c8       |
      | Monitored                | true                                    |
      | Attribute[DirectionName] | A Direction Name                        |
      | Attribute[VehicleMode]   | bus                                     |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"
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
                    <ns3:DirectionRef>aller</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:24:LOC</ns3:OriginRef>
                    <ns3:OriginName>Origin</ns3:OriginName>
                    <ns3:DestinationName>Mago-Cime OMNI</ns3:DestinationName>                    
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
                      <ns3:ActualArrivalTime>2017-01-01T13:05:00.000+02:00</ns3:ActualArrivalTime>
                      <ns3:ActualDepartureTime>2017-01-01T13:22:00.000+02:00</ns3:ActualDepartureTime>
                      <ns3:ArrivalStatus>cancelled</ns3:ArrivalStatus>
                      <ns3:DepartureStatus>cancelled</ns3:DepartureStatus>
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
    Then an audit event should exist with these attributes:
      | StopVisitUUID                 | /{test-uuid}/                                                 |
      | ActualArrivalTime             | 2017-01-01T13:05:00+02:00                                     |
      | ActualDepartureTime           | 2017-01-01T13:22:00+02:00                                     |
      | DepartureStatus               | cancelled                                                     |
      | ArrivalStatus                 | cancelled                                                     |
      | StopAreaName                  | Origin                                                        |
      | StopAreaCodes                 | [{"Kind"=>"internal", "Value"=>"NINOXE:StopPoint:SP:24:LOC"}]          |
      | StopAreaCoordinates           | POINT(2.352200 48.856600)                                     |
      | LineName                      | Ligne 3 Metro                                                 |
      | LineNumber                    | L3                                                            |
      | TransportMode                 | bus                                                           |
      | LineCodes                     | [{"Kind"=>"internal", "Value"=>"NINOXE:Line:3:LOC"}]               |
      | VehicleJourneyDirectionType   | aller                                                         |
      | VehicleJourneyDestinationName | Mago-Cime OMNI                                                |
      | VehicleJourneyOriginName      | Origin                                                        |
      | VehicleJourneyCodes           | [{"Kind"=>"internal", "Value"=>"NINOXE:VehicleJourney:201"}]       |

  @ARA-1241
  Scenario: Audit a event for a Stop Visit when departure & arrival status are set by Ara internal update mechanism
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                         | http://localhost:8090          |
      | remote_credential                  | test                           |
      | local_credential                   | NINOXE:default                 |
      | remote_code_space               | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name      | Origin                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Latitude  | 48.8566                                  |
      | Longitude | 2.3522                                   |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name      | Mago-Cime OMNI                           |
      | Codes | "internal": "NINOXE:StopPoint:SP:62:LOC" |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC", "ddip": "L3:LOC" |
      | Name      | Ligne 3 Metro                                     |
      | Number    | L3                                                |
      # "Id":"6ba7b814-9dad-11d1-5-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                              |
      | Codes                | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                   | 6ba7b814-9dad-11d1-5-00c04fd430c8       |
      | Monitored                | true                                    |
      | Attribute[DirectionName] | A Direction Name                        |
      | Attribute[VehicleMode]   | bus                                     |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-6-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T00:55:00.000+02:00                                        |
      | Schedule[aimed]#Departure       | 2017-01-01T00:59:00.000+02:00                                        |
      | DepartureStatus                 | onTime                                                               |
      | ArrivalStatus                   | onTime                                                               |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a minute has passed
    And a minute has passed
    Then an audit event should exist with these attributes:
      | StopVisitUUID       | /{test-uuid}/                                                                            |
      | AimedArrivalTime    | 2017-01-01T00:55:00+02:00                                                                |
      | AimedDepartureTime  | 2017-01-01T00:59:00+02:00                                                                |
      | ActualArrivalTime   | 2017-01-01T00:55:00+02:00                                                                |
      | ActualDepartureTime | 2017-01-01T00:59:00+02:00                                                                |
      | DepartureStatus     | departed                                                                                 |
      | ArrivalStatus       | arrived                                                                                  |
      | StopAreaName        | Mago-Cime OMNI                                                                           |
      | StopAreaCodes       | [{"Kind"=>"internal", "Value"=>"NINOXE:StopPoint:SP:62:LOC"}]                                 |
      | LineName            | Ligne 3 Metro                                                                            |
      | LineNumber          | L3                                                                                       |
      | TransportMode       | bus                                                                                      |
      | LineCodes           | [{"Kind"=>"internal", "Value"=>"NINOXE:Line:3:LOC"},{"Kind"=>"ddip", "Value"=>"L3:LOC"}]               |
      | VehicleJourneyCodes | [{"Kind"=>"internal", "Value"=>"NINOXE:VehicleJourney:201"}]                                  |

  @ARA-1241
  Scenario: Audit a event for a Stop Visit when receiving a MonitoredStopVisitCancellation
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
      | remote_code_space               | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name      | Origin                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Latitude  | 48.8566                                  |
      | Longitude | 2.3522                                   |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name      | Mago-Cime OMNI                           |
      | Codes | "internal": "NINOXE:StopPoint:SP:62:LOC" |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
      # "SubscriptionRef":"RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC"
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC", "ddip": "L3:LOC" |
      | Name      | Ligne 3 Metro                                     |
      | Number    | L3                                                |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"      
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                              |
      | Codes                | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                   | 6ba7b814-9dad-11d1-6-00c04fd430c8       |
      | Monitored                | true                                    |
      | Attribute[DirectionName] | A Direction Name                        |
      | Attribute[VehicleMode]   | bus                                     |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T00:55:00.000+02:00                                        |
      | Schedule[aimed]#Departure       | 2017-01-01T00:59:00.000+02:00                                        |
      | DepartureStatus                 | onTime                                                               |
      | ArrivalStatus                   | onTime                                                               |
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
                <ns2:MonitoredStopVisitCancellation>
                  <ns2:RecordedAtTime>2017-05-15T13:26:10.116+02:00</ns2:RecordedAtTime>
                  <ns2:ItemRef>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1</ns2:ItemRef>
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
    And a minute has passed
    Then an audit event should exist with these attributes:
      | StopVisitUUID       | /{test-uuid}/                                                                            |
      | ActualArrivalTime   | 2017-01-01T00:55:00+02:00                                                                |
      | AimedDepartureTime  | 2017-01-01T00:59:00+02:00                                                                |
      | AimedArrivalTime    | 2017-01-01T00:55:00+02:00                                                                |
      | DepartureStatus     | departed                                                                                 |
      | ArrivalStatus       | arrived                                                                                  |
      | StopAreaName        | Mago-Cime OMNI                                                                           |
      | StopAreaCodes       | [{"Kind"=>"internal", "Value"=>"NINOXE:StopPoint:SP:62:LOC"}]                                 |
      | LineName            | Ligne 3 Metro                                                                            |
      | LineNumber          | L3                                                                                       |
      | TransportMode       | bus                                                                                      |
      | LineCodes           | [{"Kind"=>"internal", "Value"=>"NINOXE:Line:3:LOC"},{"Kind"=>"ddip", "Value"=>"L3:LOC"}]               |
      | VehicleJourneyCodes | [{"Kind"=>"internal", "Value"=>"NINOXE:VehicleJourney:201"}]                                  |

  @ARA-1241
  Scenario: Audit a event for a Stop Visit when departure status is departed and the StopVisit is the NextStopVisit of a Vehicle
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
      | remote_code_space               | internal                       |
      | generators.subscription_identifier | RELAIS:Subscription::%{id}:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name      | Origin                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Latitude  | 48.8566                                  |
      | Longitude | 2.3522                                   |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name      | Mago-Cime OMNI                           |
      | Codes | "internal": "NINOXE:StopPoint:SP:62:LOC" |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                              |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
      # "SubscriptionRef":"RELAIS:Subscription::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC"
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC", "ddip": "L3:LOC" |
      | Name      | Ligne 3 Metro                                     |
      | Number    | L3                                                |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"      
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                              |
      | Codes                | "internal": "NINOXE:VehicleJourney:201" |
      | LineId                   | 6ba7b814-9dad-11d1-6-00c04fd430c8       |
      | Monitored                | true                                    |
      | Attribute[DirectionName] | A Direction Name                        |
      | Attribute[VehicleMode]   | bus                                     |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                             |
      # "Id":"6ba7b814-9dad-11d1-8-00c04fd430c8"
    And a Vehicle exists with the following attributes:
      | Codes        | "internal": "RLA290"              |
      | LineId           | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | StopAreaId       | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
      | Occupancy        | manySeatsAvailable                |
      | DriverRef        | Driver:245                        |
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
                    <ns3:DirectionRef>aller</ns3:DirectionRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>2016-09-22</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:DirectionName>Mago-Cime OMNI</ns3:DirectionName>
                    <ns3:ExternalLineRef>NINOXE:Line:3:LOC</ns3:ExternalLineRef>
                    <ns3:OriginRef>NINOXE:StopPoint:SP:24:LOC</ns3:OriginRef>
                    <ns3:OriginName>Origin</ns3:OriginName>
                    <ns3:DestinationName>Mago-Cime OMNI</ns3:DestinationName>                    
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
                      <ns3:ArrivalStatus>arrived</ns3:ArrivalStatus>
                      <ns3:DepartureStatus>departed</ns3:DepartureStatus>
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
    Then an audit event should exist with these attributes:
      | StopVisitUUID                 | /{test-uuid}/                                                                            |
      | AimedArrivalTime              | 2017-01-01T13:00:00+02:00                                                                |
      | ExpectedArrivalTime           | 2017-01-01T13:01:00+02:00                                                                |
      | DepartureStatus               | departed                                                                                 |
      | ArrivalStatus                 | arrived                                                                                  |
      | StopAreaName                  | Origin                                                                                   |
      | StopAreaCodes                 | [{"Kind"=>"internal", "Value"=>"NINOXE:StopPoint:SP:24:LOC"}]                                 |
      | StopAreaCoordinates           | POINT(2.352200 48.856600)                                                                |
      | LineName                      | Ligne 3 Metro                                                                            |
      | LineNumber                    | L3                                                                                       |
      | TransportMode                 | bus                                                                                      |
      | LineCodes                     | [{"Kind"=>"internal", "Value"=>"NINOXE:Line:3:LOC"},{"Kind"=>"ddip", "Value"=>"L3:LOC"}]              |
      | VehicleJourneyDirectionType   | aller                                                                                    |
      | VehicleJourneyDestinationName | Mago-Cime OMNI                                                                           |
      | VehicleJourneyOriginName      | Origin                                                                                   |
      | VehicleJourneyCodes           | [{"Kind"=>"internal", "Value"=>"NINOXE:VehicleJourney:201"}]                                  |
      | VehicleOccupancy              | manySeatsAvailable                                                                       |
      | VehicleDriverRef              | Driver:245                                                                               |
