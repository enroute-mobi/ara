Feature: Manager Referentials

  Scenario: Create a Referential
    When a Referential "test" is created
    Then a Referential "test" should exist

  Scenario: Destroy a Referential
    Given a Referential "test" exists
    When the Referential "test" is destroyed
    Then a Referential "test" should not exist

  Scenario: 2698 - Referential reloads Model at configured time
    Given a Referential "test" exists with the following settings:
        | model.reload_at | 01:00 |
    And a StopArea exists with the following attributes:
        | Name      | Test 1                                                                    |
        | ObjectIDs | "internal": "boaarle", "external": "RATPDev:StopPoint:Q:eeft52df543d:LOC" |
    And a Line exists with the following attributes:
        | Name      | Ligne 415                       |
        | ObjectIDs | "internal": "CdF:Line::415:LOC" |
    And a VehicleJourney exists with the following attributes:
        | ObjectIDs | "internal": "1STD721689197098"  |
    And a StopVisit exists with the following attributes:
        | ObjectIDs | "internal": "SIRI:34852540"   |
    When the time is "2017-01-02T05:00:00+01:00"
    Then a StopArea "internal":"boaarle" should exist
    And a Line "internal":"CdF:Line::415:LOC" should exist
    And a VehicleJourney "internal":"1STD721689197098" should not exist
    And a StopVisit "internal":"SIRI:34852540" should not exist
