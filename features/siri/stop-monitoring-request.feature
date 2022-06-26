Feature: Support SIRI StopMonitoring by request

  Background:
      Given a Referential "test" is created

  Scenario: 2461 - Performs a SIRI StopMonitoring request to a Partner
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:siri="http://www.siri.org.uk/siri"
                                   xmlns:ns4="http://www.ifopt.org.uk/acsb"
                                   xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                                   xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                                   xmlns:ns7="http://scma/siri"
                                   xmlns:sw="http://wsdl.siri.org.uk"
                                   xmlns:ns9="http://wsdl.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</siri:ResponseTimestamp>
        <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
        <siri:Address>http://appli.chouette.mobi/siri_france/siri</siri:Address>
        <siri:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <siri:ResponseTimestamp>2016-09-22T08:01:20.630+02:00</siri:ResponseTimestamp>
          <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
          <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
          <siri:Status>true</siri:Status>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>2016-09-22T07:56:53.000+02:00</siri:RecordedAtTime>
            <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
              <siri:DirectionRef>Left</siri:DirectionRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>2016-09-22</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
              <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
              <siri:DirectionName>Mago-Cime OMNI</siri:DirectionName>
              <siri:ExternalLineRef>NINOXE:Line:3:LOC</siri:ExternalLineRef>
              <siri:OperatorRef>NINOXE:Company:15563880:LOC</siri:OperatorRef>
              <siri:ProductCategoryRef>0</siri:ProductCategoryRef>
              <siri:VehicleFeatureRef>TRFC_M4_1</siri:VehicleFeatureRef>
              <siri:OriginRef>NINOXE:StopPoint:SP:42:LOC</siri:OriginRef>
              <siri:OriginName>Magicien Noir</siri:OriginName>
              <siri:DestinationRef>NINOXE:StopPoint:SP:62:LOC</siri:DestinationRef>
              <siri:DestinationName>Cimetière des Sauvages</siri:DestinationName>
              <siri:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</siri:OriginAimedDepartureTime>
              <siri:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</siri:DestinationAimedArrivalTime>
              <siri:Monitored>true</siri:Monitored>
              <siri:ProgressRate>normalProgress</siri:ProgressRate>
              <siri:Delay>P0Y0M0DT0H0M0.000S</siri:Delay>
              <siri:CourseOfJourneyRef>201</siri:CourseOfJourneyRef>
              <siri:VehicleRef>NINOXE:Vehicle:23:LOC</siri:VehicleRef>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:Q:50:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Elf Sylvain - Métro (R)</siri:StopPointName>
                <siri:VehicleAtStop>false</siri:VehicleAtStop>
                <siri:AimedArrivalTime>2016-09-22T07:54:00.000+02:00</siri:AimedArrivalTime>
                <siri:ActualArrivalTime>2016-09-22T07:54:00.000+02:00</siri:ActualArrivalTime>
                <siri:ArrivalStatus>arrived</siri:ArrivalStatus>
                <siri:ArrivalBoardingActivity>alighting</siri:ArrivalBoardingActivity>
                <siri:ArrivalStopAssignment>
                  <siri:AimedQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:AimedQuayRef>
                  <siri:ActualQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:ActualQuayRef>
                </siri:ArrivalStopAssignment>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                 | http://localhost:8090      |
      | remote_credential          | test                       |
      | remote_objectid_kind       | internal                   |
      | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
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
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |

    Scenario: 2461 - Performs a SIRI StopMonitoring request to a Partner which respond with multiple deliveries
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:siri="http://www.siri.org.uk/siri"
                                   xmlns:ns4="http://www.ifopt.org.uk/acsb"
                                   xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                                   xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                                   xmlns:ns7="http://scma/siri"
                                   xmlns:sw="http://wsdl.siri.org.uk"
                                   xmlns:ns9="http://wsdl.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</siri:ResponseTimestamp>
        <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
        <siri:Address>http://appli.chouette.mobi/siri_france/siri</siri:Address>
        <siri:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <siri:ResponseTimestamp>2016-09-22T08:01:20.630+02:00</siri:ResponseTimestamp>
          <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
          <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
          <siri:Status>true</siri:Status>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>2016-09-22T07:56:53.000+02:00</siri:RecordedAtTime>
            <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
              <siri:DirectionRef>Left</siri:DirectionRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>2016-09-22</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
              <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
              <siri:DirectionName>Mago-Cime OMNI</siri:DirectionName>
              <siri:ExternalLineRef>NINOXE:Line:3:LOC</siri:ExternalLineRef>
              <siri:OperatorRef>NINOXE:Company:15563880:LOC</siri:OperatorRef>
              <siri:ProductCategoryRef>0</siri:ProductCategoryRef>
              <siri:VehicleFeatureRef>TRFC_M4_1</siri:VehicleFeatureRef>
              <siri:OriginRef>NINOXE:StopPoint:SP:42:LOC</siri:OriginRef>
              <siri:OriginName>Magicien Noir</siri:OriginName>
              <siri:DestinationRef>NINOXE:StopPoint:SP:62:LOC</siri:DestinationRef>
              <siri:DestinationName>Cimetière des Sauvages</siri:DestinationName>
              <siri:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</siri:OriginAimedDepartureTime>
              <siri:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</siri:DestinationAimedArrivalTime>
              <siri:Monitored>true</siri:Monitored>
              <siri:ProgressRate>normalProgress</siri:ProgressRate>
              <siri:Delay>P0Y0M0DT0H0M0.000S</siri:Delay>
              <siri:CourseOfJourneyRef>201</siri:CourseOfJourneyRef>
              <siri:VehicleRef>NINOXE:Vehicle:23:LOC</siri:VehicleRef>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:Q:50:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Elf Sylvain - Métro (R)</siri:StopPointName>
                <siri:VehicleAtStop>false</siri:VehicleAtStop>
                <siri:AimedArrivalTime>2016-09-22T07:54:00.000+02:00</siri:AimedArrivalTime>
                <siri:ActualArrivalTime>2016-09-22T07:54:00.000+02:00</siri:ActualArrivalTime>
                <siri:ArrivalStatus>arrived</siri:ArrivalStatus>
                <siri:ArrivalBoardingActivity>alighting</siri:ArrivalBoardingActivity>
                <siri:ArrivalStopAssignment>
                  <siri:AimedQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:AimedQuayRef>
                  <siri:ActualQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:ActualQuayRef>
                </siri:ArrivalStopAssignment>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
        </siri:StopMonitoringDelivery>
        <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <siri:ResponseTimestamp>2016-09-22T08:01:20.630+02:00</siri:ResponseTimestamp>
          <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
          <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
          <siri:Status>true</siri:Status>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>2016-09-22T07:56:53.000+02:00</siri:RecordedAtTime>
            <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-3</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
              <siri:DirectionRef>Left</siri:DirectionRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>2016-09-22</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
              <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
              <siri:DirectionName>Mago-Cime OMNI</siri:DirectionName>
              <siri:ExternalLineRef>NINOXE:Line:3:LOC</siri:ExternalLineRef>
              <siri:OperatorRef>NINOXE:Company:15563880:LOC</siri:OperatorRef>
              <siri:ProductCategoryRef>0</siri:ProductCategoryRef>
              <siri:VehicleFeatureRef>TRFC_M4_1</siri:VehicleFeatureRef>
              <siri:OriginRef>NINOXE:StopPoint:SP:42:LOC</siri:OriginRef>
              <siri:OriginName>Magicien Noir</siri:OriginName>
              <siri:DestinationRef>NINOXE:StopPoint:SP:62:LOC</siri:DestinationRef>
              <siri:DestinationName>Cimetière des Sauvages</siri:DestinationName>
              <siri:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</siri:OriginAimedDepartureTime>
              <siri:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</siri:DestinationAimedArrivalTime>
              <siri:Monitored>true</siri:Monitored>
              <siri:ProgressRate>normalProgress</siri:ProgressRate>
              <siri:Delay>P0Y0M0DT0H0M0.000S</siri:Delay>
              <siri:CourseOfJourneyRef>201</siri:CourseOfJourneyRef>
              <siri:VehicleRef>NINOXE:Vehicle:23:LOC</siri:VehicleRef>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:Q:50:LOC</siri:StopPointRef>
                <siri:Order>5</siri:Order>
                <siri:StopPointName>Elf Sylvain - Métro (R)</siri:StopPointName>
                <siri:VehicleAtStop>false</siri:VehicleAtStop>
                <siri:AimedArrivalTime>2016-09-22T07:54:00.000+02:00</siri:AimedArrivalTime>
                <siri:ActualArrivalTime>2016-09-22T07:54:00.000+02:00</siri:ActualArrivalTime>
                <siri:ArrivalStatus>arrived</siri:ArrivalStatus>
                <siri:ArrivalBoardingActivity>alighting</siri:ArrivalBoardingActivity>
                <siri:ArrivalStopAssignment>
                  <siri:AimedQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:AimedQuayRef>
                  <siri:ActualQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:ActualQuayRef>
                </siri:ArrivalStopAssignment>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                 | http://localhost:8090      |
      | remote_credential          | test                       |
      | remote_objectid_kind       | internal                   |
      | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name      | Test 1                                   |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    When a minute has passed
    And the SIRI server has received a GetStopMonitoring request
    Then one StopVisit has the following attributes:
      | ObjectIDs    | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder | 4                                                                    |
    And one StopVisit has the following attributes:
      | ObjectIDs    | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-3" |
      | PassageOrder | 5                                                                    |
    And one Line has the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And one VehicleJourney has the following attributes:
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |

  Scenario: Handle a SIRI StopMonitoring request
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
      | Monitored | true                                    |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                             |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
          <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
          <siri:Status>true</siri:Status>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
            <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:ActualArrivalTime>2017-01-01T13:00:00.000Z</siri:ActualArrivalTime>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """

  Scenario: Handle a SIRI StopMonitoring request on a 'empty' StopArea
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
<S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery version='2.0:FR-IDF-2.4'>
          <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
          <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
          <siri:Status>true</siri:Status>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """

  Scenario: Handle a SIRI StopMonitoring request by returning all required attributes
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Destination                              |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:62:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Origin                                   |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:42:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Via                                       |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:256:LOC" |
      | Monitored | true                                      |
    And a Line exists with the following attributes:
      | ObjectIDs    | "internal": "NINOXE:Line:3:LOC"           |
      | Name         | Ligne 3 Metro                             |
      | OperationRef | "internal": "NINOXE:Company:15563880:LOC" |
    And a VehicleJourney exists with the following attributes:
      | ObjectIDs                              | "internal": "NINOXE:VehicleJourney:201"         |
      | Name                                   | Magicien Noir - Cimetière (OMNI)                |
      | LineId                                 | 6ba7b814-9dad-11d1-6-00c04fd430c8               |
      | Monitored                              | true                                            |
      | Occupancy                              | 1                                               |
      | Attribute[Bearing]                     | N                                               |
      | Attribute[Delay]                       | 30                                              |
      | DestinationName                        | Cimetière des Sauvages                          |
      | Attribute[DirectionName]               | Mago-Cime OMNI                                  |
      | Attribute[DirectionRef]                | Aller                                           |
      | Attribute[FirstOrLastJourney]          | first                                           |
      | Attribute[HeadwayService]              | false                                           |
      | Attribute[InCongestion]                | false                                           |
      | Attribute[InPanic]                     | false                                           |
      | Attribute[JourneyNote]                 | Note de test                                    |
      | Attribute[JourneyPatternName]          | TEST                                            |
      | Attribute[MonitoringError]             | false                                           |
      | Attribute[OriginAimedDepartureTime]    | 2016-09-22T07:54:52.977Z                        |
      | Attribute[DestinationAimedArrivalTime] | 2016-09-22T09:54:52.977Z                        |
      | OriginName                             | Magicien Noir                                   |
      | Attribute[ProductCategoryRef]          | 0                                               |
      | Attribute[ServiceFeatureRef]           | bus scolaire                                    |
      | Attribute[TrainNumberRef]              | 12345                                           |
      | Attribute[VehicleFeatureRef]           | longTrain                                       |
      | Attribute[VehicleMode]                 | bus                                             |
      | Attribute[ViaPlaceName]                | Saint Bénédicte                                 |
      | Reference[DestinationRef]#ObjectId     | "internal": "NINOXE:StopPoint:SP:62:LOC"        |
      | Reference[JourneyPatternRef]#ObjectId  | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
      | Reference[OriginRef]#ObjectId          | "internal": "NINOXE:StopPoint:SP:42:LOC"        |
      | Reference[RouteRef]#ObjectId           | "internal": "NINOXE:Route:66:LOC"               |
      | Reference[PlaceRef]#ObjectId           | "internal": "NINOXE:StopPoint:SP:256:LOC"       |
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
      | Reference[OperatorRef]#ObjectId      | "internal":"NINOXE:Company:15563880:LOC"                             |
    When I send a SIRI GetStopMonitoring request with
      | RequestorRef  | test                       |
      | MonitoringRef | NINOXE:StopPoint:SP:24:LOC |
    Then I should receive a SIRI GetStopMonitoringResponse with
      | //siri:MonitoredStopVisit[1]/siri:RecordedAtTime                                                                            | 2017-01-01T11:00:00.000Z                                    | StopVisit#RecordedAt                                  |
      | //siri:MonitoredStopVisit[1]/siri:ItemIdentifier                                                                            | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3      | StopVisit#ObjectID                                    |
      | //siri:MonitoredStopVisit[1]/siri:MonitoringRef                                                                             | NINOXE:StopPoint:SP:24:LOC                                  | StopArea#ObjectID                                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:LineRef                                                      | NINOXE:Line:3:LOC                                           | Line#ObjectID                                         |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DirectionRef                                                 | Aller                                                       | VehicleJourney#Attribute[DirectionRef]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DataFrameRef                    | RATPDev:DataFrame::2017-01-01:LOC                           | Model#Date                                            |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef          | NINOXE:VehicleJourney:201                                   | VehicleJourney#ObjectID                               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyPatternRef                                            | NINOXE:JourneyPattern:3_42_62:LOC                           | VehicleJourney#Reference[JourneyPatternRef]#ObjectId  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyPatternName                                           | TEST                                                        | VehicleJourney#Attribute[JourneyPatternName]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleMode                                                  | bus                                                         | VehicleJourney#Attribute[VehicleMode]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:PublishedLineName                                            | Ligne 3 Metro                                               | Line#Name                                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:RouteRef                                                     | RATPDev:Route::720c054714b4464d42970bda37a7edc5af8082cb:LOC | VehicleJourney#Reference[RouteRef]#ObjectId           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DirectionName                                                | Mago-Cime OMNI                                              | VehicleJourney#Attribute[DirectionName]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OperatorRef                                                  | NINOXE:Company:15563880:LOC                                 | StopVisit#Reference[OperatorRef]#ObjectId             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:ProductCategoryRef                                           | 0                                                           | VehicleJourney#Attribute[ProductCategoryRef]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:ServiceFeatureRef                                            | bus scolaire                                                | VehicleJourney#Attribute[ServiceFeatureRef]           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleFeatureRef                                            | longTrain                                                   | VehicleJourney#Attribute[VehicleFeatureRef]           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginRef                                                    | NINOXE:StopPoint:SP:42:LOC                                  | VehicleJourney#Reference[OriginRef]#ObjectId          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginName                                                   | Magicien Noir                                               | VehicleJourney#Attribute[OriginName]                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Via/siri:PlaceName                                           | Saint Bénédicte                                             | VehicleJourney#Attribute[ViaPlaceName]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Via/siri:PlaceRef                                            | NINOXE:StopPoint:SP:256:LOC                                 | VehicleJourney#Reference[PlaceRef]#ObjectId           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationRef                                               | NINOXE:StopPoint:SP:62:LOC                                  | VehicleJourney#Reference[DestinationRef]#ObjectId     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationName                                              | Cimetière des Sauvages                                      | VehicleJourney#Attribute[DestinationName]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                           | Magicien Noir - Cimetière (OMNI)                            | VehicleJourney#Name                                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyNote                                                  | Note de test                                                | VehicleJourney#Attribute[JourneyNote]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:HeadwayService                                               | false                                                       | VehicleJourney#Attribute[HeadwayService]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginAimedDepartureTime                                     | 2016-09-22T07:54:52.977Z                                    | VehicleJourney#Attribute[OriginAimedDepartureTime]    |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationAimedArrivalTime                                  | 2016-09-22T09:54:52.977Z                                    | VehicleJourney#Attribute[DestinationAimedArrivalTime] |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FirstOrLastJourney                                           | first                                                       | VehicleJourney#Attribute[FirstOrLastJourney]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Monitored                                                    | true                                                        | VehicleJourney#Attribute[Monitored]                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoringError                                              | false                                                       | VehicleJourney#Attribute[MonitoringError]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Occupancy                                                    | MANY_SEATS_AVAILABLE                                        | VehicleJourney#Occupancy                              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Delay                                                        | 30                                                          | VehicleJourney#Attribute[Delay]                       |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Bearing                                                      | N                                                           | VehicleJourney#Attribute[Bearing]                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:InPanic                                                      | false                                                       | VehicleJourney#Attribute[InPanic]                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:InCongestion                                                 | false                                                       | VehicleJourney#Attribute[InCongestion]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:TrainNumber/siri:TrainNumberRef                              | 12345                                                       | VehicleJourney#Attribute[TrainNumberRef]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:SituationRef                                                 | 1234556                                                     | TODO                                                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:StopPointRef                              | NINOXE:StopPoint:SP:24:LOC                                  | StopArea#ObjectID                                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:Order                                     | 4                                                           | StopVisit#PassageOrder                                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:VehicleAtStop                             | true                                                        | StopVisit#VehicleAtStop                               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:PlatformTraversal                         | false                                                       | StopVisit#Attribute[PlatformTraversal]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DestinationDisplay                        | Balard Terminus                                             | StopVisit#Attribute[DestinationDisplay]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedArrivalTime                          | 2017-01-01T13:00:00.000Z                                    | StopVisit#Schedule[aimed]#Arrival                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ActualArrivalTime                         | 2017-01-01T13:00:00.000Z                                    | StopVisit#Schedule[actual]#Arrival                    |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedArrivalTime                       | 2017-01-01T13:00:00.000Z                                    | StopVisit#Schedule[expected]#Arrival                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalStatus                             | onTime                                                      | StopVisit#ArrivalStatus                               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalProximyTest                        | A l'approche                                                | StopVisit#Attribute[ArrivalProximyTest]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalPlatformName                       | Platform Name                                               | StopVisit#Attribute[ArrivalPlatformName]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalStopAssignment/siri:ActualQuayName | Quay Name                                                   | StopVisit#Attribute[ActualQuayName]                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedDepartureTime                        | 2017-01-01T13:02:00.000Z                                    | StopVisit#Schedule[aimed]#Departure                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ActualDepartureTime                       | 2017-01-01T13:02:00.000Z                                    | StopVisit#Schedule[actual]#Departure                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedDepartureTime                     | 2017-01-01T13:02:00.000Z                                    | StopVisit#Schedule[expected]#Departure                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DepartureStatus                           | onTime                                                      | StopVisit#DepartureStatus                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DeparturePlatformName                     | Departure Platform Name                                     | StopVisit#Attribute[DeparturePlatformName]            |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DepartureBoardingActivity                 | boarding                                                    | StopVisit#Attribute[DepartureBoardingActivity]        |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedHeadwayInterval                      | 5                                                           | StopVisit#Attribute[AimedHeadwayInterval]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedHeadwayInterval                   | 5                                                           | StopVisit#Attribute[ExpectedHeadwayInterval]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DistanceFromStop                          | 800                                                         | StopVisit#Attribute[DistanceFromStop]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:NumberOfStopsAway                         | 1                                                           | StopVisit#Attribute[NumberOfStopsAway]                |

  Scenario: Handle a SIRI StopMonitoring request by returning all required attributes with the rewrite JourneyPatternRef setting
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                      | test     |
      | remote_objectid_kind                  | internal |
      | broadcast.rewrite_journey_pattern_ref | true     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Destination                              |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:62:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Origin                                   |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:42:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Via                                       |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:256:LOC" |
      | Monitored | true                                      |
    And a Line exists with the following attributes:
      | ObjectIDs    | "internal": "NINOXE:Line:3:LOC"           |
      | Name         | Ligne 3 Metro                             |
      | OperationRef | "internal": "NINOXE:Company:15563880:LOC" |
    And a VehicleJourney exists with the following attributes:
      | ObjectIDs                              | "internal": "NINOXE:VehicleJourney:201"         |
      | Name                                   | Magicien Noir - Cimetière (OMNI)                |
      | LineId                                 | 6ba7b814-9dad-11d1-6-00c04fd430c8               |
      | Monitored                              | true                                            |
      | Attribute[Bearing]                     | N                                               |
      | Attribute[Delay]                       | 30                                              |
      | DestinationName                        | Cimetière des Sauvages                          |
      | Attribute[DirectionName]               | Mago-Cime OMNI                                  |
      | Attribute[DirectionRef]                | Aller                                           |
      | Attribute[FirstOrLastJourney]          | first                                           |
      | Attribute[HeadwayService]              | false                                           |
      | Attribute[InCongestion]                | false                                           |
      | Attribute[InPanic]                     | false                                           |
      | Attribute[JourneyNote]                 | Note de test                                    |
      | Attribute[JourneyPatternName]          | TEST                                            |
      | Attribute[MonitoringError]             | false                                           |
      | Attribute[Occupancy]                   | seatsAvailable                                  |
      | Attribute[OriginAimedDepartureTime]    | 2016-09-22T07:54:52.977Z                        |
      | Attribute[DestinationAimedArrivalTime] | 2016-09-22T09:54:52.977Z                        |
      | OriginName                             | Magicien Noir                                   |
      | Attribute[ProductCategoryRef]          | 0                                               |
      | Attribute[ServiceFeatureRef]           | bus scolaire                                    |
      | Attribute[TrainNumberRef]              | 12345                                           |
      | Attribute[VehicleFeatureRef]           | longTrain                                       |
      | Attribute[VehicleMode]                 | bus                                             |
      | Attribute[ViaPlaceName]                | Saint Bénédicte                                 |
      | Reference[DestinationRef]#ObjectId     | "internal": "NINOXE:StopPoint:SP:62:LOC"        |
      | Reference[JourneyPatternRef]#ObjectId  | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
      | Reference[OriginRef]#ObjectId          | "internal": "NINOXE:StopPoint:SP:42:LOC"        |
      | Reference[RouteRef]#ObjectId           | "internal": "NINOXE:Route:66:LOC"               |
      | Reference[PlaceRef]#ObjectId           | "internal": "NINOXE:StopPoint:SP:256:LOC"       |
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
      | Reference[OperatorRef]#ObjectId      | "internal":"NINOXE:Company:15563880:LOC"                             |
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
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyPatternRef                                            | RATPDev:JourneyPattern::775b650b33aa71eaa01222ccf88a68ce23b58eff:LOC | VehicleJourney#Reference[JourneyPatternRef]#ObjectId  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyPatternName                                           | TEST                                                                 | VehicleJourney#Attribute[JourneyPatternName]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleMode                                                  | bus                                                                  | VehicleJourney#Attribute[VehicleMode]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:PublishedLineName                                            | Ligne 3 Metro                                                        | Line#Name                                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:RouteRef                                                     | RATPDev:Route::720c054714b4464d42970bda37a7edc5af8082cb:LOC          | VehicleJourney#Reference[RouteRef]#ObjectId           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DirectionName                                                | Mago-Cime OMNI                                                       | VehicleJourney#Attribute[DirectionName]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OperatorRef                                                  | NINOXE:Company:15563880:LOC                                          | StopVisit#Reference[OperatorRef]#ObjectId             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:ProductCategoryRef                                           | 0                                                                    | VehicleJourney#Attribute[ProductCategoryRef]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:ServiceFeatureRef                                            | bus scolaire                                                         | VehicleJourney#Attribute[ServiceFeatureRef]           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleFeatureRef                                            | longTrain                                                            | VehicleJourney#Attribute[VehicleFeatureRef]           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginRef                                                    | NINOXE:StopPoint:SP:42:LOC                                           | VehicleJourney#Reference[OriginRef]#ObjectId          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginName                                                   | Magicien Noir                                                        | VehicleJourney#Attribute[OriginName]                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Via/siri:PlaceName                                           | Saint Bénédicte                                                      | VehicleJourney#Attribute[ViaPlaceName]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Via/siri:PlaceRef                                            | NINOXE:StopPoint:SP:256:LOC                                          | VehicleJourney#Reference[PlaceRef]#ObjectId           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationRef                                               | NINOXE:StopPoint:SP:62:LOC                                           | VehicleJourney#Reference[DestinationRef]#ObjectId     |
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

  Scenario: 2466 - Don't perform StopMonitoring request for an unmonitored StopArea
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetStopMonitoringResponse xmlns:siri="http://www.siri.org.uk/siri"
                                         xmlns:ns4="http://www.ifopt.org.uk/acsb"
                                         xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                                         xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                                         xmlns:ns7="http://scma/siri"
                                         xmlns:sw="http://wsdl.siri.org.uk"
                                         xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</siri:ResponseTimestamp>
              <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
              <siri:Address>http://appli.chouette.mobi/siri_france/siri</siri:Address>
              <siri:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2016-09-22T08:01:20.630+02:00</siri:ResponseTimestamp>
                <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>2016-09-22T07:56:53.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:DirectionRef>Left</siri:DirectionRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>2016-09-22</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DirectionName>Mago-Cime OMNI</siri:DirectionName>
                    <siri:ExternalLineRef>NINOXE:Line:3:LOC</siri:ExternalLineRef>
                    <siri:OperatorRef>NINOXE:Company:15563880:LOC</siri:OperatorRef>
                    <siri:ProductCategoryRef>0</siri:ProductCategoryRef>
                    <siri:VehicleFeatureRef>TRFC_M4_1</siri:VehicleFeatureRef>
                    <siri:OriginRef>NINOXE:StopPoint:SP:42:LOC</siri:OriginRef>
                    <siri:OriginName>Magicien Noir</siri:OriginName>
                    <siri:DestinationRef>NINOXE:StopPoint:SP:62:LOC</siri:DestinationRef>
                    <siri:DestinationName>Cimetière des Sauvages</siri:DestinationName>
                    <siri:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</siri:OriginAimedDepartureTime>
                    <siri:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</siri:DestinationAimedArrivalTime>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:ProgressRate>normalProgress</siri:ProgressRate>
                    <siri:Delay>P0Y0M0DT0H0M0.000S</siri:Delay>
                    <siri:CourseOfJourneyRef>201</siri:CourseOfJourneyRef>
                    <siri:VehicleRef>NINOXE:Vehicle:23:LOC</siri:VehicleRef>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:Q:50:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Elf Sylvain - Métro (R)</siri:StopPointName>
                      <siri:VehicleAtStop>false</siri:VehicleAtStop>
                      <siri:AimedArrivalTime>2016-09-22T07:54:00.000+02:00</siri:AimedArrivalTime>
                      <siri:ActualArrivalTime>2016-09-22T07:54:00.000+02:00</siri:ActualArrivalTime>
                      <siri:ArrivalStatus>arrived</siri:ArrivalStatus>
                      <siri:ArrivalBoardingActivity>alighting</siri:ArrivalBoardingActivity>
                      <siri:ArrivalStopAssignment>
                        <siri:AimedQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:AimedQuayRef>
                        <siri:ActualQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:ActualQuayRef>
                      </siri:ArrivalStopAssignment>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
              </siri:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """
    And a Partner "source" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | source                |
      | remote_objectid_kind | internal              |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name            | arrêt 1                                  |
      | ObjectIDs       | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | CollectedAlways | false                                    |
    When a minute has passed
    Then the SIRI server should not have received a GetStopMonitoring request

  Scenario: Perform StopMonitoring request for an unmonitored StopArea
     Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <sw:GetStopMonitoringResponse xmlns:siri="http://www.siri.org.uk/siri"
                                         xmlns:ns4="http://www.ifopt.org.uk/acsb"
                                         xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                                         xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                                         xmlns:ns7="http://scma/siri"
                                         xmlns:sw="http://wsdl.siri.org.uk"
                                         xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</siri:ResponseTimestamp>
              <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
              <siri:Address>http://appli.chouette.mobi/siri_france/siri</siri:Address>
              <siri:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2016-09-22T08:01:20.630+02:00</siri:ResponseTimestamp>
                <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>2016-09-22T07:56:53.000+02:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:DirectionRef>Left</siri:DirectionRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>2016-09-22</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:JourneyPatternRef>NINOXE:JourneyPattern:3_42_62:LOC</siri:JourneyPatternRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:DirectionName>Mago-Cime OMNI</siri:DirectionName>
                    <siri:ExternalLineRef>NINOXE:Line:3:LOC</siri:ExternalLineRef>
                    <siri:OperatorRef>NINOXE:Company:15563880:LOC</siri:OperatorRef>
                    <siri:ProductCategoryRef>0</siri:ProductCategoryRef>
                    <siri:VehicleFeatureRef>TRFC_M4_1</siri:VehicleFeatureRef>
                    <siri:OriginRef>NINOXE:StopPoint:SP:42:LOC</siri:OriginRef>
                    <siri:OriginName>Magicien Noir</siri:OriginName>
                    <siri:DestinationRef>NINOXE:StopPoint:SP:62:LOC</siri:DestinationRef>
                    <siri:DestinationName>Cimetière des Sauvages</siri:DestinationName>
                    <siri:OriginAimedDepartureTime>2016-09-22T07:50:00.000+02:00</siri:OriginAimedDepartureTime>
                    <siri:DestinationAimedArrivalTime>2016-09-22T08:02:00.000+02:00</siri:DestinationAimedArrivalTime>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:ProgressRate>normalProgress</siri:ProgressRate>
                    <siri:Delay>P0Y0M0DT0H0M0.000S</siri:Delay>
                    <siri:CourseOfJourneyRef>201</siri:CourseOfJourneyRef>
                    <siri:VehicleRef>NINOXE:Vehicle:23:LOC</siri:VehicleRef>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:Q:50:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Elf Sylvain - Métro (R)</siri:StopPointName>
                      <siri:VehicleAtStop>false</siri:VehicleAtStop>
                      <siri:AimedArrivalTime>2016-09-22T07:54:00.000+02:00</siri:AimedArrivalTime>
                      <siri:ActualArrivalTime>2016-09-22T07:54:00.000+02:00</siri:ActualArrivalTime>
                      <siri:ArrivalStatus>arrived</siri:ArrivalStatus>
                      <siri:ArrivalBoardingActivity>alighting</siri:ArrivalBoardingActivity>
                      <siri:ArrivalStopAssignment>
                        <siri:AimedQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:AimedQuayRef>
                        <siri:ActualQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:ActualQuayRef>
                      </siri:ArrivalStopAssignment>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
              </siri:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetStopMonitoringResponse>
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
      | Name            | arrêt 1                                  |
      | ObjectIDs       | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | MonitoredAlways | false                                    |
      | CollectedAlways | false                                    |
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

  Scenario: 2481 - Handle a SIRI StopMonitoring request on a unknown StopArea
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    When I send a SIRI GetStopMonitoring request with
      | RequestorRef  | test    |
      | MonitoringRef | unknown |
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:StopMonitoringDelivery version='2.0:FR-IDF-2.4'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:MonitoringRef>unknown</siri:MonitoringRef>
                <siri:Status>false</siri:Status>
                <siri:ErrorCondition>
                  <siri:InvalidDataReferencesError>
                    <siri:ErrorText>StopArea not found: 'unknown'</siri:ErrorText>
                  </siri:InvalidDataReferencesError>
                </siri:ErrorCondition>
              </siri:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: Handle a SIRI StopMonitoring request with descendants
      Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
        | local_credential     | test     |
        | remote_objectid_kind | internal |
      And a StopArea exists with the following attributes:
        | Name      | Parent               |
        | ObjectIDs | "internal": "parent" |
        | Monitored | true                 |
      And a StopArea exists with the following attributes:
        | Name      | Child                             |
        | ObjectIDs | "internal": "child"               |
        | ParentId  | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
        | Monitored | true                              |
      And a Line exists with the following attributes:
        | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
        | Name      | Ligne 3 Metro                   |
      And a Line exists with the following attributes:
        | ObjectIDs | "internal": "NINOXE:Line:4:LOC" |
        | Name      | Ligne 4 Metro                   |
      And a VehicleJourney exists with the following attributes:
        | Name      | Passage 32                              |
        | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
        | LineId    | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
        | Monitored | true                                    |
      And a VehicleJourney exists with the following attributes:
        | Name      | Passage 202                             |
        | ObjectIDs | "internal": "NINOXE:VehicleJourney:202" |
        | LineId    | 6ba7b814-9dad-11d1-5-00c04fd430c8       |
        | Monitored | true                                    |
      And a StopVisit exists with the following attributes:
        | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
        | PassageOrder                    | 4                                                                    |
        | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
        | VehicleJourneyId                | 6ba7b814-9dad-11d1-6-00c04fd430c8                                    |
        | VehicleAtStop                   | true                                                                 |
        | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
        | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                             |
      And a StopVisit exists with the following attributes:
        | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:25:LOC-3" |
        | PassageOrder                    | 4                                                                    |
        | StopAreaId                      | 6ba7b814-9dad-11d1-3-00c04fd430c8                                    |
        | VehicleJourneyId                | 6ba7b814-9dad-11d1-7-00c04fd430c8                                    |
        | VehicleAtStop                   | true                                                                 |
        | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
        | Schedule[actual]#Arrival        | 2017-01-01T14:00:00.000Z                                             |
      When I send a SIRI GetStopMonitoring request with
        | RequestorRef  | test   |
        | MonitoringRef | parent |
      Then I should receive this SIRI response
        """
  <?xml version='1.0' encoding='utf-8'?>
  <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
    <S:Body>
      <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
        <ServiceDeliveryInfo>
          <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:ProducerRef>Ara</siri:ProducerRef>
          <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
          <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
        </ServiceDeliveryInfo>
        <Answer>
          <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
            <siri:MonitoringRef>parent</siri:MonitoringRef>
            <siri:Status>true</siri:Status>
            <siri:MonitoredStopVisit>
              <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
              <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
              <siri:MonitoringRef>parent</siri:MonitoringRef>
              <siri:MonitoredVehicleJourney>
                <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                <siri:FramedVehicleJourneyRef>
                  <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                  <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                </siri:FramedVehicleJourneyRef>
                <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                <siri:Monitored>true</siri:Monitored>
                <siri:MonitoredCall>
                  <siri:StopPointRef>parent</siri:StopPointRef>
                  <siri:Order>4</siri:Order>
                  <siri:StopPointName>Parent</siri:StopPointName>
                  <siri:VehicleAtStop>true</siri:VehicleAtStop>
                  <siri:ActualArrivalTime>2017-01-01T13:00:00.000Z</siri:ActualArrivalTime>
                </siri:MonitoredCall>
              </siri:MonitoredVehicleJourney>
            </siri:MonitoredStopVisit>
            <siri:MonitoredStopVisit>
              <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
              <siri:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:25:LOC-3</siri:ItemIdentifier>
              <siri:MonitoringRef>parent</siri:MonitoringRef>
              <siri:MonitoredVehicleJourney>
                <siri:LineRef>NINOXE:Line:4:LOC</siri:LineRef>
                <siri:FramedVehicleJourneyRef>
                  <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                  <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
                </siri:FramedVehicleJourneyRef>
                <siri:PublishedLineName>Ligne 4 Metro</siri:PublishedLineName>
                <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                <siri:VehicleJourneyName>Passage 202</siri:VehicleJourneyName>
                <siri:Monitored>true</siri:Monitored>
                <siri:MonitoredCall>
                  <siri:StopPointRef>child</siri:StopPointRef>
                  <siri:Order>4</siri:Order>
                  <siri:StopPointName>Child</siri:StopPointName>
                  <siri:VehicleAtStop>true</siri:VehicleAtStop>
                  <siri:ActualArrivalTime>2017-01-01T14:00:00.000Z</siri:ActualArrivalTime>
                </siri:MonitoredCall>
              </siri:MonitoredVehicleJourney>
            </siri:MonitoredStopVisit>
          </siri:StopMonitoringDelivery>
        </Answer>
        <AnswerExtension/>
      </sw:GetStopMonitoringResponse>
    </S:Body>
  </S:Envelope>
        """

  Scenario: Don't perform a SIRI StopMonitoring request when StopArea has a parent
      Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
        """
        """
      And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
        | remote_url                 | http://localhost:8090      |
        | remote_credential          | test                       |
        | remote_objectid_kind       | internal                   |
        | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
      And a minute has passed
      And a StopArea exists with the following attributes:
        | Name      | Test 1                                   |
        | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
        | ParentId  | 6ba7b814-9dad-11d1-4-00c04fd430c8        |
      And a StopArea exists with the following attributes:
        | Name      | Test 2                                   |
        | ObjectIDs | "internal": "NINOXE:StopPoint:SP:25:LOC" |
      When a minute has passed
      Then the SIRI server should not have received a GetStopMonitoring request

  Scenario: Perform a SIRI StopMonitoring request when StopArea has a parent with CollectChildren
      Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
        """
        """
      And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
        | remote_url                 | http://localhost:8090      |
        | remote_credential          | test                       |
        | remote_objectid_kind       | internal                   |
        | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
      And a minute has passed
      And a StopArea exists with the following attributes:
        | Name      | Test 1                                   |
        | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
        | ParentId  | 6ba7b814-9dad-11d1-5-00c04fd430c8        |
      And a StopArea exists with the following attributes:
        | Name            | Test 2                                   |
        | ObjectIDs       | "internal": "NINOXE:StopPoint:SP:25:LOC" |
        | CollectChildren | true                                     |
      When a minute has passed
      Then the SIRI server should have received 1 GetStopMonitoring request

  Scenario: Handle a SIRI StopMonitoring request with Operator
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | external |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "external": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | ObjectIDs | "external": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "external": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
      | Monitored | true                                    |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "external": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectId | "internal": "internalOperator"                                       |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                             |
    And an Operator exists with the following attributes:
      | Name      | Operator                                                    |
      | ObjectIDs | "internal":"internalOperator","external":"externalOperator" |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
          <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
          <siri:Status>true</siri:Status>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
            <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
              <siri:OperatorRef>externalOperator</siri:OperatorRef>
              <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:ActualArrivalTime>2017-01-01T13:00:00.000Z</siri:ActualArrivalTime>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """

  Scenario: Handle a SIRI StopMonitoring request with a not monitored StopArea
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | external |
    And a StopArea exists with the following attributes:
      | Name             | Test                                     |
      | ObjectIDs        | "external": "NINOXE:StopPoint:SP:24:LOC" |
      | Origin[partner1] | false                                    |
      | Monitored        | false                                    |
    And a Line exists with the following attributes:
      | ObjectIDs | "external": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "external": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
      | Monitored | true                                    |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "external": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectId | "internal": "operator"                                               |
      | Schedule[aimed]#Arrival         | 2017-01-01T13:00:00.000Z                                             |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:05.000Z                                             |
    And an Operator exists with the following attributes:
      | Name      | Operator                                                    |
      | ObjectIDs | "internal":"internalOperator","external":"externalOperator" |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
          <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
          <siri:Status>true</siri:Status>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
            <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
              <siri:OperatorRef>operator</siri:OperatorRef>
              <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
              <siri:Monitored>false</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:AimedArrivalTime>2017-01-01T13:00:00.000Z</siri:AimedArrivalTime>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """

  @ARA-1010
  Scenario: Handle a SIRI StopMonitoring request with a not monitored StopArea and broadcast.send_producer_unavailable_error setting
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                          | test     |
      | remote_objectid_kind                      | external |
      | broadcast.send_producer_unavailable_error | true     |
    And a StopArea exists with the following attributes:
      | Name             | Test                                     |
      | ObjectIDs        | "external": "NINOXE:StopPoint:SP:24:LOC" |
      | Origin[partner1] | false                                    |
      | Monitored        | false                                    |
    And a Line exists with the following attributes:
      | ObjectIDs | "external": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "external": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
      | Monitored | true                                    |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "external": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectId | "internal": "operator"                                               |
      | Schedule[aimed]#Arrival         | 2017-01-01T13:00:00.000Z                                             |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:05.000Z                                             |
    And an Operator exists with the following attributes:
      | Name      | Operator                                                    |
      | ObjectIDs | "internal":"internalOperator","external":"externalOperator" |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
          <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
          <siri:Status>false</siri:Status>
          <siri:ErrorCondition>
            <siri:OtherError number="1">
              <siri:ErrorText>Erreur [PRODUCER_UNAVAILABLE] : partner1 indisponible</siri:ErrorText>
            </siri:OtherError>
          </siri:ErrorCondition>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
            <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
              <siri:OperatorRef>operator</siri:OperatorRef>
              <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
              <siri:Monitored>false</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:AimedArrivalTime>2017-01-01T13:00:00.000Z</siri:AimedArrivalTime>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """

  @ARA-1044
  Scenario: Handle a SIRI StopMonitoring request with unmatching objectid kind
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test  |
      | remote_objectid_kind | wrong |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
      | Monitored | true                                    |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "other": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                 |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                 |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                 |
      | VehicleAtStop                   | true                                                              |
      | Reference[OperatorRef]#ObjectId | "other": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                          |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>Ara</siri:ProducerRef>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
          <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
          <siri:Status>false</siri:Status>
          <siri:ErrorCondition>
            <siri:InvalidDataReferencesError>
            <siri:ErrorText>StopArea not found: 'NINOXE:StopPoint:SP:24:LOC'</siri:ErrorText>
            </siri:InvalidDataReferencesError>
          </siri:ErrorCondition>
        </siri:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension/>
    </sw:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """

  @ARA-1044
  Scenario: Handle a SIRI StopMonitoring request with global setting vehicle_journey_remote_objectid_kind
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                     | test     |
      | remote_objectid_kind                 | internal |
      | vehicle_journey_remote_objectid_kind | other    |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                           |
      | ObjectIDs | "other": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8    |
      | Monitored | true                                 |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                             |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
    Then I should receive a SIRI GetStopMonitoringResponse with
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef | NINOXE:VehicleJourney:201 | VehicleJourney#ObjectID |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                  | Passage 32                | VehicleJourney#Name     |

  @ARA-1044
  Scenario: Handle a SIRI StopMonitoring request with connector setting siri-stop-monitoring-request-broadcaster.vehicle_journey_remote_objectid_kind
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                                                              | test     |
      | remote_objectid_kind                                                          | internal |
      | siri-stop-monitoring-request-broadcaster.vehicle_journey_remote_objectid_kind | other    |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                           |
      | ObjectIDs | "other": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8    |
      | Monitored | true                                 |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                             |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
    Then I should receive a SIRI GetStopMonitoringResponse with
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef | NINOXE:VehicleJourney:201 | VehicleJourney#ObjectID |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                  | Passage 32                | VehicleJourney#Name     |

  @ARA-1044
  Scenario: Handle a SIRI StopMonitoring request with multiple connector setting siri-stop-monitoring-request-broadcaster.vehicle_journey_remote_objectid_kind
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                                                              | test          |
      | remote_objectid_kind                                                          | internal      |
      | siri-stop-monitoring-request-broadcaster.vehicle_journey_remote_objectid_kind | other, other2 |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                           |
      | ObjectIDs | "other": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8    |
      | Monitored | true                                 |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectId | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                             |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
    Then I should receive a SIRI GetStopMonitoringResponse with
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef | NINOXE:VehicleJourney:201 | VehicleJourney#ObjectID |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                  | Passage 32                | VehicleJourney#Name     |

  @ARA-1044
  Scenario: Handle a SIRI StopMonitoring request with fallback on generic connector settings remote_objectid_kind
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                                              | test     |
      | remote_objectid_kind                                          | internal |
      | siri-stop-monitoring-request-broadcaster.remote_objectid_kind | other    |
    And a StopArea exists with the following attributes:
      | Name      | Test                                  |
      | ObjectIDs | "other": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                  |
    And a Line exists with the following attributes:
      | ObjectIDs | "other": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                           |
      | ObjectIDs | "other": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8    |
      | Monitored | true                                 |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "other": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                 |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                 |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                 |
      | VehicleAtStop                   | true                                                              |
      | Reference[OperatorRef]#ObjectId | "other": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival        | 2017-01-01T13:00:00.000Z                                          |
    When I send this SIRI request
      """
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
    Then I should receive a SIRI GetStopMonitoringResponse with
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef | NINOXE:VehicleJourney:201 | VehicleJourney#ObjectID |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                  | Passage 32                | VehicleJourney#Name     |
