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
      | Name            | Test 1                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    When a minute has passed
    And the SIRI server has received a GetStopMonitoring request
    Then one StopVisit has the following attributes:
      | Codes[internal]                | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3 |
      | Schedule[expected]#ArrivalTime | 2017-01-01T12:54:00+02:00                              |
      | Schedule[aimed]#ArrivalTime    | 2017-01-01T12:54:00+02:00                              |

  @nostart @database
  Scenario: Handle a Schedule Macro
    Given the table "referentials" has the following data:
      | referential_id                         | slug   | settings | tokens          |
      | '6ba7b814-9dad-11d1-0000-00c04fd430c8' | 'test' | '{}'     | '["testtoken"]' |
    And the table "macros" has the following data:
      | id                                     | referential_slug | context_id                             | position | type           | model_type       | hook | attributes                                                 |
      | '6ba7b814-9dad-11d1-0003-00c04fd430c8' | 'test'           | null                                   |        0 | 'IfAttribute'  | 'VehicleJourney' | null | '{"attribute_name": "DirectionName", "value": "Aller"}'    |
      | '6ba7b814-9dad-11d1-0004-00c04fd430c8' | 'test'           | '6ba7b814-9dad-11d1-0003-00c04fd430c8' |        0 | 'SetAttribute' | 'VehicleJourney' | null | '{"attribute_name": "DirectionType", "value": "outbound"}' |
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
              <siri:DirectionName>Aller</siri:DirectionName>
              <siri:ExternalLineRef>NINOXE:Line:3:LOC</siri:ExternalLineRef>
              <siri:OperatorRef>NINOXE:Company:15563880:LOC</siri:OperatorRef>
              <siri:OriginRef>NINOXE:StopPoint:SP:42:LOC</siri:OriginRef>
              <siri:OriginName>Magicien Noir</siri:OriginName>
              <siri:DestinationRef>NINOXE:StopPoint:SP:62:LOC</siri:DestinationRef>
              <siri:DestinationName>Cimetière des Sauvages</siri:DestinationName>
              <siri:Monitored>true</siri:Monitored>
              <siri:ProgressRate>normalProgress</siri:ProgressRate>
              <siri:CourseOfJourneyRef>201</siri:CourseOfJourneyRef>
              <siri:VehicleRef>NINOXE:Vehicle:23:LOC</siri:VehicleRef>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:Q:50:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Elf Sylvain - Métro (R)</siri:StopPointName>
                <siri:VehicleAtStop>false</siri:VehicleAtStop>
                <siri:AimedArrivalTime>2017-01-01T12:54:00.000+02:00</siri:AimedArrivalTime>
                <siri:ExpectedArrivalTime>2017-01-01T12:54:00.000+02:00</siri:ExpectedArrivalTime>
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
      | Name            | Test 1                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    When a minute has passed
    And the SIRI server has received a GetStopMonitoring request
    Then one VehicleJourney has the following attributes:
      | Codes[internal]          | NINOXE:VehicleJourney:201 |
      | Attributes[DirectionName] | Aller                     |
      | DirectionType            | outbound                  |

  @nostart @database
  Scenario: Handle a Situation Macro
    Given the table "referentials" has the following data:
      | referential_id                         | slug   | settings | tokens          |
      | '6ba7b814-9dad-11d1-0000-00c04fd430c8' | 'test' | '{}'     | '["testtoken"]' |
    And the table "macros" has the following data:
      | id                                     | referential_slug | context_id | position | type                     | model_type  | hook | attributes |
      | '6ba7b814-9dad-11d1-0003-00c04fd430c8' | 'test'           | null       |        0 | 'DefineSituationAffects' | 'Situation' | null | '{}'       |
    And a SIRI server waits GetSituationExchange request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:06.000+02:00</siri:CreationTime>
                    <siri:SituationNumber>test2</siri:SituationNumber>
                    <siri:Version>5</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-01-01T20:30:06.000+02:00</siri:EndTime>
                    </siri:ValidityPeriod>
                    <siri:UndefinedReason/>
                    <siri:Severity>noImpact</siri:Severity>
                    <siri:ReportType>general</siri:ReportType>
                    <siri:Keywords>Commercial Test2</siri:Keywords>
                    <siri:Description>carte d'abonnement</siri:Description>
                    <siri:Affects />
                    <siri:Consequences>
                      <siri:Consequence>
                        <siri:Period>
                          <siri:StartTime>2023-09-18T05:30:59.000Z</siri:StartTime>
                          <siri:EndTime>2023-09-18T08:00:54.000Z</siri:EndTime>
                        </siri:Period>
                        <siri:Severity>verySlight</siri:Severity>
                        <siri:Affects>
                          <siri:Networks>
                            <siri:AffectedNetwork>
                              <siri:AffectedLine>
                                <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                                <siri:Sections>
                                  <siri:AffectedSection>
                                    <siri:IndirectSectionRef>
                                      <siri:FirstStopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:FirstStopPointRef>
                                      <siri:LastStopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:LastStopPointRef>
                                    </siri:IndirectSectionRef>
                                  </siri:AffectedSection>
                                </siri:Sections>
                              </siri:AffectedLine>
                            </siri:AffectedNetwork>
                          </siri:Networks>
                          <siri:StopPoints>
                            <siri:AffectedStopPoint>
                              <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                              <siri:Lines>
                                  <siri:AffectedLine>
                                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                                  </siri:AffectedLine>
                                  <siri:AffectedLine>
                                    <siri:LineRef>NINOXE:Line:BP:LOC</siri:LineRef>
                                  </siri:AffectedLine>
                              </siri:Lines>
                            </siri:AffectedStopPoint>
                          </siri:StopPoints>
                        </siri:Affects>
                        <siri:Blocking>
                          <siri:JourneyPlanner>true</siri:JourneyPlanner>
                          <siri:RealTime>true</siri:RealTime>
                        </siri:Blocking>
                      </siri:Consequence>
                    </siri:Consequences>
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """
    When I start Ara
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-situation-exchange-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | ineo                  |
      | remote_code_space | external              |
    And a Line exists with the following attributes:
      | Codes[external] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a Line exists with the following attributes:
      | Codes[external] | NINOXE:Line:BP:LOC |
      | Name            | Ligne BP Metro     |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[external] | NINOXE:StopPoint:SP:24:LOC |
    And a StopArea exists with the following attributes:
      | Name            | Test last stop             |
      | Codes[external] | NINOXE:StopPoint:SP:25:LOC |
    And a StopArea exists with the following attributes:
      | Name            | Test 3534              |
      | Codes[external] | STIF:StopPoint:Q:3534" |
    And a StopArea exists with the following attributes:
      | Name            | Test 3533              |
      | Codes[external] | STIF:StopPoint:Q:3533: |
    And a minute has passed
    When a minute has passed
    And the SIRI server has received a GetSituationExchange request
    Then one Situation has the following attributes:
      | Codes[external]                                                                  | test2                                |
      | Affects[Line]                                                                    | 6ba7b814-9dad-11d1-0001-00c04fd430c8 |
      | Affects[Line=6ba7b814-9dad-11d1-0001-00c04fd430c8]/AffectedSections[0]/FirstStop | 6ba7b814-9dad-11d1-0003-00c04fd430c8 |
      | Affects[Line=6ba7b814-9dad-11d1-0001-00c04fd430c8]/AffectedSections[0]/LastStop  | 6ba7b814-9dad-11d1-0004-00c04fd430c8 |
      | Affects[StopArea]                                                                | 6ba7b814-9dad-11d1-0003-00c04fd430c8 |
      | Affects[StopArea=6ba7b814-9dad-11d1-0003-00c04fd430c8]/LineIds[0]                | 6ba7b814-9dad-11d1-0001-00c04fd430c8 |
      | Affects[StopArea=6ba7b814-9dad-11d1-0003-00c04fd430c8]/LineIds[1]                | 6ba7b814-9dad-11d1-0002-00c04fd430c8 |

  @nostart @database @ARA-1815
  Scenario: Handle Macro for VehicleJourney Cancellation
    Given the table "referentials" has the following data:
      | referential_id                         | slug   | settings | tokens          |
      | '6ba7b814-9dad-11d1-0000-00c04fd430c8' | 'test' | '{}'     | '["testtoken"]' |
    And the table "macros" has the following data:
      | id                                     | referential_slug | context_id | position | type                            | model_type       | hook | attributes |
      | '6ba7b814-9dad-11d1-0009-00c04fd430c8' | 'test'           | null       |        0 | 'SetVehicleJourneyCancellation' | 'StopVisit' | null | '{}'       |
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
              <siri:DirectionName>Aller</siri:DirectionName>
              <siri:ExternalLineRef>NINOXE:Line:3:LOC</siri:ExternalLineRef>
              <siri:OperatorRef>NINOXE:Company:15563880:LOC</siri:OperatorRef>
              <siri:OriginRef>NINOXE:StopPoint:SP:42:LOC</siri:OriginRef>
              <siri:OriginName>Magicien Noir</siri:OriginName>
              <siri:DestinationRef>NINOXE:StopPoint:SP:62:LOC</siri:DestinationRef>
              <siri:DestinationName>Cimetière des Sauvages</siri:DestinationName>
              <siri:Monitored>true</siri:Monitored>
              <siri:ProgressRate>normalProgress</siri:ProgressRate>
              <siri:CourseOfJourneyRef>201</siri:CourseOfJourneyRef>
              <siri:VehicleRef>NINOXE:Vehicle:23:LOC</siri:VehicleRef>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Test 1</siri:StopPointName>
                <siri:VehicleAtStop>false</siri:VehicleAtStop>
                <siri:AimedArrivalTime>2017-01-01T12:54:00.000+02:00</siri:AimedArrivalTime>
                <siri:ExpectedArrivalTime>2017-01-01T12:54:00.000+02:00</siri:ExpectedArrivalTime>
                <siri:ArrivalStatus>cancelled</siri:ArrivalStatus>
                <siri:ArrivalBoardingActivity>alighting</siri:ArrivalBoardingActivity>
                <siri:ArrivalStopAssignment>
                  <siri:AimedQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:AimedQuayRef>
                  <siri:ActualQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:ActualQuayRef>
                </siri:ArrivalStopAssignment>
                <siri:DepartureStatus>cancelled</siri:DepartureStatus>
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
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name                    | Passage 32                           |
      | Codes[internal]         | NINOXE:VehicleJourney:201            |
      | LineId                  | 6ba7b814-9dad-11d1-0002-00c04fd430c8 |
      | Monitored               | true                                 |
      | DestinationName         | La fin. <TER>                        |
      | HasCompleteStopSequence | true                                 |
    And a StopArea exists with the following attributes:
      | Name            | Test 1                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    And a StopVisit exists with the following attributes:
      | Codes[internal]                | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder                   | 4                                                      |
      | StopAreaId                     | 6ba7b814-9dad-11d1-0004-00c04fd430c8                   |
      | VehicleJourneyId               | 6ba7b814-9dad-11d1-0003-00c04fd430c8                   |
      | VehicleAtStop                  | true                                                   |
      | Reference[OperatorRef]#Code    | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[actual]#Arrival       | 2017-01-01T13:00:00.000Z                               |
      | Attributes[DestinationDisplay] | Cergy le haut & Arret <RER>                            |
    When a minute has passed
    Then one VehicleJourney has the following attributes:
      | Codes[internal] | NINOXE:VehicleJourney:201 |
      | Cancellation    | true                      |
