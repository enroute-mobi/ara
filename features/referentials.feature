Feature: Manager Referentials

  Scenario: Create a Referential
    When a Referential "test" is created
    Then a Referential "test" should exist

  Scenario: Destroy a Referential
    Given a Referential "test" exists
    When the Referential "test" is destroyed
    Then a Referential "test" should not exist

  @nostart @database
  Scenario: 2698 - Referential reloads Model at configured time
    Given the table "stop_areas" has the following data:
    | id                                     | referential_slug     | codes                    | model_name |
    | '6ba7b814-9dad-11d1-0011-00c04fd430c8' | 'test'               | '{"internal":"value"}'   |'2017-01-02'|
    When I start Ara
    Given a Referential "test" exists with the following settings:
        | model.reload_at | 01:00 |
    And a StopArea exists with the following attributes:
        | Name      | Test 1                                                                    |
        | Codes | "internal": "boaarle", "external": "RATPDev:StopPoint:Q:eeft52df543d:LOC" |
    And a Line exists with the following attributes:
        | Name      | Ligne 415                       |
        | Codes | "internal": "CdF:Line::415:LOC" |
    And a VehicleJourney exists with the following attributes:
        | Codes | "internal": "1STD721689197098"  |
    And a StopVisit exists with the following attributes:
        | Codes | "internal": "SIRI:34852540"   |
    When the time is "2017-01-02T05:00:00+01:00"
    And a VehicleJourney "internal":"1STD721689197098" should not exist
    And a StopVisit "internal":"SIRI:34852540" should not exist
    And a StopArea "internal":"value" should exist
