Feature: Manages Macros

  @nostart @database
  Scenario: Handle a Schedule Macro
    Given the table "referentials" has the following data:
      | referential_id                         | slug   | settings | tokens          |
      | '6ba7b814-9dad-11d1-0000-00c04fd430c8' | 'test' | '{}'     | '["testtoken"]' |
    And the table "macros" has the following data:
      | id                                     | referential_slug | context_id | position | type                        | model_type  | hook | attributes |
      | '6ba7b814-9dad-11d1-0003-00c04fd430c8' | 'test'           | null       |        0 | 'DefineAimedScheduledTimes' | 'StopVisit' | null | '{}'       |
    And a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:GetStopMonitoringResponse xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo>
        <siri:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</siri:ResponseTimestamp>
        <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
        <siri:Address>http://appli.chouette.mobi/siri_france/siri</siri:Address>
        <siri:ResponseMessageIdentifier>fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery>
          <siri:ResponseTimestamp>2016-09-22T08:01:20.630+02:00</siri:ResponseTimestamp>
          <siri:Status>true</siri:Status>
          <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
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
                <siri:ExpectedArrivalTime>2017-01-01T12:54:00.000+02:00</siri:ExpectedArrivalTime>
                <siri:ActualArrivalTime>2017-01-01T12:54:00.000+02:00</siri:ActualArrivalTime>
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
    When I start Ara
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                 | http://localhost:8090      |
      | remote_credential          | test                       |
      | remote_code_space          | internal                   |
      | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name  | Test 1                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    When a minute has passed
    And the SIRI server has received a GetStopMonitoring request
    Then one StopVisit has the following attributes:
      | Codes                          | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | Schedule[expected]#ArrivalTime |                                            2017-01-01T12:54:00+02:00 |
      | Schedule[actual]#ArrivalTime   |                                            2017-01-01T12:54:00+02:00 |
      | Schedule[aimed]#ArrivalTime    |                                            2017-01-01T12:54:00+02:00 |
