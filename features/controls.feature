Feature: Manages Controls

  @nostart @database
  Scenario: Handle a Situation Control
    Given the table "referentials" has the following data:
      | referential_id                         | slug   | settings | tokens          |
      | '6ba7b814-9dad-11d1-0000-00c04fd430c8' | 'test' | '{}'     | '["testtoken"]' |
    And the table "controls" has the following data:
      | id                                     | referential_slug | context_id | position | type    | model_type  | hook | criticity | internal_code | attributes |
      | '6ba7b814-9dad-11d1-0003-00c04fd430c8' | 'test'           | null       |        0 | 'Dummy' | 'Situation' | null | 'warning' | 'dummy'       | '{}'       |
    And a SIRI server waits SituationExchangeRequest request on "http://localhost:8090" to respond with
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
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>1</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:VersionedAtTime>2017-01-01T01:02:03.000+02:00</siri:VersionedAtTime>
                    <siri:Progress>published</siri:Progress>
                    <siri:Reality>test</siri:Reality>
                     <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-01-01T20:30:06.000+02:00</siri:EndTime>
                    </siri:ValidityPeriod>
                    <siri:AlertCause>maintenanceWork</siri:AlertCause>
                    <siri:Severity>slight</siri:Severity>
                    <siri:Summary>Nouveau pass Navigo</siri:Summary>
                    <siri:Description xml:lang="EN">The new pass is available</siri:Description>
                    <siri:Affects>
                      <siri:Networks>
                        <siri:AffectedNetwork>
                          <siri:AllLines/>
                        </siri:AffectedNetwork>
                      </siri:Networks>
                    </siri:Affects>
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
    And a minute has passed
    And a minute has passed
    Then a Control Message should exist with these attributes:
      | ControlType                      | Dummy                |
      | Criticity                        | warning              |
      | InternalCode                     | dummy                |
      | TargetModelClass                 | Situation            |
      | TargetModelUUID                  |                      |
      | Timestamp                        | 2017-01-01T12:02:00Z |
      | TranslationInfoMessageAttributes | Nouveau pass Navigo  |
      | TranslationInfoMessageKey        | dummy_Situation      |
      | UUID                             |                      |

  @nostart @database
  Scenario: 2461 - Performs a SIRI StopMonitoring request to a Partner
    Given the table "referentials" has the following data:
      | referential_id                         | slug   | settings | tokens          |
      | '6ba7b814-9dad-11d1-0000-00c04fd430c8' | 'test' | '{}'     | '["testtoken"]' |
    And the table "controls" has the following data:
      | id                                     | referential_slug | context_id | position | type         | model_type | hook          | criticity | internal_code | attributes |
      | '6ba7b814-9dad-11d1-0003-00c04fd430c8' | 'test'           | null       |        0 | 'Unexpected' | 'Line'     | 'AfterCreate' | 'warning' | 'unexpected'  | '{}'       |
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
                <siri:AimedArrivalTime>2017-01-01T12:54:00.000+02:00</siri:AimedArrivalTime>
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
      | Name            | Test 1                     |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
    When a minute has passed
    And the SIRI server has received a GetStopMonitoring request
    Then a Control Message should exist with these attributes:
      | ControlType                      | Unexpected           |
      | Criticity                        | warning              |
      | InternalCode                     | unexpected           |
      | TargetModelClass                 | Line                 |
      | TargetModelUUID                  |                      |
      | Timestamp                        | 2017-01-01T12:02:00Z |
      | TranslationInfoMessageAttributes | Ligne 3 Metro        |
      | TranslationInfoMessageKey        | unexpected_line      |
      | UUID                             |                      |

  @ARA-1775 @nostart @database
  Scenario: Handle a PassingTimeControl
    Given the table "referentials" has the following data:
      | referential_id                         | slug   | settings | tokens          |
      | '6ba7b814-9dad-11d1-0000-00c04fd430c8' | 'test' | '{}'     | '["testtoken"]' |
    And the table "controls" has the following data:
      | id                                     | referential_slug | context_id | position | type                    | model_type  | hook                    | criticity | internal_code  | attributes                   |
      | '6ba7b814-9dad-11d1-0003-00c04fd430c8' | 'test'           | null       |        0 | 'PassingTimeChronology' | 'StopVisit' | 'AfterAllStopVisitSave' | 'warning' | 'passing_time' | '{"code_space": "internal"}' |
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
                <siri:AimedArrivalTime>2017-01-01T12:54:00.000+02:00</siri:AimedArrivalTime>
                <siri:ActualArrivalTime>2017-01-01T12:54:00.000+02:00</siri:ActualArrivalTime>
                <siri:ArrivalStatus>arrived</siri:ArrivalStatus>
                <siri:ArrivalBoardingActivity>alighting</siri:ArrivalBoardingActivity>
                <siri:ArrivalStopAssignment>
                  <siri:AimedQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:AimedQuayRef>
                  <siri:ActualQuayRef>NINOXE:StopPoint:Q:50:LOC</siri:ActualQuayRef>
                </siri:ArrivalStopAssignment>
                <siri:AimedDepartureTime>2017-01-01T12:55:00.000+02:00</siri:AimedDepartureTime>
                <siri:ActualDepartureTime>2017-01-01T12:53:00.000+02:00</siri:ActualDepartureTime>
                <siri:DepartureStatus>onTime</siri:DepartureStatus>
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
    Then a Control Message should exist with these attributes:
      | ControlType                      | PassingTimeChronology                                                                 |
      | Criticity                        | warning                                                                               |
      | InternalCode                     | passing_time                                                                          |
      | TargetModelClass                 | StopVisit                                                                             |
      | TargetModelUUID                  |                                                  6ba7b814-9dad-11d1-0008-00c04fd430c8 |
      | Timestamp                        |                                                                  2017-01-01T12:02:00Z |
      | TranslationInfoMessageAttributes | NINOXE:VehicleJourney:201,2017-01-01T12:53:00.000+02:00,2017-01-01T12:54:00.000+02:00 |
      | TranslationInfoMessageKey        | actual_arrival_before_departure                                                       |
      | UUID                             |                                                                                       |

  @ARA-1775 @nostart @database
  Scenario: Handle a PassingTimeControl
    Given the table "referentials" has the following data:
      | referential_id                         | slug   | settings | tokens          |
      | '6ba7b814-9dad-11d1-0000-00c04fd430c8' | 'test' | '{}'     | '["testtoken"]' |
    And the table "controls" has the following data:
      | id                                     | referential_slug | context_id | position | type                    | model_type  | hook                    | criticity | internal_code  | attributes                   |
      | '6ba7b814-9dad-11d1-0003-00c04fd430c8' | 'test'           | null       |        0 | 'PassingTimeChronology' | 'StopVisit' | 'AfterAllStopVisitSave' | 'warning' | 'passing_time' | '{"code_space": "internal"}' |
    And a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
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
        <ns5:RequestMessageRef>Subscription:Test:0</ns5:RequestMessageRef>
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
            <ns5:SubscriberRef>subscriber</ns5:SubscriberRef>
            <ns5:SubscriptionRef>6ba7b814-9dad-11d1-0003-00c04fd430c8</ns5:SubscriptionRef>
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
    When I start Ara
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name            | Test              |
      | Codes[internal] | NINOXE:Line:3:LOC |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect             |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a minute has passed
    When I send this SIRI request
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceDeliveryInfo>
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
      <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
      <siri:RequestMessageRef></siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Notification>
      <siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <siri:ResponseTimestamp>2017-01-01T12:00:20.000Z</siri:ResponseTimestamp>
      <siri:SubscriptionRef>6ba7b814-9dad-11d1-0003-00c04fd430c8</siri:SubscriptionRef>
      <siri:Status>true</siri:Status>
      <siri:EstimatedJourneyVersionFrame>
        <siri:RecordedAtTime>2017-01-01T12:00:20.000Z</siri:RecordedAtTime>
        <siri:EstimatedVehicleJourney>
          <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
          <siri:DirectionRef>Aller</siri:DirectionRef>
          <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
          <siri:Cancellation>true</siri:Cancellation>
          <siri:DestinationRef>ThisIsTheEnd</siri:DestinationRef>
          <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
          <siri:EstimatedCalls>
            <siri:EstimatedCall>
              <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
              <siri:Order>4</siri:Order>
              <siri:StopPointName>Test</siri:StopPointName>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
              <siri:ExpectedDepartureTime>2017-01-01T15:01:02.000Z</siri:ExpectedDepartureTime>
            </siri:EstimatedCall>
            <siri:EstimatedCall>
              <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
              <siri:Order>5</siri:Order>
              <siri:StopPointName>Test2</siri:StopPointName>
              <siri:ExpectedArrivalTime>2017-01-01T15:01:01.000Z</siri:ExpectedArrivalTime>
              <siri:ArrivalStatus>delayed</siri:ArrivalStatus>
            </siri:EstimatedCall>
          </siri:EstimatedCalls>
        </siri:EstimatedVehicleJourney>
      </siri:EstimatedJourneyVersionFrame>
      </siri:EstimatedTimetableDelivery>
      </Notification>
      <SiriExtension/>
      </sw:NotifyEstimatedTimetable>
      </S:Body>
      </S:Envelope>
      """
    And 30 seconds have passed
    Then one VehicleJourney has the following attributes:
      | Codes[internal] | NINOXE:VehicleJourney:201            |
      | LineId          | 6ba7b814-9dad-11d1-0002-00c04fd430c8 |
      | DirectionType   | Aller                                |
      | Id              | 6ba7b814-9dad-11d1-0009-00c04fd430c8 |
      | Cancellation    | true                                 |
    And a Control Message should exist with these attributes:
      | ControlType                      | PassingTimeChronology                                                       |
      | Criticity                        | warning                                                                     |
      | InternalCode                     | passing_time                                                                |
      | TargetModelClass                 | StopVisit                                                                   |
      | TargetModelUUID                  |                                        6ba7b814-9dad-11d1-000a-00c04fd430c8 |
      | Timestamp                        |                                                        2017-01-01T12:01:30Z |
      | TranslationInfoMessageAttributes | NINOXE:VehicleJourney:201,2017-01-01T15:01:02.000Z,2017-01-01T15:01:01.000Z |
      | TranslationInfoMessageKey        | expected_departure_after_next_arrival                                       |
      | UUID                             |                                                                             |
