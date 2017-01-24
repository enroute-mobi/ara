Feature: Manages Lines

  Background:
    Given a Referential "test" is created

  Scenario: Create a Line
  When a Line is created with the following attributes:
  | Name      | Test               |
  | ObjectIds | "internal": "1234" |
  Then one Line has the following attributes:
  | Name      | Test               |
  | ObjectIds | "internal": "1234" |

  Scenario: Create two Lines
  When a Line is created with the following attributes:
  | Name      | Test 1             |
  | ObjectIds | "internal": "1234" |
    And a Line is created with the following attributes:
  | Name      | Test 2             |
  | ObjectIds | "internal": "2345" |
  Then one Line has the following attributes:
  | Name      | Test 1             |
  | ObjectIds | "internal": "1234" |
    And one Line has the following attributes:
  | Name      | Test 2             |
  | ObjectIds | "internal": "2345" |

  Scenario: Find Line by object id
  When a Line is created with the following attributes:
  | Name      | Test 1                                |
  | ObjectIds | "internal": "1234", "external": "abc" |
  Then a Line "internal":"1234" should exist
    And a Line "external":"abc" should exist

  Scenario: Destroy a Line
  Given a Line exists with the following attributes:
  | Name      | Test 1             |
  | ObjectIds | "internal": "1234" |
    And a Line exists with the following attributes:
  | Name      | Test 2             |
  | ObjectIds | "internal": "2345" |
  When the Line "internal":"1234" is destroyed
  Then a Line "internal":"1234" should not exist
    And one Line has the following attributes:
  | Name      | Test 2             |
  | ObjectIds | "internal": "2345" |

  Scenario: Create Lines in two Referentials
  Given a Referential "test1" exists
    And a Referential "test2" exists
  When a Line is created in Referential "test1" with the following attributes:
  | Name      | Test 1             |
  | ObjectIds | "internal": "1234" |
    And a Line is created in Referential "test2" with the following attributes:
  | Name      | Test 2             |
  | ObjectIds | "internal": "2345" |
  Then a Line "internal":"1234" should exist in Referential "test1"
    And a Line "internal":"1234" should not exist in Referential "test2"
    And a Line "internal":"2345" should exist in Referential "test2"
    And a Line "internal":"2345" should not exist in Referential "test1"