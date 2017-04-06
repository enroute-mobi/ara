Feature: Support SIRI StopMonitoring

  Background:
      Given a Referential "test" is created

  Scenario: 2461 - Performs a SIRI StopMonitoring request to a Partner
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
                                   xmlns:ns4="http://www.ifopt.org.uk/acsb"
                                   xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                                   xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                                   xmlns:ns7="http://scma/siri"
                                   xmlns:ns8="http://wsdl.siri.org.uk"
                                   xmlns:ns9="http://wsdl.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <ns3:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns3:ResponseTimestamp>
        <ns3:ProducerRef>NINOXE:default</ns3:ProducerRef>
        <ns3:Address>http://appli.chouette.mobi/siri_france/siri</ns3:Address>
        <ns3:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26</ns3:ResponseMessageIdentifier>
        <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns3:ResponseTimestamp>2016-09-22T08:01:20.630+02:00</ns3:ResponseTimestamp>
          <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
          <ns3:Status>true</ns3:Status>
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
                <ns3:AimedArrivalTime>2016-09-22T07:54:00.000+02:00</ns3:AimedArrivalTime>
                <ns3:ActualArrivalTime>2016-09-22T07:54:00.000+02:00</ns3:ActualArrivalTime>
                <ns3:ArrivalStatus>arrived</ns3:ArrivalStatus>
                <ns3:ArrivalBoardingActivity>alighting</ns3:ArrivalBoardingActivity>
                <ns3:ArrivalStopAssignment>
                  <ns3:AimedQuayRef>NINOXE:StopPoint:Q:50:LOC</ns3:AimedQuayRef>
                  <ns3:ActualQuayRef>NINOXE:StopPoint:Q:50:LOC</ns3:ActualQuayRef>
                </ns3:ArrivalStopAssignment>
              </ns3:MonitoredCall>
            </ns3:MonitoredVehicleJourney>
          </ns3:MonitoredStopVisit>
        </ns3:StopMonitoringDelivery>
      </Answer>
    </ns8:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                                   | http://localhost:8090 |
      | remote_credential                            | test                  |
      | remote_objectid_kind                         | internal              |
      | collect.include_stop_areas                   | NINOXE:StopPoint:SP:24:LOC |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name      | Test 1                                   |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    When a minute has passed
    And the SIRI server has received a GetStopMonitoring request
    Then one StopVisit has the following attributes:
      | ObjectIDs    | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder | 4                                                                    |
    And one Line has the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And one VehicleJourney has the following attributes:
      | ObjectIDs                     | "internal": "NINOXE:VehicleJourney:201" |

  Scenario: Handle a SIRI StopMonitoring request
    Given a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
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
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                               |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                             |
    And I see edwig vehicle_journeys
    And I see edwig stop_visits
    And I see edwig lines
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:ns3="http://www.ifopt.org.uk/acsb"
                           xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                           xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                           xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <ns2:RequestTimestamp>2016-09-22T07:54:52.977Z</ns2:RequestTimestamp>
        <ns2:RequestorRef>test</ns2:RequestorRef>
        <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
      </ServiceRequestInfo>

      <Request version="2.0:FR-IDF-2.4">
        <ns2:RequestTimestamp>2016-09-22T07:54:52.977Z</ns2:RequestTimestamp>
        <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
        <ns2:StartTime>2016-09-22T07:54:52.977Z</ns2:StartTime>
        <ns2:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns2:MonitoringRef>
        <ns2:StopVisitTypes>all</ns2:StopVisitTypes>
      </Request>
      <RequestExtension />
    </ns7:GetStopMonitoring>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
    xmlns:ns4="http://www.ifopt.org.uk/acsb"
    xmlns:ns5="http://www.ifopt.org.uk/ifopt"
    xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns7="http://scma/siri"
    xmlns:ns8="http://wsdl.siri.org.uk"
    xmlns:ns9="http://wsdl.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
        <ns3:ProducerRef>Edwig</ns3:ProducerRef>
        <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
        <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
          <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
          <ns3:Status>true</ns3:Status>
          <ns3:MonitoredStopVisit>
            <ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
            <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
            <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
            <ns3:MonitoredVehicleJourney>
              <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
              <ns3:FramedVehicleJourneyRef>
                <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
              </ns3:FramedVehicleJourneyRef>
              <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
              <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:LOC</ns3:OperatorRef>
              <ns3:VehicleJourneyName>Passage 32</ns3:VehicleJourneyName>
              <ns3:MonitoredCall>
                <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                <ns3:Order>4</ns3:Order>
                <ns3:StopPointName>Test</ns3:StopPointName>
                <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                <ns3:ActualArrivalTime>2017-01-01T13:00:00.000Z</ns3:ActualArrivalTime>
              </ns3:MonitoredCall>
            </ns3:MonitoredVehicleJourney>
          </ns3:MonitoredStopVisit>
        </ns3:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension />
    </ns8:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """

  Scenario: Handle a SIRI StopMonitoring request on a 'empty' StopArea
    Given a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:ns3="http://www.ifopt.org.uk/acsb"
                           xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                           xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                           xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <ns2:RequestTimestamp>2016-09-22T07:54:52.977Z</ns2:RequestTimestamp>
        <ns2:RequestorRef>test</ns2:RequestorRef>
        <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
      </ServiceRequestInfo>

      <Request version="2.0:FR-IDF-2.4">
        <ns2:RequestTimestamp>2016-09-22T07:54:52.977Z</ns2:RequestTimestamp>
        <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
        <ns2:StartTime>2016-09-22T07:54:52.977Z</ns2:StartTime>
        <ns2:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns2:MonitoringRef>
        <ns2:StopVisitTypes>all</ns2:StopVisitTypes>
      </Request>
      <RequestExtension />
    </ns7:GetStopMonitoring>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='UTF-8'?>
<S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/' xmlns:SOAP-ENV='http://schemas.xmlsoap.org/soap/envelope/'>
  <S:Body>
    <ns8:GetStopMonitoringResponse xmlns:ns3='http://www.siri.org.uk/siri' xmlns:ns4='http://www.ifopt.org.uk/acsb' xmlns:ns5='http://www.ifopt.org.uk/ifopt' xmlns:ns6='http://datex2.eu/schema/2_0RC1/2_0' xmlns:ns7='http://scma/siri' xmlns:ns8='http://wsdl.siri.org.uk' xmlns:ns9='http://wsdl.siri.org.uk/siri'>
      <ServiceDeliveryInfo>
        <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
        <ns3:ProducerRef>Edwig</ns3:ProducerRef>
        <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
        <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <ns3:StopMonitoringDelivery version='2.0:FR-IDF-2.4'>
          <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
          <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
          <ns3:Status>true</ns3:Status>
        </ns3:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </ns8:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """

  Scenario: Handle a SIRI StopMonitoring request by returning all required attributes
    Given a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Destination                              |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:62:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Origin                                   |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:42:LOC" |
    And a StopArea exists with the following attributes:
      | Name      | Via                                       |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:256:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs    | "internal": "NINOXE:Line:3:LOC"           |
      | Name         | Ligne 3 Metro                             |
      | OperationRef | "internal": "NINOXE:Company:15563880:LOC" |
    And a VehicleJourney exists with the following attributes:
      | ObjectIDs                              | "internal": "NINOXE:VehicleJourney:201"         |
      | Name                                   | Magicien Noir - Cimetière (OMNI)                |
      | LineId                                 | 6ba7b814-9dad-11d1-6-00c04fd430c8               |
      | Attribute[Bearing]                     | N                                               |
      | Attribute[Delay]                       | 30                                              |
      | Attribute[DestinationName]             | Cimetière des Sauvages                          |
      | Attribute[DirectionName]               | Mago-Cime OMNI                                  |
      | Attribute[DirectionRef]                | Aller                                           |
      | Attribute[FirstOrLastJourney]          | first                                           |
      | Attribute[HeadwayService]              | false                                           |
      | Attribute[InCongestion]                | false                                           |
      | Attribute[InPanic]                     | false                                           |
      | Attribute[JourneyNote]                 | Note de test                                    |
      | Attribute[JourneyPatternName]          | TEST                                            |
      | Attribute[Monitored]                   | true                                            |
      | Attribute[MonitoringError]             | false                                           |
      | Attribute[Occupancy]                   | seatsAvailable                                  |
      | Attribute[OriginAimedDepartureTime]    | 2016-09-22T07:54:52.977Z                        |
      | Attribute[DestinationAimedArrivalTime] | 2016-09-22T09:54:52.977Z                        |
      | Attribute[OriginName]                  | Magicien Noir                                   |
      | Attribute[ProductCategoryRef]          | 0                                               |
      | Attribute[ServiceFeatureRef]           | bus scolaire                                    |
      | Attribute[TrainNumberRef]              | 12345                                           |
      | Attribute[VehicleFeatureRef]           | longTrain                                       |
      | Attribute[VehicleMode]                 | bus                                             |
      | Attribute[ViaPlaceName]                | Saint Bénédicte                                 |
      | Reference[DestinationRef]#Id           | 6ba7b814-9dad-11d1-3-00c04fd430c8               |
      | Reference[JourneyPatternRef]#ObjectID  | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
      | Reference[OriginRef]#Id                | 6ba7b814-9dad-11d1-4-00c04fd430c8               |
      | Reference[RouteRef]#ObjectID           | "internal": "NINOXE:Route:66:LOC"               |
      | Reference[PlaceRef]#Id                 | 6ba7b814-9dad-11d1-5-00c04fd430c8               |
    And a StopVisit exists with the following attributes:
      | ArrivalStatus                        | onTime                                                               |
      | DepartureStatus                      | onTime                                                               |
      | ObjectIDs                            | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                         | 4                                                                    |
      | RecordedAt                           | 2017-01-01T11:00:00.000Z                                             |
      | Schedule[actual]#Arrival             | 2017-01-01T13:00:00.000Z                                             |
      | Schedule[actual]#Departure           | 2017-01-01T13:02:00.000Z                                             |
      | Schedule[aimed]#Arrival              | 2017-01-01T13:00:00.000Z                                             |
      | Schedule[aimed]#Departure            | 2017-01-01T13:02:00.000Z                                             |
      | Schedule[expected]#Arrival           | 2017-01-01T13:00:00.000Z                                             |
      | Schedule[expected]#Departure         | 2017-01-01T13:02:00.000Z                                             |
      | StopAreaId                           | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                     | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
      | VehicleAtStop                        | true                                                                 |
      | Attribute[AimedHeadwayInterval]      | 5                                                                    |
      | Attribute[ActualQuayName]            | Quay Name                                                            |
      | Attribute[ArrivalPlatformName]       | Platform Name                                                        |
      | Attribute[ArrivalProximyTest]        | A l'approche                                                         |
      | Attribute[DepartureBoardingActivity] | boarding                                                             |
      | Attribute[DeparturePlatformName]     | Departure Platform Name                                              |
      | Attribute[DestinationDisplay]        | Balard Terminus                                                      |
      | Attribute[DistanceFromStop]          | 800                                                                  |
      | Attribute[ExpectedHeadwayInterval]   | 5                                                                    |
      | Attribute[NumberOfStopsAway]         | 1                                                                    |
      | Attribute[PlatformTraversal]         | false                                                                |
      | Reference[OperatorRef]#ObjectID      | "internal":"NINOXE:Company:15563880:LOC"                             |
    When I send a SIRI GetStopMonitoring request with
      | RequestorRef  | test                       |
      | MonitoringRef | NINOXE:StopPoint:SP:24:LOC |
    Then I should receive a SIRI GetStopMonitoringResponse with
      | //siri:MonitoredStopVisit[1]/siri:RecordedAtTime                                                                            | 2017-01-01T11:00:00.000Z                                             | StopVisit#RecordedAt                                  |
      | //siri:MonitoredStopVisit[1]/siri:ItemIdentifier                                                                            | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3               | StopVisit#ObjectID                                    |
      | //siri:MonitoredStopVisit[1]/siri:MonitoringRef                                                                             | NINOXE:StopPoint:SP:24:LOC                                           | StopArea#ObjectID                                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:LineRef                                                      | NINOXE:Line:3:LOC                                                    | Line#ObjectID                                         |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DirectionRef                                                 | Aller                                                                | VehicleJourney#Attribute[DirectionRef]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DataFrameRef                    | RATPDev:DataFrame::2017-01-01:LOC                                    | Model#Date                                            |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef          | NINOXE:VehicleJourney:201                                            | VehicleJourney#ObjectID                               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyPatternRef                                            | RATPDev:JourneyPattern::775b650b33aa71eaa01222ccf88a68ce23b58eff:LOC | VehicleJourney#Reference[JourneyPatternRef]#ObjectID  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyPatternName                                           | TEST                                                                 | VehicleJourney#Attribute[JourneyPatternName]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleMode                                                  | bus                                                                  | VehicleJourney#Attribute[VehicleMode]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:PublishedLineName                                            | Ligne 3 Metro                                                        | Line#Name                                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:RouteRef                                                     | RATPDev:Route::720c054714b4464d42970bda37a7edc5af8082cb:LOC          | VehicleJourney#Reference[RouteRef]#ObjectID           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DirectionName                                                | Mago-Cime OMNI                                                       | VehicleJourney#Attribute[DirectionName]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OperatorRef                                                  | RATPDev:Operator::dbe9523913efc7af28fe2f166da05a9013c8a647:LOC       | StopVisit#Reference[OperatorRef]#ObjectID             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:ProductCategoryRef                                           | 0                                                                    | VehicleJourney#Attribute[ProductCategoryRef]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:ServiceFeatureRef                                            | bus scolaire                                                         | VehicleJourney#Attribute[ServiceFeatureRef]           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleFeatureRef                                            | longTrain                                                            | VehicleJourney#Attribute[VehicleFeatureRef]           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginRef                                                    | NINOXE:StopPoint:SP:42:LOC                                           | VehicleJourney#Reference[OriginRef]#ObjectID          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginName                                                   | Magicien Noir                                                        | VehicleJourney#Attribute[OriginName]                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Via/siri:PlaceName                                           | Saint Bénédicte                                                      | VehicleJourney#Attribute[ViaPlaceName]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Via/siri:PlaceRef                                            | NINOXE:StopPoint:SP:256:LOC                                          | VehicleJourney#Reference[PlaceRef]#ObjectID           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationRef                                               | NINOXE:StopPoint:SP:62:LOC                                           | VehicleJourney#Reference[DestinationRef]#ObjectID     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationName                                              | Cimetière des Sauvages                                               | VehicleJourney#Attribute[DestinationName]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                           | Magicien Noir - Cimetière (OMNI)                                     | VehicleJourney#Name                                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyNote                                                  | Note de test                                                         | VehicleJourney#Attribute[JourneyNote]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:HeadwayService                                               | false                                                                | VehicleJourney#Attribute[HeadwayService]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginAimedDepartureTime                                     | 2016-09-22T07:54:52.977Z                                             | VehicleJourney#Attribute[OriginAimedDepartureTime]    |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationAimedArrivalTime                                  | 2016-09-22T09:54:52.977Z                                             | VehicleJourney#Attribute[DestinationAimedArrivalTime] |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FirstOrLastJourney                                           | first                                                                | VehicleJourney#Attribute[FirstOrLastJourney]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Monitored                                                    | true                                                                 | VehicleJourney#Attribute[Monitored]                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoringError                                              | false                                                                | VehicleJourney#Attribute[MonitoringError]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Occupancy                                                    | seatsAvailable                                                       | VehicleJourney#Attribute[Occupancy]                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Delay                                                        | 30                                                                   | VehicleJourney#Attribute[Delay]                       |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Bearing                                                      | N                                                                    | VehicleJourney#Attribute[Bearing]                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:InPanic                                                      | false                                                                | VehicleJourney#Attribute[InPanic]                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:InCongestion                                                 | false                                                                | VehicleJourney#Attribute[InCongestion]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:TrainNumber/siri:TrainNumberRef                              | 12345                                                                | VehicleJourney#Attribute[TrainNumberRef]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:SituationRef                                                 | 1234556                                                              | TODO                                                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:StopPointRef                              | NINOXE:StopPoint:SP:24:LOC                                           | StopArea#ObjectID                                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:Order                                     | 4                                                                    | StopVisit#PassageOrder                                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:VehicleAtStop                             | true                                                                 | StopVisit#VehicleAtStop                               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:PlatformTraversal                         | false                                                                | StopVisit#Attribute[PlatformTraversal]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DestinationDisplay                        | Balard Terminus                                                      | StopVisit#Attribute[DestinationDisplay]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedArrivalTime                          | 2017-01-01T13:00:00.000Z                                             | StopVisit#Schedule[aimed]#Arrival                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ActualArrivalTime                         | 2017-01-01T13:00:00.000Z                                             | StopVisit#Schedule[actual]#Arrival                    |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedArrivalTime                       | 2017-01-01T13:00:00.000Z                                             | StopVisit#Schedule[expected]#Arrival                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalStatus                             | onTime                                                               | StopVisit#ArrivalStatus                               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalProximyTest                        | A l'approche                                                         | StopVisit#Attribute[ArrivalProximyTest]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalPlatformName                       | Platform Name                                                        | StopVisit#Attribute[ArrivalPlatformName]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalStopAssignment/siri:ActualQuayName | Quay Name                                                            | StopVisit#Attribute[ActualQuayName]                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedDepartureTime                        | 2017-01-01T13:02:00.000Z                                             | StopVisit#Schedule[aimed]#Departure                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ActualDepartureTime                       | 2017-01-01T13:02:00.000Z                                             | StopVisit#Schedule[actual]#Departure                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedDepartureTime                     | 2017-01-01T13:02:00.000Z                                             | StopVisit#Schedule[expected]#Departure                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DepartureStatus                           | onTime                                                               | StopVisit#DepartureStatus                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DeparturePlatformName                     | Departure Platform Name                                              | StopVisit#Attribute[DeparturePlatformName]            |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DepartureBoardingActivity                 | boarding                                                             | StopVisit#Attribute[DepartureBoardingActivity]        |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedHeadwayInterval                      | 5                                                                    | StopVisit#Attribute[AimedHeadwayInterval]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedHeadwayInterval                   | 5                                                                    | StopVisit#Attribute[ExpectedHeadwayInterval]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DistanceFromStop                          | 800                                                                  | StopVisit#Attribute[DistanceFromStop]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:NumberOfStopsAway                         | 1                                                                    | StopVisit#Attribute[NumberOfStopsAway]                |

  Scenario: Handle a SIRI StopMonitoring response after SM Request to a SIRI server
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
        """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
  <soap:Body>
    <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>SQYBUS</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>NAVINEO:SM:RQ:107</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:StopMonitoringDelivery version="1.3">
          <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
          <ns5:Status>true</ns5:Status>
          <ns5:MonitoredStopVisit>
            <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
            <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
              <ns5:PublishedLineName>415</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
              <ns5:DestinationRef>boabonn</ns5:DestinationRef>
              <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
              <ns5:Monitored>true</ns5:Monitored>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                <ns5:Order>44</ns5:Order>
                <ns5:StopPointName>Arletty</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2017-01-01T13:43:05.000+01:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2017-01-01T13:43:05.000+01:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2017-01-01T13:43:05.000+01:00</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2017-01-01T13:43:05.000+01:00</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney>
          </ns5:MonitoredStopVisit>
        </ns5:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
    </ns1:GetStopMonitoringResponse>
  </soap:Body>
</soap:Envelope>
        """
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ratpdev               |
      | remote_objectid_kind | internal              |
    And a Partner "stif" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | STIF                                           |
      | remote_objectid_kind | external                                       |
      | remote_credential    | RATPDev                                        |
      | local_url            | https://api.concerto.ratpdev.com/concerto/siri |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Ligne 415                                                         |
      | ObjectIDs | "internal": "CdF:Line::415:LOC", "external": "STIF:Line::C00001:" |
    And a StopArea exists with the following attributes:
      | Name      | Arletty                                                                |
      | ObjectIDs | "internal": "boaarle", "external": "STIF:StopPoint:Q:eeft52df543d:" |
    And a StopArea exists with the following attributes:
      | Name            | Test 2                                                                  |
      | ObjectIDs       | "internal": "boabonn", "external": "STIF:StopPoint:Q:875fdetgyh765:" |
      | CollectedAlways | false                                                                   |
    And a minute has passed
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:ns3="http://www.ifopt.org.uk/acsb"
                           xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                           xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                           xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
        <ns2:RequestorRef>STIF</ns2:RequestorRef>
        <ns2:MessageIdentifier>STIF:Message::2345Fsdfrg35df:LOC</ns2:MessageIdentifier>
      </ServiceRequestInfo>

      <Request version="2.0:FR-IDF-2.4">
        <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
        <ns2:MessageIdentifier>STIF:Message::2345Fsdfrg35df:LOC</ns2:MessageIdentifier>
        <ns2:MonitoringRef>STIF:StopPoint:Q:eeft52df543d:</ns2:MonitoringRef>
      </Request>
      <RequestExtension />
    </ns7:GetStopMonitoring>
  </S:Body>
</S:Envelope>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
    xmlns:ns4="http://www.ifopt.org.uk/acsb"
    xmlns:ns5="http://www.ifopt.org.uk/ifopt"
    xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
    xmlns:ns7="http://scma/siri"
    xmlns:ns8="http://wsdl.siri.org.uk"
    xmlns:ns9="http://wsdl.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <ns3:ResponseTimestamp>2017-01-01T12:02:00.000Z</ns3:ResponseTimestamp>
        <ns3:ProducerRef>RATPDev</ns3:ProducerRef>
        <ns3:Address>https://api.concerto.ratpdev.com/concerto/siri</ns3:Address>
        <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-e-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
        <ns3:RequestMessageRef>STIF:Message::2345Fsdfrg35df:LOC</ns3:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns3:ResponseTimestamp>2017-01-01T12:02:00.000Z</ns3:ResponseTimestamp>
          <ns3:RequestMessageRef>STIF:Message::2345Fsdfrg35df:LOC</ns3:RequestMessageRef>
          <ns3:Status>true</ns3:Status>
          <ns3:MonitoredStopVisit>
            <ns3:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns3:RecordedAtTime>
            <ns3:ItemIdentifier>RATPDev:Item::4d25c8186b19a5b1993e4a401aebec7fc5e8bd15:LOC</ns3:ItemIdentifier>
            <ns3:MonitoringRef>STIF:StopPoint:Q:eeft52df543d:</ns3:MonitoringRef>
            <ns3:MonitoredVehicleJourney>
              <ns3:LineRef>STIF:Line::C00001:</ns3:LineRef>
              <ns3:FramedVehicleJourneyRef>
                <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                <ns3:DatedVehicleJourneyRef>RATPDev:VehicleJourney::5d5ddf96f5db438e2f4e24af3c074e2d0733cc4e:LOC</ns3:DatedVehicleJourneyRef>
              </ns3:FramedVehicleJourneyRef>
              <ns3:JourneyPatternRef>RATPDev:JourneyPattern::983a5c43233dc44a0ed956117ee55d257fea06eb:LOC</ns3:JourneyPatternRef>
              <ns3:PublishedLineName>Ligne 415</ns3:PublishedLineName>
              <ns3:DirectionName>Aller</ns3:DirectionName>
              <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:LOC</ns3:OperatorRef>
              <ns3:DestinationRef>STIF:StopPoint:Q:875fdetgyh765:</ns3:DestinationRef>
              <ns3:DestinationName>Méliès - Croix Bonnet</ns3:DestinationName>
              <ns3:Monitored>true</ns3:Monitored>
              <ns3:MonitoredCall>
                <ns3:StopPointRef>STIF:StopPoint:Q:eeft52df543d:</ns3:StopPointRef>
                <ns3:Order>44</ns3:Order>
                <ns3:StopPointName>Arletty</ns3:StopPointName>
                <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                <ns3:DestinationDisplay>Méliès - Croix Bonnet</ns3:DestinationDisplay>
                <ns3:AimedArrivalTime>2017-01-01T13:43:05.000+01:00</ns3:AimedArrivalTime>
                <ns3:ExpectedArrivalTime>2017-01-01T13:43:05.000+01:00</ns3:ExpectedArrivalTime>
                <ns3:ArrivalStatus>onTime</ns3:ArrivalStatus>
                <ns3:AimedDepartureTime>2017-01-01T13:43:05.000+01:00</ns3:AimedDepartureTime>
                <ns3:ExpectedDepartureTime>2017-01-01T13:43:05.000+01:00</ns3:ExpectedDepartureTime>
                <ns3:DepartureStatus>onTime</ns3:DepartureStatus>
              </ns3:MonitoredCall>
            </ns3:MonitoredVehicleJourney>
          </ns3:MonitoredStopVisit>
        </ns3:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension />
    </ns8:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """

  Scenario: Handles invalid GetStopMonitoring response
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
        <html><title>Error</title></body>Error 500</body></html>
      """
    And a Partner "invalid" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | remote_objectid_kind | internal              |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | ObjectIDs | "internal": "dummy" |
    When a minute has passed
    And the SIRI server has received a GetStopMonitoring request
    Then a StopArea exists with the following attributes:
      | ObjectIDs   | "internal": "dummy" |
      | CollectedAt | -                   |

  Scenario: Handle a SIRI StopMonitoring response after SM cancellation from a SIRI server
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
 """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
  <soap:Body>
    <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>SQYBUS</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>NAVINEO:SM:RQ:107</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:StopMonitoringDelivery version="1.3">
          <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
          <ns5:Status>true</ns5:Status>
          <ns5:MonitoredStopVisit>
            <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
            <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
              <ns5:PublishedLineName>415</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
              <ns5:DestinationRef>boabonn</ns5:DestinationRef>
              <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
              <ns5:Monitored>true</ns5:Monitored>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                <ns5:Order>44</ns5:Order>
                <ns5:StopPointName>Arletty</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2017-01-01T12:43:05.000+00:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2017-01-01T12:43:05.000+00:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2017-01-01T12:43:05.000</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2017-01-01T12:43:05.000</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney>
          </ns5:MonitoredStopVisit>
        </ns5:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
    </ns1:GetStopMonitoringResponse>
  </soap:Body>
</soap:Envelope>
        """
     And a Partner "ineo" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ratpdev               |
      | remote_objectid_kind | internal              |
    And a Partner "stif" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | STIF     |
      | remote_objectid_kind | external |
      | remote_credential    | RATPDev  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Ligne 415                                                         |
      | ObjectIDs | "internal": "CdF:Line::415:LOC", "external": "STIF:Line::C00001:" |
    And a StopArea exists with the following attributes:
      | Name      | Arletty                                                                |
      | ObjectIDs | "internal": "boaarle", "external": "RATPDev:StopPoint:Q:eeft52df543d:" |
    And a StopArea exists with the following attributes:
      | Name            | Test 2                                                                  |
      | ObjectIDs       | "internal": "boabonn", "external": "RATPDev:StopPoint:Q:875fdetgyh765:" |
      | CollectedAlways | false                                                                   |
    And a minute has passed
    And the SIRI server waits GetStopMonitoring request to respond with
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>SQYBUS</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>NAVINEO:SM:RQ:107</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:ArrivalStatus>cancelled</ns5:ArrivalStatus>
                      <ns5:DepartureStatus>cancelled</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And 2 minutes have passed
    When the SIRI server has received 2 GetStopMonitoring requests
    And I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:ns3="http://www.ifopt.org.uk/acsb"
                                 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                                 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                                 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <ServiceRequestInfo>
              <ns2:RequestTimestamp>2017-01-01T12:03:00.000Z</ns2:RequestTimestamp>
              <ns2:RequestorRef>STIF</ns2:RequestorRef>
              <ns2:MessageIdentifier>STIF:Message::2345Fsdfrg35df:LOC</ns2:MessageIdentifier>
            </ServiceRequestInfo>

            <Request version="2.0:FR-IDF-2.4">
              <ns2:RequestTimestamp>2017-01-01T12:03:00.000Z</ns2:RequestTimestamp>
              <ns2:MessageIdentifier>STIF:Message::2345Fsdfrg35df:LOC</ns2:MessageIdentifier>
              <ns2:MonitoringRef>RATPDev:StopPoint:Q:eeft52df543d:</ns2:MonitoringRef>
            </Request>
            <RequestExtension />
          </ns7:GetStopMonitoring>
        </S:Body>
      </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
      xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
          xmlns:ns4="http://www.ifopt.org.uk/acsb"
          xmlns:ns5="http://www.ifopt.org.uk/ifopt"
          xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns7="http://scma/siri"
          xmlns:ns8="http://wsdl.siri.org.uk"
          xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>2017-01-01T12:04:00.000Z</ns3:ResponseTimestamp>
              <ns3:ProducerRef>RATPDev</ns3:ProducerRef>
              <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-13-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>STIF:Message::2345Fsdfrg35df:LOC</ns3:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>2017-01-01T12:04:00.000Z</ns3:ResponseTimestamp>
                <ns3:RequestMessageRef>STIF:Message::2345Fsdfrg35df:LOC</ns3:RequestMessageRef>
                <ns3:Status>true</ns3:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>RATPDev:Item::4d25c8186b19a5b1993e4a401aebec7fc5e8bd15:LOC</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>RATPDev:StopPoint:Q:eeft52df543d:</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>STIF:Line::C00001:</ns3:LineRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>RATPDev:VehicleJourney::5d5ddf96f5db438e2f4e24af3c074e2d0733cc4e:LOC</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:JourneyPatternRef>RATPDev:JourneyPattern::983a5c43233dc44a0ed956117ee55d257fea06eb:LOC</ns3:JourneyPatternRef>
                    <ns3:PublishedLineName>Ligne 415</ns3:PublishedLineName>
                    <ns3:DirectionName>Aller</ns3:DirectionName>
                    <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:LOC</ns3:OperatorRef>
                    <ns3:DestinationRef>RATPDev:StopPoint:Q:875fdetgyh765:</ns3:DestinationRef>
                    <ns3:DestinationName>Méliès - Croix Bonnet</ns3:DestinationName>
                    <ns3:Monitored>true</ns3:Monitored>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>RATPDev:StopPoint:Q:eeft52df543d:</ns3:StopPointRef>
                      <ns3:Order>44</ns3:Order>
                      <ns3:StopPointName>Arletty</ns3:StopPointName>
                      <ns3:VehicleAtStop>false</ns3:VehicleAtStop>
                      <ns3:DestinationDisplay>Méliès - Croix Bonnet</ns3:DestinationDisplay>
                      <ns3:ArrivalStatus>cancelled</ns3:ArrivalStatus>
                      <ns3:DepartureStatus>cancelled</ns3:DepartureStatus>
                      </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns3:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension />
          </ns8:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: Manage a passed StopVisit
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
    # include a MonitoredStopVisit/ItemIdentifier A at 13:00
    # include a MonitoredStopVisit/ItemIdentifier B arrival 12:02:30 / departure 12:03
    # include a MonitoredStopVisit/ItemIdentifier C at 15:00
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>SQYBUS</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>NAVINEO:SM:RQ:107</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>StopVisit:A</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T13:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T13:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T13:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T13:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>StopVisit:B</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T12:02:30.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T12:02:30.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T12:03:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T12:03:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>StopVisit:C</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | Test                  |
      | remote_objectid_kind | internal              |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name      | Arletty               |
      | ObjectIDs | "internal": "boaarle" |
    And a minute has passed
    And the SIRI server waits GetStopMonitoring request to respond with
      # include a MonitoredStopVisit/ItemIdentifier A at 14:00
      # no MonitoredStopVisit/ItemIdentifier B
      # include a MonitoredStopVisit/ItemIdentifier C at 15:00
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>SQYBUS</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>NAVINEO:SM:RQ:107</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>StopVisit:A</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T13:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T13:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T13:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T13:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>StopVisit:C</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And 90 seconds have passed
    When the SIRI server has received 2 GetStopMonitoring requests
    And the StopVisit "6ba7b814-9dad-11d1-d-00c04fd430c8" has the following attributes:
      # "internal": "A"
      | DepartureStatus   | onTime          |
      | ArrivalStatus     | onTime          |
    And the StopVisit "6ba7b814-9dad-11d1-e-00c04fd430c8" has the following attributes:
      # "internal": "B"
      | Collected   | false                |
      | CollectedAt | 2017-01-01T12:02:00Z |
    And the StopVisit "6ba7b814-9dad-11d1-f-00c04fd430c8" has the following attributes:
      # "internal": "C"
      | DepartureStatus   | onTime          |
      | ArrivalStatus     | onTime          |
    And 10 seconds have passed
    And the StopVisit "6ba7b814-9dad-11d1-e-00c04fd430c8" has the following attributes:
      # "internal": "B"
      | Collected       | false                |
      | CollectedAt     | 2017-01-01T12:02:00Z |
      | DepartureStatus | cancelled            |
      | ArrivalStatus   | cancelled            |

  Scenario: 2466 - Don't perform StopMonitoring request for an unmonitored StopArea
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
                                   xmlns:ns4="http://www.ifopt.org.uk/acsb"
                                   xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                                   xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                                   xmlns:ns7="http://scma/siri"
                                   xmlns:ns8="http://wsdl.siri.org.uk"
                                   xmlns:ns9="http://wsdl.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <ns3:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns3:ResponseTimestamp>
        <ns3:ProducerRef>NINOXE:default</ns3:ProducerRef>
        <ns3:Address>http://appli.chouette.mobi/siri_france/siri</ns3:Address>
        <ns3:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26</ns3:ResponseMessageIdentifier>
        <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns3:ResponseTimestamp>2016-09-22T08:01:20.630+02:00</ns3:ResponseTimestamp>
          <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
          <ns3:Status>true</ns3:Status>
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
                <ns3:AimedArrivalTime>2016-09-22T07:54:00.000+02:00</ns3:AimedArrivalTime>
                <ns3:ActualArrivalTime>2016-09-22T07:54:00.000+02:00</ns3:ActualArrivalTime>
                <ns3:ArrivalStatus>arrived</ns3:ArrivalStatus>
                <ns3:ArrivalBoardingActivity>alighting</ns3:ArrivalBoardingActivity>
                <ns3:ArrivalStopAssignment>
                  <ns3:AimedQuayRef>NINOXE:StopPoint:Q:50:LOC</ns3:AimedQuayRef>
                  <ns3:ActualQuayRef>NINOXE:StopPoint:Q:50:LOC</ns3:ActualQuayRef>
                </ns3:ArrivalStopAssignment>
              </ns3:MonitoredCall>
            </ns3:MonitoredVehicleJourney>
          </ns3:MonitoredStopVisit>
        </ns3:StopMonitoringDelivery>
      </Answer>
    </ns8:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "source" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | source                |
      | remote_objectid_kind | internal              |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name            | arrêt 1               |
      | ObjectIDs       | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | CollectedAlways | false |
    When a minute has passed
    Then the SIRI server should not have received a GetStopMonitoring request

  Scenario: Perform StopMonitoring request for an unmonitored StopArea
     Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
                                   xmlns:ns4="http://www.ifopt.org.uk/acsb"
                                   xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                                   xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                                   xmlns:ns7="http://scma/siri"
                                   xmlns:ns8="http://wsdl.siri.org.uk"
                                   xmlns:ns9="http://wsdl.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <ns3:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns3:ResponseTimestamp>
        <ns3:ProducerRef>NINOXE:default</ns3:ProducerRef>
        <ns3:Address>http://appli.chouette.mobi/siri_france/siri</ns3:Address>
        <ns3:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26</ns3:ResponseMessageIdentifier>
        <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns3:ResponseTimestamp>2016-09-22T08:01:20.630+02:00</ns3:ResponseTimestamp>
          <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
          <ns3:Status>true</ns3:Status>
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
                <ns3:AimedArrivalTime>2016-09-22T07:54:00.000+02:00</ns3:AimedArrivalTime>
                <ns3:ActualArrivalTime>2016-09-22T07:54:00.000+02:00</ns3:ActualArrivalTime>
                <ns3:ArrivalStatus>arrived</ns3:ArrivalStatus>
                <ns3:ArrivalBoardingActivity>alighting</ns3:ArrivalBoardingActivity>
                <ns3:ArrivalStopAssignment>
                  <ns3:AimedQuayRef>NINOXE:StopPoint:Q:50:LOC</ns3:AimedQuayRef>
                  <ns3:ActualQuayRef>NINOXE:StopPoint:Q:50:LOC</ns3:ActualQuayRef>
                </ns3:ArrivalStopAssignment>
              </ns3:MonitoredCall>
            </ns3:MonitoredVehicleJourney>
          </ns3:MonitoredStopVisit>
        </ns3:StopMonitoringDelivery>
      </Answer>
    </ns8:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "source" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | source                |
      | remote_objectid_kind | internal              |
    And a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name            | arrêt 1               |
      | ObjectIDs       | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | MonitoredAlways | false |
      | CollectedAlways | false |
    When I send a SIRI GetStopMonitoring request with
      | RequestTimestamp  | 2017-01-01T07:54:00.977Z   |
      | RequestorRef      | test                       |
      | MessageIdentifier | StopMonitoring:Test:0      |
      | StartTime         | 2017-01-01T07:54:00.977Z   |
      | MonitoringRef     | NINOXE:StopPoint:SP:24:LOC |
      | StopVisitTypes    | all                        |
    And a minute has passed
    Then the SIRI server should have received a GetStopMonitoring request with:
      | //siri:MonitoringRef | NINOXE:StopPoint:SP:24:LOC |
    # And the StopArea "arrêt 1" should have the following attributes:
    #   | CollectedUntil | ~ 07h54 |

  Scenario: 2939 - Partner Setting collect.include_stop_areas is used to select the best Partner
    Given a SIRI server "first" waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>first</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>first:ResponseMessage::6ba:LOC</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>first:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And a SIRI server "second" waits GetStopMonitoring request on "http://localhost:8091" to respond with
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>second</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>second:ResponseMessage::tf7:LOC</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:1</ns5:RequestMessageRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:33193250</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaboon</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>second:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaboon</ns5:StopPointRef>
                      <ns5:Order>45</ns5:Order>
                      <ns5:StopPointName>Charles</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:04:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:04:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:05.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:05:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And a Partner "first" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                 | http://localhost:8090 |
      | collect.include_stop_areas | first                 |
      | remote_objectid_kind       | external              |
      | remote_credential          | dummy                 |
    And a Partner "second" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                 | http://localhost:8091 |
      | collect.include_stop_areas | second                |
      | remote_objectid_kind       | external              |
      | remote_credential          | dummy                 |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | ObjectIDs       | "external": "first" |
    And a StopArea exists with the following attributes:
      | ObjectIDs       | "external": "second" |
    When a minute has passed
    Then the "first" SIRI server should have received a GetStopMonitoring request with:
      | //siri:MonitoringRef | first |
    And the "second" SIRI server should have received a GetStopMonitoring request with:
      | //siri:MonitoringRef | second |

  Scenario: 2939 - Partner Setting collect.priority is used to select the best Partner
    Given a SIRI server "first" waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>first</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>first:ResponseMessage::6ba:LOC</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>first:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And a SIRI server "second" waits GetStopMonitoring request on "http://localhost:8091" to respond with
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>first</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>first:ResponseMessage::6ba:LOC</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>first:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And a Partner "first" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | collect.priority     | 1                     |
      | remote_objectid_kind | external              |
      | remote_credential    | dummy                 |
    And a Partner "second" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8091 |
      | collect.priority     | 2                     |
      | remote_objectid_kind | external              |
      | remote_credential    | dummy                 |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | ObjectIDs       | "external": "single"     |
    When a minute has passed
    Then the "first" SIRI server should not have received a GetStopMonitoring request
    Then the "second" SIRI server should have received a GetStopMonitoring request with:
      | //siri:MonitoringRef | single |

  Scenario: 2481 - Handle a SIRI StopMonitoring request on a unknown StopArea
    Given a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    When I send a SIRI GetStopMonitoring request with
      | RequestorRef  | test    |
      | MonitoringRef | unknown |
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='UTF-8'?>
<S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/' xmlns:SOAP-ENV='http://schemas.xmlsoap.org/soap/envelope/'>
  <S:Body>
    <ns8:GetStopMonitoringResponse xmlns:ns3='http://www.siri.org.uk/siri' xmlns:ns4='http://www.ifopt.org.uk/acsb' xmlns:ns5='http://www.ifopt.org.uk/ifopt' xmlns:ns6='http://datex2.eu/schema/2_0RC1/2_0' xmlns:ns7='http://scma/siri' xmlns:ns8='http://wsdl.siri.org.uk' xmlns:ns9='http://wsdl.siri.org.uk/siri'>
      <ServiceDeliveryInfo>
        <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
        <ns3:ProducerRef>Edwig</ns3:ProducerRef>
        <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
        <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <ns3:StopMonitoringDelivery version='2.0:FR-IDF-2.4'>
          <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
          <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
          <ns3:Status>false</ns3:Status>
          <ns3:ErrorCondition>
            <ns3:InvalidDataReferencesError>
              <ns3:ErrorText>StopArea not found: 'unknown'</ns3:ErrorText>
            </ns3:InvalidDataReferencesError>
          </ns3:ErrorCondition>
        </ns3:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </ns8:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """
