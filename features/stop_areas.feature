Feature: Manager StopAreas

	Background:
	  Given a Referential "test" is created

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