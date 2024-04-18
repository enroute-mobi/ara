Feature: Support SIRI StopMonitoring by request

  Background:
    Given a Referential "test" is created

  @ARA-1240
  Scenario: Collect by using SIRI Lite Stop Monitoring
    Given a lite SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      {
      "Siri": {
      "ServiceDelivery": {
        "ResponseTimestamp": "2023-06-02T11:16:11.127Z",
        "ProducerRef": "IVTR_HET",
        "ResponseMessageIdentifier": "IVTR_HET:ResponseMessage:9bd2199f-2685-4f37-9a60-177312447a38:LOC:",
        "StopMonitoringDelivery": [
          {
            "ResponseTimestamp": "2023-06-02T11:16:11.249Z",
            "Version": "2.0",
            "Status": "true",
            "MonitoredStopVisit": [
              {
                "RecordedAtTime": "2023-06-02T01:07:19.892Z",
                "ItemIdentifier": "SNCF_ACCES_CLOUD:Item::41178_133528:LOC",
                "MonitoringRef": "STIF:StopPoint:Q:41178:",
                "MonitoredVehicleJourney": {
                  "LineRef": "STIF:Line::C01740:",
                  "OperatorRef": "SNCF_ACCES_CLOUD:Operator::SNCF:",
                  "FramedVehicleJourneyRef": {
                    "DataFrameRef": "any",
                    "DatedVehicleJourneyRef": "SNCF_ACCES_CLOUD:VehicleJourney::2e484a6e-2359-4cb2-95e1-4483d547aa5a:LOC"
                  },
                  "DestinationRef": "STIF:StopPoint:Q:41194:",
                  "DestinationName": "Gare Saint-Lazare",
                  "JourneyNote": "PASA",
                  "MonitoredCall": {
                    "StopPointName": "Gare de Saint-Cloud",
                    "VehicleAtStop": false,
                    "DestinationDisplay": "Gare Saint-Lazare",
                    "ExpectedArrivalTime": "2023-06-02T08:46:40.000Z",
                    "ExpectedDepartureTime": "2023-06-02T08:47:40.000Z",
                    "DepartureStatus": "onTime",
                    "Order": 6,
                    "AimedArrivalTime": "2023-06-02T08:46:40.000Z",
                    "ArrivalPlatformName": "2",
                    "AimedDepartureTime": "2023-06-02T08:47:40.000Z",
                    "ArrivalStatus": "onTime"
                  }
                }
              }
            ]
          }
        ]
      }
      }
      }
      """
    And a Partner "test" exists with connectors [siri-lite-stop-monitoring-request-collector] and the following settings:
      | remote_url                       | http://localhost:8090   |
      | remote_credential                | test                    |
      | remote_code_space                | internal                |
      | collect.include_stop_areas       | STIF:StopPoint:Q:41178: |
      | collect.subscriptions.persistent | true                    |
      | local_credential                 | toto                    |
      | collect.persistent               | true                    |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name  | Test 1                                |
      | Codes | "internal": "STIF:StopPoint:Q:41178:" |
    # Id 6ba7b814-9dad-11d1-2-00c04fd430c8
    When a minute has passed
    Then one StopVisit has the following attributes:
      | Codes           | "internal": "SNCF_ACCES_CLOUD:Item::41178_133528:LOC" |
      | ArrivalStatus   | onTime                                                |
      | DepartureStatus | onTime                                                |
      | DataFrameRef    | any                                                   |
      | PassageOrder    |                                                     6 |
      | StopAreaId      |                     6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleAtStop   | false                                                 |
    And one Line has the following attributes:
      | Codes | "internal": "STIF:Line::C01740:" |
    And one VehicleJourney has the following attributes:
      | Codes           | "internal": "SNCF_ACCES_CLOUD:VehicleJourney::2e484a6e-2359-4cb2-95e1-4483d547aa5a:LOC" |
      | DestinationName | Gare Saint-Lazare                                                                       |
      | Monitored       | true                                                                                    |
    And an audit event should exist with these attributes:
      | Protocol           | siri                                                                          |
      | Direction          | sent                                                                          |
      | ResponseIdentifier | IVTR_HET:ResponseMessage:9bd2199f-2685-4f37-9a60-177312447a38:LOC:            |
      | Status             | OK                                                                            |
      | Type               | StopMonitoringRequest                                                         |
      | StopAreas          | ["STIF:StopPoint:Q:41178:"]                                                   |
      | Lines              | ["STIF:Line::C01740:"]                                                        |
      | VehicleJourneys    | ["SNCF_ACCES_CLOUD:VehicleJourney::2e484a6e-2359-4cb2-95e1-4483d547aa5a:LOC"] |
      | RequestRawMessage  | MonitoringRef=STIF:StopPoint:Q:41178:                                         |

  @ARA-1240
  Scenario: Collect by using SIRI Lite Stop Monitoring with invalid SIRI Lite 'value'
    Given a lite SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      {
      "Siri": {
      "ServiceDelivery": {
        "ResponseTimestamp": "2023-06-02T11:16:11.127Z",
        "ProducerRef": "IVTR_HET",
        "ResponseMessageIdentifier": "IVTR_HET:ResponseMessage:9bd2199f-2685-4f37-9a60-177312447a38:LOC:",
        "StopMonitoringDelivery": [
          {
            "ResponseTimestamp": "2023-06-02T11:16:11.249Z",
            "Version": "2.0",
            "Status": "true",
            "MonitoredStopVisit": [
              {
                "RecordedAtTime": "2023-06-02T01:07:19.892Z",
                "ItemIdentifier": "SNCF_ACCES_CLOUD:Item::41178_133528:LOC",
                "MonitoringRef": "STIF:StopPoint:Q:41178:",
                "MonitoredVehicleJourney": {
                  "LineRef": {
                    "value": "STIF:Line::C01740:"
                    },
                  "OperatorRef": {
                    "value": "SNCF_ACCES_CLOUD:Operator::SNCF:"
                    },
                  "FramedVehicleJourneyRef": {
                    "DataFrameRef": {
                      "value": "any"
                      },
                    "DatedVehicleJourneyRef": "SNCF_ACCES_CLOUD:VehicleJourney::2e484a6e-2359-4cb2-95e1-4483d547aa5a:LOC"
                  },
                  "DestinationRef": {
                    "value": "STIF:StopPoint:Q:41194:"
                  },
                  "DestinationName": [
                    {
                      "value": "Gare Saint-Lazare"
                    }
                  ],
                  "JourneyNote": "PASA",
                  "MonitoredCall": {
                    "StopPointName": {
                      "value": "Gare de Saint-Cloud"
                      },
                    "VehicleAtStop": false,
                    "DestinationDisplay": {
                      "value": "Gare Saint-Lazare"
                      },
                    "ExpectedArrivalTime": "2023-06-02T08:46:40.000Z",
                    "ExpectedDepartureTime": "2023-06-02T08:47:40.000Z",
                    "DepartureStatus": "onTime",
                    "Order": 6,
                    "AimedArrivalTime": "2023-06-02T08:46:40.000Z",
                    "ArrivalPlatformName": "2",
                    "AimedDepartureTime": "2023-06-02T08:47:40.000Z",
                    "ArrivalStatus": "onTime"
                  }
                }
              }
            ]
          }
        ]
      }
      }
      }
      """
    And a Partner "test" exists with connectors [siri-lite-stop-monitoring-request-collector] and the following settings:
      | remote_url                       | http://localhost:8090   |
      | remote_credential                | test                    |
      | remote_code_space                | internal                |
      | collect.include_stop_areas       | STIF:StopPoint:Q:41178: |
      | collect.subscriptions.persistent | true                    |
      | local_credential                 | toto                    |
      | collect.persistent               | true                    |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name  | Test 1                                |
      | Codes | "internal": "STIF:StopPoint:Q:41178:" |
    When a minute has passed
    Then one StopVisit has the following attributes:
      | Codes           | "internal": "SNCF_ACCES_CLOUD:Item::41178_133528:LOC" |
      | ArrivalStatus   | onTime                                                |
      | DepartureStatus | onTime                                                |
      | DataFrameRef    | any                                                   |
      | PassageOrder    |                                                     6 |
      | StopAreaId      |                     6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleAtStop   | false                                                 |
    And one Line has the following attributes:
      | Codes | "internal": "STIF:Line::C01740:" |
    And one VehicleJourney has the following attributes:
      | Codes           | "internal": "SNCF_ACCES_CLOUD:VehicleJourney::2e484a6e-2359-4cb2-95e1-4483d547aa5a:LOC" |
      | DestinationName | Gare Saint-Lazare                                                                       |
      | Monitored       | true                                                                                    |
    And an audit event should exist with these attributes:
      | Protocol           | siri                                                               |
      | Direction          | sent                                                               |
      | ResponseIdentifier | IVTR_HET:ResponseMessage:9bd2199f-2685-4f37-9a60-177312447a38:LOC: |
      | Status             | OK                                                                 |
      | Type               | StopMonitoringRequest                                              |
      | StopAreas          | ["STIF:StopPoint:Q:41178:"]                                        |
      | RequestRawMessage  | MonitoringRef=STIF:StopPoint:Q:41178:                              |

  @ARA-1240
  Scenario: Handle a missing StopVisit ItemIdentifier in collect by using SIRI Lite Stop Monitoring should use DatedVehicleJourneyRef and Order to build Code
    Given a lite SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      {
      "Siri": {
      "ServiceDelivery": {
        "ResponseTimestamp": "2023-06-02T11:16:11.127Z",
        "ProducerRef": "IVTR_HET",
        "ResponseMessageIdentifier": "IVTR_HET:ResponseMessage:9bd2199f-2685-4f37-9a60-177312447a38:LOC:",
        "StopMonitoringDelivery": [
          {
            "ResponseTimestamp": "2023-06-02T11:16:11.249Z",
            "Version": "2.0",
            "Status": "true",
            "MonitoredStopVisit": [
              {
                "RecordedAtTime": "2023-06-02T01:07:19.892Z",
                "MonitoringRef": "STIF:StopPoint:Q:41178:",
                "MonitoredVehicleJourney": {
                  "LineRef": "STIF:Line::C01740:",
                  "OperatorRef": "SNCF_ACCES_CLOUD:Operator::SNCF:",
                  "FramedVehicleJourneyRef": {
                    "DataFrameRef": "any",
                    "DatedVehicleJourneyRef": "SNCF_ACCES_CLOUD:VehicleJourney::2e484a6e-2359-4cb2-95e1-4483d547aa5a:LOC"
                  },
                  "DestinationRef": "STIF:StopPoint:Q:41194:",
                  "DestinationName": "Gare Saint-Lazare",
                  "JourneyNote": "PASA",
                  "MonitoredCall": {
                    "StopPointName": "Gare de Saint-Cloud",
                    "VehicleAtStop": false,
                    "DestinationDisplay": "Gare Saint-Lazare",
                    "ExpectedArrivalTime": "2023-06-02T08:46:40.000Z",
                    "ExpectedDepartureTime": "2023-06-02T08:47:40.000Z",
                    "DepartureStatus": "onTime",
                    "Order": 6,
                    "AimedArrivalTime": "2023-06-02T08:46:40.000Z",
                    "ArrivalPlatformName": "2",
                    "AimedDepartureTime": "2023-06-02T08:47:40.000Z",
                    "ArrivalStatus": "onTime"
                  }
                }
              }
            ]
          }
        ]
      }
      }
      }
      """
    And a Partner "test" exists with connectors [siri-lite-stop-monitoring-request-collector] and the following settings:
      | remote_url                       | http://localhost:8090   |
      | remote_credential                | test                    |
      | remote_code_space                | internal                |
      | collect.include_stop_areas       | STIF:StopPoint:Q:41178: |
      | collect.subscriptions.persistent | true                    |
      | local_credential                 | toto                    |
      | collect.persistent               | true                    |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name  | Test 1                                |
      | Codes | "internal": "STIF:StopPoint:Q:41178:" |
    When a minute has passed
    Then one StopVisit has the following attributes:
      | Codes | "internal": "SNCF_ACCES_CLOUD:VehicleJourney::2e484a6e-2359-4cb2-95e1-4483d547aa5a:LOC-6" |

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
      | Codes        | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder |                                                                    4 |
    And one Line has the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And one VehicleJourney has the following attributes:
      | Codes | "internal": "NINOXE:VehicleJourney:201" |
    And an audit event should exist with these attributes:
      | Protocol           | siri                          |
      | Direction          | sent                          |
      | ResponseIdentifier | /{uuid}/                      |
      | Status             | OK                            |
      | Type               | StopMonitoringRequest         |
      | StopAreas          | ["NINOXE:StopPoint:Q:50:LOC"] |
      | VehicleJourneys    | ["NINOXE:VehicleJourney:201"] |
      | Lines              | ["NINOXE:Line:3:LOC"]         |

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
              <siri:OriginAimedDepartureTime>2017-01-01T12:50:00.000+02:00</siri:OriginAimedDepartureTime>
              <siri:DestinationAimedArrivalTime>2017-01-01T13:02:00.000+02:00</siri:DestinationAimedArrivalTime>
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
              <siri:OriginAimedDepartureTime>2017-01-01T12:50:00.000+02:00</siri:OriginAimedDepartureTime>
              <siri:DestinationAimedArrivalTime>2016-01-01T13:02:00.000+02:00</siri:DestinationAimedArrivalTime>
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
      | Codes        | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder |                                                                    4 |
    And one StopVisit has the following attributes:
      | Codes        | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:25:LOC-3" |
      | PassageOrder |                                                                    5 |
    And one Line has the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And one VehicleJourney has the following attributes:
      | Codes | "internal": "NINOXE:VehicleJourney:201" |

  Scenario: Handle a SIRI StopMonitoring request
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:201" |
      | LineId          |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                  |                                                                    4 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop                 | true                                                                 |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival      |                                             2017-01-01T13:00:00.000Z |
      | Attribute[DestinationDisplay] | Cergy le haut & Arret <RER>                                          |
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
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
              <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:DestinationDisplay>Cergy le haut &amp; Arret &lt;RER&gt;</siri:DestinationDisplay>
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
    And an audit event should exist with these attributes:
      | Protocol        | siri                           |
      | Direction       | received                       |
      | Status          | OK                             |
      | Type            | StopMonitoringRequest          |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"] |
      | VehicleJourneys | ["NINOXE:VehicleJourney:201"]  |
      | Lines           | ["NINOXE:Line:3:LOC"]          |

  Scenario: Handle a SIRI StopMonitoring request on a 'empty' StopArea
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
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
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Destination                              |
      | Codes     | "internal": "NINOXE:StopPoint:SP:62:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Origin                                   |
      | Codes     | "internal": "NINOXE:StopPoint:SP:42:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Via                                       |
      | Codes     | "internal": "NINOXE:StopPoint:SP:256:LOC" |
      | Monitored | true                                      |
    And a Line exists with the following attributes:
      | Codes        | "internal": "NINOXE:Line:3:LOC"           |
      | Name         | Ligne 3 Metro                             |
      | OperationRef | "internal": "NINOXE:Company:15563880:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Codes                                  | "internal": "NINOXE:VehicleJourney:201"         |
      | Name                                   | Magicien Noir - Cimetière (OMNI)                |
      | LineId                                 |               6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | Monitored                              | true                                            |
      | Occupancy                              | manySeatsAvailable                              |
      | Attribute[Bearing]                     | N                                               |
      | Attribute[Delay]                       |                                              30 |
      | DestinationName                        | Cimetière des Sauvages                          |
      | Attribute[DirectionName]               | Mago-Cime OMNI                                  |
      | DirectionType                          | Aller                                           |
      | Attribute[FirstOrLastJourney]          | first                                           |
      | Attribute[HeadwayService]              | false                                           |
      | Attribute[InCongestion]                | false                                           |
      | Attribute[InPanic]                     | false                                           |
      | Attribute[JourneyNote]                 | Note de test                                    |
      | Attribute[JourneyPatternName]          | TEST                                            |
      | Attribute[MonitoringError]             | false                                           |
      | Attribute[OriginAimedDepartureTime]    |                        2016-09-22T07:54:52.977Z |
      | Attribute[DestinationAimedArrivalTime] |                        2016-09-22T09:54:52.977Z |
      | OriginName                             | Magicien Noir                                   |
      | Attribute[ProductCategoryRef]          |                                               0 |
      | Attribute[ServiceFeatureRef]           | bus scolaire                                    |
      | Attribute[TrainNumberRef]              |                                           12345 |
      | Attribute[VehicleFeatureRef]           | longTrain                                       |
      | Attribute[VehicleMode]                 | bus                                             |
      | Attribute[ViaPlaceName]                | Saint Bénédicte                                 |
      | Reference[DestinationRef]#Code         | "internal": "NINOXE:StopPoint:SP:62:LOC"        |
      | Reference[JourneyPatternRef]#Code      | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
      | Reference[OriginRef]#Code              | "internal": "NINOXE:StopPoint:SP:42:LOC"        |
      | Reference[RouteRef]#Code               | "internal": "NINOXE:Route:66:LOC"               |
      | Reference[PlaceRef]#Code               | "internal": "NINOXE:StopPoint:SP:256:LOC"       |
    And a StopVisit exists with the following attributes:
      | ArrivalStatus                        | onTime                                                               |
      | DepartureStatus                      | onTime                                                               |
      | Codes                                | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                         |                                                                    4 |
      | RecordedAt                           |                                             2017-01-01T11:00:00.000Z |
      | Schedule[actual]#Arrival             |                                             2017-01-01T13:00:00.000Z |
      | Schedule[actual]#Departure           |                                             2017-01-01T13:02:00.000Z |
      | Schedule[aimed]#Arrival              |                                             2017-01-01T13:00:00.000Z |
      | Schedule[aimed]#Departure            |                                             2017-01-01T13:02:00.000Z |
      | Schedule[expected]#Arrival           |                                             2017-01-01T13:00:00.000Z |
      | Schedule[expected]#Departure         |                                             2017-01-01T13:02:00.000Z |
      | StopAreaId                           |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId                     |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop                        | true                                                                 |
      | Attribute[AimedHeadwayInterval]      |                                                                    5 |
      | Attribute[ActualQuayName]            | Quay Name                                                            |
      | Attribute[ArrivalPlatformName]       | Platform Name                                                        |
      | Attribute[ArrivalProximyTest]        | A l'approche                                                         |
      | Attribute[DepartureBoardingActivity] | boarding                                                             |
      | Attribute[DeparturePlatformName]     | Departure Platform Name                                              |
      | Attribute[DestinationDisplay]        | Balard Terminus                                                      |
      | Attribute[DistanceFromStop]          |                                                                  800 |
      | Attribute[ExpectedHeadwayInterval]   |                                                                    5 |
      | Attribute[NumberOfStopsAway]         |                                                                    1 |
      | Attribute[PlatformTraversal]         | false                                                                |
      | Reference[OperatorRef]#Code          | "internal":"NINOXE:Company:15563880:LOC"                             |
    When I send a SIRI GetStopMonitoring request with
      | RequestorRef  | test                       |
      | MonitoringRef | NINOXE:StopPoint:SP:24:LOC |
    Then I should receive a SIRI GetStopMonitoringResponse with
      | //siri:MonitoredStopVisit[1]/siri:RecordedAtTime                                                                            |                                    2017-01-01T11:00:00.000Z | StopVisit#RecordedAt                                  |
      | //siri:MonitoredStopVisit[1]/siri:ItemIdentifier                                                                            | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3      | StopVisit#Code                                        |
      | //siri:MonitoredStopVisit[1]/siri:MonitoringRef                                                                             | NINOXE:StopPoint:SP:24:LOC                                  | StopArea#Code                                         |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:LineRef                                                      | NINOXE:Line:3:LOC                                           | Line#Code                                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DirectionRef                                                 | Aller                                                       | VehicleJourney#DirectionType                          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DataFrameRef                    | RATPDev:DataFrame::2017-01-01:LOC                           | Model#Date                                            |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef          | NINOXE:VehicleJourney:201                                   | VehicleJourney#Code                                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyPatternRef                                            | NINOXE:JourneyPattern:3_42_62:LOC                           | VehicleJourney#Reference[JourneyPatternRef]#Code      |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyPatternName                                           | TEST                                                        | VehicleJourney#Attribute[JourneyPatternName]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleMode                                                  | bus                                                         | VehicleJourney#Attribute[VehicleMode]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:PublishedLineName                                            | Ligne 3 Metro                                               | Line#Name                                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:RouteRef                                                     | RATPDev:Route::720c054714b4464d42970bda37a7edc5af8082cb:LOC | VehicleJourney#Reference[RouteRef]#Code               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DirectionName                                                | Mago-Cime OMNI                                              | VehicleJourney#Attribute[DirectionName]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OperatorRef                                                  | NINOXE:Company:15563880:LOC                                 | StopVisit#Reference[OperatorRef]#Code                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:ProductCategoryRef                                           |                                                           0 | VehicleJourney#Attribute[ProductCategoryRef]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:ServiceFeatureRef                                            | bus scolaire                                                | VehicleJourney#Attribute[ServiceFeatureRef]           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleFeatureRef                                            | longTrain                                                   | VehicleJourney#Attribute[VehicleFeatureRef]           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginRef                                                    | NINOXE:StopPoint:SP:42:LOC                                  | VehicleJourney#Reference[OriginRef]#Code              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginName                                                   | Magicien Noir                                               | VehicleJourney#Attribute[OriginName]                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Via/siri:PlaceName                                           | Saint Bénédicte                                             | VehicleJourney#Attribute[ViaPlaceName]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Via/siri:PlaceRef                                            | NINOXE:StopPoint:SP:256:LOC                                 | VehicleJourney#Reference[PlaceRef]#Code               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationRef                                               | NINOXE:StopPoint:SP:62:LOC                                  | VehicleJourney#Reference[DestinationRef]#Code         |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationName                                              | Cimetière des Sauvages                                      | VehicleJourney#Attribute[DestinationName]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                           | Magicien Noir - Cimetière (OMNI)                            | VehicleJourney#Name                                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyNote                                                  | Note de test                                                | VehicleJourney#Attribute[JourneyNote]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:HeadwayService                                               | false                                                       | VehicleJourney#Attribute[HeadwayService]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginAimedDepartureTime                                     |                                    2016-09-22T07:54:52.977Z | VehicleJourney#Attribute[OriginAimedDepartureTime]    |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationAimedArrivalTime                                  |                                    2016-09-22T09:54:52.977Z | VehicleJourney#Attribute[DestinationAimedArrivalTime] |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FirstOrLastJourney                                           | first                                                       | VehicleJourney#Attribute[FirstOrLastJourney]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Monitored                                                    | true                                                        | VehicleJourney#Attribute[Monitored]                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoringError                                              | false                                                       | VehicleJourney#Attribute[MonitoringError]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Occupancy                                                    | manySeatsAvailable                                          | VehicleJourney#Occupancy                              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Delay                                                        |                                                          30 | VehicleJourney#Attribute[Delay]                       |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Bearing                                                      | N                                                           | VehicleJourney#Attribute[Bearing]                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:InPanic                                                      | false                                                       | VehicleJourney#Attribute[InPanic]                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:InCongestion                                                 | false                                                       | VehicleJourney#Attribute[InCongestion]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:TrainNumber/siri:TrainNumberRef                              |                                                       12345 | VehicleJourney#Attribute[TrainNumberRef]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:SituationRef                                                 |                                                     1234556 | TODO                                                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:StopPointRef                              | NINOXE:StopPoint:SP:24:LOC                                  | StopArea#Code                                         |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:Order                                     |                                                           4 | StopVisit#PassageOrder                                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:VehicleAtStop                             | true                                                        | StopVisit#VehicleAtStop                               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:PlatformTraversal                         | false                                                       | StopVisit#Attribute[PlatformTraversal]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DestinationDisplay                        | Balard Terminus                                             | StopVisit#Attribute[DestinationDisplay]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedArrivalTime                          |                                    2017-01-01T13:00:00.000Z | StopVisit#Schedule[aimed]#Arrival                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ActualArrivalTime                         |                                    2017-01-01T13:00:00.000Z | StopVisit#Schedule[actual]#Arrival                    |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedArrivalTime                       |                                    2017-01-01T13:00:00.000Z | StopVisit#Schedule[expected]#Arrival                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalStatus                             | onTime                                                      | StopVisit#ArrivalStatus                               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalProximyTest                        | A l'approche                                                | StopVisit#Attribute[ArrivalProximyTest]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalPlatformName                       | Platform Name                                               | StopVisit#Attribute[ArrivalPlatformName]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalStopAssignment/siri:ActualQuayName | Quay Name                                                   | StopVisit#Attribute[ActualQuayName]                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedDepartureTime                        |                                    2017-01-01T13:02:00.000Z | StopVisit#Schedule[aimed]#Departure                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ActualDepartureTime                       |                                    2017-01-01T13:02:00.000Z | StopVisit#Schedule[actual]#Departure                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedDepartureTime                     |                                    2017-01-01T13:02:00.000Z | StopVisit#Schedule[expected]#Departure                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DepartureStatus                           | onTime                                                      | StopVisit#DepartureStatus                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DeparturePlatformName                     | Departure Platform Name                                     | StopVisit#Attribute[DeparturePlatformName]            |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DepartureBoardingActivity                 | boarding                                                    | StopVisit#Attribute[DepartureBoardingActivity]        |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedHeadwayInterval                      |                                                           5 | StopVisit#Attribute[AimedHeadwayInterval]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedHeadwayInterval                   |                                                           5 | StopVisit#Attribute[ExpectedHeadwayInterval]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DistanceFromStop                          |                                                         800 | StopVisit#Attribute[DistanceFromStop]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:NumberOfStopsAway                         |                                                           1 | StopVisit#Attribute[NumberOfStopsAway]                |

  Scenario: Handle a SIRI StopMonitoring request by returning all required attributes with the rewrite JourneyPatternRef setting
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                      | test     |
      | remote_code_space                     | internal |
      | broadcast.rewrite_journey_pattern_ref | true     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Destination                              |
      | Codes     | "internal": "NINOXE:StopPoint:SP:62:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Origin                                   |
      | Codes     | "internal": "NINOXE:StopPoint:SP:42:LOC" |
      | Monitored | true                                     |
    And a StopArea exists with the following attributes:
      | Name      | Via                                       |
      | Codes     | "internal": "NINOXE:StopPoint:SP:256:LOC" |
      | Monitored | true                                      |
    And a Line exists with the following attributes:
      | Codes        | "internal": "NINOXE:Line:3:LOC"           |
      | Name         | Ligne 3 Metro                             |
      | OperationRef | "internal": "NINOXE:Company:15563880:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Codes                                  | "internal": "NINOXE:VehicleJourney:201"         |
      | Name                                   | Magicien Noir - Cimetière (OMNI)                |
      | LineId                                 |               6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | Monitored                              | true                                            |
      | Attribute[Bearing]                     | N                                               |
      | Attribute[Delay]                       |                                              30 |
      | DestinationName                        | Cimetière des Sauvages                          |
      | Attribute[DirectionName]               | Mago-Cime OMNI                                  |
      | DirectionType                          | Aller                                           |
      | Attribute[FirstOrLastJourney]          | first                                           |
      | Attribute[HeadwayService]              | false                                           |
      | Attribute[InCongestion]                | false                                           |
      | Attribute[InPanic]                     | false                                           |
      | Attribute[JourneyNote]                 | Note de test                                    |
      | Attribute[JourneyPatternName]          | TEST                                            |
      | Attribute[MonitoringError]             | false                                           |
      | Attribute[Occupancy]                   | seatsAvailable                                  |
      | Attribute[OriginAimedDepartureTime]    |                        2016-09-22T07:54:52.977Z |
      | Attribute[DestinationAimedArrivalTime] |                        2016-09-22T09:54:52.977Z |
      | OriginName                             | Magicien Noir                                   |
      | Attribute[ProductCategoryRef]          |                                               0 |
      | Attribute[ServiceFeatureRef]           | bus scolaire                                    |
      | Attribute[TrainNumberRef]              |                                           12345 |
      | Attribute[VehicleFeatureRef]           | longTrain                                       |
      | Attribute[VehicleMode]                 | bus                                             |
      | Attribute[ViaPlaceName]                | Saint Bénédicte                                 |
      | Reference[DestinationRef]#Code         | "internal": "NINOXE:StopPoint:SP:62:LOC"        |
      | Reference[JourneyPatternRef]#Code      | "internal": "NINOXE:JourneyPattern:3_42_62:LOC" |
      | Reference[OriginRef]#Code              | "internal": "NINOXE:StopPoint:SP:42:LOC"        |
      | Reference[RouteRef]#Code               | "internal": "NINOXE:Route:66:LOC"               |
      | Reference[PlaceRef]#Code               | "internal": "NINOXE:StopPoint:SP:256:LOC"       |
    And a StopVisit exists with the following attributes:
      | ArrivalStatus                        | onTime                                                               |
      | DepartureStatus                      | onTime                                                               |
      | Codes                                | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                         |                                                                    4 |
      | RecordedAt                           |                                             2017-01-01T11:00:00.000Z |
      | Schedule[actual]#Arrival             |                                             2017-01-01T13:00:00.000Z |
      | Schedule[actual]#Departure           |                                             2017-01-01T13:02:00.000Z |
      | Schedule[aimed]#Arrival              |                                             2017-01-01T13:00:00.000Z |
      | Schedule[aimed]#Departure            |                                             2017-01-01T13:02:00.000Z |
      | Schedule[expected]#Arrival           |                                             2017-01-01T13:00:00.000Z |
      | Schedule[expected]#Departure         |                                             2017-01-01T13:02:00.000Z |
      | StopAreaId                           |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId                     |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop                        | true                                                                 |
      | Attribute[AimedHeadwayInterval]      |                                                                    5 |
      | Attribute[ActualQuayName]            | Quay Name                                                            |
      | Attribute[ArrivalPlatformName]       | Platform Name                                                        |
      | Attribute[ArrivalProximyTest]        | A l'approche                                                         |
      | Attribute[DepartureBoardingActivity] | boarding                                                             |
      | Attribute[DeparturePlatformName]     | Departure Platform Name                                              |
      | Attribute[DestinationDisplay]        | Balard Terminus                                                      |
      | Attribute[DistanceFromStop]          |                                                                  800 |
      | Attribute[ExpectedHeadwayInterval]   |                                                                    5 |
      | Attribute[NumberOfStopsAway]         |                                                                    1 |
      | Attribute[PlatformTraversal]         | false                                                                |
      | Reference[OperatorRef]#Code          | "internal":"NINOXE:Company:15563880:LOC"                             |
    When I send a SIRI GetStopMonitoring request with
      | RequestorRef  | test                       |
      | MonitoringRef | NINOXE:StopPoint:SP:24:LOC |
    Then I should receive a SIRI GetStopMonitoringResponse with
      | //siri:MonitoredStopVisit[1]/siri:RecordedAtTime                                                                            |                                             2017-01-01T11:00:00.000Z | StopVisit#RecordedAt                                  |
      | //siri:MonitoredStopVisit[1]/siri:ItemIdentifier                                                                            | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3               | StopVisit#Code                                        |
      | //siri:MonitoredStopVisit[1]/siri:MonitoringRef                                                                             | NINOXE:StopPoint:SP:24:LOC                                           | StopArea#Code                                         |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:LineRef                                                      | NINOXE:Line:3:LOC                                                    | Line#Code                                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DirectionRef                                                 | Aller                                                                | VehicleJourney#DirectionType                          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DataFrameRef                    | RATPDev:DataFrame::2017-01-01:LOC                                    | Model#Date                                            |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef          | NINOXE:VehicleJourney:201                                            | VehicleJourney#Code                                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyPatternRef                                            | RATPDev:JourneyPattern::775b650b33aa71eaa01222ccf88a68ce23b58eff:LOC | VehicleJourney#Reference[JourneyPatternRef]#Code      |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyPatternName                                           | TEST                                                                 | VehicleJourney#Attribute[JourneyPatternName]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleMode                                                  | bus                                                                  | VehicleJourney#Attribute[VehicleMode]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:PublishedLineName                                            | Ligne 3 Metro                                                        | Line#Name                                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:RouteRef                                                     | RATPDev:Route::720c054714b4464d42970bda37a7edc5af8082cb:LOC          | VehicleJourney#Reference[RouteRef]#Code               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DirectionName                                                | Mago-Cime OMNI                                                       | VehicleJourney#Attribute[DirectionName]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OperatorRef                                                  | NINOXE:Company:15563880:LOC                                          | StopVisit#Reference[OperatorRef]#Code                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:ProductCategoryRef                                           |                                                                    0 | VehicleJourney#Attribute[ProductCategoryRef]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:ServiceFeatureRef                                            | bus scolaire                                                         | VehicleJourney#Attribute[ServiceFeatureRef]           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleFeatureRef                                            | longTrain                                                            | VehicleJourney#Attribute[VehicleFeatureRef]           |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginRef                                                    | NINOXE:StopPoint:SP:42:LOC                                           | VehicleJourney#Reference[OriginRef]#Code              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginName                                                   | Magicien Noir                                                        | VehicleJourney#Attribute[OriginName]                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Via/siri:PlaceName                                           | Saint Bénédicte                                                      | VehicleJourney#Attribute[ViaPlaceName]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Via/siri:PlaceRef                                            | NINOXE:StopPoint:SP:256:LOC                                          | VehicleJourney#Reference[PlaceRef]#Code               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationRef                                               | NINOXE:StopPoint:SP:62:LOC                                           | VehicleJourney#Reference[DestinationRef]#Code         |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationName                                              | Cimetière des Sauvages                                               | VehicleJourney#Attribute[DestinationName]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                           | Magicien Noir - Cimetière (OMNI)                                     | VehicleJourney#Name                                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:JourneyNote                                                  | Note de test                                                         | VehicleJourney#Attribute[JourneyNote]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:HeadwayService                                               | false                                                                | VehicleJourney#Attribute[HeadwayService]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:OriginAimedDepartureTime                                     |                                             2016-09-22T07:54:52.977Z | VehicleJourney#Attribute[OriginAimedDepartureTime]    |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DestinationAimedArrivalTime                                  |                                             2016-09-22T09:54:52.977Z | VehicleJourney#Attribute[DestinationAimedArrivalTime] |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FirstOrLastJourney                                           | first                                                                | VehicleJourney#Attribute[FirstOrLastJourney]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Monitored                                                    | true                                                                 | VehicleJourney#Attribute[Monitored]                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoringError                                              | false                                                                | VehicleJourney#Attribute[MonitoringError]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Delay                                                        |                                                                   30 | VehicleJourney#Attribute[Delay]                       |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:Bearing                                                      | N                                                                    | VehicleJourney#Attribute[Bearing]                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:InPanic                                                      | false                                                                | VehicleJourney#Attribute[InPanic]                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:InCongestion                                                 | false                                                                | VehicleJourney#Attribute[InCongestion]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:TrainNumber/siri:TrainNumberRef                              |                                                                12345 | VehicleJourney#Attribute[TrainNumberRef]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:SituationRef                                                 |                                                              1234556 | TODO                                                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:StopPointRef                              | NINOXE:StopPoint:SP:24:LOC                                           | StopArea#Code                                         |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:Order                                     |                                                                    4 | StopVisit#PassageOrder                                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:VehicleAtStop                             | true                                                                 | StopVisit#VehicleAtStop                               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:PlatformTraversal                         | false                                                                | StopVisit#Attribute[PlatformTraversal]                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DestinationDisplay                        | Balard Terminus                                                      | StopVisit#Attribute[DestinationDisplay]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedArrivalTime                          |                                             2017-01-01T13:00:00.000Z | StopVisit#Schedule[aimed]#Arrival                     |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ActualArrivalTime                         |                                             2017-01-01T13:00:00.000Z | StopVisit#Schedule[actual]#Arrival                    |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedArrivalTime                       |                                             2017-01-01T13:00:00.000Z | StopVisit#Schedule[expected]#Arrival                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalStatus                             | onTime                                                               | StopVisit#ArrivalStatus                               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalProximyTest                        | A l'approche                                                         | StopVisit#Attribute[ArrivalProximyTest]               |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalPlatformName                       | Platform Name                                                        | StopVisit#Attribute[ArrivalPlatformName]              |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ArrivalStopAssignment/siri:ActualQuayName | Quay Name                                                            | StopVisit#Attribute[ActualQuayName]                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedDepartureTime                        |                                             2017-01-01T13:02:00.000Z | StopVisit#Schedule[aimed]#Departure                   |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ActualDepartureTime                       |                                             2017-01-01T13:02:00.000Z | StopVisit#Schedule[actual]#Departure                  |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedDepartureTime                     |                                             2017-01-01T13:02:00.000Z | StopVisit#Schedule[expected]#Departure                |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DepartureStatus                           | onTime                                                               | StopVisit#DepartureStatus                             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DeparturePlatformName                     | Departure Platform Name                                              | StopVisit#Attribute[DeparturePlatformName]            |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DepartureBoardingActivity                 | boarding                                                             | StopVisit#Attribute[DepartureBoardingActivity]        |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:AimedHeadwayInterval                      |                                                                    5 | StopVisit#Attribute[AimedHeadwayInterval]             |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:ExpectedHeadwayInterval                   |                                                                    5 | StopVisit#Attribute[ExpectedHeadwayInterval]          |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:DistanceFromStop                          |                                                                  800 | StopVisit#Attribute[DistanceFromStop]                 |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:MonitoredCall/siri:NumberOfStopsAway                         |                                                                    1 | StopVisit#Attribute[NumberOfStopsAway]                |

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
      | remote_url        | http://localhost:8090 |
      | remote_credential | source                |
      | remote_code_space | internal              |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name            | arrêt 1                                  |
      | Codes           | "internal": "NINOXE:StopPoint:SP:24:LOC" |
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
      | remote_url        | http://localhost:8090 |
      | remote_credential | source                |
      | remote_code_space | internal              |
    And a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name            | arrêt 1                                  |
      | Codes           | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | MonitoredAlways | false                                    |
      | CollectedAlways | false                                    |
    When I send a SIRI GetStopMonitoring request with
      | RequestTimestamp  |   2017-01-01T07:54:00.977Z |
      | RequestorRef      | test                       |
      | MessageIdentifier | StopMonitoring:Test:0      |
      | StartTime         |   2017-01-01T07:54:00.977Z |
      | MonitoringRef     | NINOXE:StopPoint:SP:24:LOC |
      | StopVisitTypes    | all                        |
    And a minute has passed
    Then the SIRI server should have received a GetStopMonitoring request with:
      | //siri:MonitoringRef | NINOXE:StopPoint:SP:24:LOC |
    # And the StopArea "arrêt 1" should have the following attributes:
    #   | CollectedUntil | ~ 07h54 |

  Scenario: 2481 - Handle a SIRI StopMonitoring request on a unknown StopArea
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
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
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name      | Parent               |
      | Codes     | "internal": "parent" |
      | Monitored | true                 |
    And a StopArea exists with the following attributes:
      | Name      | Child                             |
      | Codes     | "internal": "child"               |
      | ParentId  | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Monitored | true                              |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:4:LOC" |
      | Name  | Ligne 4 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | Codes     | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    |       6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Monitored | true                                    |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 202                             |
      | Codes     | "internal": "NINOXE:VehicleJourney:202" |
      | LineId    |       6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored | true                                    |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival    |                                             2017-01-01T13:00:00.000Z |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:25:LOC-3" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop               | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival    |                                             2017-01-01T14:00:00.000Z |
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
      | remote_code_space          | internal                   |
      | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name     | Test 1                                   |
      | Codes    | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | ParentId |        6ba7b814-9dad-11d1-4-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | Name  | Test 2                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:25:LOC" |
    When a minute has passed
    Then the SIRI server should not have received a GetStopMonitoring request

  Scenario: Perform a SIRI StopMonitoring request when StopArea has a parent with CollectChildren
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """

      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                 | http://localhost:8090      |
      | remote_credential          | test                       |
      | remote_code_space          | internal                   |
      | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name     | Test 1                                   |
      | Codes    | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | ParentId |        6ba7b814-9dad-11d1-5-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | Name            | Test 2                                   |
      | Codes           | "internal": "NINOXE:StopPoint:SP:25:LOC" |
      | CollectChildren | true                                     |
    When a minute has passed
    Then the SIRI server should have received 1 GetStopMonitoring request

  Scenario: Handle a SIRI StopMonitoring request with Operator
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | external |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "external": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "external": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | Codes     | "external": "NINOXE:VehicleJourney:201" |
      | LineId    |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored | true                                    |
    And a StopVisit exists with the following attributes:
      | Codes                       | "external": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop               | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "internalOperator"                                       |
      | Schedule[actual]#Arrival    |                                             2017-01-01T13:00:00.000Z |
    And an Operator exists with the following attributes:
      | Name  | Operator                                                    |
      | Codes | "internal":"internalOperator","external":"externalOperator" |
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
      | local_credential  | test     |
      | remote_code_space | external |
    And a StopArea exists with the following attributes:
      | Name             | Test                                     |
      | Codes            | "external": "NINOXE:StopPoint:SP:24:LOC" |
      | Origin[partner1] | false                                    |
      | Monitored        | false                                    |
    And a Line exists with the following attributes:
      | Codes | "external": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | Codes     | "external": "NINOXE:VehicleJourney:201" |
      | LineId    |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored | true                                    |
    And a StopVisit exists with the following attributes:
      | Codes                       | "external": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop               | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "operator"                                               |
      | Schedule[aimed]#Arrival     |                                             2017-01-01T13:00:00.000Z |
      | Schedule[actual]#Arrival    |                                             2017-01-01T13:00:05.000Z |
    And an Operator exists with the following attributes:
      | Name  | Operator                                                    |
      | Codes | "internal":"internalOperator","external":"externalOperator" |
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

  @ARA-1044
  Scenario: Handle a SIRI StopMonitoring request with unmatching code kind
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test  |
      | remote_code_space | wrong |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | Codes     | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored | true                                    |
    And a StopVisit exists with the following attributes:
      | Codes                       | "other": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                |                                                                 4 |
      | StopAreaId                  |                                 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |                                 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop               | true                                                              |
      | Reference[OperatorRef]#Code | "other": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival    |                                          2017-01-01T13:00:00.000Z |
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
  Scenario: Handle a SIRI StopMonitoring request with global setting vehicle_journey_remote_code_space
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                  | test     |
      | remote_code_space                 | internal |
      | vehicle_journey_remote_code_space | other    |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                           |
      | Codes     | "other": "NINOXE:VehicleJourney:201" |
      | LineId    |    6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored | true                                 |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop               | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival    |                                             2017-01-01T13:00:00.000Z |
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
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef | NINOXE:VehicleJourney:201 | VehicleJourney#Code |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                  | Passage 32                | VehicleJourney#Name |

  @ARA-1044
  Scenario: Handle a SIRI StopMonitoring request with connector setting siri-stop-monitoring-request-broadcaster.vehicle_journey_remote_code_space
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                                                           | test     |
      | remote_code_space                                                          | internal |
      | siri-stop-monitoring-request-broadcaster.vehicle_journey_remote_code_space | other    |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                           |
      | Codes     | "other": "NINOXE:VehicleJourney:201" |
      | LineId    |    6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored | true                                 |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop               | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival    |                                             2017-01-01T13:00:00.000Z |
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
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef | NINOXE:VehicleJourney:201 | VehicleJourney#Code |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                  | Passage 32                | VehicleJourney#Name |

  @ARA-1044
  Scenario: Handle a SIRI StopMonitoring request with multiple connector setting siri-stop-monitoring-request-broadcaster.vehicle_journey_remote_code_space
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                                                           | test          |
      | remote_code_space                                                          | internal      |
      | siri-stop-monitoring-request-broadcaster.vehicle_journey_remote_code_space | other, other2 |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                           |
      | Codes     | "other": "NINOXE:VehicleJourney:201" |
      | LineId    |    6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored | true                                 |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop               | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival    |                                             2017-01-01T13:00:00.000Z |
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
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef | NINOXE:VehicleJourney:201 | VehicleJourney#Code |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                  | Passage 32                | VehicleJourney#Name |

  @ARA-1044
  Scenario: Handle a SIRI StopMonitoring request with fallback on generic connector settings remote_code_space
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                                           | test     |
      | remote_code_space                                          | internal |
      | siri-stop-monitoring-request-broadcaster.remote_code_space | other    |
    And a StopArea exists with the following attributes:
      | Name      | Test                                  |
      | Codes     | "other": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                  |
    And a Line exists with the following attributes:
      | Codes | "other": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                           |
      | Codes     | "other": "NINOXE:VehicleJourney:201" |
      | LineId    |    6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored | true                                 |
    And a StopVisit exists with the following attributes:
      | Codes                       | "other": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                |                                                                 4 |
      | StopAreaId                  |                                 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |                                 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop               | true                                                              |
      | Reference[OperatorRef]#Code | "other": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival    |                                          2017-01-01T13:00:00.000Z |
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
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:FramedVehicleJourneyRef/siri:DatedVehicleJourneyRef | NINOXE:VehicleJourney:201 | VehicleJourney#Code |
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:VehicleJourneyName                                  | Passage 32                | VehicleJourney#Name |

  @ARA-1101
  Scenario: Handle a SIRI StopMonitoring request with partner setting siri.direction_type should broadcast the DirectionRef with setting value
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential    | test                              |
      | remote_code_space   | internal                          |
      | siri.direction_type | ch:1:Direction:R,ch:1:Direction:H |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes        | "internal": "NINOXE:Line:3:LOC"           |
      | Name         | Ligne 3 Metro                             |
      | OperationRef | "internal": "NINOXE:Company:15563880:LOC" |
    And a VehicleJourney exists with the following attributes:
      | Codes                    | "internal": "NINOXE:VehicleJourney:201" |
      | Name                     | Magicien Noir - Cimetière (OMNI)        |
      | LineId                   |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored                | true                                    |
      | Occupancy                | manySeatsAvailable                      |
      | Attribute[DirectionName] | Mago-Cime OMNI                          |
      | DirectionType            | inbound                                 |
    And a StopVisit exists with the following attributes:
      | ArrivalStatus                | onTime                                                               |
      | Codes                        | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                 |                                                                    4 |
      | RecordedAt                   |                                             2017-01-01T11:00:00.000Z |
      | Schedule[expected]#Departure |                                             2017-01-01T13:02:00.000Z |
      | StopAreaId                   |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId             |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop                | true                                                                 |
      | Reference[OperatorRef]#Code  | "internal":"NINOXE:Company:15563880:LOC"                             |
    When I send a SIRI GetStopMonitoring request with
      | RequestorRef  | test                       |
      | MonitoringRef | NINOXE:StopPoint:SP:24:LOC |
    Then I should receive a SIRI GetStopMonitoringResponse with
      | //siri:MonitoredStopVisit[1]/siri:MonitoredVehicleJourney/siri:DirectionRef | ch:1:Direction:R | VehicleJourney#DirectionType |

  @ARA-1306
  Scenario: StopMonitoring request collect should send GetStopMonitoring request to partner
    Given a SIRI server on "http://localhost:8090"
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
    Then the SIRI server should have received 1 GetStopMonitoring request

  @ARA-1306
  Scenario: StopMonitoring request collect and partner CheckStatus is unavailable should not send GetStopMonitoring request to partner
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                 | http://localhost:8090      |
      | remote_credential          | test                       |
      | remote_code_space          | internal                   |
      | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name  | Test 1                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    When a minute has passed
    And the SIRI server should not have received a GetStopMonitoring request

  @ARA-1306
  Scenario: StopMonitoring request collect and partner CheckStauts is unavailable should send GetStopMonitoring request to partner whith setting collect.persistent
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                 | http://localhost:8090      |
      | remote_credential          | test                       |
      | remote_code_space          | internal                   |
      | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
      | collect.persistent         | true                       |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name  | Test 1                                   |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    When a minute has passed
    Then the SIRI server should have received 1 GetStopMonitoring request

  @ARA-1324
  Scenario: Handle a SIRI StopMonitoring request with the StopVisit beeing the next Stop Visit of a Vehicle should boradcast vehicle information
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
      # 6ba7b814-9dad-11d1-2-00c04fd430c8
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
      # 6ba7b814-9dad-11d1-3-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | Codes     | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored | true                                    |
      # 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                |                                                                    4 |
      | StopAreaId                  |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop               | true                                                                 |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival    |                                             2017-01-01T13:00:00.000Z |
      # 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "RLA290"              |
      | LineId           | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | StopAreaId       | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Bearing          |                            121.55 |
      | Latitude         |                             55.55 |
      | Longitude        |                         111.11111 |
    And I see ara vehicles
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
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:VehicleLocation>
                <siri:Longitude>111.11111</siri:Longitude>
                <siri:Latitude>55.55</siri:Latitude>
              </siri:VehicleLocation>
              <siri:Bearing>121.55</siri:Bearing>
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

  @ARA-1298
  Scenario: Handle a SIRI StopMonitoring request with Partner remote_code_space changed
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
      # "Id":"6ba7b814-9dad-11d1-2-00c04fd430c8"
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "external": "NINOXE:StopPoint:SP:25:LOC" |
      | Monitored | true                                     |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes | "external": "Test:Line:A:BUS:LOC" |
      | Name  | Ligne A Bus                       |
      # "Id":"6ba7b814-9dad-11d1-5-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:201" |
      | LineId          |       6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 33 external                     |
      | Codes           | "external": "NINOXE:VehicleJourney:202" |
      | LineId          |       6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And I see ara lines
    And I see ara stop_areas
    And I see ara vehicle_journeys
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                  |                                                                    4 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop                 | true                                                                 |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival      |                                             2017-01-01T13:00:00.000Z |
      | Attribute[DestinationDisplay] | Cergy le haut & Arret <RER>                                          |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"
    And a StopVisit exists with the following attributes:
      | Codes                         | "external": "NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:25:LOC-5" |
      | PassageOrder                  |                                                                    5 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop                 | true                                                                 |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival      |                                             2017-01-01T13:00:00.000Z |
      | Attribute[DestinationDisplay] | Montfermeil                                                          |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"
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
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
              <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:DestinationDisplay>Cergy le haut &amp; Arret &lt;RER&gt;</siri:DestinationDisplay>
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
    And the Partner "test" is updated with the following settings:
      | local_credential  | test     |
      | remote_code_space | external |
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
        <ns2:MonitoringRef>NINOXE:StopPoint:SP:25:LOC</ns2:MonitoringRef>
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
              <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-b-00c04fd430c8</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:25:LOC</siri:MonitoringRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:25:LOC-5</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:25:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>Test:Line:A:BUS:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>2017-01-01</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne A Bus</siri:PublishedLineName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
                    <siri:VehicleJourneyName>Passage 33 external</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:StopPointRef>
                      <siri:Order>5</siri:Order>
                      <siri:StopPointName>Test</siri:StopPointName>
                      <siri:VehicleAtStop>true</siri:VehicleAtStop>
                      <siri:DestinationDisplay>Montfermeil</siri:DestinationDisplay>
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

  @ARA-1363
  Scenario: Handle a SIRI StopMonitoring request using the generator setting reference_vehicle_journey_identifier
    # Setting a Partner without default generators
    Given a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential                                | test                             |
      | remote_code_space                               | internal                         |
      | generators.reference_vehicle_journey_identifier | ch:1:ServiceJourney:87_TAC:%{id} |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                                                      |
      | Codes           | "_default": "6ba7b814", "external": "NINOXE:VehicleJourney:201" |
      | LineId          |                               6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                                                            |
      | DestinationName | La fin. <TER>                                                   |
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                  |                                                                    4 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop                 | true                                                                 |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival      |                                             2017-01-01T13:00:00.000Z |
      | Attribute[DestinationDisplay] | Cergy le haut & Arret <RER>                                          |
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
        <siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:ResponseMessageIdentifier>
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
                <siri:DataFrameRef>2017-01-01</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>ch:1:ServiceJourney:87_TAC:6ba7b814</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
              <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:DestinationDisplay>Cergy le haut &amp; Arret &lt;RER&gt;</siri:DestinationDisplay>
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
    And an audit event should exist with these attributes:
      | Protocol        | siri                                    |
      | Direction       | received                                |
      | Status          | OK                                      |
      | Type            | StopMonitoringRequest                   |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"]          |
      | VehicleJourneys | ["ch:1:ServiceJourney:87_TAC:6ba7b814"] |
      | Lines           | ["NINOXE:Line:3:LOC"]                   |

  @ARA-1363
  Scenario: Handle a SIRI StopMonitoring request using the default generator should send DatedVehicleJourneyRef according to default setting
    # Setting a Partner with default generators
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                                                      |
      | Codes           | "_default": "6ba7b814", "external": "NINOXE:VehicleJourney:201" |
      | LineId          |                               6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                                                            |
      | DestinationName | La fin. <TER>                                                   |
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                  |                                                                    4 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | VehicleAtStop                 | true                                                                 |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival      |                                             2017-01-01T13:00:00.000Z |
      | Attribute[DestinationDisplay] | Cergy le haut & Arret <RER>                                          |
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
                <siri:DatedVehicleJourneyRef>RATPDev:VehicleJourney::6ba7b814:LOC</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
              <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:DestinationDisplay>Cergy le haut &amp; Arret &lt;RER&gt;</siri:DestinationDisplay>
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
    And an audit event should exist with these attributes:
      | Protocol        | siri                                     |
      | Direction       | received                                 |
      | Status          | OK                                       |
      | Type            | StopMonitoringRequest                    |
      | StopAreas       | ["NINOXE:StopPoint:SP:24:LOC"]           |
      | VehicleJourneys | ["RATPDev:VehicleJourney::6ba7b814:LOC"] |
      | Lines           | ["NINOXE:Line:3:LOC"]                    |

  @ARA-1493
  Scenario: Handle a SIRI StopMonitoring request with a Line Having a Referent
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | Codes | "external": "NINOXE:Line:4:LOC" |
      | ReferentId |6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Name  | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:201" |
      | LineId          |       6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                  |                                                                    4 |
      | StopAreaId                    |                                    6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId              |                                    6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | VehicleAtStop                 | true                                                                 |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[actual]#Arrival      |                                             2017-01-01T13:00:00.000Z |
      | Attribute[DestinationDisplay] | Cergy le haut & Arret <RER>                                          |
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
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
              <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>4</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:DestinationDisplay>Cergy le haut &amp; Arret &lt;RER&gt;</siri:DestinationDisplay>
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

  @ARA-1493
  Scenario: Handle a SIRI StopMonitoring request with a Referent Line family
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential      | test     |
      | remote_code_space     | internal |
      | sort_payload_for_test | true     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
      # "Id":"6ba7b814-9dad-11d1-2-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes | "internal": "referent" |
      | Name  | Ligne Ref              |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes      | "external": "line1"               |
      | ReferentId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Name       | Ligne 1                           |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes      | "external": "line2"               |
      | ReferentId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Name       | Ligne 2                           |
      # "Id":"6ba7b814-9dad-11d1-5-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:201" |
      | LineId          |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 33                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:202" |
      | LineId          |       6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 34                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:203" |
      | LineId          |       6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
      # "Id":"6ba7b814-9dad-11d1-8-00c04fd430c8"
    And I see ara lines
    And I see ara stop_areas
    And I see ara vehicle_journeys
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "stopVisit1"           |
      | PassageOrder                |                                  4 |
      | StopAreaId                  |  6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |  6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | true                               |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC" |
      | Schedule[actual]#Arrival    |           2017-01-01T13:00:00.000Z |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "stopVisit2"           |
      | PassageOrder                |                                  5 |
      | StopAreaId                  |  6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |  6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop               | true                               |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC" |
      | Schedule[actual]#Arrival    |           2017-01-01T13:00:00.000Z |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "stopVisit3"           |
      | PassageOrder                |                                  4 |
      | StopAreaId                  |  6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |  6ba7b814-9dad-11d1-8-00c04fd430c8 |
      | VehicleAtStop               | true                               |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC" |
      | Schedule[actual]#Arrival    |           2017-01-01T13:00:00.000Z |
    And I see ara stop_visits
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
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-c-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
            <siri:ItemIdentifier>stopVisit3</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>referent</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:203</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne Ref</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
              <siri:VehicleJourneyName>Passage 34</siri:VehicleJourneyName>
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
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
            <siri:ItemIdentifier>stopVisit2</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>referent</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne Ref</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
              <siri:VehicleJourneyName>Passage 33</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>5</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:ActualArrivalTime>2017-01-01T13:00:00.000Z</siri:ActualArrivalTime>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
            <siri:ItemIdentifier>stopVisit1</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>referent</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne Ref</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
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

  @ARA-1493
  Scenario: Handle a SIRI StopMonitoring request on a Line having a referent
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential      | test     |
      | remote_code_space     | internal |
      | sort_payload_for_test | true     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
      # "Id":"6ba7b814-9dad-11d1-2-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes | "internal": "referent" |
      | Name  | Ligne Ref              |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes      | "internal": "line1"               |
      | ReferentId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Name       | Ligne 1                           |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:201" |
      | LineId          |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
      # "Id":"6ba7b814-9dad-11d1-5-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 33                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:202" |
      | LineId          |       6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"
    And I see ara lines
    And I see ara stop_areas
    And I see ara vehicle_journeys
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "stopVisit1"           |
      | PassageOrder                |                                  4 |
      | StopAreaId                  |  6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |  6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | VehicleAtStop               | true                               |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC" |
      | Schedule[actual]#Arrival    |           2017-01-01T13:00:00.000Z |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "stopVisit2"           |
      | PassageOrder                |                                  5 |
      | StopAreaId                  |  6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |  6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | true                               |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC" |
      | Schedule[actual]#Arrival    |           2017-01-01T13:00:00.000Z |
    And I see ara stop_visits
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
        <ns2:LineRef>line1</ns2:LineRef>
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
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
            <siri:ItemIdentifier>stopVisit2</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>line1</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne 1</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
              <siri:VehicleJourneyName>Passage 33</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>5</siri:Order>
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

  @ARA-1493
  Scenario: Handle a SIRI StopMonitoring request with a Referent Line and LineFilter
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential      | test     |
      | remote_code_space     | internal |
      | sort_payload_for_test | true     |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | Codes     | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Monitored | true                                     |
      # "Id":"6ba7b814-9dad-11d1-2-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes | "internal": "referent" |
      | Name  | Ligne Ref              |
      # "Id":"6ba7b814-9dad-11d1-3-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes      | "external": "line1"               |
      | ReferentId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Name       | Ligne 1                           |
      # "Id":"6ba7b814-9dad-11d1-4-00c04fd430c8"
    And a Line exists with the following attributes:
      | Codes      | "external": "line2"               |
      | ReferentId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Name       | Ligne 2                           |
      # "Id":"6ba7b814-9dad-11d1-5-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:201" |
      | LineId          |       6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
      # "Id":"6ba7b814-9dad-11d1-6-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 33                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:202" |
      | LineId          |       6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
      # "Id":"6ba7b814-9dad-11d1-7-00c04fd430c8"
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 34                              |
      | Codes           | "internal": "NINOXE:VehicleJourney:203" |
      | LineId          |       6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Monitored       | true                                    |
      | DestinationName | La fin. <TER>                           |
      # "Id":"6ba7b814-9dad-11d1-8-00c04fd430c8"
    And I see ara lines
    And I see ara stop_areas
    And I see ara vehicle_journeys
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "stopVisit1"           |
      | PassageOrder                |                                  4 |
      | StopAreaId                  |  6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |  6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | VehicleAtStop               | true                               |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC" |
      | Schedule[actual]#Arrival    |           2017-01-01T13:00:00.000Z |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "stopVisit2"           |
      | PassageOrder                |                                  5 |
      | StopAreaId                  |  6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |  6ba7b814-9dad-11d1-7-00c04fd430c8 |
      | VehicleAtStop               | true                               |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC" |
      | Schedule[actual]#Arrival    |           2017-01-01T13:00:00.000Z |
    And a StopVisit exists with the following attributes:
      | Codes                       | "internal": "stopVisit3"           |
      | PassageOrder                |                                  4 |
      | StopAreaId                  |  6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId            |  6ba7b814-9dad-11d1-8-00c04fd430c8 |
      | VehicleAtStop               | true                               |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC" |
      | Schedule[actual]#Arrival    |           2017-01-01T13:00:00.000Z |
    And I see ara stop_visits
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
        <ns2:LineRef>referent</ns2:LineRef>
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
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-c-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
            <siri:ItemIdentifier>stopVisit3</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>referent</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:203</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne Ref</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
              <siri:VehicleJourneyName>Passage 34</siri:VehicleJourneyName>
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
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
            <siri:ItemIdentifier>stopVisit2</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>referent</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne Ref</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
              <siri:VehicleJourneyName>Passage 33</siri:VehicleJourneyName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                <siri:Order>5</siri:Order>
                <siri:StopPointName>Test</siri:StopPointName>
                <siri:VehicleAtStop>true</siri:VehicleAtStop>
                <siri:ActualArrivalTime>2017-01-01T13:00:00.000Z</siri:ActualArrivalTime>
              </siri:MonitoredCall>
            </siri:MonitoredVehicleJourney>
          </siri:MonitoredStopVisit>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
            <siri:ItemIdentifier>stopVisit1</siri:ItemIdentifier>
            <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>referent</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:PublishedLineName>Ligne Ref</siri:PublishedLineName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationName>La fin. &lt;TER&gt;</siri:DestinationName>
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
