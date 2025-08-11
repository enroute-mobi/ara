Feature: Manager Facilities

  Background:
    Given a Referential "test" is created

  @ARA-1730
  Scenario: Create a Facility
    When a Facility is created with the following attributes:
      | Codes[internal] |      1234 |
      | Status          | available |
    Then one Facility has the following attributes:
      | Codes[internal] |      1234 |
      | Status          | available |

   @ARA-1749
  Scenario: Create a Facility with missing Status
    When a Facility is created with the following attributes:
      | Codes[internal] | 1234 |
    Then one Facility has the following attributes:
      | Codes[internal] |    1234 |
      | Status          | unknown |
