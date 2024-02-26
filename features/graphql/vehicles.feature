Feature: GraphQL API for vehicles

  Background:
    Given a Referential "test" is created

  @ARA-1439
  Scenario: Handle a Vehicle Graphql query
    Given a Partner "test" exists with connectors [graphql-server] and the following settings:
      | local_credential  | test_token |
      | remote_code_space | internal   |
    And a Vehicle exists with the following attributes:
      | Codes          | "internal": "Test:Vehicle:1:LOC" |
      | Longitude      | 1.234                            |
      | Latitude       | 5.678                            |
      | Bearing        | 123                              |
      | Occupancy      | fewSeatsAvailable                |
      | Percentage     | 15.6                             |
      | RecordedAtTime | 2017-01-01T13:00:00.000Z         |
      | ValidUntilTime | 2017-01-01T14:00:00.000Z         |
    When I send this GraphQL query to the Referential "test" with token "test_token"
      """
      query {
        vehicle(code: "Test:Vehicle:1:LOC") {
          code
          longitude
          latitude
          bearing
          occupancyStatus
          occupancyRate
          recordedAt
          validUntil
        }
      }
      """
    Then the GraphQL response should contain one Vehicle with these attributes:
      | code            | Test:Vehicle:1:LOC   |
      | longitude       | 1.234                |
      | latitude        | 5.678                |
      | bearing         | 123                  |
      | occupancyStatus | fewSeatsAvailable    |
      | occupancyRate   | 15.6                 |
      | recordedAt      | 2017-01-01T13:00:00Z |
      | validUntil      | 2017-01-01T14:00:00Z |

  @ARA-1439
  Scenario: Handle a Vehicles Graphql query
    Given a Partner "test" exists with connectors [graphql-server] and the following settings:
      | local_credential  | test_token |
      | remote_code_space | internal   |
    And a Vehicle exists with the following attributes:
      | Codes          | "internal": "Test:Vehicle:1:LOC" |
      | Longitude      | 1.234                            |
      | Latitude       | 5.678                            |
      | Bearing        | 123                              |
      | Occupancy      | fewSeatsAvailable                |
      | Percentage     | 15.6                             |
      | RecordedAtTime | 2017-01-01T13:00:00.000Z         |
      | ValidUntilTime | 2017-01-01T14:00:00.000Z         |
    When I send this GraphQL query to the Referential "test" with token "test_token"
      """
      query {
        vehicles {
          code
          longitude
          latitude
          bearing
          occupancyStatus
          occupancyRate
          recordedAt
          validUntil
        }
      }
      """
    Then the GraphQL response should contain a Vehicle with these attributes:
      | code            | Test:Vehicle:1:LOC   |
      | longitude       | 1.234                |
      | latitude        | 5.678                |
      | bearing         | 123                  |
      | occupancyStatus | fewSeatsAvailable    |
      | occupancyRate   | 15.6                 |
      | recordedAt      | 2017-01-01T13:00:00Z |
      | validUntil      | 2017-01-01T14:00:00Z |

  @ARA-1439
  Scenario: Handle a Vehicle Graphql mutation without the setting
    Given a Partner "test" exists with connectors [graphql-server] and the following settings:
      | local_credential  | test_token |
      | remote_code_space | internal   |
    And a Vehicle exists with the following attributes:
      | Codes          | "internal": "Test:Vehicle:1:LOC" |
      | Longitude      | 1.234                            |
      | Latitude       | 5.678                            |
      | Bearing        | 123                              |
      | Occupancy      | fewSeatsAvailable                |
      | Percentage     | 15.6                             |
      | RecordedAtTime | 2017-01-01T13:00:00.000Z         |
      | ValidUntilTime | 2017-01-01T14:00:00.000Z         |
    When I send this GraphQL query to the Referential "test" with token "test_token"
      """
      mutation {
          updateVehicle(code: "Test:Vehicle:1:LOC", input: {occupancyStatus: "seatsAvailable", occupancyRate: 0.65}) {
            code
            longitude
            latitude
            bearing
            occupancyStatus
            occupancyRate
            recordedAt
            validUntil
          }
        }
      """
    Then the GraphQL response should contain an updated Vehicle with these attributes:
      | code            | Test:Vehicle:1:LOC   |
      | longitude       | 1.234                |
      | latitude        | 5.678                |
      | bearing         | 123                  |
      | occupancyStatus | fewSeatsAvailable    |
      | occupancyRate   | 15.6                 |
      | recordedAt      | 2017-01-01T13:00:00Z |
      | validUntil      | 2017-01-01T14:00:00Z |

  @ARA-1439
  Scenario: Handle a Vehicle Graphql mutation with the correct settings
    Given a Partner "test" exists with connectors [graphql-server] and the following settings:
      | local_credential           | test_token                                    |
      | remote_code_space          | internal                                      |
      | graphql.mutable_attributes | vehicle.occupancyStatus,vehicle.occupancyRate |
    And a Vehicle exists with the following attributes:
      | Codes          | "internal": "Test:Vehicle:1:LOC" |
      | Longitude      | 1.234                            |
      | Latitude       | 5.678                            |
      | Bearing        | 123                              |
      | Occupancy      | fewSeatsAvailable                |
      | Percentage     | 15.6                             |
      | RecordedAtTime | 2017-01-01T13:00:00.000Z         |
      | ValidUntilTime | 2017-01-01T14:00:00.000Z         |
    When I send this GraphQL query to the Referential "test" with token "test_token"
      """
      mutation {
        updateVehicle(code: "Test:Vehicle:1:LOC", input: { occupancyStatus: "seatsAvailable", occupancyRate: 0.65 }) {
          code
          longitude
          latitude
          bearing
          occupancyStatus
          occupancyRate
          recordedAt
          validUntil
        }
      }
      """
    Then the GraphQL response should contain an updated Vehicle with these attributes:
      | code            | Test:Vehicle:1:LOC   |
      | longitude       | 1.234                |
      | latitude        | 5.678                |
      | bearing         | 123                  |
      | occupancyStatus | seatsAvailable       |
      | occupancyRate   | 0.65                 |
      | recordedAt      | 2017-01-01T13:00:00Z |
      | validUntil      | 2017-01-01T14:00:00Z |
