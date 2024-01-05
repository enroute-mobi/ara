Feature: Manager StopAreas

  Background:
    Given a Referential "test" is created

  Scenario: Create a StopArea
    When a StopArea is created with the following attributes:
      | Name  | Test               |
      | Codes | "internal": "1234" |
    Then one StopArea has the following attributes:
      | Name  | Test               |
      | Codes | "internal": "1234" |

  Scenario: Create two StopAreas
    When a StopArea is created with the following attributes:
      | Name  | Test 1             |
      | Codes | "internal": "1234" |
    And a StopArea is created with the following attributes:
      | Name  | Test 2             |
      | Codes | "internal": "2345" |
    Then one StopArea has the following attributes:
      | Name  | Test 1             |
      | Codes | "internal": "1234" |
    And one StopArea has the following attributes:
      | Name  | Test 2             |
      | Codes | "internal": "2345" |

  Scenario: Find StopArea by Code
    When a StopArea is created with the following attributes:
      | Name  | Test 1                                |
      | Codes | "internal": "1234", "external": "abc" |
    Then a StopArea "internal":"1234" should exist
    And a StopArea "external":"abc" should exist

  Scenario: Destroy a StopArea
    Given a StopArea exists with the following attributes:
      | Name  | Test 1             |
      | Codes | "internal": "1234" |
    And a StopArea exists with the following attributes:
      | Name  | Test 2             |
      | Codes | "internal": "2345" |
    When the StopArea "internal":"1234" is destroyed
    Then a StopArea "internal":"1234" should not exist
    And one StopArea has the following attributes:
      | Name  | Test 2             |
      | Codes | "internal": "2345" |

  Scenario: Create StopAreas in two Referentials
    Given a Referential "test1" exists
    And a Referential "test2" exists
    When a StopArea is created in Referential "test1" with the following attributes:
      | Name  | Test 1             |
      | Codes | "internal": "1234" |
    And a StopArea is created in Referential "test2" with the following attributes:
      | Name  | Test 2             |
      | Codes | "internal": "2345" |
    Then a StopArea "internal":"1234" should exist in Referential "test1"
    And a StopArea "internal":"1234" should not exist in Referential "test2"
    And a StopArea "internal":"2345" should exist in Referential "test2"
    And a StopArea "internal":"2345" should not exist in Referential "test1"
