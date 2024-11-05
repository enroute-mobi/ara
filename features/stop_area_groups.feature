Feature: Manager StopArea Groups

  Background:
    Given a Referential "test" is created

  Scenario: Create a StopArea Groups
    When a StopArea Group is created with the following attributes:
      | Name      | Test       |
      | ShortName | short_name |
      | StopAreaIds | ["1234"]   |
    Then one StopArea Group has the following attributes:
      | Id          | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
      | Name        | Test                              |
      | ShortName   | short_name                        |
      | StopAreaIds | ["1234"]                          |
