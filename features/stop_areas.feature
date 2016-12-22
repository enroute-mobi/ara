Feature: Manager StopAreas

	Background:
	  Given a Referential "test" is created

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
