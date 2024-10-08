Feature: Support SIRI VehicleMonitoring

  Background:
    Given a Referential "test" exists

  Scenario: Handle a SIRI Lite VehicleMonitoring request
    Given a SIRI Partner "test" exists with connectors [siri-lite-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                |
      | Codes                    | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored                | true                                      |
      | Attribute[DirectionName] | Direction Name                            |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:201123:LOC" |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8     |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8     |
      | Longitude        | 1.234                                 |
      | Latitude         | 5.678                                 |
      | Bearing          | 123                                   |
      | Occupancy        | seatsAvailable                        |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z              |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z              |
      | DriverRef        | "1233"                                |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-6-00c04fd430c8     |
    And a StopArea exists with the following attributes:
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | Name  | Carabacel                                |
    # 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "Test:VehicleJourney:202:LOC-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                  | 4                                                                      |
      | VehicleAtStop                 | false                                                                  |
      | StopAreaId                    | 6ba7b814-9dad-11d1-5-00c04fd430c8                                      |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-3-00c04fd430c8                                      |
      | VehicleAtStop                 | false                                                                  |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                                     |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                                               |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                                               |
      | ArrivalStatus                 | delayed                                                                |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                                               |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                                               |
      | DepartureStatus               | delayed                                                                |
      | Attribute[DestinationDisplay] | Pouet-pouet                                                            |
    # 6ba7b814-9dad-11d1-5-00c04fd430c8
    When I send a vehicle-monitoring SIRI Lite request with the following parameters
      | Token             | test            |
      | LineRef           | Test:Line:3:LOC |
      | MessageIdentifier | Test:1234::LOC  |
    Then I should receive this SIRI Lite response
      """
{
  "Siri": {
    "ServiceDelivery": {
      "ResponseTimestamp": "2017-01-01T12:00:00Z",
      "ProducerRef": "Ara",
      "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": true,
        "VehicleActivity": [{
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T14:00:00Z",
          "VehicleMonitoringRef": "Test:Vehicle:201123:LOC",
          "MonitoredVehicleJourney": {
            "LineRef": "Test:Line:3:LOC",
            "FramedVehicleJourneyRef": {
              "DataFrameRef": "RATPDev:DataFrame::2017-01-01:LOC",
              "DatedVehicleJourneyRef": "Test:VehicleJourney:201:LOC"
            },
            "PublishedLineName": "Ligne 3 Metro",
            "DirectionName": "Direction Name",
            "Monitored": true,
            "Bearing": 123,
            "VehicleLocation": {
              "Longitude": 1.234,
              "Latitude": 5.678
            },
            "Occupancy": "seatsAvailable",
            "DriverRef": "1233",
            "MonitoredCall": {
            "StopPointRef": "NINOXE:StopPoint:SP:24:LOC",
              "StopPointName": "Carabacel",
              "DestinationDisplay": "Pouet-pouet",
              "ExpectedArrivalTime": "2017-01-01T15:02:00Z",
              "ExpectedDepartureTime": "2017-01-01T15:01:00Z",
              "DepartureStatus": "delayed",
              "Order": 4,
              "AimedArrivalTime": "2017-01-01T15:00:00Z",
              "AimedDepartureTime": "2017-01-01T15:01:00Z",
              "ArrivalStatus": "delayed",
              "ActualArrivalTime": "0001-01-01T00:00:00Z",
              "ActualDepartureTime": "0001-01-01T00:00:00Z"
            }
          }
        }]
      }
    }
  }
}
      """

  @ARA-1590
  Scenario: Handle a SIRI Lite VehicleMonitoring request with Referent Line and Referent StopArea and Fallback on vehicle remote codeSpace (RatpCap case)
    Given a SIRI Partner "test" exists with connectors [siri-lite-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential                                                           | test                |
      | remote_code_space                                                          | internal            |
      | siri-lite-vehicle-monitoring-request-broadcaster.vehicle_remote_code_space | rdmantois, rdbievre |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Referent-1" |
      | Name  | Line Referent 1          |
    # 6ba7b814-9dad-11d1-2-00c04fd430c8
    And a Line exists with the following attributes:
      | Codes      | "rdbievre": "Line-1"              |
      | Name       | Ligne 1                           |
      | ReferentId | 6ba7b814-9dad-11d1-2-00c04fd430c8 | # Line Referent 1
      # 6ba7b814-9dad-11d1-3-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                                                   |
      | Codes                    | "rdbievre": "bievre-VehicleJourney","internal": "STIF:bievre-VehicleJourney" |
      | LineId                   | 6ba7b814-9dad-11d1-3-00c04fd430c8                                            |
      | Monitored                | true                                                                         |
      | Attribute[DirectionName] | Direction Name                                                               |
    # 6ba7b814-9dad-11d1-4-00c04fd430c8
    And a Vehicle exists with the following attributes:
      | Codes            | "rdbievre": "bievre-Vehicle"      |
      | LineId           | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 123                               |
      | Occupancy        | seatsAvailable                    |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
      | DriverRef        | "1233"                            |
      | NextStopVisitId  | 6ba7b814-9dad-11d1-7-00c04fd430c8 |
    # 6ba7b814-9dad-11d1-5-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Codes      | "rdbievre": "Stop-1"              |
      | Name       | Stop 1                            |
      | ReferentId | 6ba7b814-9dad-11d1-b-00c04fd430c8 | # Stop Referent
    # 6ba7b814-9dad-11d1-6-00c04fd430c8
    And a StopVisit exists with the following attributes:
      | Codes                         | "internal": "bievre-VehicleJourney-bievre-Vehicle" |
      | PassageOrder                  | 4                                                  |
      | VehicleAtStop                 | false                                              |
      | StopAreaId                    | 6ba7b814-9dad-11d1-6-00c04fd430c8                  |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-4-00c04fd430c8                  |
      | VehicleAtStop                 | false                                              |
      | Reference[OperatorRef]#Code   | "internal": "CdF:Company::410:LOC"                 |
      | Schedule[aimed]#Arrival       | 2017-01-01T15:00:00.000Z                           |
      | Schedule[expected]#Arrival    | 2017-01-01T15:01:00.000Z                           |
      | ArrivalStatus                 | delayed                                            |
      | Schedule[aimed]#Departure     | 2017-01-01T15:01:00.000Z                           |
      | Schedule[expected]#Departure  | 2017-01-01T15:02:00.000Z                           |
      | DepartureStatus               | delayed                                            |
      | Attribute[DestinationDisplay] | Pouet-pouet                                        |
    # 6ba7b814-9dad-11d1-7-00c04fd430c8
    And a Line exists with the following attributes:
      | Codes      | "rdmantois": "Line-2"             |
      | Name       | Line 2                            |
      | ReferentId | 6ba7b814-9dad-11d1-2-00c04fd430c8 | # Line Referent 1
    # 6ba7b814-9dad-11d1-8-00c04fd430c8
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                                                       |
      | Codes                    | "rdmantois": "mantois-VehicleJourney", "internal": "STIF:mantois-VehicleJourney" |
      | LineId                   | 6ba7b814-9dad-11d1-8-00c04fd430c8                                                |
      | Monitored                | true                                                                             |
      | Attribute[DirectionName] | Another Direction Name                                                           |
    # 6ba7b814-9dad-11d1-9-00c04fd430c8
    And a Vehicle exists with the following attributes:
      | Codes            | "rdmantois": "mantois-Vehicle"    |
      | LineId           | 6ba7b814-9dad-11d1-8-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-9-00c04fd430c8 |
      | Longitude        | 3.232                             |
      | Latitude         | 8.329                             |
      | Bearing          | 355                               |
      | Occupancy        | fewSeatsAvailable                 |
      | RecordedAtTime   | 2017-01-01T14:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T15:00:00.000Z          |
      | DriverRef        | "567"                             |
    # 6ba7b814-9dad-11d1-a-00c04fd430c8
    And a StopArea exists with the following attributes:
      | Codes | "internal": "Stop-Referent-1" |
      | Name  | Stop Referent                 |
    When I send a vehicle-monitoring SIRI Lite request with the following parameters
      | Token             | test           |
      | LineRef           | Referent-1     |
      | MessageIdentifier | Test:1234::LOC |
    Then I should receive this SIRI Lite response
      """
      {
        "Siri": {
          "ServiceDelivery": {
            "ResponseTimestamp": "2017-01-01T12:00:00Z",
            "ProducerRef": "Ara",
            "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-c-00c04fd430c8:LOC",
            "RequestMessageRef": "Test:1234::LOC",
            "VehicleMonitoringDelivery": {
              "Version": "2.0:FR-IDF-2.4",
              "ResponseTimestamp": "2017-01-01T12:00:00Z",
              "RequestMessageRef": "Test:1234::LOC",
              "Status": true,
              "VehicleActivity": [
                {
                  "RecordedAtTime": "2017-01-01T13:00:00Z",
                  "ValidUntilTime": "2017-01-01T14:00:00Z",
                  "VehicleMonitoringRef": "bievre-Vehicle",
                  "MonitoredVehicleJourney": {
                    "LineRef": "Referent-1",
                    "FramedVehicleJourneyRef": {
                      "DataFrameRef": "RATPDev:DataFrame::2017-01-01:LOC",
                      "DatedVehicleJourneyRef": "STIF:bievre-VehicleJourney"
                    },
                    "PublishedLineName": "Line Referent 1",
                    "DirectionName": "Direction Name",
                    "Monitored": true,
                    "Bearing": 123,
                    "VehicleLocation": {
                      "Longitude": 1.234,
                      "Latitude": 5.678
                    },
                    "Occupancy": "seatsAvailable",
                    "DriverRef": "1233",
                    "MonitoredCall": {
                      "StopPointRef": "Stop-Referent-1",
                      "StopPointName": "Stop Referent",
                      "DestinationDisplay": "Pouet-pouet",
                      "ExpectedArrivalTime": "2017-01-01T15:02:00Z",
                      "ExpectedDepartureTime": "2017-01-01T15:01:00Z",
                      "DepartureStatus": "delayed",
                      "Order": 4,
                      "AimedArrivalTime": "2017-01-01T15:00:00Z",
                      "AimedDepartureTime": "2017-01-01T15:01:00Z",
                      "ArrivalStatus": "delayed",
                      "ActualArrivalTime": "0001-01-01T00:00:00Z",
                      "ActualDepartureTime": "0001-01-01T00:00:00Z"
                    }
                  }
                },
                {
                  "RecordedAtTime": "2017-01-01T14:00:00Z",
                  "ValidUntilTime": "2017-01-01T15:00:00Z",
                  "VehicleMonitoringRef": "mantois-Vehicle",
                  "MonitoredVehicleJourney": {
                    "LineRef": "Referent-1",
                    "FramedVehicleJourneyRef": {
                      "DataFrameRef": "RATPDev:DataFrame::2017-01-01:LOC",
                      "DatedVehicleJourneyRef": "STIF:mantois-VehicleJourney"
                    },
                    "PublishedLineName": "Line Referent 1",
                    "DirectionName": "Another Direction Name",
                    "Monitored": true,
                    "Bearing": 355,
                    "VehicleLocation": {
                      "Longitude": 3.232,
                      "Latitude": 8.329
                    },
                    "Occupancy": "fewSeatsAvailable",
                    "DriverRef": "567"
                  }
                }
              ]
            }
          }
        }

      }
      """
    Then an audit event should exist with these attributes:
        | Type            | VehicleMonitoringRequest                                      |
        | Protocol        | siri-lite                                                     |
        | Direction       | received                                                      |
        | Status          | OK                                                            |
        | Partner         | test                                                          |
        | Vehicles        | ["bievre-Vehicle", "mantois-Vehicle"]                         |
        | Lines           | ["Referent-1"]                                                |
        | VehicleJourneys | ["STIF:bievre-VehicleJourney", "STIF:mantois-VehicleJourney"] |

  Scenario: Send the correct vehicles to a SIRI Lite VehicleMonitoring request
    Given a SIRI Partner "test" exists with connectors [siri-lite-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:2:LOC" |
      | Name  | Ligne 2 Metro                 |
    And a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                                |
      | Codes     | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored | true                                      |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 33                                |
      | Codes     | "internal": "Test:VehicleJourney:202:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8         |
      | Monitored | true                                      |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:2:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 123                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:3:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Longitude        | 1.345                             |
      | Latitude         | 5.789                             |
      | Bearing          | 456                               |
      | RecordedAtTime   | 2017-01-01T14:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
    When I send a vehicle-monitoring SIRI Lite request with the following parameters
      | Token             | test            |
      | LineRef           | Test:Line:2:LOC |
      | MessageIdentifier | Test:1234::LOC  |
    Then I should receive this SIRI Lite response
      """
{
  "Siri": {
    "ServiceDelivery": {
      "ResponseTimestamp": "2017-01-01T12:00:00Z",
      "ProducerRef": "Ara",
      "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": true,
        "VehicleActivity": [{
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T14:00:00Z",
          "VehicleMonitoringRef": "Test:Vehicle:2:LOC",
          "MonitoredVehicleJourney": {
            "LineRef": "Test:Line:2:LOC",
            "FramedVehicleJourneyRef": {
              "DataFrameRef": "RATPDev:DataFrame::2017-01-01:LOC",
              "DatedVehicleJourneyRef": "Test:VehicleJourney:201:LOC"
            },
            "PublishedLineName": "Ligne 2 Metro",
            "Monitored": true,
            "Bearing": 123,
            "VehicleLocation": {
              "Longitude": 1.234,
              "Latitude": 5.678
            }
          }
        }]
      }
    }
  }
}
      """

  Scenario: Send all the vehicles to a SIRI Lite VehicleMonitoring request
    Given a SIRI Partner "test" exists with connectors [siri-lite-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential      | test     |
      | remote_code_space     | internal |
      | sort_payload_for_test | true     |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                                |
      | Codes     | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored | true                                      |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 33                                |
      | Codes     | "internal": "Test:VehicleJourney:202:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored | true                                      |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:1:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 123                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:2:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 123                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
    And a Vehicle exists with the following attributes:
      | Codes            | "internal": "Test:Vehicle:3:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 123                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z          |
    When I send a vehicle-monitoring SIRI Lite request with the following parameters
      | Token             | test            |
      | LineRef           | Test:Line:3:LOC |
      | MessageIdentifier | Test:1234::LOC  |
    Then I should receive this SIRI Lite response
      """
{
  "Siri": {
    "ServiceDelivery": {
      "ResponseTimestamp": "2017-01-01T12:00:00Z",
      "ProducerRef": "Ara",
      "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": true,
        "VehicleActivity": [{
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T14:00:00Z",
          "VehicleMonitoringRef": "Test:Vehicle:1:LOC",
          "MonitoredVehicleJourney": {
            "LineRef": "Test:Line:3:LOC",
            "FramedVehicleJourneyRef": {
              "DataFrameRef": "RATPDev:DataFrame::2017-01-01:LOC",
              "DatedVehicleJourneyRef": "Test:VehicleJourney:201:LOC"
            },
            "PublishedLineName": "Ligne 3 Metro",
            "Monitored": true,
            "Bearing": 123,
            "VehicleLocation": {
              "Longitude": 1.234,
              "Latitude": 5.678
            }
          }
        },
        {
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T14:00:00Z",
          "VehicleMonitoringRef": "Test:Vehicle:2:LOC",
          "MonitoredVehicleJourney": {
            "LineRef": "Test:Line:3:LOC",
            "FramedVehicleJourneyRef": {
              "DataFrameRef": "RATPDev:DataFrame::2017-01-01:LOC",
              "DatedVehicleJourneyRef": "Test:VehicleJourney:201:LOC"
            },
            "PublishedLineName": "Ligne 3 Metro",
            "Monitored": true,
            "Bearing": 123,
            "VehicleLocation": {
              "Longitude": 1.234,
              "Latitude": 5.678
            }
          }
        },
        {
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T14:00:00Z",
          "VehicleMonitoringRef": "Test:Vehicle:3:LOC",
          "MonitoredVehicleJourney": {
            "LineRef": "Test:Line:3:LOC",
            "FramedVehicleJourneyRef": {
              "DataFrameRef": "RATPDev:DataFrame::2017-01-01:LOC",
              "DatedVehicleJourneyRef": "Test:VehicleJourney:202:LOC"
            },
            "PublishedLineName": "Ligne 3 Metro",
            "Monitored": true,
            "Bearing": 123,
            "VehicleLocation": {
              "Longitude": 1.234,
              "Latitude": 5.678
            }
          }
        }]
      }
    }
  }
}
      """

  @ARTA-1040
  Scenario: Handle a SIRI Lite VehicleMonitoring request StopMonitoring request with unmatching code kind
   Given a SIRI Partner "test" exists with connectors [siri-lite-vehicle-monitoring-request-broadcaster] and the following settings:
     | local_credential  | test  |
     | remote_code_space | wrong |
   Given a Line exists with the following attributes:
     | Codes | "internal": "Test:Line:3:LOC" |
     | Name  | Ligne 3 Metro                 |
   And a VehicleJourney exists with the following attributes:
     | Name                     | Passage 32                                |
     | Codes                    | "internal": "Test:VehicleJourney:201:LOC" |
     | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
     | Monitored                | true                                      |
     | Attribute[DirectionName] | Direction Name                            |
   And a Vehicle exists with the following attributes:
     | Codes            | "other": "Test:Vehicle:201123:LOC" |
     | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8  |
     | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8  |
     | Longitude        | 1.234                              |
     | Latitude         | 5.678                              |
     | Bearing          | 123                                |
     | RecordedAtTime   | 2017-01-01T13:00:00.000Z           |
     | ValidUntilTime   | 2017-01-01T14:00:00.000Z           |
   When I send a vehicle-monitoring SIRI Lite request with the following parameters
     | Token             | test            |
     | LineRef           | Test:Line:3:LOC |
     | MessageIdentifier | Test:1234::LOC  |
   Then I should receive this SIRI Lite response
      """
{
  "Siri": {
    "ServiceDelivery": {
      "ResponseTimestamp": "2017-01-01T12:00:00Z",
      "ProducerRef": "Ara",
      "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": false,
        "ErrorCondition": {
          "ErrorType": "InvalidDataReferencesError",
          "ErrorNumber": 0,
          "ErrorText": "Line Test:Line:3:LOC not found"
        },
        "VehicleActivity": []
      }
    }
  }
}
      """

  @ARA-1044
  Scenario: Handle a SIRI Lite VehicleMonitoring request with multiple connector setting siri-lite-stop-monitoring-request-broadcaster.vehicle_journey_remote_code_space
   Given a SIRI Partner "test" exists with connectors [siri-lite-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential                                                           | test          |
      | remote_code_space                                                          | internal      |
      | siri-lite-vehicle-monitoring-request-broadcaster.vehicle_remote_code_space | other, other2 |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                |
      | Codes                    | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored                | true                                      |
      | Attribute[DirectionName] | Direction Name                            |
    And a Vehicle exists with the following attributes:
      | Codes            | "other": "Test:Vehicle:201123:LOC" |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8  |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8  |
      | Longitude        | 1.234                              |
      | Latitude         | 5.678                              |
      | Bearing          | 123                                |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z           |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z           |
    When I send a vehicle-monitoring SIRI Lite request with the following parameters
      | Token             | test            |
      | LineRef           | Test:Line:3:LOC |
      | MessageIdentifier | Test:1234::LOC  |
    Then I should receive this SIRI Lite response
      """
{
  "Siri": {
    "ServiceDelivery": {
      "ResponseTimestamp": "2017-01-01T12:00:00Z",
      "ProducerRef": "Ara",
      "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": true,
        "VehicleActivity": [{
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T14:00:00Z",
          "VehicleMonitoringRef": "Test:Vehicle:201123:LOC",
          "MonitoredVehicleJourney": {
            "LineRef": "Test:Line:3:LOC",
            "FramedVehicleJourneyRef": {
              "DataFrameRef": "RATPDev:DataFrame::2017-01-01:LOC",
              "DatedVehicleJourneyRef": "Test:VehicleJourney:201:LOC"
            },
            "PublishedLineName": "Ligne 3 Metro",
            "DirectionName": "Direction Name",
            "Monitored": true,
            "Bearing": 123,
            "VehicleLocation": {
              "Longitude": 1.234,
              "Latitude": 5.678
            }
          }
        }]
      }
    }
  }
}
      """
    Then an audit event should exist with these attributes:
        | Type            | VehicleMonitoringRequest        |
        | Protocol        | siri-lite                       |
        | Direction       | received                        |
        | Status          | OK                              |
        | Partner         | test                            |
        | Vehicles        | ["Test:Vehicle:201123:LOC"]     |
        | Lines           | ["Test:Line:3:LOC"]             |
        | VehicleJourneys | ["Test:VehicleJourney:201:LOC"] |

  @ARA-1044
  Scenario: Handle a SIRI Lite VehicleMonitoring request with fallback on generic connector remote_code_space
   Given a SIRI Partner "test" exists with connectors [siri-lite-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential                                                   | test     |
      | remote_code_space                                                  | internal |
      | siri-lite-vehicle-monitoring-request-broadcaster.remote_code_space | other    |
    Given a Line exists with the following attributes:
      | Codes | "other": "Test:Line:3:LOC" |
      | Name  | Ligne 3 Metro              |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                             |
      | Codes                    | "other": "Test:VehicleJourney:201:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8      |
      | Monitored                | true                                   |
      | Attribute[DirectionName] | Direction Name                         |
    And a Vehicle exists with the following attributes:
      | Codes            | "other": "Test:Vehicle:201123:LOC" |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8  |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8  |
      | Longitude        | 1.234                              |
      | Latitude         | 5.678                              |
      | Bearing          | 123                                |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z           |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z           |
    When I send a vehicle-monitoring SIRI Lite request with the following parameters
      | Token             | test            |
      | LineRef           | Test:Line:3:LOC |
      | MessageIdentifier | Test:1234::LOC  |
    Then I should receive this SIRI Lite response
      """
{
  "Siri": {
    "ServiceDelivery": {
      "ResponseTimestamp": "2017-01-01T12:00:00Z",
      "ProducerRef": "Ara",
      "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": true,
        "VehicleActivity": [{
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T14:00:00Z",
          "VehicleMonitoringRef": "Test:Vehicle:201123:LOC",
          "MonitoredVehicleJourney": {
            "LineRef": "Test:Line:3:LOC",
            "FramedVehicleJourneyRef": {
              "DataFrameRef": "RATPDev:DataFrame::2017-01-01:LOC",
              "DatedVehicleJourneyRef": "Test:VehicleJourney:201:LOC"
            },
            "PublishedLineName": "Ligne 3 Metro",
            "DirectionName": "Direction Name",
            "Monitored": true,
            "Bearing": 123,
            "VehicleLocation": {
              "Longitude": 1.234,
              "Latitude": 5.678
            }
          }
        }]
      }
    }
  }
}
      """
  @ARA-1363
  Scenario: Handle a SIRI Lite VehicleMonitoring request using the generator setting reference_vehicle_journey_identifier
    # Setting a Partner without default generators
    Given a Partner "test" exists with connectors [siri-lite-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential                                | test                             |
      | remote_code_space                            | internal                         |
      | generators.reference_vehicle_journey_identifier | ch:1:ServiceJourney:87_TAC:%{id} |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name      | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                                        |
      | Codes                | "_default": "6ba7b814", "external": "Test:VehicleJourney:201:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8                                 |
      | Monitored                | true                                                              |
      | Attribute[DirectionName] | Direction Name                                                    |
    And a Vehicle exists with the following attributes:
      | Codes        | "internal": "Test:Vehicle:201123:LOC" |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8     |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8     |
      | Longitude        | 1.234                                 |
      | Latitude         | 5.678                                 |
      | Bearing          | 123                                   |
      | Occupancy        | seatsAvailable                        |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z              |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z              |
      | DriverRef        | "1233"                                |
    When I send a vehicle-monitoring SIRI Lite request with the following parameters
      | Token             | test            |
      | LineRef           | Test:Line:3:LOC |
      | MessageIdentifier | Test:1234::LOC  |
    Then I should receive this SIRI Lite response
      """
{
  "Siri": {
    "ServiceDelivery": {
      "ResponseTimestamp": "2017-01-01T12:00:00Z",
      "ProducerRef": "Ara",
      "ResponseMessageIdentifier": "6ba7b814-9dad-11d1-5-00c04fd430c8",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": true,
        "VehicleActivity": [{
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T14:00:00Z",
          "VehicleMonitoringRef": "Test:Vehicle:201123:LOC",
          "MonitoredVehicleJourney": {
            "LineRef": "Test:Line:3:LOC",
            "FramedVehicleJourneyRef": {
              "DataFrameRef": "2017-01-01",
              "DatedVehicleJourneyRef": "ch:1:ServiceJourney:87_TAC:6ba7b814"
            },
            "PublishedLineName": "Ligne 3 Metro",
            "DirectionName": "Direction Name",
            "Monitored": true,
            "Bearing": 123,
            "VehicleLocation": {
              "Longitude": 1.234,
              "Latitude": 5.678
            },
            "Occupancy": "seatsAvailable",
            "DriverRef": "1233"
          }
        }]
      }
    }
  }
}
      """
    Then an audit event should exist with these attributes:
        | Type              | VehicleMonitoringRequest                |
        | Protocol          | siri-lite                               |
        | Direction         | received                                |
        | Status            | OK                                      |
        | Partner           | test                                    |
        | Vehicles          | ["Test:Vehicle:201123:LOC"]             |
        | RequestIdentifier | Test:1234::LOC                          |
        | Lines             | ["Test:Line:3:LOC"]                     |
        | VehicleJourneys   | ["ch:1:ServiceJourney:87_TAC:6ba7b814"] |

  @ARA-1363
  Scenario: Handle a SIRI Lite VehicleMonitoring request using the default generator should send DatedVehicleJourneyRef according to default setting
    # Setting a "SIRI Partner" with default generators
    Given a SIRI Partner "test" exists with connectors [siri-lite-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_code_space | internal |
    Given a Line exists with the following attributes:
      | Codes | "internal": "Test:Line:3:LOC" |
      | Name      | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                                        |
      | Codes                | "_default": "6ba7b814", "external": "Test:VehicleJourney:201:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8                                 |
      | Monitored                | true                                                              |
      | Attribute[DirectionName] | Direction Name                                                    |
    And a Vehicle exists with the following attributes:
      | Codes        | "internal": "Test:Vehicle:201123:LOC" |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8     |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8     |
      | Longitude        | 1.234                                 |
      | Latitude         | 5.678                                 |
      | Bearing          | 123                                   |
      | Occupancy        | seatsAvailable                        |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z              |
      | ValidUntilTime   | 2017-01-01T14:00:00.000Z              |
      | DriverRef        | "1233"                                |
    When I send a vehicle-monitoring SIRI Lite request with the following parameters
      | Token             | test            |
      | LineRef           | Test:Line:3:LOC |
      | MessageIdentifier | Test:1234::LOC  |
    Then I should receive this SIRI Lite response
      """
{
  "Siri": {
    "ServiceDelivery": {
      "ResponseTimestamp": "2017-01-01T12:00:00Z",
      "ProducerRef": "Ara",
      "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": true,
        "VehicleActivity": [{
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T14:00:00Z",
          "VehicleMonitoringRef": "Test:Vehicle:201123:LOC",
          "MonitoredVehicleJourney": {
            "LineRef": "Test:Line:3:LOC",
            "FramedVehicleJourneyRef": {
              "DataFrameRef": "RATPDev:DataFrame::2017-01-01:LOC",
              "DatedVehicleJourneyRef": "RATPDev:VehicleJourney::6ba7b814:LOC"
            },
            "PublishedLineName": "Ligne 3 Metro",
            "DirectionName": "Direction Name",
            "Monitored": true,
            "Bearing": 123,
            "VehicleLocation": {
              "Longitude": 1.234,
              "Latitude": 5.678
            },
            "Occupancy": "seatsAvailable",
            "DriverRef": "1233"
          }
        }]
      }
    }
  }
}
      """
    Then an audit event should exist with these attributes:
        | Type              | VehicleMonitoringRequest                 |
        | Protocol          | siri-lite                                |
        | Direction         | received                                 |
        | Status            | OK                                       |
        | Partner           | test                                     |
        | Vehicles          | ["Test:Vehicle:201123:LOC"]              |
        | RequestIdentifier | Test:1234::LOC                           |
        | Lines             | ["Test:Line:3:LOC"]                      |
        | VehicleJourneys   | ["RATPDev:VehicleJourney::6ba7b814:LOC"] |
