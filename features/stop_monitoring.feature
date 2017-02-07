Feature: Manage StopMonitoring

  Background:
      Given a Referential "test" is created

  Scenario: Performs a SIRI StopMonitoring request
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
      <Answer>
    </ns8:GetStopMonitoringResponse>
  </S:Body>
</S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | remote_objectid_kind | internal              |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name      | Test 1                                   |
      | ObjectIds | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    When a minute has passed
    And the SIRI server has received a GetStopMonitoring request
    Then one StopVisit has the following attributes:
      | ObjectIds    | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder |                                                                    4 |
    And one Line has the following attributes:
      | ObjectIds | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And one VehicleJourney has the following attributes:
      | ObjectIds | "internal": "NINOXE:VehicleJourney:201" |