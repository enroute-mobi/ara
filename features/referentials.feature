Feature: Manager Referentials

  Scenario: Create a Referential
    When a Referential "test" is created
    Then a Referential "test" should exist

  Scenario: Destroy a Referential
    Given a Referential "test" exists
    When the Referential "test" is destroyed
    Then a Referential "test" should not exist
