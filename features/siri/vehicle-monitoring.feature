Feature: Support SIRI VehicleMonitoring

  Background:
    Given a Referential "test" exists
    Given a Partner "test" exists with connectors [siri-lite-vehicle-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |

  Scenario: Handle a SIRI Lite VehicleMonitoring request
    Given a Line exists with the following attributes:
      | ObjectIDs | "internal": "Test:Line:3:LOC" |
      | Name      | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name                     | Passage 32                                |
      | ObjectIDs                | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored                | true                                      |
      | Attribute[DirectionName] | Direction Name                            |
    And a Vehicle exists with the following attributes:
      | ObjectIDs                | "internal": "Test:Vehicle:201123:LOC" |
      | LineId                   | 6ba7b814-9dad-11d1-2-00c04fd430c8     |
      | VehicleJourneyId         | 6ba7b814-9dad-11d1-3-00c04fd430c8     |
      | Longitude                | 1.234                                 |
      | Latitude                 | 5.678                                 |
      | Bearing                  | 123                                   |
      | RecordedAtTime           | 2017-01-01T13:00:00.000Z              |
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
      "ProducerRef": "Edwig",
      "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-5-00c04fd430c8:LOC",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": true,
        "VehicleActivity": [{
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T13:00:00Z",
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

  Scenario: Send the correct vehicles to a SIRI Lite VehicleMonitoring request
    Given a Line exists with the following attributes:
      | ObjectIDs | "internal": "Test:Line:2:LOC" |
      | Name      | Ligne 2 Metro                 |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "Test:Line:3:LOC" |
      | Name      | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                                |
      | ObjectIDs | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored | true                                      |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 33                                |
      | ObjectIDs | "internal": "Test:VehicleJourney:202:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8         |
      | Monitored | true                                      |
    And a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "Test:Vehicle:2:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 123                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
    And a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "Test:Vehicle:3:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Longitude        | 1.345                             |
      | Latitude         | 5.789                             |
      | Bearing          | 456                               |
      | RecordedAtTime   | 2017-01-01T14:00:00.000Z          |
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
      "ProducerRef": "Edwig",
      "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": true,
        "VehicleActivity": [{
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T13:00:00Z",
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

  @wip
  # See ARA-712
  Scenario: Send all the vehicles to a SIRI Lite VehicleMonitoring request
    Given a Line exists with the following attributes:
      | ObjectIDs | "internal": "Test:Line:3:LOC" |
      | Name      | Ligne 3 Metro                 |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                                |
      | ObjectIDs | "internal": "Test:VehicleJourney:201:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored | true                                      |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 33                                |
      | ObjectIDs | "internal": "Test:VehicleJourney:202:LOC" |
      | LineId    | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Monitored | true                                      |
    And a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "Test:Vehicle:1:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 123                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
    And a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "Test:Vehicle:2:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 123                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
    And a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "Test:Vehicle:3:LOC"  |
      | LineId           | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Longitude        | 1.234                             |
      | Latitude         | 5.678                             |
      | Bearing          | 123                               |
      | RecordedAtTime   | 2017-01-01T13:00:00.000Z          |
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
      "ProducerRef": "Edwig",
      "ResponseMessageIdentifier": "RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC",
      "RequestMessageRef": "Test:1234::LOC",
      "VehicleMonitoringDelivery": {
        "Version": "2.0:FR-IDF-2.4",
        "ResponseTimestamp": "2017-01-01T12:00:00Z",
        "RequestMessageRef": "Test:1234::LOC",
        "Status": true,
        "VehicleActivity": [{
          "RecordedAtTime": "2017-01-01T13:00:00Z",
          "ValidUntilTime": "2017-01-01T13:00:00Z",
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
          "ValidUntilTime": "2017-01-01T13:00:00Z",
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
          "ValidUntilTime": "2017-01-01T13:00:00Z",
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
