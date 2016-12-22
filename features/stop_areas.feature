Feature: Manager StopAreas

	Background:
	  Given a Referential "test" is created

 @server
  Scenario: Create a StopArea
    When a StopArea is created with the following attributes :
    | Name      | Test               |
    | ObjectIds | "internal": "1234" |
    Then one StopArea has the following attributes:
    | Name      | Test               |
    | ObjectIds | "internal": "1234" |

  @server
  Scenario: Create two StopAreas
    When a StopArea is created with the following attributes :
    | Name      | Test 1             |
    | ObjectIds | "internal": "1234" |
        And a StopArea is created with the following attributes :
    | Name      | Test 2             |
    | ObjectIds | "internal": "2345" |
    Then one StopArea has the following attributes:
    | Name      | Test 1             |
    | ObjectIds | "internal": "1234" |
    And one StopArea has the following attributes:
    | Name      | Test 2             |
    | ObjectIds | "internal": "2345" |

  @server
  Scenario: Find StopArea by object id
    When a StopArea is created with the following attributes :
    | Name      | Test 1                                |
    | ObjectIds | "internal": "1234", "external": "abc" |
    Then a StopArea "internal":"1234" should exist
    And a StopArea "external":"abc" should exist

  @server
  Scenario: Destroy a StopArea
    Given a StopArea exists with the following attributes :
    | Name      | Test 1             |
    | ObjectIds | "internal": "1234" |
        And a StopArea exists with the following attributes :
    | Name      | Test 2             |
    | ObjectIds | "internal": "2345" |
        When the StopArea "internal":"1234" is destroy :
    Then a StopArea "internal":"1234" should not exist
    And one StopArea has the following attributes:
    | Name      | Test 2             |
    | ObjectIds | "internal": "2345" |

  @server
  Scenario: Create StopAreas in two Referentials
    Given a Referential 'test1' exists
        And a Referential 'test2' exists
        When a StopArea is created in Referential 'test1' with the following attributes :
    | Name      | Test 1             |
    | ObjectIds | "internal": "1234" |
        And a StopArea is created in Referential 'test2' with the following attributes :
    | Name      | Test 2             |
    | ObjectIds | "internal": "2345" |
    Then a StopArea "internal":"1234" should exist in Referential 'test1'
    And a StopArea "internal":"1234" should not exist in Referential 'test2'
    And a StopArea "internal":"2345" should exist in Referential 'test2'
    And a StopArea "internal":"2345" should not exist in Referential 'test1'