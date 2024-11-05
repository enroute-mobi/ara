Feature: Manager Line Groups

  Background:
    Given a Referential "test" is created

  Scenario: Create a Line Groups
    When a Line Group is created with the following attributes:
      | Name      | Test       |
      | ShortName | short_name |
      | LineIds   | ["1234"]   |
    Then one Line Group has the following attributes:
      | Id        | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
      | Name      | Test                              |
      | ShortName | short_name                        |
      | LineIds   | ["1234"]                          |
