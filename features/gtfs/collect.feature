Feature: Collect realtime data via GTFS-RT feeds
  Background:
    Given a Referential "test" is created

  @ARA-1218
  Scenario: Collect GTFS TripUpdate (with stop_id) with stop_time_update having SKIPPED schedule_relationship
    Given a GTFS-RT server waits request on "http://localhost:8090" to respond with
      """
      header {
        gtfs_realtime_version: "2.0"
        incrementality: FULL_DATASET
        timestamp: 1630318853
      }
      entity {
        id: "trip:ORLEANS:VehicleJourney:20_R_67_13_2067_1_152701"
        trip_update {
          trip {
            trip_id: "Trip:A"
            route_id: "Line:1"
          }
          stop_time_update {
            stop_sequence: 0
            stop_id: "StopArea:A"
            arrival {
              time: 1483272000
            }
            departure {
              time: 1483272000
            }
            schedule_relationship: SKIPPED
          }
        }
      }
      """
    And a Partner "gtfs" exists with connectors [gtfs-rt-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_code_space | internal              |
    When a minute has passed
    Then one StopArea has the following attributes:
      | Codes | "internal": "StopArea:A" |
    And one StopVisit has the following attributes:
      | Codes           | "internal": "Trip:A-1" |
      | ArrivalStatus   | cancelled              |
      | DepartureStatus | cancelled              |

  @ARA-878
  Scenario: Collect GTFS TripUpdate (with stop_id)
    Given a GTFS-RT server waits request on "http://localhost:8090" to respond with
      """
      header {
        gtfs_realtime_version: "2.0"
        incrementality: FULL_DATASET
        timestamp: 1630318853
      }
      entity {
        id: "trip:ORLEANS:VehicleJourney:20_R_67_13_2067_1_152701"
        trip_update {
          trip {
            trip_id: "Trip:A"
            route_id: "Line:1"
          }
          stop_time_update {
            stop_sequence: 0
            stop_id: "StopArea:A"
            arrival {
              time: 1483272000
            }
            departure {
              time: 1483272000
            }
          }
          stop_time_update {
            stop_sequence: 1
            stop_id: "StopArea:B"
            arrival {
              time: 1483272060
            }
            departure {
              time: 1483272090
            }
          }
          stop_time_update {
            stop_sequence: 2
            stop_id: "StopArea:C"
            arrival {
              time: 1483272150
            }
            departure {
              time: 1483272150
            }
          }
        }
      }
      """
    And a Partner "gtfs" exists with connectors [gtfs-rt-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_code_space | internal              |
    When a minute has passed
    Then one StopArea has the following attributes:
      | Codes | "internal": "StopArea:A" |
    And one StopArea has the following attributes:
      | Codes | "internal": "StopArea:B" |
    And one StopArea has the following attributes:
      | Codes | "internal": "StopArea:C" |
    And one Line has the following attributes:
      | Codes | "internal": "Line:1" |
    And one VehicleJourney has the following attributes:
      | Codes  | "internal": "Trip:A"              |
      | LineId | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
    And one StopVisit has the following attributes:
      | Codes                        | "internal": "Trip:A-1"    |
      | PassageOrder                 | 1                         |
      | Schedule[expected]#Departure | 2017-01-01T13:00:00+01:00 |
      | Schedule[expected]#Arrival   | 2017-01-01T13:00:00+01:00 |
    And one StopVisit has the following attributes:
      | Codes                        | "internal": "Trip:A-2"    |
      | PassageOrder                 | 2                         |
      | Schedule[expected]#Arrival   | 2017-01-01T13:01:00+01:00 |
      | Schedule[expected]#Departure | 2017-01-01T13:01:30+01:00 |
    And one StopVisit has the following attributes:
      | Codes                        | "internal": "Trip:A-3"    |
      | PassageOrder                 | 3                         |
      | Schedule[expected]#Arrival   | 2017-01-01T13:02:30+01:00 |
      | Schedule[expected]#Departure | 2017-01-01T13:02:30+01:00 |

  @ARA-878
  Scenario: Collect GTFS TripUpdate (without stop_id)
    Given a GTFS-RT server waits request on "http://localhost:8090" to respond with
      """
      header {
        gtfs_realtime_version: "2.0"
        incrementality: FULL_DATASET
        timestamp: 1630318853
      }
      entity {
        id: "trip:ORLEANS:VehicleJourney:20_R_67_13_2067_1_152701"
        trip_update {
          trip {
            trip_id: "Trip:A"
            route_id: "Line:1"
          }
          stop_time_update {
            stop_sequence: 0
            arrival {
              time: 1483272000
            }
            departure {
              time: 1483272000
            }
          }
          stop_time_update {
            stop_sequence: 1
            arrival {
              time: 1483272060
            }
            departure {
              time: 1483272090
            }
          }
          stop_time_update {
            stop_sequence: 2
            arrival {
              time: 1483272150
            }
            departure {
              time: 1483272150
            }
          }
        }
      }
      """
    And a StopArea exists with the following attributes:
      | Codes | "internal": "StopArea:A" |
    And a StopArea exists with the following attributes:
      | Codes | "internal": "StopArea:B" |
    And a StopArea exists with the following attributes:
      | Codes | "internal": "StopArea:C" |
    And a Line exists with the following attributes:
      | Codes | "internal": "Line:1" |
    And a VehicleJourney exists with the following attributes:
      | Codes  | "internal": "Trip:A"              |
      | LineId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
    And a StopVisit exists with the following attributes:
      | Codes                      | "internal": "Trip:A-1"            |
      | PassageOrder               | 1                                 |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | StopAreaId                 | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
      | Schedule[expected]#Arrival | 2017-01-01T12:50:00+01:00         |
    And a StopVisit exists with the following attributes:
      | Codes                      | "internal": "Trip:A-2"            |
      | PassageOrder               | 2                                 |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | StopAreaId                 | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Schedule[expected]#Arrival | 2017-01-01T12:55:00+01:00         |
    And a StopVisit exists with the following attributes:
      | Codes                      | "internal": "Trip:A-3"            |
      | PassageOrder               | 3                                 |
      | VehicleJourneyId           | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | StopAreaId                 | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Schedule[expected]#Arrival | 2017-01-01T13:00:00+01:00         |
    And a Partner "gtfs" exists with connectors [gtfs-rt-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_code_space | internal              |
    When a minute has passed
    And one StopVisit has the following attributes:
      | Codes                        | "internal": "Trip:A-1"            |
      | PassageOrder                 | 1                                 |
      | Schedule[expected]#Departure | 2017-01-01T13:00:00+01:00         |
      | Schedule[expected]#Arrival   | 2017-01-01T13:00:00+01:00         |
      | StopAreaId                   | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    And one StopVisit has the following attributes:
      | Codes                        | "internal": "Trip:A-2"            |
      | PassageOrder                 | 2                                 |
      | Schedule[expected]#Arrival   | 2017-01-01T13:01:00+01:00         |
      | Schedule[expected]#Departure | 2017-01-01T13:01:30+01:00         |
      | StopAreaId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And one StopVisit has the following attributes:
      | Codes                        | "internal": "Trip:A-3"            |
      | PassageOrder                 | 3                                 |
      | Schedule[expected]#Arrival   | 2017-01-01T13:02:30+01:00         |
      | Schedule[expected]#Departure | 2017-01-01T13:02:30+01:00         |
      | StopAreaId                   | 6ba7b814-9dad-11d1-3-00c04fd430c8 |

  @ARA-1077
  Scenario: Collect GTFS VehiclePosition (with occupancy_status)
    Given a GTFS-RT server waits request on "http://localhost:8090" to respond with
      """
        header {
          gtfs_realtime_version: "2.0"
          incrementality: FULL_DATASET
          timestamp: 1284457468
        }
        entity {
          id: "entity_id"
          vehicle: {
            trip: {
              trip_id: "270856"
              start_time: "09:42:00"
              start_date: "20170313"
              schedule_relationship: SCHEDULED
            }
            position: {
              latitude : -32.92627
              longitude: 151.78036
              bearing  : 91.0
              speed    : 9.8
            }
            timestamp: 1527621931
            vehicle: {
              id   : "bus-234"
            }
            occupancy_status: FEW_SEATS_AVAILABLE
          }
        }
      """
    And a Line exists with the following attributes:
      | Codes | "internal": "1234" |
    And a VehicleJourney exists with the following attributes:
      | Codes  | "internal": "270856"              |
      | LineId | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | Codes | "internal": "1234" |
    And a Partner "gtfs" exists with connectors [gtfs-rt-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_code_space | internal              |
    When a minute has passed
    Then one Vehicle has the following attributes:
      | Codes            | "internal": "bus-234"             |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Occupancy        | fewSeatsAvailable                 |
    And one VehicleJourney has the following attributes:
      | Codes     | "internal": "270856" |
      | Occupancy | fewSeatsAvailable    |

  @ARA-1047
  Scenario: Collect GTFS VehiclePosition (with stop_id)
    Given a GTFS-RT server waits request on "http://localhost:8090" to respond with
      """
        header {
          gtfs_realtime_version: "2.0"
          incrementality: FULL_DATASET
          timestamp: 1284457468
        }
        entity {
          id: "entity_id"
          vehicle: {
            stop_id: "1234"
            trip: {
              trip_id: "270856"
              start_time: "09:42:00"
              start_date: "20170313"
              schedule_relationship: SCHEDULED
            }
            position: {
              latitude : -32.92627
              longitude: 151.78036
              bearing  : 91.0
              speed    : 9.8
            }
            timestamp: 1527621931
            vehicle: {
              id   : "bus-234"
            }
          }
        }
      """
    And a Line exists with the following attributes:
      | Codes | "internal": "1234" |
    And a VehicleJourney exists with the following attributes:
      | Codes  | "internal": "270856"              |
      | LineId | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | Codes | "internal": "1234" |
    And a Partner "gtfs" exists with connectors [gtfs-rt-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_code_space | internal              |
    When a minute has passed
    Then one Vehicle has the following attributes:
      | Codes            | "internal": "bus-234"             |
      | StopAreaId       | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Bearing          | 91.0                              |

  @ARA-1347
  Scenario: Collect GTFS VehiclePosition (with stop_id) should set the NextStopVisitId if StopVisit exists for a given VehicleJourney and StopArea
    Given a GTFS-RT server waits request on "http://localhost:8090" to respond with
      """
        header {
          gtfs_realtime_version: "2.0"
          incrementality: FULL_DATASET
          timestamp: 1284457468
        }
        entity {
          id: "entity_id"
          vehicle: {
            stop_id: "1234"
            trip: {
              trip_id: "270856"
              start_time: "09:42:00"
              start_date: "20170313"
              schedule_relationship: SCHEDULED
            }
            position: {
              latitude : -32.92627
              longitude: 151.78036
              bearing  : 91.0
              speed    : 9.8
            }
            timestamp: 1527621931
            vehicle: {
              id   : "bus-234"
            }
          }
        }
      """
    And a Line exists with the following attributes:
      | Codes | "internal": "1234" |
      # 6ba7b814-9dad-11d1-1-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Codes  | "internal": "270856"              |
      | LineId | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
      # 6ba7b814-9dad-11d1-2-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Codes | "internal": "1234" |
      # 6ba7b814-9dad-11d1-3-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes                   | "internal": "270856-1234"         |
      | StopAreaId              | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | LineId                  | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
      | VehicleJourneyId        | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Schedule[aimed]#Arrival | 2017-01-01T12:55:00.000+02:00     |
      # 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a Partner "gtfs" exists with connectors [gtfs-rt-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_code_space | internal              |
    When a minute has passed
    Then one Vehicle has the following attributes:
      | Codes            | "internal": "bus-234"             |
      | StopAreaId       | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Bearing          | 91.0                              |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-4-00c04fd430c8 |

  @ARA-1347
  Scenario: Collect GTFS VehiclePosition (with stop_id) should not set the NextStopVisitId if multiple StopVisit exists for a given VehicleJourney and StopArea
    Given a GTFS-RT server waits request on "http://localhost:8090" to respond with
      """
        header {
          gtfs_realtime_version: "2.0"
          incrementality: FULL_DATASET
          timestamp: 1284457468
        }
        entity {
          id: "entity_id"
          vehicle: {
            stop_id: "1234"
            trip: {
              trip_id: "270856"
              start_time: "09:42:00"
              start_date: "20170313"
              schedule_relationship: SCHEDULED
            }
            position: {
              latitude : -32.92627
              longitude: 151.78036
              bearing  : 91.0
              speed    : 9.8
            }
            timestamp: 1527621931
            vehicle: {
              id   : "bus-234"
            }
          }
        }
      """
    And a Line exists with the following attributes:
      | Codes | "internal": "1234" |
      # 6ba7b814-9dad-11d1-1-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Codes  | "internal": "270856"              |
      | LineId | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
      # 6ba7b814-9dad-11d1-2-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Codes | "internal": "1234" |
      # 6ba7b814-9dad-11d1-3-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes                   | "internal": "270856-1234-6"       |
      | StopAreaId              | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | LineId                  | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
      | VehicleJourneyId        | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | PassageOrder            | 6                                 |
      | Schedule[aimed]#Arrival | 2017-01-01T12:55:00.000+02:00     |
      # 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes                   | "internal": "270856-1234-22"      |
      | StopAreaId              | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | LineId                  | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
      | VehicleJourneyId        | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | PassageOrder            | 22                                |
      | Schedule[aimed]#Arrival | 2017-01-01T13:55:00.000+02:00     |
      # 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a Partner "gtfs" exists with connectors [gtfs-rt-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_code_space | internal              |
    When a minute has passed
    Then one Vehicle has the following attributes:
      | Codes            | "internal": "bus-234"             |
      | StopAreaId       | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Bearing          | 91.0                              |
    Then No Vehicle exists with the following attributes:
      | NextStopVisitId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
    Then No Vehicle exists with the following attributes:
      | NextStopVisitId | 6ba7b814-9dad-11d1-5-00c04fd430c8 |

  @ARA-1449
  Scenario: Collect GTFS VehiclePosition without occupancy should not change the occupancy status
    Given a GTFS-RT server waits request on "http://localhost:8090" to respond with
      """
        header {
          gtfs_realtime_version: "2.0"
          incrementality: FULL_DATASET
          timestamp: 1284457468
        }
        entity {
          id: "entity_id"
          vehicle: {
            trip: {
              trip_id: "270856"
              start_time: "09:42:00"
              start_date: "20170313"
              schedule_relationship: SCHEDULED
            }
            position: {
              latitude : -32.92627
              longitude: 151.78036
              bearing  : 91.0
              speed    : 9.8
            }
            timestamp: 1527621931
            vehicle: {
              id   : "bus-234"
            }
          }
        }
      """
    And a Line exists with the following attributes:
      | Codes | "internal": "1234" |
    And a VehicleJourney exists with the following attributes:
      | Codes     | "internal": "270856"              |
      | LineId    | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
      | Occupancy | fewSeatsAvailable                 |
    And a StopArea exists with the following attributes:
      | Codes | "internal": "1234" |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "bus-234"             |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Occupancy        | fewSeatsAvailable                 |
    And a Partner "gtfs" exists with connectors [gtfs-rt-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_code_space | internal              |
    When a minute has passed
    Then one Vehicle has the following attributes:
      | Codes            | "internal": "bus-234"             |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Occupancy        | fewSeatsAvailable                 |
    Then one VehicleJourney has the following attributes:
      | Codes     | "internal": "270856" |
      | Occupancy | fewSeatsAvailable    |

    Scenario: Collect GTFS ServiceAlert
    Given a GTFS-RT server waits request on "http://localhost:8090" to respond with
      """
        header {
          gtfs_realtime_version: "2.0"
          incrementality: 0
          timestamp: 1284457468
        }
        entity {
          id: "entity_id"
          alert: {
            active_period: {
                start: 1677664800
                end: 32503698000
              }
            informed_entity: {
                route_id: "Tram:A"
              }
            informed_entity: {
                route_id: "Bus:A"
                stop_id: "94"
              }
            cause: 6
            effect: 4
            header_text: {
              translation: {
                  text: "Route Bus:A detour"
                  language: "fr"
                }
            }
            description_text: {
              translation: {
                  text: "Route Bus:A EASTBOUND is on detour between King at Claren"
                  language: "en"
                }
              translation: {
                  text: "Route Bus:A RETOUR effectue un détour entre King et Claren"
                  language: "fr"
                }
            }
            severity_level: 3
          }
        }
      """
    And a Partner "gtfs" exists with connectors [gtfs-rt-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_code_space | internal              |
    And a Line exists with the following attributes:
      | Codes | "internal": "Tram:A" |
    And a Line exists with the following attributes:
      | Codes | "internal": "Bus:A" |
    And a StopArea exists with the following attributes:
      | Codes | "internal": "94" |
    When a minute has passed
    Then one Situation has the following attributes:
      | Codes                                                          | "internal" : "entity_id"                                   |
      | Origin                                                         | gtfs                                                       |
      | RecordedAt                                                     | 2017-01-01T12:01:00Z                                       |
      | Version                                                        | 5542761353770188895                                        |
      | VersionedAt                                                    | 2010-09-14T11:44:28+02:00                                  |
      | Progress                                                       | published                                                  |
      | ValidityPeriods[0]#EndTime                                     | 3000-01-01T06:00:00+01:00                                  |
      | ValidityPeriods[0]#StartTime                                   | 2023-03-01T11:00:00+01:00                                  |
      | AlertCause                                                     | accident                                                   |
      | Severity                                                       | normal                                                     |
      | Summary[DefaultValue]                                          | Route Bus:A detour                                         |
      | Description[DefaultValue]                                      | Route Bus:A RETOUR effectue un détour entre King et Claren |
      | Affects[StopArea]                                              | 6ba7b814-9dad-11d1-4-00c04fd430c8                          |
      | Affects[StopArea=6ba7b814-9dad-11d1-4-00c04fd430c8]/LineIds[0] | 6ba7b814-9dad-11d1-3-00c04fd430c8                          |
      | Affects[Line]                                                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                          |
    And the Situation "internal":"entity_id" has a Consequence with the following attributes:
      | Condition | diverted |
    And an audit event should exist with these attributes:
      | Protocol  | gtfs                |
      | Direction | sent                |
      | Type      | GtfsRequest         |
      | StopAreas | ["94"]              |
      | Lines     | ["Tram:A", "Bus:A"] |
