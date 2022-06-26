Feature: Support GTFS-RT feeds
  Background:
    Given a Referential "test" is created

  Scenario: Provide a public GTFS-RT feed
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" without token
    Then I should receive a GTFS-RT response

  Scenario: Provide a authenticated GTFS-RT feed
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credential     | secret   |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response

  Scenario: Provide a authenticated GTFS-RT feed (multiple credentials)
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credentials    | secret1,secret2 |
      | remote_objectid_kind | internal        |
    When I send a GTFS-RT request to the Referential "test" with token "secret1"
    Then I should receive a GTFS-RT response
    When I send a GTFS-RT request to the Referential "test" with token "secret2"
    Then I should receive a GTFS-RT response

  Scenario: Forbid authorized request on GTFS-RT feed (no token)
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credential     | secret   |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" without token
    Then I should not receive a GTFS-RT but an unauthorized client error status

  Scenario: Forbid authorized request on GTFS-RT feed (wrong token)
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credential     | secret   |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" with token "wrong"
    Then I should not receive a GTFS-RT but an unauthorized client error status

  Scenario: Retrieve Vehicle Positions
    Given a Line exists with the following attributes:
      | Name      | Test               |
      | ObjectIDs | "internal": "1234" |
    And a VehicleJourney exists with the following attributes:
      | ObjectIDs | "internal": "2345"                |
      | LineId    | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | ObjectIDs | "internal": "4567" |
    And a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "3456"                |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | StopAreaId       | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
    And a Partner "test" exists with connectors [gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | local_credential     | secret   |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should contain a Vehicle Position with these attributes:
      | vehicle_id | 3456 |
      | stop_id    | 4567 |
      | trip_id    | 2345 |
      | route_id   | 1234 |

  @ARA-872
  Scenario: Retrieve Vehicle Positions with unmatching objectid kind
    Given a Line exists with the following attributes:
      | Name      | Test            |
      | ObjectIDs | "other": "1234" |
    Given a VehicleJourney exists with the following attributes:
      | ObjectIDs | "internal": "2345"                |
      | LineId    | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    Given a Vehicle exists with the following attributes:
      | ObjectIDs        | "other": "3456"                   |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a Partner "test" exists with connectors [gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | local_credential     | secret   |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should not contain Vehicle Positions

  @ARA-872
  Scenario: Retrieve Vehicle Positions with connector setting gtfs-rt-vehicle-positions-broadcaster.vehicle_remote_objectid_kind
    Given a Line exists with the following attributes:
      | Name      | Test               |
      | ObjectIDs | "internal": "1234" |
    Given a VehicleJourney exists with the following attributes:
      | ObjectIDs | "internal": "2345"                |
      | LineId    | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    Given a Vehicle exists with the following attributes:
      | ObjectIDs        | "other": "3456"                   |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a Partner "test" exists with connectors [gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | local_credential                                                   | secret   |
      | remote_objectid_kind                                               | internal |
      | gtfs-rt-vehicle-positions-broadcaster.vehicle_remote_objectid_kind | other    |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should contain a Vehicle Position with these attributes:
      | vehicle_id | 3456 |
      | trip_id    | 2345 |
      | route_id   | 1234 |

  @ARA-1044
  Scenario: Retrieve Vehicle Positions with connector setting gtfs-rt-vehicle-positions-broadcaster.vehicle_journey_remote_objectid_kind
    Given a Line exists with the following attributes:
      | Name      | Test               |
      | ObjectIDs | "internal": "1234" |
    Given a VehicleJourney exists with the following attributes:
      | ObjectIDs | "other": "2345"                   |
      | LineId    | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    Given a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "3456"                |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a Partner "test" exists with connectors [gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | local_credential                                                           | secret   |
      | remote_objectid_kind                                                       | internal |
      | gtfs-rt-vehicle-positions-broadcaster.vehicle_journey_remote_objectid_kind | other    |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should contain a Vehicle Position with these attributes:
      | vehicle_id | 3456 |
      | trip_id    | 2345 |
      | route_id   | 1234 |

  @ARA-1044
  Scenario: Retrieve Vehicle Positions with multiple setting gtfs-rt-vehicle-positions-broadcaster.vehicle_journey_remote_objectid_kind
    Given a Line exists with the following attributes:
      | Name      | Test               |
      | ObjectIDs | "internal": "1234" |
    Given a VehicleJourney exists with the following attributes:
      | ObjectIDs | "other": "2345"                   |
      | LineId    | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    Given a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "3456"                |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a Partner "test" exists with connectors [gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | local_credential                                                           | secret        |
      | remote_objectid_kind                                                       | internal      |
      | gtfs-rt-vehicle-positions-broadcaster.vehicle_journey_remote_objectid_kind | other, other2 |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should contain a Vehicle Position with these attributes:
      | vehicle_id | 3456 |
      | trip_id    | 2345 |
      | route_id   | 1234 |

  @ARA-1044
  Scenario: Retrieve Vehicle Positions with global setting vehicle_remote_objectid_kind
    Given a Line exists with the following attributes:
      | Name      | Test               |
      | ObjectIDs | "internal": "1234" |
    Given a VehicleJourney exists with the following attributes:
      | ObjectIDs | "internal": "2345"                |
      | LineId    | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    Given a Vehicle exists with the following attributes:
      | ObjectIDs        | "other": "3456"                   |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a Partner "test" exists with connectors [gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | local_credential             | secret   |
      | remote_objectid_kind         | internal |
      | vehicle_remote_objectid_kind | other    |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should contain a Vehicle Position with these attributes:
      | vehicle_id | 3456 |
      | trip_id    | 2345 |
      | route_id   | 1234 |

  @ARA-1044
  Scenario: Retrieve Vehicle Positions with fallback on generic connector settings
    Given a Line exists with the following attributes:
      | Name      | Test            |
      | ObjectIDs | "other": "1234" |
    Given a VehicleJourney exists with the following attributes:
      | ObjectIDs | "other": "2345"                   |
      | LineId    | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    Given a Vehicle exists with the following attributes:
      | ObjectIDs        | "other": "3456"                   |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a Partner "test" exists with connectors [gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | local_credential                                           | secret   |
      | remote_objectid_kind                                       | internal |
      | gtfs-rt-vehicle-positions-broadcaster.remote_objectid_kind | other    |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should contain a Vehicle Position with these attributes:
      | vehicle_id | 3456 |
      | trip_id    | 2345 |
      | route_id   | 1234 |

  @ARA-1077
  Scenario: Retrieve Vehicle Positions with OccupancyStatus
    Given a Line exists with the following attributes:
      | Name      | Test               |
      | ObjectIDs | "internal": "1234" |
    Given a VehicleJourney exists with the following attributes:
      | ObjectIDs | "internal": "2345"                |
      | LineId    | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    Given a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "3456"                |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Occupancy        | 2                                 |
    And a Partner "test" exists with connectors [gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | local_credential     | secret   |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should contain a Vehicle Position with these attributes:
      | vehicle_id       | 3456                |
      | trip_id          | 2345                |
      | route_id         | 1234                |
      | occupancy_status | FEW_SEATS_AVAILABLE |

  @ARA-1047
  Scenario: Broadcast after a VehicleMonitoring request collect
    Given a SIRI server waits GetVehicleMonitoring request on "http://localhost:8090" to respond with
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns1:GetVehicleMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>RLA_Bus</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>RLA_Bus:ResponseMessage::23833:LOC</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:VehicleMonitoringDelivery version="2.0:FR-IDF-2.4">
          <ns5:ResponseTimestamp>2021-08-02T08:50:49.660+02:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>Test:Message::1234:LOC</ns5:RequestMessageRef>
          <ns5:Status>true</ns5:Status>
          <ns5:VehicleActivity>
            <ns5:RecordedAtTime>2021-08-02T08:50:27.733+02:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>290</ns5:ItemIdentifier>
            <ns5:ValidUntilTime>2021-08-02T09:50:27.733+02:00</ns5:ValidUntilTime>
            <ns5:VehicleMonitoringRef>290</ns5:VehicleMonitoringRef>
            <ns5:ProgressBetweenStops>
              <ns5:LinkDistance>349.0</ns5:LinkDistance>
              <ns5:Percentage>70.0</ns5:Percentage>
            </ns5:ProgressBetweenStops>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>RLA_Bus:Line::05:LOC</ns5:LineRef>
              <ns5:DirectionRef>Aller</ns5:DirectionRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>RLA_Bus:DataFrame::1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>RLA_Bus:VehicleJourney::2978464:LOC</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>RLA_Bus:JourneyPattern::L05P99:LOC</ns5:JourneyPatternRef>
              <ns5:JourneyPatternName>L05P99</ns5:JourneyPatternName>
              <ns5:PublishedLineName>05</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>RLA_Bus:Operator::RLA:LOC</ns5:OperatorRef>
              <ns5:OriginRef>RLA_Bus:StopPoint:BP:DELOY0:LOC</ns5:OriginRef>
              <ns5:OriginName>Deloye / Dubouchage</ns5:OriginName>
              <ns5:DestinationRef>RLA_Bus:StopPoint:BP:RIMIE9:LOC</ns5:DestinationRef>
              <ns5:DestinationName>Rimiez Saint-George</ns5:DestinationName>
              <ns5:Monitored>false</ns5:Monitored>
              <ns5:VehicleLocation srsName="EPSG:2154">
                <ns5:Coordinates>1044593 6298716</ns5:Coordinates>
              </ns5:VehicleLocation>
              <ns5:Bearing>287.0</ns5:Bearing>
              <ns5:VehicleRef>RLA290</ns5:VehicleRef>
              <ns5:DriverRef>5753</ns5:DriverRef>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>RLA_Bus:StopPoint:BP:PASTO8:LOC</ns5:StopPointRef>
                <ns5:Order>6</ns5:Order>
                <ns5:StopPointName>Carabacel</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Rimiez Saint-George</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2021-08-02T07:38:42.000+02:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>delayed</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2021-08-02T07:38:42.000+02:00</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2021-08-02T08:50:51.000+02:00</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>delayed</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney><ns5:Extensions/></ns5:VehicleActivity>
        </ns5:VehicleMonitoringDelivery>
      </Answer><AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/></ns1:GetVehicleMonitoringResponse>
  </soap:Body>
</soap:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-vehicle-monitoring-request-collector, gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_objectid_kind  | internal              |
      | local_credential      | secret                |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Test 1                             |
      | ObjectIDs | "internal": "RLA_Bus:Line::05:LOC" |
    When a minute has passed
    And the SIRI server has received a GetVehicleMonitoring request
    And I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should contain a Vehicle Position with these attributes:
      | vehicle_id | RLA290                              |
      | trip_id    | RLA_Bus:VehicleJourney::2978464:LOC |
      | route_id   | RLA_Bus:Line::05:LOC                |
      | stop_id    | RLA_Bus:StopPoint:BP:PASTO8:LOC     |
