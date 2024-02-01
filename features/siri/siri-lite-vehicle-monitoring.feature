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
