Feature: Manager StopAreaGroups

  Background:
    Given a Referential "test" is created

  Scenario: Create a StopAreaGroups
    When a StopAreaGroup is created with the following attributes:
      | Name        | Test       |
      | ShortName   | short_name |
      | StopAreaIds | ["1234"]   |
    Then the StopAreaGroup "6ba7b814-9dad-11d1-1-00c04fd430c8" has the following attributes:
      | Name        | Test       |
      | ShortName   | short_name |
      | StopAreaIds | ["1234"]   |
