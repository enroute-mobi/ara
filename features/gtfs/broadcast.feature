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
      | local_credential | secret |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response

  Scenario: Provide a authenticated GTFS-RT feed (multiple credentials)
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credentials | secret1,secret2 |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" with token "secret1"
    Then I should receive a GTFS-RT response
    When I send a GTFS-RT request to the Referential "test" with token "secret2"
    Then I should receive a GTFS-RT response

  Scenario: Forbid authorized request on GTFS-RT feed (no token)
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credential | secret |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" without token
    Then I should not receive a GTFS-RT but an unauthorized client error status

  Scenario: Forbid authorized request on GTFS-RT feed (wrong token)
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credential | secret |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" with token "wrong"
    Then I should not receive a GTFS-RT but an unauthorized client error status

  Scenario: Retrieve Vehicle Positions
    Given a Line exists with the following attributes:
      | Name      | Test               |
      | ObjectIDs | "internal": "1234" |
    Given a VehicleJourney exists with the following attributes:
      | ObjectIDs | "internal": "2345" |
      | LineId           | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    Given a Vehicle exists with the following attributes:
      | ObjectIDs        | "internal": "3456"                |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a Partner "test" exists with connectors [gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | local_credential     | secret   |
      | remote_objectid_kind | internal |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should contain a Vehicle Position with these attributes:
      | vehicle_id | 3456 |
      | trip_id    | 2345 |
      | route_id   | 1234 |

  @ARA-872
  Scenario: Retrieve Vehicle Positions with unmatching objectid kind
    Given a Line exists with the following attributes:
      | Name      | Test               |
      | ObjectIDs | "other": "1234"    |
    Given a VehicleJourney exists with the following attributes:
      | ObjectIDs | "internal": "2345" |
      | LineId           | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
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
  Scenario: Retrieve Vehicle Positions with setting gtfs-rt-vehicle-positions-broadcaster.vehicle_remote_objectid_kind
    Given a Line exists with the following attributes:
      | Name      | Test               |
      | ObjectIDs | "internal": "1234" |
    Given a VehicleJourney exists with the following attributes:
      | ObjectIDs | "internal": "2345" |
      | LineId           | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
    Given a Vehicle exists with the following attributes:
      | ObjectIDs        | "other": "3456"                   |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a Partner "test" exists with connectors [gtfs-rt-vehicle-positions-broadcaster] and the following settings:
      | local_credential     | secret   |
      | remote_objectid_kind | internal |
      | gtfs-rt-vehicle-positions-broadcaster.vehicle_remote_objectid_kind | other |
    When I send a GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response
    And this GTFS-RT response should contain a Vehicle Position with these attributes:
      | vehicle_id | 3456 |
      | trip_id    | 2345 |
      | route_id   | 1234 |
