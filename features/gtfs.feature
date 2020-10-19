Feature: Support GTFS-RT feeds
  Background:
    Given a Referential "test" is created

  Scenario: Provide a public GTFS-RT feed
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | remote_objectid_kind | internal |
    When I send this GTFS-RT request to the Referential "test" without token
    Then I should receive a GTFS-RT response

  Scenario: Provide a authenticated GTFS-RT feed
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credential | secret |
      | remote_objectid_kind | internal |
    When I send this GTFS-RT request to the Referential "test" with token "secret"
    Then I should receive a GTFS-RT response

  Scenario: Provide a authenticated GTFS-RT feed (multiple credentials)
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credentials | secret1,secret2 |
      | remote_objectid_kind | internal |
    When I send this GTFS-RT request to the Referential "test" with token "secret1"
    Then I should receive a GTFS-RT response
    When I send this GTFS-RT request to the Referential "test" with token "secret2"
    Then I should receive a GTFS-RT response

  Scenario: Forbid authorized request on GTFS-RT feed (no token)
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credential | secret |
      | remote_objectid_kind | internal |
    When I send this GTFS-RT request to the Referential "test" without token
    Then I should not receive a GTFS-RT but an unauthorized client error status

  Scenario: Forbid authorized request on GTFS-RT feed (wrong token)
    Given a Partner "test" exists with connectors [gtfs-rt-trip-updates-broadcaster] and the following settings:
      | local_credential | secret |
      | remote_objectid_kind | internal |
    When I send this GTFS-RT request to the Referential "test" with token "wrong"
    Then I should not receive a GTFS-RT but an unauthorized client error status
