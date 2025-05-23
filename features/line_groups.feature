Feature: Manager LineGroups

  Background:
    Given a Referential "test" is created

  Scenario: Create a LineGroup
    When a LineGroup is created with the following attributes:
      | Name      | Test       |
      | ShortName | short_name |
      | LineIds   | ["1234"]   |
    Then the LineGroup "6ba7b814-9dad-11d1-1-00c04fd430c8" has the following attributes:
      | Name      | Test                              |
      | ShortName | short_name                        |
      | LineIds   | ["1234"]                          |
