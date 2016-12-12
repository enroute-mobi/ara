Feature: Manager Referentials

  @server
  Scenario: Create a Referential
    When a Referential 'test' is created
    Then one Referential 'test' should exist

  @server
  Scenario: Destroy a Referential
    Given a Referential 'test' exists
    When the Referential 'test' is destroyed
    Then a Referential 'test' should not exists
